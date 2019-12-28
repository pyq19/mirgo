package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

const (
	LOGIN = iota
	SELECT
	GAME
	DISCONNECTED
)

// Player ...
type Player struct {
	Map       *Map
	AccountID int
	GameStage int
	Session   *cellnet.Session
	Character *common.Character
	Magics    []*common.MagicInfo
	UserItems []*common.UserItem
}

func (p *Player) Point() common.Point {
	x := int(p.Character.CurrentLocationX)
	y := int(p.Character.CurrentLocationY)
	return *common.NewPoint(x, y)
}

func (p *Player) Send(msg interface{}) {
	(*p.Session).Send(msg)
}

func (p *Player) NotifySurroundingPlayer(msg interface{}) {
	p.Map.Submit(NewTask(func(args ...interface{}) {
		grids := p.Map.AOI.GetSurroundGridsByCoordinate(p.Point().Coordinate())
		for i := range grids {
			grids[i].Players.Range(func(k, v interface{}) bool {
				o := v.(*Player)
				if p.Character.ID != o.Character.ID {
					o.Send(msg)
				}
				return true
			})
		}
	}))
}

// TODO
func (p *Player) Turn(direction common.MirDirection) {
	p.NotifySurroundingPlayer(server.ObjectTurn{
		ObjectID:  uint32(p.Character.ID),
		Location:  p.Point(),
		Direction: direction,
	})
}

// TODO
func (p *Player) Walk(direction common.MirDirection) {

}
