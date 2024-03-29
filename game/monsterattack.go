package game

import (
	"time"
)

// 专用于大刀卫士攻击
func (m *Monster) GuardAttack() {
	if !m.Target.IsAttackTarget(m) {
		return
	}

	target := ObjectBack(m.Target)

	dir := DirectionFromPoint(target, m.Target.GetPoint())

	m.Broadcast(&SM_ObjectAttack{
		ObjectID:  m.GetID(),
		LocationX: int32(target.X),
		LocationY: int32(target.Y),
		Direction: dir,
		Spell:     SpellNone,
		Level:     uint8(0),
		Type:      uint8(0),
	})
	m.Broadcast(&SM_ObjectTurn{
		ObjectID:  m.GetID(),
		Direction: m.Direction,
		Location:  m.CurrentLocation,
	})

	now := time.Now()
	// ActionTime = Envir.Time + 300;
	m.AttackTime = now.Add(time.Duration(m.AttackSpeed) * time.Millisecond)

	damage := m.GetAttackPower(int(m.MinDC), int(m.MaxDC))

	if m.Target.GetRace() == ObjectTypePlayer {
		damage = int(^uint(0) >> 1) // INTMAX
	}

	if damage <= 0 {
		return
	}

	m.Target.Attacked(m, damage, DefenceTypeAgility, false)
}

func (m *Monster) SpittingSpiderAttack() {

}
