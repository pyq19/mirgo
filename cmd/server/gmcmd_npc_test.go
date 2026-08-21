package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// NPC/怪物/地图模式族测试（gmcmd_npc.go）。
// Delphi 参照行号见各命令注释。

// newNpcFamilyFixture 复用传送族 fixture：GM + 地图 "0" + 引擎。
func newNpcFamilyFixture(t *testing.T) (*PlayObject, *Environment, *UserEngine, *MapManager) {
	t.Helper()
	return newMoveFamilyFixture(t)
}

// TestNpcFamilyRegistration 注册表字面量/权限
//（g_GameCommand，M2Share.pas:2793-2795,2818-2819,2847-2848,2915-2916,2923）。
func TestNpcFamilyRegistration(t *testing.T) {
	buildGMCommandTable(nil)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"mobnpc", "MobNpc", 10},
		{"delnpc", "DelNpc", 10},
		{"npcscript", "NpcScript", 10},
		{"mission", "Mission", 10},
		{"mobplace", "MobPlace", 10},
		// 表权限 10 但实现内 <6 静默退出 → perm 6（ObjBase.pas:11867/11889/12868）。
		{"showmapmode", "ShowMapMode", 6},
		{"setmapmode", "SetMapMode", 6},
		{"fireburn", "FireBurn", 6},
		{"testfire", "TestFire", 6},
		{"mobfireburn", "MobFireBurn", 6},
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

// TestNpcFamilyDispatch 权限拒绝与 "?" 帮助均被消费（含 permFailSilent 命令）。
func TestNpcFamilyDispatch(t *testing.T) {
	buildGMCommandTable(nil)
	p, server := newTestPlayer()
	// 普通玩家（perm 0）：permFailLow 与 permFailSilent 都消费。
	for _, cmd := range []string{"mobnpc a b", "mission 1 2", "showmapmode 0", "fireburn 1 2 3"} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: permission-denied should be consumed", cmd)
		}
	}
	p.Permission = 10
	for _, cmd := range []string{
		"mobnpc ?", "npcscript ?", "mission ?", "mobplace ?",
		"showmapmode ?", "setmapmode ?", "fireburn ?", "testfire ?",
		"mobfireburn ?", // usage 为空：handler 内帮助分支
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: usage help should be consumed", cmd)
		}
	}
}

// TestCmdMobNpc 面前生成 NPC（CmdMobNpc，ObjBase.pas:13777-13815）。
func TestCmdMobNpc(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}
	p.Dir = 0 // 上 → (5,4)

	if !p.cmdMobNpc(server, []string{"测试商人", "testscript", "11", "1"}) {
		t.Fatal("mobnpc should be consumed")
	}
	if len(engine.Npcs) != 1 {
		t.Fatalf("engine.Npcs = %d, want 1", len(engine.Npcs))
	}
	npc := engine.Npcs[0]
	if npc.Name != "测试商人" || npc.Appr != 11 || !npc.Castle || !npc.IsMerchant {
		t.Errorf("npc fields: name=%s appr=%d castle=%v merchant=%v",
			npc.Name, npc.Appr, npc.Castle, npc.IsMerchant)
	}
	if npc.CurrX != 5 || npc.CurrY != 4 || npc.envir != env || npc.MapName != "0" {
		t.Errorf("npc placement = (%d,%d) env=%v map=%s, want (5,4) on map 0",
			npc.CurrX, npc.CurrY, npc.envir != nil, npc.MapName)
	}
	if got := env.GetMovingObject(5, 4); got != interface{}(npc) {
		t.Error("npc not inserted into map cell (5,4)")
	}
	// 缺参：帮助分支，不创建。
	if !p.cmdMobNpc(server, []string{"只有名字"}) || len(engine.Npcs) != 1 {
		t.Error("missing-param mobnpc should only show help")
	}
	// 脚本名按 npc_scripts/<名>-<地图>.txt 解析。
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "npc_scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	scriptPath := filepath.Join(dir, "npc_scripts", "quest-0.txt")
	if err := os.WriteFile(scriptPath, []byte("[@main]\n#IF\n#ACT\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	engine.npcConfigDir = dir
	p.Dir = 2 // 右 → (6,5)
	p.cmdMobNpc(server, []string{"任务NPC", "quest", "15", "0"})
	if len(engine.Npcs) != 2 {
		t.Fatalf("engine.Npcs = %d, want 2", len(engine.Npcs))
	}
	if engine.Npcs[1].Script != scriptPath {
		t.Errorf("script resolved to %q, want %q", engine.Npcs[1].Script, scriptPath)
	}
}

// TestCmdDelNpc 删除面对的 NPC（CmdDelNpc，ObjBase.pas:12368-12414）。
func TestCmdDelNpc(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}

	// 面前无 NPC：g_sGameCommandDelNpcMsg 分支（M2Share.pas:3306）。
	if !p.cmdDelNpc(server, nil) {
		t.Error("delnpc without facing npc should be consumed")
	}

	npc := NewNpcObject("待删NPC", 7001, 10)
	npc.MapName, npc.CurrX, npc.CurrY, npc.envir = env.Name, 5, 4, env
	env.AddObject(5, 4, OS_MOVINGOBJECT, npc)
	engine.Npcs = append(engine.Npcs, npc)

	p.Dir = 0 // 上 → (5,4)
	if !p.cmdDelNpc(server, nil) {
		t.Fatal("delnpc should be consumed")
	}
	if len(engine.Npcs) != 0 {
		t.Errorf("engine.Npcs = %d after delete, want 0", len(engine.Npcs))
	}
	if env.GetMovingObject(5, 4) == npc {
		t.Error("npc still in map cell after delete")
	}
	if !npc.Ghost {
		t.Error("deleted npc should be ghosted")
	}
	// '?' 帮助分支（空 usage，ObjBase.pas:12380-12382）。
	if !p.cmdDelNpc(server, []string{"?"}) {
		t.Error("delnpc ? should be consumed")
	}
}

