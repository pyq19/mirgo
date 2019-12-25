package main

import "os"

var addr = "0.0.0.0:7000"
var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"

type Config struct {
	Addr  string
	MirDB string
}

func main() {
	gopath := os.Getenv("GOPATH")
	conf := Config{addr, gopath + mirDB}
	g := NewGame(conf)

	/*
		v, _ := g.Env.Maps.Load(1)
		count := 0
		v.(*Map).Cells.Range(func (k, m interface{}) bool {
			c := m.(*Cell)
			if c.Object != nil {
				count += 1
			}
			return true
		})
		log.Debugln("!!!", count)
	*/

	g.ServerStart()
}
