package main

import (
	"sync"
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/setting"
)

type Character struct {
	MapObject
	Player             *Player
	HP                 uint16
	MP                 uint16
	Level              uint16
	Experience         int64
	Gold               uint64
	GuildName          string
	GuildRankName      string
	Class              common.MirClass
	Gender             common.MirGender
	Hair               uint8
	Light              uint8
	Inventory          []common.UserItem // 46
	Equipment          []common.UserItem // 14
	QuestInventory     []common.UserItem // 40
	Trade              []common.UserItem // 10
	Refine             []common.UserItem // 16
	LooksArmour        int
	LooksWings         int
	LooksWeapon        int
	LooksWeaponEffect  int
	SendItemInfo       []common.ItemInfo
	CurrentBagWeight   int
	MaxHP              uint16
	MaxMP              uint16
	MinAC              uint16 // 物理防御力
	MaxAC              uint16
	MinMAC             uint16 // 魔法防御力
	MaxMAC             uint16
	MinDC              uint16 // 攻击力
	MaxDC              uint16
	MinMC              uint16 // 魔法力
	MaxMC              uint16
	MinSC              uint16 // 道术力
	MaxSC              uint16
	MaxExperience      int64
	Accuracy           uint8
	Agility            uint8
	CriticalRate       uint8
	CriticalDamage     uint8
	MaxBagWeight       uint16 //Other Stats;
	MaxWearWeight      uint16
	MaxHandWeight      uint16
	ASpeed             int8
	Luck               int8
	LifeOnHit          uint8
	HpDrainRate        uint8
	Reflect            uint8
	MagicResist        uint8
	PoisonResist       uint8
	HealthRecovery     uint8
	SpellRecovery      uint8
	PoisonRecovery     uint8
	Holy               uint8
	Freezing           uint8
	PoisonAttack       uint8
	ExpRateOffset      float32
	ItemDropRateOffset float32
	MineRate           uint8
	GemRate            uint8
	FishRate           uint8
	CraftRate          uint8
	GoldDropRateOffset float32
	AttackBonus        uint8
	Magics             []common.UserMagic
	ActionList         *sync.Map // map[uint32]DelayedAction
	Health             Health    // 状态恢复
}

type Health struct {
	// 生命药水回复
	HPPotValue    int           // 回复总值
	HPPotPerValue int           // 一次回复多少
	HPPotNextTime *time.Time    // 下次生效时间
	HPPotDuration time.Duration // 两次生效时间间隔
	HPPotTickNum  int           // 总共跳几次
	HPPotTickTime int           // 当前第几跳
	// 魔法药水回复
	MPPotValue    int
	MPPotPerValue int
	MPPotNextTime *time.Time
	MPPotDuration time.Duration
	MPPotTickNum  int
	MPPotTickTime int
	// 角色生命/魔法回复
	HealNextTime *time.Time
	HealDuration time.Duration
}

