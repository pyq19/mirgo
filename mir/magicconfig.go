package mir

import "github.com/yenkeia/mirgo/common"

func init() {
	add := func(sp common.Spell, c *MagicConfig) {

		c.Spell = sp
		if c.ItemCost == nil {
			c.ItemCost = Cost_None
		}
		if c.DelayAt == nil {
			c.DelayAt = DelayAt_Player
		}

		configsMap[sp] = c
	}

	add(common.SpellFireBall, &MagicConfig{
		SelectType: Select_Point | Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
	})

	add(common.SpellGreatFireBall, &MagicConfig{
		SelectType: Select_Point | Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
	})

	add(common.SpellFrostCrunch, &MagicConfig{
		SelectType: Select_Point | Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_FrostCrunch,
		Formula:    Formula_MC,
	})

	add(common.SpellThunderBolt, &MagicConfig{
		SelectType: Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
	})

	add(common.SpellSoulFireBall, &MagicConfig{
		SelectType: Select_Enemy | Select_Point,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
		ItemCost:   Cost_Amulet,
	})

	add(common.SpellHealing, &MagicConfig{
		SelectType: Select_Self | Select_Friend,
		TargetType: Target_Friend | Target_Self,
		Action:     Action_HealingTarget,
		Formula:    Formula_MC,
	})

}
