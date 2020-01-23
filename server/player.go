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
	AccountID int
	GameStage int
	Session   *cellnet.Session
	MapObject
	Character
}

func (p *Player) Enqueue(msg interface{}) {
	if msg == nil {
		log.Errorln("warning: enqueue nil message")
		return
	}
	(*p.Session).Send(msg)
}

func (p *Player) ReceiveChat(text string, ct common.ChatType) {
	p.Enqueue(&server.Chat{
		Message: text,
		Type:    ct,
	})
}

func (p *Player) Broadcast(msg interface{}) {
	p.Map.Submit(NewTask(func(args ...interface{}) {
		grids := p.Map.AOI.GetSurroundGridsByCoordinate(p.Point().Coordinate())
		for i := range grids {
			areaPlayers := grids[i].GetAllPlayer()
			for i := range areaPlayers {
				if p.GetID() == areaPlayers[i].GetID() {
					continue
				}
				areaPlayers[i].Enqueue(msg)
			}
		}
	}))
}

func (p *Player) Point() common.Point {
	return p.GetPoint()
}

func (p *Player) GetID() uint32 {
	return p.ID
}

func (p *Player) GetRace() common.ObjectType {
	return common.ObjectTypePlayer
}

func (p *Player) GetCoordinate() string {
	return p.GetPoint().Coordinate()
}

func (p *Player) GetPoint() common.Point {
	return p.CurrentLocation
}

func (p *Player) GetCell() *Cell {
	return p.Map.GetCell(p.GetCoordinate())
}

func (p *Player) GetDirection() common.MirDirection {
	return p.CurrentDirection
}

func (p *Player) GetInfo() interface{} {
	res := &server.ObjectPlayer{
		ObjectID:         p.ID,
		Name:             p.Name,
		GuildName:        p.GuildName,
		GuildRankName:    p.GuildRankName,
		NameColor:        p.NameColor.ToInt32(),
		Class:            p.Class,
		Gender:           p.Gender,
		Level:            p.Level,
		Location:         p.GetPoint(),
		Direction:        p.GetDirection(),
		Hair:             p.Hair,
		Light:            p.Light,
		Weapon:           int16(p.LooksWeapon),
		WeaponEffect:     int16(p.LooksWeaponEffect),
		Armour:           int16(p.LooksArmour),
		Poison:           common.PoisonTypeNone, // TODO
		Dead:             p.IsDead(),
		Hidden:           p.IsHidden(),
		Effect:           common.SpellEffectNone, // TODO
		WingEffect:       uint8(p.LooksWings),
		Extra:            false,                      // TODO
		MountType:        0,                          // TODO
		RidingMount:      false,                      // TODO
		Fishing:          false,                      // TODO
		TransformType:    0,                          // TODO
		ElementOrbEffect: 0,                          // TODO
		ElementOrbLvl:    0,                          // TODO
		ElementOrbMax:    0,                          // TODO
		Buffs:            make([]common.BuffType, 0), // TODO
		LevelEffects:     common.LevelEffectsNone,    // TODO
	}
	return res
}

func (p *Player) GetCurrentGrid() *Grid {
	return p.Map.AOI.GetGridByPoint(p.Point())
}

func (p *Player) EnqueueAreaObjects(oldGrid, newGrid *Grid) {
	oldAreaGrids := make([]*Grid, 0)
	if oldGrid != nil {
		oldAreaGrids = p.Map.AOI.GetSurroundGridsByGridID(oldGrid.GID)
	}
	newAreaGrids := p.Map.AOI.GetSurroundGridsByGridID(newGrid.GID)
	send := make(map[int]bool)
	for i := range newAreaGrids {
		ng := newAreaGrids[i]
		send[ng.GID] = true
		for j := range oldAreaGrids {
			og := oldAreaGrids[j]
			if ng.GID == og.GID {
				send[ng.GID] = false
			}
		}
	}
	newAreaObjects := make([]IMapObject, 0)
	for i := range newAreaGrids {
		ng := newAreaGrids[i]
		if send[ng.GID] {
			newAreaObjects = append(newAreaObjects, ng.GetAllObjects()...)
		}
	}
	for i := range newAreaObjects {
		if o := newAreaObjects[i]; o.GetID() != p.GetID() {
			p.Enqueue(ServerMessage{}.Object(o))
		}
	}
	drop := make(map[int]bool)
	for i := range oldAreaGrids {
		og := oldAreaGrids[i]
		drop[og.GID] = true
		for j := range newAreaGrids {
			ng := newAreaGrids[j]
			if og.GID == ng.GID {
				drop[og.GID] = false
			}
		}
	}
	oldAreaObjects := make([]IMapObject, 0)
	for i := range oldAreaGrids {
		og := oldAreaGrids[i]
		if drop[og.GID] {
			oldAreaObjects = append(oldAreaObjects, og.GetAllObjects()...)
		}
	}
	for i := range oldAreaObjects {
		if o := oldAreaObjects[i]; o.GetID() != p.GetID() {
			p.Enqueue(ServerMessage{}.ObjectRemove(o))
		}
	}
}

