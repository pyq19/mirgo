package main

// 特化怪物特效数据表（数据化方案，仿 magicEffectParams 模式）。
// 权威来源：Delphi AxeMon.pas 34 个特化渲染类 + HerbActor.pas 10 个物体类，
// 工厂映射 PlayScn.pas:2046-2121。本文件按 raceImg 描述特效层，
// actorbase.go 的 calcMonsterFrame/Run/drawBody 统一消费，不逐类硬移植。

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pyq0109/mirgo/internal/engine"
	"github.com/pyq0109/mirgo/internal/protocol"
	"github.com/pyq0109/mirgo/internal/wil"
)

// EffectDirMode 特效帧索引的方向组合方式。
type EffectDirMode int

const (
	DirNone   EffectDirMode = iota // 帧索引无方向分量
	Dir8x10                        // + dir*10（多数 8 方向特效）
	Dir8x8                         // + dir*8（带刀卫士 AxeMon.pas:962-968）
	Dir16x10                       // + fire16dir*10（僵尸1攻击指向目标 AxeMon.pas:970-982）
	Dir10x4                        // + dir*4（魔龙教主死亡/黄泉教主）
	Dir10x2                        // + dir*2（黄泉教主 STRUCK）
)

// WIL 文件选择特殊值；>=0 时为 resources.Mon 索引（Mon{N}.wil = Mon[N-1]）。
const (
	WilMonSelf = -1 // 与身体同一 WIL 文件
	WilDragon  = -2 // Dragon.wil
	WilEffect  = -3 // Effect.wil
)

// EffectLayer 描述一个加法混合特效层（Delphi EffectSurface/
// AttackEffectSurface/DieEffectSurface 的统一抽象）。
type EffectLayer struct {
	Wil         int           // WIL 选择（见 Wil* 常量）
	Base        int           // 帧基址
	BaseByAppr  map[int]int   // 按 appearance 覆盖基址（城墙：901→224 左墙，余→240）
	DirMode     EffectDirMode // 方向分量
	FixedFrames int           // >0 固定帧数（如雕像攻击 9 帧/石碑 20 帧）；0=随本体
	FrameTime   int           // >0 独立帧时(ms)；0=与本体同帧时
	SyncBody    bool          // true=帧随本体 CurrentFrame-StartFrame 联动；false=独立特效时钟
	StartRel    int           // SyncBody 时特效起始相对帧（触龙神：本体帧≥5 才显示）
}

// ProjectileSpec 投射物规格（Delphi TFlyingAxe/TFlyingArrow/
// TFlyingFireBall/TFlyingBug，AxeMon.pas:578-594/702-718/2019-2032/
// 2550-2563/2345-2357）。
type ProjectileSpec struct {
	TriggerMsg    int // protocol.SMFlyAxe / protocol.SMLighting
	RelFrame      int // 本体动画相对帧触发（2 或 4）
	Wil           int
	Base          int
	FrameTime     int // 0=默认 50ms
	CycleFrames   int // 飞行帧循环数（斧/虫 3 帧循环；0=不循环）
	ArrowDir      bool // 箭：帧索引用 Dir16（非 Dir16*10）
	ExplosionBase int // 命中爆炸帧基址（0=无）
	GlobalTarget  bool
}

// CastSpec 施法规格（Delphi TBanyaGuardMon.Run 帧4 NewMagic，
// AxeMon.pas:1530-1548；TDragonStatue :2494-2497）。
type CastSpec struct {
	RelFrame int      // 本体动画相对帧触发
	Kind     CastKind // 特效类型
	Sound    int      // 音效号
	AtSelf   bool     // 特效位置=施法者格（false=目标格）
	EffNum   int      // EffFly 用的魔法动画号
}

// CastKind 施法特效类型（复用 EffectManager 现有特效）。
type CastKind int

const (
	CastThunder  CastKind = iota // mtThunder 雷电（目标点）
	CastFly                      // mtFly 飞行魔法（施法者→目标）
	CastThunder2                 // mt13 TThuderEffect(140)+Magic2
)

