package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// 封禁/发言过滤族测试（gmcmd_ban.go + usrengine.go 持久化 +
// main.go 登录拦截判定 + chatsystem.go 消费端）。
// Delphi 参照行号见各测试注释。

// TestBanFamilyRegistration 注册表字面量/键/权限对照
// （g_GameCommand，M2Share.pas:2838-2840,2903-2912；
// DENYIPLOGON 的 Command.ini 键为 DenyIPaddrLogon，M2Share.pas:6724）。
// Shutup/ReleaseShutup/ShutupList 实现于 gmcommands.go（早期批次），
// 一并纳入本族注册表对照。
func TestBanFamilyRegistration(t *testing.T) {
	ensureGMTable(t)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"shutup", "Shutup", 10},
		{"releaseshutup", "ReleaseShutup", 10},
		{"shutuplist", "ShutupList", 10},
		{"denyiplogon", "DenyIPaddrLogon", 10},
		{"denyaccountlogon", "DenyAccountLogon", 10},
		{"denycharnamelogon", "DenyCharNameLogon", 10},
		{"deldenyiplogon", "DelDenyIPLogon", 10},
		{"deldenyaccountlogon", "DelDenyAccountLogon", 10},
		{"deldenycharnamelogon", "DelDenyCharNameLogon", 10},
		{"showdenyiplogon", "ShowDenyIPLogon", 10},
		{"showdenyaccountlogon", "ShowDenyAccountLogon", 6},
		{"showdenycharnamelogon", "ShowDenyCharNameLogon", 10},
		{"disablesendmsg", "DisableSendMsg", 10},
		{"enablesendmsg", "EnableSendMsg", 10},
		{"disablesendmsglist", "DisableSendMsgList", 10},
		{"viewwhisper", "ViewWhisper", 10},
	}
	for _, c := range cases {
		def := gmCmdTable[c.literal]
		if def == nil {
			t.Errorf("@%s: not registered", c.literal)
			continue
		}
		if def.key != c.key {
			t.Errorf("@%s key = %s, want %s", c.literal, def.key, c.key)
		}
		if def.resolvedPerm != c.perm {
			t.Errorf("@%s perm = %d, want %d", c.literal, def.resolvedPerm, c.perm)
		}
	}
	// @ShowDenyAccountLogon 权限不足静默（Delphi 硬编码 m_btPermission<6 Exit，
	// ObjBase.pas:12692）；其余封禁命令走常规"权限不够"。
	if def := gmCmdTable["showdenyaccountlogon"]; def != nil && def.permFail != permFailSilent {
		t.Error("@ShowDenyAccountLogon must fail silently (permFailSilent)")
	}
	for _, lit := range []string{"denyiplogon", "showdenyiplogon", "disablesendmsg"} {
		if def := gmCmdTable[lit]; def != nil && def.permFail != permFailLow {
			t.Errorf("@%s must use permFailLow", lit)
		}
	}
}

// TestBanFamilyDispatch 权限拒绝与 "?" 帮助均被消费。
func TestBanFamilyDispatch(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	for _, cmd := range []string{
		"denyiplogon 1.2.3.4", "deldenyaccountlogon abc", "showdenycharnamelogon",
		"disablesendmsg abc", "viewwhisper abc",
		// <6 静默分支（ObjBase.pas:12692）同样被消费。
		"showdenyaccountlogon",
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: permission-denied should be consumed", cmd)
		}
	}
	// 权限 6-9：@ShowDenyAccountLogon 可用（Delphi 硬编码门槛 6），
	// 其余封禁命令仍被权限拦截。
	p.Permission = 6
	if !p.HandleGMCommand("showdenyaccountlogon", server) {
		t.Error("perm-6 GM may use @ShowDenyAccountLogon")
	}
	if !p.cmdShowDenyAccountLogon(server, nil) {
		t.Error("@ShowDenyAccountLogon handler must be consumed")
	}
	p.Permission = 10
	for _, cmd := range []string{
		"denyiplogon ?", "denyaccountlogon ?", "denycharnamelogon ?",
		"deldenyiplogon ?", "deldenyaccountlogon ?", "deldenycharnamelogon ?",
		"disablesendmsg ?", "enablesendmsg ?", "viewwhisper ?",
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: usage help should be consumed", cmd)
		}
	}
}

