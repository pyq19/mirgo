package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
	"github.com/pyq0109/mirgo/internal/storage"
)

type CastleWarState int

const (
	CastlePeace       CastleWarState = iota
	CastleWarDeclared
	CastleWarActive
)

// CastleWarDeclaration — Delphi m_AttackWarList 条目（Castle.pas:341/446）：
// 宣战时登记行会名与预约攻城日期，到期日 WarStartHour 点自动开战。
type CastleWarDeclaration struct {
	GuildName  string `json:"guildName"`
	AttackDate string `json:"attackDate"` // YYYY-MM-DD
}

type CastleConfig struct {
	Name           string `json:"name"`
	MapName        string `json:"map"`
	PalaceX        int    `json:"palaceX"`
	PalaceY        int    `json:"palaceY"`
	PalaceRadius   int    `json:"palaceRadius"`
	DoorMaxHP      int    `json:"doorMaxHP"`
	WallMaxHP      int    `json:"wallMaxHP"`
	DoorRepairCost int    `json:"doorRepairCost"`
	WallRepairCost int    `json:"wallRepairCost"`
	WarStartHour   int    `json:"warStartHour"`
	WarDurationMin int    `json:"warDurationMin"`
	GracePeriodSec int    `json:"gracePeriodSec"`
	TaxRate        int    `json:"taxRate"`
	MaxTaxRate     int    `json:"maxTaxRate"`
	ArcherCost     int    `json:"archerCost"`
	MaxGuards      int    `json:"maxGuards"`
}

func DefaultCastleConfig() *CastleConfig {
	return &CastleConfig{
		Name:           "沙巴克",
		MapName:        "3",
		PalaceX:        100,
		PalaceY:        100,
		PalaceRadius:   5,
		DoorMaxHP:      5000,
		WallMaxHP:      10000,
		DoorRepairCost: 100000,
		WallRepairCost: 200000,
		WarStartHour:   20,
		WarDurationMin: 180,
		GracePeriodSec: 300,
		TaxRate:        5,
		MaxTaxRate:     20,
		ArcherCost:     50000,
		MaxGuards:      8,
	}
}

func LoadCastleConfig(configDir string) (*CastleConfig, error) {
	cfg := DefaultCastleConfig()

	data, err := os.ReadFile(filepath.Join(configDir, "castle", "sabuk_wall.jsonc"))
	if err != nil {
		return cfg, nil // 配置文件不存在时使用默认值
	}

	if err := parseJSONC(data, cfg); err != nil {
		return nil, fmt.Errorf("parse castle config: %w", err)
	}
	return cfg, nil
}

type CastleObject struct {
	mu           sync.Mutex
	Config       CastleConfig
	OwnerGuild   string
	Gold         int64
	TaxRate      int
	DoorHP       int
	WallHP       int
	WarState     CastleWarState
	WarStartTick int64
	WarEndTick   int64
	AttackGuilds []string
	GuardIDs     []int32
	DoorOpen     bool

	// 城门/城墙可视化 actor（批6：Delphi TCastleDoor/TWallStructure
	// 由怪物实体渲染，ObjMon2.pas:1013-1096）
	structuresSpawned    bool
	DoorActorID          int32
	WallActorIDs         []int32
	doorVisualState      int  // -1 未同步；0-2 关闭破损度；7 开启；9 破坏
	wallsVisualDestroyed bool // 城墙当前可视化状态是否破坏

	// Delphi: 预约攻城列表 m_AttackWarList（Castle.pas:644-679）
	WarDeclarations []CastleWarDeclaration
	lastWarCheckDay string // 当日已完成 WarStartHour 开战检查（防重）

	// Delphi: 城堡经济 (Castle.pas:67-80)
	TechLevel   int   // 科技等级
	Power       int   // 能源
	TodayIncome int64 // 今日收入

	// Delphi: 修理冷却 (Castle.pas:1150-1219)
	DoorStruckTick int64 // 城门最后被攻击时间
	WallStruckTick int64 // 城墙最后被攻击时间
}

