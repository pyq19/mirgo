package main

import (
	"testing"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// 批1回归测试：尸体可见性（Bug 2）与杀鸡挖肉机制（Bug 3）。
// 对齐 Delphi ObjBase.pas:25697-25726（视野内尸体通告）、
// 20510-20515（MakeGhost）、17458-17504（ClientGetButchItem）、
// UsrEngn.pas:1536-1605（MonGetRandomItems）。

func placeNamedMonster(env *Environment, name string, id int32, race byte, x, y int) *MonsterObject {
	mon := NewMonsterObject(name, id, race, 11, 160, 100, 1400, 2000, 0)
	mon.MapName = env.Name
	mon.CurrX, mon.CurrY = x, y
	mon.envir = env
	env.AddObject(x, y, OS_MOVINGOBJECT, mon)
	return mon
}

func newButchPlayer(env *Environment, id int32, x, y int) *PlayObject {
	p := &PlayObject{BaseObject: NewBaseObject("tester", id)}
	p.Session = &netserver.Session{ID: int64(id)}
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	p.envir = env
	p.CurrX, p.CurrY = x, y
	p.ViewRange = 12
	p.VisibleActors = map[int32]*VisibleEntry{}
	env.AddObject(x, y, OS_MOVINGOBJECT, p)
	return p
}

// TestSearchViewRangeKeepsDeadMonster 回归（Bug 2 根因）：死亡怪物不再被
// 视野扫描排除，后进入视野的玩家能看到尸体。
func TestSearchViewRangeKeepsDeadMonster(t *testing.T) {
	env := newTestEnv(30, 30)
	p := newButchPlayer(env, 1, 15, 15)
	mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
	mon.Death = true
	mon.DeathTick = time.Now().UnixMilli()

	p.SearchViewRange(&netserver.TCPServer{})

	if _, ok := p.VisibleActors[mon.ID]; !ok {
		t.Fatal("死亡怪物应保留在 VisibleActors（旧实现 2~4 秒后发 SMDisappear 删尸体）")
	}

	// 死亡玩家同样保留。
	p2 := newButchPlayer(env, 3, 17, 15)
	p2.Death = true
	p.SearchViewRange(&netserver.TCPServer{})
	if _, ok := p.VisibleActors[p2.ID]; !ok {
		t.Fatal("死亡玩家应保留在 VisibleActors")
	}

	// ghost 后不再可见。
	mon.Ghost = true
	p.SearchViewRange(&netserver.TCPServer{})
	if _, ok := p.VisibleActors[mon.ID]; ok {
		t.Fatal("ghost 后的怪物不应再可见")
	}
}

// TestNewVisibleEntryMsg 新入视野对象的消息选择（Delphi ObjBase.pas:25697-25726）：
// 尸体→SM_DEATH、骨架→SM_SKELETON、活体→SM_TURN、死亡玩家→SM_DEATH。
func TestNewVisibleEntryMsg(t *testing.T) {
	env := newTestEnv(30, 30)
	p := newButchPlayer(env, 1, 15, 15)
	mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)

	mon.Death = true
	resp, body := p.newVisibleEntryMsg(mon)
	if resp.Ident != protocol.SMDeath {
		t.Errorf("死亡怪物新入视野 ident = %d, want SMDeath(%d)", resp.Ident, protocol.SMDeath)
	}
	if len(body) != len(protocol.EncodeBuffer(make([]byte, 8))) {
		t.Errorf("尸体消息 body 应为 8 字节 TCharDesc（不带名字），实际 %d 字节", len(body))
	}

	mon.Skeleton = true
	resp, _ = p.newVisibleEntryMsg(mon)
	if resp.Ident != protocol.SMSkeleton {
		t.Errorf("骨架新入视野 ident = %d, want SMSkeleton(%d)", resp.Ident, protocol.SMSkeleton)
	}

	mon.Death, mon.Skeleton = false, false
	resp, _ = p.newVisibleEntryMsg(mon)
	if resp.Ident != protocol.SMTurn {
		t.Errorf("活体新入视野 ident = %d, want SMTurn(%d)", resp.Ident, protocol.SMTurn)
	}

	p2 := newButchPlayer(env, 3, 17, 15)
	p2.Death = true
	resp, _ = p.newVisibleEntryMsg(p2)
	if resp.Ident != protocol.SMDeath {
		t.Errorf("死亡玩家新入视野 ident = %d, want SMDeath(%d)", resp.Ident, protocol.SMDeath)
	}
}

