package behavior

// 追杀
type ChaseAndAttackNode struct {
	Node
}

func ChaseAndAttack() INode {
	return &ChaseAndAttackNode{}
}

func (n *ChaseAndAttackNode) Visit(c *BT) {
	if n.Status() == READY {

		// if c.Monster.Target == nil {
		// 	c.Monster.FindTarget()
		// }

		if c.Monster.Target == nil {
			n.status = FAILED
		} else {
			n.status = RUNNING
		}
	}

	if n.Status() == RUNNING {

		if c.Monster.Target == nil || c.Monster.Target.IsDead() {
			n.status = SUCCESS
			c.Monster.Target = nil
		} else {
			if c.Monster.InAttackRange() {
				if c.Monster.CanAttack() {
					c.Monster.Attack()
				}

				if c.Monster.Target == nil || c.Monster.Target.IsDead() {
					c.Monster.Target = nil
				}
			} else {
				c.Monster.MoveTo(c.Monster.Target.GetPoint())
			}
		}
	}
}