func NewCastleObject(cfg CastleConfig) *CastleObject {
	return &CastleObject{
		Config:  cfg,
		TaxRate: cfg.TaxRate,
		DoorHP:  cfg.DoorMaxHP,
		WallHP:  cfg.WallMaxHP,
	}
}

func (c *CastleObject) IsAtWar() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.WarState == CastleWarActive
}

func (c *CastleObject) IsAttackingGuild(guildName string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, g := range c.AttackGuilds {
		if g == guildName {
			return true
		}
	}
	return false
}

func (c *CastleObject) IsDefendingGuild(guildName string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.OwnerGuild != "" && c.OwnerGuild == guildName
}

// IsAttackAllyGuild — Delphi Castle.pas:768-783：是否为攻方行会的联盟行会。
func (c *CastleObject) IsAttackAllyGuild(engine *UserEngine, guildName string) bool {
	if guildName == "" || engine == nil {
		return false
	}
	c.mu.Lock()
	attackers := make([]string, len(c.AttackGuilds))
	copy(attackers, c.AttackGuilds)
	owner := c.OwnerGuild
	c.mu.Unlock()
	for _, name := range attackers {
		if name == owner {
			continue
		}
		if g := engine.FindGuild(name); g != nil && g.IsAllyGuild(guildName) {
			return true
		}
	}
	return false
}

// IsDefenseAllyGuild — Delphi Castle.pas:896-902：是否为城主（守方）的联盟行会，仅战争期间有效。
func (c *CastleObject) IsDefenseAllyGuild(engine *UserEngine, guildName string) bool {
	if guildName == "" || engine == nil || !c.IsAtWar() {
		return false
	}
	owner := c.GetOwnerGuild()
	if owner == "" {
		return false
	}
	if g := engine.FindGuild(owner); g != nil {
		return g.IsAllyGuild(guildName)
	}
	return false
}

// DeclareWar — Delphi 宣战（攻城）语义：登记预约战条目到 m_AttackWarList，
// 到期日 WarStartHour 点自动开战（Castle.pas:644-679）。
// attackDate 为空默认明天；不允许早于当天。
func (c *CastleObject) DeclareWar(guildName, attackDate string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if guildName == "" || guildName == c.OwnerGuild {
		return false
	}
	if c.WarState == CastleWarActive {
		return false
	}
	today := time.Now().Format("2006-01-02")
	if attackDate == "" {
		attackDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	} else {
		if _, err := time.Parse("2006-01-02", attackDate); err != nil {
			return false
		}
		if attackDate < today {
			return false
		}
	}
	for _, w := range c.WarDeclarations {
		if w.GuildName == guildName {
			return false
		}
	}

	c.WarDeclarations = append(c.WarDeclarations, CastleWarDeclaration{GuildName: guildName, AttackDate: attackDate})
	if c.WarState == CastlePeace {
		c.WarState = CastleWarDeclared
	}
	log.Logf(log.LevelInfo, "Castle", "%s declared war on castle %s for %s", guildName, c.Config.Name, attackDate)
	return true
}

// GetWarDeclarations 返回当前预约战列表副本（GM 查询/客户端展示用）。
func (c *CastleObject) GetWarDeclarations() []CastleWarDeclaration {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]CastleWarDeclaration, len(c.WarDeclarations))
	copy(out, c.WarDeclarations)
	return out
}

func (c *CastleObject) StartWar(server *netserver.TCPServer, engine *UserEngine) {
	c.mu.Lock()
	if c.WarState != CastleWarDeclared || len(c.AttackGuilds) == 0 {
		c.mu.Unlock()
		return
	}
	now := time.Now().UnixMilli()
	c.WarState = CastleWarActive
	c.WarStartTick = now
	c.WarEndTick = now + int64(c.Config.WarDurationMin)*60*1000
	attackers := make([]string, len(c.AttackGuilds))
	copy(attackers, c.AttackGuilds)
	c.mu.Unlock()

	// Delphi: 攻城开始时关门 (Castle.pas:697)
	c.ToggleDoor(false)
	c.SyncDoorVisual(server, engine)

	// Delphi: 生成弓箭守卫 (Castle.pas:260-290)
	c.spawnCastleArchers(engine, now)

	text := fmt.Sprintf("[攻城战] %s 攻城战开始！进攻方: %s，防守方: %s",
		c.Config.Name, strings.Join(attackers, ","), c.OwnerGuild)
	c.broadcastSysMsg(server, engine, text)
	c.SaveState(engine.DB())
	log.Logf(log.LevelInfo, "Castle", "castle war started: attackers=%v defender=%s", attackers, c.OwnerGuild)
}

