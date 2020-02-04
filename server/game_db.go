package main

import (
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// GameDB ...
type GameDB struct {
	Basic              common.Basic
	GameShopItems      []common.GameShopItem
	ItemInfos          []common.ItemInfo
	MagicInfos         []common.MagicInfo
	MapInfos           []common.MapInfo
	MonsterInfos       []common.MonsterInfo
	MovementInfos      []common.MovementInfo
	NpcInfos           []common.NpcInfo
	QuestInfos         []common.QuestInfo
	RespawnInfos       []common.RespawnInfo
	SafeZoneInfos      []common.SafeZoneInfo
	UserMagics         []common.UserMagic
	MapIDInfoMap       *sync.Map // key: MapID, value: MapInfo
	ItemIDInfoMap      *sync.Map // key: ItemID, value: ItemInfo
	ItemNameInfoMap    *sync.Map // key: ItemName, value: ItemInfo
	MonsterIDInfoMap   *sync.Map // key: MonsterID, value: MonsterInfo
	MonsterNameInfoMap *sync.Map // key: Monster.Name value: MonsterInfo
	DropInfoMap        *sync.Map // key: MonsterName, value: []common.DropInfo
	MagicIDInfoMap     *sync.Map // key: MagicInfo.ID, value: MagicInfo
}

// GetMapInfoByID
func (db *GameDB) GetMapInfoByID(mapID int) *common.MapInfo {
	v, ok := db.MapIDInfoMap.Load(mapID)
	if !ok {
		return nil
	}
	return v.(*common.MapInfo)
}

// GetItemInfoByID
func (db *GameDB) GetItemInfoByID(itemID int) *common.ItemInfo {
	v, ok := db.ItemIDInfoMap.Load(itemID)
	if !ok {
		return nil
	}
	return v.(*common.ItemInfo)
}

// GetItemInfoByName
func (db *GameDB) GetItemInfoByName(itemName string) *common.ItemInfo {
	v, ok := db.ItemNameInfoMap.Load(itemName)
	if !ok {
		return nil
	}
	return v.(*common.ItemInfo)
}

// GetMonsterInfoByID
func (db *GameDB) GetMonsterInfoByID(monsterID int) *common.MonsterInfo {
	v, ok := db.MonsterIDInfoMap.Load(monsterID)
	if !ok {
		return nil
	}
	return v.(*common.MonsterInfo)
}

// GetMonsterInfoByName
func (db *GameDB) GetMonsterInfoByName(monsterName string) *common.MonsterInfo {
	v, ok := db.MonsterNameInfoMap.Load(monsterName)
	if !ok {
		return nil
	}
	return v.(*common.MonsterInfo)
}

// GetMagicInfoByID
func (db *GameDB) GetMagicInfoByID(magicID int) *common.MagicInfo {
	v, ok := db.MagicIDInfoMap.Load(magicID)
	if !ok {
		return nil
	}
	return v.(*common.MagicInfo)
}
