package tetris

import (
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Move int

const (
  PushLeft Move = iota
  PushRight
  PushDown
  RotateLeft
  RotateRight
)

type  int

const (
  h keyState = iota
  l
  j
  leftArrow
  rightArrow
  downArrow

)
