package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
)

// StdItem 表示数据库中的物品定义。
// Expands 对应 Expand1–Expand5，Elements 对应 Element1–Element24。
type StdItem struct {
	Idx               int    `json:"idx"`
	Name              string `json:"name"`
	StdMode           int    `json:"stdMode"`
	Shape             int    `json:"shape"`
	Weight            int    `json:"weight"`
	AniCount          int    `json:"aniCount"`
	Source            int    `json:"source"`
	Reserved          int    `json:"reserved"`
	Looks             int    `json:"looks"`
	DuraMax           int    `json:"duraMax"`
	AC                int    `json:"ac"`
	ACMax             int    `json:"acMax"`
	MAC               int    `json:"mac"`
	MACMax            int    `json:"macMax"`
	DC                int    `json:"dc"`
	DCMax             int    `json:"dcMax"`
	MC                int    `json:"mc"`
	MCMax             int    `json:"mcMax"`
	SC                int    `json:"sc"`
	SCMax             int    `json:"scMax"`
	Need              int    `json:"need"`
	NeedLevel         int    `json:"needLevel"`
	Price             int    `json:"price"`
	Stock             int    `json:"stock"`
	Color             int    `json:"color"`
	OverLap           int    `json:"overLap"`
	HP                int    `json:"hp"`
	MP                int    `json:"mp"`
	Light             int    `json:"light"`
	Horse             int    `json:"horse"`
	Element           int    `json:"element"`
	Expands           []int  `json:"expands"`
	InsuranceCurrency int    `json:"insuranceCurrency"`
	InsuranceGold     int    `json:"insuranceGold"`
	Elements          []int  `json:"elements"`
}

// MonsterDef 表示数据库中的怪物定义。
type MonsterDef struct {
	Name               string `json:"name"`
	Race               int    `json:"race"`
	RaceImg            int    `json:"raceImg"`
	Appr               int    `json:"appr"`
	Lvl                int    `json:"lvl"`
	Undead             int    `json:"undead"`
	CoolEye            int    `json:"coolEye"`
	Exp                int    `json:"exp"`
	HP                 int    `json:"hp"`
	MP                 int    `json:"mp"`
	AC                 int    `json:"ac"`
	MAC                int    `json:"mac"`
	DC                 int    `json:"dc"`
	DCMax              int    `json:"dcMax"`
	MC                 int    `json:"mc"`
	SC                 int    `json:"sc"`
	Speed              int    `json:"speed"`
	Hit                int    `json:"hit"`
	WalkSpeed          int    `json:"walkSpeed"`
	WalkStep           int    `json:"walkStep"`
	WalkWait           int    `json:"walkWait"`
	AttackSpd          int    `json:"attackSpd"`
	AttackState        int    `json:"attackState"`
	ExploreItem        int    `json:"exploreItem"`
	AttackSource       int    `json:"attackSource"`
	DisableSimpleActor int    `json:"disableSimpleActor"`
}