// TestMonsterGhostBroadcastsDisappear 尸体到 3 分钟 ghost 时广播 RM_DISAPPEAR
//（Delphi MakeGhost→DisappearA，ObjBase.pas:20510-20515, 20144-20148）。
func TestMonsterGhostBroadcastsDisappear(t *testing.T) {
	env := newTestEnv(30, 30)
	p := newButchPlayer(env, 1, 15, 15)
	mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
	mon.Animal = true
	mon.Death = true
	mon.LootDropped = true
	mon.DeathTick = time.Now().UnixMilli() - 200000 // 超过 corpseDelay 默认 180s

	e := NewUserEngine(nil, NewMapManager(""))
	e.Config = &ServerConfig{}
	e.Monsters = []*MonsterObject{mon}

	e.ProcessMonsters(&netserver.TCPServer{}, time.Now().UnixMilli())

	if !mon.Ghost {
		t.Fatal("超过 corpseDelay 的尸体应 ghost")
	}
	found := false
	for _, m := range p.msgList {
		if m.Ident == RM_DISAPPEAR && m.SourceID == mon.ID {
			found = true
		}
	}
	if !found {
		t.Fatal("ghost 时应向附近玩家广播 RM_DISAPPEAR")
	}
}

// TestAnimalLootStaysInCorpse 动物死亡不散落物品（Delphi ObjBase.pas:20983
// not m_boAnimal 条件），非动物照常掉落。
func TestAnimalLootStaysInCorpse(t *testing.T) {
	env := newTestEnv(30, 30)
	itemDB := newTestItemDB(
		ItemDef{Idx: 57, Name: "鸡肉", StdMode: 40, DuraMax: 4000},
		ItemDef{Idx: 1, Name: "金创药(小量)", StdMode: 0, DuraMax: 100},
	)
	dt := &DropTable{tables: map[string][]DropEntry{
		"鸡":    {{ItemName: "鸡肉", Chance: 1, Count: 1}},
		"稻草人": {{ItemName: "金创药(小量)", Chance: 1, Count: 1}},
	}}
	e := NewUserEngine(nil, NewMapManager(""))
	e.Config = &ServerConfig{}
	e.ItemDB = itemDB
	e.DropTables = dt

	chicken := placeNamedMonster(env, "鸡", 2, 51, 15, 15)
	chicken.Animal = true
	chicken.Death = true
	scarecrow := placeNamedMonster(env, "稻草人", 3, 83, 20, 20)
	scarecrow.Death = true
	e.Monsters = []*MonsterObject{chicken, scarecrow}

	e.ProcessMonsters(&netserver.TCPServer{}, time.Now().UnixMilli())

	if !chicken.LootDropped || !scarecrow.LootDropped {
		t.Fatal("死亡处理后 LootDropped 应置位")
	}
	for _, gi := range env.GroundItems {
		if gi.Name == "鸡肉" {
			t.Fatal("鸡肉不应散落地面（应留尸体等挖肉）")
		}
	}
	foundItem := false
	for _, gi := range env.GroundItems {
		if gi.Name == "金创药(小量)" {
			foundItem = true
		}
	}
	if !foundItem {
		t.Fatal("非动物（稻草人）掉落应照常散落地面")
	}
}

