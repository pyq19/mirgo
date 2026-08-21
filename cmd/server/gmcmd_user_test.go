package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/protocol"
)

// TestUserCommandRegistration 用户命令族注册校验
//（字面量与权限对齐 g_GameCommand 表，M2Share.pas:2725-2925）。
func TestUserCommandRegistration(t *testing.T) {
	ensureGMTable(t)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"date", "Date", 0},
		{"prvmsg", "PrvMsg", 0},
		{"allowmsg", "AllowMsg", 0},
		{"letshout", "LetShout", 0},
		{"lettrade", "LetTrade", 0},
		{"banguildchat", "BanGuildChat", 0},
		{"letguild", "LetGuild", 0},
		{"endguild", "EndGuild", 0},
		{"authally", "AuthAlly", 0},
		{"allowgrouprecall", "AllowGroupCall", 0},
		{"grouprecall", "GroupCall", 0},
		{"allowguildrecall", "AllowGuildReCall", 0},
		{"联盟", "Auth", 0},
		{"取消联盟", "AuthCancel", 0},
		{"memberfunc", "MemberFunc", 0},
		{"memberfuncex", "MemberFuncEx", 0},
		{"locklogin", "LockLogin", 0},
		{"clearbag", "ClearBag", 10},
		{"showuseiteminfo", "ShowUseItemInfo", 10},
		{"binduseitem", "BindUseItem", 10},
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

// TestCmdUserToggles 用户开关命令切换字段与文案分支
//（ObjBase.pas:7624-7660,7687-7693 内联段）。
func TestCmdUserToggles(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()

	p.HearWhisper = true
	p.cmdAllowMsg(server, nil)
	if p.HearWhisper {
		t.Error("@AllowMsg should toggle HearWhisper off")
	}
	p.cmdAllowMsg(server, nil)
	if !p.HearWhisper {
		t.Error("@AllowMsg should toggle HearWhisper back on")
	}

	p.BanShout = true
	p.cmdLetShout(server, nil)
	if p.BanShout {
		t.Error("@LetShout should toggle BanShout off")
	}

	p.AllowDeal = true
	p.cmdLetTrade(server, nil)
	if p.AllowDeal {
		t.Error("@LetTrade should toggle AllowDeal off")
	}

	p.BanGuildChat = true
	p.cmdBanGuildChat(server, nil)
	if p.BanGuildChat {
		t.Error("@BanGuildChat should toggle BanGuildChat off")
	}

	p.AllowGuildInvite = false
	p.cmdLetGuild(server, nil)
	if !p.AllowGuildInvite {
		t.Error("@LetGuild should toggle AllowGuildInvite on")
	}

	p.AllowGuildReCall = false
	p.cmdAllowGuildRecall(server, nil)
	if !p.AllowGuildReCall {
		t.Error("@AllowGuildRecall should toggle AllowGuildReCall on")
	}

	// @AllowGroupRecall：切换 + '?' 自定义帮助（ObjBase.pas:11355-11368）。
	p.AllowGroupReCall = false
	p.cmdAllowGroupRecall(server, nil)
	if !p.AllowGroupReCall {
		t.Error("@AllowGroupRecall should toggle AllowGroupReCall on")
	}
	p.cmdAllowGroupRecall(server, []string{"?"})
	if !p.AllowGroupReCall {
		t.Error("? help must not toggle AllowGroupReCall")
	}
}

// TestCmdPrvMsg 私聊拉黑切换（Delphi CmdPrvMsg，ObjBase.pas:13969-13996）。
func TestCmdPrvMsg(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()

	p.cmdPrvMsg(server, []string{"ZhangSan"})
	if len(p.BlockWhisper) != 1 || !p.isBlockWhisper("zhangsan") {
		t.Fatalf("BlockWhisper after add = %v", p.BlockWhisper)
	}
	// 大小写不敏感匹配（Delphi CompareText）：异体大小写同名即移除。
	p.cmdPrvMsg(server, []string{"zhangsan"})
	if len(p.BlockWhisper) != 0 {
		t.Errorf("BlockWhisper after remove = %v, want empty", p.BlockWhisper)
	}
	if p.isBlockWhisper("ZhangSan") {
		t.Error("removed name should not be blocked")
	}
	// 无参数只出帮助文案，不入列表。
	p.cmdPrvMsg(server, nil)
	if len(p.BlockWhisper) != 0 {
		t.Errorf("no-param call must not add entries, got %v", p.BlockWhisper)
	}
}

