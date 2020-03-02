package mir

import (
	"github.com/jinzhu/gorm"
)

type DB struct {
	db *gorm.DB
}

// adb 保存游戏运行时生成的数据
var adb *DB

func (d *DB) setCharacterAttr(p *Player, attr string, value interface{}) {
	d.db.Table("character").Where("id = ?", p.ID).Update(attr, value)
}

func (d *DB) SyncLevel(p *Player) {
	d.setCharacterAttr(p, "level", p.Level)
}

func (d *DB) SyncGold(p *Player) {
	d.setCharacterAttr(p, "gold", p.Gold)
}
