package game

import (
	"container/list"
	"time"

	"github.com/yenkeia/mirgo/game/cm"
)

type Buff struct {
	ObjectID   uint32
	Type       cm.BuffType
	Caster     IMapObject
	Visible    bool      // 是否可见
	ExpireTime time.Time // 过期时间️
	Values     []int32   // public int[] Values
	Infinite   bool      // 是否永久
	Paused     bool
}

func NewBuff(buffType cm.BuffType, caster IMapObject, expireTime int, values []int32) *Buff {
	return &Buff{
		Type:       buffType,
		Caster:     caster,
		Visible:    false,
		ObjectID:   0,
		ExpireTime: time.Now().Add(time.Duration(expireTime) * time.Millisecond),
		Values:     values,
		Infinite:   false,
		Paused:     false,
	}
}

type BuffList struct {
	List *list.List
}

func NewBuffList() *BuffList {
	ret := &BuffList{}
	ret.List = list.New()
	return ret
}

func (bl *BuffList) AddBuff(b *Buff) {
	for it := bl.List.Front(); it != nil; it = it.Next() {
		buf := it.Value.(*Buff)
		if buf.Type != b.Type {
			continue
		}

		// 新的替换旧的
		b.Paused = false
		bl.List.InsertBefore(b, it)
		bl.List.Remove(it)
		return
	}

	bl.List.PushBack(b)
}

func (bl *BuffList) RemoveBuff(t cm.BuffType) {
	for it := bl.List.Front(); it != nil; it = it.Next() {
		buf := it.Value.(*Buff)
		if buf.Type != t {
			continue
		}

		buf.Infinite = false
		buf.ExpireTime = time.Now()
	}
}

func (bl *BuffList) Has(f func(b *Buff) bool) bool {
	for it := bl.List.Front(); it != nil; it = it.Next() {

		if f(it.Value.(*Buff)) {
			return true
		}
	}
	return false
}
