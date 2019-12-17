package main

import (
	"os"
	"sync"
)

// Environ ...
type Environ struct {
	Game               *Game
	GameDB             *GameDB
	SessionIDPlayerMap *sync.Map // map[int64]Player
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

// InitMaps ...
func (e *Environ) InitMaps() {
	mapDirPath := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/"
	//e.Maps = make([]Map, 386)
	e.Maps = make([]Map, 1)
	for _, mi := range e.GameDB.MapInfos {
		if mi.Filename == "0" {
			m := GetMapV1(GetMapBytes(mapDirPath + mi.Filename + ".map"))
			e.Maps[0] = *m
			break
		}
	}
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
