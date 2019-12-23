package main

import "github.com/yenkeia/mirgo/common"

// Respawn = Monster Obj
type Respawn struct {
	Cell *Cell
	Info *common.RespawnInfo
}

func (r *Respawn) Point() *common.Point {
	x := r.Info.LocationX
	y := r.Info.LocationY
	return common.NewPoint(x, y)
}

//func (r *Respawn) SetCell(c *Cell) {
//	r.Cell = c
//	r.Info.LocationX = int(c.Point.X)
//	r.Info.LocationY = int(c.Point.Y)
//}
