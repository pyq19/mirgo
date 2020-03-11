package mir

import (
	"container/list"
	"time"

	"github.com/yenkeia/mirgo/common"
)

type BuffList struct {
	List *list.List
}

type Buff struct {
	BuffType   common.BuffType
	Caster     IMapObject
	Visible    bool // 是否可见
	ObjectID   uint32
	ExpireTime time.Time // 过期时间️
	Values     []int     // public int[] Values
	Infinite   bool      // 是否永久
	Paused     bool
}

func NewBuff(buffType common.BuffType, caster IMapObject, expireTime int, values []int) *Buff {
	return &Buff{
		BuffType:   buffType,
		Caster:     caster,
		Visible:    false,
		ObjectID:   0,
		ExpireTime: time.Now().Add(time.Duration(expireTime) * time.Millisecond),
		Values:     values,
		Infinite:   false,
		Paused:     false,
	}
}

func NewBuffList() *BuffList {
	ret := &BuffList{}
	ret.List = list.New()
	return ret
}
