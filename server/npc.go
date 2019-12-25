package main

import "github.com/yenkeia/mirgo/common"

type NPC struct {
	ID   string
	Info *common.NpcInfo
}

func (n *NPC) Point() *common.Point {
	x := n.Info.LocationX
	y := n.Info.LocationY
	return common.NewPoint(x, y)
}
