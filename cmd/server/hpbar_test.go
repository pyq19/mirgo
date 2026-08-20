package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// TestObjectCharStatusShowHPBit 验证 ShowHP 时 charStatus 并入
// STATE_OPENHEATH（GetCharStatus，ObjBase.pas:20074-20088）。
func TestObjectCharStatusShowHPBit(t *testing.T) {
	p := &PlayObject{BaseObject: NewBaseObject("target", 1)}
	if s := objectCharStatus(p); s&STATE_OPENHEATH != 0 {
		t.Fatalf("ShowHP=false 时不应含状态位: %x", s)
	}
	p.ShowHP = true
	if s := objectCharStatus(p); s&STATE_OPENHEATH == 0 {
		t.Fatalf("ShowHP=true 时应含 STATE_OPENHEATH: %x", s)
	}
}

// TestBroadcastStruckMsgCarriesHP 验证受击广播按 Delphi 语义携带
// 受击后 HP/MaxHP/伤害/魔法标志（RM_STRUCK → SM_STRUCK 的入参，
// ObjBase.pas:5508-5518）。
func TestBroadcastStruckMsgCarriesHP(t *testing.T) {
	env := newTestEnv(40, 40)

	victim := placeMonster(env, 2, 20, 20)
	victim.WAbil.HP = 60
	victim.WAbil.MaxHP = 100

	observer := &PlayObject{BaseObject: NewBaseObject("observer", 1)}
	observer.Session = &netserver.Session{ID: 1}
	observer.envir = env
	observer.CurrX, observer.CurrY = 21, 20
	env.AddObject(21, 20, OS_MOVINGOBJECT, observer)

	env.broadcastStruckMsg(victim.BaseObject, victim.ID, int(victim.WAbil.HP), int(victim.WAbil.MaxHP), 40, 1)

	msg, ok := observer.GetMsg()
	if !ok {
		t.Fatal("观察者应收到 RM_STRUCK")
	}
	if msg.Ident != RM_STRUCK {
		t.Fatalf("Ident = %d, want RM_STRUCK(%d)", msg.Ident, RM_STRUCK)
	}
	if msg.SourceID != victim.ID {
		t.Errorf("SourceID = %d, want 受害者 ID %d", msg.SourceID, victim.ID)
	}
	if msg.Param1 != 60 || msg.Param2 != 100 || msg.Param3 != 40 || msg.Dir != 1 {
		t.Errorf("Param1/2/3/Dir = %d/%d/%d/%d, want 60/100/40/1",
			msg.Param1, msg.Param2, msg.Param3, msg.Dir)
	}
}

// TestMakeBreakOpenHealthFlags 验证开关方法的标志位语义
//（ObjBase.pas:3595-3612）。无 envir 时广播安全地跳过。
func TestMakeBreakOpenHealthFlags(t *testing.T) {
	p := &PlayObject{BaseObject: NewBaseObject("target", 1)}
	p.MakeOpenHealth() // envir=nil：仅置位不广播
	if !p.ShowHP {
		t.Fatal("MakeOpenHealth 应置 ShowHP")
	}
	p.BreakOpenHealth()
	if p.ShowHP {
		t.Fatal("BreakOpenHealth 应清 ShowHP")
	}
	// 重复关闭应为幂等无操作
	p.BreakOpenHealth()
}

// TestGetShowHPOnStruckDefault 验证受击开血条配置的默认语义：
// 未配置(nil)默认开启；显式 false 关闭（Go 专属配置项）。
func TestGetShowHPOnStruckDefault(t *testing.T) {
	if got := (&ServerConfig{}).GetShowHPOnStruck(); !got {
		t.Fatal("未配置时默认应开启")
	}
	off := false
	on := true
	cfg := &ServerConfig{}
	cfg.Game.ShowHPOnStruck = &off
	if cfg.GetShowHPOnStruck() {
		t.Fatal("显式 false 应关闭")
	}
	cfg.Game.ShowHPOnStruck = &on
	if !cfg.GetShowHPOnStruck() {
		t.Fatal("显式 true 应开启")
	}
}

