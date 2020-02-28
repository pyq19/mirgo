package behavior

import (
	"time"
)

// 默认
func DefaultBrain() INode {

	root := Priority(1*time.Second,
		AttackWall(),
		ChaseAndAttack(),
		Wander(),
	)

	return root
}

// 鹿
func DeerBrain() INode {

	root := Priority(1*time.Second,
		While(HasTarget, ChaseAndAttack()),
		Wander(),
	)

	return root
}

// 守卫（大刀卫士）
func GuardBrain() INode {

	root := Priority(1*time.Second,
		While(FindMonsterInViewRange, GuardAttack()),
	)

	return root
}

// 弓箭守卫
func ArcherGuardBrain() INode {

	root := Priority(1*time.Second,
		While(FindPlayerByPKPoints200, WatchAndShoot()),
	)

	return root
}
