package common

// DropInfo 一个 DropInfo 对象对应 database/Envir/Drops 文件夹内一个文件的一行
type DropInfo struct {
	Chance int
	Gold   int
	ItemID int // ItemInfo.ID
}

func NewDropInfo(line string) *DropInfo {
	return nil
}

// GetDropInfosByMonsterName 加载怪物掉落物品
func GetDropInfosByMonsterName(dropDirPath, monsterName string) []DropInfo {
	return nil
}