// TestStruckOpenHealth 验证受击开启血条：置位+设置到期参数，
// 重复受击仅刷新到期时刻不改变标志。
func TestStruckOpenHealth(t *testing.T) {
	p := &PlayObject{BaseObject: NewBaseObject("target", 1)}
	if p.ShowHPInterval != 0 || p.ShowHPTick != 0 {
		t.Fatal("初始不应有到期参数")
	}
	p.StruckOpenHealth() // envir=nil：仅置位
	if !p.ShowHP {
		t.Fatal("StruckOpenHealth 应置 ShowHP")
	}
	if p.ShowHPInterval != showHPStruckDurationMs {
		t.Errorf("ShowHPInterval = %d, want %d", p.ShowHPInterval, showHPStruckDurationMs)
	}
	first := p.ShowHPTick
	if first == 0 {
		t.Fatal("ShowHPTick 应置为当前时刻")
	}
	// 重复受击：标志保持，到期时刻刷新（不回退）
	p.StruckOpenHealth()
	if !p.ShowHP || p.ShowHPTick < first {
		t.Error("重复受击应保持 ShowHP 并刷新到期时刻")
	}
}

// TestGiveItemReturnsInstance 验证 GiveItem 返回创建的 UserItem
//（拾取流程用它构建 SM_ADDITEM 的 TClientItem body）。
func TestGiveItemReturnsInstance(t *testing.T) {
	p := &PlayObject{BaseObject: NewBaseObject("picker", 1)}
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	p.ItemDB = &ItemDB{byIdx: map[int]*ItemDef{5: {Idx: 5, Name: "金创药", DuraMax: 1000}}}

	item, ok := p.GiveItem(5)
	if !ok || item == nil {
		t.Fatal("GiveItem 应成功并返回实例")
	}
	if item.WIndex != 5 || item.DuraMax != 1000 {
		t.Errorf("item = WIndex %d DuraMax %d, want 5/1000", item.WIndex, item.DuraMax)
	}
	if len(p.ItemList) != 1 || p.ItemList[0] != item {
		t.Error("返回的实例应已入包")
	}

	// 不存在的物品
	if _, ok := p.GiveItem(999); ok {
		t.Error("不存在的物品不应给出")
	}
}

// TestGetRPow 验证 Delphi GetRPow（Magic.pas:247-252）打包字取值。
func TestGetRPow(t *testing.T) {
	// 未打包（高字=0）→ 直接取低字
	if got := getRPow(35); got != 35 {
		t.Errorf("getRPow(35) = %d, want 35", got)
	}
	// 打包 [10, 20] → 结果在区间内
	for i := 0; i < 1000; i++ {
		got := getRPow(10 | 20<<16)
		if got < 10 || got > 20 {
			t.Fatalf("getRPow([10,20]) = %d, 超出区间", got)
		}
	}
}

// TestMagicDuration13ZeroPower 验证技能 28（心灵启示）时长计算在
// defPower=defMaxPower=0（magic_db.jsonc 实际值）时不 panic
//（Delphi Random(0)=0，Go rand.Intn(0) 会 panic，必须 guard）。
func TestMagicDuration13ZeroPower(t *testing.T) {
	p := &PlayObject{BaseObject: NewBaseObject("taoist", 1)}
	p.MagicDB = &MagicDB{
		Magics: []MagicDef{{MagID: 28, MagName: "心灵启示", MaxTrainLv: 3}},
		byID:   map[int]*MagicDef{28: {MagID: 28, MagName: "心灵启示", MaxTrainLv: 3}},
	}
	p.LearnedMagics = []*PlayerMagic{{MagID: 28, Level: 1}}
	p.WAbil.SC = 5 | 8<<16

	d := p.magicDuration13(28)
	if d <= 0 {
		t.Fatalf("duration = %d, want > 0", d)
	}
	// 无 MagicDB 时回退默认参数也不 panic
	p2 := &PlayObject{BaseObject: NewBaseObject("taoist2", 2)}
	p2.WAbil.SC = 5
	if d2 := p2.magicDuration13(28); d2 <= 0 {
		t.Fatalf("fallback duration = %d, want > 0", d2)
	}
}
