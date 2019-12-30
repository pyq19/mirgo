package common

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X uint32
	Y uint32
}

func (p *Point) String() string {
	return p.Coordinate()
}

// Coordinate 点的坐标
func (p Point) Coordinate() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

// NewPointByCoordinate 坐标转换成点
func NewPointByCoordinate(coordinate string) *Point {
	strArr := strings.Split(coordinate, ",")
	x, _ := strconv.Atoi(strArr[0])
	y, _ := strconv.Atoi(strArr[1])
	return &Point{X: uint32(x), Y: uint32(y)}
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

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (c Color) ToInt() int {
	return 4294967295 // 255, 255, 255, 255
}

func (c Color) ToInt32() int32 {
	return int32(c.ToUint32())
}

func (c Color) ToUint32() uint32 {
	return 4294967295
}
