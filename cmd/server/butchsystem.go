package main

import (
	"math/rand"
	"time"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// initMeatQuality 刷怪时初始化肉质/坚韧度（Delphi UsrEngn.pas:1835-1864）：
// 鸡 MQ=rand(3500)+3000、BL=50；鹿 1/30 稀有 rand(20000)+10000（TChickenDeer），
// 否则 rand(8000)+8000、BL=150；狼/猎鹰 rand(8000)+8000、BL=150；
// 蝎子（race 84）保持 0——BL=0 第一刀即出货（TScorpion 忠实原版）。
func initMeatQuality(mon *MonsterObject) {
	if !mon.Animal {
		return
	}
	switch mon.Race {
	case 51: // ANIMAL_CHICKEN（UsrEngn.pas:1835-1841）
		mon.MeatQuality = rand.Intn(3500) + 3000
		mon.BodyLeathery = 50
	case 52: // ANIMAL_DEER（UsrEngn.pas:1842-1857）
		if rand.Intn(30) == 0 {
			mon.MeatQuality = rand.Intn(20000) + 10000
		} else {
			mon.MeatQuality = rand.Intn(8000) + 8000
		}
		mon.BodyLeathery = 150
	case 53: // ANIMAL_WOLF（UsrEngn.pas:1858-1864）
		mon.MeatQuality = rand.Intn(8000) + 8000
		mon.BodyLeathery = 150
	}
}

// rollCorpseItems 刷怪时把掉落表 roll 进尸体背包（Delphi MonGetRandomItems，
// UsrEngn.pas:1536-1605，AddBaseObject 内调用 :1951）。动物死亡不散落，
// 物品留尸体由 CM_BUTCH 挖出（ObjBase.pas:20982-20998 not m_boAnimal）。
// 金币条目跳过：Delphi 金币 roll 进 m_nGold（:1563-1565），但动物死亡
// 不 ScatterGolds、挖肉只取 m_ItemList，玩家永远拿不到，本就是死数据。
func (e *UserEngine) rollCorpseItems(mon *MonsterObject) {
	if e.DropTables == nil || e.ItemDB == nil {
		return
	}
	for _, entry := range e.DropTables.GetDrops(mon.Name) {
		if entry.Chance <= 0 || rand.Intn(entry.Chance) != 0 {
			continue
		}
		if entry.Gold > 0 {
			continue
		}
		_, ui := createDropItem(e.ItemDB, entry.ItemName, &e.nextItemID, e.Config, false)
		if ui != nil {
			mon.CorpseItems = append(mon.CorpseItems, ui)
		}
	}
}

// applyMeatQuality 把玩家背包内所有肉（StdMode 40）的耐久设为给定品质。
// 注意：Delphi 字面实现（ApplyMeatQuality，ObjBase.pas:20517-20535）在
// ClientGetButchItem 中无修饰符调用 = Self（玩家），会把背包肉的耐久设为
// 玩家自己的 m_nMeatQuality（恒 0），属原版源码缺陷。Go 此处为有意修正：
// 用尸体的 MeatQuality 覆写，使"肉质"玩法成立。
func (p *PlayObject) applyMeatQuality(quality int) {
	if p.ItemDB == nil {
		return
	}
	for _, item := range p.ItemList {
		if item == nil {
			continue
		}
		if def := p.ItemDB.GetByIdx(int(item.WIndex)); def != nil && def.StdMode == 40 {
			item.Dura = uint16(quality)
		}
	}
}

// takeCorpseItems 逐件把尸体物品移入玩家背包（Delphi TakeBagItems，
// ObjBase.pas:20537-20556）：背包满即停、剩余留在尸体（下次再挖或换人挖）；
// 至少拿到 1 件返回 true。
func (p *PlayObject) takeCorpseItems(server *netserver.TCPServer, target *MonsterObject) bool {
	got := false
	for len(target.CorpseItems) > 0 {
		if p.Engine == nil || len(p.ItemList) >= p.Engine.Config.GetMaxBagSlots() {
			break
		}
		p.ItemList = append(p.ItemList, target.CorpseItems[0])
		target.CorpseItems = target.CorpseItems[1:]
		got = true
	}
	if got {
		p.SendBagItemsFull(server)
	}
	return got
}

// broadcastButch 广播挖肉动作（Delphi RM_BUTCH→SM_BUTCH=637，
// ObjBase.pas:17503；附近玩家看到屠宰者下蹲，ClMain.pas:4059-4071）。
func (p *PlayObject) broadcastButch(server *netserver.TCPServer) {
	butchMsg := protocol.MakeDefaultMsg(protocol.SMButch, p.ID, uint16(p.CurrX), uint16(p.CurrY), uint16(p.Dir))
	objs := p.envir.GetRangeObjects(p.CurrX, p.CurrY, viewRange)
	for _, obj := range objs {
		if other, ok := obj.(*PlayObject); ok && !other.Ghost {
			server.Send(other.Session.ID, butchMsg, "")
		}
	}
}

// HandleButch 处理屠宰动物尸体（CMButch=1007, Recog=目标怪物ID）。
// Delphi ClientGetButchItem（ObjBase.pas:17458-17504）：
// 频限 dwTurnIntervalTime=600ms；距离≤2 格；条件 Death && !Skeleton && Animal；
// 每次挖扣坚韧度 rand(16)+5、肉质 rand(201)+100；坚韧度耗尽时
// race∈[50,80) 尸体骨架化（RM_SKELETON）+ 尸体物品入背包 + 坚韧度重置 50；
// RM_BUTCH 广播在频限之后无条件执行（距离外也播下蹲动作）。
func (p *PlayObject) HandleButch(msg SendMessage, server *netserver.TCPServer) {
	if p.envir == nil || p.Death {
		return
	}
	targetID := msg.SourceID
	if targetID == 0 {
		targetID = int32(msg.Param1)
	}
	target, ok := p.envir.getObjectByID(targetID).(*MonsterObject)
	if !ok || target.Ghost || !target.Death {
		return
	}

	// 频限：超限静默丢弃（ObjBase.pas:17465-17471）。不复用
	// checkActionSpeed/TurnTick——它带超速计数与 FAIL 回包，
	// Delphi 屠宰频限只是静默跳过。
	now := time.Now().UnixMilli()
	interval := int64(600)
	if p.Engine != nil && p.Engine.Config != nil {
		interval = p.Engine.Config.GetButchInterval()
	}
	if now-p.butchTick < interval {
		return
	}
	p.butchTick = now

	// RM_BUTCH 在距离判断之外无条件广播（ObjBase.pas:17503）。
	defer p.broadcastButch(server)

	// 距离 ≤2（ObjBase.pas:17472），用目标实际坐标（等价 IsValidObject）。
	if abs(p.CurrX-target.CurrX) > 2 || abs(p.CurrY-target.CurrY) > 2 {
		return
	}

	if !target.Skeleton && target.Animal { // ObjBase.pas:17476
		target.BodyLeathery -= rand.Intn(16) + 5   // :17478
		target.MeatQuality -= rand.Intn(201) + 100 // :17479
		if target.MeatQuality < 0 {
			target.MeatQuality = 0 // :17482
		}
		if target.BodyLeathery <= 0 { // :17483
			if target.Race >= 50 && target.Race < 80 { // RC_ANIMAL..RC_MONSTER，:17485
				target.Skeleton = true
				target.envir.broadcastSkeletonMsg(target.BaseObject, target.ID, target.CurrX, target.CurrY, target.Dir)
			}
			// 先给物后对齐肉质，使新挖到的肉也带上品质。
			got := p.takeCorpseItems(server, target)
			if target.MeatQuality > 0 {
				p.applyMeatQuality(target.MeatQuality)
			}
			if !got {
				p.sysMsg(server, "你没有找到任何物品") // sYouFoundNothing，:17493
			}
			target.BodyLeathery = 50 // :17495
		}
	}

	log.Logf(log.LevelDebug, "Butch", "%s butchered %s (id=%d), corpse items left=%d", p.Name, target.Name, target.ID, len(target.CorpseItems))
}
