package behavior

import (
	"github.com/yenkeia/mirgo/game"
	"github.com/yenkeia/mirgo/game/cm"
)

func HasTarget(c *BT) bool {
	return c.Monster.Target != nil
}

func FindMonsterInViewRange(c *BT) bool {
	c.Monster.FindTarget()
	return c.Monster.Target != nil
}

func FindPlayerByPKPoints200(c *BT) bool {
	return FindPlayerByPKPoints(c, 200)
}

func FindPlayerByPKPoints(c *BT, pkpoints int) bool {
	m := c.Monster

	m.Map.RangeObject(m.CurrentLocation, m.ViewRange, func(o game.IMapObject) bool {

		if o == m {
			return true
		}

		switch o.GetRace() {
		case cm.ObjectTypePlayer:
			p := o.(*game.Player)
			if !p.IsAttackTarget(m) {
				return true // continue
			}

			// if (playerob.PKPoints < 200 || ob.Hidden && (!CoolEye || Level < ob.Level)) continue;
			if p.PKPoints < pkpoints {
				return true
			}

			m.Target = o

			return false
		}

		return true
	})

	return m.Target != nil
}
