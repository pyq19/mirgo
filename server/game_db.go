package main

import (
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// GameDB ...
type GameDB struct {
	Basic            common.Basic
	GameShopItems    []common.GameShopItem
	ItemInfos        []common.ItemInfo
	MagicInfos       []common.MagicInfo
	MapInfos         []common.MapInfo
	MonsterInfos     []common.MonsterInfo
	MovementInfos    []common.MovementInfo
	NpcInfos         []common.NpcInfo
	QuestInfos       []common.QuestInfo
	RespawnInfos     []common.RespawnInfo
	SafeZoneInfos    []common.SafeZoneInfo
	UserMagics       []common.UserMagic
	MapIDInfoMap     *sync.Map // key: MapID, value: MapInfo
	ItemIDInfoMap    *sync.Map // key: ItemID, value: ItemInfo
	MonsterIDInfoMap *sync.Map // key: MonsterID, value: MonsterInfo
}

func (db *GameDB) Init() {
	db.MapIDInfoMap = new(sync.Map)
	db.ItemIDInfoMap = new(sync.Map)
	db.MonsterIDInfoMap = new(sync.Map)
	for i := range db.MapInfos {
		v := db.MapInfos[i]
		db.MapIDInfoMap.Store(v.ID, &v)
	}
	for i := range db.ItemInfos {
		v := db.ItemInfos[i]
		db.ItemIDInfoMap.Store(int(v.ID), &v)
	}
	for i := range db.MonsterInfos {
		v := db.MonsterInfos[i]
		db.MonsterIDInfoMap.Store(v.ID, &v)
	}
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

// GetMonsterInfoByID
func (db *GameDB) GetMonsterInfoByID(monsterID int) *common.MonsterInfo {
	v, ok := db.MonsterIDInfoMap.Load(monsterID)
	if !ok {
		return nil
	}
	return v.(*common.MonsterInfo)
}
