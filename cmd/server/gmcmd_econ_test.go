package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// 属性/经济/技能族测试（gmcmd_econ.go）。
// Delphi 参照行号见各命令注释。

// newEconMagicDB 测试技能库（名字对齐 MagIC.db 常见技能，仅作夹具）。
func newEconMagicDB() *MagicDB {
	db := &MagicDB{Magics: []MagicDef{
		{MagID: 1, MagName: "火球术"},
		{MagID: 6, MagName: "雷电术"},
		{MagID: 15, MagName: "抗拒火环"},
	}}
	db.byID = make(map[int]*MagicDef, len(db.Magics))
	db.byName = make(map[string]*MagicDef, len(db.Magics))
	for i := range db.Magics {
		m := &db.Magics[i]
		db.byID[m.MagID] = m
		db.byName[m.MagName] = m
	}
	return db
}

// TestEconFamilyRegistration 注册表字面量/权限
//（g_GameCommand，M2Share.pas:2800-2862；调试段命令注册 6+silent
// 对齐 Delphi 硬编码 m_btPermission<6 静默退出）。
func TestEconFamilyRegistration(t *testing.T) {
	buildGMCommandTable(nil)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"luckypoint", "LuckPoint", 10},
		{"changejob", "ChangeJob", 10},
		{"changegender", "ChangeGender", 10},
		{"hair", "Hair", 10},
		{"bonuspoint", "BonusPoint", 10},
		{"delbonuspoint", "DelBonuPoint", 10},
		{"restbonuspoint", "RestBonuPoint", 10},
		{"setpermission", "SetPermission", 10},
		{"renewlevel", "ReNewLevel", 6},
		{"gamegold", "GameGold", 10},
		{"gamepoint", "GamePoint", 10},
		{"creditpoint", "CreditPoint", 10},
		{"training", "Training", 10},
		{"deleteskill", "DeleteSkill", 10},
		{"trainingskill", "OpTraining", 10},
		{"delgamepoint", "DelGameGold", 6},
		{"addgamepoint", "AddGamePoint", 6},
		{"adjustlevel", "AdjuestTLevel", 6},
		{"adjustexp", "AdjuestExp", 6},
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

// TestEconFamilyDispatch 权限拒绝与 "?" 帮助均被消费（不静默）。
func TestEconFamilyDispatch(t *testing.T) {
	buildGMCommandTable(nil)
	p, server := newTestPlayer()
	// permFailSilent 命令：权限不足仍被消费。
	for _, cmd := range []string{
		"renewlevel foo 1", "adjustlevel foo 10", "adjustexp foo 10",
		"addgamepoint foo 100", "delgamepoint foo 100",
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: silent permission-denied should be consumed", cmd)
		}
	}
	// perm 10 命令：低权限走"权限不够..."仍被消费。
	for _, cmd := range []string{"changejob foo Wizard", "gamegold foo = 1"} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: permission-denied should be consumed", cmd)
		}
	}
	p.Permission = 10
	for _, cmd := range []string{
		"changejob ?", "changegender ?", "hair ?", "bonuspoint ?",
		"setpermission ?", "renewlevel ?", "gamegold ?", "gamepoint ?",
		"creditpoint ?", "training ?", "deleteskill ?", "trainingskill ?",
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: usage help should be consumed", cmd)
		}
	}
}

// TestCmdChangeJob 职业字面量与客户端刷新（CmdChangeJob，ObjBase.pas:10818-10847）。
func TestCmdChangeJob(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	if !p.cmdChangeJob(server, []string{"victim1", "Wizard"}) {
		t.Fatal("changejob should be consumed")
	}
	if victim.Job != 1 {
		t.Errorf("job = %d, want 1 (Wizard)", victim.Job)
	}
	// 帮助文案短名 + 大小写不敏感（sTaos='Taoist'，M2Share.pas:180）。
	p.cmdChangeJob(server, []string{"victim1", "taos"})
	if victim.Job != 2 {
		t.Errorf("job = %d, want 2 (Taoist)", victim.Job)
	}
	p.cmdChangeJob(server, []string{"victim1", "Warrior"})
	if victim.Job != 0 {
		t.Errorf("job = %d, want 0 (Warrior)", victim.Job)
	}
	// 未知职业名：不改动（Delphi 原版误报成功，Go 输出用法）。
	p.cmdChangeJob(server, []string{"victim1", "foo"})
	if victim.Job != 0 {
		t.Errorf("unknown job changed job to %d", victim.Job)
	}
	// 离线分支消费。
	if !p.cmdChangeJob(server, []string{"nobody", "Wizard"}) {
		t.Error("changejob offline should be consumed")
	}
}

