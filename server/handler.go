package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/codec/mircodec"
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
		c1.Gender = common.MirGenderFemale
		c1.Class = common.MirClassArcher
		res.Characters = append(res.Characters, *c1)

		c2 := new(common.SelectInfo)
		c2.Name = "测试登陆2"
		c2.Index = 2
		c2.Gender = common.MirGenderFemale
		c2.Class = common.MirClassAssassin
		res.Characters = append(res.Characters, *c2)

		ev.Session().Send(res)

	// TODO 创建角色成功
	case *client.NewCharacter:
		log.Debugln(msg.Name, msg.Class, msg.Gender)
		res := new(server.NewCharacterSuccess)
		res.CharInfo.Index = 0
		res.CharInfo.Name = msg.Name
		res.CharInfo.Class = msg.Class
		res.CharInfo.Gender = msg.Gender
		ev.Session().Send(res)

	// TODO
	case *client.StartGame:
		// SetConcentration
		sc := new(server.SetConcentration)
		sc.ObjectID = 66432
		sc.Enabled = false
		sc.Interrupted = false
		ev.Session().Send(sc)

		// StartGame
		sg := new(server.StartGame)
		sg.Result = 4
		sg.Resolution = 1024
		ev.Session().Send(sg)

		// MapInformation
		bytes := []byte{
			1, 48,
			14, 66, 105, 99, 104, 111, 110, 80, 114, 111, 118, 105, 110, 99, 101,
			101, 0,
			135, 0,
			0, 0,
			0,
			0,
			0}
		r := new(server.MapInformation)
		codec := new(mircodec.MirCodec)
		if err := codec.Decode(bytes, r); err != nil {
			panic(err)
		}
		ev.Session().Send(r)

		// NewItemInfo
		bytes1 := []byte{146, 2, 0, 0, 13, 40, 72, 80, 41, 68, 114, 117, 103, 83, 109, 97, 108, 108, 13, 0, 0, 31, 3, 0, 0, 0, 1, 0, 0, 142, 1, 0, 0, 20, 0, 0, 0, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 1, 0}
		item1 := new(server.NewItemInfo)
		codec.Decode(bytes1, item1)
		ev.Session().Send(item1)

		bytes2 := []byte{235, 4, 0, 0, 16, 84, 101, 115, 116, 83, 101, 114, 118, 101, 114, 83, 99, 114, 111, 108, 108, 21, 4, 0, 31, 3, 0, 1, 0, 1, 0, 0, 254, 6, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 23, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 66, 97, 115, 105, 99, 32, 84, 101, 115, 116, 32, 83, 101, 114, 118, 101, 114, 32, 83, 99, 114, 111, 108, 108, 32, 119, 104, 105, 99, 104, 32, 103, 105, 118, 101, 115, 32, 105, 110, 102, 111, 114, 109, 97, 116, 105, 111, 110, 46}
		item2 := new(server.NewItemInfo)
		codec.Decode(bytes2, item2)
		ev.Session().Send(item2)

		bytes3 := []byte{221, 0, 0, 0, 11, 87, 111, 111, 100, 101, 110, 83, 119, 111, 114, 100, 1, 1, 0, 7, 3, 0, 0, 0, 4, 0, 1, 30, 0, 160, 15, 1, 0, 0, 0, 50, 0, 0, 0, 0, 0, 0, 0, 2, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0}
		item3 := new(server.NewItemInfo)
		codec.Decode(bytes3, item3)
		ev.Session().Send(item3)

		bytes4 := []byte{61, 1, 0, 0, 12, 66, 97, 115, 101, 68, 114, 101, 115, 115, 40, 77, 41, 2, 1, 0, 31, 1, 0, 1, 0, 5, 0, 1, 60, 0, 136, 19, 1, 0, 0, 0, 120, 0, 0, 0, 2, 2, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1, 1, 0}
		item4 := new(server.NewItemInfo)
		codec.Decode(bytes4, item4)
		ev.Session().Send(item4)

		bytes5 := []byte{210, 2, 0, 0, 6, 67, 97, 110, 100, 108, 101, 12, 0, 0, 31, 3, 0, 0, 0, 1, 38, 0, 130, 0, 64, 31, 1, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 1, 0}
		item5 := new(server.NewItemInfo)
		codec.Decode(bytes5, item5)
		ev.Session().Send(item5)

		uiCodec := new(mircodec.MirUserInformationCodec)
		bytes6 := []byte{128, 3, 1, 0, 1, 0, 0, 0, 6, 99, 99, 99, 99, 99, 99, 0, 0, 255, 255, 255, 255, 1, 1, 1, 0, 28, 1, 0, 0, 96, 2, 0, 0, 1, 8, 15, 0, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 1, 46, 0, 0, 0, 1, 3, 0, 0, 0, 0, 0, 0, 0, 146, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 54, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 14, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 62, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		item6 := new(server.UserInformation)
		uiCodec.Decode(bytes6, item6)
		ev.Session().Send(item6)

	default:
		log.Debugln("default:", msg)
	}
}
