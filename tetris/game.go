package tetris

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	assets "github.com/oscisn93/tetris/assets"
)

// Screen & Board Dimensions
const (
	ScreenWidth  = 480
	ScreenHeight = 800
	BoardWidth   = 320
	BoardHeight  = 640
	BoardCellsX  = 10
	BoardCellsY  = 20
)

func createSplashScreen(screen *ebiten.Image) {
  img, _, err := image.Decode(bytes.NewReader(assets.EbitenginePng))
  if err != nil {
    panic(err)
  }
  splashImage := ebiten.NewImageFromImage(img)
  scale := splashImage.Bounds().Dx() / ScreenWidth

  x := ScreenWidth / 2 - (splashImage.Bounds().Dx() / 2 * scale)
  y := ScreenHeight / 2 - (splashImage.Bounds().Dy() / 2 * scale)

  options := &ebiten.DrawImageOptions{}
  options.GeoM.Scale(float64(1 / scale), float64(1 / scale))
  options.GeoM.Translate(float64(x), float64(y))

  screen.DrawImage(splashImage, options)
}

type Game struct {
	CurrentScene Scene
	State        *GameState
}

func (game *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (game *Game) Update() error {
  if game.CurrentScene != GameScene {
    return nil
  }

  err := game.State.Update(game)

  return err
}

func (game *Game) Draw(screen *ebiten.Image) {
  if game.CurrentScene == SplashScene {
    createSplashScreen(screen)
  }
}