// MagicDef 表示数据库中的魔法定义。
type MagicDef struct {
	MagID        int    `json:"magId"`
	MagName      string `json:"magName"`
	EffectType   int    `json:"effectType"`
	Effect       int    `json:"effect"`
	Spell        int    `json:"spell"`
	Power        int    `json:"power"`
	MaxPower     int    `json:"maxPower"`
	DefSpell     int    `json:"defSpell"`
	DefPower     int    `json:"defPower"`
	DefMaxPower  int    `json:"defMaxPower"`
	Job          int    `json:"job"`
	NeedL1       int    `json:"needL1"`
	L1Train      int    `json:"l1Train"`
	NeedL2       int    `json:"needL2"`
	L2Train      int    `json:"l2Train"`
	NeedL3       int    `json:"needL3"`
	L3Train      int    `json:"l3Train"`
	Delay        int    `json:"delay"`
	Descr        string `json:"descr"`
	NeedL4       int    `json:"needL4"`
	L4Train      int    `json:"l4Train"`
	NeedL5       int    `json:"needL5"`
	L5Train      int    `json:"l5Train"`
	NeedL6       int    `json:"needL6"`
	L6Train      int    `json:"l6Train"`
	NeedL7       int    `json:"needL7"`
	L7Train      int    `json:"l7Train"`
	NeedL8       int    `json:"needL8"`
	L8Train      int    `json:"l8Train"`
	NeedL9       int    `json:"needL9"`
	L9Train      int    `json:"l9Train"`
	NeedL10      int    `json:"needL10"`
	L10Train     int    `json:"l10Train"`
	NeedL11      int    `json:"needL11"`
	L11Train     int    `json:"l11Train"`
	NeedL12      int    `json:"needL12"`
	L12Train     int    `json:"l12Train"`
	NeedL13      int    `json:"needL13"`
	L13Train     int    `json:"l13Train"`
	NeedL14      int    `json:"needL14"`
	L14Train     int    `json:"l14Train"`
	NeedL15      int    `json:"needL15"`
	L15Train     int    `json:"l15Train"`
	MaxTrainLv   int    `json:"maxTrainLv"`
	CanUpgrade   int    `json:"canUpgrade"`
	MaxUpgradeLv int    `json:"maxUpgradeLv"`
}

// ConvertDatabase 将 SQLite 数据库转换为 JSONC 文件。
func ConvertDatabase(inputDir, outputDir string) error {
	dbFile := filepath.Join(inputDir, "数据库", "GEEM2.db")

	if !FileExists(dbFile) {
		return fmt.Errorf("database file not found: %s", dbFile)
	}

	db, err := ParseSQLite(dbFile)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// 转换 StdItems
	if err := convertStdItems(db, outputDir); err != nil {
		return fmt.Errorf("converting StdItems: %w", err)
	}

	// 转换 Monster
	if err := convertMonster(db, outputDir); err != nil {
		return fmt.Errorf("converting Monster: %w", err)
	}

	// 转换 Magic
	if err := convertMagic(db, outputDir); err != nil {
		return fmt.Errorf("converting Magic: %w", err)
	}

	return nil
}

func convertStdItems(db *sql.DB, outputDir string) error {
	query := `SELECT Idx, COALESCE(Name,''), Stdmode, Shape, Weight,
		COALESCE(Anicount,0), COALESCE(Source,0), COALESCE(Reserved,0),
		Looks, DuraMax, Ac, Ac2, Mac, Mac2, Dc, Dc2, Mc, Mc2, Sc, Sc2,
		Need, NeedLevel, Price,
		COALESCE(Stock,0), COALESCE(Color,0), COALESCE(OverLap,0),
		COALESCE(HP,0), COALESCE(MP,0), COALESCE(Light,0), COALESCE(Horse,0),
		COALESCE(Element,0),
		COALESCE(Expand1,0), COALESCE(Expand2,0), COALESCE(Expand3,0), COALESCE(Expand4,0), COALESCE(Expand5,0),
		COALESCE(InsuranceCurrency,0), COALESCE(InsuranceGold,0),
		COALESCE(Element1,0), COALESCE(Element2,0), COALESCE(Element3,0), COALESCE(Element4,0),
		COALESCE(Element5,0), COALESCE(Element6,0), COALESCE(Element7,0), COALESCE(Element8,0),
		COALESCE(Element9,0), COALESCE(Element10,0), COALESCE(Element11,0), COALESCE(Element12,0),
		COALESCE(Element13,0), COALESCE(Element14,0), COALESCE(Element15,0), COALESCE(Element16,0),
		COALESCE(Element17,0), COALESCE(Element18,0), COALESCE(Element19,0), COALESCE(Element20,0),
		COALESCE(Element21,0), COALESCE(Element22,0), COALESCE(Element23,0), COALESCE(CAST(Element24 AS INTEGER),0)
		FROM StdItems`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []StdItem
	for rows.Next() {
		var item StdItem
		var expands [5]int
		var elements [24]int
		args := []interface{}{
			&item.Idx, &item.Name, &item.StdMode, &item.Shape, &item.Weight,
			&item.AniCount, &item.Source, &item.Reserved,
			&item.Looks, &item.DuraMax, &item.AC, &item.ACMax, &item.MAC, &item.MACMax,
			&item.DC, &item.DCMax, &item.MC, &item.MCMax, &item.SC, &item.SCMax,
			&item.Need, &item.NeedLevel, &item.Price,
			&item.Stock, &item.Color, &item.OverLap, &item.HP, &item.MP, &item.Light, &item.Horse, &item.Element,
		}
		for i := 0; i < 5; i++ {
			args = append(args, &expands[i])
		}
		args = append(args, &item.InsuranceCurrency, &item.InsuranceGold)
		for i := 0; i < 24; i++ {
			args = append(args, &elements[i])
		}
		if err := rows.Scan(args...); err != nil {
			return err
		}
		item.Expands = expands[:]
		item.Elements = elements[:]
		items = append(items, item)
	}

	result := map[string]interface{}{
		"items": items,
	}

	data, err := MarshalJSONC(result, stdItemComments, false)
	if err != nil {
		return fmt.Errorf("marshaling std items: %w", err)
	}

	outputFile := filepath.Join(outputDir, "items", "std_items.jsonc")
	comment := fmt.Sprintf("物品定义表（所有物品的基础属性模板）\n来源: asset/server/数据库/GEEM2.db → StdItems 表\n数量: %d 个物品\n", len(items)) + stdItemFieldDocs

	return WriteJSONC(outputFile, data, comment)
}

