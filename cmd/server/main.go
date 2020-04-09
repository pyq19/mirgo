package main

import (
	"github.com/yenkeia/mirgo/game"
	_ "github.com/yenkeia/mirgo/game/behavior"
	_ "github.com/yenkeia/mirgo/game/mircodec"
)

func main() {
	game.NewEnviron().ServerStart()
}
