package main

import (
	"os"
)

func _GetMapForTest() *Map {
	gopath := os.Getenv("GOPATH")
	var addr = "0.0.0.0:7000"
	var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	conf := Config{addr, gopath + mirDB}
	g := NewGame(conf)
	v, _ := g.Env.Maps.Load(1)
	m := v.(*Map)
	return m
}