// stdItemFieldDocs 是 std_items.jsonc 头部注释追加的字段说明表，
// 按 StdItem 结构体字段顺序，每个字段一行。
const stdItemFieldDocs = `字段说明（按结构体字段顺序）:
idx               数据库索引（物品实例 WIndex 即指向该值）
name              物品名称
stdMode           物品类型（武器 5/6、衣服 10/11、首饰 15-26、肉 40、矿石 43/45 等）
shape             外形/子类型（StdMode 31 打包物品按 Shape 查 UnbindList 解包）
weight            重量
aniCount          动画帧数（Delphi AniCount）
source            Delphi 原义为武器神圣值（Delphi Source）
reserved          保留字段（Delphi reserved）
looks             外观，Items.WIL 中的图标索引（Delphi Looks）
duraMax           耐久上限（千分单位，1000 = 10 点耐久）
ac                防御下限（AC，Delphi Ac 列）
acMax             防御上限（Delphi Ac2 列）
mac               魔御下限（MAC，Delphi Mac 列）
macMax            魔御上限（Delphi Mac2 列）
dc                攻击下限（DC，Delphi Dc 列）
dcMax             攻击上限（Delphi Dc2 列）
mc                魔法下限（MC）
mcMax             魔法上限（Delphi Mc2 列）
sc                道术下限（SC）
scMax             道术上限（Delphi Sc2 列）
need              佩戴需求类型（0=等级 1=攻击 2=魔法 3=道术，Delphi Need）
needLevel         佩戴需求数值（分支 10-13/40-44/70/81/82 按高低字打包双条件，见 itemdb.go）
price             价格（金币）
stock             NPC 商店库存（Delphi stock）
color             颜色（GEEM2 扩展列，Go 服务端暂未消费）
overLap           可堆叠标记（GEEM2 扩展列，已读取未消费）
hp                附加 HP（GEEM2 扩展列，已读取未消费）
mp                附加 MP（GEEM2 扩展列，已读取未消费）
light             光源标记（GEEM2 扩展列，已读取未消费）
horse             坐骑标记（GEEM2 扩展列，Go 服务端暂未消费）
element           元素属性（GEEM2 扩展列，Go 服务端暂未消费）
expands           扩展属性组 Expand1-5（GEEM2 扩展列，Go 服务端暂未消费）
insuranceCurrency 保险货币（GEEM2 装备保险，Go 服务端暂未消费）
insuranceGold     保险金币（GEEM2 装备保险，Go 服务端暂未消费）
elements          元素属性组 Element1-24（GEEM2 扩展列，Go 服务端暂未消费）`

