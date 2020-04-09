package behavior

import (
	"time"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/util"
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
			if util.RandomNext(10) <= 3 {
				n.PickNewDirection(c)
			}
		}
	}
}

func (n *WanderNode) WaitTo(t time.Duration) {
	n.waittime = t
}

func (n *WanderNode) PickNewDirection(c *BT) {

	switch util.RandomNext(3) {
	case 0:
		c.Monster.Turn(cm.RandomDirection())
	default:
		c.Monster.Walk(c.Monster.Direction)
	}
}