func NewCharacter(g *Game, p *Player, c *common.Character) *Character {
	userItemIDIndexMap := make(map[int]int)
	cui := make([]common.CharacterUserItem, 0, 100)
	g.DB.Table("character_user_item").Where("character_id = ?", c.ID).Find(&cui)
	is := make([]int, 0, 46)
	es := make([]int, 0, 14)
	qs := make([]int, 0, 40)
	for _, i := range cui {
		switch common.UserItemType(i.Type) {
		case common.UserItemTypeInventory:
			is = append(is, i.UserItemID)
		case common.UserItemTypeEquipment:
			es = append(es, i.UserItemID)
		case common.UserItemTypeQuestInventory:
			qs = append(qs, i.UserItemID)
		}
		userItemIDIndexMap[i.UserItemID] = i.Index
	}
	inventory := make([]common.UserItem, 46)
	equipment := make([]common.UserItem, 14)
	questInventory := make([]common.UserItem, 40)
	trade := make([]common.UserItem, 0)
	refine := make([]common.UserItem, 0)
	uii := make([]common.UserItem, 0, 46)
	uie := make([]common.UserItem, 0, 14)
	uiq := make([]common.UserItem, 0, 40)
	g.DB.Table("user_item").Where("id in (?)", is).Find(&uii)
	g.DB.Table("user_item").Where("id in (?)", es).Find(&uie)
	g.DB.Table("user_item").Where("id in (?)", qs).Find(&uiq)
	for _, v := range uii {
		inventory[userItemIDIndexMap[int(v.ID)]] = v
	}
	for _, v := range uie {
		equipment[userItemIDIndexMap[int(v.ID)]] = v
	}
	for _, v := range uiq {
		questInventory[userItemIDIndexMap[int(v.ID)]] = v
	}
	magics := make([]common.UserMagic, 0)
	g.DB.Table("user_magic").Where("character_id = ?", c.ID).Find(&magics)
	healNextTime := time.Now().Add(10 * time.Second)
	return &Character{
		Player:         p,
		HP:             c.HP,
		MP:             c.MP,
		Level:          c.Level,
		Experience:     c.Experience,
		Gold:           c.Gold,
		GuildName:      "", // TODO
		GuildRankName:  "", // TODO
		Class:          c.Class,
		Gender:         c.Gender,
		Hair:           c.Hair,
		Inventory:      inventory,
		Equipment:      equipment,
		QuestInventory: questInventory,
		Trade:          trade,
		Refine:         refine,
		SendItemInfo:   make([]common.ItemInfo, 0),
		MaxExperience:  100,
		Magics:         magics,
		ActionList:     new(sync.Map),
		Health: Health{
			HPPotNextTime: new(time.Time),
			HPPotDuration: 3 * time.Second,
			MPPotNextTime: new(time.Time),
			MPPotDuration: 3 * time.Second,
			HealNextTime:  &healNextTime,
			HealDuration:  10 * time.Second,
		},
	}
}

// IsAttackTarget 判断玩家是否是攻击者的攻击对象
func (p *Character) IsAttackTarget(attacker IMapObject) bool {
	return false
}

func (p *Character) IsFriendlyTarget(attacker IMapObject) bool {
	return true
}

