package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/server/script"
	"github.com/yenkeia/mirgo/setting"
	"time"
)

type NPC struct {
	MapObject
	Image    int
	Light    uint8
	TurnTime time.Time
	Script   *script.Script
}

func NewNPC(m *Map, ni *common.NpcInfo) *NPC {
	sc, err := script.LoadFile(setting.Conf.NPCDirPath + ni.Filename + ".txt")
	if err != nil {
		log.Warnf("NPC %s 脚本加载失败: %s\n", ni.Name, err.Error())
	}
	return &NPC{
		MapObject: MapObject{
			ID:               m.Env.NewObjectID(),
			Name:             ni.Name,
			NameColor:        common.Color{R: 255, G: 255, B: 255},
			Map:              m,
			CurrentLocation:  common.NewPoint(ni.LocationX, ni.LocationY),
			CurrentDirection: common.MirDirection(G_Rand.RandInt(0, 2)),
		},
		Image:    ni.Image,
		Light:    0, // TODO
		TurnTime: time.Now(),
		Script:   sc,
	}
}

func (n *NPC) CallScript(p *Player, key string) ([]string, error) {
	say, err := n.Script.Call(n, p, key)
	if err != nil {
		return nil, err
	}
	return say, nil
}

func (n *NPC) GetID() uint32 {
	return n.ID
}

func (n *NPC) GetRace() common.ObjectType {
	return common.ObjectTypeMerchant
}

func (n *NPC) GetPoint() common.Point {
	return n.CurrentLocation
}

func (n *NPC) GetCell() *Cell {
	return n.Map.GetCell(n.CurrentLocation)
}

func (n *NPC) GetDirection() common.MirDirection {
	return n.CurrentDirection
}

func (n *NPC) GetInfo() interface{} {
	res := &server.ObjectNPC{
		ObjectID:  n.ID,
		Name:      n.Name,
		NameColor: -16711936, // TODO
		Image:     uint16(n.Image),
		Color:     0, // TODO
		Location:  n.GetPoint(),
		Direction: n.GetDirection(),
		QuestIDs:  []int32{}, // TODO
	}
	return res
}

func (n *NPC) IsAttackTarget(IMapObject) bool {
	return false
}

func (n *NPC) IsFriendlyTarget(attacker IMapObject) bool {
	return true
}

func (n *NPC) GetBaseStats() BaseStats {
	return BaseStats{}
}

func (n *NPC) String() string {
	return fmt.Sprintf("NPC pos: %s, ID: %d, name: %s\n", n.GetPoint(), n.ID, n.Name)
}

func (n *NPC) Broadcast(msg interface{}) {
	n.Map.Submit(NewTask(func(args ...interface{}) {
		grids := n.Map.AOI.GetSurroundGrids(n.CurrentLocation)
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
		n.TurnTime = time.Now().Add(time.Second * time.Duration(G_Rand.RandInt(20, 60)))
		n.CurrentDirection = common.MirDirection(G_Rand.RandInt(0, 2))
		n.Broadcast(ServerMessage{}.ObjectTurn(n))
	}
}
