package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// 物品族测试（gmcmd_item.go）。
// Delphi 参照行号见各命令注释。

// newItemFamilyFixture GM + 地图 "0" + 引擎 + 物品库的最小可运行环境。
func newItemFamilyFixture(t *testing.T) (*PlayObject, *UserEngine) {
	t.Helper()
	env := newTestEnv(20, 20)
	env.Name = "0"
	mm := NewMapManager("")
	mm.maps["0"] = env
	engine := NewUserEngine(nil, mm)
	engine.Config = &ServerConfig{}
	gm := &PlayObject{BaseObject: NewBaseObject("gm", 1)}
	gm.Session = &netserver.Session{ID: 1}
	gm.Permission = 10
	gm.Engine = engine
	gm.MapMgr = mm
	gm.envir = env
	gm.MapName = env.Name
	gm.CurrX, gm.CurrY = 5, 5
	gm.WAbil.Level = 10
	gm.WAbil.HP, gm.WAbil.MaxHP = 100, 100
	gm.ItemDB = newTestItemDB(
		ItemDef{Idx: 1, Name: "木剑", StdMode: 5, DuraMax: 8000, DC: 1, DCMax: 3},
		ItemDef{Idx: 2, Name: "金创药", StdMode: 0},
		ItemDef{Idx: 3, Name: "IronSword", StdMode: 5, DuraMax: 10000},
	)
	env.AddObject(5, 5, OS_MOVINGOBJECT, gm)
	engine.AddPlayer(gm)
	return gm, engine
}

// newItemFamilyVictim 在引擎中注册一名带背包的目标玩家。
func newItemFamilyVictim(t *testing.T, name string, id int32, engine *UserEngine, db *ItemDB) *PlayObject {
	t.Helper()
	pl := &PlayObject{BaseObject: NewBaseObject(name, id)}
	pl.Session = &netserver.Session{ID: int64(id)}
	pl.Engine = engine
	pl.ItemDB = db
	pl.WAbil.HP, pl.WAbil.MaxHP = 50, 50
	engine.AddPlayer(pl)
	return pl
}

// resetCustomItemNames 清空自定义物品名登记表（测试隔离）。
func resetCustomItemNames() {
	gmCustomItemNameMu.Lock()
	gmCustomItemNames = make(map[int32]string)
	gmCustomItemNameMu.Unlock()
}

// TestItemFamilyRegistration 注册表字面量/权限/失败模式
//（g_GameCommand，M2Share.pas:2749,2821,2837,2843,2858）。
func TestItemFamilyRegistration(t *testing.T) {
	buildGMCommandTable(nil)
	cases := []struct {
		literal string
		key     string
		perm    int
		fail    gmPermFailMode
	}{
		{"supermake", "SuperMake", 10, permFailLow},
		{"deleteitem", "DeleteItem", 10, permFailLow},
		{"clrpassword", "StorageClearPassword", 10, permFailLow},
		{"refineweapon", "RefineWeapon", 10, permFailLow},
		// Delphi 实现内 <6 硬编码静默（ObjBase.pas:10795）。
		{"changeitemname", "ChangeItemName", 6, permFailSilent},
	}
	for _, c := range cases {
		def := gmCmdTable[c.literal]
		if def == nil {
			t.Errorf("%s: not registered", c.literal)
			continue
		}
		if def.key != c.key || def.resolvedPerm != c.perm || def.permFail != c.fail {
			t.Errorf("%s: key=%s perm=%d fail=%d, want key=%s perm=%d fail=%d",
				c.literal, def.key, def.resolvedPerm, def.permFail, c.key, c.perm, c.fail)
		}
	}
}

// TestItemFamilyDispatch 权限拒绝与 "?" 帮助均被消费（不静默）；
// ChangeItemName 低权限静默消费（permFailSilent）。
func TestItemFamilyDispatch(t *testing.T) {
	buildGMCommandTable(nil)
	p, server := newTestPlayer()
	for _, cmd := range []string{
		"supermake 1 0 5", "deleteitem a b", "clrpassword a",
		"refineweapon 1 1 1 1", "changeitemname 1 1 x",
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: permission-denied should be consumed", cmd)
		}
	}
	p.Permission = 10
	for _, cmd := range []string{
		"supermake ?", "deleteitem ?", "clrpassword ?", "refineweapon ?", "changeitemname ?",
	} {
		if !p.HandleGMCommand(cmd, server) {
			t.Errorf("%q: usage help should be consumed", cmd)
		}
	}
}

