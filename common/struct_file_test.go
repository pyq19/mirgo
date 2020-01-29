package common

import (
	"os"
	"testing"
)

func TestNewDropInfo(t *testing.T) {
	di1 := NewDropInfo("1/10 Gold 500")
	t.Log(di1)
	di2 := NewDropInfo("1/5 (MP)DrugLarge")
	t.Log(di2)
}

func TestGetDropInfosByMonsterName(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	monsterName := "Ancient_AxeSkeleton"
	dropDirPath := gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Envir/Drops/"
	drops := GetDropInfosByMonsterName(dropDirPath, monsterName)
	t.Log(drops)
}
