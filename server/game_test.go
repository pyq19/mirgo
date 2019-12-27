package main

import (
	"os"
	"testing"
)

func TestGame(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	var addr = "0.0.0.0:7000"
	var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	conf := Config{addr, gopath + mirDB}
	g := NewGame(conf)

	//v, _ := g.Env.Maps.Load(1)
	//count1 := 0
	//v.(*Map).Cells.Range(func(k, m interface{}) bool {
	//	c := m.(*Cell)
	//	if c.Object != nil {
	//		count1 += 1
	//	}
	//	return true
	//})
	//log.Debugln("!!!", count1)

	v, _ := g.Env.Maps.Load(1)
	count2 := 0
	v.(*Map).AOI.grids.Range(func(k, v interface{}) bool {
		g := v.(*Grid)
		t.Log(g.String())
		g.Monsters.Range(func(k, v interface{}) bool {
			o := v.(*Monster)
			if o != nil {
				count2 += 1
				t.Logf("Coordinate: %s, MonsterID: %s, ptr: %p", o.CurrentLocation.Coordinate(), o.ID, o)
			}
			return true
		})
		return true
	})
	log.Debugln("!!!", count2)
}

func TestInterface(t *testing.T) {
	m := new(Monster)
	t.Logf("%p\n", m)
	i(t, m)
}

func i(t *testing.T, o interface{}) {
	t.Logf("%p\n", o)
}