// MonEffectSpec 一个 raceImg 的完整特化规格。
type MonEffectSpec struct {
	DeathEffect          *EffectLayer // SM_NOWDEATH 触发
	AttackEffect         *EffectLayer // SM_HIT 触发
	AttackEffectLighting *EffectLayer // SM_LIGHTING 触发
	AttackEffectFlyAxe   *EffectLayer // SM_FLYAXE 触发
	StruckEffect         *EffectLayer // SM_STRUCK 触发（黄泉教主）
	WalkEffect           *EffectLayer // SM_WALK 触发（黄泉教主）
	TurnEffect           *EffectLayer // SM_TURN 触发（石碑）
	SuppressDraw         bool         // TCatMon 系：特效加载但不绘制（race 17/19/30/31/32/74）
	HitDisablesEffect    bool         // 月魔蜘蛛：SM_HIT 强制关特效（AxeMon.pas:1794-1796）
	LightMin             int          // 攻击特效期间最低光照（火焰沃玛/沃玛教主=2）
	LockDir0             bool         // 强制 dir=0（物体类/龙类/石碑）
	// 帧段特例（Delphi 各类 CalcActorFrame 覆写）
	LightingUsesCritical bool // SM_LIGHTING 用 ActCritical（60/62/71）
	FlyAxeUsesCritical   bool // SM_FLYAXE 用 ActCritical（63）
	LightingAttackShift  int  // SM_LIGHTING 用 ActAttack+偏移帧（63=+80）
	HitUsesCritical      bool // SM_HIT 用 ActCritical（触龙神 33，HerbActor.pas:335-349）
	// 物体类行为（HerbActor.pas）
	DigUpUsesWalk   bool // SM_DIGUP 用 ActWalk 出现动画（食人花/矿石等 :155-166）
	DigDownDeletes  bool // SM_DIGDOWN 用 ActDeath 且播完自删（:199-206）
	StandNoDir      bool // SM_TURN/SM_DIGUP 帧段无方向偏移（食人花 :146-154）
	AllNoDir        bool // 所有动作无方向乘子（角蝇巢 TBeeQueen :238-292）
	BlendBody       bool // 本体整体加法混合重绘（龙身 DrawEff :803-815）
	FixedIdleFrames bool // 固定 2 帧动画 0..1@400ms（龙身 CalcActorFrame :781-801）
	// 城门/城墙（TCastleDoor/TCastleWall，HerbActor.pas:365-683）
	DigUpUsesAttack     bool // SM_DIGUP=开门 ActAttack(56)（城门 :442-449）
	DigDownUsesCritical bool // SM_DIGDOWN=关门 ActCritical(64)（城门 :450-459）
	DoorTurnBehavior    bool // SM_TURN：dir<3 关闭破损帧 / dir>=3 开启定格（城门 :467-488）
	TurnDirAdditive     bool // SM_TURN 帧索引用 +dir 而非 dir*(frame+skip)（城墙 :605-614）
	BrokenOverlay       bool // 死亡时常驻碎片层（城墙 BrokenSurface 帧 8+dir :618-636）
	Projectile      *ProjectileSpec
	Cast            *CastSpec
}

