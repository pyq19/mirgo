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
	env.Players = make([]*Player, 0)
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
	gdb.Init()
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
	p.Map.AddObject(p)
}

func (e *Environ) GetPlayer(ID uint32) *Player {
	e.lock.Lock()
	for i := 0; i < len(e.Players); i++ {
		o := e.Players[i]
		if ID == o.ID {
			return o
		}
	}
	e.lock.Unlock()
	return nil
}

func (e *Environ) DeletePlayer(p *Player) {
	e.lock.Lock()
	for i := 0; i < len(e.Players); i++ {
		o := e.Players[i]
		if o == nil {
			continue
		}
		if p.ID == o.ID {
			e.Players[i] = e.Players[len(e.Players)-1]
			e.Players = e.Players[:len(e.Players)-1]
			break
		}
	}
	e.lock.Unlock()
	p.Map.DeleteObject(p)
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
	systemBroadcastTicker := time.NewTicker(1 * time.Hour)

	debugTicker := time.NewTicker(10 * time.Second)

	// 地图事件 刷怪 地图物品

	// 玩家事件 buff 等状态改变

	// 怪物事件 移动 buff
	monsterProcessTicker := time.NewTicker(500 * time.Millisecond)

	// NPC
	npcProcessTicker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-systemBroadcastTicker.C:
			e.Submit(NewTask(e.SystemBroadcast))
		case <-debugTicker.C:
			e.Debug()
		case <-monsterProcessTicker.C:
			e.Submit(NewTask(e.MonstersProcess))
		case <-npcProcessTicker.C:
			e.Submit(NewTask(e.NPCsProcess))
		}
	}
}

func (e *Environ) SystemBroadcast(...interface{}) {
	envPlayerCount := e.GetPlayersCount()
	text := "当前在线玩家人数: " + strconv.Itoa(envPlayerCount)
	(*e.Game.Peer).(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
		ses.Send(&server.Chat{
			Message: text,
			Type:    common.ChatTypeSystem,
		})
		return true
	})

}

func (e *Environ) Debug() {
	envPlayerCount := e.GetPlayersCount()
	allPlayer := make([]*Player, 0)
	e.Maps.Range(func(k, v interface{}) bool {
		m := v.(*Map)
		allPlayer = append(allPlayer, m.GetAllPlayers()...)
		return true
	})
	if len(allPlayer) != envPlayerCount {
		log.Errorf("!!! warning envPlayerCount: %d != map allPlayer: %d\n", envPlayerCount, len(allPlayer))
	} else {
		log.Debugf("envPlayerCount: %d, map allPlayer: %d\n", envPlayerCount, len(allPlayer))
	}
}

func (e *Environ) GetActiveObjects() (monster []*Monster, npc []*NPC) {
	e.lock.Lock()
	defer e.lock.Unlock()
	gridMap := make(map[int]*Grid)
	for i := range e.Players {
		g := e.Players[i].GetCurrentGrid()
		gridMap[g.GID] = g
	}
	grids := make([]*Grid, 0)
	for _, g := range gridMap {
		grids = append(grids, g)
	}
	for i := range grids {
		g := grids[i]
		objs := g.GetAllObjects()
		for i := range objs {
			o := objs[i]
			switch o.GetRace() {
			case common.ObjectTypeMonster:
				monster = append(monster, o.(*Monster))
			case common.ObjectTypeMerchant:
				npc = append(npc, o.(*NPC))
			}
		}
	}
	return
}

func (e *Environ) MonstersProcess(...interface{}) {
	monsters, _ := e.GetActiveObjects()
	for i := range monsters {
		monsters[i].Process()
	}
}

func (e *Environ) NPCsProcess(...interface{}) {
	_, npcs := e.GetActiveObjects()
	for i := range npcs {
		npcs[i].Process()
	}
}