// TestCmdDenyLogonAddDel 封禁增删（CmdDenyIPaddrLogon/CmdDelDenyIPaddrLogon，
// ObjBase.pas:12524-12551/12649-12685）：第二参首字符 '1' 为永久并落盘，
// 删除永久条目回写，删除不存在条目不改列表。
func TestCmdDenyLogonAddDel(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Permission = 10
	engine := &UserEngine{Config: &ServerConfig{}}
	engine.denyListPath = filepath.Join(t.TempDir(), "deny_lists.json")
	p.Engine = engine

	// 临时封禁：入表不落盘。
	p.cmdDenyIPLogon(server, []string{"10.0.0.1"})
	if len(engine.DenyIPList) != 1 || engine.DenyIPList[0].Value != "10.0.0.1" || engine.DenyIPList[0].Permanent {
		t.Fatalf("temp deny entry wrong: %+v", engine.DenyIPList)
	}
	if _, err := os.Stat(engine.denyListPath); !os.IsNotExist(err) {
		t.Error("temp deny must not write deny_lists.json")
	}

	// 永久封禁：入表且落盘。
	p.cmdDenyIPLogon(server, []string{"10.0.0.2", "1"})
	if len(engine.DenyIPList) != 2 || !engine.DenyIPList[1].Permanent {
		t.Fatalf("permanent deny entry wrong: %+v", engine.DenyIPList)
	}
	data, err := os.ReadFile(engine.denyListPath)
	if err != nil {
		t.Fatalf("permanent deny must write deny_lists.json: %v", err)
	}
	if !strings.Contains(string(data), "10.0.0.2") || strings.Contains(string(data), "10.0.0.1") {
		t.Errorf("deny_lists.json must contain only permanent entries, got %s", data)
	}

	// 帐号/角色名表同结构。
	p.cmdDenyAccountLogon(server, []string{"acct1", "1"})
	p.cmdDenyCharNameLogon(server, []string{"角色甲"})
	if len(engine.DenyAccountList) != 1 || !engine.DenyAccountList[0].Permanent {
		t.Errorf("account deny wrong: %+v", engine.DenyAccountList)
	}
	if len(engine.DenyCharList) != 1 || engine.DenyCharList[0].Permanent {
		t.Errorf("char deny wrong: %+v", engine.DenyCharList)
	}

	// 删除（CompareText 不敏感）：永久条目删除后回写。
	p.cmdDelDenyIPLogon(server, []string{"10.0.0.2"})
	if len(engine.DenyIPList) != 1 || engine.DenyIPList[0].Value != "10.0.0.1" {
		t.Fatalf("after del: %+v", engine.DenyIPList)
	}
	data, _ = os.ReadFile(engine.denyListPath)
	if strings.Contains(string(data), "10.0.0.2") {
		t.Error("removed permanent entry must be rewritten out of deny_lists.json")
	}
	p.cmdDelDenyIPLogon(server, []string{"10.0.0.9"}) // 不存在：文案分支，列表不变
	if len(engine.DenyIPList) != 1 {
		t.Error("deleting unknown entry must not change list")
	}
	// 大小写不敏感删除。
	p.cmdDelDenyAccountLogon(server, []string{"ACCT1"})
	if len(engine.DenyAccountList) != 0 {
		t.Error("EqualFold delete failed")
	}
}

// TestCmdShowDenyLogon 列表打印分支（ObjBase.pas:12687-12759）。
func TestCmdShowDenyLogon(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Permission = 10
	engine := &UserEngine{Config: &ServerConfig{}}
	p.Engine = engine

	// 空表分支（含 @ShowDenyIPLogon 复用角色列表文案的怪癖，12749）。
	if !p.cmdShowDenyIPLogon(server, nil) || !p.cmdShowDenyAccountLogon(server, nil) || !p.cmdShowDenyCharNameLogon(server, nil) {
		t.Error("show commands must be consumed")
	}
	// 非空分支。
	engine.DenyIPList = []DenyEntry{{Value: "1.1.1.1"}}
	engine.DenyAccountList = []DenyEntry{{Value: "acc"}}
	engine.DenyCharList = []DenyEntry{{Value: "角色"}}
	if !p.cmdShowDenyIPLogon(server, nil) || !p.cmdShowDenyAccountLogon(server, nil) || !p.cmdShowDenyCharNameLogon(server, nil) {
		t.Error("show commands with entries must be consumed")
	}
	if msgDenyIPEmpty != msgDenyCharEmpty {
		t.Error("Delphi quirk: ShowDenyIPLogon empty text reuses char-list text")
	}
}

