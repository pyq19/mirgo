package main

import "github.com/yenkeia/mirgo/common"

type Character struct {
	Player            *Player
	HP                uint16
	MP                uint16
	Experience        int64
	Gold              uint16
	GuildName         string
	GuildRankName     string
	Class             common.MirClass
	Gender            common.MirGender
	Hair              uint8
	Inventory         []common.UserItem
	Equipment         []common.UserItem
	QuestInventory    []common.UserItem
	Trade             []common.UserItem
	Refine            []common.UserItem
	LooksArmour       int
	LooksWings        int
	LooksWeapon       int
	LooksWeaponEffect int
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
		c.EnqueueItemInfo(itemInfos[i])
	}
}

func (c *Character) EnqueueItemInfo(i *common.ItemInfo) {
	// TODO
	//if (Connection.SentItemInfo.Contains(info)) return;
	//Enqueue(new S.NewItemInfo { Info = info });
	//Connection.SentItemInfo.Add(info)
	c.Player.Enqueue(ServerMessage{}.NewItemInfo(i))
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
