package main

import (
	"github.com/davyxu/cellnet"
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/proto/server"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Environ ...
type Environ struct {
	Game               *Game
	GameDB             *GameDB
	SessionIDPlayerMap *sync.Map // map[int64]*Player
	Maps               *sync.Map // map[int]*Map	// mapID: Map
	ObjectID           uint32
	Players            []*Player
	lock               *sync.Mutex
}

// NewEnviron ...
func NewEnviron(g *Game) (env *Environ) {
	env = new(Environ)
	env.Game = g
	env.InitGameDB()
	env.InitMaps()
	env.ObjectID = 100000
	env.Players = make([]*Player, 10)
	env.lock = new(sync.Mutex)
	err := env.InitObjects()
	if err != nil {
		panic(err)
	}
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
	var um []common.UserMagic
	db.Table("user_magic").Find(&um)
	gdb.UserMagics = um
}

// InitMaps ...
func (e *Environ) InitMaps() {
	//e.Maps = make([]Map, 386)
	e.Maps = new(sync.Map)
	for _, mi := range e.GameDB.MapInfos {
		mi := mi
		if mi.ID == 1 {
			m := GetMapV1(GetMapBytes(e.Game.Conf.MapDirPath + mi.Filename + ".map"))
			m.Env = e
			m.Info = &mi
			e.Maps.Store(1, m)
			break
		}
	}
}

func (e *Environ) NewObjectID() uint32 {
	return atomic.AddUint32(&e.ObjectID, 1)
}

// InitObjects 初始化地图
func (e *Environ) InitObjects() (err error) {
	var maps []*Map
	e.Maps.Range(func(k, v interface{}) bool {
		maps = append(maps, v.(*Map))
		return true
	})
	for _, m := range maps {
		err = m.InitMonsters()
		if err != nil {
			return err
		}
		err = m.InitNPCs()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Environ) AddPlayer(p *Player) {
	e.lock.Lock()
	e.Players = append(e.Players, p)
	e.lock.Unlock()
}

func (e *Environ) DeletePlayer(p *Player) {
	e.lock.Lock()
	for i := 0; i < len(e.Players); i++ {
		o := e.Players[i]
		if p.ID == o.ID {
			e.Players[i] = e.Players[len(e.Players)-1]
			e.Players = e.Players[:len(e.Players)-1]
			break
		}
	}
	e.lock.Unlock()
}

func (e *Environ) GetPlayersCount() int {
	e.lock.Lock()
	c := 0
	for i := range e.Players {
		if e.Players[i] != nil {
			c++
		}
	}
	e.lock.Unlock()
	return c
}

func (e *Environ) GetMap(mapID int) *Map {
	v, ok := e.Maps.Load(mapID)
	if !ok {
		return nil
	}
	return v.(*Map)
}

func (e *Environ) Submit(t *Task) {
	e.Game.Pool.EntryChan <- t
}

func (e *Environ) Broadcast(msg interface{}) {
	e.Maps.Range(func(k, v interface{}) bool {
		v.(*Map).Broadcast(msg)
		return true
	})
}

// StartLoop
func (e *Environ) StartLoop() {
	go e.TimeTick()
	go e.Game.Pool.Run()
}

func (e *Environ) TimeTick() {
	// 系统事件 广播 存档
	systemBroadcastTicker := time.NewTicker(10 * time.Second)

	// 玩家事件 buff 等状态改变

	// 地图事件 刷怪 地图物品

	// 怪物事件 移动 buff
	monsterMoveTicker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-systemBroadcastTicker.C:
			e.Submit(NewTask(e.systemBroadcast))
		case <-monsterMoveTicker.C:
			e.Submit(NewTask(e.monsterMove))
		}
	}
}

func (e *Environ) systemBroadcast(...interface{}) {
	text := "当前在线玩家人数: " + strconv.Itoa(e.GetPlayersCount())
	(*e.Game.Peer).(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
		ses.Send(&server.Chat{
			Message: text,
			Type:    common.ChatTypeSystem,
		})
		return true
	})
}

func (e *Environ) monsterMove(...interface{}) {

}
