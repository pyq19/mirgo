package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
)

// Game ...
type Game struct {
	Conf Config
	DB   *gorm.DB
	Pool *Pool
	Env  *Environ
	Peer *cellnet.GenericPeer
}

// NewGame ...
func NewGame(conf Config) *Game {
	g := new(Game)
	g.Conf = conf
	db, err := gorm.Open("sqlite3", conf.MirDB)
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()
	g.DB = db
	g.Pool = NewPool(10)
	g.Env = NewEnviron(g)
	g.Env.StartLoop()
	return g
}

// ServerStart ...
func (g *Game) ServerStart() {
	queue := cellnet.NewEventQueue()
	p := peer.NewGenericPeer("tcp.Acceptor", "server", g.Conf.Addr, queue)
	g.Peer = &p
	proc.BindProcessorHandler(p, "mir.server.tcp", g.HandleEvent)
	p.Start()         // 开始侦听
	queue.StartLoop() // 事件队列开始循环
	queue.Wait()      // 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
}
