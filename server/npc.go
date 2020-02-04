package server

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/setting"
	"os"
	"time"
)

type NPC struct {
	MapObject
	FilePath string
	Image    int
	Light    uint8
	TurnTime time.Time
}

func NewNPC(m *Map, ni *common.NpcInfo) *NPC {
	filePath := setting.Conf.NPCDirPath + ni.Filename + ".txt"
	if _, err := os.Stat(filePath); err != nil {
		log.Warnf("NPC %s 文件 %s 不存在\n", ni.Name, filePath)
		return nil
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
		FilePath: filePath,
		Image:    ni.Image,
		Light:    0, // TODO
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
	return fmt.Sprintf("NPC Coordinate: %s, ID: %d, name: %s, filepath: %s\n", n.GetPoint().Coordinate(), n.ID, n.Name, n.FilePath)
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
		n.TurnTime = time.Now().Add(time.Second * time.Duration(G_Rand.RandInt(20, 60)))
		n.CurrentDirection = common.MirDirection(G_Rand.RandInt(0, 2))
		n.Broadcast(ServerMessage{}.ObjectTurn(n))
	}
}