// TestCmdNpcScript 打印面对 NPC 脚本（CmdNpcScript，ObjBase.pas:13857-13935）。
func TestCmdNpcScript(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}

	// 面前无 NPC：面对面提示分支。
	if !p.cmdNpcScript(server, []string{"x"}) {
		t.Error("npcscript without facing npc should be consumed")
	}

	scriptFile := filepath.Join(t.TempDir(), "npc.txt")
	content := "[@main]\n#IF\nCHECK level 7\n#ACT\nGOTO @exit\n"
	if err := os.WriteFile(scriptFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	npc := NewNpcObject("脚本NPC", 7002, 10)
	npc.MapName, npc.CurrX, npc.CurrY, npc.envir = env.Name, 5, 4, env
	npc.Script = scriptFile
	env.AddObject(5, 4, OS_MOVINGOBJECT, npc)
	engine.Npcs = append(engine.Npcs, npc)

	p.Dir = 0
	if !p.cmdNpcScript(server, []string{"x"}) {
		t.Fatal("npcscript should be consumed")
	}
	// 脚本路径不可读：'读取脚本文件错误: ' 分支。
	npc.Script = filepath.Join(t.TempDir(), "不存在.txt")
	if !p.cmdNpcScript(server, []string{"x"}) {
		t.Error("npcscript with unreadable script should be consumed")
	}
	// 缺参：帮助分支（g_sGameCommandNpcScriptHelpMsg='？？？？'，M2Share.pas:3305）。
	if !p.cmdNpcScript(server, nil) {
		t.Error("npcscript without params should be consumed")
	}
}