// TestCmdChangeGenderHair 性别/发型（CmdChangeGender 10749-10790、
// CmdHair 13092-13118）。
func TestCmdChangeGenderHair(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	if !p.cmdChangeGender(server, []string{"victim1", "Woman"}) {
		t.Fatal("changegender should be consumed")
	}
	if victim.Gender != 1 {
		t.Errorf("gender = %d, want 1", victim.Gender)
	}
	// 已是该性别 → "未改变"分支。
	p.cmdChangeGender(server, []string{"victim1", "Female"})
	if victim.Gender != 1 {
		t.Errorf("gender = %d, want 1 unchanged", victim.Gender)
	}
	// 非法字面量 → 用法分支。
	p.cmdChangeGender(server, []string{"victim1", "X"})
	if victim.Gender != 1 {
		t.Errorf("invalid literal changed gender to %d", victim.Gender)
	}
	// 离线分支文案 '%s没有在线！！！'。
	if !p.cmdChangeGender(server, []string{"nobody", "Man"}) {
		t.Error("changegender offline should be consumed")
	}

	if !p.cmdHair(server, []string{"victim1", "3"}) {
		t.Fatal("hair should be consumed")
	}
	if victim.Hair != 3 {
		t.Errorf("hair = %d, want 3", victim.Hair)
	}
	// 负值 → 用法分支；不带类型值按 0（Str_ToInt(sParam2, 0)）。
	p.cmdHair(server, []string{"victim1", "-1"})
	if victim.Hair != 3 {
		t.Errorf("negative hair changed hair to %d", victim.Hair)
	}
	p.cmdHair(server, []string{"victim1"})
	if victim.Hair != 0 {
		t.Errorf("hair without value = %d, want 0", victim.Hair)
	}
}

// TestCmdBonusPointFamily 属性点三命令（CmdBonuPoint 10618-10660、
// CmdDelBonuPoint 11526-11554、CmdRestBonuPoint 11582-11623）。
func TestCmdBonusPointFamily(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	if !p.cmdBonusPoint(server, []string{"victim1", "12"}) {
		t.Fatal("bonuspoint set should be consumed")
	}
	if victim.BonusPoint != 12 {
		t.Errorf("bonuspoint = %d, want 12", victim.BonusPoint)
	}
	// 无点数参数 = 查看分支。
	if !p.cmdBonusPoint(server, []string{"victim1"}) {
		t.Error("bonuspoint view should be consumed")
	}
	if victim.BonusPoint != 12 {
		t.Errorf("view changed bonuspoint to %d", victim.BonusPoint)
	}
	if !p.cmdDelBonusPoint(server, []string{"victim1"}) {
		t.Fatal("delbonuspoint should be consumed")
	}
	if victim.BonusPoint != 0 {
		t.Errorf("delbonuspoint left %d, want 0", victim.BonusPoint)
	}
	// Go 无分项分配模型：复位不改变未分配点数。
	victim.BonusPoint = 7
	if !p.cmdRestBonusPoint(server, []string{"victim1"}) {
		t.Fatal("restbonuspoint should be consumed")
	}
	if victim.BonusPoint != 7 {
		t.Errorf("restbonuspoint changed points to %d, want 7", victim.BonusPoint)
	}
}

