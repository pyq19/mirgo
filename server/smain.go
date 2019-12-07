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
	conf := &Config{addr, gopath + mirDB}
	g := NewGame(*conf)
	g.ServerStart()
}