// TestWhisperBlock @AllowMsg/@PrvMsg 的消费端（Delphi Whisper 转发检查，
// ObjBase.pas:2181-2185）：目标关私聊或拉黑发送者时不转发。
func TestWhisperBlock(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	target, _ := newTestPlayer()
	target.Name = "收信人"
	target.ID = 2
	target.Session.ID = 2
	p.Engine = &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}

	// 目标禁止私聊：静默不转发（发送者收到"拒绝私聊"反馈）。
	target.HearWhisper = false
	p.HandleWhisper(SendMessage{Msg: "收信人/你好"}, server)

	// 目标允许私聊但拉黑了发送者。
	target.HearWhisper = true
	target.BlockWhisper = []string{p.Name}
	p.HandleWhisper(SendMessage{Msg: "收信人/你好"}, server)
}

// TestCmdEndGuild 退出行会（Delphi CmdEndGuild，ObjBase.pas:12841-12863）。
func TestCmdEndGuild(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	guild := &Guild{Name: "沙巴克", Leader: "掌门", Members: []string{"掌门", p.Name}, Ranks: map[string]string{"掌门": "掌门人", p.Name: "成员"}}
	p.Engine = &UserEngine{Config: &ServerConfig{}, Guilds: []*Guild{guild}}
	p.GuildName = guild.Name
	p.GuildRank = "成员"

	p.cmdEndGuild(server, nil)
	if guild.IsMember(p.Name) {
		t.Error("member should be removed from guild")
	}
	if p.GuildName != "" || p.GuildRank != "" {
		t.Error("player guild state should be cleared")
	}

	// 掌门人不能退出。
	leader, _ := newTestPlayer()
	leader.Name = "掌门"
	leader.Engine = p.Engine
	leader.cmdEndGuild(server, nil)
	if !guild.IsMember("掌门") {
		t.Error("leader must not leave guild via @EndGuild")
	}

	// 未加入行会。
	nobody, _ := newTestPlayer()
	nobody.Name = "路人"
	nobody.Engine = p.Engine
	nobody.cmdEndGuild(server, nil)
}

// TestGuildInviteBlocked @LetGuild 的消费端（Delphi ClientGuildAddMember，
// ObjBase.pas:18095-18112）：目标 AllowGuildInvite=false → 失败码 5。
func TestGuildInviteBlocked(t *testing.T) {
	ensureGMTable(t)
	leader, server := newTestPlayer()
	leader.Name = "掌门"
	target, _ := newTestPlayer()
	target.Name = "新人"
	target.ID = 2
	target.Session.ID = 2
	target.AllowGuildInvite = false
	guild := &Guild{Name: "测试行会", Leader: leader.Name, Members: []string{leader.Name}, Ranks: map[string]string{leader.Name: "掌门人"}}
	leader.Engine = &UserEngine{
		Config:         &ServerConfig{},
		Guilds:         []*Guild{guild},
		PlayObjectList: map[int32]*PlayObject{target.ID: target},
	}
	leader.GuildName = guild.Name
	leader.GuildRank = "掌门人"

	leader.HandleGuildAddMember(SendMessage{Msg: "新人"}, server)
	if guild.IsMember("新人") {
		t.Error("invite must be rejected when target disallows guild invites")
	}

	// 打开开关后允许加入。
	target.AllowGuildInvite = true
	leader.HandleGuildAddMember(SendMessage{Msg: "新人"}, server)
	if !guild.IsMember("新人") {
		t.Error("invite should succeed after target allows guild invites")
	}
}

