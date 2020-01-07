package main

import "os"

type Config struct {
	Addr          string
	DBPath        string
	MapDirPath    string
	ScriptDirPath string
}

func main() {
	gopath := os.Getenv("GOPATH")
	g := NewGame(Config{
		Addr:          "0.0.0.0:7000",
		DBPath:        gopath + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite",
		MapDirPath:    gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/",
		ScriptDirPath: gopath + "/src/github.com/yenkeia/mirgo/script/",
	})
	g.ServerStart()
}
