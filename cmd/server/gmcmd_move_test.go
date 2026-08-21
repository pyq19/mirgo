package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// 传送与人员管理族测试（gmcmd_move.go）。
// Delphi 参照行号见各命令注释。

// newMoveFamilyFixture GM + 地图 "0" + 引擎的最小可运行环境。
func newMoveFamilyFixture(t *testing.T) (*PlayObject, *Environment, *UserEngine, *MapManager) {
	t.Helper()
	env := newTestEnv(20, 20)
	env.Name = "0"
	mm := NewMapManager("")
	mm.maps["0"] = env
	engine := NewUserEngine(nil, mm)
	engine.Config = &ServerConfig{}
	gm := &PlayObject{BaseObject: NewBaseObject("gm", 1)}
	gm.Session = &netserver.Session{ID: 1}
	gm.Permission = 10
	gm.Engine = engine
	gm.MapMgr = mm
	gm.envir = env
	gm.MapName = env.Name
	gm.CurrX, gm.CurrY = 5, 5
	gm.WAbil.HP, gm.WAbil.MaxHP = 100, 100
	env.AddObject(5, 5, OS_MOVINGOBJECT, gm)
	engine.AddPlayer(gm)
	return gm, env, engine, mm
}

func newMoveFamilyVictim(t *testing.T, name string, id int32, env *Environment, engine *UserEngine, x, y int) *PlayObject {
	t.Helper()
	pl := &PlayObject{BaseObject: NewBaseObject(name, id)}
	pl.Session = &netserver.Session{ID: int64(id)}
	pl.Engine = engine
	pl.envir = env
	pl.MapName = env.Name
	pl.CurrX, pl.CurrY = x, y
	pl.WAbil.HP, pl.WAbil.MaxHP = 50, 50
	env.AddObject(x, y, OS_MOVINGOBJECT, pl)
	engine.AddPlayer(pl)
	return pl
}

// TestMoveFamilyRegistration 注册表字面量/权限/usage
// （g_GameCommand，M2Share.pas:2766,2768,2770-2773,2775-2777,2841）。
func TestMoveFamilyRegistration(t *testing.T) {
	buildGMCommandTable(nil)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"humanlocal", "HumanLocal", 3},
		{"positionmove", "PositionMove", 3},
		{"moblevel", "MobLevel", 3},
		{"mobcount", "MobCount", 3},
		{"humancount", "HumanCount", 3},
		{"ting", "Ting", 10},
		{"superting", "SuperTing", 10},
		{"mapmove", "MapMove", 10},
		{"kill", "Kill", 10},
	}
	for _, c := range cases {
		def := gmCmdTable[c.literal]
		if def == nil {
			t.Errorf("%s: not registered", c.literal)
			continue
		}
		if def.key != c.key || def.resolvedPerm != c.perm {
			t.Errorf("%s: key=%s perm=%d, want key=%s perm=%d",
				c.literal, def.key, def.resolvedPerm, c.key, c.perm)
		}
	}
}

// TestMoveFamilyDispatch 权限拒绝与 "?" 帮助均被消费（不静默）。
func TestMoveFamilyDispatch(t *testing.T) {
	buildGMCommandTable(nil)
	p, server := newTestPlayer()
	for _, cmd := range []string{"ting abc", "kill abc", "positionmove 0 10 10"} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: permission-denied should be consumed", cmd)
		}
	}
	p.Permission = 10
	for _, cmd := range []string{
		"humanlocal ?", "positionmove ?", "mobcount ?", "humancount ?",
		"ting ?", "superting ?", "mapmove ?", "kill ?",
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: usage help should be consumed", cmd)
		}
	}
}