// TestCryBanShout @LetShout 的消费端（Delphi CryCry 接收方检查，
// UsrEngn.pas:1516-1522）：BanShout=false 的玩家收不到喊话。
func TestCryBanShout(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	env := newTestEnv(20, 20)
	p.envir = env
	p.CurrX, p.CurrY = 10, 10
	env.AddObject(10, 10, OS_MOVINGOBJECT, p)

	other, _ := newTestPlayer()
	other.Name = "听众"
	other.ID = 2
	other.Session.ID = 2
	other.envir = env
	other.CurrX, other.CurrY = 11, 10
	env.AddObject(11, 10, OS_MOVINGOBJECT, other)

	p.Engine = &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{p.ID: p, other.ID: other}}

	other.BanShout = false
	p.HandleCry("喊话内容", server) // 不应 panic，接收方被跳过
	other.BanShout = true
	p.HandleCry("喊话内容", server)
}

// TestCmdClearBag 清空背包（Delphi CmdClearBagItem，ObjBase.pas:11048-11087）。
func TestCmdClearBag(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	target, _ := newTestPlayer()
	target.Name = "背包客"
	target.Session.ID = 2
	target.ItemList = []*protocol.UserItem{{MakeIndex: 1, WIndex: 10}}
	p.Engine = &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}

	p.cmdClearBag(server, []string{"背包客"})
	if len(target.ItemList) != 0 {
		t.Errorf("ItemList after clear = %d items, want 0", len(target.ItemList))
	}
	// 不在线目标。
	p.cmdClearBag(server, []string{"不存在"})
}

// TestCmdShowUseItemInfo 穿戴物品明细（Delphi CmdShowUseItemInfo，
// ObjBase.pas:14437-14474）；槽名映射 GetUseItemName。
func TestCmdShowUseItemInfo(t *testing.T) {
	if useItemSlotName(protocol.UWeapon) != "Weapon" || useItemSlotName(protocol.UDress) != "Dress" {
		t.Error("slot names must match Delphi U_*NAME literals")
	}
	if useItemSlotName(protocol.UBelt) != "" {
		t.Error("1.76 belt slot name should be empty (commented out in Delphi)")
	}
	if useItemSlotIdx("weapon") != protocol.UWeapon || useItemSlotIdx("Bujuk") != protocol.UBujuk {
		t.Error("slot index lookup must be case-insensitive")
	}
	if useItemSlotIdx("腰带") != -1 {
		t.Error("1.76 disabled slots must not resolve")
	}

	ensureGMTable(t)
	p, server := newTestPlayer()
	target, _ := newTestPlayer()
	target.Name = "模特"
	target.Session.ID = 2
	target.UseItems[protocol.UWeapon] = &protocol.UserItem{MakeIndex: 7, WIndex: 10, Dura: 500, DuraMax: 1000}
	p.Engine = &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}
	p.cmdShowUseItemInfo(server, []string{"模特"})
}

// TestCmdBindUseItem 装备绑定（Delphi CmdBindUseItem，ObjBase.pas:14476-14637；
// Go 内存登记表最小闭环）。
func TestCmdBindUseItem(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	target, _ := newTestPlayer()
	target.Name = "模特"
	target.Session.ID = 2
	target.AccountName = "acct1"
	target.UseItems[protocol.UWeapon] = &protocol.UserItem{MakeIndex: 42, WIndex: 10, Dura: 500, DuraMax: 1000}
	p.Engine = &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}
	defer func() {
		itemBindMu.Lock()
		delete(itemBindRegistry, 42)
		itemBindMu.Unlock()
	}()

	// 空槽位。
	p.cmdBindUseItem(server, []string{"模特", "Dress", "人物"})
	if len(itemBindRegistry) != 0 {
		t.Error("empty slot must not register a bind")
	}
	// 参数错误。
	p.cmdBindUseItem(server, []string{"模特", "Weapon"})
	p.cmdBindUseItem(server, []string{"模特", "Weapon", "邮箱"})
	if len(itemBindRegistry) != 0 {
		t.Error("bad params must not register a bind")
	}
	// 成功绑定（帐号）。
	p.cmdBindUseItem(server, []string{"模特", "Weapon", "帐号"})
	if got := itemBindRegistry[42]; got != "帐号:acct1" {
		t.Errorf("bind registry = %q, want 帐号:acct1", got)
	}
	// 重复绑定。
	p.cmdBindUseItem(server, []string{"模特", "Weapon", "人物"})
	if got := itemBindRegistry[42]; got != "帐号:acct1" {
		t.Errorf("rebind must not overwrite, got %q", got)
	}
}