func (p *Player) StartGame() {
	p.ReceiveChat("这是一个以学习为目的传奇服务端", common.ChatTypeSystem)
	p.ReceiveChat("如有任何建议、疑问欢迎交流", common.ChatTypeSystem)
	p.ReceiveChat("源码地址 https://github.com/yenkeia/mirgo", common.ChatTypeSystem)
	p.EnqueueItemInfos()
	p.RefreshStats()
	p.EnqueueQuestInfo()
	p.Enqueue(ServerMessage{}.MapInformation(p.Map.Info))
	p.Enqueue(ServerMessage{}.UserInformation(p))
	p.Enqueue(ServerMessage{}.TimeOfDay(common.LightSettingDay))
	p.EnqueueAreaObjects(nil, p.Map.AOI.GetGridByPoint(p.GetPoint()))
	p.Enqueue(ServerMessage{}.NPCResponse([]string{}))
	p.Broadcast(ServerMessage{}.ObjectPlayer(p))
}

func (p *Player) StopGame(reason int) {
	p.Broadcast(ServerMessage{}.ObjectRemove(p))
}

func (p *Player) Turn(direction common.MirDirection) {
	if !p.CanMove() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectTurn(p))
}

func (p *Player) Walk(direction common.MirDirection) {
	if !p.CanMove() || !p.CanWalk() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	n := p.Point().NextPoint(direction, 1)
	ok := p.Map.UpdateObject(p, n)
	if !ok {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.CurrentLocation = n
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectWalk(p))
}

func (p *Player) Run(direction common.MirDirection) {
	n1 := p.Point().NextPoint(direction, 1)
	n2 := p.Point().NextPoint(direction, 2)
	if ok := p.Map.UpdateObject(p, n1, n2); !ok {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.CurrentLocation = n2
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectRun(p))
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
	msg := ServerMessage{}.ObjectChat(p, message, common.ChatTypeNormal)
	p.Enqueue(msg)
	p.Broadcast(msg)
}

func (p *Player) MoveItem(mirGridType common.MirGridType, from int32, to int32) {
	msg := &server.MoveItem{
		Grid:    mirGridType,
		From:    from,
		To:      to,
		Success: false,
	}
	switch mirGridType {
	case common.MirGridTypeInventory:
		l := len(p.Inventory)
		if from > 0 && to > 0 && int(from) < l && int(to) < l {
			array := p.Inventory
			i := array[to]
			array[to] = array[from]
			array[from] = i
			msg.Success = true
		}
	case common.MirGridTypeStorage:
		// TODO
	case common.MirGridTypeTrade:
		// TODO
	case common.MirGridTypeRefine:
		// TODO
	}
	p.Enqueue(msg)
}

// TODO
func (p *Player) StoreItem(from int32, to int32) {
	msg := &server.StoreItem{
		From:    from,
		To:      to,
		Success: false,
	}
	p.Enqueue(msg)
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

func (p *Player) EquipItem(mirGridType common.MirGridType, id uint64, to int32) {
	var msg = &server.EquipItem{
		Grid:     mirGridType,
		UniqueID: id,
		To:       to,
		Success:  false,
	}
	if l := len(p.Equipment); to < 0 || int(to) >= l {
		p.Enqueue(msg)
		return
	}
	switch mirGridType {
	case common.MirGridTypeInventory:
		index, item := p.GetUserItemByID(mirGridType, id)
		if item == nil {
			p.Enqueue(msg)
			return
		}
		p.Inventory[index] = p.Equipment[to]
		p.Equipment[to] = *item
		msg.Success = true
		p.RefreshStats()
	case common.MirGridTypeStorage:
		// TODO
	}
	p.Enqueue(msg)
}

// TODO
func (p *Player) RemoveItem(mirGridType common.MirGridType, id uint64, to int32) {
	msg := &server.RemoveItem{
		Grid:     mirGridType,
		UniqueID: id,
		To:       to,
		Success:  false,
	}
	p.Enqueue(msg)
}

func (p *Player) RemoveSlotItem(grid common.MirGridType, id uint64, to int32, to2 common.MirGridType) {

}

func (p *Player) SplitItem(grid common.MirGridType, id uint64, count uint32) {

}

func (p *Player) UseItem(id uint64) {

}

func (p *Player) DropItem(id uint64, count uint32) {
	msg := &server.DropItem{
		UniqueID: id,
		Count:    count,
		Success:  false,
	}
	index, item := p.GetUserItemByID(common.MirGridTypeInventory, id)
	if item == nil || item.ID == 0 {
		p.Enqueue(msg)
		return
	}
	obj := &Item{
		MapObject: MapObject{
			ID:  p.Map.Env.NewObjectID(),
			Map: p.Map,
		},
		Gold:     0,
		UserItem: item,
	}
	if dropMsg, ok := obj.Drop(p.GetPoint(), 1); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}
	if count >= item.Count {
		p.Inventory[index] = common.UserItem{}
	} else {
		p.Inventory[index].Count -= count
	}
	p.RefreshBagWeight()
	msg.Success = true
	p.Enqueue(msg)
}

