package tetris

import "github.com/hajimehoshi/ebiten/v2"

const (
  ScreenWidth  = 1476
  ScreenHeight = 716
)

type Game struct {
  // input        *Input
  // board        *Board
  // boardImage   *ebiten.Image
}

// func NewGame() (*Game, error) {
//   game := &Game{
//     input: NewInput(),
//   }

//   var err error

//   game.board, err = NewBoard()
//   if err != nil {
//     return nil, err
//   }
//   return game, nil
// }

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  return outsideWidth, outsideHeight
}

func (game *Game) Update() error {
  // game.input.Update()
  // if err := game.board.Update(game.input); err != nil {
  //   return err
  // }
  return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
  // if game.boardImage == nil {
    // game.boardImage = ebiten.NewImage(game.board.Size())
  // }
  // screen.Fill(backgroundColor)
  // game.board.Draw(game.boardImage)
  options := &ebiten.DrawImageOptions{}
  screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
  // boardWidth, boardHeight := game.boardImage.Bounds().Dx(), game.boardImage.Bounds().Dy()
  // x := (screenWidth - boardWidth) / 2
  x := screenWidth / 2
  // y := (screenHeight - boardHeight) / 2
  y := screenHeight / 2
  options.GeoM.Translate(float64(x), float64(y))
  // screen.DrawImage(game.boardImage, options)
}