// TestCmdAuthAlly 行会联盟授权（Go 最小闭环：掌门人查询联盟列表）。
func TestCmdAuthAlly(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	guild := &Guild{Name: "沙巴克", Leader: p.Name, Members: []string{p.Name}, AllyGuilds: []string{"兄弟行会"}}
	p.Engine = &UserEngine{Config: &ServerConfig{}, Guilds: []*Guild{guild}}
	p.cmdAuthAlly(server, nil) // 掌门人：输出联盟列表，不 panic

	member, _ := newTestPlayer()
	member.Name = "成员"
	member.Engine = p.Engine
	member.cmdAuthAlly(server, nil) // 非掌门人：静默退出
}

// TestCmdLockLogin Go 未实装登录锁，恒输出未启用文案。
func TestCmdLockLogin(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	if !p.cmdLockLogin(server, nil) {
		t.Error("@LockLogin should be consumed")
	}
}

// TestCmdGuildAlly @联盟：双掌门面对面结盟（Delphi ClientGuildAlly，
// ObjBase.pas:18217-18260）。
func TestCmdGuildAlly(t *testing.T) {
	ensureGMTable(t)
	env := newTestEnv(20, 20)

	p, server := newTestPlayer()
	p.Name = "掌门甲"
	p.envir = env
	p.CurrX, p.CurrY, p.Dir = 10, 10, 2 // 面向右
	env.AddObject(10, 10, OS_MOVINGOBJECT, p)
	guildA := &Guild{Name: "行会甲", Leader: p.Name, Members: []string{p.Name}}

	other, _ := newTestPlayer()
	other.Name = "掌门乙"
	other.ID = 2
	other.Session.ID = 2
	other.envir = env
	other.CurrX, other.CurrY, other.Dir = 17, 17, 6 // 先在远处
	env.AddObject(17, 17, OS_MOVINGOBJECT, other)
	guildB := &Guild{Name: "行会乙", Leader: other.Name, Members: []string{other.Name}}

	engine := &UserEngine{Config: &ServerConfig{}, Guilds: []*Guild{guildA, guildB}}
	p.Engine = engine
	other.Engine = engine

	// 无面对面目标：失败码 -1，不结盟。
	p.cmdGuildAlly(server, nil)
	if len(guildA.AllyGuilds) != 0 {
		t.Error("no facing target must not ally")
	}

	// 面对面但对方非掌门人：失败码 -3。
	env.RemoveObject(17, 17, OS_MOVINGOBJECT, other)
	other.CurrX, other.CurrY, other.Dir = 11, 10, 6 // 移到 p 面前，面向左
	env.AddObject(11, 10, OS_MOVINGOBJECT, other)
	guildB.Leader = "别人"
	p.cmdGuildAlly(server, nil)
	if len(guildA.AllyGuilds) != 0 {
		t.Error("non-master facing target must not ally")
	}
	guildB.Leader = other.Name

	// 行会战中：失败码 -2。
	guildA.WarGuilds = []GuildWar{{GuildName: guildB.Name, EndTick: 9e18}}
	p.cmdGuildAlly(server, nil)
	if len(guildA.AllyGuilds) != 0 {
		t.Error("warring guilds must not ally")
	}
	guildA.WarGuilds = nil

	// 成功结盟：双向记录。
	p.cmdGuildAlly(server, nil)
	if !guildA.IsAllyGuild("行会乙") || !guildB.IsAllyGuild("行会甲") {
		t.Errorf("ally not recorded both ways: %v / %v", guildA.AllyGuilds, guildB.AllyGuilds)
	}
	// 重复结盟被拒（Go 扩展防重复）。
	p.cmdGuildAlly(server, nil)
	if len(guildA.AllyGuilds) != 1 {
		t.Errorf("duplicate ally should be rejected, got %v", guildA.AllyGuilds)
	}
}

