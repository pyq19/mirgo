package mir

import (
	"container/list"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"weshare/ut"

	"github.com/davyxu/cellnet"
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/mir/script"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/setting"
)

// Environ ...
type Environ struct {
	Game               *Game
	GameDB             *GameDB
	SessionIDPlayerMap *sync.Map    // map[int64]*Player
	Maps               map[int]*Map // mapID: Map
	ObjectID           uint32
	Players            []*Player
	lock               *sync.Mutex

	DefaultNPC *NPC

	msgMutex sync.Mutex
	MsgList  list.List
}

var env *Environ

func (e *Environ) PushMsg(f func()) {
	e.msgMutex.Lock()
	e.MsgList.PushBack(f)
	e.msgMutex.Unlock()
}

func (e *Environ) Loop() {

	fpsicker := time.NewTicker(time.Second / time.Duration(60))
	var lastFrame = time.Now()
	var now time.Time

	for {

		select {
		case <-fpsicker.C:
			e.msgMutex.Lock()
			if e.MsgList.Len() > 0 {
				for it := e.MsgList.Front(); it != nil; {
					curr := it
					it = it.Next()
					e.MsgList.Remove(curr).(func())()
				}
			}
			e.msgMutex.Unlock()

			now = time.Now()
			dt := now.Sub(lastFrame)
			lastFrame = now

			for _, m := range e.Maps {
				m.Frame(dt)
			}
		}
	}
}

// NewEnviron ...
func NewEnviron(g *Game) *Environ {
	env = new(Environ)
	env.Game = g

	script.SearchPaths = []string{
		filepath.Join(setting.Conf.EnvirPath, "NPCs"),
		setting.Conf.EnvirPath,
	}

	env.InitGameDB()
	env.InitMonsterDrop()
	env.InitMaps()

	env.DefaultNPC = NewNPC(nil, env.NewObjectID(), &common.NpcInfo{
		Name:     "DefaultNPC",
		Filename: "00Default",
	})

	env.ObjectID = 100000
	env.Players = make([]*Player, 0)
	env.lock = new(sync.Mutex)
	env.SessionIDPlayerMap = new(sync.Map)
	PrintEnviron(env)
	return env
}

func PrintEnviron(env *Environ) {
	mapCount := 0
	monsterCount := 0
	npcCount := 0
	for _, m := range env.Maps {
		mapCount++
		monsterCount += len(m.monsters)
		npcCount += len(m.npcs)
	}
	log.Debugf("共加载了 %d 张地图，%d 怪物，%d NPC\n", mapCount, monsterCount, npcCount)
}

// InitGameDB ...
func (e *Environ) InitGameDB() {
	gdb := new(GameDB)
	e.GameDB = gdb
	db := e.Game.DB

	db.Table("basic").First(&gdb.Basic)
	db.Table("game_shop_item").Find(&gdb.GameShopItems)
	db.Table("item").Find(&gdb.ItemInfos)
	db.Table("magic").Find(&gdb.MagicInfos)
	db.Table("map").Find(&gdb.MapInfos)
	db.Table("monster").Find(&gdb.MonsterInfos)
	db.Table("movement").Find(&gdb.MovementInfos)
	db.Table("npc").Find(&gdb.NpcInfos)
	db.Table("quest").Find(&gdb.QuestInfos)
	db.Table("respawn").Find(&gdb.RespawnInfos)
	db.Table("safe_zone").Find(&gdb.SafeZoneInfos)

	gdb.MapIDInfoMap = new(sync.Map)
	gdb.ItemIDInfoMap = new(sync.Map)
	gdb.ItemNameInfoMap = new(sync.Map)
	gdb.MonsterIDInfoMap = new(sync.Map)
	gdb.MonsterNameInfoMap = new(sync.Map)
	gdb.MagicIDInfoMap = new(sync.Map)
	for i := range gdb.MapInfos {
		v := gdb.MapInfos[i]
		gdb.MapIDInfoMap.Store(v.ID, v)
	}
	for i := range gdb.ItemInfos {
		v := gdb.ItemInfos[i]
		gdb.ItemIDInfoMap.Store(int(v.ID), v)
		gdb.ItemNameInfoMap.Store(v.Name, v)
	}
	for i := range gdb.MonsterInfos {
		v := gdb.MonsterInfos[i]
		gdb.MonsterNameInfoMap.Store(v.Name, v)
		gdb.MonsterIDInfoMap.Store(v.ID, v)
	}
	for i := range gdb.MagicInfos {
		v := gdb.MagicInfos[i]
		gdb.MagicIDInfoMap.Store(v.ID, v)
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

func (e *Environ) CreateDropItem(m *Map, userItem *common.UserItem, gold uint64) *Item {
	return &Item{
		MapObject: MapObject{
			ID:  e.NewObjectID(),
			Map: m,
		},
		Gold:     gold,
		UserItem: userItem,
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

	uppercaseNameRealNameMap := map[string]string{}
	files := ut.GetFiles(setting.Conf.MapDirPath, []string{".map"})

	for _, f := range files {
		uppercaseNameRealNameMap[strings.ToUpper(filepath.Base(f))] = f
	}

	e.Maps = map[int]*Map{}
	for _, mi := range e.GameDB.MapInfos {
		// FIXME 开发只加载部分地图
		if mi.ID != 1 && mi.ID != 384 {
			continue
		}

		m := LoadMap(uppercaseNameRealNameMap[strings.ToUpper(mi.Filename+".map")])
		mi.Filename = strings.ToUpper(mi.Filename)
		m.Info = mi
		if err := m.InitMonsters(); err != nil {
			panic(err)
		}
		if err := m.InitNPCs(); err != nil {
			panic(err)
		}
		e.Maps[mi.ID] = m
	}
}

func (e *Environ) NewObjectID() uint32 {
	return atomic.AddUint32(&e.ObjectID, 1)
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

func (e *Environ) GetPlayerByName(name string) *Player {
	e.lock.Lock()
	for i := 0; i < len(e.Players); i++ {
		o := e.Players[i]
		if name == o.Name {
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

func (e *Environ) GetMapByName(filename string) *Map {

	for _, m := range e.Maps {
		if m.Info.Filename == filename {
			return m
		}
	}
	return nil
}

func (e *Environ) GetMap(mapID int) *Map {
	v, ok := e.Maps[mapID]
	if !ok {
		return nil
	}
	return v
}

func (e *Environ) Broadcast(msg interface{}) {
	(*e.Game.Peer).(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
		ses.Send(msg)
		return true
	})
}

func (e *Environ) SystemBroadcast(...interface{}) {
	envPlayerCount := e.GetPlayersCount()
	text := "当前在线玩家人数: " + strconv.Itoa(envPlayerCount)
	e.Broadcast(&server.Chat{
		Message: text,
		Type:    common.ChatTypeSystem,
	})
}

func (e *Environ) Debug() {
	envPlayerCount := e.GetPlayersCount()
	nplayers := 0
	for _, m := range e.Maps {
		nplayers += len(m.GetAllPlayers())
	}
	if nplayers != envPlayerCount {
		log.Errorf("!!! warning envPlayerCount: %d != map allPlayer: %d\n", envPlayerCount, nplayers)
	} else {
		// log.Debugf("envPlayerCount: %d, map allPlayer: %d\n", envPlayerCount, len(allPlayer))
	}
}
