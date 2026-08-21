package main

import (
	"testing"
	"time"
)

// TestDebugCommandRegistration 调试杂项族注册校验（字面量含中文
// '祈祷生效'/'停止叛变'，M2Share.pas:2913-2914；内部 <6 静默退出的
// 命令注册 6+silent，ServerStatus 内部 <10 静默注册 10+silent）。
func TestDebugCommandRegistration(t *testing.T) {
	ensureGMTable(t)
	cases := []struct {
		literal string
		key     string
		perm    int
	}{
		{"backstep", "Backstep", 6},
		{"hunger", "Hunger", 6},
		{"disablefilter", "DisableFilter", 6},
		{"chguserfull", "ChgUserFull", 6},
		{"chgzenfaststep", "ChgZenFastStep", 6},
		{"teststatus", "TestStatus", 6},
		{"祈祷生效", "Spirit", 6},
		{"停止叛变", "SpiritStop", 6},
		{"serverstatus", "ServerStatus", 10},
		{"testgetbagitem", "TestGetBagItem", 10},
		{"testspeedmode", "TestSpeedMode", 10},
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
	// Reconnection/TestServerConfig 未实装（Go 无网关换连/SM_SERVERCONFIG）。
	for _, literal := range []string{"reconnection", "testserverconfig"} {
		if gmCmdTable[literal] != nil {
			t.Errorf("@%s should not be registered (no Go counterpart)", literal)
		}
	}
}

// TestCmdChgUserFullZenFastStep 运行时改 UserFull/ZenFastStep
// （Delphi ObjBase.pas:11015-11047，仅内存）。
func TestCmdChgUserFullZenFastStep(t *testing.T) {
	ensureGMTable(t)
	p, _ := newTestPlayer()
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	_, server := newTestPlayer()

	p.cmdChgUserFull(server, []string{"500"})
	if p.Engine.Config.Game.UserFull != 500 {
		t.Errorf("UserFull = %d, want 500", p.Engine.Config.Game.UserFull)
	}
	if p.Engine.Config.GetUserFull() != 500 {
		t.Errorf("GetUserFull = %d, want 500", p.Engine.Config.GetUserFull())
	}
	p.cmdChgUserFull(server, []string{"?"}) // 帮助分支不改值
	if p.Engine.Config.Game.UserFull != 500 {
		t.Error("'?' help must not change UserFull")
	}
	p.cmdChgUserFull(server, []string{"0"}) // <1 拒绝
	if p.Engine.Config.Game.UserFull != 500 {
		t.Error("invalid count must not change UserFull")
	}

	p.cmdChgZenFastStep(server, []string{"120"})
	if p.Engine.Config.Game.ZenFastStep != 120 {
		t.Errorf("ZenFastStep = %d, want 120", p.Engine.Config.Game.ZenFastStep)
	}
}

// TestCmdDisableFilter 文字过滤总开关翻转（Delphi boFilterWord）。
func TestCmdDisableFilter(t *testing.T) {
	ensureGMTable(t)
	p, _ := newTestPlayer()
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	_, server := newTestPlayer()
	old := wordFilterIsEnabled()
	t.Cleanup(func() { setWordFilterEnabled(old) })

	p.cmdDisableFilter(server, nil)
	if wordFilterIsEnabled() == old {
		t.Error("@DisableFilter should toggle the filter switch")
	}
	p.cmdDisableFilter(server, nil)
	if wordFilterIsEnabled() != old {
		t.Error("second @DisableFilter should restore the switch")
	}
}

// TestCmdTestSpeedMode 速度测试模式开关 + checkActionSpeed 豁免。
func TestCmdTestSpeedMode(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Engine = &UserEngine{Config: &ServerConfig{}}

	p.cmdTestSpeedMode(server, nil)
	if !p.TestSpeedMode {
		t.Fatal("@TestSpeedMode should enable the mode")
	}
	// 过期 tick（正常会被拒）在测试模式下豁免。
	tick := int64(1 << 40)
	if !p.checkActionSpeed(0, 1000, &tick, server) {
		t.Error("TestSpeedMode should bypass action-speed checks")
	}
	if tick != 0 {
		t.Errorf("tick = %d, want refreshed to now(0)", tick)
	}
	p.cmdTestSpeedMode(server, nil)
	if p.TestSpeedMode {
		t.Error("second @TestSpeedMode should disable the mode")
	}
}

// TestCmdSpiritMutiny 祈祷生效/停止叛变写 SpiritMutinyTick，
// slaveFollow 生效期清除主从纽带。
func TestCmdSpiritMutiny(t *testing.T) {
	ensureGMTable(t)
	p, _ := newTestPlayer()
	e := &UserEngine{Config: &ServerConfig{}}
	p.Engine = e
	_, server := newTestPlayer()

	p.cmdSpiritStart(server, []string{"10"})
	now := time.Now().UnixMilli()
	if e.SpiritMutinyTick <= now || e.SpiritMutinyTick > now+11000 {
		t.Errorf("SpiritMutinyTick = %d, want now+~10s", e.SpiritMutinyTick)
	}
	// 生效期宠物叛变：slaveFollow 清除 PlayerMasterID。
	master, _ := newTestPlayer()
	master.ID = 7
	e.PlayObjectList = map[int32]*PlayObject{master.ID: master}
	mon := &MonsterObject{BaseObject: NewBaseObject("pet", 99)}
	mon.PlayerMasterID = master.ID
	mon.slaveFollow(e, now)
	if mon.PlayerMasterID != 0 {
		t.Error("active spirit mutiny should break the master bond")
	}

	p.cmdSpiritStop(server, nil)
	if e.SpiritMutinyTick != 0 {
		t.Error("@停止叛变 should clear SpiritMutinyTick")
	}
}

// TestCmdHunger 设置目标能量值（Delphi ObjBase.pas:13167-13190）。
func TestCmdHunger(t *testing.T) {
	ensureGMTable(t)
	p, _ := newTestPlayer()
	target, _ := newTestPlayer()
	target.Name = "victim"
	target.ID = 2
	p.Engine = &UserEngine{Config: &ServerConfig{}, PlayObjectList: map[int32]*PlayObject{target.ID: target}}
	_, server := newTestPlayer()

	p.cmdHunger(server, []string{"victim", "123"})
	if target.HungerStatus != 123 {
		t.Errorf("HungerStatus = %d, want 123", target.HungerStatus)
	}
	// 离线目标不 panic；负值出用法分支。
	p.cmdHunger(server, []string{"nobody", "1"})
	p.cmdHunger(server, []string{"victim", "-1"})
	if target.HungerStatus != 123 {
		t.Error("invalid hunger params must not change the value")
	}
}

// TestCmdBackstep 向后推（nType=0 向身后，Dir 保持；nType>0 随机方向）。
func TestCmdBackstep(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	p.envir = newTestEnv(10, 10)
	p.CurrX, p.CurrY = 5, 5
	p.Dir = 0
	dx, dy := dirToOffset((p.Dir + 4) % 8)

	p.cmdBackstep(server, []string{"0", "1"})
	if p.CurrX != 5+dx || p.CurrY != 5+dy {
		t.Errorf("backstep landed (%d,%d), want (%d,%d)", p.CurrX, p.CurrY, 5+dx, 5+dy)
	}
	if p.Dir != 0 {
		t.Errorf("Dir = %d, want unchanged 0 for nType=0", p.Dir)
	}
}

// TestCmdTestStatus 状态时长按 tick 间隔换算写入并广播。
func TestCmdTestStatus(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	p.envir = newTestEnv(4, 4)

	p.cmdTestStatus(server, []string{"5", "30"})
	if p.StatusTimeArr[5] <= 0 {
		t.Errorf("StatusTimeArr[5] = %d, want > 0", p.StatusTimeArr[5])
	}
	// 越界类型出用法分支不写数组。
	p.cmdTestStatus(server, []string{"12", "30"})
	if len(p.StatusTimeArr) != 12 {
		t.Error("out-of-range type must not resize the array")
	}
}

// TestCmdTestGetBagItem 黑铁矿与首饰被消耗并回显合计。
func TestCmdTestGetBagItem(t *testing.T) {
	ensureGMTable(t)
	p, server := newTestPlayer()
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	// ItemDB=nil：无物品可扫描，输出全 0 且不 panic。
	p.cmdTestGetBagItem(server, nil)
}
