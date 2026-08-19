package main

import (
	"testing"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// TestPickupOnlyAtFeet 回归：只能拾取人物脚下格子的物品。
// Delphi ClientPickUpItem 只取 GetItem(m_nCurrX, m_nCurrY)
//（ObjBase.pas:1696），不扫描周围格子；旧 Go 实现扫 3×3，
// 站着不动就能吸走身边物品。
func TestPickupOnlyAtFeet(t *testing.T) {
	env := newTestEnv(40, 40)
	p := &PlayObject{BaseObject: NewBaseObject("picker", 1)}
	p.Session = &netserver.Session{ID: 1}
	p.Engine = &UserEngine{Config: &ServerConfig{}}
	p.envir = env
	p.CurrX, p.CurrY = 20, 20

	now := time.Now().UnixMilli()
	server := &netserver.TCPServer{}

	// 相邻格的物品不可拾取
	adjacent := &GroundItem{ID: 1, Name: "金创药", X: 21, Y: 20, DropTick: now,
		UserItem: &protocol.UserItem{MakeIndex: 1, WIndex: 1}}
	if env.AddGroundItem(adjacent) == nil {
		t.Fatal("setup: 相邻格物品放置失败")
	}
	p.HandlePickup(SendMessage{Ident: int(protocol.CMPickup)}, server)
	if len(p.ItemList) != 0 {
		t.Errorf("相邻格物品不应被拾取，背包却有 %d 件", len(p.ItemList))
	}
	if env.GetGroundItemAt(21, 20) == nil {
		t.Error("相邻格物品应仍在原地")
	}

	// 脚下格的物品可拾取
	feet := &GroundItem{ID: 2, Name: "金创药", X: 20, Y: 20, DropTick: now,
		UserItem: &protocol.UserItem{MakeIndex: 2, WIndex: 1}}
	if env.AddGroundItem(feet) == nil {
		t.Fatal("setup: 脚下格物品放置失败")
	}
	p.HandlePickup(SendMessage{Ident: int(protocol.CMPickup)}, server)
	if len(p.ItemList) != 1 {
		t.Fatalf("脚下格物品应被拾取，背包有 %d 件", len(p.ItemList))
	}
	if env.GetGroundItemAt(20, 20) != nil {
		t.Error("脚下格物品应已从地面移除")
	}
}
