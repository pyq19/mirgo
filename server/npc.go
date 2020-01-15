package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"time"
)

type NPC struct {
	MapObject
	Image    int
	TurnTime time.Time
}

func NewNPC(m *Map, ni *common.NpcInfo) *NPC {
	return &NPC{
		MapObject: MapObject{
			ID:               m.Env.NewObjectID(),
			Name:             ni.Name,
			Map:              m,
			CurrentLocation:  common.NewPoint(ni.LocationX, ni.LocationY),
			CurrentDirection: common.MirDirectionDown,
			Light:            0, // TODO
		},
		Image:    ni.Image,
		TurnTime: time.Now(),
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

func (n *NPC) GetDirection() common.MirDirection {
	return n.CurrentDirection
}

func (n *NPC) GetInfo() interface{} {
	return ServerMessage{}.ObjectNPC(n)
}

func (n *NPC) String() string {
	return fmt.Sprintf("NPC Coordinate: %s, ID: %d, name: %s\n", n.GetPoint().Coordinate(), n.ID, n.Name)
}

func (n *NPC) Broadcast(msg interface{}) {
	n.Map.Submit(NewTask(func(args ...interface{}) {
		grids := n.Map.AOI.GetSurroundGridsByCoordinate(n.GetCoordinate())
		for i := range grids {
			areaPlayers := grids[i].GetAllPlayer()
			for i := range areaPlayers {
				areaPlayers[i].Enqueue(msg)
			}
		}
	}))
}

func (n *NPC) Process() {
	if n.TurnTime.Before(time.Now()) {
		n.TurnTime = time.Now().Add(time.Second * time.Duration(G_Rand.RandInt(5, 15)))
		n.CurrentDirection = common.MirDirection(G_Rand.RandInt(0, 7))
		n.Broadcast(ServerMessage{}.ObjectTurn(n))
	}
}
