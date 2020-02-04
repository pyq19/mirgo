package main

import (
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/server"
)

func main() {
	server.NewGame().ServerStart()
}
