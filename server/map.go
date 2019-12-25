package main

import (
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// Map ...
type Map struct {
	Env    *Environ
	Id     int
	Width  uint16 // 测试用
	Height uint16 // 测试用
	Info   *common.MapInfo
	AOI    *AOIManager
	Cells  *sync.Map // key=Cell.Coordinate  value=*Cell
}

func (m *Map) Submit(t *Task) {
	m.Env.Game.Pool.EntryChan <- t
}

func (m *Map) GetCell(coordinate string) *Cell {
	v, ok := m.Cells.Load(coordinate)
	if !ok {
		return nil
	}
	return v.(*Cell)
}

func (m *Map) AddObject(obj interface{}) {
	switch o := obj.(type) {
	case *Player:
		coordinate := o.Point().Coordinate()
		grid := m.AOI.GetGridByCoordinate(coordinate)
		grid.AddPlayer(o)
	case *Respawn:
		coordinate := o.Point().Coordinate()
		grid := m.AOI.GetGridByCoordinate(coordinate)
		grid.AddRespawn(o)
	case *NPC:
		coordinate := o.Point().Coordinate()
		grid := m.AOI.GetGridByCoordinate(coordinate)
		grid.AddNPC(o)
	}
}
