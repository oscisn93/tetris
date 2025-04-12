package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oscisn93/tetris/tetris"
)

func main() {
	ebiten.SetWindowSize(tetris.ScreenWidth, tetris.ScreenHeight)
	ebiten.SetWindowTitle("Tetris")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&tetris.Game{}); err != nil {
		log.Fatal(err)
	}
}
