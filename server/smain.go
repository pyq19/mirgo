package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/proto/client"
	"github.com/yenkeia/mirgo/proto/server"

	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

var log = golog.New("server")

func main() {

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Acceptor", "server", "127.0.0.1:7000", queue)

	proc.BindProcessorHandler(p, "mir.server.tcp", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {

		// 有新的连接
		case *cellnet.SessionAccepted:
			connected := server.Connected{}
			ev.Session().Send(&connected)

		// 有连接断开
		case *cellnet.SessionClosed:
			log.Debugln("session closed: ", ev.Session().ID())

		case *client.ClientVersion:
			clientVersion := server.ClientVersion{Result: 1} // TODO
			ev.Session().Send(&clientVersion)

		case *client.KeepAlive:
			keepAlive := server.KeepAlive{Time: 1000} // TODO
			ev.Session().Send(keepAlive)

		default:
			log.Debugln(msg)
		}
	})

	// 开始侦听
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	queue.Wait()

}
