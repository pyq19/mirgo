package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pyq0109/mirgo/internal/protocol"
	"github.com/pyq0109/mirgo/internal/wil"
)

// 批3回归：特化怪物死亡特效族（Delphi AxeMon.pas）。

// TestScarecrowDeathFireEffect 稻草人（raceImg 18）死亡火焰特效
//（THuSuABi，AxeMon.pas:782-787，DEATHFIREEFFECTBASE=2860）：
// SM_NOWDEATH 开启、SM_DEATH 直接尸体不开、其它种族不开。
func TestScarecrowDeathFireEffect(t *testing.T) {
	newScarecrow := func() *Actor {
		a := NewActor(1, 0, 0, 0)
		a.Type = ActorMonster
		a.Race = 18
		a.Appearance = 27
		a.MonAction = GetRaceByPM(18, 27)
		return a
	}

	a := newScarecrow()
	a.CurrentAction = protocol.SMNowDeath
	a.calcMonsterFrame()
	if !a.UseEffect || a.EffSpec == nil {
		t.Fatal("稻草人 SM_NOWDEATH 应开启死亡特效")
	}
	if a.EffSpec.Base != 2860 || a.EffSpec.Wil != 2 {
		t.Errorf("特效规格 = Wil%d Base%d, want Wil2(Mon3.wil) Base2860", a.EffSpec.Wil, a.EffSpec.Base)
	}
	if !a.EffSpec.SyncBody {
		t.Error("稻草人火焰特效应随本体帧联动")
	}

	// SM_DEATH（后进入视野直接尸体）不播特效
	b := newScarecrow()
	b.CurrentAction = protocol.SMDeath
	b.calcMonsterFrame()
	if b.UseEffect {
		t.Fatal("SM_DEATH 直接尸体不应开启特效")
	}

	// 无特效表的种族不开
	c := NewActor(2, 0, 0, 0)
	c.Type = ActorMonster
	c.Race = 19
	c.Appearance = 100
	c.MonAction = GetRaceByPM(19, 100)
	c.CurrentAction = protocol.SMNowDeath
	c.calcMonsterFrame()
	if c.UseEffect {
		t.Fatal("半兽人（raceImg 19）无死亡特效表项，不应开启特效")
	}
}

// TestDeathEffectFrameIndex 特效帧索引随本体帧联动（AxeMon.pas:786
// DEATHFIREEFFECTBASE + m_nCurrentFrame - m_nStartFrame）。
func TestDeathEffectFrameIndex(t *testing.T) {
	a := NewActor(1, 0, 0, 0)
	a.Type = ActorMonster
	a.Race = 18
	a.Appearance = 27
	a.MonAction = GetRaceByPM(18, 27)
	a.CurrentAction = protocol.SMNowDeath
	a.calcMonsterFrame()

	if a.CurrentFrame != a.StartFrame {
		t.Fatalf("动画起始 CurrentFrame = %d, want StartFrame %d", a.CurrentFrame, a.StartFrame)
	}
	// 模拟播到第 k 帧
	for k := 0; k < 10; k++ {
		a.CurrentFrame = a.StartFrame + k
		rel := a.CurrentFrame - a.StartFrame
		idx := a.EffSpec.Base + effectDirOffset(a.EffSpec.DirMode, a.Dir, a.Fire16Dir) + rel
		if want := 2860 + k; idx != want {
			t.Errorf("第 %d 帧特效索引 = %d, want %d", k, idx, want)
		}
	}
}

// TestIndependentClockDeathEffect 独立时钟特效（石碑 raceImg 75：
// 20帧@80ms，TStoneMonster AxeMon.pas:2109-2120）。
func TestIndependentClockDeathEffect(t *testing.T) {
	spec := monEffectSpec(75)
	if spec == nil || spec.DeathEffect == nil {
		t.Fatal("raceImg 75 应有死亡特效")
	}
	de := spec.DeathEffect
	if de.SyncBody {
		t.Fatal("石碑死亡特效应为独立时钟")
	}
	if de.FixedFrames != 20 || de.FrameTime != 80 {
		t.Errorf("石碑死亡特效 = %d帧@%dms, want 20帧@80ms", de.FixedFrames, de.FrameTime)
	}
}

