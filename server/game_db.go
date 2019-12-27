package main

import "github.com/yenkeia/mirgo/common"

var G_GameDB *GameDB

// GameDB ...
type GameDB struct {
	Basic         common.Basic
	GameShopItems []common.GameShopItem
	ItemInfos     []common.ItemInfo
	MagicInfos    []common.MagicInfo
	MapInfos      []common.MapInfo
	MonsterInfos  []common.MonsterInfo
	MovementInfos []common.MovementInfo
	NpcInfos      []common.NpcInfo
	QuestInfos    []common.QuestInfo
	RespawnInfos  []common.RespawnInfo
	SafeZoneInfos []common.SafeZoneInfo
}

// GetMapInfoByID FIXME 改成从 map 取出
func (db *GameDB) GetMapInfoByID(mapID int) *common.MapInfo {
	for _, v := range db.MapInfos {
		if v.ID == mapID {
			return &v
		}
	}
	return nil
}

// GetItemInfoByID FIXME 改成从 map 取出
func (db *GameDB) GetItemInfoByID(itemID int) *common.ItemInfo {
	for _, v := range db.ItemInfos {
		if v.ID == int32(itemID) {
			return &v
		}
	}
	return nil
}

// GetMonsterInfoByID FIXME 改成从 map 取出
func (db *GameDB) GetMonsterInfoByID(monsterID int) *common.MonsterInfo {
	for _, v := range db.MonsterInfos {
		if v.ID == monsterID {
			return &v
		}
	}
	return nil
}