// stdItemComments 是 std_items.jsonc 的中文注释表。
// 键为点号路径；commentAll=false 时元素注释只渲染在首条记录。
// 文案来源：cmd/server/itemdb.go 消费逻辑、Delphi Grobal2.pas TStdItem
// 与 LocalDB.pas LoadItemsDB；GEEM2 扩展列服务端未消费者如实标注。
var stdItemComments = map[string]string{
	"items":                     "物品定义（StdItems 表记录，字段说明见文件头）",
	"items[].idx":               "数据库索引（物品实例 WIndex）",
	"items[].name":              "物品名称",
	"items[].stdMode":           "物品类型（武器 5/6、衣服 10/11、首饰 15-26、肉 40、矿石 43/45）",
	"items[].shape":             "外形/子类型",
	"items[].weight":            "重量",
	"items[].aniCount":          "动画帧数",
	"items[].source":            "Delphi 原义: 武器神圣值",
	"items[].reserved":          "保留字段",
	"items[].looks":             "外观（Items.WIL 图标索引）",
	"items[].duraMax":           "耐久上限（1000 = 10 点耐久）",
	"items[].ac":                "防御下限（AC）",
	"items[].acMax":             "防御上限",
	"items[].mac":               "魔御下限（MAC）",
	"items[].macMax":            "魔御上限",
	"items[].dc":                "攻击下限（DC）",
	"items[].dcMax":             "攻击上限",
	"items[].mc":                "魔法下限（MC）",
	"items[].mcMax":             "魔法上限",
	"items[].sc":                "道术下限（SC）",
	"items[].scMax":             "道术上限",
	"items[].need":              "佩戴需求类型（0=等级 1=攻击 2=魔法 3=道术）",
	"items[].needLevel":         "佩戴需求数值（部分分支高低字打包双条件）",
	"items[].price":             "价格（金币）",
	"items[].stock":             "NPC 商店库存",
	"items[].color":             "颜色（GEEM2 扩展，未消费）",
	"items[].overLap":           "可堆叠标记（GEEM2 扩展，未消费）",
	"items[].hp":                "附加 HP（GEEM2 扩展，未消费）",
	"items[].mp":                "附加 MP（GEEM2 扩展，未消费）",
	"items[].light":             "光源标记（GEEM2 扩展，未消费）",
	"items[].horse":             "坐骑标记（GEEM2 扩展，未消费）",
	"items[].element":           "元素属性（GEEM2 扩展，未消费）",
	"items[].expands":           "扩展属性组 Expand1-5（GEEM2 扩展，未消费）",
	"items[].insuranceCurrency": "保险货币（GEEM2 装备保险，未消费）",
	"items[].insuranceGold":     "保险金币（GEEM2 装备保险，未消费）",
	"items[].elements":          "元素属性组 Element1-24（GEEM2 扩展，未消费）",
}

