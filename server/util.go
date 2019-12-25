package main

import (
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
	m := new(Map)
	m.CoordinateCellMap = new(sync.Map)
	offset := 21
	w := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	xor := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	h := common.BytesToUint16(bytes[offset : offset+2])
	width := w ^ xor
	height := h ^ xor
	offset = 54
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			p := common.Point{X: uint32(i), Y: uint32(j)}
			c := new(Cell)
			c.Map = m
			c.Coordinate = p.Coordinate()
			if (common.BytesToUint32(bytes[offset:offset+4])^0xAA38AA38)&0x20000000 != 0 {
				c.Attribute = common.CellAttributeHighWall
			}
			if ((common.BytesToUint16(bytes[offset+6:offset+8]) ^ xor) & 0x8000) != 0 {
				c.Attribute = common.CellAttributeLowWall
			}
			m.CoordinateCellMap.Store(p.Coordinate(), c)
			offset += 15
		}
	}
	m.Width = width
	m.Height = height
	return m
}
