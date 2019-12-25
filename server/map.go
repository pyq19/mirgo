package main

import (
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// Map ...
type Map struct {
	Env               *Environ
	Id                int
	Width             uint16 // 测试用
	Height            uint16 // 测试用
	Info              *common.MapInfo
	Cells             []Cell
	CoordinateCellMap *sync.Map // map[string]*Cell
	AOI               *AOIManager
}

func (m *Map) getCell(coordinate string) *Cell {
	v, ok := m.CoordinateCellMap.Load(coordinate)
	if !ok {
		return nil
	}
	return v.(*Cell)
}

func (m *Map) AddObject(obj interface{}) {
	switch o := obj.(type) {
	case *Player:
		p := common.NewPoint(int(o.Character.CurrentLocationX), int(o.Character.CurrentLocationY))
		c := m.getCell(p.Coordinate())
		c.SetPlayer(o)
	case *Respawn:
		m.getCell(o.Point().Coordinate()).SetRespawn(o)
	case *NPC:
	}
}

func (m *Map) DeleteObject(obj interface{}) {

}

func (m *Map) UpdateObject(obj interface{}) {

}