// spawnCastleArchers — Delphi: 攻城战时在城堡地图生成弓箭手 (Castle.pas:260-290)
func (c *CastleObject) spawnCastleArchers(engine *UserEngine, now int64) {
	if engine == nil || engine.mapMgr == nil {
		return
	}
	env := engine.mapMgr.FindMap(c.Config.MapName)
	if env == nil {
		return
	}
	// 在皇宫周围生成弓箭手
	count := c.Config.MaxGuards
	if count <= 0 {
		count = 8
	}
	for i := 0; i < count; i++ {
		ax := c.Config.PalaceX + (i%4)*3 - 4
		ay := c.Config.PalaceY + (i/4)*3 - 2
		mon := engine.SpawnMonsterByName(c.Config.MapName, ax, ay, "弓箭守卫", now)
		if mon != nil {
			mon.StickMode = true
			mon.AIBehavior = AIGuard
			mon.ViewRange = 12
			c.mu.Lock()
			c.GuardIDs = append(c.GuardIDs, mon.ID)
			c.mu.Unlock()
		}
	}
}

func (c *CastleObject) EndWar(server *netserver.TCPServer, engine *UserEngine) {
	c.mu.Lock()
	if c.WarState != CastleWarActive {
		c.mu.Unlock()
		return
	}
	// 仍有后续预约战 → 回到 Declared 等待下一场（Castle.pas 预约列表语义）
	if len(c.WarDeclarations) > 0 {
		c.WarState = CastleWarDeclared
	} else {
		c.WarState = CastlePeace
	}
	c.AttackGuilds = nil
	c.WarStartTick = 0
	c.WarEndTick = 0
	c.mu.Unlock()

	text := fmt.Sprintf("[攻城战] %s 攻城战结束", c.Config.Name)
	c.broadcastSysMsg(server, engine, text)
	c.SaveState(engine.DB())
	log.Logf(log.LevelInfo, "Castle", "castle war ended")
}

// ChangeOwner 更改城堡所属行会（Delphi Castle.GetCastle，
// Castle.pas:828-846：设 m_MasterGuild/m_sOwnGuild、广播占领消息、
// SaveConfigFile；@ChangeSabukLord 调用）。
// 偏差：Delphi 新旧行会 RefMemberName 成员称号刷新在 Go 无对应，跳过。
func (c *CastleObject) ChangeOwner(server *netserver.TCPServer, engine *UserEngine, guildName string) {
	c.mu.Lock()
	c.OwnerGuild = guildName
	c.mu.Unlock()
	// Delphi 广播 '[%s 已被 %s 占领]'（Castle.pas:838）。
	c.broadcastSysMsg(server, engine, fmt.Sprintf("[%s 已被 %s 占领]", c.Config.Name, guildName))
	c.SaveState(engine.DB())
	log.Logf(log.LevelInfo, "Castle", "castle owner changed to %s", guildName)
}

// ForceStartWar 强制开始攻城战（@ForcedWallconquestWar，Delphi
// StartWallconquestWar，Castle.pas:848-861）：无需宣战预约、无进攻方
// 名单，任意状态直接开战（区别于 StartWar 的预约战前置校验）。
func (c *CastleObject) ForceStartWar(server *netserver.TCPServer, engine *UserEngine) {
	c.mu.Lock()
	if c.WarState == CastleWarActive {
		c.mu.Unlock()
		return
	}
	now := time.Now().UnixMilli()
	c.WarState = CastleWarActive
	c.WarStartTick = now
	c.WarEndTick = now + int64(c.Config.WarDurationMin)*60*1000
	c.mu.Unlock()

	c.ToggleDoor(false) // Delphi MainDoorControl(True) 关城门（Castle.pas:1136-1149）
	c.SyncDoorVisual(server, engine)
	c.spawnCastleArchers(engine, now)
	// Delphi 广播 '[%s攻城战已经开始]'（Castle.pas:856）。
	c.broadcastSysMsg(server, engine, fmt.Sprintf("[%s攻城战已经开始]", c.Config.Name))
	c.SaveState(engine.DB())
	log.Logf(log.LevelInfo, "Castle", "castle war force-started")
}