// TestRollCorpseItems 刷怪时 roll 尸体背包（Delphi MonGetRandomItems，
// UsrEngn.pas:1536-1605）：鸡 1/1 恰 1 件、耐久为 DuraMax 20%~100%、
// 金币条目不入尸体、肉不做落地-2000 惩罚。
func TestRollCorpseItems(t *testing.T) {
	itemDB := newTestItemDB(
		ItemDef{Idx: 57, Name: "鸡肉", StdMode: 40, DuraMax: 4000},
		ItemDef{Idx: 3, Name: "肉", StdMode: 40, DuraMax: 4000},
	)
	e := NewUserEngine(nil, NewMapManager(""))
	e.Config = &ServerConfig{}
	e.ItemDB = itemDB
	e.DropTables = &DropTable{tables: map[string][]DropEntry{
		"鸡": {{ItemName: "鸡肉", Chance: 1, Count: 1}},
		"鹿": {{ItemName: "肉", Chance: 1, Count: 1}, {ItemName: "肉", Chance: 10, Count: 1}},
		"猎鹰": {{ItemName: "金币", Chance: 4, Count: 70, Gold: 70}},
	}}

	chicken := NewMonsterObject("鸡", 1, 51, 11, 160, 5, 1400, 2000, 9)
	chicken.Animal = true
	e.rollCorpseItems(chicken)
	if len(chicken.CorpseItems) != 1 {
		t.Fatalf("鸡尸体应有 1 件物品，实际 %d", len(chicken.CorpseItems))
	}
	meat := chicken.CorpseItems[0]
	if meat.WIndex != 57 {
		t.Errorf("尸体物品 WIndex = %d, want 57（鸡肉）", meat.WIndex)
	}
	// Dura = DuraMax/100*(20+rand(80))，无落地-2000 惩罚 → [800, 4000]
	if int(meat.Dura) < 800 || int(meat.Dura) > 4000 {
		t.Errorf("尸体肉耐久 = %d, want [800,4000]（不应有落地惩罚）", meat.Dura)
	}

	deer := NewMonsterObject("鹿", 2, 52, 11, 161, 25, 1400, 2000, 18)
	deer.Animal = true
	e.rollCorpseItems(deer)
	if len(deer.CorpseItems) < 1 || len(deer.CorpseItems) > 2 {
		t.Fatalf("鹿尸体应有 1~2 件肉，实际 %d", len(deer.CorpseItems))
	}

	falcon := NewMonsterObject("猎鹰", 3, 53, 11, 161, 25, 1400, 2000, 18)
	falcon.Animal = true
	for i := 0; i < 20; i++ { // 金币 1/4，跑 20 次必然命中过
		falcon.CorpseItems = nil
		e.rollCorpseItems(falcon)
		if len(falcon.CorpseItems) != 0 {
			t.Fatal("金币条目不应进入尸体背包（Delphi 动物金币玩家永远拿不到）")
		}
	}
}

// newButchSetup 搭一个玩家+死鸡的屠宰场景。
func newButchSetup(t *testing.T) (*PlayObject, *MonsterObject) {
	t.Helper()
	env := newTestEnv(30, 30)
	p := newButchPlayer(env, 1, 15, 15)
	p.ItemDB = newTestItemDB(ItemDef{Idx: 57, Name: "鸡肉", StdMode: 40, DuraMax: 4000})
	mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
	mon.Animal = true
	mon.Death = true
	mon.MeatQuality = 5000
	mon.BodyLeathery = 50
	mon.CorpseItems = []*protocol.UserItem{
		{MakeIndex: 99, WIndex: 57, Dura: 3000, DuraMax: 4000},
	}
	return p, mon
}