func convertMonster(db *sql.DB, outputDir string) error {
	query := `SELECT Name, Race, RaceImg, Appr, Lvl, Undead, COALESCE(CoolEye,0), Exp, HP, COALESCE(MP,0),
		AC, MAC, DC, DCMAX, MC, SC, SPEED, HIT,
		COALESCE(WALK_SPD,0), COALESCE(WalkStep,0), COALESCE(WalkWait,0), COALESCE(ATTACK_SPD,0),
		COALESCE(AttackState,0), COALESCE(ExploreItem,0), COALESCE(AttackSource,0), COALESCE(DisableSimpleActor,0)
		FROM Monster`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var monsters []MonsterDef
	for rows.Next() {
		var m MonsterDef
		err := rows.Scan(
			&m.Name, &m.Race, &m.RaceImg, &m.Appr, &m.Lvl, &m.Undead, &m.CoolEye, &m.Exp, &m.HP, &m.MP,
			&m.AC, &m.MAC, &m.DC, &m.DCMax, &m.MC, &m.SC, &m.Speed, &m.Hit,
			&m.WalkSpeed, &m.WalkStep, &m.WalkWait, &m.AttackSpd,
			&m.AttackState, &m.ExploreItem, &m.AttackSource, &m.DisableSimpleActor,
		)
		if err != nil {
			return err
		}
		monsters = append(monsters, m)
	}

	result := map[string]interface{}{
		"monsters": monsters,
	}

	data, err := MarshalJSONC(result, monsterComments, false)
	if err != nil {
		return fmt.Errorf("marshaling monsters: %w", err)
	}

	outputFile := filepath.Join(outputDir, "monsters", "monster_db.jsonc")
	comment := fmt.Sprintf("怪物定义表（所有怪物的基础属性模板）\n来源: asset/server/数据库/GEEM2.db → Monster 表\n数量: %d 个怪物\n", len(monsters)) + monsterFieldDocs

	return WriteJSONC(outputFile, data, comment)
}

// monsterFieldDocs 是 monster_db.jsonc 头部注释追加的字段说明表，
// 按 MonsterDef 结构体字段顺序，每个字段一行。
const monsterFieldDocs = `字段说明（按结构体字段顺序）:
name               怪物名称（刷怪点与掉落表的索引键）
race               种族/AI 类型（Delphi btRace；Go 按它分发 AI 初始状态，如 51/53/84 动物、52 鹿、110/111 城门/城墙）
raceImg            种族形象集编号（Delphi btRaceImg）
appr               外观编号（Delphi wAppr，Mon.wil 形象索引）
lvl                怪物等级（Delphi wLevel）
undead             不死属性（Delphi btLifeAttrib；>0 为不死系）
coolEye            视野 CoolEye（Delphi wCoolEye；影响对隐身目标的感知）
exp                击杀给予经验（Delphi dwExp）
hp                 HP 上限（Delphi wHP；城门/城墙 race 110/111 的 HP 控制建筑显示）
mp                 MP 上限（Delphi wMP；Go 服务端暂未消费）
ac                 防御（AC，Delphi wAC）
mac                魔御（MAC，Delphi wMAC）
dc                 攻击下限（DC，Delphi wDC）
dcMax              攻击上限（Delphi DCMAX 列 → wMaxDC）
mc                 魔法（MC，Delphi wMC）
sc                 道术（SC，Delphi wSC）
speed              速度点数（Delphi wSpeed → SpeedPoint）
hit                准确（Delphi wHitPoint → HitPoint）
walkSpeed          行走速度 ms/步（Delphi WALK_SPD 列；下限 200，默认 1400）
walkStep           单次连续行走步数（Delphi WalkStep；默认 3）
walkWait           行走停顿 ms（Delphi WalkWait；默认 1000）
attackSpd          攻击速度 ms/次（Delphi ATTACK_SPD 列；Go 暂未读取，用 walkSpeed+attackSpeedOffset）
attackState        保留字段（GEEM2 扩展列，Go 服务端暂未消费）
exploreItem        保留字段（GEEM2 扩展列，Go 服务端暂未消费）
attackSource       保留字段（GEEM2 扩展列，Go 服务端暂未消费）
disableSimpleActor 保留字段（GEEM2 扩展列，Go 服务端暂未消费）`

