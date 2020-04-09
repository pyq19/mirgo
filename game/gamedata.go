package game

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/util"
)

// GameDB 服务端数据，启动时加载
type GameData struct {
	// Basic              cm.Basic
	GameShopItems      []*cm.GameShopItem
	ItemInfos          []*cm.ItemInfo
	StartItems         []*cm.ItemInfo // 新玩家初始装备
	MagicInfos         []*cm.MagicInfo
	MapInfos           []*cm.MapInfo
	MonsterInfos       []*cm.MonsterInfo
	MovementInfos      []*cm.MovementInfo
	NpcInfos           []*cm.NpcInfo
	QuestInfos         []*cm.QuestInfo
	RespawnInfos       []*cm.RespawnInfo
	SafeZoneInfos      []*cm.SafeZoneInfo
	StartPoints        []*cm.SafeZoneInfo // 起始点
	MapIDInfoMap       map[int]*cm.MapInfo
	ItemIDInfoMap      map[int]*cm.ItemInfo
	ItemNameInfoMap    map[string]*cm.ItemInfo
	MonsterIDInfoMap   map[int]*cm.MonsterInfo
	MonsterNameInfoMap map[string]*cm.MonsterInfo
	DropInfoMap        map[string][]*cm.DropInfo
	MagicIDInfoMap     map[int]*cm.MagicInfo
	ExpList            []int
}

var data *GameData

func NewGameData(db *gorm.DB) *GameData {
	d := &GameData{}
	d.MapIDInfoMap = map[int]*cm.MapInfo{}
	d.ItemIDInfoMap = map[int]*cm.ItemInfo{}
	d.ItemNameInfoMap = map[string]*cm.ItemInfo{}
	d.MonsterIDInfoMap = map[int]*cm.MonsterInfo{}
	d.MonsterNameInfoMap = map[string]*cm.MonsterInfo{}
	d.DropInfoMap = map[string][]*cm.DropInfo{}
	d.MagicIDInfoMap = map[int]*cm.MagicInfo{}

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

	d.StartItems = []*cm.ItemInfo{}

	for _, v := range d.ItemInfos {
		v := v
		// ClassBased = (bools & 0x04) == 0x04;
		b1 := v.Bools & 0x04
		if b1 == 0x04 {
			v.ClassBased = true
		}
		// LevelBased = (bools & 0x08) == 0x08;
		b2 := v.Bools & 0x08
		if b2 == 0x08 {
			v.LevelBased = true
		}
		if v.StartItem {
			d.StartItems = append(d.StartItems, v)
		}
	}

	d.StartPoints = []*cm.SafeZoneInfo{}
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

func (d *GameData) RandomStartPoint() *cm.SafeZoneInfo {
	return d.StartPoints[util.RandomNext(len(d.StartPoints))]
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

func (d *GameData) loadDropFile(filename string) ([]*cm.DropInfo, error) {
	lines, err := util.ReadLines(filename)
	if err != nil {
		return nil, err
	}

	chanceReg := regexp.MustCompile(`(\d+)/(\d+)`)

	lineError := func(line int, detail string) error {
		return fmt.Errorf("DropInfo 格式不正确，%s行%d:%s %s", filename, line, lines[line], detail)
	}

	ret := []*cm.DropInfo{}
	for i, line := range lines {

		line = strings.TrimSpace(line)
		line = util.RemoveBOM(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		res := util.SplitString(line)

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

		info := &cm.DropInfo{Low: low, High: high, ItemName: res[1], Count: 1}

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
func (d *GameData) GetMapInfoByID(mapID int) *cm.MapInfo {
	v, ok := d.MapIDInfoMap[mapID]
	if !ok {
		return nil
	}
	return v
}

// GetItemInfoByID ...
func (d *GameData) GetItemInfoByID(itemID int) *cm.ItemInfo {
	v, ok := d.ItemIDInfoMap[itemID]
	if !ok {
		return nil
	}
	return v
}

// GetItemInfoByName ...
func (d *GameData) GetItemInfoByName(itemName string) *cm.ItemInfo {
	v, ok := d.ItemNameInfoMap[itemName]
	if !ok {
		return nil
	}
	return v
}

// GetMonsterInfoByID ...
func (d *GameData) GetMonsterInfoByID(monsterID int) *cm.MonsterInfo {
	v, ok := d.MonsterIDInfoMap[monsterID]
	if !ok {
		return nil
	}
	return v
}

// GetMonsterInfoByName ...
func (d *GameData) GetMonsterInfoByName(monsterName string) *cm.MonsterInfo {
	v, ok := d.MonsterNameInfoMap[monsterName]
	if !ok {
		return nil
	}
	return v
}

// GetMagicInfoByID ...
func (d *GameData) GetMagicInfoByID(magicID int) *cm.MagicInfo {
	v, ok := d.MagicIDInfoMap[magicID]
	if !ok {
		return nil
	}
	return v
}

func (d *GameData) GetMagicInfoBySpell(spell cm.Spell) *cm.MagicInfo {
	for _, v := range d.MagicIDInfoMap {
		if v.Spell == int(spell) {
			return v
		}
	}
	return nil
}

func (d *GameData) GetMagicInfoByName(name string) *cm.MagicInfo {
	for _, v := range d.MagicIDInfoMap {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (d *GameData) GetRealItem(origin *cm.ItemInfo, level uint16, job cm.MirClass, itemList []*cm.ItemInfo) *cm.ItemInfo {
	if origin.ClassBased && origin.LevelBased {
		return d.GetClassAndLevelBasedItem(origin, job, level, itemList)
	}
	if origin.ClassBased {
		return d.GetClassBasedItem(origin, job, itemList)
	}
	if origin.LevelBased {
		return d.GetLevelBasedItem(origin, level, itemList)
	}
	return origin
}

func (d *GameData) GetLevelBasedItem(origin *cm.ItemInfo, level uint16, itemList []*cm.ItemInfo) *cm.ItemInfo {
	output := origin
	for i := 0; i < len(itemList); i++ {
		info := itemList[i]
		// if info.Name.StartsWith(Origin.Name) {
		if strings.HasPrefix(info.Name, origin.Name) {
			if (info.RequiredType == cm.RequiredTypeLevel) && (uint16(info.RequiredAmount) <= level) && (output.RequiredAmount < info.RequiredAmount) && (origin.RequiredGender == info.RequiredGender) {
				output = info
			}
		}
	}
	return output
}

func (d *GameData) GetClassBasedItem(origin *cm.ItemInfo, job cm.MirClass, itemList []*cm.ItemInfo) *cm.ItemInfo {
	for i := 0; i < len(itemList); i++ {
		info := itemList[i]
		if strings.HasPrefix(info.Name, origin.Name) {
			if (uint8(info.RequiredClass) == (1 << uint8(job))) && (origin.RequiredGender == info.RequiredGender) {
				return info
			}
		}
	}
	return origin
}

func (d *GameData) GetClassAndLevelBasedItem(origin *cm.ItemInfo, job cm.MirClass, level uint16, itemList []*cm.ItemInfo) *cm.ItemInfo {
	output := origin
	for i := 0; i < len(itemList); i++ {
		info := itemList[i]
		if strings.HasPrefix(info.Name, origin.Name) {
			if uint8(info.RequiredClass) == (1 << uint8(job)) {
				if (info.RequiredType == cm.RequiredTypeLevel) && (uint16(info.RequiredAmount) <= level) && (output.RequiredAmount <= info.RequiredAmount) && (origin.RequiredGender == info.RequiredGender) {
					output = info
				}
			}
		}
	}
	return output
}