// TestAttackEffectTrigger 攻击特效按消息触发（Delphi TGasKuDeGi/
// TBanyaGuardMon 等 CalcActorFrame 置 m_boUseEffect）。
func TestAttackEffectTrigger(t *testing.T) {
	newMon := func(race, appr, action int) *Actor {
		a := NewActor(1, 0, 0, 0)
		a.Type = ActorMonster
		a.Race = race
		a.Appearance = appr
		a.MonAction = GetRaceByPM(race, appr)
		a.CurrentAction = action
		a.calcMonsterFrame()
		return a
	}

	// 火焰沃玛 SMHit → Mon4[1800+dir*10]
	a := newMon(20, 31, protocol.SMHit)
	if !a.UseEffect || a.EffSpec == nil || a.EffSpec.Base != 1800 || a.EffSpec.Wil != 3 {
		t.Fatal("火焰沃玛 SMHit 应开启 Mon4 1800 攻击特效")
	}
	// 僵尸1 SMLighting → Mon5[350+fire16dir*10]
	b := newMon(40, 110, protocol.SMLighting)
	if !b.UseEffect || b.EffSpec == nil || b.EffSpec.Base != 350 || b.EffSpec.DirMode != Dir16x10 {
		t.Fatal("僵尸1 SMLighting 应开启 Mon5 350 Dir16x10 攻击特效")
	}
	// 月魔蜘蛛 SMHit 强制关特效（AxeMon.pas:1794-1796）
	c := newMon(36, 120, protocol.SMHit)
	if c.UseEffect {
		t.Fatal("月魔蜘蛛 SMHit 不应开启特效（HitDisablesEffect）")
	}
	// 虹魔蝎卫 SMLighting 用 ActCritical 帧段
	d := newMon(60, 141, protocol.SMLighting)
	wantStart, _ := CalcFrame(d.MonAction.ActCritical, d.Dir)
	if d.StartFrame != wantStart {
		t.Errorf("虹魔蝎卫 SMLighting StartFrame = %d, want ActCritical %d", d.StartFrame, wantStart)
	}
	if !d.UseEffect || d.EffSpec == nil || d.EffSpec.Base != 430 {
		t.Fatal("虹魔蝎卫 SMLighting 应开启 Mon19 430 特效")
	}
	// 黄泉教主 SMLighting 用 ActAttack+80 帧偏移
	e := newMon(63, 148, protocol.SMLighting)
	atkStart, _ := CalcFrame(e.MonAction.ActAttack, e.Dir)
	if e.StartFrame != atkStart+80 {
		t.Errorf("黄泉教主 SMLighting StartFrame = %d, want ActAttack+80 = %d", e.StartFrame, atkStart+80)
	}
	// SMDeath 直接尸体不开攻击特效
	f := newMon(20, 31, protocol.SMDeath)
	if f.UseEffect {
		t.Fatal("SMDeath 不应开启特效")
	}
}

// TestProjectileSpawnFrames 投射物在指定相对帧触发（Delphi RunFrameAction
// NewFlyObject：掷斧骷髅 rel2、弓箭手 rel4，AxeMon.pas:578-594/702-718）。
func TestProjectileSpawnFrames(t *testing.T) {
	newMon := func(race, appr int) *Actor {
		a := NewActor(1, 0, 0, 0)
		a.Type = ActorMonster
		a.Race = race
		a.Appearance = appr
		a.MonAction = GetRaceByPM(race, appr)
		a.CurrentAction = protocol.SMFlyAxe
		a.calcMonsterFrame()
		return a
	}

	// 掷斧骷髅 rel=2
	a := newMon(15, 21)
	for k := 0; k < 2; k++ {
		a.CurrentFrame = a.StartFrame + k
		a.effSpawned = false
		a.PendingProjectile = nil
		a.checkEffectSpawns()
		if a.PendingProjectile != nil {
			t.Fatalf("掷斧骷髅 rel=%d 不应触发投射物", k)
		}
	}
	a.CurrentFrame = a.StartFrame + 2
	a.effSpawned = false
	a.PendingProjectile = nil
	a.checkEffectSpawns()
	if a.PendingProjectile == nil || a.PendingProjectile.Base != 447 {
		t.Fatal("掷斧骷髅 rel=2 应触发 447 飞斧")
	}
	// 同一动作只触发一次：PlayScene 消费后不再重复产生
	a.PendingProjectile = nil
	a.checkEffectSpawns()
	if a.PendingProjectile != nil {
		t.Fatal("投射物同一动作只应触发一次")
	}

	// 弓箭手 rel=4
	b := newMon(45, 110)
	b.CurrentFrame = b.StartFrame + 4
	b.checkEffectSpawns()
	if b.PendingProjectile == nil || b.PendingProjectile.Base != 272 || !b.PendingProjectile.ArrowDir {
		t.Fatal("弓箭手 rel=4 应触发 272 箭（ArrowDir）")
	}
}