// TestCmdSuperMake 穿戴物品属性调整（CmdSmakeItem，ObjBase.pas:14899-14952）。
func TestCmdSuperMake(t *testing.T) {
	buildGMCommandTable(nil)
	p, _ := newItemFamilyFixture(t)
	server := &netserver.TCPServer{}

	p.UseItems[protocol.UWeapon] = &protocol.UserItem{
		MakeIndex: 100, WIndex: 1, Dura: 5000, DuraMax: 8000,
	}

	// 属性位 0-13 → BtValue 直映射。
	if !p.cmdSuperMake(server, []string{"1", "0", "7"}) {
		t.Fatal("supermake should be consumed")
	}
	if p.UseItems[protocol.UWeapon].BtValue[0] != 7 {
		t.Errorf("BtValue[0] = %d, want 7", p.UseItems[protocol.UWeapon].BtValue[0])
	}
	if p.UseItems[protocol.UWeapon].BtValue[5] != 0 {
		t.Error("supermake pos 0 must not touch other slots")
	}

	// 属性位 14=持久：上限 65 → 65×1000。
	p.cmdSuperMake(server, []string{"1", "14", "70"})
	if p.UseItems[protocol.UWeapon].Dura != 65000 {
		t.Errorf("Dura = %d, want 65000 (min(65,70)*1000)", p.UseItems[protocol.UWeapon].Dura)
	}
	// 属性位 15=持久上限。
	p.cmdSuperMake(server, []string{"1", "15", "20"})
	if p.UseItems[protocol.UWeapon].DuraMax != 20000 {
		t.Errorf("DuraMax = %d, want 20000", p.UseItems[protocol.UWeapon].DuraMax)
	}

	// 空槽位：'身上没指定物品！！！' 分支，正常消费。
	if !p.cmdSuperMake(server, []string{"3", "0", "5"}) {
		t.Error("empty slot should be consumed")
	}

	// 越界参数：静默退出，不改动。
	before := p.UseItems[protocol.UWeapon].BtValue[0]
	for _, args := range [][]string{
		{"13", "0", "5"}, // 部位越界
		{"1", "16", "5"}, // 属性位越界
		{"1", "0", "256"}, // 值越界
		{"-1", "0", "5"},
	} {
		p.cmdSuperMake(server, args)
	}
	if p.UseItems[protocol.UWeapon].BtValue[0] != before {
		t.Error("out-of-range params must not modify item")
	}

	// wIndex 有效但物品库无定义：静默退出（Delphi GetStdItem=nil）。
	p.UseItems[protocol.URingL] = &protocol.UserItem{MakeIndex: 101, WIndex: 99, Dura: 1000, DuraMax: 1000}
	if !p.cmdSuperMake(server, []string{"7", "0", "9"}) {
		t.Error("unknown wIndex should be consumed")
	}
	if p.UseItems[protocol.URingL].BtValue[0] != 0 {
		t.Error("unknown wIndex must not be modified")
	}

	// 缺参：静默消费。
	if !p.cmdSuperMake(server, []string{"1", "0"}) {
		t.Error("missing params should be consumed")
	}
}

// TestCmdDeleteItem 按名称删除目标背包物品（CmdDeleteItem，ObjBase.pas:12255-12295）。
func TestCmdDeleteItem(t *testing.T) {
	buildGMCommandTable(nil)
	p, engine := newItemFamilyFixture(t)
	server := &netserver.TCPServer{}

	victim := newItemFamilyVictim(t, "victim1", 2, engine, p.ItemDB)
	for mi := int32(201); mi <= 203; mi++ {
		victim.ItemList = append(victim.ItemList, &protocol.UserItem{MakeIndex: mi, WIndex: 2}) // 金创药
	}
	victim.ItemList = append(victim.ItemList, &protocol.UserItem{MakeIndex: 204, WIndex: 1}) // 木剑
	victim.ItemList = append(victim.ItemList, &protocol.UserItem{MakeIndex: 205, WIndex: 3}) // IronSword

	// 删 2 件金创药（共 3 件 → 剩 1 件）。
	if !p.cmdDeleteItem(server, []string{"victim1", "金创药", "2"}) {
		t.Fatal("deleteitem should be consumed")
	}
	potions := 0
	for _, it := range victim.ItemList {
		if it.WIndex == 2 {
			potions++
		}
	}
	if potions != 1 || len(victim.ItemList) != 3 {
		t.Errorf("after delete: potions=%d total=%d, want 1/3", potions, len(victim.ItemList))
	}

	// 名称大小写不敏感（Delphi CompareText）。
	p.cmdDeleteItem(server, []string{"victim1", "ironsword", "1"})
	for _, it := range victim.ItemList {
		if it.WIndex == 3 {
			t.Error("case-insensitive delete should remove IronSword")
		}
	}

	// 数量超过持有量：全部删完不越界。
	p.cmdDeleteItem(server, []string{"victim1", "金创药", "99"})
	if len(victim.ItemList) != 1 { // 只剩木剑
		t.Errorf("bag size = %d, want 1", len(victim.ItemList))
	}

	// 不存在的物品名：无变化。
	p.cmdDeleteItem(server, []string{"victim1", "不存在", "1"})
	if len(victim.ItemList) != 1 {
		t.Error("non-matching name must not delete")
	}

	// 目标不在线：文案分支，正常消费。
	if !p.cmdDeleteItem(server, []string{"nobody", "木剑", "1"}) {
		t.Error("offline target should be consumed")
	}
	// 缺参：格式提示分支。
	if !p.cmdDeleteItem(server, []string{"victim1"}) {
		t.Error("missing params should be consumed")
	}
}