// TestCmdPositionMove 定点传送（CmdPositionMove，ObjBase.pas:13480-13523）。
func TestCmdPositionMove(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, _, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	p.Permission = 3 // 按 nPermissionMin 起步，验证禁传检查

	if !p.cmdPositionMove(server, []string{"0", "10", "10"}) {
		t.Fatal("positionmove should be consumed")
	}
	if p.CurrX != 10 || p.CurrY != 10 {
		t.Errorf("after teleport = (%d,%d), want (10,10)", p.CurrX, p.CurrY)
	}

	// 越界坐标：拒绝且不移动（'无法移动到地图...' 分支）。
	p.cmdPositionMove(server, []string{"0", "99", "99"})
	if p.CurrX != 10 || p.CurrY != 10 {
		t.Errorf("invalid target moved player to (%d,%d)", p.CurrX, p.CurrY)
	}

	// 地图禁传标志：权限<6 拒绝，权限≥6 豁免（nPermissionMax=6）。
	env.Flag.NoPositionMove = true
	p.cmdPositionMove(server, []string{"0", "3", "3"})
	if p.CurrX != 10 || p.CurrY != 10 {
		t.Error("NoPositionMove map should block perm-3 GM")
	}
	p.Permission = 6
	p.cmdPositionMove(server, []string{"0", "3", "3"})
	if p.CurrX != 3 || p.CurrY != 3 {
		t.Errorf("perm>=6 should bypass NoPositionMove, got (%d,%d)", p.CurrX, p.CurrY)
	}
	env.Flag.NoPositionMove = false
}

// TestCountMapObjects MobCount/HumanCount 统计口径
// （GetMapMonster/GetMapHuman，UsrEngn.pas:2879/2920）。
func TestCountMapObjects(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)

	mon := placeMonster(env, 900, 8, 8)
	engine.Monsters = append(engine.Monsters, mon)
	dead := placeMonster(env, 901, 9, 9)
	dead.Death = true
	engine.Monsters = append(engine.Monsters, dead)

	newMoveFamilyVictim(t, "victim1", 2, env, engine, 7, 7)

	if got := countMapMonsters(engine, env); got != 1 {
		t.Errorf("countMapMonsters = %d, want 1 (alive only)", got)
	}
	// fixture GM + victim1。
	if got := countMapHumans(engine, env); got != 2 {
		t.Errorf("countMapHumans = %d, want 2", got)
	}

	server := &netserver.TCPServer{}
	if !p.cmdMobCount(server, []string{"0"}) || !p.cmdHumanCount(server, []string{"0"}) {
		t.Error("count commands should be consumed")
	}
	// 不存在的地图同样消费（红字提示分支）。
	if !p.cmdMobCount(server, []string{"nomap"}) || !p.cmdHumanCount(server, []string{"nomap"}) {
		t.Error("map-not-found should be consumed")
	}
}

// TestCmdMobLevel 身边 2 格对象信息（CmdMobLevel，ObjBase.pas:13751-13776）。
func TestCmdMobLevel(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}

	placeMonster(env, 900, 6, 5) // 半径 2 内
	far := placeMonster(env, 901, 15, 15)
	engine.Monsters = append(engine.Monsters, far)

	if !p.cmdMobLevel(server, nil) {
		t.Error("moblevel should be consumed")
	}
}

// TestCmdTingSuperTingMapMove 回城/遣送/整图迁移。
func TestCmdTingSuperTingMapMove(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, mm := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}

	home := newTestEnv(30, 30)
	home.Name = "3"
	mm.maps["3"] = home
	engine.Config.Game.HomeMap = "3"

	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	// @Ting：目标回城。
	if !p.cmdTing(server, []string{"victim1"}) {
		t.Fatal("ting should be consumed")
	}
	if victim.envir != home || victim.MapName != "3" {
		t.Errorf("@Ting victim on map %s, want 3", victim.MapName)
	}
	// 不在线文案分支。
	if !p.cmdTing(server, []string{"nobody"}) {
		t.Error("ting offline should be consumed")
	}

	// @SuperTing：把 GM 与 victim 拉回源图，范围内另一玩家被遣送。
	teleportPlayerTo(p, server, env, 5, 5)
	teleportPlayerTo(victim, server, env, 6, 5)
	near := newMoveFamilyVictim(t, "near", 3, env, engine, 7, 5)
	farPl := newMoveFamilyVictim(t, "far", 4, env, engine, 18, 18)
	if !p.cmdSuperTing(server, []string{"victim1", "2"}) {
		t.Fatal("superting should be consumed")
	}
	if near.envir != home {
		t.Error("@SuperTing should send nearby player home")
	}
	if farPl.envir != env {
		t.Error("@SuperTing should not touch out-of-range player")
	}
	if p.envir != env {
		t.Error("@SuperTing should not move the GM")
	}

	// @MapMove：整图迁移（除 GM 外）。
	dst := newTestEnv(30, 30)
	dst.Name = "D001"
	mm.maps["D001"] = dst
	teleportPlayerTo(victim, server, env, 6, 5)
	if !p.cmdMapMove(server, []string{"0", "D001"}) {
		t.Fatal("mapmove should be consumed")
	}
	if victim.envir != dst {
		t.Error("@MapMove should move source-map players to destination")
	}
	if p.envir != env {
		t.Error("@MapMove should not move the GM")
	}
	// 源/目标图不存在分支。
	if !p.cmdMapMove(server, []string{"nomap", "D001"}) || !p.cmdMapMove(server, []string{"0", "nomap"}) {
		t.Error("mapmove map-not-found should be consumed")
	}
}

