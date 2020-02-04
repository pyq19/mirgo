package server

import (
	"math/rand"
	"sync"
	"time"
)

var (
	r      *rand.Rand
	G_Rand *RandGenerator
)

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
	G_Rand = &RandGenerator{
		used: make(map[string]bool),
		lock: sync.RWMutex{},
	}
}

type RandGenerator struct {
	used map[string]bool
	lock sync.RWMutex // 保护 used map 的锁
}

// RandString 生成随机字符串 https://www.jianshu.com/p/ff8539aff912
func (g *RandGenerator) RandString(length int) string {
	g.lock.Lock()
	defer g.lock.Unlock()

	res := ""
	for {
		res = randString(length)
		if g.used[res] == true {
			continue
		}
		g.used[res] = true
		break
	}
	return res
}

func randString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func (g *RandGenerator) RandInt(min, max int) int {
	if min <= 0 {
		min = 0
	}
	return rand.Intn(max-min) + min
}