// TestNpcScriptLines 脚本行格式化：'%d,%s'、空格→逗号、超行截断。
func TestNpcScriptLines(t *testing.T) {
	got := npcScriptLines("[@main]\nCHECK level 7\n", 100)
	want := []string{"0,[@main]", "1,CHECK,level,7"}
	if len(got) != len(want) {
		t.Fatalf("lines = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
	// 截断保护：>100 行追加省略提示（Go 聊天窗口保护）。
	var sb strings.Builder
	for i := 0; i < 120; i++ {
		fmt.Fprintf(&sb, "line %d\n", i)
	}
	got = npcScriptLines(sb.String(), 100)
	if len(got) != 101 {
		t.Fatalf("truncated lines = %d, want 101", len(got))
	}
	if !strings.Contains(got[100], "120") || !strings.Contains(got[100], "100") {
		t.Errorf("truncation notice %q should mention total/limit", got[100])
	}
}

// TestCmdMission 设置怪物攻城集中点（CmdMission，ObjBase.pas:13593-13614）。
func TestCmdMission(t *testing.T) {
	buildGMCommandTable(nil)
	p, _, engine, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}

	if !p.cmdMission(server, []string{"12", "34"}) {
		t.Fatal("mission should be consumed")
	}
	engine.mu.RLock()
	m, x, y := engine.MissionMap, engine.MissionX, engine.MissionY
	engine.mu.RUnlock()
	if m != "0" || x != 12 || y != 34 {
		t.Errorf("mission = %s(%d:%d), want 0(12:34)", m, x, y)
	}
	// 覆盖设定。
	p.cmdMission(server, []string{"1", "2"})
	engine.mu.RLock()
	x, y = engine.MissionX, engine.MissionY
	engine.mu.RUnlock()
	if x != 1 || y != 2 {
		t.Errorf("mission overwrite = (%d:%d), want (1:2)", x, y)
	}
	// 缺参帮助分支。
	if !p.cmdMission(server, []string{"1"}) {
		t.Error("mission missing param should be consumed")
	}
}

// TestCmdMobPlace 集中点刷怪（CmdMobPlace，ObjBase.pas:13816-13856）。
func TestCmdMobPlace(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}
	engine.MonsterDB = &MonsterDB{} // GetByName 返回 nil → 走缺省怪定义

	// 未设集中点：双行提示分支。
	if !p.cmdMobPlace(server, []string{"10", "10", "鸡", "3"}) {
		t.Fatal("mobplace without mission should be consumed")
	}
	if len(engine.Monsters) != 0 {
		t.Fatalf("monsters = %d before mission set, want 0", len(engine.Monsters))
	}

	p.cmdMission(server, []string{"8", "8"})
	if !p.cmdMobPlace(server, []string{"10", "10", "鸡", "3"}) {
		t.Fatal("mobplace should be consumed")
	}
	if len(engine.Monsters) != 3 {
		t.Fatalf("monsters = %d after mobplace, want 3", len(engine.Monsters))
	}
	for _, mon := range engine.Monsters {
		if mon.envir != env {
			t.Errorf("monster on env %v, want mission map", mon.envir)
		}
		if abs(mon.CurrX-8) > 9 || abs(mon.CurrY-8) > 9 {
			t.Errorf("monster at (%d,%d) not scattered near mission point (8,8)",
				mon.CurrX, mon.CurrY)
		}
	}

	// 数量上限 500（Delphi _MIN(500, ...)，ObjBase.pas:13828）。
	before := len(engine.Monsters)
	p.cmdMobPlace(server, []string{"10", "10", "鸡", "999"})
	spawned := len(engine.Monsters) - before
	if spawned <= 0 || spawned > 500 {
		t.Errorf("mobplace 999 spawned %d, want (0,500]", spawned)
	}
	// 非法参数帮助分支。
	if !p.cmdMobPlace(server, []string{"0", "10", "鸡", "3"}) ||
		!p.cmdMobPlace(server, []string{"10", "10", "", "3"}) {
		t.Error("invalid-param mobplace should be consumed")
	}
}

// TestGetEnvirInfo 地图标志列表输出（Delphi GetEnvirInfo，Envir.pas:1499-1550）。
func TestGetEnvirInfo(t *testing.T) {
	env := newTestEnv(5, 5)
	env.Flag.MusicID = -1 // parseMapFlag 缺省（mapmanager.go:117）
	if got := getEnvirInfo(env); got != "" {
		t.Errorf("empty flags info = %q, want empty", got)
	}
	env.Flag.Safe = true
	env.Flag.Dark = true
	env.Flag.NoReconnect = true
	env.Flag.ReconnectMap = "3"
	env.Flag.MusicID = 5
	env.Flag.ExpRate = 2
	env.Flag.DecHPPoint, env.Flag.DecHPTime = 10, 60
	env.Flag.NoChat = true
	got := getEnvirInfo(env)
	for _, want := range []string{"SAFE", "DARK", "NORECONNECT(3)", "MUSIC(5)", "EXPRATE(2)", "DECHP(10/60)", "NOCHAT"} {
		if !strings.Contains(got, want) {
			t.Errorf("getEnvirInfo %q missing %q", got, want)
		}
	}
}

// TestCmdShowMapMode 显示地图模式（CmdShowMapMode，ObjBase.pas:11862-11881）。
func TestCmdShowMapMode(t *testing.T) {
	buildGMCommandTable(nil)
	p, _, _, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}

	if !p.cmdShowMapMode(server, []string{"0"}) {
		t.Error("showmapmode should be consumed")
	}
	if !p.cmdShowMapMode(server, []string{"nomap"}) {
		t.Error("showmapmode map-not-found should be consumed")
	}
	if !p.cmdShowMapMode(server, nil) {
		t.Error("showmapmode missing param should be consumed")
	}
}

