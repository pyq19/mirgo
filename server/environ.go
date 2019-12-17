package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// Environ ...
type Environ struct {
	Game               *Game
	GameDB             *GameDB
	SessionIDPlayerMap *sync.Map // map[int64]Player
	// CommonEventChan  chan interface{} // 系统事件
	// PlayerEventChan  chan interface{} // 玩家事件
	// MonsterEventChan chan interface{} // 怪物事件
	// MapEventChan     chan interface{} // 地图事件
}

// GameDB ...
type GameDB struct {
	Basic         common.Basic
	GameShopItems []common.GameShopItem
	ItemInfos     []common.ItemInfo
	MagicInfos    []common.MagicInfo
	MapInfos      []common.MapInfo
	MonsterInfos  []common.MonsterInfo
	MovementInfos []common.MovementInfo
	NpcInfos      []common.NpcInfo
	QuestInfos    []common.QuestInfo
	RespawnInfos  []common.RespawnInfo
	SafeZoneInfos []common.SafeZoneInfo
}

const (
	LOGIN = iota
	SELECT
	GAME
	DISCONNECTED
)

// Player ...
type Player struct {
	AccountId int
	GameStage int
	Session   *cellnet.Session
	Character *common.Character
	Magics    *[]common.MagicInfo
	UserItems *[]common.UserItem
}

// NewEnviron ...
func (g *Game) NewEnviron() (env *Environ) {
	env = new(Environ)
	env.Game = g
	env.LoadGameDB()
	env.InitMaps()
	env.InitMonsters()
	env.SessionIDPlayerMap = new(sync.Map)
	return
}

// LoadGameDB ...
func (e *Environ) LoadGameDB() {
	gdb := new(GameDB)
	e.GameDB = gdb
	db := e.Game.DB
	b := new(common.Basic)
	db.Table("basic").Find(b)
	gdb.Basic = *b
	gsi := make([]common.GameShopItem, 106)
	db.Table("game_shop_item").Find(&gsi)
	gdb.GameShopItems = gsi
	ii := make([]common.ItemInfo, 1346)
	db.Table("item").Find(&ii)
	gdb.ItemInfos = ii
	mi := make([]common.MagicInfo, 105)
	db.Table("magic").Find(&mi)
	gdb.MagicInfos = mi
	mp := make([]common.MapInfo, 386)
	db.Table("map").Find(&mp)
	gdb.MapInfos = mp
	ms := make([]common.MonsterInfo, 506)
	db.Table("monster").Find(&ms)
	gdb.MonsterInfos = ms
	mm := make([]common.MovementInfo, 1837)
	db.Table("movement").Find(&mm)
	gdb.MovementInfos = mm
	ni := make([]common.NpcInfo, 293)
	db.Table("npc").Find(&ni)
	gdb.NpcInfos = ni
	qi := make([]common.QuestInfo, 157)
	db.Table("quest").Find(&qi)
	gdb.QuestInfos = qi
	ri := make([]common.RespawnInfo, 5931)
	db.Table("respawn").Find(&ri)
	gdb.RespawnInfos = ri
	si := make([]common.SafeZoneInfo, 19)
	db.Table("safe_zone").Find(&si)
	gdb.SafeZoneInfos = si
}

// InitMaps ...
func (e *Environ) InitMaps() {

}

// InitMonsters ...
func (e *Environ) InitMonsters() {

}

// StartLoop ...
func (e *Environ) StartLoop() {
	// p := *e.Game.Peer
	// p.(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
	// 	ses.Send(&ack)
	// 	return true
	// })
}
