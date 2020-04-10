package game

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/util"
)

func TestMapAbsPath(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, _ := gorm.Open("sqlite3", gopath+mirDB)

	mp := make([]cm.MapInfo, 386)
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
	m := LoadMap(mapAbsPath)
	// t.Log(m)
	str := ""
	for i := 0; i < int(m.Width); i++ {
		for j := 0; j < int(m.Height); j++ {
			c := m.GetCellXY(i, j)
			if c.Attribute == cm.CellAttributeWalk {
				str = str + "0"
			} else if c.Attribute == cm.CellAttributeHighWall {
				str = str + "1"
			} else if int(c.Attribute) == cm.CellAttributeLowWall {
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
	m := LoadMap(os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/0.map")
	c := &Cell{
		Point:     cm.Point{100, 200},
		Attribute: 0,
		// Objects:   nil,
	}
	t.Log(c.Point)
	for i := 0; i < 8; i++ {
		//MirDirectionUp        MirDirection = 0
		//MirDirectionUpRight                = 1
		//MirDirectionRight                  = 2
		//MirDirectionDownRight              = 3
		//MirDirectionDown                   = 4
		//MirDirectionDownLeft               = 5
		//MirDirectionLeft                   = 6
		//MirDirectionUpLeft                 = 7
		nc := m.GetNextCell(c, cm.MirDirection(i), 3)
		t.Log(nc.Point)
	}
}

func TestEnviron_LoadAllMap(t *testing.T) {
	mapDirPath := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/"
	uppercaseNameRealNameMap := make(map[string]string) // 目录下的文件名大写与该文件的真实文件名对应关系
	f, err := os.OpenFile(mapDirPath, os.O_RDONLY, os.ModeDir)
	if err != nil {
		panic(err)
	}
	fileInfo, _ := f.Readdir(-1)
	for _, info := range fileInfo {
		if !info.IsDir() {
			uppercaseNameRealNameMap[strings.ToUpper(info.Name())] = info.Name()
		}
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	for k, v := range uppercaseNameRealNameMap {
		t.Log(k, v)
	}
}

func TestMapRange(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	mapAbsPath := gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/0.map"
	m := LoadMap(mapAbsPath)

	p := cm.Point{X: 1, Y: 1}

	var printpos = func(c *Cell, _, _ int) bool {
		if c != nil {
			fmt.Println(c.Point)
		} else {
			fmt.Println("nil nil")
		}
		return true
	}

	fmt.Println("---", m.Width, m.Height)
	m.RangeCell(p, 0, printpos)

	fmt.Println("---")
	m.RangeCell(p, 1, printpos)

	fmt.Println("---")
	m.RangeCell(p, 2, printpos)
}

func TestAllMaps(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	mappath := gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/"

	maps := util.GetFiles(mappath, []string{".map"})

	mark := map[byte]bool{}

	for _, m := range maps {
		bytes, _ := ioutil.ReadFile(m)
		mark[DetectMapVersion(bytes)] = true
	}

	for k, v := range mark {
		fmt.Println(k, v)
	}
}
