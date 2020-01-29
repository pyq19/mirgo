package common

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// DropInfo 一个 DropInfo 对象对应 database/Envir/Drops 文件夹内一个文件的一行
type DropInfo struct {
	Chance   int // 1/3 中的 3，1/10 中的 10
	Gold     int
	ItemName string // ItemInfo.Name
}

func NewDropInfo(line string) *DropInfo {
	line = strings.TrimSpace(line)
	res := strings.Split(line, " ")
	chance, err := strconv.Atoi(strings.Split(res[0], "/")[1])
	if err != nil {
		panic("NewDropInfo 格式不正确" + err.Error())
	}
	if len(res) == 3 {
		// 1/10 Gold 500
		gold, _ := strconv.Atoi(res[2])
		return &DropInfo{Chance: chance, Gold: gold, ItemName: "Gold"}
	} else if len(res) == 2 {
		// 1/5 (MP)DrugLarge
		return &DropInfo{Chance: chance, Gold: 0, ItemName: res[1]}
	} else {
		panic("NewDropInfo 格式不正确")
	}
}

// GetDropInfosByMonsterName 加载怪物掉落物品
func GetDropInfosByMonsterName(dropDirPath, monsterName string) (res []DropInfo) {
	filename := dropDirPath + monsterName + ".txt"
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil
	}
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		line := fscanner.Text()
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		res = append(res, *NewDropInfo(line))
	}
	return
}
