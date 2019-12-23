package main

import (
	"github.com/yenkeia/mirgo/common"
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

func (m *Map) GetCell(coordinate string) *Cell {
	if v, ok := m.CoordinateCellMap.Load(coordinate); ok {
		return v.(*Cell)
	}
	return nil
}