// TestCheckDenyLogon 登录拦截判定纯函数（Delphi CheckDenyLogon，
// ObjBase.pas:26763-26783；CompareText 不敏感，顺序 IP→帐号→角色名）。
func TestCheckDenyLogon(t *testing.T) {
	engine := &UserEngine{
		DenyIPList:      []DenyEntry{{Value: "192.168.1.10"}},
		DenyAccountList: []DenyEntry{{Value: "BadAccount"}},
		DenyCharList:    []DenyEntry{{Value: "坏角色"}},
	}
	cases := []struct {
		ip, account, char string
		want              string
	}{
		{"192.168.1.10", "ok", "ok", "你当前登录的IP地址已被禁止登录了！！！"},
		{"1.2.3.4", "badaccount", "ok", "你当前登录的帐号已被禁止登录了！！！"},
		{"1.2.3.4", "ok", "坏角色", "你当前登录的人物已被禁止登录了！！！"},
		{"1.2.3.4", "ok", "ok", ""},
		{"", "", "", ""},
		// IP 命中优先于帐号/角色（Delphi else-if 链顺序）。
		{"192.168.1.10", "BadAccount", "坏角色", "你当前登录的IP地址已被禁止登录了！！！"},
	}
	for _, c := range cases {
		if got := engine.CheckDenyLogon(c.ip, c.account, c.char); got != c.want {
			t.Errorf("CheckDenyLogon(%q,%q,%q) = %q, want %q", c.ip, c.account, c.char, got, c.want)
		}
	}
	// 大小写不敏感（CompareText）。
	if engine.CheckDenyLogon("192.168.1.10", "", "") == "" {
		t.Error("IP match must be case-insensitive")
	}
	if engine.CheckDenyLogon("", "BADACCOUNT", "") == "" {
		t.Error("account match must be case-insensitive")
	}
}

// TestDenyListsLoadSave LoadDenyLists/SaveDenyLists 往返（Go 替代 Delphi
// Envir/Deny*List.txt，M2Share.pas:4755-4972）：只持久化永久条目。
func TestDenyListsLoadSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "deny_lists.json")

	e1 := &UserEngine{}
	e1.DenyIPList = []DenyEntry{{Value: "9.9.9.9", Permanent: true}, {Value: "8.8.8.8"}}
	e1.DenyAccountList = []DenyEntry{{Value: "permAcc", Permanent: true}}
	e1.DenyCharList = []DenyEntry{{Value: "tempChar"}}
	e1.denyListPath = path
	e1.SaveDenyLists()

	var doc denyListFile
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("SaveDenyLists: %v", err)
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("deny_lists.json invalid: %v", err)
	}
	if len(doc.IP) != 1 || doc.IP[0] != "9.9.9.9" {
		t.Errorf("ip = %v, want [9.9.9.9] (permanent only)", doc.IP)
	}
	if len(doc.Account) != 1 || doc.Account[0] != "permAcc" {
		t.Errorf("account = %v, want [permAcc]", doc.Account)
	}
	if len(doc.Char) != 0 {
		t.Errorf("char = %v, want empty", doc.Char)
	}

	// 加载：恢复为永久条目。
	e2 := &UserEngine{}
	e2.LoadDenyLists(path)
	if len(e2.DenyIPList) != 1 || e2.DenyIPList[0].Value != "9.9.9.9" || !e2.DenyIPList[0].Permanent {
		t.Errorf("loaded ip = %+v", e2.DenyIPList)
	}
	if len(e2.DenyAccountList) != 1 || !e2.DenyAccountList[0].Permanent {
		t.Errorf("loaded account = %+v", e2.DenyAccountList)
	}
	if len(e2.DenyCharList) != 0 {
		t.Errorf("loaded char = %+v", e2.DenyCharList)
	}

	// 文件不存在：空表、无错误；此后 SaveDenyLists 使用记录的路径。
	e3 := &UserEngine{}
	e3.LoadDenyLists(filepath.Join(t.TempDir(), "absent.json"))
	if len(e3.DenyIPList)+len(e3.DenyAccountList)+len(e3.DenyCharList) != 0 {
		t.Error("absent file must load empty lists")
	}

	// 未加载过（无路径）：保存为无操作。
	e4 := &UserEngine{}
	e4.DenyIPList = []DenyEntry{{Value: "x", Permanent: true}}
	e4.SaveDenyLists() // 不应 panic
}

