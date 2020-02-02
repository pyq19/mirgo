package main

import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/yenkeia/mirgo/common"
)

func TestGameMonsters(t *testing.T) {
	g := NewGame()

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
		g.Objects.Range(func(k, v interface{}) bool {
			if v.(IMapObject).GetRace() == common.ObjectTypeMonster {
				o := v.(*Monster)
				if o != nil {
					count2 += 1
					t.Logf("Coordinate: %s, MonsterID: %d, ptr: %p", o.CurrentLocation.Coordinate(), o.ID, o)
				}
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
	g := NewGame()
	count := 0
	g.Env.Maps.Range(func(k, v interface{}) bool {
		v.(*Map).AOI.grids.Range(func(k, v interface{}) bool {
			g := v.(*Grid)
			//t.Log(g.String())
			g.Objects.Range(func(k, v interface{}) bool {
				if v.(IMapObject).GetRace() == common.ObjectTypeMerchant {
					n := v.(*NPC)
					if n != nil {
						//t.Log(n.String())
						count += 1
						if _, err := os.Stat(n.FilePath); err != nil {
							t.Logf("文件: %s 不存在\n", n.Name)
						}
					}
				}
				return true
			})
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

func TestGameDB_GetItemInfoByID(t *testing.T) {
	g := NewGame()
	gdb := g.Env.GameDB
	i1 := gdb.GetItemInfoByID(1)
	i2 := gdb.GetItemInfoByID(2)
	t.Log(i1)
	t.Log(i2)
}

func TestGameDB_GetMonsterInfoByID(t *testing.T) {
	g := NewGame()
	gdb := g.Env.GameDB
	m1 := gdb.GetMonsterInfoByID(1)
	m2 := gdb.GetMonsterInfoByID(2)
	t.Log(m1)
	t.Log(m2)
}

func TestStringSplit(t *testing.T) {
	sl := []string{"@", "@ ", "@a", "@a ", "@a b"}
	for i := range sl {
		res := strings.Split(sl[i][1:], " ")
		t.Log("----", res, len(res))
		for j := range res {
			t.Log(res[j])
		}
	}
}
