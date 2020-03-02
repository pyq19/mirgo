package mir

import (
	"github.com/jinzhu/gorm"
	"github.com/yenkeia/mirgo/common"
)

type DB struct {
	db *gorm.DB
}

// adb 保存游戏运行时生成的数据
var adb *DB

func (d *DB) DelItem(p *Player, item *common.UserItem) {
	d.db.Table("user_item").Where("id = ?", item.ID).Delete(&common.UserItem{})
	d.db.Table("character_user_item").Where("useritemid = ?", item.ID).Delete(&common.CharacterUserItem{})
}

func (d *DB) AddItem(p *Player, t common.UserItemType, index int, item *common.UserItem) {
	d.db.Table("user_item").Create(item)

	d.db.Table("character_user_item").Create(&common.CharacterUserItem{
		CharacterID: int(p.ID),
		UserItemID:  int(item.ID),
		Type:        int(t),
		Index:       index,
	})
}

func (d *DB) setCharacterAttr(p *Player, attr string, value interface{}) {
	d.db.Table("character").Where("id = ?", p.ID).Update(attr, value)
}

func (d *DB) SyncLevel(p *Player) {
	d.setCharacterAttr(p, "level", p.Level)
}

func (d *DB) SyncGold(p *Player) {
	d.setCharacterAttr(p, "gold", p.Gold)
}