// TestCmdKill 无参杀面前对象、有参杀指定玩家（CmdKill，ObjBase.pas:13243-13270）。
func TestCmdKill(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}

	// 无参：面前一格（Dir=0 上 → (5,4)）的怪物被正常死亡结算杀死。
	mon := placeMonster(env, 900, 5, 4)
	mon.WAbil.HP, mon.WAbil.MaxHP = 100, 100
	engine.Monsters = append(engine.Monsters, mon)
	p.Dir = 0
	if !p.cmdKill(server, nil) {
		t.Fatal("kill should be consumed")
	}
	if !mon.Death || mon.WAbil.HP != 0 {
		t.Error("@Kill should kill faced monster via normal death path")
	}

	// 无参且面前无对象：用法提示分支。
	p.Dir = 4 // 下 → (5,6) 空格
	if !p.cmdKill(server, nil) {
		t.Error("kill without faced target should be consumed")
	}

	// 有参：同图玩家走 applyDamage 死亡路径。
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 8, 8)
	if !p.cmdKill(server, []string{"victim1"}) {
		t.Fatal("kill by name should be consumed")
	}
	if !victim.Death || victim.WAbil.HP != 0 {
		t.Error("@Kill by name should kill same-map player")
	}

	// 有参：跨图玩家走目标地图结算分支。
	home := newTestEnv(30, 30)
	home.Name = "3"
	remote := newMoveFamilyVictim(t, "remote", 3, home, engine, 4, 4)
	if !p.cmdKill(server, []string{"remote"}) {
		t.Fatal("kill remote should be consumed")
	}
	if !remote.Death || remote.WAbil.HP != 0 {
		t.Error("@Kill by name should kill cross-map player")
	}

	// 离线目标：不在线文案分支。
	if !p.cmdKill(server, []string{"nobody"}) {
		t.Error("kill offline should be consumed")
	}
}

// TestCmdKillSupermanExempt 无敌模式豁免（Delphi Die 检查 m_boSuperMan，
// ObjBase.pas:20782）。
func TestCmdKillSupermanExempt(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}

	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 8, 8)
	victim.SuperManMode = true
	p.cmdKill(server, []string{"victim1"})
	if victim.Death {
		t.Error("@Kill should not kill superman-mode player")
	}
}

// TestGmLethalDamage 致死量需压过魔法盾部分吸收上限（MaxMP×2/3）。
func TestGmLethalDamage(t *testing.T) {
	base := NewBaseObject("t", 1)
	base.WAbil.MaxHP = 65535
	base.WAbil.MaxMP = 65535
	dmg := gmLethalDamage(base)
	absorbed := int(base.WAbil.MaxMP) * 2 / 3
	if dmg-absorbed <= int(base.WAbil.MaxHP) {
		t.Errorf("lethal damage %d minus max shield absorb %d not fatal vs HP %d",
			dmg, absorbed, base.WAbil.MaxHP)
	}
}
