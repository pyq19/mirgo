package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
	"strings"
)

const (
	LOGIN = iota
	SELECT
	GAME
	DISCONNECTED
)

// Player ...
type Player struct {
	MapObject
	AccountID int
	GameStage int
	Session   *cellnet.Session
	Magics    []*common.MagicInfo
	UserItems []*common.UserItem
}

func (p *Player) Point() common.Point {
	return *p.CurrentLocation
}

func (p *Player) CurrentCell() *Cell {
	return p.Map.GetCell(p.Point().Coordinate())
}

func (p *Player) Coordinate() string {
	return p.Point().Coordinate()
}

func (p *Player) Enqueue(msg interface{}) {
	(*p.Session).Send(msg)
}

func (p *Player) ReceiveChat(text string, ct common.ChatType) {
	p.Enqueue(server.Chat{
		Message: text,
		Type:    ct,
	})
}

func (p *Player) Broadcast(msg interface{}) {
	p.Map.Submit(NewTask(func(args ...interface{}) {
		grids := p.Map.AOI.GetSurroundGridsByCoordinate(p.Point().Coordinate())
		for i := range grids {
			grids[i].Players.Range(func(k, v interface{}) bool {
				o := v.(*Player)
				if p.ID != o.ID {
					o.Enqueue(msg)
				}
				return true
			})
		}
	}))
}

func (p *Player) Process() {

}

func (p *Player) canMove() bool {
	return true
}

func (p *Player) canWalk() bool {
	return true
}

func (p *Player) Turn(direction common.MirDirection) {
	if p.canMove() {
		p.Broadcast(server.ObjectTurn{
			ObjectID:  p.ID,
			Location:  p.Point(),
			Direction: direction,
		})
		p.CurrentDirection = direction
	}
	p.Enqueue(server.UserLocation{
		Location:  p.Point(),
		Direction: p.CurrentDirection,
	})
}

func (p *Player) Walk(direction common.MirDirection) {
	if !p.canMove() || !p.canWalk() {
		p.Enqueue(server.UserLocation{
			Location:  p.Point(),
			Direction: p.CurrentDirection,
		})
		return
	}
	n := p.Point().NextPoint(direction, 1)
	ok := p.Map.UpdateObject(p, n)
	if !ok {
		p.Enqueue(server.UserLocation{
			Location:  p.Point(),
			Direction: p.CurrentDirection,
		})
		return
	}
	p.Broadcast(server.ObjectWalk{
		ObjectID:  p.ID,
		Location:  p.Point(),
		Direction: direction,
	})
	p.CurrentDirection = direction
	p.CurrentLocation = n
}

func (p *Player) Run(direction common.MirDirection) {
	n1 := p.Point().NextPoint(direction, 1)
	n2 := p.Point().NextPoint(direction, 2)
	if ok := p.Map.UpdateObject(p, n1, n2); !ok {
		p.Enqueue(server.UserLocation{
			Location:  p.Point(),
			Direction: p.CurrentDirection,
		})
		return
	}
	p.Broadcast(server.ObjectRun{
		ObjectID:  p.ID,
		Location:  p.Point(),
		Direction: direction,
	})
	p.CurrentDirection = direction
	p.CurrentLocation = n2
}

func (p *Player) Chat(message string) {
	// private message
	if strings.HasPrefix(message, "/") {
		return
	}
	// group
	if strings.HasPrefix(message, "!!") {
		return
	}
	message = p.Name + ":" + message
	msg := server.ObjectChat{
		ObjectID: p.ID,
		Text:     message,
		Type:     common.ChatTypeNormal,
	}
	p.Enqueue(msg)
	p.Broadcast(msg)
}

func (p *Player) MoveItem(grid common.MirGridType, from int32, to int32) {

}

func (p *Player) StoreItem(from int32, to int32) {

}

func (p *Player) DepositRefineItem(from int32, to int32) {

}

func (p *Player) RetrieveRefineItem(from int32, to int32) {

}

func (p *Player) RefineCancel() {

}

func (p *Player) RefineItem(id uint64) {

}

func (p *Player) CheckRefine(id uint64) {

}

func (p *Player) ReplaceWeddingRing(id uint64) {

}

func (p *Player) DepositTradeItem(from int32, to int32) {

}

func (p *Player) RetrieveTradeItem(from int32, to int32) {

}

func (p *Player) TakeBackItem(from int32, to int32) {

}

func (p *Player) MergeItem(from common.MirGridType, to common.MirGridType, from2 uint64, to2 uint64) {

}

func (p *Player) EquipItem(grid common.MirGridType, id uint64, to int32) {

}

func (p *Player) RemoveItem(grid common.MirGridType, id uint64, to int32) {

}

func (p *Player) RemoveSlotItem(grid common.MirGridType, id uint64, to int32, to2 common.MirGridType) {

}

func (p *Player) SplitItem(grid common.MirGridType, id uint64, count uint32) {

}

func (p *Player) UseItem(id uint64) {

}

func (p *Player) DropItem(id uint64, count uint32) {

}

func (p *Player) DropGold(amount uint32) {

}

func (p *Player) PickUp() {

}
