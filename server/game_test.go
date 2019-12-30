package main

import (
	"os"
	"sync"
	"testing"
)

func TestGameMonsters(t *testing.T) {
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
				t.Logf("Coordinate: %s, MonsterID: %d, ptr: %p", o.CurrentLocation.Coordinate(), o.ID, o)
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

func TestGameNPCs(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	var addr = "0.0.0.0:7000"
	var mirDB = "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	conf := Config{addr, gopath + mirDB}
	g := NewGame(conf)
	v, _ := g.Env.Maps.Load(1)
	count := 0
	v.(*Map).AOI.grids.Range(func(k, v interface{}) bool {
		g := v.(*Grid)
		//t.Log(g.String())
		g.NPCs.Range(func(k, v interface{}) bool {
			n := v.(*NPC)
			if n != nil {
				t.Log(n.String())
				count += 1
			}
			return true
		})
		return true
	})
	t.Log(count)
}

func TestEnviron_NewObjectID(t *testing.T) {
	var wg sync.WaitGroup
	e := new(Environ)
	e.ObjectID = 0
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			e.NewObjectID()
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log(e.ObjectID)
}
