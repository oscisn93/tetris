package tetris

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

type Tetrimino struct {
	Tshape Shape
	X, Y   float64
}
