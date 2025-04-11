package tetris

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const maxFlushCount = 20

type Board struct {
	Cells               [BoardCellsX][BoardCellsY]Block
	flushCount          int
	OnEndFlushAnimating func(int)
}

func (board *Board) IsBlocked(x, y int) bool {
	if x < 0 || BoardCellsX <= x {
		return true
	}
	if y < 0 {
		return false
	}
	return board.Cells[x][y] != EmptyBlock
}

func (board *Board) MoveToLeft(tetrimino *Tetrimino, x, y int, theta Theta) int {
	if tetrimino.collides(board, x-1, y, theta) {
		return x
	}
	return x - 1
}

func (board *Board) MoveToRight(tetrimino *Tetrimino, x, y int, theta Theta) int {
	if tetrimino.collides(board, x-1, y, theta) {
		return x
	}
	return x + 1
}

func (board *Board) IsTetriminoDroppable(tetrimino *Tetrimino, x, y int, theta Theta) bool {
	return !tetrimino.collides(board, x, y+1, theta)
}

func (board *Board) DropTetrimino(tetrimino *Tetrimino, x, y int, theta Theta) int {
	if tetrimino.collides(board, x, y+1, theta) {
		return y
	}
	return y + 1
}

func (board *Board) Rotate(tetrimino *Tetrimino, x, y int, theta Theta) Theta {
	rotated := theta.Rotate()
	if tetrimino.collides(board, x, y, rotated) {
		return theta
	}
	return rotated
}

func (board *Board) AbsorbTetrimino(tetrimino *Tetrimino, x, y int, theta Theta) {
	tetrimino.AbsorbInto(board, x, y, theta)
	if board.flushable() {
		board.flushCount = maxFlushCount
	}
}

func (board *Board) IsFlushAnimating() bool {
	return 0 < board.flushCount
}

func (board *Board) SetEndFlushAnimating(fn func(line int)) {
	board.OnEndFlushAnimating = fn
}

func (board *Board) flushable() bool {
	for j := BoardCellsY - 1; 0 < j; j-- {
		if board.flushableLine(j) {
			return true
		}
	}
	return false
}

func (board *Board) flushableLine(y int) bool {
	for i := 0; i < BoardCellsX; i++ {
		if board.Cells[i][y] == EmptyBlock {
			return false
		}
	}
	return true
}

func (board *Board) setBlock(x, y int, blockType Block) {
	board.Cells[x][y] = blockType
}

func (board *Board) endFlushAnimating() int {
	flushedLines := 0
	for j := BoardCellsY - 1; 0 <= j; j-- {
		if board.flushLine(j + flushedLines) {
			flushedLines++
		}
	}
	return flushedLines
}

func (board *Board) flushLine(y int) bool {
	for i := 0; i < BoardCellsX; i++ {
		if board.Cells[i][y] == EmptyBlock {
			return false
		}
	}

	for yPrime := y; 1 <= yPrime; yPrime-- {
		for i := 0; i < BoardCellsX; i++ {
			board.Cells[i][yPrime] = board.Cells[i][yPrime-1]
		}
	}

	for i := 0; i < BoardCellsX; i++ {
		board.Cells[i][0] = EmptyBlock
	}
	return true
}

func (board *Board) Update() {
	if board.flushCount == 0 {
		return
	}
	board.flushCount--
	if board.flushCount > 0 {
		return
	}
	if board.OnEndFlushAnimating != nil {
		board.OnEndFlushAnimating(board.endFlushAnimating())
	}
}

func flushingColor(rate float64) colorm.ColorM {
	var clr colorm.ColorM
	alpha := min(1, rate*2)
	clr.Scale(1, 1, 1, alpha)
	red := min(1, (1-rate)*2)
	clr.Translate(red, 0, 0, 0)
	return clr
}

func (board *Board) Draw(rect *ebiten.Image, x, y int) {
	flushColor := flushingColor(float64(board.flushCount) / maxFlushCount)
	for j := 0; j < BoardCellsY; j++ {
		if board.flushableLine(j) {
			for i := 0; i < BoardCellsX; i++ {
				drawBlock(rect, board.Cells[i][j], i*BlockWidth+x, j*BlockHeight+y, flushColor)
			}
		} else {
			for i := 0; i < BoardCellsX; i++ {
				drawBlock(rect, board.Cells[i][j], i*BlockWidth+x, j*BlockHeight+y, colorm.ColorM{})
			}
		}
	}
}
