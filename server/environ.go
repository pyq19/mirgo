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
	Maps               *sync.Map // map[int]Map	// mapId: Map
}

// NewEnviron ...
func (g *Game) NewEnviron() (env *Environ) {
	env = new(Environ)
	env.Game = g
	env.InitGameDB()
	env.InitMaps()
	env.InitNPCs()
	env.InitRespawns()
	env.SessionIDPlayerMap = new(sync.Map)
	return
}

// p := *e.Game.Peer
// p.(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
// 	ses.Send(&ack)
// 	return true
// })
// StartLoop 玩家事件 / 怪物事件 / 地图事件 ...
func (e *Environ) StartLoop() {

	go e.Game.Pool.Run()
}

// GetMapInfoById FIXME 改成从 map 取出
func (e *Environ) GetMapInfoById(mapId int) *common.MapInfo {
	for _, v := range e.GameDB.MapInfos {
		if v.Id == mapId {
			return &v
		}
	}
	return nil
}

// GetItemInfoById FIXME 改成从 map 取出
func (e *Environ) GetItemInfoById(itemId int) *common.ItemInfo {
	for _, v := range e.GameDB.ItemInfos {
		if v.Id == int32(itemId) {
			return &v
		}
	}
	return nil
}

// InitNPCs 初始化地图上的 NPC
func (e *Environ) InitNPCs() {

}

// InitRespawns 初始化地图上的怪物
func (e *Environ) InitRespawns() {
	for _, ri := range e.GameDB.RespawnInfos {
		// FIXME 只加载第一张地图 测试用
		if ri.MapId != 1 {
			continue
		}
		v, ok := e.Maps.Load(ri.MapId)
		if !ok {
			continue
		}
		r := new(Respawn)
		r.Info = &ri
		m := v.(*Map)
		m.GetCell(common.NewPoint(ri.LocationX, ri.LocationY).String()).SetRespawn(r)
	}
}

func (e *Environ) GetMap(mapId int) *Map {
	v, ok := e.Maps.Load(mapId)
	if !ok {
		return nil
	}
	return v.(*Map)
}
