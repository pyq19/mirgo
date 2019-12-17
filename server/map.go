package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"io/ioutil"
	"sync"
)

// Map ...
type Map struct {
	Width             uint16
	Height            uint16
	Cells             []Cell
	WalkableCells     []Cell
	CoordinateCellMap *sync.Map // map[string]*Cell
	NPCs              []NPC
	Players           []Player
	Respawns          []Respawn
}

type MapObj interface {
}

type NPC struct {
	MapObj
}

const (
	LOGIN = iota
	SELECT
	GAME
	DISCONNECTED
)

// Player ...
type Player struct {
	MapObj
	AccountId int
	GameStage int
	Session   *cellnet.Session
	Character *common.Character
	Magics    *[]common.MagicInfo
	UserItems *[]common.UserItem
}

type Respawn struct {
	MapObj
}

type Cell struct {
	Coordinate string       // 坐标 x,y
	Point      common.Point // 坐标点
	Attribute  common.CellAttribute
	Objects    []MapObj
}

type Door struct{}

func (m *Map) GetCell(coordinate string) *Cell {
	if v, ok := m.CoordinateCellMap.Load(coordinate); ok {
		return v.(*Cell)
	}
	return nil
}

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
	count := int(width) * int(height)
	cells := make([]Cell, 0, count)
	walkableCells := make([]Cell, 0, count / 3)
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			p := common.Point{X: uint32(i), Y: uint32(j)}
			//c := Cell{Coordinate: p.String(), Point: p}
			c := new(Cell)
			c.Coordinate = p.String()
			c.Point = p
			c.Objects = []MapObj{}
			if (common.BytesToUint32(bytes[offset:offset+4])^0xAA38AA38)&0x20000000 != 0 {
				c.Attribute = common.CellAttributeHighWall
			}
			if ((common.BytesToUint16(bytes[offset+6:offset+8]) ^ xor) & 0x8000) != 0 {
				c.Attribute = common.CellAttributeLowWall
			}
			if c.Attribute == common.CellAttributeWalk {
				walkableCells = append(walkableCells, *c)
			}
			cells = append(cells, *c)
			m.CoordinateCellMap.Store(p.String(), c)
			offset += 15
		}
	}
	m.Width = width
	m.Height = height
	m.Cells = cells
	m.WalkableCells = walkableCells
	return m
}
