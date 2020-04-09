package behavior

import (
	"time"

	"github.com/yenkeia/mirgo/mir"
	"github.com/yenkeia/mirgo/ut"
)

// 游荡
type WanderNode struct {
	Node
	waittime time.Duration
}

func Wander() INode {
	return &WanderNode{}
}

func (n *WanderNode) Visit(c *BT) {
	if n.Status() == READY {

		n.status = RUNNING
		n.WaitTo(1*time.Second + c.GetTime())

	} else if n.Status() == RUNNING {

		if c.GetTime() > n.waittime {
			if ut.RandomNext(10) <= 3 {
				n.PickNewDirection(c)
			}
		}
	}
}

func (n *WanderNode) WaitTo(t time.Duration) {
	n.waittime = t
}

func (n *WanderNode) PickNewDirection(c *BT) {

	switch ut.RandomNext(3) {
	case 0:
		c.Monster.Turn(mir.RandomDirection())
	default:
		c.Monster.Walk(c.Monster.Direction)
	}
}