// monEffectTable raceImg→特化规格。
// 死亡特效族（SM_NOWDEATH，加法混合）——AxeMon.pas 逐类核实：
//   - 14/15: TSkeletonOma.LoadSurface Mon3[340+rel]（:437-441，race 22 除外由标志门控）
//   - 18: THuSuABi Mon3[2860+rel]（:782-787，DEATHFIREEFFECTBASE）
//   - 36: TExplosionSpider Mon14[730+rel] 自爆与死亡动画同步（:1797-1815）
//   - 37: TFlyingSpider Mon12[1420] 20帧死亡烟花（:1779-1785；Delphi 为地图
//     特效层 TNormalDrawEffect，Go 近似为本体叠加层）
//   - 40: TGasKuDeGi 死亡 Mon5[340+rel]（:970-982）
//   - 62: Mon19[1470+rel] 无方向（:1680-1721）
//   - 63: Mon20[3400+dir*4+rel]（:1910-1971）
//   - 65/66/68: Mon20[350/1600/1600+rel]（:1004-1027）
//   - 67: Mon20[1160+dir*10+rel] 唯一带方向死亡特效（:1016-1021）
//   - 69: TSkeletonArcherMon bo260 Mon20[1600+rel]（:1742-1750）
//   - 70/71: Mon21[2320+rel / 2870+dir*10+rel]（:1425-1443）
//   - 75/77: TStoneMonster Mon22[2530/2660+rel] 20帧@80ms 独立时钟（:2109-2120）
//   - 78: Mon22[3120+dir*4+rel]（:1425-1443）
//   - 81: TAngel 继承 TSkeletonArcherMon 默认 Mon20[1600+rel]
var monEffectTable = map[int]*MonEffectSpec{
	14: {DeathEffect: &EffectLayer{Wil: 2, Base: 340, SyncBody: true}},
	15: {DeathEffect: &EffectLayer{Wil: 2, Base: 340, SyncBody: true},
		Projectile: &ProjectileSpec{TriggerMsg: protocol.SMFlyAxe, RelFrame: 2, Wil: 2, Base: 447, CycleFrames: 3}}, // 飞斧 FLYOMAAXEBASE
	16: {AttackEffect: &EffectLayer{Wil: 2, Base: 1444, DirMode: Dir8x10, SyncBody: true}},                       // 洞蛆喷气 KUDEGIGASBASE-1
	17: {SuppressDraw: true},                                                                                      // TCatMon：特效不绘制
	18: {DeathEffect: &EffectLayer{Wil: 2, Base: 2860, SyncBody: true}},                                          // 稻草人火焰
	19: {SuppressDraw: true},
	20: {AttackEffect: &EffectLayer{Wil: 3, Base: 1800, DirMode: Dir8x10, SyncBody: true}, LightMin: 2}, // 火焰沃玛 COWMONFIREBASE
	21: {AttackEffect: &EffectLayer{Wil: 3, Base: 1900, DirMode: Dir8x10, SyncBody: true}, LightMin: 2}, // 沃玛教主 COWMONLIGHTBASE
	22: {Projectile: &ProjectileSpec{TriggerMsg: protocol.SMFlyAxe, RelFrame: 2, Wil: 2, Base: 2967, CycleFrames: 3}}, // 暗黑战士荆棘 THORNBASE
	24: {AttackEffect: &EffectLayer{Wil: 0, Base: 760, DirMode: Dir8x8, SyncBody: true}},                               // 带刀卫士 SUPERIORGUARDBASE
	30: {SuppressDraw: true},
	31: {SuppressDraw: true},
	32: {SuppressDraw: true},
	36: {DeathEffect: &EffectLayer{Wil: 13, Base: 730, SyncBody: true},
		HitDisablesEffect: true},
	37: {DeathEffect: &EffectLayer{Wil: 11, Base: 1420, FixedFrames: 20, FrameTime: 140, SyncBody: false}},
	40: {DeathEffect: &EffectLayer{Wil: 4, Base: 340, SyncBody: true},
		AttackEffect:         &EffectLayer{Wil: 4, Base: 350, DirMode: Dir16x10, SyncBody: true}, // ZOMBILIGHTINGBASE 指向目标
		AttackEffectLighting: &EffectLayer{Wil: 4, Base: 350, DirMode: Dir16x10, SyncBody: true}},
	45: {Projectile: &ProjectileSpec{TriggerMsg: protocol.SMFlyAxe, RelFrame: 4, Wil: WilEffect, Base: 272, FrameTime: 30, ArrowDir: true}}, // 弓箭手 ARCHERBASE2
	49: {AttackEffect: &EffectLayer{Wil: 6, Base: 1680, DirMode: Dir8x10, FixedFrames: 9, SyncBody: true}},                                  // 祖玛教主 SCULPTUREFIREBASE 固定9帧
	52: {AttackEffect: &EffectLayer{Wil: 3, Base: 3590, DirMode: Dir8x10, SyncBody: true}},                                                  // 楔蛾 MOTHPOISONGASBASE
	53: {AttackEffect: &EffectLayer{Wil: 2, Base: 3590, DirMode: Dir8x10, SyncBody: true}},                                                  // 粪虫 DUNGPOISONGASBASE
	55: {AttackEffect: &EffectLayer{Wil: 17, Base: 821, DirMode: Dir8x10, FixedFrames: 5, FrameTime: 100, SyncBody: false}},                 // 神兽 WARRIORELFFIREBASE+1（地图特效近似）
	60: {AttackEffectLighting: &EffectLayer{Wil: 18, Base: 430, DirMode: Dir8x10, SyncBody: true},
		LightingUsesCritical: true}, // 虹魔蝎卫 SM_LIGHTING 用 ActCritical
	61: {AttackEffect: &EffectLayer{Wil: 18, Base: 860, DirMode: Dir8x10, SyncBody: true}}, // 白野猪王
	62: {DeathEffect: &EffectLayer{Wil: 18, Base: 1470, SyncBody: true},
		AttackEffect:         &EffectLayer{Wil: 18, Base: 1490, DirMode: Dir8x10, SyncBody: true},
		AttackEffectLighting: &EffectLayer{Wil: 18, Base: 1380, DirMode: Dir8x10, SyncBody: true},
		LightingUsesCritical: true}, // 虹魔教主
	63: {DeathEffect: &EffectLayer{Wil: 19, Base: 3400, DirMode: Dir10x4, SyncBody: true},
		AttackEffect:         &EffectLayer{Wil: 19, Base: 3140, DirMode: Dir8x10, SyncBody: true},
		AttackEffectLighting: &EffectLayer{Wil: 19, Base: 3220, DirMode: Dir8x10, SyncBody: true},
		AttackEffectFlyAxe:   &EffectLayer{Wil: 19, Base: 3300, DirMode: Dir8x10, SyncBody: true},
		StruckEffect:         &EffectLayer{Wil: 19, Base: 3380, DirMode: Dir10x2, SyncBody: true},
		WalkEffect:           &EffectLayer{Wil: 19, Base: 3060, DirMode: Dir8x10, SyncBody: true},
		LightingUsesCritical: false,
		LightingAttackShift:  80, // SM_LIGHTING 用 ActAttack+80
		FlyAxeUsesCritical:   true,
		Projectile:           &ProjectileSpec{TriggerMsg: protocol.SMFlyAxe, RelFrame: 4, Wil: 19, Base: 3573, FrameTime: 40}}, // 火球
	64: {AttackEffect: &EffectLayer{Wil: 19, Base: 720, DirMode: Dir8x10, SyncBody: true}}, // 恶灵尸王
	65: {DeathEffect: &EffectLayer{Wil: 19, Base: 350, SyncBody: true}},
	66: {DeathEffect: &EffectLayer{Wil: 19, Base: 1600, SyncBody: true}},
	67: {DeathEffect: &EffectLayer{Wil: 19, Base: 1160, DirMode: Dir8x10, SyncBody: true}},
	68: {DeathEffect: &EffectLayer{Wil: 19, Base: 1600, SyncBody: true}},
	69: {DeathEffect: &EffectLayer{Wil: 19, Base: 1600, SyncBody: true},
		Projectile: &ProjectileSpec{TriggerMsg: protocol.SMFlyAxe, RelFrame: 4, Wil: WilEffect, Base: 272, FrameTime: 30, ArrowDir: true}}, // 骷髅弓箭手
	70: {DeathEffect: &EffectLayer{Wil: 20, Base: 2320, SyncBody: true},
		AttackEffect: &EffectLayer{Wil: 20, Base: 2230, DirMode: Dir8x10, SyncBody: true},
		Cast:         &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 10112}}, // 牛魔法师雷电
	71: {DeathEffect: &EffectLayer{Wil: 20, Base: 2870, DirMode: Dir8x10, SyncBody: true},
		AttackEffect:         &EffectLayer{Wil: 20, Base: 2780, DirMode: Dir8x10, SyncBody: true},
		AttackEffectLighting: &EffectLayer{Wil: 20, Base: 2960, DirMode: Dir8x10, SyncBody: true},
		LightingUsesCritical: true,
		Cast:                 &CastSpec{RelFrame: 4, Kind: CastFly, Sound: 10012, EffNum: 1}}, // 牛魔祭司飞行魔法
	72: {AttackEffect: &EffectLayer{Wil: 20, Base: 3490, DirMode: Dir8x10, SyncBody: true},
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder2, Sound: 2276}}, // 牛魔王 mt13（目标点）
	73: {SuppressDraw: true,
		Projectile: &ProjectileSpec{TriggerMsg: protocol.SMFlyAxe, RelFrame: 4, Wil: 21, Base: 350, FrameTime: 50, CycleFrames: 3, ExplosionBase: 430}}, // 飞虫+爆炸
	74: {SuppressDraw: true},
	75: {DeathEffect: &EffectLayer{Wil: 21, Base: 2530, FixedFrames: 20, FrameTime: 80, SyncBody: false},
		AttackEffect: &EffectLayer{Wil: 21, Base: 2500, FixedFrames: 26, FrameTime: 150, SyncBody: false},
		TurnEffect:   &EffectLayer{Wil: 21, Base: 2490, FixedFrames: 26, FrameTime: 300, SyncBody: false},
		LockDir0:     true}, // 魔龙石碑
	77: {DeathEffect: &EffectLayer{Wil: 21, Base: 2660, FixedFrames: 20, FrameTime: 80, SyncBody: false},
		AttackEffect: &EffectLayer{Wil: 21, Base: 2630, FixedFrames: 26, FrameTime: 150, SyncBody: false},
		TurnEffect:   &EffectLayer{Wil: 21, Base: 2620, FixedFrames: 26, FrameTime: 300, SyncBody: false},
		LockDir0:     true},
	78: {DeathEffect: &EffectLayer{Wil: 21, Base: 3120, DirMode: Dir10x4, SyncBody: true},
		AttackEffect: &EffectLayer{Wil: 21, Base: 3440, DirMode: Dir8x10, SyncBody: true},
		Cast:         &CastSpec{RelFrame: 4, Kind: CastThunder2, Sound: 2396, AtSelf: true}}, // 魔龙教主 mt13 自身
	79: {SuppressDraw: true,
		Projectile: &ProjectileSpec{TriggerMsg: protocol.SMFlyAxe, RelFrame: 4, Wil: 21, Base: 1989, FrameTime: 50, CycleFrames: 3, GlobalTarget: true}},
	81: {DeathEffect: &EffectLayer{Wil: 19, Base: 1600, SyncBody: true},
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 10112}}, // 天使雷电
	83: {AttackEffect: &EffectLayer{Wil: WilDragon, Base: 60, SyncBody: true},
		LockDir0: true}, // 火龙教主（本体动作 81/82/83 特例未移植，近似处理）
	84: {AttackEffectLighting: &EffectLayer{Wil: WilDragon, Base: 310, SyncBody: true}, LockDir0: true,
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 8222}}, // 龙柱雷电
	85: {AttackEffectLighting: &EffectLayer{Wil: WilDragon, Base: 310, SyncBody: true}, LockDir0: true,
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 8222}},
	86: {AttackEffectLighting: &EffectLayer{Wil: WilDragon, Base: 310, SyncBody: true}, LockDir0: true,
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 8222}},
	87: {AttackEffectLighting: &EffectLayer{Wil: WilDragon, Base: 330, SyncBody: true}, LockDir0: true,
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 8222}},
	88: {AttackEffectLighting: &EffectLayer{Wil: WilDragon, Base: 330, SyncBody: true}, LockDir0: true,
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 8222}},
	89: {AttackEffectLighting: &EffectLayer{Wil: WilDragon, Base: 330, SyncBody: true}, LockDir0: true,
		Cast: &CastSpec{RelFrame: 4, Kind: CastThunder, Sound: 8222}},
	// 物体类（HerbActor.pas，服务端 race 85/103/107/115/116 等）
	13: {DigUpUsesWalk: true, DigDownDeletes: true, StandNoDir: true}, // 食人花 TKillingHerb
	33: {AttackEffect: &EffectLayer{Wil: 14, Base: 100, FixedFrames: 10, FrameTime: 62, SyncBody: true, StartRel: 5},
		HitUsesCritical: true, LockDir0: true, DigUpUsesWalk: true, DigDownDeletes: true}, // 触龙神
	34: {LockDir0: true, DigUpUsesWalk: true, DigDownDeletes: true},                       // 赤月恶魔
	35: {LockDir0: true, DigUpUsesWalk: true, DigDownDeletes: true},                       // 幻影蜘蛛
	43: {StandNoDir: true, AllNoDir: true},                                                // 角蝇巢：所有动作无方向乘子（HerbActor.pas:238-292）
	80: {DigUpUsesWalk: true, DigDownDeletes: true, StandNoDir: true},                     // 矿石 TMineMon
	90: {LockDir0: true, BlendBody: true, FixedIdleFrames: true},                          // 龙身 TDragonBody
	// 城门/城墙（服务端 race 110/111，帧布局 Effect.wil nrace90 槽位）
	98: {DeathEffect: &EffectLayer{Wil: WilEffect, Base: 240,
		BaseByAppr: map[int]int{901: 224}, SyncBody: true}, // WALLLEFTBROKENEFFECTBASE=224/RIGHT=240
		DigUpUsesAttack: true, // 修复动画播 ActDie 近似（Delphi SM_DIGUP=破坏动画复用）
		TurnDirAdditive: true,
		BrokenOverlay:   true},
	99: {DeathEffect: &EffectLayer{Wil: WilEffect, Base: 120, FixedFrames: 10, SyncBody: true}, // DOORDEATHEFFECTBASE
		DigUpUsesAttack:     true,
		DigDownUsesCritical: true,
		DoorTurnBehavior:    true},
}

