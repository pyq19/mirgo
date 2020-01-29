package main

import (
	"github.com/davyxu/cellnet"
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/setting"
	"os"
	"strconv"
	"strings"
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
	env.InitMonsterDrop()
	env.InitMaps()
	env.ObjectID = 100000
	env.Players = make([]*Player, 0)
	env.lock = new(sync.Mutex)
	err := env.InitObjects()
	if err != nil {
		panic(err)
	}
	env.SessionIDPlayerMap = new(sync.Map)
	PrintEnviron(env)
	return
}

func PrintEnviron(env *Environ) {
	mapCount := 0
	monsterCount := 0
	npcCount := 0
	env.Maps.Range(func(k, v interface{}) bool {
		mapCount++
		m := v.(*Map)
		m.AOI.grids.Range(func(k, v interface{}) bool {
			g := v.(*Grid)
			objs := g.GetAllObjects()
			for i := range objs {
				o := objs[i]
				switch o.GetRace() {
				case common.ObjectTypeMonster:
					monsterCount++
				case common.ObjectTypeMerchant:
					npcCount++
				}
			}
			return true
		})
		return true
	})
	log.Debugf("共加载了 %d 张地图，%d 怪物，%d NPC\n", mapCount, monsterCount, npcCount)
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
	gdb.MapIDInfoMap = new(sync.Map)
	gdb.ItemIDInfoMap = new(sync.Map)
	gdb.ItemNameInfoMap = new(sync.Map)
	gdb.MonsterIDInfoMap = new(sync.Map)
	for i := range gdb.MapInfos {
		v := gdb.MapInfos[i]
		gdb.MapIDInfoMap.Store(v.ID, &v)
	}
	for i := range gdb.ItemInfos {
		v := gdb.ItemInfos[i]
		gdb.ItemIDInfoMap.Store(int(v.ID), &v)
		gdb.ItemNameInfoMap.Store(v.Name, &v)
	}
	for i := range gdb.MonsterInfos {
		v := gdb.MonsterInfos[i]
		gdb.MonsterIDInfoMap.Store(v.ID, &v)
	}
}

func (e *Environ) InitMonsterDrop() {
	gdb := e.GameDB
	itemMap := make(map[string]int32)
	for i := range gdb.ItemInfos {
		v := gdb.ItemInfos[i]
		itemMap[v.Name] = v.ID
	}
	gdb.DropInfoMap = new(sync.Map)
	for i := range gdb.MonsterInfos {
		v := gdb.MonsterInfos[i]
		dropInfos, err := common.GetDropInfosByMonsterName(setting.Conf.DropDirPath, v.Name)
		if err != nil {
			log.Warnln("加载怪物掉落错误", v.Name, err.Error())
			continue
		}
		gdb.DropInfoMap.Store(v.Name, dropInfos)
	}
}

func (e *Environ) NewUserItem(i *common.ItemInfo) *common.UserItem {
	res := &common.UserItem{
		ID:             uint64(e.NewObjectID()),
		ItemID:         i.ID,
		CurrentDura:    100,
		MaxDura:        100,
		Count:          1,
		AC:             i.MinAC,
		MAC:            i.MinMAC,
		DC:             i.MinDC,
		MC:             i.MinMC,
		SC:             i.MinSC,
		Accuracy:       i.Accuracy,
		Agility:        i.Agility,
		HP:             0,
		MP:             0,
		AttackSpeed:    i.AttackSpeed,
		Luck:           i.Luck,
		SoulBoundId:    0,
		Bools:          0,
		Strong:         0,
		MagicResist:    0,
		PoisonResist:   0,
		HealthRecovery: 0,
		ManaRecovery:   0,
		PoisonRecovery: 0,
		CriticalRate:   0,
		CriticalDamage: 0,
		Freezing:       0,
		PoisonAttack:   0,
	}
	return res
}

// InitMaps ...
func (e *Environ) InitMaps() {
	mapDirPath := setting.Conf.MapDirPath
	uppercaseNameRealNameMap := make(map[string]string) // 目录下的文件名大写与该文件的真实文件名对应关系
	f, err := os.OpenFile(mapDirPath, os.O_RDONLY, os.ModeDir)
	if err != nil {
		panic(err)
	}
	fileInfo, _ := f.Readdir(-1)
	for _, info := range fileInfo {
		if !info.IsDir() {
			uppercaseNameRealNameMap[strings.ToUpper(info.Name())] = info.Name()
		}
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	// FIXME get map v2 v3 ??
	skipMap := map[string]bool{
		"R05":   true,
		"R07":   true,
		"R08":   true,
		"R10":   true,
		"R11":   true,
		"EM000": true,
		"EM001": true,
		"EM002": true,
		"EM003": true,
	}
	//e.Maps = make([]Map, 386)
	e.Maps = new(sync.Map)
	for i := range e.GameDB.MapInfos {
		mi := e.GameDB.MapInfos[i]
		if skipMap[strings.ToUpper(mi.Filename)] {
			continue
		}
		// FIXME 开发只加载第一张地图
		if mi.ID != 1 {
			continue
		}
		m := GetMapV1(GetMapBytes(mapDirPath + uppercaseNameRealNameMap[strings.ToUpper(mi.Filename+".map")]))
		m.Env = e
		m.Info = &mi
		e.Maps.Store(mi.ID, m)
		break
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
		if o == nil || o.ID == 0 {
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
		// log.Debugf("envPlayerCount: %d, map allPlayer: %d\n", envPlayerCount, len(allPlayer))
	}
}

// TODO 待优化
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
