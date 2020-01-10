package main

import (
	"os"
)

func _GetMapForTest() *Map {
	gopath := os.Getenv("GOPATH")
	conf := Config{
		Addr:          "0.0.0.0:7000",
		DBPath:        gopath + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite",
		MapDirPath:    gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/",
		ScriptDirPath: gopath + "/src/github.com/yenkeia/mirgo/script/",
	}
	g := NewGame(conf)
	v, _ := g.Env.Maps.Load(1)
	m := v.(*Map)
	return m
}
