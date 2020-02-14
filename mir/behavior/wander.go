package behavior

import (
	. "github.com/yenkeia/mirgo/mir"
	"time"
)

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
			if RandomNext(10) == 1 {
				n.PickNewDirection(c)
			}
		}
	}
}

func (n *WanderNode) WaitTo(t time.Duration) {
	n.waittime = t
}

func (n *WanderNode) PickNewDirection(c *BT) {

	switch RandomNext(3) {
	case 0:
		c.Monster.Turn(RandomDirection())
	default:
		c.Monster.Walk(c.Monster.CurrentDirection)
	}
}