func monEffectSpec(raceImg int) *MonEffectSpec {
	return monEffectTable[raceImg]
}

// effectDirOffset 按方向模式计算帧索引分量。
func effectDirOffset(mode EffectDirMode, dir, fire16dir int) int {
	switch mode {
	case Dir8x10:
		return dir * 10
	case Dir8x8:
		return dir * 8
	case Dir16x10:
		return fire16dir * 10
	case Dir10x4:
		return dir * 4
	case Dir10x2:
		return dir * 2
	}
	return 0
}

// effectWilFile 按 EffectLayer.Wil 选择 WIL 文件。
func (a *Actor) effectWilFile(spec *EffectLayer, resources *engine.ResourceManager) *wil.File {
	switch {
	case spec.Wil == WilMonSelf:
		return a.getMonFile(resources)
	case spec.Wil == WilDragon:
		return resources.Dragon
	case spec.Wil == WilEffect:
		return resources.Effect
	case spec.Wil >= 0 && spec.Wil < len(resources.Mon):
		return resources.Mon[spec.Wil]
	}
	return nil
}

// drawMonsterEffect 绘制加法混合特效层（Delphi AxeMon.pas DrawChr/DrawEff
// 的 DrawBlend(...,1)，AxeMon.pas:526-532/1117-1139）。位置用特效图自身
// 偏移（HotX/HotY，即 Delphi GetCachedImage 返回的 ax/ay）。
func (a *Actor) drawMonsterEffect(glState *engine.GLState, resources *engine.ResourceManager, screenX, screenY float32, proj [16]float32) {
	spec := a.EffSpec
	if spec == nil || a.SuppressEffDraw {
		return
	}
	file := a.effectWilFile(spec, resources)
	if file == nil {
		return
	}
	rel := a.EffFrame
	if spec.SyncBody {
		rel = a.CurrentFrame - a.StartFrame - spec.StartRel
	}
	if rel < 0 {
		return
	}
	if spec.FixedFrames > 0 && rel >= spec.FixedFrames {
		return
	}
	base := spec.Base
	if len(spec.BaseByAppr) > 0 {
		if b, ok := spec.BaseByAppr[a.Appearance]; ok {
			base = b
		}
	}
	idx := base + effectDirOffset(spec.DirMode, a.Dir, a.Fire16Dir) + rel
	if idx < 0 || idx >= file.Count {
		return
	}
	img := file.GetImage(idx)
	if img == nil || img.RGBA == nil {
		return
	}
	tex := resources.GetTexture(file, idx)
	if tex == 0 {
		return
	}
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	glState.DrawQuad(tex, screenX+float32(img.HotX), screenY+float32(img.HotY),
		float32(img.Width), float32(img.Height), proj)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
