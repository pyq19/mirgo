package mir

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/setting"
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
	db, err := gorm.Open("sqlite3", setting.Conf.DBPath)
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()
	g.DB = db
	g.Env = NewEnviron(g)
	go g.Env.Loop()

	return g
}

// ServerStart ...
func (g *Game) ServerStart() {
	queue := cellnet.NewEventQueue()
	p := peer.NewGenericPeer("tcp.Acceptor", "server", setting.Conf.Addr, queue)
	g.Peer = &p
	proc.BindProcessorHandler(p, "mir.server.tcp", g.HandleEvent)
	p.Start()         // 开始侦听
	queue.StartLoop() // 事件队列开始循环
	queue.Wait()      // 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
}
