package common

import (
	"fmt"
	"io/ioutil"
)

type Point struct {
	X uint32
	Y uint32
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

type SelectInfo struct {
	Index      uint32
	Name       string
	Level      uint16
	Class      MirClass
	Gender     MirGender
	LastAccess int64
}

type Cell struct {
	Coordinate string // 坐标 x,y
	Point      Point  // 坐标点
	Attribute  CellAttribute
}

type Door struct{}

type Map struct {
	Width  uint16
	Height uint16
	Cells  []Cell
}

func GetMapBytes(mapAbsPath string) []byte {
	fileBytes, err := ioutil.ReadFile(mapAbsPath)
	if err != nil {
		panic(err)
	}
	return fileBytes
}

func GetMapV1(bytes []byte) *Map {
	offset := 21
	w := BytesToUint16(bytes[offset : offset+2])
	offset += 2
	xor := BytesToUint16(bytes[offset : offset+2])
	offset += 2
	h := BytesToUint16(bytes[offset : offset+2])
	width := w ^ xor
	height := h ^ xor
	offset = 54
	cellCount := int(width) * int(height)
	cells := make([]Cell, 0, cellCount)
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			p := Point{X: uint32(i), Y: uint32(j)}
			c := Cell{Coordinate: p.String(), Point: p}
			if (BytesToUint32(bytes[offset:offset+4])^0xAA38AA38)&0x20000000 != 0 {
				c.Attribute = CellAttributeHighWall
				cells = append(cells, c)
			}
			if ((BytesToUint16(bytes[offset+6:offset+8]) ^ xor) & 0x8000) != 0 {
				c.Attribute = CellAttributeLowWall
				cells = append(cells, c)
			}
			offset += 15
		}
	}
	m := new(Map)
	m.Width = width
	m.Height = height
	m.Cells = cells
	return m
}
