package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
)

type NPC struct {
	MapObject
}

func NewNPC(m *Map, ni *common.NpcInfo) *NPC {
	return &NPC{
		MapObject: MapObject{
			ID:               m.Env.NewObjectID(),
			Name:             ni.Name,
			Map:              m,
			CurrentLocation:  common.NewPoint(ni.LocationX, ni.LocationY),
			CurrentDirection: common.MirDirectionDown,
		},
	}
}

func (n *NPC) GetID() uint32 {
	return n.ID
}

func (n *NPC) GetRace() common.ObjectType {
	return common.ObjectTypeMerchant
}

func (n *NPC) GetCoordinate() string {
	return n.GetPoint().Coordinate()
}

func (n *NPC) GetPoint() common.Point {
	return n.CurrentLocation
}

func (n *NPC) GetCell() *Cell {
	return n.Map.GetCell(n.GetCoordinate())
}

func (n *NPC) String() string {
	return fmt.Sprintf("NPC Coordinate: %s, ID: %d, name: %s\n", n.GetPoint().Coordinate(), n.ID, n.Name)
}

func (n *NPC) Broadcast(msg interface{}) {

}

func (n *NPC) Process() {

}
