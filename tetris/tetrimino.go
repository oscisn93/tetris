package tetris

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	assets "github.com/oscisn93/tetris/assets"
)

type Shape int

const (
	ShapeO Shape = iota
	ShapeI
	ShapeJ
	ShapeL
	ShapeZ
	ShapeT
	ShapeS
)

var blocksImage *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(assets.BlocksPng))
	if err != nil {
		panic(err)
	}
	blocksImage = ebiten.NewImageFromImage(img)
}

type Theta int

const (
	TwoPI Theta = iota
	HalfPI
	OnePI
	OneHalfPI
)

func (theta Theta) Rotate() Theta {
	if theta == OneHalfPI {
		return TwoPI
	}
	return theta + 1
}

type Block int

const (
	RedBlock Block = iota
	YellowBlock
	GreenBlock
	BlueBlock
	PurpleBlock
	PinkBlock
	GrayBlock
	EmptyBlock
)

type Tetrimino struct {
	BlockType Block
	Body      [][]bool
}

var Tetriminos map[Block]*Tetrimino

func init() {
	Tetriminos = map[Block]*Tetrimino{
		RedBlock: {
			BlockType: RedBlock,
			Body: [][]bool{
        {false, false, false, false },
        {true,  true,  true,  true  },
        {false, false, false, false },
        {false, false, false, false },
			},
		},
		YellowBlock: {
			BlockType: YellowBlock,
			Body: [][]bool{
        { false, true,  false },
        { true,  true,  true  },
        { false, false, false },
			},
		},
		GreenBlock: {
			BlockType: GreenBlock,
			Body: [][]bool{
        { true,  false, false },
        { true,  true,  true  },
        { false, false, false },
			},
		},
		BlueBlock: {
			BlockType: BlueBlock,
			Body: [][]bool{
        { false, false, true  },
        { true,  true,  true  },
        { false, false, false },
			},
		},
		PurpleBlock: {
			BlockType: PurpleBlock,
			Body: [][]bool{
        { true,  true,  false },
        { false, true,  true  },
        { false, false, false },
			},
		},
		PinkBlock: {
			BlockType: PinkBlock,
			Body: [][]bool{
        { false, true,  true  },
        { true,  true,  false },
        { false, false, false },
			},
		},
		GrayBlock: {
			BlockType: GrayBlock,
			Body: [][]bool{
        { true, true },
        { true, true },
			},
		},
	}
}

const (
	BlockWidth  = 32
	BlockHeight = 32
)

func drawBlock(rect *ebiten.Image, block Block, x, y int, clr colorm.ColorM) {
	options := &colorm.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))
	sourceX := int(block) * BlockWidth
	blockImage := blocksImage.SubImage(image.Rect(sourceX, 0, sourceX+BlockWidth, BlockHeight)).(*ebiten.Image)
	colorm.DrawImage(rect, blockImage, clr, options)
}

func (tetrimino *Tetrimino) InitialPosition() (int, int) {
	size := len(tetrimino.Body)
	x := (BoardCellsX - size) / 2
	y := 0
Loop:
	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			if tetrimino.Body[i][j] {
				break Loop
			}
		}
		y--
	}
	return x, y
}

func (tetrimino *Tetrimino) isBlocked(i, j int, theta Theta) bool {
	size := len(tetrimino.Body)
	i2, j2 := i, j
	switch theta {
	case TwoPI:
	case HalfPI:
		i2 = j
		j2 = size - 1 - i
	case OnePI:
		i2 = size - 1 - i
		j2 = size - 1 - j
	case OneHalfPI:
		i2 = size - 1 - j
		j2 = i
	}
	return tetrimino.Body[i2][j2]
}

func (tetrimino *Tetrimino) collides(board *Board, x, y int, theta Theta) bool {
	size := len(tetrimino.Body)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board.IsBlocked(x+i, y+j) && tetrimino.isBlocked(i, j, theta) {
				return true
			}
		}
	}
	return false
}

func (tetrimino *Tetrimino) AbsorbInto(board *Board, x, y int, theta Theta) {
	size := len(tetrimino.Body)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if tetrimino.isBlocked(i, j, theta) {
				board.setBlock(x+i, y+j, tetrimino.BlockType)
			}
		}
	}
}

func (tetrimino *Tetrimino) DrawAtCenter(rect *ebiten.Image, x, y, width, height int, theta Theta) {
	x += (width - len(tetrimino.Body[0])*BlockWidth) / 2
	y += (height - len(tetrimino.Body)*BlockHeight) / 2
	tetrimino.Draw(rect, x, y, theta)
}

func (tetrimino *Tetrimino) Draw(rect *ebiten.Image, x, y int, theta Theta) {
	for i := range tetrimino.Body {
		for j := range tetrimino.Body[i] {
			if tetrimino.isBlocked(i, j, theta) {
				drawBlock(rect, tetrimino.BlockType, i*BlockWidth+x, i*BlockHeight+y, colorm.ColorM{})
			}
		}
	}
}
