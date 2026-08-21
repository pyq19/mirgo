package main

import (
	"path/filepath"
	"testing"
)

// TestReloadCommandRegistration 重载族注册校验（字面量/权限对齐
// g_GameCommand 表 M2Share.pas:2802-2804,2859-2869；内部 <6 静默退出
// 的命令注册 6+silent，分发即实现的三条注册 5+fallback）。
func TestReloadCommandRegistration(t *testing.T) {
	ensureGMTable(t)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"reloadguild", "ReloadGuild", 10},
		{"reloadlinenotice", "ReloadLineNotice", 10},
		{"reloadabuse", "ReloadAbuse", 10},
		{"reloadadmin", "ReloadAdmin", 6},
		{"reloadmanage", "ReloadManage", 10},
		{"reloadmonitems", "ReloadMonItems", 6},
		{"reloaditemdb", "ReloadItemDB", 5},
		{"reloadmonsterdb", "ReloadMonsterDB", 5},
		{"reloadminmap", "ReLoadMinMap", 5},
	}
	for _, c := range cases {
		def := gmCmdTable[c.literal]
		if def == nil {
			t.Errorf("command @%s not registered", c.literal)
			continue
		}
		if def.key != c.key {
			t.Errorf("@%s key = %s, want %s", c.literal, def.key, c.key)
		}
		if def.resolvedPerm != c.perm {
			t.Errorf("@%s perm = %d, want %d", c.literal, def.resolvedPerm, c.perm)
		}
	}
	// ReloadRobot/ReloadRobotManage 未实装（Go 无机器人子系统）。
	for _, literal := range []string{"reloadrobot", "reloadrobotmanage"} {
		if gmCmdTable[literal] != nil {
			t.Errorf("@%s should not be registered (no robot subsystem)", literal)
		}
	}
}

// gmReloadFixture 重载族夹具：临时配置目录 + 空引擎。
func gmReloadFixture(t *testing.T) (*PlayObject, *UserEngine) {
	t.Helper()
	p, _ := newTestPlayer()
	e := &UserEngine{Config: &ServerConfig{}, configDir: t.TempDir()}
	p.Engine = e
	return p, e
}

// TestCmdReloadLineNotice 公告重载回退默认文案
// （loadNoticeText 文件缺失回退，main.go:38-54）。
func TestCmdReloadLineNotice(t *testing.T) {
	ensureGMTable(t)
	p, e := gmReloadFixture(t)
	_, server := newTestPlayer()
	if !p.cmdReloadLineNotice(server, nil) {
		t.Fatal("@ReloadLineNotice should consume")
	}
	if e.GetNoticeText() != defaultNoticeText {
		t.Errorf("notice = %q, want default fallback", e.GetNoticeText())
	}
}

// TestCmdReloadAbuse 脏话过滤表重载（文件缺失不 panic）。
func TestCmdReloadAbuse(t *testing.T) {
	ensureGMTable(t)
	p, _ := gmReloadFixture(t)
	_, server := newTestPlayer()
	if !p.cmdReloadAbuse(server, nil) {
		t.Fatal("@ReloadAbuse should consume")
	}
}

// TestCmdReloadDBFailure 物品/怪物库文件缺失走失败分支不 panic。
func TestCmdReloadDBFailure(t *testing.T) {
	ensureGMTable(t)
	p, e := gmReloadFixture(t)
	_, server := newTestPlayer()
	e.itemDBPath = filepath.Join(e.configDir, "missing_items.jsonc")
	e.monsterDBPath = filepath.Join(e.configDir, "missing_monsters.jsonc")
	if !p.cmdReloadItemDB(server, nil) {
		t.Fatal("@ReloadItemDB should consume")
	}
	if e.ItemDB != nil {
		t.Error("failed reload must not replace ItemDB")
	}
	if !p.cmdReloadMonsterDB(server, nil) {
		t.Fatal("@ReloadMonsterDB should consume")
	}
	if e.MonsterDB != nil {
		t.Error("failed reload must not replace MonsterDB")
	}
}

// TestCmdReloadAdmin 管理员表重载（空路径返回 0 条不 panic）。
func TestCmdReloadAdmin(t *testing.T) {
	ensureGMTable(t)
	p, _ := gmReloadFixture(t)
	_, server := newTestPlayer()
	if !p.cmdReloadAdmin(server, nil) {
		t.Fatal("@ReloadAdmin should consume")
	}
}

// TestCmdReloadGuild @ReloadGuild 行会不存在分支（大小写敏感匹配，
// Delphi FindGuild '=' 怪癖）与 DB 缺失时保留内存行会。
func TestCmdReloadGuild(t *testing.T) {
	ensureGMTable(t)
	p, e := gmReloadFixture(t)
	_, server := newTestPlayer()
	e.Guilds = []*Guild{{Name: "test", Leader: "boss", Members: []string{"boss"}}}

	if !p.cmdReloadGuild(server, []string{"ghost"}) {
		t.Fatal("@ReloadGuild should consume")
	}
	// db=nil：内存行会保留，仍报成功（Delphi FindGuild 命中后恒成功）。
	if !p.cmdReloadGuild(server, []string{"test"}) {
		t.Fatal("@ReloadGuild should consume")
	}
	if len(e.Guilds) != 1 {
		t.Errorf("guilds = %d, want 1", len(e.Guilds))
	}
}

// TestCmdReloadMonItems 爆物品表重载（空目录返回空表不 panic）。
func TestCmdReloadMonItems(t *testing.T) {
	ensureGMTable(t)
	p, e := gmReloadFixture(t)
	_, server := newTestPlayer()
	e.dropTablesDir = filepath.Join(e.configDir, "mon_items")
	if !p.cmdReloadMonItems(server, nil) {
		t.Fatal("@ReloadMonItems should consume")
	}
	if e.DropTables == nil {
		t.Error("@ReloadMonItems should install a (possibly empty) drop table")
	}
}

// TestCmdReloadManage NPC 脚本作废按类别分流（Go 无 QManage/QFunction
// 单例的退化语义）。
func TestCmdReloadManage(t *testing.T) {
	ensureGMTable(t)
	p, e := gmReloadFixture(t)
	_, server := newTestPlayer()
	e.Npcs = []*NpcObject{NewNpcObject("admin-npc", 1, 0), NewNpcObject("shop", 2, 0)}
	e.Npcs[1].IsMerchant = true
	if !p.cmdReloadManage(server, nil) {
		t.Fatal("@ReloadManage should consume")
	}
	if !p.cmdReloadManage(server, []string{"x"}) {
		t.Fatal("@ReloadManage with param should consume")
	}
}

// TestCmdReloadMinMap 小地图配置重载（文件缺失仅日志不 panic）。
func TestCmdReloadMinMap(t *testing.T) {
	ensureGMTable(t)
	p, e := gmReloadFixture(t)
	_, server := newTestPlayer()
	p.MapMgr = &MapManager{maps: map[string]*Environment{"0": newTestEnv(8, 8)}}
	e.mapMgr = p.MapMgr
	if !p.cmdReloadMinMap(server, nil) {
		t.Fatal("@ReLoadMinMap should consume")
	}
}
