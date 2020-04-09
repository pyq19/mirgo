package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/game/cm"
	_ "github.com/yenkeia/mirgo/game/mircodec"
	_ "github.com/yenkeia/mirgo/game/mirtcp"
	"github.com/yenkeia/mirgo/game/proto/client"
	"github.com/yenkeia/mirgo/game/proto/server"
)

var log = golog.New("client")

func readConsole(callback func(string)) {
	for {
		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)
		callback(text)
	}
}

func parse(str string) interface{} {
	args := strings.Split(str, " ")
	switch args[0] {
	case "ClientVersion", "CLIENT_VERSION":
		return &client.ClientVersion{VersionHash: []uint8{16, 0, 0, 0, 86, 92, 129, 20, 102, 64, 159, 148, 125, 97, 112, 85, 237, 250, 133, 162}}
	case "KeepAlive", "KEEP_ALIVE":
		return &client.KeepAlive{Time: 0}
	case "NewAccount":
		return &client.NewAccount{
			AccountID:      args[1],
			Password:       args[2],
			DateTime:       0,
			UserName:       "",
			SecretQuestion: "",
			SecretAnswer:   "",
			EMailAddress:   "",
		}
	case "ChangePassword":
	case "Login", "login": // 5
		return &client.Login{
			AccountID: args[1],
			Password:  args[2],
		}
	case "NewCharacter":
		return &client.NewCharacter{
			Name:   args[1],
			Gender: cm.MirGenderMale,
			Class:  cm.MirClassTaoist,
		}
	case "DeleteCharacter":
		i, _ := strconv.Atoi(args[1])
		return &client.DeleteCharacter{CharacterIndex: int32(i)}
	case "StartGame", "start", "startgame": // 8
		i, _ := strconv.Atoi(args[1])
		return &client.StartGame{CharacterIndex: int16(i)}
	case "LogOut":
	case "Turn", "turn":
		return &client.Turn{Direction: direction(args[1])}
	case "Walk", "walk":
		return &client.Walk{Direction: direction(args[1])}
	case "Run", "run":
		return &client.Run{Direction: direction(args[1])}
	case "Chat", "chat":
		return &client.Chat{Message: args[1]}
	case "MoveItem":
	case "StoreItem":
	case "DepositRefineItem":
	case "RetrieveRefineItem":
	case "RefineCancel":
	case "RefineItem":
	case "CheckRefine":
	case "ReplaceWedRing":
	case "DepositTradeItem":
	case "RetrieveTradeItem":
	case "TakeBackItem":
	case "MergeItem":
	case "EquipItem":
	case "RemoveItem":
	case "RemoveSlotItem":
	case "SplitItem":
	case "UseItem":
	case "DropItem":
	case "DropGold":
	case "PickUp":
	case "Inspect":
	case "ChangeAMode":
	case "ChangePMode":
	case "ChangeTrade":
	case "Attack":
		return &client.Attack{
			Direction: direction(args[1]),
			Spell:     cm.SpellNone,
		}
	case "RangeAttack":
	case "Harvest":
	case "CallNPC":
	case "TalkMonsterNPC":
	case "BuyItem":
	default:
		return nil
	}
	return nil
}

func direction(s string) cm.MirDirection {
	switch s {
	case "up":
		return cm.MirDirectionUp
	case "upright":
		return cm.MirDirectionUpRight
	case "right":
		return cm.MirDirectionRight
	case "downright":
		return cm.MirDirectionDownRight
	case "down":
		return cm.MirDirectionDown
	case "downleft":
		return cm.MirDirectionDownLeft
	case "left":
		return cm.MirDirectionLeft
	case "upleft":
		return cm.MirDirectionUpLeft
	default:
		return cm.MirDirectionUp
	}
}
func directionIcon(d cm.MirDirection) string {
	switch d {
	case cm.MirDirectionUp:
		return "↑↑↑"
	case cm.MirDirectionUpRight:
		return "↗↗↗"
	case cm.MirDirectionRight:
		return "→→→"
	case cm.MirDirectionDownRight:
		return "↘↘↘"
	case cm.MirDirectionDown:
		return "↓↓↓"
	case cm.MirDirectionDownLeft:
		return "↙↙↙"
	case cm.MirDirectionLeft:
		return "←←←"
	case cm.MirDirectionUpLeft:
		return "↖↖↖"
	}
	return fmt.Sprintf("!!! error direction: %d", d)
}

func main() {
	queue := cellnet.NewEventQueue()
	//addr := "192.168.31.242:7000"
	addr := "127.0.0.1:7000"
	p := peer.NewGenericPeer("tcp.Connector", "client", addr, queue)
	session := p.(interface{ Session() cellnet.Session }).Session()
	proc.BindProcessorHandler(p, "mir.client.tcp", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			//log.Debugln("client connected")
		case *cellnet.SessionClosed:
			//log.Debugln("client error")
		case *server.Connected:
			//log.Infof("<--- server.Connected")
			session.Send(&client.ClientVersion{VersionHash: []uint8{16, 0, 0, 0, 86, 92, 129, 20, 102, 64, 159, 148, 125, 97, 112, 85, 237, 250, 133, 162}})
		case *server.ClientVersion:
			//log.Infof("<--- server.ClientVersion")
		case *server.KeepAlive:
			//log.Infof("<--- server.KeepAlive")
			session.Send(&client.KeepAlive{Time: 0})
		case *server.Chat:
		case *server.ObjectTurn:
			log.Infof("<--- server.ObjectTurn %d, %s\n", msg.ObjectID, directionIcon(msg.Direction))
		case *server.ObjectWalk:
			log.Infof("<--- server.ObjectWalk %d, %s\n", msg.ObjectID, directionIcon(msg.Direction))
		case *server.ObjectRun:
			log.Infof("<--- server.ObjectRun %d, %s\n", msg.ObjectID, directionIcon(msg.Direction))
		default:
			//log.Debugf("default: 客户端收到: %s", msg)
			_ = msg
		}
	})
	p.Start()
	queue.StartLoop()
	log.Debugln("Ready to chat!")
	readConsole(func(str string) {
		msg := parse(str)
		if msg == nil {
			return
		}
		session.Send(msg)
	})
}
