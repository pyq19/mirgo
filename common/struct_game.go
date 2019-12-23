package common

import (
	"fmt"
)

type Point struct {
	X uint32
	Y uint32
}

// Coordinate 点的坐标
func (p Point) Coordinate() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func NewPoint(x, y int) *Point {
	return &Point{uint32(x), uint32(y)}
}

type SelectInfo struct {
	Index      uint32
	Name       string
	Level      uint16
	Class      MirClass
	Gender     MirGender
	LastAccess int64
}
