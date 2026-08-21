package main

import "testing"

// TestQuestCommandRegistration 任务标志族注册校验
// （字面量/权限对齐 g_GameCommand 表 M2Share.pas:2788-2790,2823-2826,2899；
// Set* 三条为调试段内联分发，表值 10 不生效，实际门槛 >4，见
// gmcmd_quest.go init 注释）。
func TestQuestCommandRegistration(t *testing.T) {
	ensureGMTable(t)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"showflag", "ShowFlag", 10},
		{"showopen", "ShowOpen", 10},
		{"showunit", "ShowUnit", 10},
		{"clearmission", "ClearMission", 10},
		{"startquest", "StartQuest", 10},
		{"setflag", "SetFlag", 5},
		{"setopen", "SetOpen", 5},
		{"setunit", "SetUnit", 5},
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

// gmQuestFixture GM + 在线目标玩家夹具（目标 ID=2 防 objIndex 冲突）。
func gmQuestFixture(t *testing.T) (*PlayObject, *PlayObject, *UserEngine) {
	t.Helper()
	p, _ := newTestPlayer()
	target, _ := newTestPlayer()
	target.Name = "victim"
	target.ID = 2
	e := &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}
	p.Engine = e
	return p, target, e
}

// TestCmdSetShowFlag SetFlag/ShowFlag 读写闭环（Delphi 内联段
// ObjBase.pas:8290-8309 + CmdShowHumanFlag 11749-11780）。
func TestCmdSetShowFlag(t *testing.T) {
	ensureGMTable(t)
	p, target, _ := gmQuestFixture(t)
	_, server := newTestPlayer()

	if !p.cmdSetFlag(server, []string{"victim", "100", "1"}) {
		t.Fatal("@SetFlag should consume the line")
	}
	if !questGetBit(&target.QuestFlag, 100) {
		t.Error("@SetFlag 100 1 should set bit 100")
	}
	p.cmdShowFlag(server, []string{"victim", "100"}) // 输出 ON 文案
	p.cmdSetFlag(server, []string{"victim", "100", "0"})
	if questGetBit(&target.QuestFlag, 100) {
		t.Error("@SetFlag 100 0 should clear bit 100")
	}
}

// TestCmdSetOpenUnit SetOpen/SetUnit 分别写 QuestUnitOpen/QuestUnit。
func TestCmdSetOpenUnit(t *testing.T) {
	ensureGMTable(t)
	p, target, _ := gmQuestFixture(t)
	_, server := newTestPlayer()

	p.cmdSetOpen(server, []string{"victim", "5", "1"})
	p.cmdSetUnit(server, []string{"victim", "7", "1"})
	if !questGetBit(&target.QuestUnitOpen, 5) {
		t.Error("@SetOpen should set QuestUnitOpen bit 5")
	}
	if !questGetBit(&target.QuestUnit, 7) {
		t.Error("@SetUnit should set QuestUnit bit 7")
	}
	if questGetBit(&target.QuestFlag, 5) || questGetBit(&target.QuestFlag, 7) {
		t.Error("SetOpen/SetUnit must not touch QuestFlag")
	}
}

// TestCmdSetQuestBitMissingTarget 目标不存在时输出用法不 panic
// （Delphi GetPlayObject=nil 分支，ObjBase.pas:8302-8306）。
func TestCmdSetQuestBitMissingTarget(t *testing.T) {
	ensureGMTable(t)
	p, _, _ := gmQuestFixture(t)
	_, server := newTestPlayer()
	if !p.cmdSetFlag(server, []string{"nobody", "1", "1"}) {
		t.Error("@SetFlag with missing target should still consume")
	}
	if !p.cmdSetOpen(server, nil) {
		t.Error("@SetOpen with no params should still consume")
	}
	// 只有目标名（Delphi sParam2/sParam3 空串按 0）不越界。
	if !p.cmdSetFlag(server, []string{"victim"}) {
		t.Error("@SetFlag with only target should still consume")
	}
}

// TestCmdClearMissionFlagOnly 只清 QuestFlag，不清 QuestUnit/QuestUnitOpen
// （Delphi FillChar 仅 m_QuestFlag，ObjBase.pas:11193）。
func TestCmdClearMissionFlagOnly(t *testing.T) {
	ensureGMTable(t)
	p, target, _ := gmQuestFixture(t)
	_, server := newTestPlayer()

	questSetBit(&target.QuestFlag, 10)
	questSetBit(&target.QuestUnit, 20)
	questSetBit(&target.QuestUnitOpen, 30)
	p.cmdClearMission(server, []string{"victim"})
	if questGetBit(&target.QuestFlag, 10) {
		t.Error("@ClearMission should clear QuestFlag")
	}
	if !questGetBit(&target.QuestUnit, 20) || !questGetBit(&target.QuestUnitOpen, 30) {
		t.Error("@ClearMission must not touch QuestUnit/QuestUnitOpen")
	}
}

// TestCmdStartQuestNoParams 无参输出用法不 panic；有参无脚本命中静默。
func TestCmdStartQuestNoParams(t *testing.T) {
	ensureGMTable(t)
	p, _, _ := gmQuestFixture(t)
	_, server := newTestPlayer()
	if !p.cmdStartQuest(server, nil) {
		t.Error("@StartQuest no-params should consume")
	}
	if !p.cmdStartQuest(server, []string{"somequest"}) {
		t.Error("@StartQuest with label should consume (silent no-op)")
	}
}
