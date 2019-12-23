package main

import (
	"github.com/jinzhu/gorm"
	"github.com/yenkeia/mirgo/common"
	"io/ioutil"
	"os"
	"testing"
)

func TestMapAbsPath(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, _ := gorm.Open("sqlite3", gopath+mirDB)

	mp := make([]common.MapInfo, 386)
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
	mapAbsPath := gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/0.map"
	m := GetMapV1(GetMapBytes(mapAbsPath))
	// t.Log(m)
	str := ""
	for i := 0; i < int(m.Width); i++ {
		for j := 0; j < int(m.Height); j++ {
			v, _ := m.CoordinateCellMap.Load(common.Point{uint32(i), uint32(j)}.Coordinate())
			c := v.(*Cell)
			if c.Attribute == common.CellAttributeWalk {
				str = str + "0"
			} else if c.Attribute == common.CellAttributeHighWall {
				str = str + "1"
			} else if int(c.Attribute) == common.CellAttributeLowWall {
				str = str + "2"
			} else {
				str = str + "?"
			}
		}
		str = str + "\n"
	}
	ioutil.WriteFile(filePath, []byte(str), 0644)
}
