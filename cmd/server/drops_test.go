package main

import (
	"testing"
	"time"
)

// getDropPosition 语义回归（Delphi GetDropPosition，ObjBase.pas:1534-1584）。

// 空地：返回第一圈扫描次序的首格（dy 外层、dx 内层，确定性落点）。
func TestGetDropPositionEmpty(t *testing.T) {
	env := newTestEnv(20, 20)
	x, y := getDropPosition(env, 10, 10, 1)
	if x != 9 || y != 9 {
		t.Errorf("getDropPosition = (%d,%d), want (9,9)", x, y)
	}
}

// 第一圈（含原点，Delphi 扫描含 (0,0)）被物品占满：落到第二圈首个空格。
func TestGetDropPositionSkipsOccupiedRing(t *testing.T) {
	env := newTestEnv(20, 20)
	now := time.Now().UnixMilli()
	id := int32(1)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			gi := &GroundItem{ID: id, Name: "金创药", X: 10 + dx, Y: 10 + dy, DropTick: now}
			if env.AddGroundItem(gi) == nil {
				t.Fatalf("setup: 放置 (%d,%d) 失败", 10+dx, 10+dy)
			}
			id++
		}
	}
	x, y := getDropPosition(env, 10, 10, 3)
	if x != 8 || y != 8 {
		t.Errorf("getDropPosition = (%d,%d), want (8,8)", x, y)
	}
}

// 墙格（不可行走）被跳过。
func TestGetDropPositionSkipsWalls(t *testing.T) {
	env := newTestEnv(20, 20)
	env.Cells[9*env.Width+9].Flag = 1 // 扫描首格 (9,9) 设为墙
	x, y := getDropPosition(env, 10, 10, 1)
	if x != 10 || y != 9 {
		t.Errorf("getDropPosition = (%d,%d), want (10,9)", x, y)
	}
}

// 范围内无空格：退回物品最少的格子。
func TestGetDropPositionFallbackLeastItems(t *testing.T) {
	env := newTestEnv(20, 20)
	now := time.Now().UnixMilli()
	id := int32(1)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			n := 2
			if dx == 1 && dy == 1 {
				n = 1 // (11,11) 物品最少
			}
			for k := 0; k < n; k++ {
				gi := &GroundItem{ID: id, Name: "金创药", X: 10 + dx, Y: 10 + dy, DropTick: now}
				if env.AddGroundItem(gi) == nil {
					t.Fatalf("setup: 放置 (%d,%d) 失败", 10+dx, 10+dy)
				}
				id++
			}
		}
	}
	x, y := getDropPosition(env, 10, 10, 1)
	if x != 11 || y != 11 {
		t.Errorf("getDropPosition = (%d,%d), want (11,11) 物品最少的格", x, y)
	}
}

// 范围内全部堆满上限：退回原点。
func TestGetDropPositionFallbackOrigin(t *testing.T) {
	env := newTestEnv(20, 20)
	now := time.Now().UnixMilli()
	id := int32(1)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			for k := 0; k < maxGroundItemsPerCell; k++ {
				gi := &GroundItem{ID: id, Name: "金创药", X: 10 + dx, Y: 10 + dy, DropTick: now}
				if env.AddGroundItem(gi) == nil {
					t.Fatalf("setup: 放置 (%d,%d) 失败", 10+dx, 10+dy)
				}
				id++
			}
		}
	}
	x, y := getDropPosition(env, 10, 10, 1)
	if x != 10 || y != 10 {
		t.Errorf("getDropPosition = (%d,%d), want (10,10) 原点兜底", x, y)
	}
}

// 回归（"背包扔物品全部堆叠"）：连续掉落逐格散开，落点互不重复。
func TestConsecutiveDropsScatter(t *testing.T) {
	env := newTestEnv(20, 20)
	now := time.Now().UnixMilli()
	seen := make(map[[2]int]bool)
	for i := int32(0); i < 8; i++ {
		x, y := getDropPosition(env, 10, 10, 3)
		if seen[[2]int{x, y}] {
			t.Fatalf("第 %d 件落在已有物品的格子 (%d,%d)", i+1, x, y)
		}
		seen[[2]int{x, y}] = true
		gi := &GroundItem{ID: i + 1, Name: "金创药", X: x, Y: y, DropTick: now}
		if env.AddGroundItem(gi) == nil {
			t.Fatalf("第 %d 件落地失败", i+1)
		}
	}
}