// TestObjectClassBehaviors 物体类行为（HerbActor.pas）：
// 食人花出土/回潜自删、触龙神 ActCritical+延迟特效、角蝇巢无方向。
func TestObjectClassBehaviors(t *testing.T) {
	newMon := func(race, appr, action int) *Actor {
		a := NewActor(1, 0, 0, 0)
		a.Type = ActorMonster
		a.Race = race
		a.Appearance = appr
		a.MonAction = GetRaceByPM(race, appr)
		a.CurrentAction = action
		a.calcMonsterFrame()
		return a
	}

	// 食人花 SM_DIGUP = ActWalk 出现动画（HerbActor.pas:155-166）
	herb := newMon(13, 10, protocol.SMDigUp)
	wantStart, wantEnd := CalcFrame(herb.MonAction.ActWalk, 0) // StandNoDir → 无方向
	if herb.StartFrame != wantStart || herb.EndFrame != wantEnd {
		t.Errorf("食人花 DIGUP 帧 = %d..%d, want ActWalk %d..%d", herb.StartFrame, herb.EndFrame, wantStart, wantEnd)
	}

	// 食人花 SM_DIGDOWN = ActDeath + 播完自删（HerbActor.pas:199-206）
	herb2 := newMon(13, 10, protocol.SMDigDown)
	if !herb2.DelAfterAction {
		t.Fatal("食人花 DIGDOWN 应置 DelAfterAction")
	}
	// 播完动画 → DelActor
	now := int64(1000)
	herb2.LastFrameTick = now
	for i := 0; i < 15 && !herb2.DelActor; i++ {
		now += 150
		herb2.Run(now)
	}
	if !herb2.DelActor {
		t.Fatal("DIGDOWN 动画播完应置 DelActor（ActorManager 随后移除）")
	}

	// 触龙神 SM_HIT 用 ActCritical + 特效本体帧≥5 才显示
	centi := newMon(33, 140, protocol.SMHit)
	critStart, _ := CalcFrame(centi.MonAction.ActCritical, 0) // LockDir0
	if centi.StartFrame != critStart {
		t.Errorf("触龙神 SMHit StartFrame = %d, want ActCritical %d", centi.StartFrame, critStart)
	}
	if !centi.UseEffect || centi.EffSpec == nil || centi.EffSpec.StartRel != 5 {
		t.Fatal("触龙神 SMHit 应开启 Mon15[100] 特效且 StartRel=5")
	}

	// 角蝇巢所有动作无方向乘子（HerbActor.pas:238-292）
	hive := newMon(43, 41, protocol.SMHit)
	hiveNoDir, _ := CalcFrame(hive.MonAction.ActAttack, 0)
	hiveWithDir, _ := CalcFrame(hive.MonAction.ActAttack, 3)
	if hive.StartFrame != hiveNoDir {
		t.Errorf("角蝇巢 SMHit StartFrame = %d, want 无方向 %d（有方向为 %d）", hive.StartFrame, hiveNoDir, hiveWithDir)
	}
}

// TestCastleFrameBehavior 城门/城墙帧行为（HerbActor.pas:410-666）：
// 城门 DIGUP=开门 ActAttack、DIGDOWN=关门 ActCritical、TURN dir>=3 开启定格；
// 城墙 TURN 帧索引 +dir（加法）。
func TestCastleFrameBehavior(t *testing.T) {
	newMon := func(race, appr, action int) *Actor {
		a := NewActor(1, 0, 0, 0)
		a.Type = ActorMonster
		a.Race = race
		a.Appearance = appr
		a.MonAction = GetRaceByPM(race, appr)
		a.CurrentAction = action
		a.calcMonsterFrame()
		return a
	}

	// 城门开门动画 ActAttack(56)
	door := newMon(99, 900, protocol.SMDigUp)
	want, _ := CalcFrame(door.MonAction.ActAttack, door.Dir)
	if door.StartFrame != want {
		t.Errorf("城门 DIGUP StartFrame = %d, want ActAttack %d", door.StartFrame, want)
	}

	// 城门关门动画 ActCritical(64)
	door2 := newMon(99, 900, protocol.SMDigDown)
	want2, _ := CalcFrame(door2.MonAction.ActCritical, door2.Dir)
	if door2.StartFrame != want2 {
		t.Errorf("城门 DIGDOWN StartFrame = %d, want ActCritical %d", door2.StartFrame, want2)
	}

	// 城门静态开启态（dir>=3）定格 ActCritical.start
	door3 := newMon(99, 900, protocol.SMTurn)
	door3.Dir = 7
	door3.calcMonsterFrame()
	if door3.StartFrame != door3.MonAction.ActCritical.Start || door3.EndFrame != door3.StartFrame {
		t.Errorf("城门开启定格 = %d..%d, want %d 单帧", door3.StartFrame, door3.EndFrame, door3.MonAction.ActCritical.Start)
	}

	// 城门静态关闭破损态 dir<3 → ActStand.start + dir*(1+7)
	door4 := newMon(99, 900, protocol.SMTurn)
	door4.Dir = 2
	door4.calcMonsterFrame()
	wantStand, _ := CalcFrame(door4.MonAction.ActStand, 2)
	if door4.StartFrame != wantStand {
		t.Errorf("城门关闭破损帧 = %d, want %d", door4.StartFrame, wantStand)
	}

	// 城墙静态帧：ActStand.start + dir（加法，非乘法）
	wall := newMon(98, 901, protocol.SMTurn)
	wall.Dir = 2
	wall.calcMonsterFrame()
	if wall.StartFrame != wall.MonAction.ActStand.Start+2 {
		t.Errorf("城墙静态帧 = %d, want ActStand.start+2 = %d", wall.StartFrame, wall.MonAction.ActStand.Start+2)
	}
	wallSpec := monEffectSpec(98)
	if wallSpec == nil || wallSpec.DeathEffect == nil || !wallSpec.BrokenOverlay {
		t.Fatal("城墙应有破坏特效与碎片层规格")
	}
	spec := monEffectSpec(99)
	if spec == nil || spec.DeathEffect == nil || spec.DeathEffect.Base != 120 {
		t.Fatal("城门应有 DOORDEATHEFFECTBASE=120 破坏特效")
	}
}

