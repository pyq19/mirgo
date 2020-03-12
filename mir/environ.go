package mir

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/timer"
	"github.com/davyxu/golog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/mir/script"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/setting"
	"github.com/yenkeia/mirgo/ut"
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
	Players            []*Player
	lock               *sync.Mutex
	lastFrame          time.Time

	DefaultNPC *NPC
}

func (e *Environ) Loop() {
	now := time.Now()
	dt := now.Sub(e.lastFrame)
	e.lastFrame = now

	for _, m := range e.Maps {
		m.Frame(dt)
	}
}

// ServerStart ...
func (g *Environ) ServerStart() {

	// 这里用cellnet 单线程模式。消息处理都在queue线程。无需再另开线程
	queue := cellnet.NewEventQueue()
	p := peer.NewGenericPeer("tcp.Acceptor", "server", settings.Addr, queue)
	proc.BindProcessorHandler(p, "mir.server.tcp", g.Game.HandleEvent)

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

	gameData, _ := gorm.Open("sqlite3", settings.DBPath)
	data = NewGameData(gameData)

	accountData, _ := gorm.Open("sqlite3", settings.AccountDBPath)
	adb = NewAccountDB(accountData)

	e := new(Environ)
	env = e

	e.lastFrame = time.Now()

	script.SearchPaths = []string{
		filepath.Join(settings.EnvirPath, "NPCs"),
		settings.EnvirPath,
	}

	e.DefaultNPC = NewNPC(nil, e.NewObjectID(), &common.NpcInfo{
		Name:     "DefaultNPC",
		Filename: "00Default",
	})

	e.InitMaps()

	e.ObjectID = 100000
	e.Players = make([]*Player, 0)
	e.lock = new(sync.Mutex)
	e.SessionIDPlayerMap = new(sync.Map)
	e.Game = &Game{}

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
		Info:           i,
	}
	return res
}

// InitMaps ...
func (e *Environ) InitMaps() {

	uppercaseNameRealNameMap := map[string]string{}
	files := ut.GetFiles(settings.MapDirPath, []string{".map"})

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
	e.Peer.(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {
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
