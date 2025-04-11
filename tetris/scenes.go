package tetris

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math/rand/v2"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	assets "github.com/oscisn93/tetris/assets"
)


type Scene int

const (
	// the splash screen
	SplashScene Scene = iota
	// the main game screen
	GameScene
	GameOver
)

var imageGameBackground *ebiten.Image
var imageWindows = ebiten.NewImage(ScreenWidth, ScreenHeight)
var imageGameOver = ebiten.NewImage(ScreenWidth, ScreenHeight)

func boardWindowPosition() (x, y int) {
	return 20, 20
}

func queuedWindowLabelPosition() (x, y int) {
	x, y = boardWindowPosition()
	return x + BoardWidth + 2*BlockWidth, y
}

func queuedWindowPosition() (x, y int) {
	x, y = queuedWindowLabelPosition()
	return x, y + BlockHeight
}

func textBoxWidth() int {
	x, _ := queuedWindowPosition()
	return ScreenWidth - 2*BlockWidth - x
}

func scoreTextBoxPosition() (x, y int) {
	x, y = queuedWindowPosition()
	return x, y + 6*BlockHeight
}

func levelTextBoxPosition() (x, y int) {
	x, y = scoreTextBoxPosition()
	return x, y + 4*BlockHeight
}

func linesTextBoxPosition() (x, y int) {
	x, y = levelTextBoxPosition()
	return x, y + 4*BlockHeight
}

var fontColor = color.RGBA{0x40, 0x40, 0xff, 0xff}
var shadowColor = color.RGBA{0, 0, 0, 0x80}

const arcadeFontBaseSize = 8

var arcadeFaceSource *text.GoTextFaceSource

func initFont() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(assets.PerssStart2PTtf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = src
}

func drawWindow(rect *ebiten.Image, x, y, width, height int) {
	vector.DrawFilledRect(rect, float32(x), float32(y), float32(width), float32(height), color.RGBA{0, 0, 0, 0xc0}, false)
}

func drawTextWithShadow(rect *ebiten.Image, str string, x, y, scale int, clr color.Color, primaryAlign, secondaryAlign text.Align) {
	options := &text.DrawOptions{}

	options.GeoM.Translate(float64(x)+1, float64(y)+1)
	options.ColorScale.ScaleWithColor(shadowColor)
	options.LineSpacing = arcadeFontBaseSize * float64(scale)
	options.PrimaryAlign = primaryAlign
	options.SecondaryAlign = secondaryAlign

	text.Draw(rect, str, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   arcadeFontBaseSize * float64(scale),
	}, options)

	options.GeoM.Reset()
	options.GeoM.Translate(float64(x), float64(y))

	options.ColorScale.Reset()
	options.ColorScale.ScaleWithColor(clr)

	text.Draw(rect, str, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   arcadeFontBaseSize * float64(scale),
	}, options)
}

func drawTextBox(rect *ebiten.Image, label string, x, y, width int) {
	drawTextWithShadow(rect, label, x, y, 1, fontColor, text.AlignStart, text.AlignStart)
	y += BlockHeight
	drawWindow(rect, x, y, width, 2*BlockHeight)
}

func drawTextBoxContent(rect *ebiten.Image, content string, x, y, width int) {
	y += BlockHeight
	drawTextWithShadow(rect, content, x+width-2*BlockHeight/4, y+2*BlockHeight/2, 1, color.White, text.AlignEnd, text.AlignEnd)
}

func initGameWindow() {
	initFont()

	img, _, err := image.Decode(bytes.NewReader(assets.LogoPng))
	if err != nil {
		panic(err)
	}
	imageGameBackground = ebiten.NewImageFromImage(img)
}

var lightGray colorm.ColorM

var maxQueuedTetrinimos = 3

type CurrentTetrimino struct {
	terimino    *Tetrimino
	x, y, carry int
	theta       Theta
}

type GameState struct {
	board    *Board
	current  *CurrentTetrimino
	held     *Tetrimino
	queued   []*Tetrimino
	landing  int
	score    int
	lines    int
	gameover bool
}

func NewGameState() *GameState {
	return &GameState{
		board: &Board{},
	}
}

func init() {
	var id colorm.ColorM
	var mono colorm.ColorM
	mono.ChangeHSV(0, 0, 1)

	for j := 0; j < colorm.Dim-1; j++ {
		for i := 0; i < colorm.Dim-1; i++ {
			lightGray.SetElement(1, j, mono.Element(i, j)*0.7+id.Element(i, j)*0.3)
		}
	}

	lightGray.Translate(0.3, 0.3, 0.3, 0)
}

func (state *GameState) drawBackground(rect *ebiten.Image) {
	rect.Fill(color.White)

	width, height := imageGameBackground.Bounds().Dx(), imageGameBackground.Bounds().Dy()

	scaleWidth := ScreenWidth / float64(width)
	scaleHeight := ScreenHeight / float64(height)

	scale := scaleWidth

	if scale < scaleHeight {
		scale = scaleHeight
	}

	options := &colorm.DrawImageOptions{}

	options.GeoM.Translate(-float64(width)/2, -float64(height)/2)
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	options.Filter = ebiten.FilterLinear

	colorm.DrawImage(rect, imageGameBackground, lightGray, options)
}

func (state *GameState) chooseTetrimino() *Tetrimino {
	max := int(GrayBlock)
	blockType := Block(rand.IntN(max))
	return Tetriminos[blockType]
}

