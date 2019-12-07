package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/proto/client"
	"github.com/yenkeia/mirgo/proto/server"
)

var log = golog.New("server.handler")

func (g *Game) EventHandler(ev cellnet.Event) {

	switch msg := ev.Message().(type) {

	// 有新的连接
	case *cellnet.SessionAccepted:
		connected := server.Connected{}
		ev.Session().Send(&connected)

	// 有连接断开
	case *cellnet.SessionClosed:
		log.Debugln("session closed: ", ev.Session().ID())

	case *client.ClientVersion:
		clientVersion := server.ClientVersion{Result: 1}
		ev.Session().Send(&clientVersion)

	case *client.KeepAlive:
		keepAlive := server.KeepAlive{Time: 0}
		ev.Session().Send(keepAlive)

	default:
		log.Debugln(msg)
	}
}