func (p *Character) GetInfo() interface{} {
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

func (p *Character) GetBaseStats() BaseStats {
	return BaseStats{
		MinAC:    p.MinAC,
		MaxAC:    p.MaxAC,
		MinMAC:   p.MinMAC,
		MaxMAC:   p.MaxMAC,
		MinDC:    p.MinDC,
		MaxDC:    p.MaxDC,
		MinMC:    p.MinMC,
		MaxMC:    p.MaxMC,
		MinSC:    p.MinSC,
		MaxSC:    p.MaxSC,
		Accuracy: p.Accuracy,
		Agility:  p.Agility,
	}
}

func (p *Character) GetID() uint32 {
	return p.ID
}

func (p *Character) Point() common.Point {
	return p.GetPoint()
}

func (p *Character) GetRace() common.ObjectType {
	return common.ObjectTypePlayer
}

func (p *Character) GetPoint() common.Point {
	return p.CurrentLocation
}

func (p *Character) GetCell() *Cell {
	return p.Map.GetCell(p.CurrentLocation)
}

func (p *Character) GetDirection() common.MirDirection {
	return p.CurrentDirection
}

func (p *Character) Enqueue(msg interface{}) {
	p.Player.Enqueue(msg)
}

func (p *Character) ReceiveChat(text string, ct common.ChatType) {
	p.Player.ReceiveChat(text, ct)
}

func (p *Character) GetCurrentGrid() *Grid {
	return p.Map.AOI.GetGridByPoint(p.Point())
}

func (p *Character) BroadcastDamageIndicator(typ common.DamageType, dmg int) {
	msg := ServerMessage{}.DamageIndicator(int32(dmg), typ, p.GetID())
	p.Enqueue(msg)
	p.Broadcast(msg)
}

func (p *Character) Broadcast(msg interface{}) {
	p.Map.Submit(NewTask(func(args ...interface{}) {
		grids := p.Map.AOI.GetSurroundGrids(p.Point())
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

func (c *Character) Process() {
	finishID := make([]uint32, 0)
	now := time.Now()
	c.ActionList.Range(func(k, v interface{}) bool {
		action := v.(*DelayedAction)
		if action.Finish || now.Before(action.ActionTime) {
			return true
		}
		action.Task.Execute()
		action.Finish = true
		if action.Finish {
			finishID = append(finishID, action.ID)
		}
		return true
	})
	for i := range finishID {
		c.ActionList.Delete(finishID[i])
	}
	ch := &c.Health
	if ch.HPPotValue != 0 && ch.HPPotNextTime.Before(now) {
		c.ChangeHP(ch.HPPotPerValue)
		ch.HPPotTickTime += 1
		if ch.HPPotTickTime >= ch.HPPotTickNum {
			ch.HPPotValue = 0
		} else {
			*ch.HPPotNextTime = now.Add(ch.HPPotDuration)
		}
	}
	if ch.MPPotValue != 0 && ch.MPPotNextTime.Before(now) {
		c.ChangeMP(ch.MPPotPerValue)
		ch.MPPotTickTime += 1
		if ch.MPPotTickTime >= ch.MPPotTickNum {
			ch.MPPotValue = 0
		} else {
			*ch.MPPotNextTime = now.Add(ch.MPPotDuration)
		}
	}
	if ch.HealNextTime.Before(now) {
		*ch.HealNextTime = now.Add(ch.HealDuration)
		c.ChangeHP(int(float32(c.MaxHP)*0.03) + 1)
		c.ChangeMP(int(float32(c.MaxMP)*0.03) + 1)
	}
}

func (c *Character) CompleteAttack(args ...interface{})          {}
func (c *Character) CompleteMapMovement(args ...interface{})     {}
func (c *Character) CompleteMine(args ...interface{})            {}
func (c *Character) CompleteNPC(args ...interface{})             {}
func (c *Character) CompletePoison(args ...interface{})          {}
func (c *Character) CompleteDamageIndicator(args ...interface{}) {}

func (c *Character) NewObjectID() uint32 {
	return c.Map.Env.NewObjectID()
}

func (c *Character) IsDead() bool {
	return false
}

func (c *Character) IsHidden() bool {
	return false
}

func (c *Character) CanMove() bool {
	return true
}

func (c *Character) CanWalk() bool {
	return true
}

func (c *Character) CanRun() bool {
	return true
}

func (c *Character) CanAttack() bool {
	return true
}

func (c *Character) CanRegen() bool {
	return true
}

func (c *Character) CanCast() bool {
	return true
}

func (c *Character) CanUseItem(item *common.UserItem) bool {
	return true
}

func (c *Character) EnqueueItemInfos() {
	gdb := c.Map.Env.GameDB
	itemInfos := make([]*common.ItemInfo, 0)
	for i := range c.Inventory {
		itemID := int(c.Inventory[i].ItemID)
		if itemID == 0 {
			continue
		}
		itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
	}
	for i := range c.Equipment {
		itemID := int(c.Equipment[i].ItemID)
		if itemID == 0 {
			continue
		}
		itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
	}
	for i := range c.QuestInventory {
		itemID := int(c.QuestInventory[i].ItemID)
		if itemID == 0 {
			continue
		}
		itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
	}
	for i := range itemInfos {
		c.EnqueueItemInfo(itemInfos[i].ID)
	}
}

func (c *Character) EnqueueItemInfo(itemID int32) {
	for m := range c.SendItemInfo {
		s := c.SendItemInfo[m]
		if s.ID == itemID {
			return
		}
	}
	item := c.Map.Env.GameDB.GetItemInfoByID(int(itemID))
	if item == nil {
		return
	}
	c.Player.Enqueue(ServerMessage{}.NewItemInfo(item))
	c.SendItemInfo = append(c.SendItemInfo, *item)
}

func (c *Character) EnqueueQuestInfo() {

}

func (c *Character) RefreshStats() {
	c.RefreshLevelStats()
	c.RefreshBagWeight()
	c.RefreshEquipmentStats()
	c.RefreshItemSetStats()
	c.RefreshMirSetStats()
	c.RefreshSkills()
	c.RefreshBuffs()
	c.RefreshStatCaps()
	c.RefreshMountStats()
	c.RefreshGuildBuffs()
}

func (c *Character) RefreshLevelStats() {
	baseStats := setting.BaseStats[c.Class]
	c.Accuracy = uint8(baseStats.StartAccuracy)
	c.Agility = uint8(baseStats.StartAgility)
	c.CriticalRate = uint8(baseStats.StartCriticalRate)
	c.CriticalDamage = uint8(baseStats.StartCriticalDamage)
	c.MaxExperience = 100
	c.MaxHP = uint16(14 + (float32(c.Level)/baseStats.HpGain+baseStats.HpGainRate)*float32(c.Level))
	c.MinAC = 0
	if baseStats.MinAc > 0 {
		c.MinAC = uint16(int(c.Level) / baseStats.MinAc)
	}
	c.MaxAC = 0
	if baseStats.MaxAc > 0 {
		c.MaxAC = uint16(int(c.Level) / baseStats.MaxAc)
	}
	c.MinMAC = 0
	if baseStats.MinMac > 0 {
		c.MinMAC = uint16(int(c.Level) / baseStats.MinMac)
	}
	c.MaxMAC = 0
	if baseStats.MaxMac > 0 {
		c.MaxMAC = uint16(int(c.Level) / baseStats.MaxMac)
	}
	c.MinDC = 0
	if baseStats.MinDc > 0 {

		c.MinDC = uint16(int(c.Level) / baseStats.MinDc)
	}
	c.MaxDC = 0
	if baseStats.MaxDc > 0 {
		c.MaxDC = uint16(int(c.Level) / baseStats.MaxDc)
	}
	c.MinMC = 0
	if baseStats.MinMc > 0 {
		c.MinMC = uint16(int(c.Level) / baseStats.MinMc)
	}
	c.MaxMC = 0
	if baseStats.MaxMc > 0 {
		c.MaxMC = uint16(int(c.Level) / baseStats.MaxMc)
	}
	c.MinSC = 0
	if baseStats.MinSc > 0 {
		c.MinSC = uint16(int(c.Level) / baseStats.MinSc)
	}
	c.MaxSC = 0
	if baseStats.MaxSc > 0 {
		c.MaxSC = uint16(int(c.Level) / baseStats.MaxSc)
	}
	c.CriticalRate = 0
	if baseStats.CritialRateGain > 0 {
		c.CriticalRate = uint8(float32(c.CriticalRate) + (float32(c.Level) / baseStats.CritialRateGain))
	}
	c.CriticalDamage = 0
	if baseStats.CriticalDamageGain > 0 {
		c.CriticalDamage = uint8(float32(c.CriticalDamage) + (float32(c.Level) / baseStats.CriticalDamageGain))
	}
	c.MaxBagWeight = uint16(50.0 + float32(c.Level)/baseStats.BagWeightGain*float32(c.Level))
	c.MaxWearWeight = uint16(15.0 + float32(c.Level)/baseStats.WearWeightGain*float32(c.Level))
	c.MaxHandWeight = uint16(12.0 + float32(c.Level)/baseStats.HandWeightGain*float32(c.Level))
	switch c.Class {
	case common.MirClassWarrior:
		c.MaxHP = uint16(14.0 + (float32(c.Level)/baseStats.HpGain+baseStats.HpGainRate+float32(c.Level)/20.0)*float32(c.Level))
		c.MaxMP = uint16(11.0 + (float32(c.Level) * 3.5) + (float32(c.Level) * baseStats.MpGainRate))
	case common.MirClassWizard:
		c.MaxMP = uint16(13.0 + (float32(c.Level/5.0+2.0) * 2.2 * float32(c.Level)) + (float32(c.Level) * baseStats.MpGainRate))
	case common.MirClassTaoist:
		c.MaxMP = uint16((13 + float32(c.Level)/8.0*2.2*float32(c.Level)) + (float32(c.Level) * baseStats.MpGainRate))
	}
}

func (c *Character) RefreshBagWeight() {
	c.CurrentBagWeight = 0
	for i := range c.Inventory {
		ui := c.Inventory[i]
		if ui.ID != 0 {
			it := c.Map.Env.GameDB.GetItemInfoByID(int(ui.ItemID))
			c.CurrentBagWeight += int(it.Weight)
		}
	}
}

func (c *Character) RefreshEquipmentStats() {
	gdb := c.Map.Env.GameDB
	for i := range c.Equipment {
		e := gdb.GetItemInfoByID(int(c.Equipment[i].ItemID))
		if e == nil {
			continue
		}
		switch e.Type {
		case common.ItemTypeArmour:
			c.LooksArmour = int(e.Shape)
			c.LooksWings = int(e.Effect)
		case common.ItemTypeWeapon:
			c.LooksWeapon = int(e.Shape)
			c.LooksWeaponEffect = int(e.Effect)
		}
	}
}

func (c *Character) RefreshItemSetStats() {

}

func (c *Character) RefreshMirSetStats() {

}

func (c *Character) RefreshSkills() {

}

func (c *Character) RefreshBuffs() {

}

func (c *Character) RefreshStatCaps() {

}

func (c *Character) RefreshMountStats() {

}

func (c *Character) RefreshGuildBuffs() {

}

// GetUserItemByID 获取物品，返回该物品在容器的索引和是否成功
func (c *Character) GetUserItemByID(mirGridType common.MirGridType, id uint64) (index int, item *common.UserItem) {
	var arr []common.UserItem
	switch mirGridType {
	case common.MirGridTypeInventory:
		arr = c.Inventory
	case common.MirGridTypeEquipment:
		arr = c.Equipment
	default:
		panic("error mirGridType")
	}
	for i := range arr {
		item := arr[i]
		if item.ID == id {
			return i, &item
		}
	}
	return -1, nil
}

// GainItem 为玩家增加物品，增加成功返回 true
func (c *Character) GainItem(ui *common.UserItem) bool {
	item := c.Map.Env.GameDB.GetItemInfoByID(int(ui.ItemID))
	if item == nil {
		return false
	}
	i, j := 6, 46
	if item.Type == common.ItemTypePotion ||
		item.Type == common.ItemTypeScroll ||
		item.Type == common.ItemTypeScript ||
		item.Type == common.ItemTypeAmulet {
		i = 0
		j = 4
	} else if item.Type == common.ItemTypeAmulet {
		i = 4
		j = 6
	} else {
		i = 6
	}
	for i < j {
		if c.Inventory[i].ID != 0 {
			i++
			continue
		}
		c.Inventory[i] = *ui
		break
	}
	c.EnqueueItemInfo(ui.ItemID)
	c.Player.Enqueue(ServerMessage{}.GainedItem(ui))
	c.RefreshBagWeight()
	return true
}

// GainGold 为玩家增加金币
func (c *Character) GainGold(gold uint64) {
	if gold <= 0 {
		return
	}
	c.Gold += gold
	c.Player.Enqueue(ServerMessage{}.GainedGold(gold))
}

func (c *Character) UpdateConcentration() {
	c.Player.Enqueue(ServerMessage{}.SetConcentration(c.Player))
	c.Broadcast(ServerMessage{}.SetObjectConcentration(c.Player))
}

func (c *Character) GetAttackPower(min, max int) int {
	if min < 0 {
		min = 0
	}
	if max < min {
		max = min
	}
	// TODO luck
	return G_Rand.RandInt(min, max+1)
}

// TODO
func (c *Character) Attacked(attacker IMapObject, damageFinal int, defenceType common.DefenceType, damageWeapon bool) {

}

func (c *Character) GainExp(amount uint32) {
	c.Experience += int64(amount)
	c.Player.Enqueue(ServerMessage{}.GainExperience(amount))
	if c.Experience < c.MaxExperience {
		return
	}
	c.Experience -= c.MaxExperience
	c.Level++
	c.LevelUp()
}

func (c *Character) SetHP(amount uint32) {
	c.HP = uint16(amount)
	msg := ServerMessage{}.HealthChanged(c.HP, c.MP)
	c.Player.Enqueue(msg)
	c.Broadcast(msg)
}

func (c *Character) SetMP(amount uint32) {
	c.MP = uint16(amount)
	msg := ServerMessage{}.HealthChanged(c.HP, c.MP)
	c.Player.Enqueue(msg)
	c.Broadcast(msg)
}

func (c *Character) ChangeHP(amount int) {
	if amount == 0 || c.IsDead() || c.HP >= c.MaxHP {
		return
	}
	c.SetHP(uint32(int(c.HP) + amount))
}

func (c *Character) ChangeMP(amount int) {
	if amount == 0 || c.IsDead() || c.MP >= c.MaxMP {
		return
	}
	c.SetMP(uint32(int(c.MP) + amount))
}

func (c *Character) LevelUp() {
	c.RefreshStats()
	c.SetHP(uint32(c.MaxHP))
	c.SetMP(uint32(c.MaxMP))
	c.Player.Enqueue(ServerMessage{}.LevelChanged(c.Level, c.Experience, c.MaxExperience))
	c.Broadcast(ServerMessage{}.ObjectLeveled(c.GetID()))
}

func (c *Character) Die() {

}

func (c *Character) Teleport(m *Map, p common.Point) {

}

func (p *Character) Magic(spell common.Spell, direction common.MirDirection, targetID uint32, targetLocation common.Point) {
	if !p.CanCast() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	userMagic := p.GetMagic(spell)
	if userMagic == nil {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	info := p.Map.Env.GameDB.GetMagicInfoByID(userMagic.MagicID)
	cost := info.BaseCost + info.LevelCost*userMagic.Level
	if uint16(cost) > p.MP {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.ChangeMP(-cost)
	target := p.Map.GetObjectInAreaByID(targetID, targetLocation)
	cast, targetID := p.UseMagic(spell, userMagic, target)
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Enqueue(ServerMessage{}.Magic(spell, targetID, targetLocation, cast, userMagic.Level))
	p.Broadcast(ServerMessage{}.ObjectMagic(p, spell, targetID, targetLocation, cast, userMagic.Level))
}

func (p *Character) Attack(direction common.MirDirection, spell common.Spell) {
	if !p.CanAttack() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectAttack(p, common.SpellNone, 0, 0))
	target := p.GetPoint().NextPoint(p.GetDirection(), 1)
	damageBase := p.GetAttackPower(int(p.MinDC), int(p.MaxDC)) // = the original damage from your gear (+ bonus from moonlight and darkbody)
	damageFinal := damageBase                                  // = the damage you're gonna do with skills added
	cell := p.Map.GetCell(target)
	if !cell.CanWalk() {
		return
	}
	cell.Objects.Range(func(k, v interface{}) bool {
		o := v.(IMapObject)
		if !o.IsAttackTarget(p) {
			return true
		}
		switch o.GetRace() {
		case common.ObjectTypePlayer:
			o.(*Character).Attacked(p, damageFinal, common.DefenceTypeAgility, false)
		case common.ObjectTypeMonster:
			o.(*Monster).Attacked(p, damageFinal, common.DefenceTypeAgility, false)
		}
		return true
	})
}

func (p *Character) UseItem(id uint64) {
	msg := &server.UseItem{UniqueID: id, Success: false}
	if p.IsDead() {
		p.Enqueue(msg)
		return
	}
	index, item := p.GetUserItemByID(common.MirGridTypeInventory, id)
	if item == nil || item.ID == 0 || !p.CanUseItem(item) {
		p.Enqueue(msg)
		return
	}
	ph := &p.Health
	info := p.Map.Env.GameDB.GetItemInfoByID(int(item.ItemID))
	switch info.Type {
	case common.ItemTypePotion:
		switch info.Shape {
		case 0: // NormalPotion 一般药水
			if info.HP > 0 {
				ph.HPPotValue = int(info.HP)                         // 回复总值
				ph.HPPotPerValue = int(info.HP / 3)                  // 一次回复多少
				*ph.HPPotNextTime = time.Now().Add(ph.HPPotDuration) // 下次生效时间
				ph.HPPotTickNum = 3                                  // 总共跳几次
				ph.HPPotTickTime = 0                                 // 当前第几跳
			}
			if info.MP > 0 {
				ph.MPPotValue = int(info.MP)
				ph.MPPotPerValue = int(info.MP / 3)
				*ph.MPPotNextTime = time.Now().Add(ph.MPPotDuration)
				ph.MPPotTickNum = 3
				ph.MPPotTickTime = 0
			}
		case 1: // SunPotion 太阳水
			p.ChangeHP(int(info.HP))
			p.ChangeMP(int(info.MP))
		case 2: // TODO MysteryWater
		case 3: // TODO Buff
		case 4: // TODO Exp 经验
		}
	case common.ItemTypeScroll:
	case common.ItemTypeBook:
	case common.ItemTypeScript:
	case common.ItemTypeFood:
	case common.ItemTypePets:
	case common.ItemTypeTransform: //Transforms
	default:
		p.Enqueue(msg)
		return
	}
	if item.Count > 1 {
		item.Count--
	} else {
		p.Inventory[index] = common.UserItem{}
	}
	p.RefreshBagWeight()
	msg.Success = true
	p.Enqueue(msg)
}

func (p *Character) DropItem(id uint64, count uint32) {
	msg := &server.DropItem{
		UniqueID: id,
		Count:    count,
		Success:  false,
	}
	index, userItem := p.GetUserItemByID(common.MirGridTypeInventory, id)
	if userItem == nil || userItem.ID == 0 {
		p.Enqueue(msg)
		return
	}
	obj := p.Map.Env.CreateDropItem(p.Map, userItem, 0)
	if dropMsg, ok := obj.Drop(p.GetPoint(), 1); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}
	if count >= userItem.Count {
		p.Inventory[index] = common.UserItem{}
	} else {
		p.Inventory[index].Count -= count
	}
	p.RefreshBagWeight()
	msg.Success = true
	p.Enqueue(msg)
}

func (p *Character) DropGold(gold uint64) {
	if p.Gold < gold {
		return
	}
	obj := p.Map.Env.CreateDropItem(p.Map, nil, gold)
	if dropMsg, ok := obj.Drop(p.GetPoint(), 3); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}
	p.Gold -= gold
	p.Enqueue(&server.LoseGold{Gold: uint32(gold)})
}

func (p *Character) PickUp() {
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

func (p *Character) Inspect(id uint32) {
	o := p.Map.Env.GetPlayer(id)
	for i := range o.Equipment {
		item := p.Map.Env.GameDB.GetItemInfoByID(int(o.Equipment[i].ItemID))
		if item != nil {
			p.EnqueueItemInfo(item.ID)
		}
	}
	p.Enqueue(ServerMessage{}.PlayerInspect(o))
}

func (p *Character) Walk(direction common.MirDirection) {
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

func (p *Character) Run(direction common.MirDirection) {
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

func (p *Character) EquipItem(mirGridType common.MirGridType, id uint64, to int32) {
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
	case common.MirGridTypeStorage:
		// TODO
	}
	msg.Success = true
	p.RefreshStats()
	p.Enqueue(msg)
	p.UpdateConcentration()
	p.Broadcast(ServerMessage{}.PlayerUpdate(p))
}

func (p *Character) RemoveItem(mirGridType common.MirGridType, id uint64, to int32) {
	msg := &server.RemoveItem{
		Grid:     mirGridType,
		UniqueID: id,
		To:       to,
		Success:  false,
	}
	if l := len(p.Inventory); to < 0 || int(to) >= l {
		p.Enqueue(msg)
		return
	}
	switch mirGridType {
	case common.MirGridTypeInventory:
		index, item := p.GetUserItemByID(common.MirGridTypeEquipment, id)
		if item == nil {
			p.Enqueue(msg)
			return
		}
		invItem := p.Inventory[to]
		if invItem.ID == 0 {
			p.Inventory[to], p.Equipment[index] = p.Equipment[index], p.Inventory[to]
			break
		}
		for i := range p.Inventory[6:] {
			tmp := p.Inventory[6:][i]
			if tmp.ID != 0 {
				continue
			}
			p.Inventory[6:][i], p.Equipment[index] = p.Equipment[index], p.Inventory[6:][i]
			break
		}
	case common.MirGridTypeStorage:
		// TODO
	}
	msg.Success = true
	p.RefreshStats()
	p.Enqueue(msg)
	p.UpdateConcentration()
	p.Broadcast(ServerMessage{}.PlayerUpdate(p))
}

func (p *Character) MoveItem(mirGridType common.MirGridType, from int32, to int32) {
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

func (p *Character) Turn(direction common.MirDirection) {
	if !p.CanMove() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectTurn(p))
}

func (p *Character) EnqueueAreaObjects(oldGrid, newGrid *Grid) {
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

// TODO
func (p *Character) StoreItem(from int32, to int32) {
	msg := &server.StoreItem{
		From:    from,
		To:      to,
		Success: false,
	}
	p.Enqueue(msg)
}
