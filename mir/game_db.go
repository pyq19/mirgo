package mir

import (
	"github.com/jinzhu/gorm"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/setting"
)

// GameDB ...
type GameData struct {
	Basic              common.Basic
	GameShopItems      []*common.GameShopItem
	ItemInfos          []*common.ItemInfo
	MagicInfos         []*common.MagicInfo
	MapInfos           []*common.MapInfo
	MonsterInfos       []*common.MonsterInfo
	MovementInfos      []*common.MovementInfo
	NpcInfos           []*common.NpcInfo
	QuestInfos         []*common.QuestInfo
	RespawnInfos       []*common.RespawnInfo
	SafeZoneInfos      []*common.SafeZoneInfo
	UserMagics         []*common.UserMagic
	MapIDInfoMap       map[int]*common.MapInfo
	ItemIDInfoMap      map[int]*common.ItemInfo
	ItemNameInfoMap    map[string]*common.ItemInfo
	MonsterIDInfoMap   map[int]*common.MonsterInfo
	MonsterNameInfoMap map[string]*common.MonsterInfo
	DropInfoMap        map[string][]common.DropInfo
	MagicIDInfoMap     map[int]*common.MagicInfo
	ExpList            []int
}

var data = NewGameData()

func NewGameData() *GameData {
	d := &GameData{}
	d.MapIDInfoMap = map[int]*common.MapInfo{}
	d.ItemIDInfoMap = map[int]*common.ItemInfo{}
	d.ItemNameInfoMap = map[string]*common.ItemInfo{}
	d.MonsterIDInfoMap = map[int]*common.MonsterInfo{}
	d.MonsterNameInfoMap = map[string]*common.MonsterInfo{}
	d.DropInfoMap = map[string][]common.DropInfo{}
	d.MagicIDInfoMap = map[int]*common.MagicInfo{}
	return d
}

func (d *GameData) Load(db *gorm.DB) {
	db.Table("basic").First(&d.Basic)
	db.Table("game_shop_item").Find(&d.GameShopItems)
	db.Table("item").Find(&d.ItemInfos)
	db.Table("magic").Find(&d.MagicInfos)
	db.Table("map").Find(&d.MapInfos)
	db.Table("monster").Find(&d.MonsterInfos)
	db.Table("movement").Find(&d.MovementInfos)
	db.Table("npc").Find(&d.NpcInfos)
	db.Table("quest").Find(&d.QuestInfos)
	db.Table("respawn").Find(&d.RespawnInfos)
	db.Table("safe_zone").Find(&d.SafeZoneInfos)

	for _, v := range d.MapInfos {
		d.MapIDInfoMap[v.ID] = v
	}
	for _, v := range d.ItemInfos {
		d.ItemIDInfoMap[int(v.ID)] = v
		d.ItemNameInfoMap[v.Name] = v
	}
	for _, v := range d.MonsterInfos {
		d.MonsterIDInfoMap[v.ID] = v
		d.MonsterNameInfoMap[v.Name] = v
	}
	for _, v := range d.MagicInfos {
		d.MagicIDInfoMap[v.ID] = v
	}

	d.LoadMonsterDrop()
	d.LoadExpList()
}

func (d *GameData) LoadMonsterDrop() {
	itemMap := map[string]int32{}

	for i := range d.ItemInfos {
		v := d.ItemInfos[i]
		itemMap[v.Name] = v.ID
	}

	for i := range d.MonsterInfos {
		v := d.MonsterInfos[i]
		dropInfos, err := common.GetDropInfosByMonsterName(setting.Conf.DropDirPath, v.Name)
		if err != nil {
			log.Warnln("加载怪物掉落错误", v.Name, err.Error())
			continue
		}
		d.DropInfoMap[v.Name] = dropInfos
	}
}

func (d *GameData) LoadExpList() {

}

// GetMapInfoByID ...
func (db *GameData) GetMapInfoByID(mapID int) *common.MapInfo {
	v, ok := db.MapIDInfoMap[mapID]
	if !ok {
		return nil
	}
	return v
}

// GetItemInfoByID ...
func (db *GameData) GetItemInfoByID(itemID int) *common.ItemInfo {
	v, ok := db.ItemIDInfoMap[itemID]
	if !ok {
		return nil
	}
	return v
}

// GetItemInfoByName ...
func (db *GameData) GetItemInfoByName(itemName string) *common.ItemInfo {
	v, ok := db.ItemNameInfoMap[itemName]
	if !ok {
		return nil
	}
	return v
}

// GetMonsterInfoByID ...
func (db *GameData) GetMonsterInfoByID(monsterID int) *common.MonsterInfo {
	v, ok := db.MonsterIDInfoMap[monsterID]
	if !ok {
		return nil
	}
	return v
}

// GetMonsterInfoByName ...
func (db *GameData) GetMonsterInfoByName(monsterName string) *common.MonsterInfo {
	v, ok := db.MonsterNameInfoMap[monsterName]
	if !ok {
		return nil
	}
	return v
}

// GetMagicInfoByID ...
func (db *GameData) GetMagicInfoByID(magicID int) *common.MagicInfo {
	v, ok := db.MagicIDInfoMap[magicID]
	if !ok {
		return nil
	}
	return v
}
