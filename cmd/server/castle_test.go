package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// 批6回归：城门/城墙可视化（Delphi TCastleDoor/TWallStructure，
// ObjMon2.pas:1013-1096 / HerbActor.pas:365-683）。

func newCastleSetup(t *testing.T) (*CastleObject, *Environment, *UserEngine) {
	t.Helper()
	env := newTestEnv(40, 40)
	mm := &MapManager{maps: map[string]*Environment{"3": env}}
	e := &UserEngine{Config: &ServerConfig{}, mapMgr: mm}
	cfg := DefaultCastleConfig()
	cfg.PalaceX = 20
	cfg.PalaceY = 20
	cfg.PalaceRadius = 3
	e.Castle = NewCastleObject(*cfg)
	env.Castle = e.Castle
	return e.Castle, env, e
}

func TestDoorVisualDir(t *testing.T) {
	c, _, _ := newCastleSetup(t)
	max := c.Config.DoorMaxHP

	c.DoorHP = max
	if d := c.doorVisualDirLocked(); d != 0 {
		t.Errorf("满血关闭 dir = %d, want 0", d)
	}
	c.DoorHP = max * 2 / 3
	if d := c.doorVisualDirLocked(); d != 1 {
		t.Errorf("2/3 血 dir = %d, want 1", d)
	}
	c.DoorHP = max / 3
	if d := c.doorVisualDirLocked(); d != 2 {
		t.Errorf("1/3 血 dir = %d, want 2", d)
	}
	c.DoorHP = max
	c.DoorOpen = true
	if d := c.doorVisualDirLocked(); d != 7 {
		t.Errorf("开启 dir = %d, want 7", d)
	}
}

func TestSyncDoorVisualTransitions(t *testing.T) {
	c, env, e := newCastleSetup(t)
	server := &netserver.TCPServer{}
	doorX, doorY := c.DoorPosition()

	door := placeNamedMonster(env, "SabukDoor", 900, 110, doorX, doorY)
	c.DoorActorID = door.ID
	c.doorVisualState = 0

	p := newButchPlayer(env, 1, doorX+2, doorY)
	hasMsg := func(ident int) bool {
		for _, m := range p.msgList {
			if m.Ident == ident && m.SourceID == door.ID {
				return true
			}
		}
		return false
	}

	// 开门 → RM_DIGUP(dir=7)
	c.ToggleDoor(true)
	c.SyncDoorVisual(server, e)
	if door.Dir != 7 || !hasMsg(RM_DIGUP) {
		t.Fatal("开门应广播 RM_DIGUP 且 door.Dir=7")
	}

	// 关门 → RM_DIGDOWN(dir=破损度0)
	c.ToggleDoor(false)
	c.SyncDoorVisual(server, e)
	if !hasMsg(RM_DIGDOWN) {
		t.Fatal("关门应广播 RM_DIGDOWN")
	}

	// 破坏 → 死亡广播
	c.DoorHP = 0
	c.SyncDoorVisual(server, e)
	if !door.Death || !hasMsg(RM_DEATH) {
		t.Fatal("城门破坏应置 Death 并广播 RM_DEATH")
	}

	// 修复 → RM_ALIVE
	c.DoorHP = c.Config.DoorMaxHP
	c.SyncDoorVisual(server, e)
	if door.Death || !hasMsg(RM_ALIVE) {
		t.Fatal("城门修复应清除 Death 并广播 RM_ALIVE")
	}

	// 重复同步不重复广播
	before := len(p.msgList)
	c.SyncDoorVisual(server, e)
	if len(p.msgList) != before {
		t.Fatal("状态未变化不应重复广播")
	}
}

func TestStructureBlocking(t *testing.T) {
	c, env, _ := newCastleSetup(t)
	doorX, doorY := c.DoorPosition()
	door := placeNamedMonster(env, "SabukDoor", 900, 110, doorX, doorY)
	door.IsStructure = true
	door.Race = 110
	wall := placeNamedMonster(env, "SabukW1", 901, 111, 17, 17)
	wall.IsStructure = true
	wall.Race = 111

	// 城门实体本身不阻挡（通行由 IsBlockedByDoor 判定）
	if blocksMovement(door) {
		t.Fatal("城门实体不应直接阻挡（由 IsBlockedByDoor 按状态判定）")
	}
	if !c.IsBlockedByDoor(doorX, doorY) {
		t.Fatal("关闭城门应阻挡通行")
	}
	c.ToggleDoor(true)
	if c.IsBlockedByDoor(doorX, doorY) {
		t.Fatal("开启城门不应阻挡通行")
	}

	// 城墙存活阻挡，破坏后放行
	if !blocksMovement(wall) {
		t.Fatal("完整城墙应阻挡通行")
	}
	wall.Death = true
	if blocksMovement(wall) {
		t.Fatal("破坏的城墙不应阻挡通行")
	}
}

func TestStructureNoDropNoGhost(t *testing.T) {
	c, env, e := newCastleSetup(t)
	e.MonsterDB = &MonsterDB{}
	doorX, doorY := c.DoorPosition()
	door := placeNamedMonster(env, "SabukDoor", 900, 110, doorX, doorY)
	door.IsStructure = true
	door.Death = true
	door.DeathTick = 1 // 远超 corpseDelay
	e.Monsters = []*MonsterObject{door}

	e.ProcessMonsters(&netserver.TCPServer{}, 200000)

	if door.Ghost || door.LootDropped {
		t.Fatal("城门尸体不应 ghost、不应走掉落")
	}
}
