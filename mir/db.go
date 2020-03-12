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

func NewAccountDB(db *gorm.DB) *DB {

	r := &DB{db: db}

	db.SingularTable(true)
	db.AutoMigrate(
		&common.Basic{},
		&common.Account{},
		&common.AccountCharacter{},
		&common.Character{},
		&common.CharacterUserItem{},
		&common.UserItem{},
		&common.UserMagic{},
	)

	return r
}

func (d *DB) Table(name string) *gorm.DB {
	return d.db.Table(name)
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

func (d *DB) SyncPosition(p *Player) {
	d.setCharacterAttr(p, "current_map_id", p.Map.Info.ID)
	d.setCharacterAttr(p, "direction", p.GetDirection())
	d.setCharacterAttr(p, "current_location_x", p.GetPoint().X)
	d.setCharacterAttr(p, "current_location_y", p.GetPoint().Y)
}

func (d *DB) SyncMagicKey(p *Player, spell common.Spell, key uint8) {
	d.db.Table("user_magic").Where("spell = ?", spell).Update("magic_key", key)
}

func (d *DB) AddSkill(p *Player, magic *common.UserMagic) {
	d.db.Table("user_magic").Create(magic)
}

func (d *DB) GetObjectID() uint32 {
	var basic common.Basic
	d.db.Table("basic").First(&basic)
	return basic.ObjectID
}

func (d *DB) SyncObjectID(id uint32) {
	d.db.Table("basic").Where(common.Basic{ID: 1}).Update(map[string]interface{}{"object_id": id})
}
