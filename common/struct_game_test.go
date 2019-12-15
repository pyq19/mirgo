package common

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestMapAbsPath(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, _ := gorm.Open("sqlite3", gopath+mirDB)

	mp := make([]MapInfo, 386)
	db.Table("map").Find(&mp)

	mapDirPath := "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/"
	fileName := mp[0].Filename
	mapAbsPath := gopath + mapDirPath + fileName + ".map"
	t.Log(mapAbsPath)
}

func TestLoadMap(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	mapPath := "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/0.map"
	mapAbsPath := gopath + mapPath
	t.Log(mapAbsPath)
}

func TestSaveMapText(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	filePath := gopath + "/src/github.com/yenkeia/mirgo/01.txt"
	t.Log(filePath)
}
