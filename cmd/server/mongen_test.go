package main

import (
	"os"
	"path/filepath"
	"testing"
)

// mon_gen.jsonc 的 interval 单位是分钟（Delphi LocalDB.pas:1237
// dwZenTime := Str_ToInt(sData,-1)*60*1000），防止再被误按秒解释。
func TestMonGenIntervalIsMinutes(t *testing.T) {
	dir := t.TempDir()
	content := `{
  "spawns": [
    {"mapName": "0", "x": 100, "y": 100, "name": "鸡", "range": 5, "count": 3, "interval": 10},
    {"mapName": "0", "x": 200, "y": 200, "name": "鹿", "range": 5, "count": 2, "interval": 0}
  ]
}`
	path := filepath.Join(dir, "mon_gen.jsonc")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	e := &UserEngine{monGenPath: path} // mapMgr 为 nil 时不过滤地图
	if !e.loadMonGenFromFile("0") {
		t.Fatal("loadMonGenFromFile 失败")
	}
	if len(e.MonGenList) != 2 {
		t.Fatalf("期望 2 个刷新点，实际 %d", len(e.MonGenList))
	}
	if got := e.MonGenList[0].ZenTime; got != 10*60*1000 {
		t.Errorf("interval=10 分钟应为 600000ms，实际 %d", got)
	}
	// interval<=0 回退默认 10 分钟
	if got := e.MonGenList[1].ZenTime; got != 10*60*1000 {
		t.Errorf("interval=0 应回退默认 10 分钟，实际 %d", got)
	}
}
