package common

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestDB(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var basic Basic
	db.Table("basic").Where("game_version = ?", 65).Find(&basic)
	t.Log(basic.GameVersion, basic.MapIndex, basic.RespawnIndex)

	var gameShopItem GameShopItem
	db.Table("game_shop_item").Where("id = ?", 2).Find(&gameShopItem)
	t.Log(gameShopItem.GoldPrice, gameShopItem.ID)

	var magicInfo MagicInfo
	db.Table("magic").Where("name = ?", "Fencing").Find(&magicInfo)
	t.Log(magicInfo.Name, magicInfo.Icon)

	var mapInfo MapInfo
	db.Table("map").Where("id = ?", 1).Find(&mapInfo)
	t.Log(mapInfo.Title)

	//var mineZone com.MineZone
	//db.Table("mine_zone").Where("map_index = ?", )

	var monsterInfo MonsterInfo
	db.Table("monster").Where("id = ?", 1).Find(&monsterInfo)
	t.Log(monsterInfo.Name)

	var movementInfo MovementInfo
	db.Table("movement").Where("map_id = ?", 2).First(&movementInfo)
	// t.Log(movementInfo.MapID, movementInfo.ConquestIndex, movementInfo.DestinationX, movementInfo.DestinationY)

	var npcInfo NpcInfo
	db.Table("npc").Where("id = ?", 1).Find(&npcInfo)
	t.Log(npcInfo.Filename)

	var questInfo QuestInfo
	db.Table("quest").Where("id = ?", 1).Find(&questInfo)
	t.Log(questInfo.Name)

	var respawnInfo RespawnInfo
	db.Table("respawn").Where("location_x = ?", 350).Find(&respawnInfo)
	t.Log(respawnInfo.MapID, respawnInfo.RespawnIndex, respawnInfo.Count)

	var safeZoneInfo SafeZoneInfo
	db.Table("safe_zone").Where("map_id = ?", 1).Find(&safeZoneInfo)
	t.Log(safeZoneInfo.MapID, safeZoneInfo.LocationX, safeZoneInfo.LocationY)
}

func TestAccountDB(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var account Account
	db.Table("account").Where("username = ?", "不存在").Find(&account)
	t.Log(account)
}

func TestCharacter(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	ac := make([]AccountCharacter, 3)
	db.Table("account_character").Where("account_id = ?", 1).Limit(3).Find(&ac)
	t.Log(ac, len(ac))

	ids := make([]int, 3)
	for _, c := range ac {
		ids = append(ids, c.ID)
	}
	cs := make([]Character, 3)
	db.Table("character").Where("id in (?)", ids).Find(&cs)
	t.Log(cs)
}

func TestUserMagic(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var um UserMagic
	db.Table("user_magic").Where("id = ?", "不存在").Find(&um)
	t.Log(um)
}
