package behavior

import (
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/mir"
)

func init() {
	mir.SetMonsterBehaviorFactory(NewBehavior)
}

func NewBehavior(id int, mon *mir.Monster) mir.IBehavior {
	return NewBehaviorDefault(mon)
}

type BehaviorDefault struct {
	Mon        *mir.Monster
	SearchTime time.Time // 怪物下一次搜索目标的时间
	RoamTime   time.Time
}

func NewBehaviorDefault(mon *mir.Monster) *BehaviorDefault {
	m := &BehaviorDefault{Mon: mon}
	now := time.Now()
	m.SearchTime = now
	m.RoamTime = now
	return m
}

func (b *BehaviorDefault) Process() {
	if b.Mon.IsDead() {
		return
	}
	b.ProcessSearch()
	b.ProcessRoam()
	b.ProcessTarget()
}

// ProcessSearch 寻找目标
func (b *BehaviorDefault) ProcessSearch() {
	now := time.Now()
	if b.SearchTime.After(now) {
		return
	}
	b.SearchTime = now.Add(1 * time.Second)

	if b.Mon.CanMove() && b.Mon.CheckStacked() {

		// Walk Randomly
		if !b.Mon.Walk(b.Mon.CurrentDirection) {

			dir := b.Mon.CurrentDirection

			switch mir.RandomNext(3) {
			case 0:
				for i := 0; i < common.MirDirectionCount; i++ {
					dir = mir.NextDirection(dir)

					if b.Mon.Walk(dir) {
						break
					}
				}
			default:
				for i := 0; i < common.MirDirectionCount; i++ {
					dir = mir.NextDirection(dir)

					if b.Mon.Walk(dir) {
						break
					}
				}
			}
		}
	}

	if b.Mon.Target == nil {
		b.Mon.FindTarget()
	}
}

// 巡逻
func (b *BehaviorDefault) ProcessRoam() {
	now := time.Now()
	if b.RoamTime.After(now) {
		return
	}
	b.RoamTime = now.Add(1 * time.Second)

	if mir.RandomNext(10) != 0 {
		return
	}

	switch mir.RandomNext(3) {
	case 0:
		b.Mon.Turn(mir.RandomDirection())
	default:
		b.Mon.Walk(b.Mon.CurrentDirection)
	}
}

func (b *BehaviorDefault) ProcessTarget() {
	if b.Mon.Target == nil || !b.Mon.CanAttack() {
		return
	}
	if b.Mon.InAttackRange() {
		b.Mon.Attack()
		if b.Mon.Target.IsDead() {
			b.Mon.FindTarget()
		}
		return
	}
	b.Mon.MoveTo(b.Mon.Target.GetPoint())
}

// type BehaviorDeer struct {
// 	RunAway bool
// }

// func NewBehaviorDeer(mon *Monster) *BehaviorDeer {
// 	return &BehaviorDeer{
// 		RunAway: RandomNext(7) == 0,
// 	}
// }

// func (b *BehaviorDeer) Process() {

// }