// TestCmdSetPermissionReNewLevel 权限/转生（CmdSetPermission 11718-11748、
// CmdReNewLevel 11555-11581）。
func TestCmdSetPermissionReNewLevel(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	if !p.cmdSetPermission(server, []string{"victim1", "5"}) {
		t.Fatal("setpermission should be consumed")
	}
	if victim.Permission != 5 {
		t.Errorf("permission = %d, want 5", victim.Permission)
	}
	// 越界等级 → 用法分支。
	p.cmdSetPermission(server, []string{"victim1", "11"})
	if victim.Permission != 5 {
		t.Errorf("out-of-range permission changed to %d", victim.Permission)
	}

	if !p.cmdReNewLevel(server, []string{"victim1", "3"}) {
		t.Fatal("renewlevel set should be consumed")
	}
	if victim.ReNewLevel != 3 {
		t.Errorf("renewlevel = %d, want 3", victim.ReNewLevel)
	}
	// 查看分支。
	if !p.cmdReNewLevel(server, []string{"victim1"}) {
		t.Error("renewlevel view should be consumed")
	}
	if victim.ReNewLevel != 3 {
		t.Errorf("view changed renewlevel to %d", victim.ReNewLevel)
	}
	// 离线分支文案 '%s 没在线上！！！'。
	if !p.cmdReNewLevel(server, []string{"nobody"}) {
		t.Error("renewlevel offline should be consumed")
	}
}

// TestCmdGameGoldPointCredit 经济三命令（CmdGameGold 10280-10331、
// CmdGamePoint 10332-10384、CmdCreditPoint 10385-10435）。
func TestCmdGameGoldPointCredit(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	victim.GameGold = 100
	if !p.cmdGameGold(server, []string{"victim1", "+", "50"}) {
		t.Fatal("gamegold add should be consumed")
	}
	if victim.GameGold != 150 {
		t.Errorf("gamegold = %d, want 150", victim.GameGold)
	}
	p.cmdGameGold(server, []string{"victim1", "=", "1000"})
	if victim.GameGold != 1000 {
		t.Errorf("gamegold '=' = %d, want 1000", victim.GameGold)
	}
	p.cmdGameGold(server, []string{"victim1", "-", "400"})
	if victim.GameGold != 600 {
		t.Errorf("gamegold '-' = %d, want 600", victim.GameGold)
	}
	// 非法控制符/超输入上限 → 用法分支，数值不变。
	p.cmdGameGold(server, []string{"victim1", "*", "5"})
	p.cmdGameGold(server, []string{"victim1", "+", "300000000"})
	if victim.GameGold != 600 {
		t.Errorf("invalid args changed gamegold to %d", victim.GameGold)
	}

	victim.GamePoint = 10
	if !p.cmdGamePoint(server, []string{"victim1", "+", "5"}) {
		t.Fatal("gamepoint add should be consumed")
	}
	if victim.GamePoint != 15 {
		t.Errorf("gamepoint = %d, want 15", victim.GamePoint)
	}
	p.cmdGamePoint(server, []string{"victim1", "=", "99"})
	if victim.GamePoint != 99 {
		t.Errorf("gamepoint '=' = %d, want 99", victim.GamePoint)
	}
	p.cmdGamePoint(server, []string{"victim1", "+", "200000000"})
	if victim.GamePoint != 99 {
		t.Errorf("over-cap input changed gamepoint to %d", victim.GamePoint)
	}

	victim.CreditPoint = 250
	if !p.cmdCreditPoint(server, []string{"victim1", "+", "10"}) {
		t.Fatal("creditpoint add should be consumed")
	}
	if victim.CreditPoint != 255 {
		t.Errorf("creditpoint = %d, want 255 (clamped)", victim.CreditPoint)
	}
	// Delphi 输入上限 255（High(Byte)）："- 300" 会被拒绝；
	// 合法入参 255 使结果 -5，钳制到 0。
	p.cmdCreditPoint(server, []string{"victim1", "-", "300"})
	if victim.CreditPoint != 255 {
		t.Errorf("over-cap input changed creditpoint to %d", victim.CreditPoint)
	}
	p.cmdCreditPoint(server, []string{"victim1", "-", "255"})
	if victim.CreditPoint != 0 {
		t.Errorf("creditpoint = %d, want 0 (clamped)", victim.CreditPoint)
	}
	p.cmdCreditPoint(server, []string{"victim1", "=", "120"})
	if victim.CreditPoint != 120 {
		t.Errorf("creditpoint '=' = %d, want 120", victim.CreditPoint)
	}
	p.cmdCreditPoint(server, []string{"victim1", "+", "256"})
	if victim.CreditPoint != 120 {
		t.Errorf("over-cap input changed creditpoint to %d", victim.CreditPoint)
	}
}

