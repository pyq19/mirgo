package mir

import (
	"container/list"

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
	ExpireTime int   // time.Time // 过期时间️ 毫秒
	Values     []int // public int[] Values
	Infinite   bool  // 是否永久
	Paused     bool
}

func NewBuffList() *BuffList {
	ret := &BuffList{}
	ret.List = list.New()
	return ret
}