// monsterComments 是 monster_db.jsonc 的中文注释表。
// 键为点号路径；commentAll=false 时元素注释只渲染在首条记录。
// 文案来源：cmd/server/mongen.go initMonsterFromDef 消费逻辑、
// Delphi LocalDB.pas LoadMonsterDB。
var monsterComments = map[string]string{
	"monsters":                      "怪物定义（Monster 表记录，字段说明见文件头）",
	"monsters[].name":               "怪物名称（刷怪点/掉落表索引键）",
	"monsters[].race":               "种族/AI 类型（51/53/84 动物、52 鹿、110/111 城门/城墙等）",
	"monsters[].raceImg":            "种族形象集编号",
	"monsters[].appr":               "外观编号（Mon.wil 形象索引）",
	"monsters[].lvl":                "怪物等级",
	"monsters[].undead":             "不死属性（>0 为不死系）",
	"monsters[].coolEye":            "视野 CoolEye（影响隐身感知）",
	"monsters[].exp":                "击杀给予经验",
	"monsters[].hp":                 "HP 上限（城门/城墙的 HP 控制建筑显示）",
	"monsters[].mp":                 "MP 上限（未消费）",
	"monsters[].ac":                 "防御（AC）",
	"monsters[].mac":                "魔御（MAC）",
	"monsters[].dc":                 "攻击下限（DC）",
	"monsters[].dcMax":              "攻击上限",
	"monsters[].mc":                 "魔法（MC）",
	"monsters[].sc":                 "道术（SC）",
	"monsters[].speed":              "速度点数（SpeedPoint）",
	"monsters[].hit":                "准确（HitPoint）",
	"monsters[].walkSpeed":          "行走速度 ms/步（下限 200，默认 1400）",
	"monsters[].walkStep":           "单次连续行走步数（默认 3）",
	"monsters[].walkWait":           "行走停顿 ms（默认 1000）",
	"monsters[].attackSpd":          "攻击速度 ms/次（Go 暂未消费，用 walkSpeed+偏移）",
	"monsters[].attackState":        "保留字段（GEEM2 扩展，未消费）",
	"monsters[].exploreItem":        "保留字段（GEEM2 扩展，未消费）",
	"monsters[].attackSource":       "保留字段（GEEM2 扩展，未消费）",
	"monsters[].disableSimpleActor": "保留字段（GEEM2 扩展，未消费）",
}

func convertMagic(db *sql.DB, outputDir string) error {
	query := `SELECT MagID, MagName, EffectType, Effect, Spell, Power, MaxPower,
		DefSpell, DefPower, DefMaxPower, Job,
		NeedL1, L1Train, NeedL2, L2Train, NeedL3, L3Train, Delay, COALESCE(Descr,''),
		NeedL4, L4Train, NeedL5, L5Train, NeedL6, L6Train, NeedL7, L7Train,
		NeedL8, L8Train, NeedL9, L9Train, NeedL10, L10Train, NeedL11, L11Train,
		NeedL12, L12Train, NeedL13, L13Train, NeedL14, L14Train, NeedL15, L15Train,
		MaxTrainLv, CanUpgrade, MaxUpgradeLv
		FROM Magic`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var magics []MagicDef
	for rows.Next() {
		var m MagicDef
		err := rows.Scan(
			&m.MagID, &m.MagName, &m.EffectType, &m.Effect, &m.Spell, &m.Power, &m.MaxPower,
			&m.DefSpell, &m.DefPower, &m.DefMaxPower, &m.Job,
			&m.NeedL1, &m.L1Train, &m.NeedL2, &m.L2Train, &m.NeedL3, &m.L3Train, &m.Delay, &m.Descr,
			&m.NeedL4, &m.L4Train, &m.NeedL5, &m.L5Train, &m.NeedL6, &m.L6Train, &m.NeedL7, &m.L7Train,
			&m.NeedL8, &m.L8Train, &m.NeedL9, &m.L9Train, &m.NeedL10, &m.L10Train, &m.NeedL11, &m.L11Train,
			&m.NeedL12, &m.L12Train, &m.NeedL13, &m.L13Train, &m.NeedL14, &m.L14Train, &m.NeedL15, &m.L15Train,
			&m.MaxTrainLv, &m.CanUpgrade, &m.MaxUpgradeLv,
		)
		if err != nil {
			return err
		}
		magics = append(magics, m)
	}

	result := map[string]interface{}{
		"magic": magics,
	}

	data, err := MarshalJSONC(result, magicComments, false)
	if err != nil {
		return fmt.Errorf("marshaling magic: %w", err)
	}

	outputFile := filepath.Join(outputDir, "magic", "magic_db.jsonc")
	comment := fmt.Sprintf("魔法定义表（所有魔法的基础属性模板）\n来源: asset/server/数据库/GEEM2.db → Magic 表\n数量: %d 个魔法\n", len(magics)) + magicFieldDocs

	return WriteJSONC(outputFile, data, comment)
}

