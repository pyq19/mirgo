package main

import (
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"sync"
)

// Environ ...
type Environ struct {
	Game               *Game
	GameDB             *GameDB
	SessionIDPlayerMap *sync.Map // map[int64]*Player
	Maps               []Map
}

// NewEnviron ...
func (g *Game) NewEnviron() (env *Environ) {
	env = new(Environ)
	env.Game = g
	env.InitGameDB()
	env.InitMaps()
	env.SessionIDPlayerMap = new(sync.Map)
	return
}

// FIXME 改成从 map 取出
func (e *Environ) GetMapInfoById(mapId int) *common.MapInfo {
	for _, v := range e.GameDB.MapInfos {
		if v.Id == mapId {
			return &v
		}
	}
	return nil
}

// FIXME 改成从 map 取出
func (e *Environ) GetItemInfoById(itemId int) *common.ItemInfo {
	for _, v := range e.GameDB.ItemInfos {
		if v.Id == int32(itemId) {
			return &v
		}
	}
	return nil
}

// p := *e.Game.Peer
// p.(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
// 	ses.Send(&ack)
// 	return true
// })
// CommonEventChan  chan interface{} // 系统事件
// PlayerEventChan  chan interface{} // 玩家事件
// MonsterEventChan chan interface{} // 怪物事件
// MapEventChan     chan interface{} // 地图事件
// StartLoop ...
func (e *Environ) StartLoop() {

	go e.Game.Pool.Run()
}
