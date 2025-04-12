package tetris

import (
	"github.com/hajimehoshi/ebiten/v2"
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

// Vim Keybindings
const (
	MoveRightKey = ebiten.KeyL
	MoveLeftKey = ebiten.KeyH
	MoveDownKey = ebiten.KeyJ
	RotateKey = ebiten.KeyK
	SlamKey = ebiten.KeySpace
	PlayGameKey = ebiten.KeyI
  SuspendGameKey = ebiten.KeyEscape
  ExitKey = ebiten.KeyQ
)

type Block int

const (
  
)

type Board struct {
  cells [BoardCellsX][BoardCellsY]Block
}

type Game struct {

}

func (game *Game) Layout(width, height int) (screenWidth, screenHeight int) {
  return ScreenWidth, ScreenHeight 
}

func (game *Game) Update() error {
  return nil
}

func (game *Game) Draw(screen *ebiten.Image) {

}
