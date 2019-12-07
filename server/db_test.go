package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yenkeia/mirgo/com"
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db := NewDB(path)
	defer db.Close()

	var basic com.Basic
	db.Table("basic").Where("game_version = ?", 65).Find(&basic)

	t.Log(basic.GameVersion, basic.MapIndex, basic.RespawnIndex)
}
