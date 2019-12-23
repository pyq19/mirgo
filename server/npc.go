package main

import "github.com/yenkeia/mirgo/common"

type NPC struct {
	Cell *Cell
	Info *common.NpcInfo
}

func (n *NPC) Point() *common.Point {
	x := n.Info.LocationX
	y := n.Info.LocationY
	return common.NewPoint(x, y)
}