// TestCmdSetMapMode 运行时地图标志修改（CmdSetMapMode，ObjBase.pas:11883-12253）：
// 带值参数置位、无参数清除（非翻转）。
func TestCmdSetMapMode(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, _, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}

	set := func(args ...string) {
		t.Helper()
		if !p.cmdSetMapMode(server, args) {
			t.Fatalf("setmapmode %v should be consumed", args)
		}
	}

	// 布尔开关族：带参置位 / 无参清除。
	boolModes := []struct {
		mode string
		get  func() bool
	}{
		{"SAFE", func() bool { return env.Flag.Safe }},
		{"DARK", func() bool { return env.Flag.Dark }},
		{"FIGHT", func() bool { return env.Flag.Fight }},
		{"FIGHT3", func() bool { return env.Flag.Fight3 }},
		{"DAY", func() bool { return env.Flag.DayLight }},
		{"QUIZ", func() bool { return env.Flag.Quiz }},
		{"NEEDHOLE", func() bool { return env.Flag.NeedHole }},
		{"NORECALL", func() bool { return env.Flag.NoRecall }},
		{"NORANDOMMOVE", func() bool { return env.Flag.NoRandomMove }},
		{"NODRUG", func() bool { return env.Flag.NoDrug }},
		{"MINE", func() bool { return env.Flag.Mine }},
		{"MINE2", func() bool { return env.Flag.Mine2 }},
		{"NOTHROWITEM", func() bool { return env.Flag.NoThrowItem }},
		{"NODROPITEM", func() bool { return env.Flag.NoDropItem }},
		{"NOPOSITIONMOVE", func() bool { return env.Flag.NoPositionMove }},
		{"NOHORSE", func() bool { return env.Flag.NoHorse }},
		{"NOCHAT", func() bool { return env.Flag.NoChat }},
	}
	for _, c := range boolModes {
		set("0", c.mode, "1")
		if !c.get() {
			t.Errorf("%s with param should be set", c.mode)
		}
		set("0", c.mode)
		if c.get() {
			t.Errorf("%s without param should be cleared", c.mode)
		}
	}

	// NORECONNECT：带目标地图名（ObjBase.pas:11961-11966）。
	set("0", "noreconnect", "3") // CompareText 大小写不敏感
	if !env.Flag.NoReconnect || env.Flag.ReconnectMap != "3" {
		t.Errorf("NORECONNECT(3) = %v/%q", env.Flag.NoReconnect, env.Flag.ReconnectMap)
	}
	set("0", "NORECONNECT")
	if env.Flag.NoReconnect {
		t.Error("NORECONNECT without param should be cleared")
	}

	// MUSIC/EXPRATE 数值族。
	set("0", "MUSIC", "12")
	if env.Flag.MusicID != 12 {
		t.Errorf("MUSIC(12) = %d", env.Flag.MusicID)
	}
	set("0", "MUSIC")
	if env.Flag.MusicID != -1 {
		t.Errorf("MUSIC cleared = %d, want -1", env.Flag.MusicID)
	}
	set("0", "EXPRATE", "2")
	if env.Flag.ExpRate != 2 {
		t.Errorf("EXPRATE(2) = %d", env.Flag.ExpRate)
	}
	set("0", "EXPRATE")
	if env.Flag.ExpRate != 0 {
		t.Errorf("EXPRATE cleared = %d, want 0", env.Flag.ExpRate)
	}

	// DECHP：sParam1=时间 sParam2=点数，缺一清除（ObjBase.pas:12038-12044）。
	set("0", "DECHP", "60", "5")
	if env.Flag.DecHPTime != 60 || env.Flag.DecHPPoint != 5 {
		t.Errorf("DECHP(60,5) = time %d point %d", env.Flag.DecHPTime, env.Flag.DecHPPoint)
	}
	set("0", "DECHP", "60")
	if env.Flag.DecHPTime != 0 || env.Flag.DecHPPoint != 0 {
		t.Error("DECHP with single param should be cleared")
	}

	// Go MapFlag 无对应字段的模式：'暂不支持: ' 分支，不 panic。
	for _, mode := range []string{"PKWINLEVEL", "PKWINEXP", "PKLOSTLEVEL", "PKLOSTEXP",
		"DECGAMEGOLD", "INCGAMEGOLD", "INCGAMEPOINT", "RUNHUMAN", "RUNMON"} {
		set("0", mode, "5")
	}
	// 未知模式：Delphi if 链落空静默。
	set("0", "不存在的模式", "1")

	// 地图不存在与缺参分支。
	set("nomap", "SAFE", "1")
	set("0")
}