// TestMuteCommands 禁言三命令（CmdDisableSendMsg/CmdEnableSendMsg/
// CmdDisableSendMsgList，ObjBase.pas:12762-12840）。
func TestMuteCommands(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Permission = 10
	target, _ := newTestPlayer()
	target.Name = "禁言对象"
	target.ID = 2
	engine := &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}
	p.Engine = engine

	// 在线目标：设 FilterSendMsg + 入名单。
	p.cmdDisableSendMsg(server, []string{"禁言对象"})
	if !target.FilterSendMsg {
		t.Error("@DisableSendMsg must set FilterSendMsg on online target")
	}
	if !engine.isMsgDisabled("禁言对象") {
		t.Error("muted name must be in DisableSendMsg list")
	}

	// 离线目标：仍入名单（登录时恢复，ObjBase.pas:16628）。
	p.cmdDisableSendMsg(server, []string{"离线玩家"})
	if !engine.isMsgDisabled("离线玩家") {
		t.Error("offline name must also be recorded")
	}

	// 列表命令消费（非空/空分支）。
	if !p.cmdDisableSendMsgList(server, nil) {
		t.Error("@DisableSendMsgList must be consumed")
	}

	// 解除：清 FilterSendMsg + 出名单（CompareText 不敏感）。
	p.cmdEnableSendMsg(server, []string{"禁言对象"})
	if target.FilterSendMsg {
		t.Error("@EnableSendMsg must clear FilterSendMsg")
	}
	if engine.isMsgDisabled("禁言对象") {
		t.Error("unmuted name must leave DisableSendMsg list")
	}
	// 未禁言分支（c_Red 文案）。
	if !p.cmdEnableSendMsg(server, []string{"路人"}) {
		t.Error("enable on non-muted name must be consumed")
	}
	// 空名单分支。
	p.cmdEnableSendMsg(server, []string{"离线玩家"})
	if !p.cmdDisableSendMsgList(server, nil) {
		t.Error("@DisableSendMsgList empty branch must be consumed")
	}
}

// TestShutupFamily @Shutup/@ReleaseShutup/@ShutupList 行为
// （Delphi CmdShutup/CmdShutupList/CmdReleaseShutup，ObjBase.pas:
// 14800-14831/14832-14862/14863-14898；实现见 gmcommands.go，
// HandleSay 的 ShutupTick 消费见 chatsystem.go:38-41）。
func TestShutupFamily(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Permission = 10
	target, _ := newTestPlayer()
	target.Name = "话痨"
	target.ID = 2
	engine := &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{p.ID: p, target.ID: target}}
	p.Engine = engine
	target.Engine = engine

	// 默认禁言 5 分钟。
	now := time.Now().UnixMilli()
	p.cmdShutup(server, []string{"话痨"})
	if target.ShutupTick < now+4*60000 || target.ShutupTick > now+6*60000 {
		t.Errorf("ShutupTick = %d, want ~now+5min (%d)", target.ShutupTick, now+5*60000)
	}
	// 禁言期间发言被 HandleSay 拦截（早退，不进分发）。
	target.HandleSay(SendMessage{Msg: "还能说吗"}, server)
	target.HandleSay(SendMessage{Msg: "@shutuplist"}, server) // GM 命令同样被拦
	// 自定义分钟数。
	p.cmdShutup(server, []string{"话痨", "3"})
	if target.ShutupTick < now+2*60000 || target.ShutupTick > now+4*60000 {
		t.Errorf("custom minutes ShutupTick = %d, want ~now+3min", target.ShutupTick)
	}
	// 列表命令消费（非空）。
	if !p.cmdShutupList(server, nil) {
		t.Error("@ShutupList must be consumed")
	}
	// 解除禁言。
	p.cmdReleaseShutup(server, []string{"话痨"})
	if target.ShutupTick != 0 {
		t.Error("@ReleaseShutup must clear ShutupTick")
	}
	// 空列表分支。
	if !p.cmdShutupList(server, nil) {
		t.Error("@ShutupList empty branch must be consumed")
	}
	// 不在线目标分支。
	if !p.cmdShutup(server, []string{"不存在"}) || !p.cmdReleaseShutup(server, []string{"不存在"}) {
		t.Error("offline target branches must be consumed")
	}
}