// TestCmdAddDelGamePoint 游戏点增扣（CmdAddGameGold 10222-10252、
// CmdDelGameGold 10253-10279；上限 2000000）。
func TestCmdAddDelGamePoint(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	victim.GamePoint = 1999990
	if !p.cmdAddGamePoint(server, []string{"victim1", "100"}) {
		t.Fatal("addgamepoint should be consumed")
	}
	if victim.GamePoint != 2000000 {
		t.Errorf("gamepoint = %d, want cap 2000000", victim.GamePoint)
	}
	// 已在封顶：实际增量为 0。
	p.cmdAddGamePoint(server, []string{"victim1", "50"})
	if victim.GamePoint != 2000000 {
		t.Errorf("gamepoint over cap = %d", victim.GamePoint)
	}
	// 非法数量 → 用法/静默分支。
	p.cmdAddGamePoint(server, []string{"victim1", "0"})
	if victim.GamePoint != 2000000 {
		t.Errorf("zero add changed gamepoint to %d", victim.GamePoint)
	}

	p.cmdDelGamePoint(server, []string{"victim1", "150"})
	if victim.GamePoint != 1999850 {
		t.Errorf("gamepoint after del = %d, want 1999850", victim.GamePoint)
	}
	// 超额扣减只扣到 0。
	p.cmdDelGamePoint(server, []string{"victim1", "99999999"})
	if victim.GamePoint != 0 {
		t.Errorf("gamepoint after over-del = %d, want 0", victim.GamePoint)
	}
	// 非法参数静默（Delphi ObjBase.pas:10259）。
	if !p.cmdDelGamePoint(server, []string{"victim1", "-5"}) {
		t.Error("delgamepoint invalid should be consumed silently")
	}
}

// TestCmdTrainingSkillFamily 技能三命令（CmdTrainingSkill 10177-10221、
// CmdTrainingMagic 10130-10176、CmdDelSkill 12415-12465）。
func TestCmdTrainingSkillFamily(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	db := newEconMagicDB()
	p.MagicDB = db
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)
	victim.MagicDB = db
	victim.LearnedMagics = []*PlayerMagic{{MagID: 1, Level: 0}}

	// @Training：调整已学技能（技能名查 MagicDB）。
	if !p.cmdTraining(server, []string{"victim1", "火球术", "9"}) {
		t.Fatal("training should be consumed")
	}
	if got := victim.findMagic(1); got == nil || got.Level != 3 {
		t.Errorf("training level = %v, want clamped 3", got)
	}
	// 未学技能 → 静默退出。
	if !p.cmdTraining(server, []string{"victim1", "雷电术", "2"}) {
		t.Error("training unlearned skill should be consumed silently")
	}
	if victim.findMagic(6) != nil {
		t.Error("training must not add new skills")
	}

	// @TrainingSkill：新增技能。
	if !p.cmdTrainingSkill(server, []string{"victim1", "雷电术", "1"}) {
		t.Fatal("trainingskill should be consumed")
	}
	if got := victim.findMagic(6); got == nil || got.Level != 1 {
		t.Errorf("new skill = %v, want level 1", got)
	}
	// 重复修炼分支：等级不变。
	p.cmdTrainingSkill(server, []string{"victim1", "雷电术", "2"})
	if got := victim.findMagic(6); got == nil || got.Level != 1 {
		t.Errorf("duplicate training changed level to %v", got)
	}
	// 技能名不正确分支。
	p.cmdTrainingSkill(server, []string{"victim1", "不存在", "1"})
	if len(victim.LearnedMagics) != 2 {
		t.Errorf("bad skill name added magic, list = %d", len(victim.LearnedMagics))
	}
	// 等级越界 → 用法分支。
	p.cmdTrainingSkill(server, []string{"victim1", "抗拒火环", "4"})
	if victim.findMagic(15) != nil {
		t.Error("out-of-range level should not add skill")
	}

	// @DeleteSkill：单个删除。
	if !p.cmdDeleteSkill(server, []string{"victim1", "火球术"}) {
		t.Fatal("deleteskill should be consumed")
	}
	if victim.findMagic(1) != nil || len(victim.LearnedMagics) != 1 {
		t.Errorf("deleteskill left %d magics", len(victim.LearnedMagics))
	}
	// All = 全删。
	if !p.cmdDeleteSkill(server, []string{"victim1", "all"}) {
		t.Fatal("deleteskill all should be consumed")
	}
	if len(victim.LearnedMagics) != 0 {
		t.Errorf("deleteskill all left %d magics", len(victim.LearnedMagics))
	}
}