// TestHandleButchLifecycle 屠宰全流程（Delphi ClientGetButchItem，
// ObjBase.pas:17458-17504）：频限/距离/皮量减量/骨架化/背包满/皮量重置。
func TestHandleButchLifecycle(t *testing.T) {
	server := &netserver.TCPServer{}
	msg := SendMessage{Ident: int(protocol.CMButch), SourceID: 2}

	t.Run("每次挖扣皮量5~20、肉质100~300", func(t *testing.T) {
		p, mon := newButchSetup(t)
		p.HandleButch(msg, server)
		blLoss := 50 - mon.BodyLeathery
		if blLoss < 5 || blLoss > 20 {
			t.Errorf("皮量减量 = %d, want [5,20]", blLoss)
		}
		mqLoss := 5000 - mon.MeatQuality
		if mqLoss < 100 || mqLoss > 300 {
			t.Errorf("肉质减量 = %d, want [100,300]", mqLoss)
		}
		if len(mon.CorpseItems) != 1 || len(p.ItemList) != 0 {
			t.Fatal("皮量未耗尽时不应出货")
		}
	})

	t.Run("频限内第二次无效", func(t *testing.T) {
		p, mon := newButchSetup(t)
		p.HandleButch(msg, server)
		bl := mon.BodyLeathery
		p.HandleButch(msg, server)
		if mon.BodyLeathery != bl {
			t.Fatal("600ms 频限内第二次挖肉不应生效")
		}
	})

	t.Run("距离超过2格无效", func(t *testing.T) {
		p, mon := newButchSetup(t)
		p.CurrX, p.CurrY = 20, 20 // 距尸体 4 格
		p.HandleButch(msg, server)
		if mon.BodyLeathery != 50 {
			t.Fatal("距离超过 2 格不应扣皮量")
		}
	})

	t.Run("皮量耗尽出货+骨架化+重置", func(t *testing.T) {
		p, mon := newButchSetup(t)
		mon.BodyLeathery = 5 // 一刀必耗尽
		p.HandleButch(msg, server)
		if !mon.Skeleton {
			t.Fatal("race 51 皮量耗尽应骨架化")
		}
		if mon.BodyLeathery != 50 {
			t.Errorf("出货后皮量应重置 50，实际 %d", mon.BodyLeathery)
		}
		if len(mon.CorpseItems) != 0 {
			t.Fatal("尸体物品应已移出")
		}
		if len(p.ItemList) != 1 || p.ItemList[0].WIndex != 57 {
			t.Fatal("鸡肉应入背包")
		}
		// 肉质联动：背包肉的耐久 = 尸体 MeatQuality
		if int(p.ItemList[0].Dura) != mon.MeatQuality {
			t.Errorf("背包肉耐久 = %d, want 尸体肉质 %d", p.ItemList[0].Dura, mon.MeatQuality)
		}
		// 附近玩家（含自己）应收到 RM_SKELETON
		found := false
		for _, m := range p.msgList {
			if m.Ident == RM_SKELETON && m.SourceID == mon.ID {
				found = true
			}
		}
		if !found {
			t.Fatal("骨架化应广播 RM_SKELETON")
		}
		// 骨架后再挖不再扣皮量
		p.butchTick = 0
		bl := mon.BodyLeathery
		p.HandleButch(msg, server)
		if mon.BodyLeathery != bl {
			t.Fatal("骨架尸体再挖不应扣皮量")
		}
	})

	t.Run("背包满物品留尸体", func(t *testing.T) {
		p, mon := newButchSetup(t)
		mon.BodyLeathery = 5
		p.Engine.Config.Player.MaxBagSlots = 2
		p.ItemList = []*protocol.UserItem{{WIndex: 1}, {WIndex: 1}}
		p.HandleButch(msg, server)
		if len(mon.CorpseItems) != 1 {
			t.Fatal("背包满时物品应留在尸体")
		}
		if len(p.ItemList) != 2 {
			t.Fatal("背包满时不应塞入物品")
		}
		if mon.BodyLeathery != 50 {
			t.Errorf("Delphi 语义：无论是否拿到物品皮量都重置，实际 %d", mon.BodyLeathery)
		}
	})

	t.Run("非动物无法屠宰", func(t *testing.T) {
		env := newTestEnv(30, 30)
		p := newButchPlayer(env, 1, 15, 15)
		mon := placeNamedMonster(env, "稻草人", 2, 83, 16, 15)
		mon.Death = true
		mon.BodyLeathery = 50
		p.HandleButch(SendMessage{Ident: int(protocol.CMButch), SourceID: 2}, server)
		if mon.BodyLeathery != 50 {
			t.Fatal("非动物不应扣皮量")
		}
	})
}

