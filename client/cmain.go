package main

import (
	"bufio"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/proto/client"
	"github.com/yenkeia/mirgo/proto/server"
	"os"
	"strconv"
	"strings"
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
	case "CLIENT_VERSION":
		return &client.ClientVersion{VersionHash: []uint8{16, 0, 0, 0, 86, 92, 129, 20, 102, 64, 159, 148, 125, 97, 112, 85, 237, 250, 133, 162}}
	case "KEEP_ALIVE":
		return &client.KeepAlive{Time: 0}
	case "LOGIN": // 5
		return &client.Login{
			AccountID: args[1],
			Password:  args[2],
		}
	case "NEW_CHARACTER":
		return &client.NewCharacter{
			Name:   args[1],
			Gender: common.MirGenderMale,
			Class:  common.MirClassTaoist,
		}
	case "START_GAME": // 8
		i, _ := strconv.Atoi(args[1])
		return &client.StartGame{CharacterIndex: int16(i)}
	case "DELETE_CHARACTER":
		i, _ := strconv.Atoi(args[1])
		return &client.DeleteCharacter{CharacterIndex: int16(i)}
	default:
		return nil
	}
}

func main() {
	queue := cellnet.NewEventQueue()
	//addr := "192.168.31.242:7000"
	addr := "127.0.0.1:7000"
	p := peer.NewGenericPeer("tcp.Connector", "client", addr, queue)
	proc.BindProcessorHandler(p, "mir.client.tcp", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			//log.Debugln("client connected")
		case *cellnet.SessionClosed:
			//log.Debugln("client error")
		case *server.Connected:
			//log.Infof("<--- server.Connected")
		case *server.ClientVersion:
			//log.Infof("<--- server.ClientVersion")
		case *server.KeepAlive:
			//log.Infof("<--- server.KeepAlive")
		case *server.Chat:
			log.Infoln("<--- server.Chat. Type:", msg.Type, ", Text: ", msg.Message)
		default:
			log.Debugf("default: 客户端收到: %s", msg)
		}
	})
	p.Start()
	queue.StartLoop()
	log.Debugln("Ready to chat!")
	readConsole(func(str string) {
		session := p.(interface{ Session() cellnet.Session }).Session()
		msg := parse(str)
		if msg == nil {
			return
		}
		session.Send(msg)
	})
}
