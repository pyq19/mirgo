package main

import (
	"github.com/yenkeia/mirgo/server"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
)

func main() {
	server.NewGame().ServerStart()
}