func (c *CastleObject) DamageDoor(dmg int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DoorHP -= dmg
	c.DoorStruckTick = time.Now().UnixMilli()
	if c.DoorHP <= 0 {
		c.DoorHP = 0
		log.Logf(log.LevelInfo, "Castle", "castle door destroyed")
		return true
	}
	return false
}

func (c *CastleObject) DamageWall(dmg int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Delphi: 非攻城期间城墙无敌 (Castle.pas:720-726, m_boStoneMode=True)
	if c.WarState != CastleWarActive {
		return false
	}
	c.WallHP -= dmg
	c.WallStruckTick = time.Now().UnixMilli()
	if c.WallHP <= 0 {
		c.WallHP = 0
		log.Logf(log.LevelInfo, "Castle", "castle wall destroyed")
		return true
	}
	return false
}

// RepairDoor — Delphi (Castle.pas:1150-1182): 攻城期间不可修，被攻击60秒后才能修
func (c *CastleObject) RepairDoor() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.WarState == CastleWarActive {
		return false
	}
	if c.DoorHP >= c.Config.DoorMaxHP {
		return false
	}
	now := time.Now().UnixMilli()
	if c.DoorStruckTick > 0 && now-c.DoorStruckTick < 60000 {
		return false
	}
	if c.Gold < int64(c.Config.DoorRepairCost) {
		return false
	}
	c.Gold -= int64(c.Config.DoorRepairCost)
	c.DoorHP = c.Config.DoorMaxHP
	log.Logf(log.LevelInfo, "Castle", "castle door repaired (treasury: %d)", c.Gold)
	return true
}

// RepairWall — Delphi (Castle.pas:1184-1219): 攻城期间不可修，被攻击60秒后才能修
func (c *CastleObject) RepairWall() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.WarState == CastleWarActive {
		return false
	}
	if c.WallHP >= c.Config.WallMaxHP {
		return false
	}
	now := time.Now().UnixMilli()
	if c.WallStruckTick > 0 && now-c.WallStruckTick < 60000 {
		return false
	}
	if c.Gold < int64(c.Config.WallRepairCost) {
		return false
	}
	c.Gold -= int64(c.Config.WallRepairCost)
	c.WallHP = c.Config.WallMaxHP
	log.Logf(log.LevelInfo, "Castle", "castle wall repaired (treasury: %d)", c.Gold)
	return true
}

func (c *CastleObject) CheckCapture(engine *UserEngine) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.WarState != CastleWarActive {
		return ""
	}

	radius := c.Config.PalaceRadius
	if radius <= 0 {
		radius = 5
	}

	guildCounts := make(map[string]int)
	for _, p := range engine.allPlayers() {
		if p.Ghost || p.Death || p.MapName != c.Config.MapName {
			continue
		}
		dx := abs(p.CurrX - c.Config.PalaceX)
		dy := abs(p.CurrY - c.Config.PalaceY)
		if dx > radius || dy > radius {
			continue
		}
		if p.GuildName == "" {
			continue
		}
		if !c.isAttackingGuildLocked(p.GuildName) {
			continue
		}
		guildCounts[p.GuildName]++
	}

	for guildName, count := range guildCounts {
		if count >= 3 {
			c.OwnerGuild = guildName
			if len(c.WarDeclarations) > 0 {
				c.WarState = CastleWarDeclared
			} else {
				c.WarState = CastlePeace
			}
			c.AttackGuilds = nil
			c.SaveState(engine.DB())
			log.Logf(log.LevelInfo, "Castle", "castle captured by %s (%d members in palace)", guildName, count)
			return guildName
		}
	}
	return ""
}

func (c *CastleObject) CollectTax(amount int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Gold += amount
	c.TodayIncome += amount
}

func (c *CastleObject) WithdrawGold(amount int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if amount <= 0 || amount > c.Gold {
		return false
	}
	c.Gold -= amount
	return true
}

