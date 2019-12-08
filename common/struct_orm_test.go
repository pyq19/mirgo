package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"testing"
)

var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"

func TestItemInfo(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	db, err := gorm.Open("sqlite3", gopath+mirDB)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var itemInfo ItemInfo
	db.Table("item_info").Where("item_index = ?", "1").Find(&itemInfo)
	t.Log(itemInfo)
}
