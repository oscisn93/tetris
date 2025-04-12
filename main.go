package main

import (
  "log"

  "github.com/hajimehoshi/ebiten/v2"
  "github.com/oscisn93/tetris/tetris"
)

func main() {
  game := tetris.Game{}

  ebiten.SetWindowSize(tetris.ScreenWidth, tetris.ScreenHeight)
  ebiten.SetWindowTitle("Tetris")

  if err := ebiten.RunGame(game); err != nil {
    log.Fatal(err)
  }
}

