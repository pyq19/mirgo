package main

import (
	"github.com/yenkeia/mirgo/common"
	"io/ioutil"
	"os"
	"sync"
)

// Map ...
type Map struct {
	Id                int
	Info              *common.MapInfo
	Cells             []Cell
	CoordinateCellMap *sync.Map // map[string]*Cell
	//WalkableCells     []Cell
}

// InitMaps ...
func (e *Environ) InitMaps() {
	mapDirPath := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/"
	//e.Maps = make([]Map, 386)
	e.Maps = new(sync.Map)
	for _, mi := range e.GameDB.MapInfos {
		if mi.Id == 1 {
			m := GetMapV1(GetMapBytes(mapDirPath + mi.Filename + ".map"))
			m.Id = mi.Id
			e.Maps.Store(1, m)
			break
		}
	}
}

type Cell struct {
	Coordinate string       // 坐标 x,y
	Point      common.Point // 坐标点
	Attribute  common.CellAttribute
	Respawn    *Respawn
	NPC        *NPC
	Player     *Player
}

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
	//walkableCells := make([]Cell, 0, count/3)
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			p := common.Point{X: uint32(i), Y: uint32(j)}
			//c := Cell{Coordinate: p.String(), Point: p}
			c := new(Cell)
			c.Coordinate = p.String()
			c.Point = p
			if (common.BytesToUint32(bytes[offset:offset+4])^0xAA38AA38)&0x20000000 != 0 {
				c.Attribute = common.CellAttributeHighWall
			}
			if ((common.BytesToUint16(bytes[offset+6:offset+8]) ^ xor) & 0x8000) != 0 {
				c.Attribute = common.CellAttributeLowWall
			}
			//if c.Attribute == common.CellAttributeWalk {
			//	walkableCells = append(walkableCells, *c)
			//}
			cells = append(cells, *c)
			m.CoordinateCellMap.Store(p.String(), c)
			offset += 15
		}
	}
	m.Cells = cells
	//m.WalkableCells = walkableCells
	return m
}

func (m *Map) AddRespawn(r *Respawn) {
	c := m.GetCell(common.NewPoint(r.Info.LocationY, r.Info.LocationY).String())
	c.Respawn = r
}

func (m *Map) AddNPC(n *NPC) {

}

func (m *Map) AddPlayer(p *Player) {

}
