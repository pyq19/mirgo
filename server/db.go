package main

import "github.com/yenkeia/mirgo/common"

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

// GetMapInfoById FIXME 改成从 map 取出
func (db *GameDB) GetMapInfoById(mapId int) *common.MapInfo {
	for _, v := range db.MapInfos {
		if v.Id == mapId {
			return &v
		}
	}
	return nil
}

// GetItemInfoById FIXME 改成从 map 取出
func (db *GameDB) GetItemInfoById(itemId int) *common.ItemInfo {
	for _, v := range db.ItemInfos {
		if v.Id == int32(itemId) {
			return &v
		}
	}
	return nil
}
