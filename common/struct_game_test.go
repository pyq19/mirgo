package common

import (
	"io/ioutil"
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
	mapAbsPath := gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/0.map"
	m := GetMapV1(GetMapBytes(mapAbsPath))
	// t.Log(m)
	str := ""
	for i := 0; i < int(m.Width); i++ {
		for j := 0; j < int(m.Height); j++ {
			c := m.CoordinateCellMap[Point{uint32(i), uint32(j)}.String()]
			if int(c.Attribute) == 0 {
				str = str + "#"
			} else if int(c.Attribute) == 1 {
				str = str + "1"
			} else if int(c.Attribute) == 2 {
				str = str + "2"
			} else {
				str = str + "?"
			}
		}
		str = str + "\n"
	}
	ioutil.WriteFile(filePath, []byte(str), 0644)
}
