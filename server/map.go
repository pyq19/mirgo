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

func (m *Map) AddObject(obj interface{}) {
	m.AOI.AddObject(obj)
}

func (m *Map) GetCell(coordinate string) *Cell {
	v, ok := m.Cells.Load(coordinate)
	if !ok {
		return nil
	}
	return v.(*Cell)
}