// DepositGold — Delphi: 城主存款 (Castle.pas:120)
func (c *CastleObject) DepositGold(amount int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if amount <= 0 {
		return false
	}
	c.Gold += amount
	return true
}

func (c *CastleObject) SetTaxRate(rate int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if rate < 0 || rate > c.Config.MaxTaxRate {
		return false
	}
	c.TaxRate = rate
	return true
}

func (c *CastleObject) AddGuard(id int32) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.GuardIDs) >= c.Config.MaxGuards {
		return false
	}
	c.GuardIDs = append(c.GuardIDs, id)
	return true
}

func (c *CastleObject) RemoveGuard(id int32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, gid := range c.GuardIDs {
		if gid == id {
			c.GuardIDs = append(c.GuardIDs[:i], c.GuardIDs[i+1:]...)
			return
		}
	}
}

func (c *CastleObject) GuardCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.GuardIDs)
}

func (c *CastleObject) GetGold() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Gold
}

func (c *CastleObject) GetDoorHP() (int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.DoorHP, c.Config.DoorMaxHP
}

func (c *CastleObject) GetWallHP() (int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.WallHP, c.Config.WallMaxHP
}

func (c *CastleObject) GetWarState() CastleWarState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.WarState
}

func (c *CastleObject) GetOwnerGuild() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.OwnerGuild
}

func (c *CastleObject) GetTaxRate() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.TaxRate
}

func (c *CastleObject) ToggleDoor(open bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DoorOpen = open
}

func (c *CastleObject) IsDoorOpen() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.DoorOpen
}

// DoorPosition 返回城门的地图坐标（皇宫入口前方）。
func (c *CastleObject) DoorPosition() (int, int) {
	return c.Config.PalaceX, c.Config.PalaceY + c.Config.PalaceRadius + 1
}

// WallPositions 返回城墙段的坐标列表（皇宫周围）。
func (c *CastleObject) WallPositions() [][2]int {
	r := c.Config.PalaceRadius
	px, py := c.Config.PalaceX, c.Config.PalaceY
	var walls [][2]int
	for dx := -r; dx <= r; dx++ {
		walls = append(walls, [2]int{px + dx, py - r})
		walls = append(walls, [2]int{px + dx, py + r})
	}
	for dy := -r + 1; dy < r; dy++ {
		walls = append(walls, [2]int{px - r, py + dy})
		walls = append(walls, [2]int{px + r, py + dy})
	}
	return walls
}

// SpawnCastleStructures 生成城门/城墙静态可视化 actor（Delphi 城堡建筑
// 以怪物实体渲染：SabukDoor race110 appr900，SabukW1-3 race111 appr901-903，
// 帧布局在 Effect.wil nrace90 槽位，Actor.pas:1147-1152）。只执行一次。
func (c *CastleObject) SpawnCastleStructures(engine *UserEngine, now int64) {
	c.mu.Lock()
	if c.structuresSpawned {
		c.mu.Unlock()
		return
	}
	c.structuresSpawned = true
	c.mu.Unlock()

	if engine == nil || engine.mapMgr == nil {
		return
	}
	if env := engine.mapMgr.FindMap(c.Config.MapName); env == nil {
		return
	}

	// 城门
	doorX, doorY := c.DoorPosition()
	if door := engine.SpawnCastleStructure(c.Config.MapName, doorX, doorY, "SabukDoor", now); door != nil {
		c.mu.Lock()
		door.Dir = c.doorVisualDirLocked()
		c.DoorActorID = door.ID
		c.doorVisualState = door.Dir
		c.mu.Unlock()
	}

	// 城墙（W1/W2/W3 按序循环取外观）
	wallNames := []string{"SabukW1", "SabukW2", "SabukW3"}
	ids := make([]int32, 0)
	for i, w := range c.WallPositions() {
		mon := engine.SpawnCastleStructure(c.Config.MapName, w[0], w[1], wallNames[i%3], now)
		if mon != nil {
			mon.Dir = 0
			ids = append(ids, mon.ID)
		}
	}
	c.mu.Lock()
	c.WallActorIDs = ids
	c.mu.Unlock()
	log.Logf(log.LevelInfo, "Castle", "spawned castle structures: door id=%d walls=%d", c.DoorActorID, len(ids))
}