// magicFieldDocs 是 magic_db.jsonc 头部注释追加的字段说明表，
// 按 MagicDef 结构体字段顺序，每个字段一行。
const magicFieldDocs = `字段说明（按结构体字段顺序）:
magId        魔法编号（Delphi wMagicId；魔法表索引键）
magName      魔法名称（Delphi sMagicName）
effectType   客户端特效类型（Delphi btEffectType；客户端 NewMagic 按它分发）
effect       特效动画编号（Delphi btEffect；Go 默认特效号取该值）
spell        施法 MP 消耗（Delphi wSpell；Go 施法直接扣该值）
power        威力成长值（按技能等级缩放，Delphi wPower）
maxPower     威力成长上限（Delphi wMaxPower）
defSpell     基础 MP 消耗（与 spell 的等级缩放部分相加，Delphi btDefSpell；Go 暂未消费）
defPower     基础威力下限（Delphi btDefPower；威力随机下限）
defMaxPower  基础威力上限（Delphi btDefMaxPower；威力随机上限）
job          职业（0=战士 1=法师 2=道士，Delphi btJob）
needL1       1 级技能学习等级（Delphi TrainLevel[0]；Go 学书校验用）
l1Train      1 级技能所需熟练度（Delphi MaxTrain[0]；Go 暂用 config trainThreshold0-2）
needL2       2 级技能需求等级（Delphi TrainLevel[1]；Go 暂未消费）
l2Train      2 级技能所需熟练度（Go 暂未消费）
needL3       3 级技能需求等级（Delphi TrainLevel[2]；Go 暂未消费）
l3Train      3 级技能所需熟练度（Go 暂未消费）
delay        施法冷却 ms（Delphi dwDelayTime；Go 暂未消费）
descr        魔法描述（Delphi sDescr；Go 暂未消费）
needL4       4 级技能需求等级（GEEM2 高阶技能等级；Go 暂未消费）
l4Train      4 级技能所需熟练度（Go 暂未消费）
needL5       5 级技能需求等级（Go 暂未消费）
l5Train      5 级技能所需熟练度（Go 暂未消费）
needL6       6 级技能需求等级（Go 暂未消费）
l6Train      6 级技能所需熟练度（Go 暂未消费）
needL7       7 级技能需求等级（Go 暂未消费）
l7Train      7 级技能所需熟练度（Go 暂未消费）
needL8       8 级技能需求等级（Go 暂未消费）
l8Train      8 级技能所需熟练度（Go 暂未消费）
needL9       9 级技能需求等级（Go 暂未消费）
l9Train      9 级技能所需熟练度（Go 暂未消费）
needL10      10 级技能需求等级（Go 暂未消费）
l10Train     10 级技能所需熟练度（Go 暂未消费）
needL11      11 级技能需求等级（Go 暂未消费）
l11Train     11 级技能所需熟练度（Go 暂未消费）
needL12      12 级技能需求等级（Go 暂未消费）
l12Train     12 级技能所需熟练度（Go 暂未消费）
needL13      13 级技能需求等级（Go 暂未消费）
l13Train     13 级技能所需熟练度（Go 暂未消费）
needL14      14 级技能需求等级（Go 暂未消费）
l14Train     14 级技能所需熟练度（Go 暂未消费）
needL15      15 级技能需求等级（Go 暂未消费）
l15Train     15 级技能所需熟练度（Go 暂未消费）
maxTrainLv   技能等级上限（Delphi btTrainLv；Go 心灵启示等时长计算用）
canUpgrade   技能可否升级（GEEM2 扩展列，Go 暂未消费）
maxUpgradeLv 技能升级等级上限（GEEM2 扩展列，Go 暂未消费）`

