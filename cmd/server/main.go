package main

import (
	"github.com/pyq19/mirgo/game"
	_ "github.com/pyq19/mirgo/game/behavior"
	_ "github.com/pyq19/mirgo/game/mircodec"
)

func main() {
	game.NewEnviron().ServerStart()
}