// TestCmdFireBurn 身边放火（CmdFireBurn，ObjBase.pas:12864-12876）：
// 复用魔法火墙同款 AddFireEvent（mapevent.go）。
func TestCmdFireBurn(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, _, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}

	// 零参帮助分支。
	if !p.cmdFireBurn(server, nil) {
		t.Fatal("fireburn without params should be consumed")
	}
	if len(env.Events) != 0 {
		t.Fatalf("events = %d after help branch, want 0", len(env.Events))
	}

	if !p.cmdFireBurn(server, []string{"5", "10000", "10"}) {
		t.Fatal("fireburn should be consumed")
	}
	if len(env.Events) != 1 {
		t.Fatalf("events = %d, want 1", len(env.Events))
	}
	ev := env.Events[0]
	if ev.Type != protocol.ETFire || ev.X != p.CurrX || ev.Y != p.CurrY {
		t.Errorf("event type=%d at (%d,%d), want ETFire at (%d,%d)",
			ev.Type, ev.X, ev.Y, p.CurrX, p.CurrY)
	}
	if ev.Damage != 10 || ev.Duration != 10000 || ev.OwnerID != p.ID {
		t.Errorf("event damage=%d duration=%d owner=%d, want 10/10000/%d",
			ev.Damage, ev.Duration, ev.OwnerID, p.ID)
	}
}

// TestCmdTestFire 方形环火墙（CmdTestFire，ObjBase.pas:15080-15104）。
func TestCmdTestFire(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, _, _ := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}

	if !p.cmdTestFire(server, []string{"2", "5", "10", "8"}) {
		t.Fatal("testfire should be consumed")
	}
	// 边长 2*2+1 的方形环 = 16 格。
	if len(env.Events) != 16 {
		t.Fatalf("events = %d, want 16 (5x5 perimeter)", len(env.Events))
	}
	for _, ev := range env.Events {
		dx, dy := ev.X-p.CurrX, ev.Y-p.CurrY
		if abs(dx) > 2 || abs(dy) > 2 {
			t.Errorf("event at (%d,%d) outside range 2", ev.X, ev.Y)
		}
		if abs(dx) != 2 && abs(dy) != 2 {
			t.Errorf("event at (%d,%d) not on perimeter", ev.X, ev.Y)
		}
		if ev.Damage != 8 || ev.Duration != 10000 {
			t.Errorf("event damage=%d duration=%d, want 8/10000", ev.Damage, ev.Duration)
		}
	}

	// range=0：脚下 1 格。
	env.Events = nil
	p.cmdTestFire(server, []string{"0", "5", "1", "1"})
	if len(env.Events) != 1 || env.Events[0].X != p.CurrX || env.Events[0].Y != p.CurrY {
		t.Errorf("range-0 events = %d, want 1 at feet", len(env.Events))
	}
}

// TestCmdMobFireBurn 指定地图放火（CmdMobFireBurn，ObjBase.pas:13707-13750）。
func TestCmdMobFireBurn(t *testing.T) {
	buildGMCommandTable(nil)
	p, _, _, mm := newNpcFamilyFixture(t)
	server := &netserver.TCPServer{}
	other := newTestEnv(30, 30)
	other.Name = "3"
	mm.maps["3"] = other

	// 正常：目标地图新增火事件（nTime 秒×1000）。
	if !p.cmdMobFireBurn(server, []string{"3", "10", "12", "5", "20", "7"}) {
		t.Fatal("mobfireburn should be consumed")
	}
	if len(other.Events) != 1 {
		t.Fatalf("target map events = %d, want 1", len(other.Events))
	}
	ev := other.Events[0]
	if ev.X != 10 || ev.Y != 12 || ev.Damage != 7 || ev.Duration != 20000 || ev.Type != protocol.ETFire {
		t.Errorf("event = (%d,%d) dmg %d dur %d type %d, want (10,12) 7 20000 ETFire",
			ev.X, ev.Y, ev.Damage, ev.Duration, ev.Type)
	}

	// nPoint 缺省钳制为 1（ObjBase.pas:13731）。
	p.cmdMobFireBurn(server, []string{"3", "5", "5", "5", "5"})
	if n := len(other.Events); n != 2 || other.Events[1].Damage != 1 {
		t.Errorf("missing-point event damage should clamp to 1 (events=%d)", n)
	}

	// 地图不存在：'地图%s 不存在'（M2Share.pas:3355）。
	if !p.cmdMobFireBurn(server, []string{"nomap", "1", "1", "5", "5", "1"}) {
		t.Error("mobfireburn map-not-found should be consumed")
	}
	// 参数缺失/非法：帮助分支。
	for _, args := range [][]string{nil, {"3"}, {"3", "x", "1", "5", "5", "1"}} {
		if !p.cmdMobFireBurn(server, args) {
			t.Errorf("mobfireburn %v should be consumed by help branch", args)
		}
	}
}