// TestViewWhisper 侦听私聊切换（CmdViewWhisper，ObjBase.pas:15237-15268）。
func TestViewWhisper(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Permission = 10
	target, _ := newTestPlayer()
	target.Name = "密语者"
	target.ID = 2
	p.Engine = &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}

	// 开始侦听（Msg2）。
	p.cmdViewWhisper(server, []string{"密语者"})
	if p.ViewWhisperName != "密语者" {
		t.Fatalf("ViewWhisperName = %q, want 密语者", p.ViewWhisperName)
	}
	// 再输同一名取消（Msg1）。
	p.cmdViewWhisper(server, []string{"密语者"})
	if p.ViewWhisperName != "" {
		t.Error("same name must toggle listening off")
	}
	// 空参取消（Go 模型：无侦听目标时出参数帮助）。
	p.cmdViewWhisper(server, []string{"密语者"})
	p.cmdViewWhisper(server, nil)
	if p.ViewWhisperName != "" {
		t.Error("empty param must cancel listening")
	}
	if !p.cmdViewWhisper(server, nil) {
		t.Error("empty param without listening must be consumed (usage help)")
	}
	// 目标不在线（g_sNowNotOnLineOrOnOtherServer）。
	if !p.cmdViewWhisper(server, []string{"不存在"}) {
		t.Error("offline target must be consumed")
	}
	if p.ViewWhisperName != "" {
		t.Error("offline target must not arm listening")
	}
}

// banWireFixture 真实 TCP 会话夹具：起一个 TCPServer，拨入 n 个
// mockClient，返回对应 PlayObject（同图相邻，可互相收到附近聊天）。
func banWireFixture(t *testing.T, n int) (*netserver.TCPServer, *UserEngine, []*PlayObject, []*mockClient) {
	t.Helper()
	// 找空闲端口（与 newTestServer 同款做法，main_test.go:191-196）。
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	srv := netserver.NewTCPServer(addr)
	if err := srv.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(srv.Stop)

	env := newTestEnv(20, 20)
	env.Name = "0"
	mm := NewMapManager("")
	mm.maps["0"] = env
	engine := NewUserEngine(nil, mm)
	engine.Config = &ServerConfig{}

	players := make([]*PlayObject, 0, n)
	clients := make([]*mockClient, 0, n)
	for i := 1; i <= n; i++ {
		c := newMockClient(t, addr)
		clients = append(clients, c)
		deadline := time.Now().Add(2 * time.Second)
		for srv.GetSessionCount() < i && time.Now().Before(deadline) {
			time.Sleep(5 * time.Millisecond)
		}
		sess := srv.GetSession(int64(i))
		if sess == nil {
			t.Fatalf("session %d not registered", i)
		}
		name := fmt.Sprintf("玩家%d", i)
		pl := &PlayObject{BaseObject: NewBaseObject(name, int32(i))}
		pl.Session = sess
		pl.Engine = engine
		pl.envir = env
		pl.MapName = "0"
		pl.CurrX, pl.CurrY = 5+i, 5
		pl.WAbil.HP, pl.WAbil.MaxHP = 100, 100
		// NewPlayObject 默认值（playobject.go:288）：夹具未经构造函数创建。
		pl.HearWhisper = true
		env.AddObject(pl.CurrX, pl.CurrY, OS_MOVINGOBJECT, pl)
		engine.AddPlayer(pl)
		players = append(players, pl)
	}
	return srv, engine, players, clients
}

