package main

import (
	"testing"
)

func ensureGMTable(t *testing.T) {
	t.Helper()
	buildGMCommandTable(nil)
}

// TestSplitGMParams GetValidStr3 语义（HUtil32.pas:695）：
// 空格/冒号/逗号/TAB 分隔、连续分隔符合并、上限 8 段。
func TestSplitGMParams(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"move 10 20", []string{"move", "10", "20"}},
		{"make 木剑,2", []string{"make", "木剑", "2"}},
		{"make 木剑:2", []string{"make", "木剑", "2"}},
		{"a  b,,c", []string{"a", "b", "c"}},
		{"a\tb c", []string{"a", "b", "c"}},
		{"", nil},
	}
	for _, c := range cases {
		got := splitGMParams(c.in)
		if len(got) == 0 && len(c.want) == 0 {
			continue
		}
		if len(got) != len(c.want) {
			t.Errorf("splitGMParams(%q) = %v, want %v", c.in, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("splitGMParams(%q) = %v, want %v", c.in, got, c.want)
				break
			}
		}
	}
	if got := splitGMParams("1 2 3 4 5 6 7 8 9 10"); len(got) != 8 {
		t.Errorf("param cap = %d fields, want 8", len(got))
	}
}

// TestBuildGMCommandTable 默认字面量/权限解析与 config 覆盖。
func TestBuildGMCommandTable(t *testing.T) {
	ensureGMTable(t)

	if def := gmCmdTable["move"]; def == nil || def.key != "UserMove" {
		t.Fatalf("move should resolve to USERMOVE def, got %+v", def)
	}
	if def := gmCmdTable["骑马"]; def == nil || def.key != "TakeOnHorse" {
		t.Fatalf("骑马 should resolve to TakeOnHorse def, got %+v", def)
	}
	if def := gmCmdTable["pkpoint"]; def == nil || def.key != "PkPoint" {
		t.Fatalf("pkpoint should resolve to PkPoint def, got %+v", def)
	}
	// GM 命令权限默认值（M2Share.pas:2724-2926 对应项）。
	if def := gmCmdTable["kick"]; def == nil || def.resolvedPerm != 10 {
		t.Errorf("kick perm = %v, want 10", def)
	}
	if def := gmCmdTable["info"]; def == nil || def.resolvedPerm != 3 {
		t.Errorf("info perm = %v, want 3", def)
	}
	if got := gmPermByKey("Move", -1); got != 3 {
		t.Errorf("gmPermByKey(Move) = %d, want 3", got)
	}

	// config 覆盖：改名 + 改权限。
	cfg := DefaultConfig()
	cfg.Commands.Names = map[string]string{"Kick": "踢人"}
	cfg.Commands.Permissions = map[string]int{"Kick": 8}
	buildGMCommandTable(cfg)
	if def := gmCmdTable["踢人"]; def == nil || def.key != "Kick" || def.resolvedPerm != 8 {
		t.Fatalf("config override failed: %+v", gmCmdTable["踢人"])
	}
	if _, ok := gmCmdTable["kick"]; ok {
		t.Error("old literal should be replaced after rename")
	}
	buildGMCommandTable(nil) // 还原
}

// TestHandleGMCommandDispatch 分发兜底行为：
// 未知命令与权限不足都被消费（返回 true），不再静默。
func TestHandleGMCommandDispatch(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()

	if !p.HandleGMCommand("不存在的命令 xyz", server) {
		t.Error("unknown command should be consumed (fallback message)")
	}
	// 普通玩家用 GM 命令：被消费并提示"权限不够..."。
	if !p.HandleGMCommand("kick 某人", server) {
		t.Error("permission-denied command should be consumed")
	}
	// 空输入不消费。
	if p.HandleGMCommand("   ", server) {
		t.Error("blank command should not be consumed")
	}
}

// TestGMRedBroadcastPerm "@!..." 红字广播权限门槛（默认 6）。
func TestGMRedBroadcastPerm(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	// 权限不足：消费并回兜底提示（Delphi 落回普通分发→兜底）。
	if !p.HandleGMCommand("!全服喊话", server) {
		t.Error("red broadcast under perm should be consumed")
	}
}

// TestCmdAttackMode 攻击模式命令（Delphi CmdChangeAttackMode，
// ObjBase.pas:10679-10708）：无参循环、有参指定、越界回退循环。
func TestCmdAttackMode(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()

	p.AttackMode = AttackModePeace
	p.cmdAttackMode(server, nil)
	if p.AttackMode != AttackModeGroup {
		t.Errorf("cycle from peace = %d, want group", p.AttackMode)
	}
	p.AttackMode = AttackModePK
	p.cmdAttackMode(server, nil)
	if p.AttackMode != AttackModePeace {
		t.Errorf("wrap from pk = %d, want peace", p.AttackMode)
	}
	p.cmdAttackMode(server, []string{"4"})
	if p.AttackMode != AttackModePK {
		t.Errorf("set 4 = %d, want pk", p.AttackMode)
	}
	p.cmdAttackMode(server, []string{"9"})
	if p.AttackMode != AttackModePeace {
		t.Errorf("invalid param should cycle, got %d", p.AttackMode)
	}
}

// TestCmdGoldClamp 金币命令钳制（Delphi nHumanMaxGold）。
func TestCmdGoldClamp(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Engine = nil
	p.Gold = 0
	// Engine 为 nil 时 clampGold 用默认上限，负数钳 0。
	p.cmdGold(server, []string{"-100"})
	if p.Gold != 0 {
		t.Errorf("gold after -100 = %d, want 0", p.Gold)
	}
}
