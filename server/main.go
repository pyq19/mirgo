package main

import (
	"github.com/yenkeia/mirgo/mir"
	_ "github.com/yenkeia/mirgo/mir/behavior"
)

func main() {
	mir.NewEnviron().ServerStart()
}