// magicComments 是 magic_db.jsonc 的中文注释表。
// 键为点号路径；commentAll=false 时元素注释只渲染在首条记录。
// 文案来源：cmd/server/magicdb.go + magicsystem.go 消费逻辑、
// Delphi Magic.pas/ObjBase.pas 的威力与耗蓝公式、LocalDB.pas LoadMagicDB。
var magicComments = map[string]string{
	"magic":                "魔法定义（Magic 表记录，字段说明见文件头）",
	"magic[].magId":        "魔法编号（索引键）",
	"magic[].magName":      "魔法名称",
	"magic[].effectType":   "客户端特效类型（Delphi btEffectType）",
	"magic[].effect":       "特效动画编号（Go 默认特效号取该值）",
	"magic[].spell":        "施法 MP 消耗（Go 直接扣该值）",
	"magic[].power":        "威力成长值（按技能等级缩放）",
	"magic[].maxPower":     "威力成长上限",
	"magic[].defSpell":     "基础 MP 消耗（未消费）",
	"magic[].defPower":     "基础威力下限（威力随机下限）",
	"magic[].defMaxPower":  "基础威力上限（威力随机上限）",
	"magic[].job":          "职业（0=战士 1=法师 2=道士）",
	"magic[].needL1":       "1 级技能学习等级（Go 学书校验用）",
	"magic[].l1Train":      "1 级技能所需熟练度（Go 暂用 config trainThreshold）",
	"magic[].needL2":       "2 级技能需求等级（未消费）",
	"magic[].l2Train":      "2 级技能所需熟练度（未消费）",
	"magic[].needL3":       "3 级技能需求等级（未消费）",
	"magic[].l3Train":      "3 级技能所需熟练度（未消费）",
	"magic[].delay":        "施法冷却 ms（未消费）",
	"magic[].descr":        "魔法描述（未消费）",
	"magic[].needL4":       "4 级技能需求等级（未消费）",
	"magic[].l4Train":      "4 级技能所需熟练度（未消费）",
	"magic[].needL5":       "5 级技能需求等级（未消费）",
	"magic[].l5Train":      "5 级技能所需熟练度（未消费）",
	"magic[].needL6":       "6 级技能需求等级（未消费）",
	"magic[].l6Train":      "6 级技能所需熟练度（未消费）",
	"magic[].needL7":       "7 级技能需求等级（未消费）",
	"magic[].l7Train":      "7 级技能所需熟练度（未消费）",
	"magic[].needL8":       "8 级技能需求等级（未消费）",
	"magic[].l8Train":      "8 级技能所需熟练度（未消费）",
	"magic[].needL9":       "9 级技能需求等级（未消费）",
	"magic[].l9Train":      "9 级技能所需熟练度（未消费）",
	"magic[].needL10":      "10 级技能需求等级（未消费）",
	"magic[].l10Train":     "10 级技能所需熟练度（未消费）",
	"magic[].needL11":      "11 级技能需求等级（未消费）",
	"magic[].l11Train":     "11 级技能所需熟练度（未消费）",
	"magic[].needL12":      "12 级技能需求等级（未消费）",
	"magic[].l12Train":     "12 级技能所需熟练度（未消费）",
	"magic[].needL13":      "13 级技能需求等级（未消费）",
	"magic[].l13Train":     "13 级技能所需熟练度（未消费）",
	"magic[].needL14":      "14 级技能需求等级（未消费）",
	"magic[].l14Train":     "14 级技能所需熟练度（未消费）",
	"magic[].needL15":      "15 级技能需求等级（未消费）",
	"magic[].l15Train":     "15 级技能所需熟练度（未消费）",
	"magic[].maxTrainLv":   "技能等级上限（Go 心灵启示等时长计算用）",
	"magic[].canUpgrade":   "技能可否升级（GEEM2 扩展，未消费）",
	"magic[].maxUpgradeLv": "技能升级等级上限（GEEM2 扩展，未消费）",
}