func (p *Player) DropGold(gold uint64) {
	if p.Gold < gold {
		return
	}
	obj := &Item{
		MapObject: MapObject{
			ID:  p.Map.Env.NewObjectID(),
			Map: p.Map,
		},
		Gold:     gold,
		UserItem: nil,
	}
	if dropMsg, ok := obj.Drop(p.GetPoint(), 3); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}
	p.Gold -= gold
	p.Enqueue(&server.LoseGold{Gold: uint32(gold)})
}

func (p *Player) PickUp() {
	if p.IsDead() {
		return
	}
	c := p.GetCell()
	if c == nil {
		return
	}
	items := make([]*Item, 0)
	c.Objects.Range(func(k, v interface{}) bool {
		if o, ok := v.(*Item); ok {
			if o.UserItem == nil {
				p.GainGold(o.Gold)
				items = append(items, o)
			} else {
				if p.GainItem(o.UserItem) {
					items = append(items, o)
				}
			}
		}
		return true
	})
	for i := range items {
		o := items[i]
		p.Map.DeleteObject(o)
		o.Broadcast(ServerMessage{}.ObjectRemove(o))
	}
}

func (p *Player) Inspect(id uint32) {
	o := p.Map.Env.GetPlayer(id)
	for i := range o.Equipment {
		item := p.Map.Env.GameDB.GetItemInfoByID(int(o.Equipment[i].ItemID))
		if item != nil {
			p.EnqueueItemInfo(item.ID)
		}
	}
	p.Enqueue(ServerMessage{}.PlayerInspect(o))
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
	if !p.CanAttack() {
		p.Enqueue(&server.UserLocation{
			Location:  p.Point(),
			Direction: p.CurrentDirection,
		})
		return
	}
	p.CurrentDirection = direction
	p.Enqueue(&server.UserLocation{
		Location:  p.Point(),
		Direction: direction,
	})
	p.Broadcast(&server.ObjectAttack{
		ObjectID:  p.ID,
		Location:  p.Point(),
		Direction: p.CurrentDirection,
		Spell:     common.SpellNone,
		Level:     0,
		Type:      0,
	})
	//target := p.Point().NextPoint(p.CurrentDirection, 1)
	//c := p.Map.GetCell(target.Coordinate())
	//if c == nil || c.IsEmpty() {
	//	return
	//}
	//damageBase := p.getAttackPower(p.MinDC, p.MaxDC)
	//damageFinal := damageBase // TODO
	//for i := range c.Objects {
	//	o := c.Objects[i]
	//	switch c.GetRace(o) {
	//	case common.ObjectTypePlayer:
	//		ob := o.(*Player)
	//		if !ob.isAttackTarget(p) {
	//			continue
	//		}
	//		ob.attacked(p, damageFinal, common.DefenceTypeAgility, false)
	//	case common.ObjectTypeMonster:
	//		ob := o.(*Monster)
	//		if !ob.isAttackTarget(p) {
	//			continue
	//		}
	//		ob.attacked(p, damageFinal, common.DefenceTypeAgility, false)
	//	}
	//}
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

func (p *Player) getMagic(spell common.Spell) *common.UserMagic {
	return nil
}

func (p *Player) Magic(spell common.Spell, direction common.MirDirection, id uint32, location common.Point) {
	var (
		um *common.UserMagic
	)
	if !p.CanCast() {
		goto err
	}
	um = p.getMagic(spell)
	if um == nil {
		goto err
	}
	// TODO
err:
	p.Enqueue(&server.UserLocation{
		Location:  p.Point(),
		Direction: p.CurrentDirection,
	})
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
