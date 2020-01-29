package main

import "github.com/yenkeia/mirgo/common"

import "github.com/yenkeia/mirgo/setting"

type Character struct {
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
}

func NewCharacter(g *Game, p *Player, c *common.Character) Character {
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
	return Character{
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
	}
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
	gdb := c.Player.Map.Env.GameDB
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
	item := c.Player.Map.Env.GameDB.GetItemInfoByID(int(itemID))
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
	c.MinAC = uint16(int(c.Level) / baseStats.MinAc)
	c.MaxAC = uint16(int(c.Level) / baseStats.MaxAc)
	c.MinMAC = uint16(int(c.Level) / baseStats.MinMac)
	c.MaxMAC = uint16(int(c.Level) / baseStats.MaxMac)
	c.MinDC = uint16(int(c.Level) / baseStats.MinDc)
	c.MaxDC = uint16(int(c.Level) / baseStats.MaxDc)
	c.MinMC = uint16(int(c.Level) / baseStats.MinMc)
	c.MaxMC = uint16(int(c.Level) / baseStats.MaxMc)
	c.MinSC = uint16(int(c.Level) / baseStats.MinSc)
	c.MaxSC = uint16(int(c.Level) / baseStats.MaxSc)
	c.CriticalRate = uint8(float32(c.CriticalRate) + (float32(c.Level) / baseStats.CritialRateGain))
	c.CriticalDamage = uint8(float32(c.CriticalDamage) + (float32(c.Level) / baseStats.CriticalDamageGain))
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
			it := c.Player.Map.Env.GameDB.GetItemInfoByID(int(ui.ItemID))
			c.CurrentBagWeight += int(it.Weight)
		}
	}
}

func (c *Character) RefreshEquipmentStats() {
	gdb := c.Player.Map.Env.GameDB
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
	item := c.Player.Map.Env.GameDB.GetItemInfoByID(int(ui.ItemID))
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
	c.Player.Broadcast(ServerMessage{}.SetObjectConcentration(c.Player))
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
	c.Player.Broadcast(msg)
}

func (c *Character) SetMP(amount uint32) {
	c.MP = uint16(amount)
	msg := ServerMessage{}.HealthChanged(c.HP, c.MP)
	c.Player.Enqueue(msg)
	c.Player.Broadcast(msg)
}

func (c *Character) ChangeHP(amount int) {
	if amount == 0 || c.IsDead() {
		return
	}
	c.SetHP(uint32(int(c.HP) + amount))
}

func (c *Character) ChangeMP(amount int) {
	if amount == 0 || c.IsDead() {
		return
	}
	c.SetMP(uint32(int(c.MP) + amount))
}

func (c *Character) LevelUp() {
	c.RefreshStats()
	c.SetHP(uint32(c.MaxHP))
	c.SetMP(uint32(c.MaxMP))
	c.Player.Enqueue(ServerMessage{}.LevelChanged(c.Level, c.Experience, c.MaxExperience))
	c.Player.Broadcast(ServerMessage{}.ObjectLeveled(c.Player.GetID()))
}

func (c *Character) GetMagic(spell common.Spell) *common.UserMagic {
	return nil
}