// TestMonEffectTableConsistency 表内条目结构自检：
// 独立时钟必须有帧时，SyncBody 不要求；方向模式合法。
func TestMonEffectTableConsistency(t *testing.T) {
	for race, spec := range monEffectTable {
		layers := []*EffectLayer{spec.DeathEffect, spec.AttackEffect,
			spec.AttackEffectLighting, spec.AttackEffectFlyAxe,
			spec.StruckEffect, spec.WalkEffect, spec.TurnEffect}
		for _, l := range layers {
			if l == nil {
				continue
			}
			if !l.SyncBody && l.FrameTime <= 0 {
				t.Errorf("raceImg %d: 独立时钟特效缺少 FrameTime", race)
			}
			if l.DirMode < DirNone || l.DirMode > Dir10x2 {
				t.Errorf("raceImg %d: 非法方向模式 %d", race, l.DirMode)
			}
		}
	}
}

// TestMonEffectTableWilBounds 启动自检：每个特效基址+最大方向偏移+帧数
// 必须在对应 WIL 的 imgCount 范围内，防越界读串图。
// 依赖 asset/client/Data 素材，缺失时跳过。
func TestMonEffectTableWilBounds(t *testing.T) {
	dataDir := filepath.Join("..", "..", "asset", "client", "Data")
	if _, err := os.Stat(dataDir); err != nil {
		t.Skip("素材目录不存在，跳过 WIL 边界自检")
	}
	load := func(name string) *wil.File {
		f, err := wil.Load(filepath.Join(dataDir, name))
		if err != nil {
			t.Logf("加载 %s 失败（表内相关条目跳过自检）: %v", name, err)
			return nil
		}
		return f
	}
	mons := make(map[int]*wil.File)
	for _, i := range []int{0, 2, 3, 4, 6, 11, 13, 17, 18, 19, 20, 21} {
		mons[i] = load(fmt.Sprintf("Mon%d.wil", i+1))
	}
	dragon := load("Dragon.wil")
	effect := load("Effect.wil")

	maxDirOff := func(mode EffectDirMode) int {
		switch mode {
		case Dir8x10:
			return 7 * 10
		case Dir8x8:
			return 7 * 8
		case Dir16x10:
			return 15 * 10
		case Dir10x4:
			return 7 * 4
		case Dir10x2:
			return 7 * 2
		}
		return 0
	}
	for race, spec := range monEffectTable {
		layers := []*EffectLayer{spec.DeathEffect, spec.AttackEffect,
			spec.AttackEffectLighting, spec.AttackEffectFlyAxe,
			spec.StruckEffect, spec.WalkEffect, spec.TurnEffect}
		for _, l := range layers {
			if l == nil {
				continue
			}
			var f *wil.File
			switch {
			case l.Wil == WilDragon:
				f = dragon
			case l.Wil == WilEffect:
				f = effect
			case l.Wil == WilMonSelf:
				continue // 与身体同文件，帧段随本体模板已验证
			default:
				f = mons[l.Wil]
			}
			if f == nil {
				t.Logf("raceImg %d: WIL 索引 %d 加载失败，跳过该条目自检", race, l.Wil)
				continue
			}
			frames := l.FixedFrames
			if frames == 0 {
				frames = 10 // SyncBody 死亡特效上限 = ActDie 10 帧
			}
			maxIdx := l.Base + maxDirOff(l.DirMode) + frames - 1
			if maxIdx >= f.Count {
				t.Errorf("raceImg %d: 特效帧越界 maxIdx=%d >= imgCount=%d (Wil%d Base%d dirMode%d frames%d)",
					race, maxIdx, f.Count, l.Wil, l.Base, l.DirMode, frames)
			}
		}
	}
}
