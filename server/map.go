package main

import (
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// Map ...
type Map struct {
	Id                int
	Width             uint16 // 测试用
	Height            uint16 // 测试用
	Info              *common.MapInfo
	Cells             []Cell
	CoordinateCellMap *sync.Map // map[string]*Cell
	AOIManager        *AOIManager
	//WalkableCells     []Cell
}

func (m *Map) GetCell(coordinate string) *Cell {
	v, ok := m.CoordinateCellMap.Load(coordinate)
	if !ok {
		return nil
	}
	return v.(*Cell)
}

func (m *Map) AddObject(obj interface{}) {

}

func (m *Map) DeleteObject(obj interface{}) {

}

func (m *Map) UpdateObject(obj interface{}) {

}
