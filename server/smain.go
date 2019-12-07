package main

var addr = "0.0.0.0:7000"
var mirDB = "../dotnettools/mir.sqlite"
var accountDB = "../dotnettools/account.sqlite"

type Config struct {
	Addr      string
	MirDB     string
	AccountDB string
}

func main() {
	conf := &Config{addr, mirDB, accountDB}
	g := NewGame(*conf)
	g.ServerStart()
}
