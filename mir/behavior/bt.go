package behavior

import (
	. "github.com/yenkeia/mirgo/mir"
	"time"
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
	Monster *Monster
}

func (c *BT) GetTime() time.Duration {
	return c.timer
}

func (c *BT) Process() {
	c.timer += 1 * time.Second
	c.Root.Visit(c)
	c.Root.Step()
}

func init() {
	SetMonsterBehaviorFactory(NewBehavior)
}

func NewBehavior(id int, mon *Monster) IBehavior {

	var root INode

	switch id {
	case 2:
		root = DeerBrain()
	default:
		root = DefaultBrain()
	}

	bt := &BT{
		Root:    root,
		Monster: mon,
	}

	return bt
}
