package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
)

var log = golog.New("server")

func main() {

	queue := cellnet.NewEventQueue()

	addr := "0.0.0.0:7000"
	p := peer.NewGenericPeer("tcp.Acceptor", "server", addr, queue)

	proc.BindProcessorHandler(p, "mir.server.tcp", eventHandler)

	// 开始侦听
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	queue.Wait()

}
