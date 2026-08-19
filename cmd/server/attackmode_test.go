package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// TestDropItemConsecutive 回归："站在同一个位置只能丢一次物品"。
// Delphi 的 m_DealLastTick 只由交易操作写入，ClientDropItem 的 3 秒
// 检查仅禁止交易后立刻丢弃（ObjBase.pas:16244），连续丢弃不节流。
// 旧 Go 实现每次丢弃都刷新 lastDealTick，导致 3 秒内只能丢一次。
func TestDropItemConsecutive(t *testing.T) {
	env := newTestEnv(40, 40)
	p := &PlayObject{BaseObject: NewBaseObject("dropper", 1)}
	p.Session = &netserver.Session{ID: 1}
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	p.envir = env
	p.CurrX, p.CurrY = 20, 20

	const n = 3
	for i := 0; i < n; i++ {
		p.ItemList = append(p.ItemList, &protocol.UserItem{MakeIndex: int32(100 + i), WIndex: 1})
	}

	server := &netserver.TCPServer{}
	for i := 0; i < n; i++ {
		before := len(env.GroundItems)
		p.HandleDropItem(SendMessage{Ident: int(protocol.CMDropItem), Param1: 100 + i}, server)
		if len(env.GroundItems) != before+1 {
			t.Fatalf("第 %d 次丢弃未落地（地面物品 %d -> %d）", i+1, before, len(env.GroundItems))
		}
	}
	if len(p.ItemList) != 0 {
		t.Errorf("背包应清空，仍剩 %d 件", len(p.ItemList))
	}
	// 落点应散开（范围 1 逐圈找空格），不全堆同一格。
	seen := make(map[[2]int]bool)
	for _, gi := range env.GroundItems {
		seen[[2]int{gi.X, gi.Y}] = true
	}
	if len(seen) != n {
		t.Errorf("%d 件物品落在 %d 个格子上（%v），应各不相同", n, len(seen), seen)
	}
}