// TestCmdGuildBreakAlly @取消联盟：按行会名解除联盟（Delphi
// ClientGuildBreakAlly，ObjBase.pas:18270-18300）。
func TestCmdGuildBreakAlly(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Name = "掌门甲"
	guildA := &Guild{Name: "行会甲", Leader: p.Name, Members: []string{p.Name}, AllyGuilds: []string{"行会乙"}}
	guildB := &Guild{Name: "行会乙", Leader: "掌门乙", Members: []string{"掌门乙"}, AllyGuilds: []string{"行会甲"}}
	p.Engine = &UserEngine{Config: &ServerConfig{}, Guilds: []*Guild{guildA, guildB}}

	// 不存在的行会：失败回包，联盟不变。
	p.cmdGuildBreakAlly(server, []string{"不存在"})
	if len(guildA.AllyGuilds) != 1 {
		t.Error("unknown guild must not change alliance")
	}
	// 成功解除：双向移除。
	p.cmdGuildBreakAlly(server, []string{"行会乙"})
	if guildA.IsAllyGuild("行会乙") || guildB.IsAllyGuild("行会甲") {
		t.Errorf("alliance should be broken both ways: %v / %v", guildA.AllyGuilds, guildB.AllyGuilds)
	}
}

// TestCmdGroupRecall @GroupRecall 编组传送执行（Delphi CmdGroupRecall，
// ObjBase.pas:12942-12996）：召回许可成员、拒绝成员提示、冷却生效。
func TestCmdGroupRecall(t *testing.T) {
	ensureGMTable(t)
	env := newTestEnv(30, 30)
	p, server := newTestPlayer()
	p.Name = "队长"
	p.envir = env
	p.CurrX, p.CurrY = 10, 10
	p.HasRecallSuite = true
	env.AddObject(10, 10, OS_MOVINGOBJECT, p)

	okMember, _ := newTestPlayer()
	okMember.Name = "许可成员"
	okMember.ID = 2
	okMember.Session.ID = 2
	okMember.envir = env
	okMember.CurrX, okMember.CurrY = 25, 25
	okMember.AllowGroupReCall = true
	env.AddObject(25, 25, OS_MOVINGOBJECT, okMember)

	noMember, _ := newTestPlayer()
	noMember.Name = "拒绝成员"
	noMember.ID = 3
	noMember.Session.ID = 3
	noMember.envir = env
	noMember.CurrX, noMember.CurrY = 26, 26
	noMember.AllowGroupReCall = false
	env.AddObject(26, 26, OS_MOVINGOBJECT, noMember)

	p.Engine = &UserEngine{
		Config:         &ServerConfig{},
		PlayObjectList: map[int32]*PlayObject{p.ID: p, okMember.ID: okMember, noMember.ID: noMember},
		Parties:        map[int32]*Party{p.ID: {Leader: p.ID, Members: []int32{p.ID, okMember.ID, noMember.ID}}},
	}

	// 无套装无权限：拒绝。
	p.HasRecallSuite = false
	p.cmdGroupRecall(server, nil)
	if okMember.CurrX != 25 {
		t.Error("without recall suite group recall must not fire")
	}
	p.HasRecallSuite = true

	p.cmdGroupRecall(server, nil)
	if d := absInt(okMember.CurrX-p.CurrX) + absInt(okMember.CurrY-p.CurrY); d > 6 {
		t.Errorf("allowed member should be recalled near leader, dist=%d", d)
	}
	if noMember.CurrX != 26 {
		t.Error("rejecting member must not be recalled")
	}
	if p.GroupRecallTick == 0 {
		t.Error("cooldown should be armed after use")
	}
	// 冷却中：位置不再变化。
	before := p.GroupRecallTick
	okMember.CurrX, okMember.CurrY = 28, 28
	p.cmdGroupRecall(server, nil)
	if okMember.CurrX != 28 {
		t.Error("cooldown should block second recall")
	}
	if p.GroupRecallTick != before {
		t.Error("cooldown must not be extended while active")
	}
}

