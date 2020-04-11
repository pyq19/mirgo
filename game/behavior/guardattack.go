package behavior

// 卫士攻击逻辑
type GuardAttackNode struct {
	Node
}

func GuardAttack() INode {
	return &GuardAttackNode{}
}

func (n *GuardAttackNode) Visit(c *BT) {

	if n.Status() == READY {
		if c.Monster.Target == nil || c.Monster.Target.IsDead() {
			n.status = FAILED
			return
		} else {
			n.status = RUNNING
		}
	}

	if n.Status() == RUNNING {
		if c.Monster.Target == nil || c.Monster.Target.IsDead() {
			n.status = SUCCESS
			c.Monster.Target = nil
		} else if c.Monster.CanAttack() {
			c.Monster.GuardAttack()
			if c.Monster.Target == nil || c.Monster.Target.IsDead() {
				c.Monster.Target = nil
			}
		} else {
			n.status = FAILED
		}
	}
}