// doorVisualDirLocked 城门展示 dir：0-2=关闭三档破损度，7=开启
//（Delphi 关门 dir=3-round(HP/MaxHP*3)，开门 dir=7，ObjMon2.pas:1042/1053）。
func (c *CastleObject) doorVisualDirLocked() int {
	if c.DoorOpen || c.DoorHP <= 0 {
		return 7
	}
	max := c.Config.DoorMaxHP
	if max <= 0 {
		return 0
	}
	dmgLevel := 3 - (c.DoorHP*3+max-1)/max // 破损度 0(完好)..2
	if dmgLevel < 0 {
		dmgLevel = 0
	}
	if dmgLevel > 2 {
		dmgLevel = 2
	}
	return dmgLevel
}

// castleEnv 返回城堡所在地图。
func (c *CastleObject) castleEnv(engine *UserEngine) *Environment {
	if engine == nil || engine.mapMgr == nil {
		return nil
	}
	return engine.mapMgr.FindMap(c.Config.MapName)
}

// SyncDoorVisual 按城门当前状态广播门 actor 动画（Delphi TCastleDoor
// Open/Close/RefStatus，ObjMon2.pas:1039-1096）：
// 开启→RM_DIGUP(dir=7)、关闭→RM_DIGDOWN(dir=破损度)、破坏→RM_DEATH、
// 修复→RM_ALIVE、破损度变化→RM_TURN。
func (c *CastleObject) SyncDoorVisual(server *netserver.TCPServer, engine *UserEngine) {
	env := c.castleEnv(engine)
	if env == nil {
		return
	}
	c.mu.Lock()
	doorID := c.DoorActorID
	destroyed := c.DoorHP <= 0
	prev := c.doorVisualState
	state := c.doorVisualDirLocked()
	if destroyed {
		state = 9
	}
	if state == prev {
		c.mu.Unlock()
		return
	}
	c.doorVisualState = state
	c.mu.Unlock()

	door, ok := env.getObjectByID(doorID).(*MonsterObject)
	if !ok || door == nil {
		return
	}
	switch {
	case state == 9 && !door.Death:
		// 破坏：死亡动画（客户端 ActDie+Effect.wil 120..129 破坏特效）
		door.Death = true
		door.DeathTick = time.Now().UnixMilli()
		env.broadcastDeathMsg(door.BaseObject, door.ID, door.CurrX, door.CurrY, door.Dir, true)
	case state != 9 && door.Death:
		// 修复：复活 + 关闭态
		door.Death = false
		door.DeathTick = 0
		door.Dir = state
		door.SendRefMsg(RM_ALIVE, door.Dir, door.CurrX, door.CurrY, "")
	case state == 7:
		door.Dir = 7
		door.SendRefMsg(RM_DIGUP, 7, door.CurrX, door.CurrY, "")
	case prev == 7:
		// 开门→关闭：关门动画（dir 携带破损度）
		door.Dir = state
		door.SendRefMsg(RM_DIGDOWN, state, door.CurrX, door.CurrY, "")
	default:
		// 关闭中破损度变化
		door.Dir = state
		door.SendRefMsg(RM_TURN, state, door.CurrX, door.CurrY, "")
	}
}

// SyncWallVisual 按城墙当前状态广播墙 actor 动画（破坏→RM_DEATH，
// 修复→RM_ALIVE；客户端 ActDie+破损特效/碎片层）。
func (c *CastleObject) SyncWallVisual(server *netserver.TCPServer, engine *UserEngine) {
	env := c.castleEnv(engine)
	if env == nil {
		return
	}
	c.mu.Lock()
	ids := append([]int32(nil), c.WallActorIDs...)
	destroyed := c.WallHP <= 0
	prev := c.wallsVisualDestroyed
	if destroyed == prev {
		c.mu.Unlock()
		return
	}
	c.wallsVisualDestroyed = destroyed
	c.mu.Unlock()

	for _, id := range ids {
		wall, ok := env.getObjectByID(id).(*MonsterObject)
		if !ok || wall == nil {
			continue
		}
		if destroyed && !wall.Death {
			wall.Death = true
			wall.DeathTick = time.Now().UnixMilli()
			env.broadcastDeathMsg(wall.BaseObject, wall.ID, wall.CurrX, wall.CurrY, wall.Dir, true)
		} else if !destroyed && wall.Death {
			wall.Death = false
			wall.DeathTick = 0
			wall.SendRefMsg(RM_ALIVE, wall.Dir, wall.CurrX, wall.CurrY, "")
		}
	}
}

