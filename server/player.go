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
	MapObject
	AccountID int
	GameStage int
	Session   *cellnet.Session
	Character *common.Character // 仅在与数据库交互时使用
	Magics    []*common.MagicInfo
	UserItems []*common.UserItem
}

func (p *Player) Point() common.Point {
	return *p.CurrentLocation
}

func (p *Player) Coordinate() string {
	return p.Point().Coordinate()
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
				if p.ID != o.ID {
					o.Send(msg)
				}
				return true
			})
		}
	}))
}

func (p *Player) CurrentCell() *Cell {
	return p.Map.GetCell(p.Point().Coordinate())
}

// TODO
func (p *Player) Turn(direction common.MirDirection) {
	p.NotifySurroundingPlayer(server.ObjectTurn{
		ObjectID:  uint32(p.ID),
		Location:  p.Point(),
		Direction: direction,
	})
}

func (p *Player) Walk(direction common.MirDirection, point *common.Point) {
	p.NotifySurroundingPlayer(&server.ObjectWalk{
		ObjectID:  uint32(p.ID),
		Location:  p.Point(),
		Direction: direction,
	})
	p.CurrentDirection = direction
	p.CurrentLocation = point
}

func (p *Player) Run(direction common.MirDirection, point *common.Point) {
	p.NotifySurroundingPlayer(&server.ObjectRun{
		ObjectID:  uint32(p.ID),
		Location:  p.Point(),
		Direction: direction,
	})
	p.CurrentDirection = direction
	p.CurrentLocation = point
}

func (p *Player) Chat(message string) {

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
