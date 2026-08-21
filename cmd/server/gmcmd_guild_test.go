package main

import (
	"path/filepath"
	"testing"

	"github.com/pyq0109/mirgo/internal/storage"
)

// openTestGuildDB 行会持久化测试库（@AddGuild/@DelGuild 走 SaveGuilds/
// DeleteGuild 真实落盘路径）。
func openTestGuildDB(t *testing.T) *storage.Database {
	t.Helper()
	db, err := storage.Open(filepath.Join(t.TempDir(), "guild_test.db"))
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// TestGuildCommandRegistration 行会/城堡族注册校验（字面量/权限对齐
// g_GameCommand 表 M2Share.pas:2785,2831-2834,2873-2876）。
func TestGuildCommandRegistration(t *testing.T) {
	ensureGMTable(t)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"sabukwallgold", "SabukWallGold", 10},
		{"contestpoint", "ContestPoint", 10},
		{"startcontest", "StartContest", 10},
		{"endcontest", "EndContest", 10},
		{"announcement", "Announcement", 10},
		{"addguild", "AddGuild", 10},
		{"delguild", "DelGuild", 10},
		{"changesabuklord", "ChangeSabukLord", 10},
		{"forcedwallconquestwar", "ForcedWallconquestWar", 10},
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
}

// TestCmdAddGuild 为在线玩家建行会（Delphi ObjBase.pas:10483-10534）。
func TestCmdAddGuild(t *testing.T) {
	ensureGMTable(t)
	p, _ := newTestPlayer()
	chief, _ := newTestPlayer()
	chief.Name = "chief"
	chief.ID = 2
	e := &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{chief.ID: chief}}
	e.db = openTestGuildDB(t)
	p.Engine = e
	_, server := newTestPlayer()

	p.cmdAddGuild(server, []string{"testguild", "chief"})
	if len(e.Guilds) != 1 || e.Guilds[0].Name != "testguild" || e.Guilds[0].Leader != "chief" {
		t.Fatalf("guild not created: %+v", e.Guilds)
	}
	if chief.GuildName != "testguild" || chief.GuildRank != "掌门人" {
		t.Errorf("chief guild state = %q/%q", chief.GuildName, chief.GuildRank)
	}
	// 掌门人已入行会 → 静默不建（Delphi 原版行为）。
	p.cmdAddGuild(server, []string{"other", "chief"})
	// 行会重名 → 静默不建。
	other, _ := newTestPlayer()
	other.Name = "other"
	other.ID = 3
	e.PlayObjectList[other.ID] = other
	p.cmdAddGuild(server, []string{"testguild", "other"})
	if len(e.Guilds) != 1 {
		t.Errorf("guilds = %d, want 1 (silent no-op branches)", len(e.Guilds))
	}
	// 目标离线 → 不建。
	p.cmdAddGuild(server, []string{"ghostguild", "nobody"})
	if len(e.Guilds) != 1 {
		t.Errorf("guilds = %d, want 1", len(e.Guilds))
	}
}

// TestCmdDelGuild 删除行会并清理在线成员（Delphi ObjBase.pas:12343-12368）。
func TestCmdDelGuild(t *testing.T) {
	ensureGMTable(t)
	p, _ := newTestPlayer()
	member, _ := newTestPlayer()
	member.Name = "member"
	member.ID = 2
	member.GuildName = "doomed"
	member.GuildRank = "成员"
	e := &UserEngine{
		Config:         &ServerConfig{},
		PlayObjectList: map[int32]*PlayObject{member.ID: member},
		Guilds:         []*Guild{{Name: "doomed", Leader: "boss", Members: []string{"boss", "member"}}},
	}
	e.db = openTestGuildDB(t)
	p.Engine = e
	_, server := newTestPlayer()

	// 未找到（CompareText 不敏感查找，'nomatch'）。
	p.cmdDelGuild(server, []string{"nomatch"})
	if len(e.Guilds) != 1 {
		t.Fatal("nonexistent guild must not delete anything")
	}
	// 大小写不敏感命中并删除。
	p.cmdDelGuild(server, []string{"DOOMED"})
	if len(e.Guilds) != 0 {
		t.Fatalf("guilds = %d, want 0", len(e.Guilds))
	}
	if member.GuildName != "" || member.GuildRank != "" {
		t.Errorf("member guild state = %q/%q, want cleared", member.GuildName, member.GuildRank)
	}
}

// TestCmdChangeSabukLord 更改城堡所属行会（Delphi ObjBase.pas:10940-10984）。
func TestCmdChangeSabukLord(t *testing.T) {
	ensureGMTable(t)
	castle, _, e := newCastleSetup(t)
	e.Guilds = []*Guild{{Name: "owners", Leader: "boss", Members: []string{"boss"}}}
	e.MonsterDB = &MonsterDB{}
	p, _ := newTestPlayer()
	p.Engine = e
	_, server := newTestPlayer()

	p.cmdChangeSabukLord(server, []string{"沙巴克", "owners"})
	if castle.OwnerGuild != "owners" {
		t.Errorf("OwnerGuild = %q, want owners", castle.OwnerGuild)
	}
	// 城堡名不匹配 / 行会不存在分支。
	p.cmdChangeSabukLord(server, []string{"不存在", "owners"})
	p.cmdChangeSabukLord(server, []string{"沙巴克", "ghost"})
	if castle.OwnerGuild != "owners" {
		t.Error("failed lookups must not change the owner")
	}
}