// HandleStructureDamage 攻城战期间玩家攻击城门/城墙时调用。
// 返回是否命中了建筑结构；命中后同步建筑可视化状态。
func (c *CastleObject) HandleStructureDamage(server *netserver.TCPServer, engine *UserEngine, x, y, damage int) bool {
	if !c.IsAtWar() {
		return false
	}
	doorX, doorY := c.DoorPosition()
	if abs(x-doorX) <= 1 && abs(y-doorY) <= 1 {
		destroyed := c.DamageDoor(damage)
		if destroyed {
			c.ToggleDoor(true)
		}
		c.SyncDoorVisual(server, engine)
		return true
	}
	for _, w := range c.WallPositions() {
		if abs(x-w[0]) <= 0 && abs(y-w[1]) <= 0 {
			c.DamageWall(damage)
			c.SyncWallVisual(server, engine)
			return true
		}
	}
	return false
}

// IsBlockedByDoor 检查指定位置是否被关闭的城门阻挡。
func (c *CastleObject) IsBlockedByDoor(x, y int) bool {
	if c.IsDoorOpen() || c.DoorHP <= 0 {
		return false
	}
	doorX, doorY := c.DoorPosition()
	return abs(x-doorX) <= 1 && abs(y-doorY) <= 0
}

func (c *CastleObject) isAttackingGuildLocked(guildName string) bool {
	for _, g := range c.AttackGuilds {
		if g == guildName {
			return true
		}
	}
	return false
}

func (c *CastleObject) ProcessCastleTick(engine *UserEngine, server *netserver.TCPServer, now int64) {
	// 城门/城墙静态 actor 只生成一次（地图加载后首 tick）
	c.mu.Lock()
	spawned := c.structuresSpawned
	c.mu.Unlock()
	if !spawned {
		c.SpawnCastleStructures(engine, now)
	}

	state := c.GetWarState()

	switch state {
	case CastleWarDeclared:
		// Delphi Castle.pas:644-679：每日 WarStartHour 点检查预约战列表
		c.checkScheduledWar(server, engine)
	case CastleWarActive:
		c.mu.Lock()
		endTick := c.WarEndTick
		c.mu.Unlock()
		if now >= endTick {
			c.EndWar(server, engine)
		} else {
			prevOwner := c.GetOwnerGuild()
			captured := c.CheckCapture(engine)
			if captured != "" {
				text := fmt.Sprintf("[攻城战] %s 被 %s 占领！", c.Config.Name, captured)
				c.broadcastSysMsg(server, engine, text)
				// 原城主行会成员：Need 7/70 城堡装备需求失效自动脱下。
				if prevOwner != "" && prevOwner != captured {
					for _, p := range engine.allPlayers() {
						if p.Ghost || p.GuildName != prevOwner {
							continue
						}
						p.checkAutoTakeOff(server)
					}
				}
			}
		}
	}

	// 税收：每 60 秒向城主行会金库收取税金
	if state == CastlePeace && c.OwnerGuild != "" {
		c.collectPeriodicTax(engine, now)
	}
}

var lastTaxTick int64

func (c *CastleObject) collectPeriodicTax(engine *UserEngine, now int64) {
	if now-lastTaxTick < 60000 {
		return
	}
	lastTaxTick = now

	rate := c.GetTaxRate()
	if rate <= 0 {
		return
	}

	var totalTax int64
	for _, p := range engine.allPlayers() {
		if p.Ghost || p.Death {
			continue
		}
		tax := int64(p.Gold) * int64(rate) / 100
		if tax > 0 && tax <= int64(p.Gold) {
			p.Gold -= int(tax)
			totalTax += tax
		}
	}
	if totalTax > 0 {
		c.CollectTax(totalTax)
	}
}

