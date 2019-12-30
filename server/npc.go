package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
)

type NPC struct {
	ID   uint32
	Map  *Map
	Info *common.NpcInfo
}

func (n *NPC) Point() common.Point {
	x := n.Info.LocationX
	y := n.Info.LocationY
	return *common.NewPoint(x, y)
}

func NewNPC(m *Map, ni *common.NpcInfo) *NPC {
	return &NPC{
		ID:   m.Env.NewObjectID(),
		Map:  m,
		Info: ni,
	}
}

func (n *NPC) String() string {
	return fmt.Sprintf("NPC Coordinate: %s, ID: %d, name: %s\n", n.Point().Coordinate(), n.ID, n.Info.Name)
}
