package behavior

import "time"

// 鹿
func DeerBrain() INode {

	root := Priority(1*time.Second,
		While(HasTarget, ChaseAndAttack()),
		Wander(),
	)

	return root
}

func HasTarget(c *BT) bool {
	return c.Monster.Target != nil
}
