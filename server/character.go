package main

import "github.com/yenkeia/mirgo/common"

type Character struct {
	Player            *Player
	HP                uint16
	MP                uint16
	Experience        int64
	Gold              uint64
	GuildName         string
	GuildRankName     string
	Class             common.MirClass
	Gender            common.MirGender
	Hair              uint8
	Inventory         []common.UserItem // 46
	Equipment         []common.UserItem // 14
	QuestInventory    []common.UserItem // 40
	Trade             []common.UserItem // 10
	Refine            []common.UserItem // 16
	LooksArmour       int
	LooksWings        int
	LooksWeapon       int
	LooksWeaponEffect int
	SendItemInfo      []common.ItemInfo
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

// GainGold 为玩家增加金币，增加成功返回 true
func (c *Character) GainGold(gold uint64) bool {
	return true
}
