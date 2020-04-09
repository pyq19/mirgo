package behavior

// 警戒
type AttackWallNode struct {
	Node
}

func (n *AttackWallNode) Visit(c *BT) {
	if n.Status() == READY {
		c.Monster.FindTarget()

		if c.Monster.Target != nil {
			n.status = RUNNING
		} else {
			n.status = FAILED
		}
	}

	if n.Status() == RUNNING {
		if c.Monster.CanAttack() && c.Monster.InAttackRange() {
			c.Monster.Attack()
			n.status = SUCCESS
		} else {
			n.status = FAILED
		}
	}
}

func AttackWall() INode {
	return &AttackWallNode{}
}
