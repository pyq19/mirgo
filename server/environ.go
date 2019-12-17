package main

import (
	"github.com/yenkeia/mirgo/common"
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

func (e *Environ) GetMapInfoById(mapId int) *common.MapInfo {
	for _, m := range e.GameDB.MapInfos {
		if m.Id == mapId {
			return &m
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

}
