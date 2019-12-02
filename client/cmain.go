package main

import (
	"bufio"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/proto/client"
	"github.com/yenkeia/mirgo/proto/server"
	"os"
	"strconv"
	"strings"

	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
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

func main() {

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "client", "127.0.0.1:7000", queue)

	proc.BindProcessorHandler(p, "mir.client.tcp", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			log.Debugln("client connected")
		case *cellnet.SessionClosed:
			log.Debugln("client error")
		case *server.Connected:
			log.Infof("<--- server.Connected")
		case *server.ClientVersion:
			log.Infof("<--- server.ClientVersion")
		case *server.KeepAlive:
			log.Infof("<--- server.KeepAlive")
		default:
			log.Debugln(msg)
		}
	})

	p.Start()

	queue.StartLoop()

	log.Debugln("Ready to chat!")

	readConsole(func(str string) {
		session := p.(interface{ Session() cellnet.Session }).Session()

		id, _ := strconv.Atoi(str)
		id = id + 1000
		idStr := strconv.Itoa(id)
		switch id {
		case client.CLIENT_VERSION:
			log.Debugln(idStr + " CLIENT_VERSION")
			session.Send(&client.ClientVersion{
				VersionHash: []uint8{}, // TODO
			})
		case client.KEEP_ALIVE:
			log.Debugln(idStr + " KEEP_ALIVE")
			session.Send(&client.KeepAlive{
				Time: 1000, // TODO
			})
		default:
			log.Debugln(idStr + " default")
		}
	})

}
