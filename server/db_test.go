package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yenkeia/mirgo/com"
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var basic com.Basic
	db.Table("basic").Where("game_version = ?", 65).Find(&basic)
	t.Log(basic.GameVersion, basic.MapIndex, basic.RespawnIndex)

	var gameShopItem com.GameShopItem
	db.Table("game_shop_item").Where("game_shop_item_index = ?", 2).Find(&gameShopItem)
	t.Log(gameShopItem.GoldPrice)

	var magicInfo com.MagicInfo
	db.Table("magic_info").Where("name = ?", "Fencing").Find(&magicInfo)
	t.Log(magicInfo.Name, magicInfo.Icon)

	var mapInfo com.MapInfo
	db.Table("map_info").Where("map_index = ?", 1).Find(&mapInfo)
	t.Log(mapInfo.Title)

	//var mineZone com.MineZone
	//db.Table("mine_zone").Where("map_index = ?", )

	var monsterInfo com.MonsterInfo
	db.Table("monster_info").Where("monster_index = ?", 1).Find(&monsterInfo)
	t.Log(monsterInfo.Name)

	var movementInfo com.MovementInfo
	db.Table("movement_info").Where("map_index = ?", 2).First(&movementInfo)
	t.Log(movementInfo.MapIndex, movementInfo.ConquestIndex, movementInfo.DestinationX, movementInfo.DestinationY)

	var npcInfo com.NpcInfo
	db.Table("npc_info").Where("npc_index = ?", 1).Find(&npcInfo)
	t.Log(npcInfo.Filename)

	var questInfo com.QuestInfo
	db.Table("quest_info").Where("quest_index = ?", 1).Find(&questInfo)
	t.Log(questInfo.Name)

	var respawnInfo com.RespawnInfo
	db.Table("respawn_info").Where("location_x = ?", 350).Find(&respawnInfo)
	t.Log(respawnInfo.MapIndex, respawnInfo.RespawnIndex, respawnInfo.Count)

	var safeZoneInfo com.SafeZoneInfo
	db.Table("safe_zone_info").Where("map_index = ?", 1).Find(&safeZoneInfo)
	t.Log(safeZoneInfo.MapIndex, safeZoneInfo.LocationX, safeZoneInfo.LocationY)
}
