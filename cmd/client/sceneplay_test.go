package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/engine"
)

// TestObjectsCacheLazyInit 防 "assignment to entry in nil map" 回归：
// frontImageData 直接把 s.objectsCaches[area] 传给 getTex 写入
//（sceneplay.go），getObjectsLoader 必须保证该 area 的缓存 map 已创建，
// 无论 LoadMap 是否跑过、resources 是否就绪。
func TestObjectsCacheLazyInit(t *testing.T) {
	s := NewPlayScene(nil, &engine.ResourceManager{}, "", nil)
	loader := s.getObjectsLoader(0)
	if loader != nil {
		t.Fatalf("无资源时 loader 应为 nil, got %v", loader)
	}
	if s.objectsCaches[0] == nil {
		t.Fatal("objectsCaches[0] 未初始化 — getTex 写入会 panic")
	}
	// 模拟 getTex 的写入路径（崩溃点）。
	s.objectsCaches[0][1382] = 42
	if got := s.objectsCaches[0][1382]; got != 42 {
		t.Fatalf("缓存写入读回 = %d, want 42", got)
	}
	// 二次调用不得重置已有缓存（纹理 ID 跨地图复用）。
	s.getObjectsLoader(0)
	if got := s.objectsCaches[0][1382]; got != 42 {
		t.Fatalf("二次 getObjectsLoader 后缓存被重置: got %d, want 42", got)
	}
}
