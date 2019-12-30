package main

import (
	"github.com/jinzhu/gorm"
	"github.com/yenkeia/mirgo/common"
	"io/ioutil"
	"os"
	"sync"
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
			v, _ := m.Cells.Load(common.Point{uint32(i), uint32(j)}.Coordinate())
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

func TestMap_GetNextCell(t *testing.T) {
	m := GetMapV1(GetMapBytes(os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/0.map"))
	c := &Cell{
		Map:        m,
		Coordinate: "100,200",
		Attribute:  0,
		Object:     nil,
		lock:       sync.RWMutex{},
	}
	t.Log(c.Coordinate)
	for i := 0; i < 8; i++ {
		//MirDirectionUp        MirDirection = 0
		//MirDirectionUpRight                = 1
		//MirDirectionRight                  = 2
		//MirDirectionDownRight              = 3
		//MirDirectionDown                   = 4
		//MirDirectionDownLeft               = 5
		//MirDirectionLeft                   = 6
		//MirDirectionUpLeft                 = 7
		nc := m.GetNextCell(c, common.MirDirection(i), 3)
		t.Log(nc.Coordinate)
	}
}