// TestMeatQualityDegradation 肉质下降三路径（Delphi ObjBase.pas:2807-2811
// 物理 -rand(300)、:4517 魔法 -伤害×1000、:4260 绿毒每跳 -1000）。
func TestMeatQualityDegradation(t *testing.T) {
	server := &netserver.TCPServer{}

	t.Run("物理攻击扣rand(300)", func(t *testing.T) {
		env := newTestEnv(30, 30)
		p := newButchPlayer(env, 1, 15, 15)
		mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
		mon.Animal = true
		mon.MeatQuality = 5000
		mon.WAbil.HP = 100
		p.applyDamage(server, mon.BaseObject, 10, 0, false)
		loss := 5000 - mon.MeatQuality
		if loss < 0 || loss >= 300 {
			t.Errorf("物理肉质减量 = %d, want [0,300)", loss)
		}
	})

	t.Run("魔法攻击扣伤害x1000", func(t *testing.T) {
		env := newTestEnv(30, 30)
		p := newButchPlayer(env, 1, 15, 15)
		mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
		mon.Animal = true
		mon.MeatQuality = 5000
		mon.WAbil.HP = 100
		p.applyDamage(server, mon.BaseObject, 3, 0, true)
		if mon.MeatQuality != 2000 {
			t.Errorf("魔法后肉质 = %d, want 2000（5000-3*1000）", mon.MeatQuality)
		}
	})

	t.Run("肉质不为负", func(t *testing.T) {
		env := newTestEnv(30, 30)
		p := newButchPlayer(env, 1, 15, 15)
		mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
		mon.Animal = true
		mon.MeatQuality = 500
		mon.WAbil.HP = 10000
		p.applyDamage(server, mon.BaseObject, 3, 0, true)
		if mon.MeatQuality != 0 {
			t.Errorf("肉质应 clamp 到 0，实际 %d", mon.MeatQuality)
		}
	})

	t.Run("非动物不扣", func(t *testing.T) {
		env := newTestEnv(30, 30)
		p := newButchPlayer(env, 1, 15, 15)
		mon := placeNamedMonster(env, "稻草人", 2, 83, 16, 15)
		mon.MeatQuality = 5000
		mon.WAbil.HP = 100
		p.applyDamage(server, mon.BaseObject, 10, 0, false)
		if mon.MeatQuality != 5000 {
			t.Fatal("非动物不应扣肉质")
		}
	})

	t.Run("绿毒每跳扣1000", func(t *testing.T) {
		env := newTestEnv(30, 30)
		mon := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
		mon.Animal = true
		mon.MeatQuality = 5000
		mon.WAbil.HP = 100
		mon.MaxHP = 100
		mon.StatusTimeArr[POISON_DECHEALTH] = 5
		mon.StatusTimeArr[POISON_STONE] = 5 // 石化让 Run 在毒 tick 后提前返回
		mon.Run(server, time.Now().UnixMilli(), nil)
		if mon.MeatQuality != 4000 {
			t.Errorf("绿毒一跳后肉质 = %d, want 4000", mon.MeatQuality)
		}
	})

	t.Run("怪物互殴物理扣肉质", func(t *testing.T) {
		env := newTestEnv(30, 30)
		target := placeNamedMonster(env, "鸡", 2, 51, 16, 15)
		target.Animal = true
		target.MeatQuality = 5000
		target.WAbil.HP = 10000
		target.SpeedPoint = 1
		attacker := placeNamedMonster(env, "狼", 3, 53, 15, 15)
		attacker.HitPoint = 100
		attacker.WAbil.DC = 10 | 10<<16
		attacker.attackMonsterTarget(server, target, time.Now().UnixMilli())
		loss := 5000 - target.MeatQuality
		if loss < 0 || loss >= 300 {
			t.Errorf("怪物互殴肉质减量 = %d, want [0,300)", loss)
		}
	})
}
