package main

import (
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"os"
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
func NewEnviron(g *Game) (env *Environ) {
	env = new(Environ)
	env.Game = g
	env.InitGameDB()
	env.InitMaps()
	env.InitNPCs()
	env.InitRespawns()
	env.SessionIDPlayerMap = new(sync.Map)
	return
}

// InitGameDB ...
func (e *Environ) InitGameDB() {
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
	mapDirPath := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/"
	//e.Maps = make([]Map, 386)
	e.Maps = new(sync.Map)
	for _, mi := range e.GameDB.MapInfos {
		if mi.Id == 1 {
			m := GetMapV1(GetMapBytes(mapDirPath + mi.Filename + ".map"))
			m.Id = mi.Id
			e.Maps.Store(1, m)
			break
		}
	}
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
		v.(*Map).AddObject(r)
	}
}

func (e *Environ) GetMap(mapId int) *Map {
	v, ok := e.Maps.Load(mapId)
	if !ok {
		return nil
	}
	return v.(*Map)
}
