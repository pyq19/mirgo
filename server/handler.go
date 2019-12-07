package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/common"
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

	// TODO 保存新账号
	case *client.NewAccount:
		log.Debugln(msg.AccountID, msg.Password)
		ev.Session().Send(server.NewAccount{8})

	// TODO 登陆
	case *client.Login:
		res := new(server.LoginSuccess)

		c1 := new(common.SelectInfo)
		c1.Name = "测试登陆1"
		c1.Index = 1
		c1.Gender = common.Female
		c1.Class = common.Archer
		res.Characters = append(res.Characters, *c1)

		c2 := new(common.SelectInfo)
		c2.Name = "测试登陆2"
		c2.Index = 2
		c2.Gender = common.Female
		c2.Class = common.Taoist
		res.Characters = append(res.Characters, *c2)

		ev.Session().Send(res)

	// TODO 创建角色成功
	case *client.NewCharacter:
		res := new(server.NewCharacterSuccess)
		res.CharInfo.Index = 1
		res.CharInfo.Name = "创建角色成功了"
		ev.Session().Send(res)

	default:
		log.Debugln("default:", msg)
	}
}
