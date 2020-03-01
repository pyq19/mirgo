package mir

import (
	"time"

	"github.com/yenkeia/mirgo/common"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/timer"
	"github.com/davyxu/golog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
)

var log = golog.New("server.game")

// Game ...
type Game struct {
	DB   *gorm.DB
	Env  *Environ
	Peer *cellnet.GenericPeer
}

// NewGame ...
func NewGame() *Game {
	g := new(Game)
	gameData, _ := gorm.Open("sqlite3", settings.DBPath)
	accountData, _ := gorm.Open("sqlite3", settings.AccountDBPath)
	adb = &DB{db: accountData}
	adb.db.AutoMigrate(
		&common.Account{},
		&common.AccountCharacter{},
		&common.Character{},
		&common.CharacterUserItem{},
		&common.UserItem{},
		&common.UserMagic{},
	)

	g.DB = gameData
	g.Env = NewEnviron(g)

	return g
}

// ServerStart ...
func (g *Game) ServerStart() {

	// 这里用cellnet 单线程模式。消息处理都在queue线程。无需再另开线程
	queue := cellnet.NewEventQueue()
	p := peer.NewGenericPeer("tcp.Acceptor", "server", settings.Addr, queue)
	g.Peer = &p
	proc.BindProcessorHandler(p, "mir.server.tcp", g.HandleEvent)

	timer.NewLoop(queue, time.Second/time.Duration(60), func(*timer.Loop) {
		g.Env.Loop()
	}, nil).Start()

	p.Start()         // 开始侦听
	queue.StartLoop() // 事件队列开始循环
	queue.Wait()      // 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
}