// TestCmdForcedWallconquestWar 强制开/停攻城战切换
// （Delphi ObjBase.pas:12877-12916）。
func TestCmdForcedWallconquestWar(t *testing.T) {
	ensureGMTable(t)
	castle, _, e := newCastleSetup(t)
	e.MonsterDB = &MonsterDB{}
	p, _ := newTestPlayer()
	p.Engine = e
	_, server := newTestPlayer()

	p.cmdForcedWallconquestWar(server, []string{"沙巴克"})
	if !castle.IsAtWar() {
		t.Fatal("first call should force-start the war")
	}
	p.cmdForcedWallconquestWar(server, []string{"沙巴克"})
	if castle.IsAtWar() {
		t.Fatal("second call should end the war")
	}
}

// TestCmdSabukWallGold 金库查询与 =/+/- 调整（'-' 修正为减 nGold，
// Delphi 原版 Dec 不带步长 bug 不继承）。
func TestCmdSabukWallGold(t *testing.T) {
	ensureGMTable(t)
	castle, _, e := newCastleSetup(t)
	p, _ := newTestPlayer()
	p.Engine = e
	_, server := newTestPlayer()

	// 查询分支（无参）。
	if !p.cmdSabukWallGold(server, nil) {
		t.Fatal("@SabukWallGold query should consume")
	}
	// '=' 设定。
	p.cmdSabukWallGold(server, []string{"沙巴克", "=", "5000"})
	if castle.GetGold() != 5000 {
		t.Errorf("gold = %d, want 5000", castle.GetGold())
	}
	// '+' 增加。
	p.cmdSabukWallGold(server, []string{"沙巴克", "+", "100"})
	if castle.GetGold() != 5100 {
		t.Errorf("gold = %d, want 5100", castle.GetGold())
	}
	// '-' 减全额（原版 bug 修正）。
	p.cmdSabukWallGold(server, []string{"沙巴克", "-", "200"})
	if castle.GetGold() != 4900 {
		t.Errorf("gold = %d, want 4900", castle.GetGold())
	}
	// 结果钳制到 0。
	p.cmdSabukWallGold(server, []string{"沙巴克", "-", "99999999"})
	if castle.GetGold() != 0 {
		t.Errorf("gold = %d, want clamped 0", castle.GetGold())
	}
	// 非法控制符/越界值走用法分支不改值。
	p.cmdSabukWallGold(server, []string{"沙巴克", "*", "1"})
	p.cmdSabukWallGold(server, []string{"沙巴克", "=", "100000001"})
	if castle.GetGold() != 0 {
		t.Error("invalid params must not change gold")
	}
}

// TestCmdContestPoint 争霸赛得分查询双分支
// （行会不存在也用 c_Green 为原版怪癖）。
func TestCmdContestPoint(t *testing.T) {
	ensureGMTable(t)
	p, _ := newTestPlayer()
	e := &UserEngine{Config: &ServerConfig{}, Guilds: []*Guild{
		{Name: "fighters", Leader: "boss", Members: []string{"boss"}, ContestPoint: 42},
	}}
	p.Engine = e
	_, server := newTestPlayer()

	if !p.cmdContestPoint(server, []string{"fighters"}) {
		t.Fatal("@ContestPoint found branch should consume")
	}
	if !p.cmdContestPoint(server, []string{"ghost"}) {
		t.Fatal("@ContestPoint not-found branch should consume")
	}
	if !p.cmdContestPoint(server, nil) {
		t.Fatal("@ContestPoint usage branch should consume")
	}
}

// TestCmdContestLifecycle FIGHT3 门槛 + 争霸赛开始/结束生命周期。
func TestCmdContestLifecycle(t *testing.T) {
	ensureGMTable(t)
	env := newTestEnv(20, 20)
	member, _ := newTestPlayer()
	member.Name = "fighter"
	member.ID = 2
	member.envir = env
	guild := &Guild{Name: "fighters", Leader: "fighter", Members: []string{"fighter"}}
	e := &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{member.ID: member}, Guilds: []*Guild{guild}}

	p, _ := newTestPlayer()
	p.Engine = e
	p.envir = env
	_, server := newTestPlayer()

	// 非 FIGHT3 地图拒绝。
	p.cmdStartContest(server, nil)
	if guild.TeamFightActive {
		t.Fatal("StartContest must reject non-FIGHT3 maps")
	}
	env.Flag.Fight3 = true

	p.cmdStartContest(server, nil)
	if !guild.TeamFightActive {
		t.Fatal("@StartContest should activate the team fight")
	}
	if _, ok := guild.TeamFightMembers["fighter"]; !ok {
		t.Error("@StartContest should register map members")
	}
	p.cmdAnnouncement(server, []string{"fighters"})
	p.cmdEndContest(server, nil)
	if guild.TeamFightActive {
		t.Fatal("@EndContest should deactivate the team fight")
	}
	// 行会不存在的公告分支不 panic（原版此处 AV）。
	p.cmdAnnouncement(server, []string{"ghost"})
}
