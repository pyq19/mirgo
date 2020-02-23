package behavior

import (
	"time"

	"github.com/yenkeia/mirgo/mir"
)

type Status uint8

const (
	SUCCESS Status = 1
	FAILED  Status = 2
	RUNNING Status = 3
	READY   Status = 4
)

type INode interface {
	Visit(*BT)
	Reset()
	Step()
	Status() Status
}

type BT struct {
	timer   time.Duration // 记录从启动开始的时间
	Root    INode
	Monster *mir.Monster
}

func (c *BT) GetTime() time.Duration {
	return c.timer
}

func (c *BT) Process(dt time.Duration) {
	c.timer += dt
	c.Root.Visit(c)
	c.Root.Step()
}

func init() {
	mir.SetMonsterBehaviorFactory(NewBehavior)
}

func NewBehavior(id int, mon *mir.Monster) mir.IBehavior {

	var root INode

	switch id {
	case 2:
		root = DeerBrain()
	case 6:
		root = GuardBrain()
	default:
		root = DefaultBrain()
	}

	bt := &BT{
		Root:    root,
		Monster: mon,
	}

	return bt
}
