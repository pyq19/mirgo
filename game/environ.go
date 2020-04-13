package game

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/timer"
	"github.com/davyxu/golog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yenkeia/mirgo/game/cm"
	_ "github.com/yenkeia/mirgo/game/mirtcp"
	"github.com/yenkeia/mirgo/game/proto/server"
	"github.com/yenkeia/mirgo/game/script"
	"github.com/yenkeia/mirgo/game/util"
	"github.com/yenkeia/mirgo/setting"
)

var env *Environ
var log = golog.New("server")

// Environ ...
type Environ struct {
	Game               *Game
	Peer               cellnet.GenericPeer
	SessionIDPlayerMap *sync.Map    // map[int64]*Player
	Maps               map[int]*Map // mapID: Map
	ObjectID           uint32
	ObjectIDChan       chan uint32
	Players            *PlayerList
	lastFrame          time.Time

	DefaultNPC *NPC
	Lights     cm.LightSetting
	GuildList  *GuildList
}

func (e *Environ) Loop() {
	now := time.Now()
	dt := now.Sub(e.lastFrame)

	lastHour := e.lastFrame.Hour()
	nowHour := now.Hour()
	if lastHour != nowHour {
		e.AdjustLights()
		e.Broadcast(&server.TimeOfDay{Lights: e.Lights})
	}

	for _, m := range e.Maps {
		m.Frame(dt)
	}
	e.lastFrame = now
}

// ServerStart ...
func (g *Environ) ServerStart() {

	// 这里用cellnet 单线程模式。消息处理都在queue线程。无需再另开线程
	queue := cellnet.NewEventQueue()

	acceptor := "tcp.Acceptor"
	if settings.Acceptor == "websocket" {
		acceptor = "gorillaws.Acceptor"
	} else {
		settings.Acceptor = "tcp"
	}

	p := peer.NewGenericPeer(acceptor, "server", settings.Addr, queue)
	proc.BindProcessorHandler(p, "mir.server."+settings.Acceptor, g.Game.HandleEvent)

	timer.NewLoop(queue, time.Second/time.Duration(60), func(*timer.Loop) {
		env.Loop()
	}, nil).Start()

	env.Peer = p

	p.Start()         // 开始侦听
	queue.StartLoop() // 事件队列开始循环
	queue.Wait()      // 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
}

// NewEnviron ...
func NewEnviron() *Environ {
	settings = setting.Must()

	gameData, err := gorm.Open("sqlite3", settings.DBPath)
	if err != nil {
		panic(err)
	}
	data = NewGameData(gameData)

	accountData, err := gorm.Open("sqlite3", settings.AccountDBPath)
	if err != nil {
		panic(err)
	}
	adb = NewAccountDB(accountData)

	e := new(Environ)
	env = e

	e.lastFrame = time.Now()

	e.ObjectIDChan = make(chan uint32, 100000)
	id := adb.GetObjectID()
	if id == 0 {
		id = 100000
		adb.Table("basic").Create(&cm.Basic{ID: 1, ObjectID: id})
	}
	e.ObjectID = id
	go func() {
		for id := range e.ObjectIDChan {
			adb.SyncObjectID(id)
		}
	}()

	script.SearchPaths = []string{
		filepath.Join(settings.EnvirPath, "NPCs"),
		settings.EnvirPath,
	}

	e.DefaultNPC = NewNPC(nil, e.NewObjectID(), &cm.NpcInfo{
		Name:     "DefaultNPC",
		Filename: "00Default",
	})

	e.InitMaps()

	e.Players = NewPlayerList()
	e.SessionIDPlayerMap = new(sync.Map)
	e.Game = &Game{}

	e.AdjustLights()

	PrintEnviron(e)

	return env
}

func PrintEnviron(env *Environ) {
	banner := `
• ▌ ▄ ·. ▪  ▄▄▄   ▄▄ •       
·██ ▐███▪██ ▀▄ █·▐█ ▀ ▪▪     
▐█ ▌▐▌▐█·▐█·▐▀▀▄ ▄█ ▀█▄ ▄█▀▄ 
██ ██▌▐█▌▐█▌▐█•█▌▐█▄▪▐█▐█▌.▐▌
▀▀  █▪▀▀▀▀▀▀.▀  ▀·▀▀▀▀  ▀█▄▀▪
`
	mapCount := 0
	monsterCount := 0
	npcCount := 0
	for _, m := range env.Maps {
		mapCount++
		monsterCount += len(m.monsters)
		npcCount += len(m.npcs)
	}
	fmt.Println(banner)
	log.Debugf("共加载了 %d 张地图，%d 怪物，%d NPC\n", mapCount, monsterCount, npcCount)
}

func (e *Environ) NewUserItem(i *cm.ItemInfo) *cm.UserItem {
	res := &cm.UserItem{
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
		Info:           i,
	}
	return res
}

// InitMaps ...
func (e *Environ) InitMaps() {

	uppercaseNameRealNameMap := map[string]string{}
	files := util.GetFiles(settings.MapDirPath, []string{".map"})

	for _, f := range files {
		uppercaseNameRealNameMap[strings.ToUpper(filepath.Base(f))] = f
	}

	// FIXME 开发只加载部分地图
	allowarr := []int{1, 2, 3, 4, 6, 7, 8, 10, 11, 12, 13, 15, 16, 17, 18, 19, 20, 21, 22, 24, 26, 27, 28, 29, 30, 31, 32, 25, 144, 384}
	allow := map[int]bool{}
	for _, v := range allowarr {
		allow[v] = true
	}

	e.Maps = map[int]*Map{}
	for _, mi := range data.MapInfos {

		if _, ok := allow[mi.ID]; !ok {
			continue
		}

		m := LoadMap(uppercaseNameRealNameMap[strings.ToUpper(mi.Filename+".map")])
		mi.Filename = strings.ToUpper(mi.Filename)
		m.Info = mi
		if err := m.InitAll(); err != nil {
			panic(err)
		}

		e.Maps[mi.ID] = m
	}
}

func (e *Environ) NewObjectID() uint32 {
	id := atomic.AddUint32(&e.ObjectID, 1)
	e.ObjectIDChan <- id
	return id
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
	e.Peer.(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
		ses.Send(msg)
		return true
	})
}

func (e *Environ) AdjustLights() {
	oldLights := e.Lights
	hours := (time.Now().Hour() * 2) % 24
	if hours == 6 || hours == 7 {
		e.Lights = cm.LightSettingDawn
	} else if hours >= 8 && hours <= 15 {
		e.Lights = cm.LightSettingDay
	} else if hours == 16 || hours == 17 {
		e.Lights = cm.LightSettingEvening
	} else {
		// e.Lights = cm.LightSettingNight
		e.Lights = cm.LightSettingEvening
	}
	if oldLights == e.Lights {
		return
	}
	// e.Broadcast(&server.TimeOfDay{Lights: e.Lights})
}

// GetGuild 通过工会名称获取工会
func (e *Environ) GetGuild(name string) *Guild {
	for e := e.GuildList.List.Front(); e != nil; e = e.Next() {
		guild := e.Value.(*Guild)
		if guild.Name == name {
			return guild
		}
	}
	return nil
}

func (e *Environ) GetPlayer(name string) *Player {
	return e.Players.GetPlayerByName(name)
}
