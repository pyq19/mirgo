package mir

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/ut"
)

// GameDB 服务端数据，启动时加载
type GameData struct {
	// Basic              common.Basic
	GameShopItems      []*common.GameShopItem
	ItemInfos          []*common.ItemInfo
	StartItems         []*common.ItemInfo // 新玩家初始装备
	MagicInfos         []*common.MagicInfo
	MapInfos           []*common.MapInfo
	MonsterInfos       []*common.MonsterInfo
	MovementInfos      []*common.MovementInfo
	NpcInfos           []*common.NpcInfo
	QuestInfos         []*common.QuestInfo
	RespawnInfos       []*common.RespawnInfo
	SafeZoneInfos      []*common.SafeZoneInfo
	StartPoints        []*common.SafeZoneInfo // 起始点
	MapIDInfoMap       map[int]*common.MapInfo
	ItemIDInfoMap      map[int]*common.ItemInfo
	ItemNameInfoMap    map[string]*common.ItemInfo
	MonsterIDInfoMap   map[int]*common.MonsterInfo
	MonsterNameInfoMap map[string]*common.MonsterInfo
	DropInfoMap        map[string][]*common.DropInfo
	MagicIDInfoMap     map[int]*common.MagicInfo
	ExpList            []int
}

var data *GameData

func NewGameData(db *gorm.DB) *GameData {
	d := &GameData{}
	d.MapIDInfoMap = map[int]*common.MapInfo{}
	d.ItemIDInfoMap = map[int]*common.ItemInfo{}
	d.ItemNameInfoMap = map[string]*common.ItemInfo{}
	d.MonsterIDInfoMap = map[int]*common.MonsterInfo{}
	d.MonsterNameInfoMap = map[string]*common.MonsterInfo{}
	d.DropInfoMap = map[string][]*common.DropInfo{}
	d.MagicIDInfoMap = map[int]*common.MagicInfo{}

	d.Load(db)

	return d
}

func (d *GameData) Load(db *gorm.DB) {
	// db.Table("basic").First(&d.Basic)
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

	d.StartItems = []*common.ItemInfo{}

	for _, v := range d.ItemInfos {
		if v.StartItem {
			d.StartItems = append(d.StartItems, v)
		}
	}

	d.StartPoints = []*common.SafeZoneInfo{}
	for _, v := range d.SafeZoneInfos {
		if v.StartPoint != 0 {
			d.StartPoints = append(d.StartPoints, v)
		}
	}

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

func (d *GameData) RandomStartPoint() *common.SafeZoneInfo {
	return d.StartPoints[ut.RandomNext(len(d.StartPoints))]
}

func (d *GameData) LoadMonsterDrop() {
	itemMap := map[string]int32{}

	for i := range d.ItemInfos {
		v := d.ItemInfos[i]
		itemMap[v.Name] = v.ID
	}

	for i := range d.MonsterInfos {
		v := d.MonsterInfos[i]
		dropInfos, err := d.loadDropFile(filepath.Join(settings.DropDirPath, v.Name+".txt"))
		if err != nil {
			log.Warnln(err.Error())
			continue
		}
		d.DropInfoMap[v.Name] = dropInfos
	}
}

func (d *GameData) loadDropFile(filename string) ([]*common.DropInfo, error) {
	lines, err := ut.ReadLines(filename)
	if err != nil {
		return nil, err
	}

	chanceReg := regexp.MustCompile(`(\d+)/(\d+)`)

	lineError := func(line int, detail string) error {
		return fmt.Errorf("DropInfo 格式不正确，%s行%d:%s %s", filename, line, lines[line], detail)
	}

	ret := []*common.DropInfo{}
	for i, line := range lines {

		line = strings.TrimSpace(line)
		line = ut.RemoveBOM(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		res := ut.SplitString(line)

		if len(res) != 3 && len(res) != 2 {
			log.Errorf("line: %s; line[0]: %s", line, line[0])
			return nil, lineError(i, "参数个数")
		}

		match := chanceReg.FindStringSubmatch(res[0])
		low, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, lineError(i, "分子错误")
		}

		high, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, lineError(i, "分母错误")
		}

		info := &common.DropInfo{Low: low, High: high, ItemName: res[1], Count: 1}

		if len(res) == 3 { // 1/10 Gold 500
			if strings.ToUpper(res[2]) == "Q" {
				info.QuestRequired = true
			} else {
				count, err := strconv.Atoi(res[2])
				info.Count = count
				if err != nil {
					for i, v := range res {
						fmt.Println(i, v)
					}
					return nil, lineError(i, "参数错误")
				}
			}
		}
		ret = append(ret, info)
	}

	return ret, nil
}

func (d *GameData) LoadExpList() {
	t, err := toml.LoadFile(filepath.Join(settings.ConfigsPath, "ExpList.ini"))
	if err != nil {
		panic("load explist error  " + err.Error())
	}

	t = t.Get("Exp").(*toml.Tree)

	d.ExpList = make([]int, 499)

	for i := 1; i < 500; i++ {
		d.ExpList[i-1] = int(t.Get(fmt.Sprintf("Level%d", i)).(int64))
	}
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

func (db *GameData) GetMagicInfoBySpell(spell common.Spell) *common.MagicInfo {
	for _, v := range db.MagicIDInfoMap {
		if v.Spell == int(spell) {
			return v
		}
	}
	return nil
}

func (db *GameData) GetMagicInfoByName(name string) *common.MagicInfo {
	for _, v := range db.MagicIDInfoMap {
		if v.Name == name {
			return v
		}
	}
	return nil
}