// TestFilterSendMsgSayWire 禁言消费端（Delphi m_boFilterSendMsg，
// ObjBase.pas:8730-8735）：被禁言者发言只回显给自己，不广播附近。
func TestFilterSendMsgSayWire(t *testing.T) {
	ensureGMTable(t)
	srv, engine, players, clients := banWireFixture(t, 2)
	sender := players[0]
	senderClient, otherClient := clients[0], clients[1]

	sender.FilterSendMsg = true
	sender.HandleSay(SendMessage{Msg: "只有自己能看见"}, srv)

	msg, body, err := senderClient.recvTimeout(2 * time.Second)
	if err != nil {
		t.Fatalf("muted sender self-echo: %v", err)
	}
	if msg.Ident != protocol.SMHear || !strings.Contains(body, "只有自己能看见") {
		t.Errorf("self-echo = ident %d body %q, want SMHear with text", msg.Ident, body)
	}
	if got, body2, err := otherClient.recvTimeout(300 * time.Millisecond); err == nil {
		t.Errorf("muted say must not broadcast, other got ident %d body %q", got.Ident, body2)
	}

	// 解除后正常广播。
	sender.FilterSendMsg = false
	sender.HandleSay(SendMessage{Msg: "大家都能看见"}, srv)
	msg, body, err = otherClient.recvTimeout(2 * time.Second)
	if err != nil {
		t.Fatalf("normal say broadcast: %v", err)
	}
	if msg.Ident != protocol.SMHear || !strings.Contains(body, "大家都能看见") {
		t.Errorf("broadcast = ident %d body %q, want SMHear with text", msg.Ident, body)
	}

	// GM 命令不受禁言过滤（消费端在 '@' 分支之后）：
	// 用 @denyiplogon 的可观察副作用（入 DenyIPList）证明仍被分发。
	sender.Permission = 10
	sender.FilterSendMsg = true
	sender.HandleSay(SendMessage{Msg: "@denyiplogon 5.5.5.5 1"}, srv)
	if !denyEntryHit(engine.DenyIPList, "5.5.5.5") {
		t.Error("@ command must still be dispatched while muted")
	}
}

// TestWhisperMuteAndWatchWire 私聊消费端：禁言者不能发起私聊
// （ObjBase.pas:8674-8676）；侦听 GM 收到转发副本
// （ObjBase.pas:2190-2202，Go 扩展格式 '[侦听] A→B: 内容'）。
func TestWhisperMuteAndWatchWire(t *testing.T) {
	ensureGMTable(t)
	srv, _, players, clients := banWireFixture(t, 3)
	sender, target, gm := players[0], players[1], players[2]
	targetClient, gmClient := clients[1], clients[2]

	// 禁言者私聊被静默丢弃。
	sender.FilterSendMsg = true
	sender.HandleWhisper(SendMessage{Msg: target.Name + "/悄悄话"}, srv)
	if got, body, err := targetClient.recvTimeout(300 * time.Millisecond); err == nil {
		t.Errorf("muted whisper must be dropped, target got ident %d body %q", got.Ident, body)
	}
	sender.FilterSendMsg = false

	// 正常私聊 + GM 侦听目标侧。
	gm.ViewWhisperName = target.Name
	sender.HandleWhisper(SendMessage{Msg: target.Name + "/秘密内容"}, srv)

	msg, body, err := targetClient.recvTimeout(2 * time.Second)
	if err != nil {
		t.Fatalf("whisper to target: %v", err)
	}
	if msg.Ident != protocol.SMWhisper || !strings.Contains(body, "秘密内容") {
		t.Errorf("target recv = ident %d body %q", msg.Ident, body)
	}
	msg, body, err = gmClient.recvTimeout(2 * time.Second)
	if err != nil {
		t.Fatalf("watcher copy: %v", err)
	}
	if msg.Ident != protocol.SMSysMessage || !strings.Contains(body, "[侦听]") || !strings.Contains(body, "秘密内容") {
		t.Errorf("watcher recv = ident %d body %q, want [侦听] copy", msg.Ident, body)
	}

	// 侦听发起者侧同样收到副本（Delphi 两侧 m_GetWhisperHuman 检查）。
	gm.ViewWhisperName = sender.Name
	sender.HandleWhisper(SendMessage{Msg: target.Name + "/第二句"}, srv)
	msg, body, err = gmClient.recvTimeout(2 * time.Second)
	if err != nil {
		t.Fatalf("watcher copy (sender side): %v", err)
	}
	if !strings.Contains(body, "第二句") {
		t.Errorf("watcher sender-side copy body %q", body)
	}

	// 未侦听的私聊不打扰 GM。
	gm.ViewWhisperName = ""
	sender.HandleWhisper(SendMessage{Msg: target.Name + "/第三句"}, srv)
	if got, body, err := gmClient.recvTimeout(300 * time.Millisecond); err == nil {
		t.Errorf("no watcher must not forward, gm got ident %d body %q", got.Ident, body)
	}
}
