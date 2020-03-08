package mir

import (
	"time"

	"github.com/yenkeia/mirgo/common"
)

type Buff struct {
	ObjectID   uint32
	BuffType   common.BuffType
	Visible    bool      // 是否可见
	Infinite   bool      // 是否永久
	Values     int       // public int[] Values
	ExpireTime time.Time // 过期时间️
}

func NewBuff(id uint32, typ common.BuffType, value int, expire time.Time) *Buff {
	return &Buff{
		ObjectID:   id,
		BuffType:   typ,
		Visible:    false,
		Infinite:   false,
		Values:     value,
		ExpireTime: expire,
	}
}