// TestCmdAdjustLevelExp 等级/经验调整（CmdAdjuestLevel 10568-10598、
// CmdAdjuestExp 10535-10567）。
func TestCmdAdjustLevelExp(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)

	if !p.cmdAdjustLevel(server, []string{"victim1", "10"}) {
		t.Fatal("adjustlevel should be consumed")
	}
	if victim.WAbil.Level != 10 {
		t.Errorf("level = %d, want 10", victim.WAbil.Level)
	}
	if victim.WAbil.HP != victim.WAbil.MaxHP {
		t.Errorf("HP = %d/%d after adjust", victim.WAbil.HP, victim.WAbil.MaxHP)
	}
	// 下限钳制 1（_MAX(1, _MIN(MAXUPLEVEL, n))）。
	p.cmdAdjustLevel(server, []string{"victim1", "0"})
	if victim.WAbil.Level != 1 {
		t.Errorf("level = %d, want clamp 1", victim.WAbil.Level)
	}
	// 上限 MAXUPLEVEL=65535。
	p.cmdAdjustLevel(server, []string{"victim1", "999999"})
	if victim.WAbil.Level != 65535 {
		t.Errorf("level = %d, want clamp 65535", victim.WAbil.Level)
	}
	// 离线分支。
	if !p.cmdAdjustLevel(server, []string{"nobody", "5"}) {
		t.Error("adjustlevel offline should be consumed")
	}

	victim.WAbil.Level = 1
	victim.WAbil.Exp = 0
	victim.BonusPoint = 0
	if !p.cmdAdjustExp(server, []string{"victim1", "100000"}) {
		t.Fatal("adjustexp should be consumed")
	}
	// 升级链路：100000 经验远超 1 级需求，应多级提升且余数 < 当前级需求。
	if victim.WAbil.Level <= 1 {
		t.Errorf("level = %d after 100000 exp, want > 1", victim.WAbil.Level)
	}
	if victim.WAbil.Exp >= victim.GetMaxExp() {
		t.Errorf("exp = %d still above max %d", victim.WAbil.Exp, victim.GetMaxExp())
	}
	if victim.BonusPoint == 0 {
		t.Error("level-ups should grant bonus points")
	}
	// 归零经验不改变等级。
	lv := victim.WAbil.Level
	p.cmdAdjustExp(server, []string{"victim1", "0"})
	if victim.WAbil.Exp != 0 || victim.WAbil.Level != lv {
		t.Errorf("exp=%d level=%d after reset, want exp 0 level %d",
			victim.WAbil.Exp, victim.WAbil.Level, lv)
	}
}

// TestCmdLuckyPoint 幸运点查看（CmdLuckPoint，ObjBase.pas:13328-13357）。
func TestCmdLuckyPoint(t *testing.T) {
	buildGMCommandTable(nil)
	p, env, engine, _ := newMoveFamilyFixture(t)
	server := &netserver.TCPServer{}
	victim := newMoveFamilyVictim(t, "victim1", 2, env, engine, 6, 5)
	victim.Luck = 2

	if !p.cmdLuckyPoint(server, []string{"victim1"}) {
		t.Fatal("luckypoint should be consumed")
	}
	// 离线分支。
	if !p.cmdLuckyPoint(server, []string{"nobody"}) {
		t.Error("luckypoint offline should be consumed")
	}
}
