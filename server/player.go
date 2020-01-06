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

func (p *Player) canRun() bool {
	return true
}

func (p *Player) canAttack() bool {
	return true
}

func (p *Player) canRegen() bool {
	return true
}

func (p *Player) canCast() bool {
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

func (p *Player) Inspect(id uint32) {

}

func (p *Player) ChangeAMode(mode common.AttackMode) {

}

func (p *Player) ChangePMode(mode common.AttackMode) {

}

func (p *Player) ChangeTrade(trade bool) {

}

func (p *Player) getAttackPower(minDC, maxDC uint16) int {
	return 0
}

func (p *Player) isAttackTarget(attacker *Player) bool {
	return true
}

func (p *Player) attacked(attacker *Player, finalDamage int, defenceType common.DefenceType, damageWeapon bool) {

}

func (p *Player) Attack(direction common.MirDirection, spell common.Spell) {
	if !p.canAttack() {
		p.Enqueue(server.UserLocation{
			Location:  p.Point(),
			Direction: p.CurrentDirection,
		})
		return
	}
	p.CurrentDirection = direction
	p.Enqueue(server.UserLocation{
		Location:  p.Point(),
		Direction: direction,
	})
	p.Broadcast(server.ObjectAttack{
		ObjectID:  p.ID,
		Location:  p.Point(),
		Direction: p.CurrentDirection,
		Spell:     common.SpellNone,
		Level:     0,
		Type:      0,
	})
	target := p.Point().NextPoint(p.CurrentDirection, 1)
	c := p.Map.GetCell(target.Coordinate())
	if c == nil || c.IsEmpty() {
		return
	}
	damageBase := p.getAttackPower(p.MinDC, p.MaxDC)
	damageFinal := damageBase // TODO
	for i := range c.Objects {
		o := c.Objects[i]
		switch c.GetRace(o) {
		case common.ObjectTypePlayer:
			ob := o.(*Player)
			if !ob.isAttackTarget(p) {
				continue
			}
			ob.attacked(p, damageFinal, common.DefenceTypeAgility, false)
		case common.ObjectTypeMonster:
			ob := o.(*Monster)
			if !ob.isAttackTarget(p) {
				continue
			}
			ob.attacked(p, damageFinal, common.DefenceTypeAgility, false)
		}
	}
}

func (p *Player) RangeAttack(direction common.MirDirection, location common.Point, id uint32) {

}

func (p *Player) Harvest(direction common.MirDirection) {

}

func (p *Player) CallNPC(id uint32, key string) {

}

func (p *Player) TalkMonsterNPC(id uint32) {

}

func (p *Player) BuyItem(index uint64, count uint32, panelType common.PanelType) {

}

func (p *Player) CraftItem() {

}

func (p *Player) SellItem(id uint64, count uint32) {

}

func (p *Player) RepairItem(id uint64) {

}

func (p *Player) BuyItemBack(id uint64, count uint32) {

}

func (p *Player) SRepairItem(id uint64) {

}

func (p *Player) MagicKey(spell common.Spell, key uint8) {

}

func (p *Player) Magic(spell common.Spell, direction common.MirDirection, id uint32, location common.Point) {

}

func (p *Player) SwitchGroup(group bool) {

}

func (p *Player) AddMember(name string) {

}

func (p *Player) DelMember(name string) {

}

func (p *Player) GroupInvite(invite bool) {

}

func (p *Player) TownRevive() {

}

func (p *Player) SpellToggle(spell common.Spell, use bool) {

}

func (p *Player) ConsignItem(id uint64, price uint32) {

}

func (p *Player) MarketSearch(match string) {

}

func (p *Player) MarketRefresh() {

}

func (p *Player) MarketPage(page int32) {

}

func (p *Player) MarketBuy(id uint64) {

}

func (p *Player) MarketGetBack(id uint64) {

}

func (p *Player) RequestUserName(id uint32) {

}

func (p *Player) RequestChatItem(id uint64) {

}

func (p *Player) EditGuildMember(name string, name2 string, index uint8, changeType uint8) {

}
