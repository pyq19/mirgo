package main

import "github.com/yenkeia/mirgo/common"

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
}

func NewCharacter(g *Game, p *Player, c *common.Character) Character {
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
	for i, v := range uii {
		inventory[i] = v
	}
	for i, v := range uie {
		equipment[i] = v
	}
	for i, v := range uiq {
		questInventory[i] = v
	}
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

}

func (c *Character) RefreshBagWeight() {

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