// TestCmdRefineWeapon 武器属性调整（CmdRefineWeapon，ObjBase.pas:14106-14124）。
func TestCmdRefineWeapon(t *testing.T) {
	buildGMCommandTable(nil)
	p, _ := newItemFamilyFixture(t)
	server := &netserver.TCPServer{}

	// 无武器：静默退出。
	if !p.cmdRefineWeapon(server, []string{"1", "1", "1", "1"}) {
		t.Error("no weapon should be consumed")
	}

	p.UseItems[protocol.UWeapon] = &protocol.UserItem{
		MakeIndex: 300, WIndex: 1, Dura: 8000, DuraMax: 8000,
	}
	p.UseItems[protocol.UWeapon].BtValue[0] = 99 // 预置脏值，验证被覆盖

	if !p.cmdRefineWeapon(server, []string{"3", "2", "1", "5"}) {
		t.Fatal("refineweapon should be consumed")
	}
	w := p.UseItems[protocol.UWeapon]
	if w.BtValue[0] != 3 || w.BtValue[1] != 2 || w.BtValue[2] != 1 || w.BtValue[5] != 5 {
		t.Errorf("BtValue = %v, want [0]=3 [1]=2 [2]=1 [5]=5", w.BtValue)
	}

	// DC+MC+SC > 10：静默退出，不改动。
	p.cmdRefineWeapon(server, []string{"5", "5", "1", "9"})
	if w.BtValue[0] != 3 || w.BtValue[5] != 5 {
		t.Error("sum>10 must not modify weapon")
	}

	// 缺参按 0（Delphi Str_ToInt 默认 0）：清零。
	p.cmdRefineWeapon(server, []string{"1"})
	if w.BtValue[0] != 1 || w.BtValue[1] != 0 || w.BtValue[2] != 0 || w.BtValue[5] != 0 {
		t.Errorf("missing params should default to 0, got %v", w.BtValue)
	}
}

// TestCmdChangeItemName 自定义物品名（CmdChangeItemName，ObjBase.pas:10791-10817）。
func TestCmdChangeItemName(t *testing.T) {
	buildGMCommandTable(nil)
	resetCustomItemNames()
	p, _ := newItemFamilyFixture(t)
	server := &netserver.TCPServer{}

	// 首次设置成功。
	if !p.cmdChangeItemName(server, []string{"100", "1", "屠龙"}) {
		t.Fatal("changeitemname should be consumed")
	}
	gmCustomItemNameMu.Lock()
	name := gmCustomItemNames[100]
	gmCustomItemNameMu.Unlock()
	if name != "屠龙" {
		t.Errorf("custom name = %q, want 屠龙", name)
	}

	// 重复设置：'此物品，已经设置了其它的名称！！！' 分支，不覆盖。
	p.cmdChangeItemName(server, []string{"100", "1", "倚天"})
	gmCustomItemNameMu.Lock()
	name = gmCustomItemNames[100]
	gmCustomItemNameMu.Unlock()
	if name != "屠龙" {
		t.Errorf("duplicate set must not overwrite, got %q", name)
	}

	// 非法编号：格式提示分支，不登记。
	p.cmdChangeItemName(server, []string{"0", "1", "裁决"})  // makeIndex<=0
	p.cmdChangeItemName(server, []string{"5", "-1", "裁决"}) // itemIndex<0
	p.cmdChangeItemName(server, []string{"x", "1", "裁决"})  // 解析失败
	gmCustomItemNameMu.Lock()
	_, exists5 := gmCustomItemNames[5]
	_, exists0 := gmCustomItemNames[0]
	gmCustomItemNameMu.Unlock()
	if exists5 || exists0 {
		t.Error("invalid params must not register")
	}

	// 缺参：格式提示分支。
	if !p.cmdChangeItemName(server, []string{"100", "1"}) {
		t.Error("missing params should be consumed")
	}
}

// TestCmdClrPassword 清除仓库密码（CmdClearHumanPassword，ObjBase.pas:11088-11110）。
func TestCmdClrPassword(t *testing.T) {
	buildGMCommandTable(nil)
	p, engine := newItemFamilyFixture(t)
	server := &netserver.TCPServer{}

	victim := newItemFamilyVictim(t, "victim1", 2, engine, p.ItemDB)
	victim.StoragePassword = "1234"
	victim.StoragePwdLocked = true
	victim.storagePwdFail = 2

	if !p.cmdClrPassword(server, []string{"victim1"}) {
		t.Fatal("clrpassword should be consumed")
	}
	if victim.StoragePassword != "" || victim.StoragePwdLocked || victim.storagePwdFail != 0 {
		t.Errorf("password state not cleared: pwd=%q locked=%v fail=%d",
			victim.StoragePassword, victim.StoragePwdLocked, victim.storagePwdFail)
	}

	// 目标不在线：静默退出（Delphi GetPlayObject=nil → Exit）。
	if !p.cmdClrPassword(server, []string{"nobody"}) {
		t.Error("offline target should be consumed")
	}
	// 缺参：双行提示分支。
	if !p.cmdClrPassword(server, nil) {
		t.Error("missing params should be consumed")
	}
}
