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

		// UserInformation
		bytes = []byte{
			128, 3, 1, 0, // ObjectID
			1, 0, 0, 0, // RealId
			10, 104, 101, 108, 108, 111, 119, 111, 114, 108, 100,
			0,                  // GuildName
			0,                  // GuildRank
			255, 255, 255, 255, // Color	uint32	[4]byte{225, 225, 225, 225}
			0,    // MirClass
			0,    // MirGender
			1, 0, // Level ushort uint16
			28, 1, 0, 0, // Point.X uint32
			97, 2, 0, 0, // Point.Y uint32
			6,
			5,
			17, 0,
			14, 0,
			0, 0, 0, 0, 0, 0, 0, 0, // Experience long int42
			100, 0, 0, 0, 0, 0, 0, 0, // MaxExperience
			0, // LevelEffects
			1, 46,
			0, 0, 0, 1, 3, 0, 0, 0,
			0, 0, 0, 0,
			146, 2,
			0, 0,
			0, 0, 0, 0,
			1, 0, 0, 0,
			0, // AC
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,                  // Luck
			255, 255, 255, 255, // SoulBoundId
			0,
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 1, 6, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 14, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 61, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 1, 4, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 49, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 0}
		res := new(server.UserInformation)
		res.Bytes = bytes
		ev.Session().Send(res)

	default:
		log.Debugln("default:", msg)
	}
}