func (state *GameState) initCurrentTerinimo(tetrinimo *Tetrimino) {
	state.current.terimino = tetrinimo

	x, y := state.current.terimino.InitialPosition()

	state.current.x = x
	state.current.y = y
	state.current.carry = 0
	state.current.theta = TwoPI
}

func (state *GameState) level() int {
	return state.lines / 10
}

func (state *GameState) addScore(lines int) {
	base := 0

	switch lines {
	case 1:
		base = 100
	case 2:
		base = 300
	case 3:
		base = 600
	case 4:
		base = 1000
	default:
		panic("NOT_REACHED")
	}

	state.score += (state.level() + 1) * base
}

// Vim Keybindings
const (
	MoveRightKey   = ebiten.KeyL
	MoveLeftKey    = ebiten.KeyH
	MoveDownKey    = ebiten.KeyJ
	RotateKey      = ebiten.KeyK
	SlamKey        = ebiten.KeySpace
	PlayGameKey    = ebiten.KeyI
	SuspendGameKey = ebiten.KeyEscape
	ExitKey        = ebiten.KeyQ
)

func (state *GameState) Update(game *Game) error {
	state.board.Update()

	if state.gameover {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			game.CurrentScene = GameOver 
			game.State = nil
		}
		return nil
	}

	maxLandingCount := ebiten.TPS()

	if state.current.terimino == nil {
		state.initCurrentTerinimo(state.chooseTetrimino())
	}

	if state.queued == nil {
		state.queued = []*Tetrimino{}
		for i := 0; i < maxQueuedTetrinimos; i++ {
			state.queued = append(state.queued, state.chooseTetrimino())
		}
	}

	moved := false
	tetrimino := state.current.terimino
	theta := state.current.theta

	if !state.board.IsFlushAnimating() {
		tetrimino := state.current.terimino

		x := state.current.x
		y := state.current.y

		if inpututil.IsKeyJustPressed(RotateKey) {
			state.current.theta = state.board.Rotate(tetrimino, x, y, theta)
			moved = theta != state.current.theta
		} else if dxl := inpututil.KeyPressDuration(MoveLeftKey); dxl > 0 && dxl == 1 || (10 <= dxl && dxl%2 == 0) {
			state.current.x = state.board.MoveToLeft(tetrimino, x, y, theta)
			moved = x != state.current.x
		} else if dxr := inpututil.KeyPressDuration(MoveRightKey); dxr > 0 && dxr == 1 || (10 <= dxr && dxr%2 == 0) {
			state.current.x = state.board.MoveToRight(tetrimino, x, y, theta)
			moved = x != state.current.x
		} else if dy := inpututil.KeyPressDuration(MoveDownKey); (dy-1)%2 == 0 {
			state.current.y = state.board.DropTetrimino(tetrimino, x, y, theta)
			moved = y != state.current.y
			if moved {
				state.score++
			}
		}
	}

	if !state.board.IsFlushAnimating() {
		theta := state.current.theta

		state.current.carry += 2*state.level() + 1

		const maxCarry = 60

		for maxCarry <= state.current.carry {
			state.current.carry -= maxCarry
			state.current.y = state.board.DropTetrimino(tetrimino, state.current.x, state.current.y, theta)
		}
	}

	if !state.board.IsFlushAnimating() && !state.board.IsTetriminoDroppable(tetrimino, state.current.x, state.current.y, theta) {
		if 0 < inpututil.KeyPressDuration(MoveDownKey) {
			state.landing += 10
		} else {
			state.landing -= 10
		}

		if maxLandingCount <= state.landing {
			state.board.AbsorbTetrimino(tetrimino, state.current.x, state.current.y, theta)
			if state.board.IsFlushAnimating() {
				state.board.SetEndFlushAnimating(func(lines int) {
					state.lines += lines
					if 0 < lines {
						state.addScore(lines)
					}
					state.getNextTetrimino()
				})
			} else {
				state.getNextTetrimino()
			}
		}
	}
	return nil
}

func (state *GameState) getNextTetrimino() {
	state.initCurrentTerinimo(state.queued[0])
	state.queued = state.queued[1:]
	state.queued = append(state.queued, state.chooseTetrimino())
	state.landing = 0

	if state.current.terimino.collides(state.board, state.current.x, state.current.y, state.current.theta) {
		state.gameover = true
	}
}

func (state *GameState) Draw(rect *ebiten.Image) {
	state.drawBackground(rect)

	rect.DrawImage(imageWindows, nil)

	x, y := scoreTextBoxPosition()
	drawTextBoxContent(rect, strconv.Itoa(state.score), x, y, textBoxWidth())

	x, y = levelTextBoxPosition()
	drawTextBoxContent(rect, strconv.Itoa(state.level()), x, y, textBoxWidth())

	x, y = linesTextBoxPosition()
	drawTextBoxContent(rect, strconv.Itoa(state.lines), x, y, textBoxWidth())

	boardX, boardY := boardWindowPosition()
	state.board.Draw(rect, boardX, boardY)

	if state.current.terimino != nil && !state.board.IsFlushAnimating() {
		x := boardX + state.current.x*BlockWidth
		y := boardY + state.current.y*BlockHeight
		state.current.terimino.Draw(rect, x, y, state.current.theta)
	}

	if state.queued != nil && len(state.queued) > 0 {
		x := boardX + BoardWidth + BlockWidth*2
		y := boardY + BoardHeight
		state.queued[0].DrawAtCenter(rect, x, y, BlockWidth*5, BlockHeight*5, 0)
	}

	if state.gameover {
		rect.DrawImage(imageGameOver, nil)
	}
}
