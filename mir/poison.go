package mir

import (
	"container/list"
	"time"

	"github.com/yenkeia/mirgo/common"
)

type PoisonList struct {
	List *list.List
}

type Poison struct {
	Owner     IMapObject
	PType     common.PoisonType
	Value     int           // 效果总数
	Duration  time.Duration // 持续多久（秒）
	TickSpeed time.Duration
	TickNum   int // 总共跳几次
	TickTime  int // 当前第几跳
}

func NewPoison(duration int, owner IMapObject, ptype common.PoisonType, tickSpeed int, value int) *Poison {
	d := time.Duration(duration) * time.Second       // 持续多少秒
	t := time.Duration(tickSpeed) * time.Millisecond // 两次间隔多少毫秒
	tickNum := int(d / t)                            // 总共跳几次
	return &Poison{
		Owner:     owner,
		PType:     ptype,
		Value:     value,
		Duration:  d,
		TickSpeed: t,
		TickNum:   tickNum,
		TickTime:  0,
	}
}

func NewPoisonList() *PoisonList {
	ret := &PoisonList{}
	ret.List = list.New()
	return ret
}

func (ls *PoisonList) Execute() {

}
