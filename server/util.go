package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"io/ioutil"
	"sync"
)

func GetMapBytes(mapAbsPath string) []byte {
	fileBytes, err := ioutil.ReadFile(mapAbsPath)
	if err != nil {
		panic(err)
	}
	return fileBytes
}

func GetMapV1(bytes []byte) *Map {
	offset := 21
	w := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	xor := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	h := common.BytesToUint16(bytes[offset : offset+2])
	width := int(w ^ xor)
	height := int(h ^ xor)

	// m := new(Map)
	// m.cells = make([]*Cell, width*height) //new(sync.Map)
	// aoi := newAOI(m, width, height)
	// m.AOI = aoi
	m := NewMap(width, height)
	fmt.Println(width, height)

	offset = 54
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			p := common.Point{X: uint32(i), Y: uint32(j)}
			c := new(Cell)
			c.Map = m
			c.Point = p
			c.Objects = new(sync.Map)
			if (common.BytesToUint32(bytes[offset:offset+4])^0xAA38AA38)&0x20000000 != 0 {
				c.Attribute = common.CellAttributeHighWall
			}
			if ((common.BytesToUint16(bytes[offset+6:offset+8]) ^ xor) & 0x8000) != 0 {
				c.Attribute = common.CellAttributeLowWall
			}
			if c.Attribute == common.CellAttributeWalk {
				m.SetCell(p, c)
			}
			offset += 15
		}
	}
	m.Width = width
	m.Height = height
	return m
}

func newAOI(m *Map, width int, height int) (aoi *AOIManager) {
	cntX := width / 20
	cntY := height / 20
	if width < 20 {
		cntX = 1
	}
	if height < 20 {
		cntY = 1
	}
	aoi = NewAOIManager(m, 0, width, cntX, 0, height, cntY)
	return
}
