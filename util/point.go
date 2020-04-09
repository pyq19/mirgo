package util

import "fmt"

type Point struct {
	X uint32
	Y uint32
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (p Point) NextPoint(direction MirDirection, step uint32) Point {
	x := p.X
	y := p.Y
	switch direction {
	case MirDirectionUp:
		y = y - step
	case MirDirectionUpRight:
		x = x + step
		y = y - step
	case MirDirectionRight:
		x = x + step
	case MirDirectionDownRight:
		x = x + step
		y = y + step
	case MirDirectionDown:
		y = y + step
	case MirDirectionDownLeft:
		x = x - step
		y = y + step
	case MirDirectionLeft:
		x = x - step
	case MirDirectionUpLeft:
		x = x - step
		y = y - step
	}
	return Point{X: x, Y: y}
}

func NewPoint(x, y int) Point {
	return Point{uint32(x), uint32(y)}
}

func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) EqualXY(x, y int) bool {
	return p.X == uint32(x) && p.Y == uint32(y)
}

func InRange(a, b Point, i int) bool {
	return AbsInt(int(a.X)-int(b.X)) <= i && AbsInt(int(a.Y)-int(b.Y)) <= i
}

func InRangeXY(a Point, x, y, i int) bool {
	return AbsInt(int(a.X)-x) <= i && AbsInt(int(a.Y)-y) <= i
}