func (c *CastleObject) broadcastSysMsg(server *netserver.TCPServer, engine *UserEngine, text string) {
	msg := protocol.MakeDefaultMsg(protocol.SMSysMessage, 0, chatColorSys, 0, 0)
	for _, p := range engine.allPlayers() {
		if !p.Ghost {
			server.Send(p.Session.ID, msg, protocol.EncodeString(text))
		}
	}
}

// checkScheduledWar — Delphi Castle.pas:644-679：每日 WarStartHour 点扫描预约战列表，
// 攻击日期为当日的条目转入攻方列表并开战；无当日条目且列表清空则回到和平。
func (c *CastleObject) checkScheduledWar(server *netserver.TCPServer, engine *UserEngine) {
	t := time.Now()
	if t.Hour() != c.Config.WarStartHour {
		return
	}
	today := t.Format("2006-01-02")
	c.mu.Lock()
	if c.lastWarCheckDay == today {
		c.mu.Unlock()
		return
	}
	c.lastWarCheckDay = today
	var remaining []CastleWarDeclaration
	started := false
	for _, w := range c.WarDeclarations {
		if w.AttackDate == today {
			c.AttackGuilds = append(c.AttackGuilds, w.GuildName)
			started = true
		} else {
			remaining = append(remaining, w)
		}
	}
	c.WarDeclarations = remaining
	if !started {
		if len(c.WarDeclarations) == 0 {
			c.WarState = CastlePeace
		}
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()
	c.StartWar(server, engine)
}

// SaveState 持久化城堡状态（城主/金库/税率/城门墙血量/预约战）。
func (c *CastleObject) SaveState(db *storage.Database) {
	if db == nil {
		return
	}
	c.mu.Lock()
	rec := storage.CastleRecord{
		OwnerGuild: c.OwnerGuild,
		Gold:       c.Gold,
		DoorHP:     c.DoorHP,
		WallHP:     c.WallHP,
		WarState:   int(c.WarState),
		TaxRate:    c.TaxRate,
	}
	if len(c.WarDeclarations) > 0 {
		if buf, err := json.Marshal(c.WarDeclarations); err == nil {
			rec.Declarations = buf
		}
	}
	c.mu.Unlock()
	if err := db.SaveCastle(rec); err != nil {
		log.Logf(log.LevelError, "Castle", "failed to save castle state: %v", err)
	}
}

// LoadState 启动时恢复城堡状态（Delphi 城堡文件加载语义）。
// 进行中的战争不恢复为 Active（重启视为停战，预约战保留）。
func (c *CastleObject) LoadState(db *storage.Database) {
	if db == nil {
		return
	}
	rec, err := db.LoadCastle()
	if err != nil {
		log.Logf(log.LevelError, "Castle", "failed to load castle state: %v", err)
		return
	}
	if rec.OwnerGuild == "" && rec.Gold == 0 && len(rec.Declarations) == 0 {
		return
	}
	c.mu.Lock()
	c.OwnerGuild = rec.OwnerGuild
	c.Gold = rec.Gold
	c.TaxRate = rec.TaxRate
	if rec.DoorHP > 0 {
		c.DoorHP = rec.DoorHP
	}
	if rec.WallHP > 0 {
		c.WallHP = rec.WallHP
	}
	if len(rec.Declarations) > 0 {
		json.Unmarshal(rec.Declarations, &c.WarDeclarations)
	}
	if len(c.WarDeclarations) > 0 {
		c.WarState = CastleWarDeclared
	}
	c.mu.Unlock()
	log.Logf(log.LevelInfo, "Castle", "loaded castle state: owner=%s gold=%d declarations=%d",
		c.OwnerGuild, c.Gold, len(c.WarDeclarations))
}

func (e *UserEngine) allPlayers() []*PlayObject {
	e.mu.RLock()
	defer e.mu.RUnlock()
	players := make([]*PlayObject, 0, len(e.PlayObjectList))
	for _, p := range e.PlayObjectList {
		players = append(players, p)
	}
	return players
}
