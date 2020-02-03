package main

import (
	"github.com/yenkeia/mirgo/common"
)

func (c *Character) GetMagic(spell common.Spell) *common.UserMagic {
	for i := range c.Magics {
		userMagic := c.Magics[i]
		if userMagic.Spell == spell {
			return &userMagic
		}
	}
	return nil
}

func (c *Character) GetClientMagics() []common.ClientMagic {
	gdb := c.Player.Map.Env.GameDB
	res := make([]common.ClientMagic, 0)
	for i := range c.Magics {
		userMagic := c.Magics[i]
		info := gdb.GetMagicInfoByID(userMagic.MagicID)
		res = append(res, userMagic.GetClientMagic(info))
	}
	return res
}

func (c *Character) UseMagic(spell common.Spell, magic *common.UserMagic, target IMapObject) (cast bool, targetID uint32) {
	cast = true
	switch spell {
	case common.SpellFireBall, common.SpellGreatFireBall, common.SpellFrostCrunch:
		if ok := c.Fireball(target, magic); !ok {
			targetID = 0
		}
	case common.SpellHealing:
		if target == nil {
			target = c.Player
			targetID = c.Player.GetID()
		}
		c.Healing(target, magic)
	case common.SpellRepulsion, common.SpellEnergyRepulsor, common.SpellFireBurst:
		c.Repulsion(magic)
	case common.SpellElectricShock:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 500, magic, target as MonsterObject));
		action := NewDelayedAction(c.NewObjectID(), DelayedTypeMagic, NewTask(c.CompleteMagic, magic, target))
		c.ActionList.Store(action.ID, action)
	case common.SpellPoisoning:
		if !c.Poisoning(target, magic) {
			cast = false
		}
	case common.SpellHellFire:
		c.HellFire(magic)
	case common.SpellThunderBolt:
		c.ThunderBolt(target, magic)
	case common.SpellSoulFireBall:
		// if (!SoulFireball(target, magic, out cast)) targetID = 0;
		if !c.SoulFireball(target, magic) {
			targetID = 0
			cast = false
		}
	case common.SpellSummonSkeleton:
		c.SummonSkeleton(magic)
	case common.SpellTeleport, common.SpellBlink:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 200, magic, location));
		action := NewDelayedAction(c.NewObjectID(), DelayedTypeMagic, NewTask(c.CompleteMagic, magic, c.Player.GetPoint()))
		c.ActionList.Store(action.ID, action)
	case common.SpellHiding:
		c.Hiding(magic)
	case common.SpellHaste, common.SpellLightBody:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 500, magic));
		action := NewDelayedAction(c.NewObjectID(), DelayedTypeMagic, NewTask(c.CompleteMagic, magic))
		c.ActionList.Store(action.ID, action)
	case common.SpellFury:
		cast = c.FurySpell(magic)
	case common.SpellImmortalSkin:
		cast = c.ImmortalSkin(magic)
	case common.SpellFireBang, common.SpellIceStorm:
		// FireBang(magic, target == null ? location : target.CurrentLocation);
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		c.FireBang(magic, location)
	case common.SpellMassHiding:
		// MassHiding(magic, target == null ? location : target.CurrentLocation, out cast);
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		cast = c.MassHiding(magic, location)
	case common.SpellSoulShield, common.SpellBlessedArmour:
		// SoulShield(magic, target == null ? location : target.CurrentLocation, out cast);
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		cast = c.SoulShield(magic, location)
	case common.SpellFireWall:
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		c.FireWall(magic, location)
	case common.SpellLightning:
		c.Lightning(magic)
	case common.SpellHeavenlySword:
		c.HeavenlySword(magic)
	case common.SpellMassHealing:
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		c.MassHealing(magic, location)
	case common.SpellShoulderDash:
		c.ShoulderDash(magic)
	case common.SpellThunderStorm, common.SpellFlameField, common.SpellStormEscape:
		/*
			ThunderStorm(magic);
			if (spell == Spell.FlameField)
				SpellTime = Envir.Time + 2500; //Spell Delay
			if (spell == Spell.StormEscape)
				//Start teleport.
				ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 750, magic, location));
		*/
	case common.SpellMagicShield:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 500, magic, magic.GetPower(GetAttackPower(MinMC, MaxMC) + 15)));
		action := NewDelayedAction(c.NewObjectID(), DelayedTypeMagic, NewTask(c.CompleteMagic, magic, magic.GetPower(c.GetAttackPower(int(c.MinMC), int(c.MaxMC))+15)))
		c.ActionList.Store(action.ID, action)
	case common.SpellFlameDisruptor:
		c.FlameDisruptor(target, magic)
	case common.SpellTurnUndead:
		c.TurnUndead(target, magic)
	case common.SpellMagicBooster:
		c.MagicBooster(magic)
	case common.SpellVampirism:
		c.Vampirism(target, magic)
	case common.SpellSummonShinsu:
		c.SummonShinsu(magic)
	case common.SpellPurification:
		/*
			if (target == null)
			{
				target = this;
				targetID = ObjectID;
			}
			Purification(target, magic);
		*/
	case common.SpellLionRoar, common.SpellBattleCry:
		// CurrentMap.ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 500, this, magic, CurrentLocation));
	case common.SpellRevelation:
		c.Revelation(target, magic)
	case common.SpellPoisonCloud:
		cast = c.PoisonCloud(magic, c.Player.GetPoint())
	case common.SpellEntrapment:
		c.Entrapment(target, magic)
	case common.SpellBladeAvalanche:
		c.BladeAvalanche(magic)
	case common.SpellSlashingBurst:
		cast = c.SlashingBurst(magic)
	case common.SpellRage:
		c.Rage(magic)
	case common.SpellMirroring:
		c.Mirroring(magic)
	case common.SpellBlizzard:
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		cast = c.Blizzard(magic, location)
	case common.SpellMeteorStrike:
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		cast = c.MeteorStrike(magic, location)
	case common.SpellIceThrust:
		c.IceThrust(magic)
	case common.SpellProtectionField:
		c.ProtectionField(magic)
	case common.SpellPetEnhancer:
		cast = c.PetEnhancer(target, magic)
	case common.SpellTrapHexagon:
		cast = c.TrapHexagon(magic, target)
	case common.SpellReincarnation:
		// Reincarnation(magic, target == null ? null : target as PlayerObject, out cast);
		if target != nil {
			target = c.Player
		}
		cast = c.Reincarnation(magic, target)
	case common.SpellCurse:
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		cast = c.Curse(magic, location)
	case common.SpellSummonHolyDeva:
		c.SummonHolyDeva(magic)
	case common.SpellHallucination:
		c.Hallucination(target, magic)
	case common.SpellEnergyShield:
		cast = c.EnergyShield(target, magic)
	case common.SpellUltimateEnhancer:
		cast = c.UltimateEnhancer(target, magic)
	case common.SpellPlague:
		location := target.GetPoint()
		if target == nil {
			location = c.Player.GetPoint()
		}
		cast = c.Plague(magic, location)
	default:
		cast = false
	}
	return
}

func (c *Character) CompleteMagic(args ...interface{}) {
	userMagic := args[0].(*common.UserMagic)
	switch userMagic.Spell {
	case common.SpellFireBall, common.SpellGreatFireBall, common.SpellThunderBolt, common.SpellSoulFireBall, common.SpellFlameDisruptor, common.SpellStraightShot, common.SpellDoubleShot:
		value := args[1].(int)
		target := args[2].(IMapObject)
		if target == nil || !target.IsAttackTarget(c.Player) {
			return
		}
		if target.GetRace() == common.ObjectTypePlayer {
			target.(*Player).Attacked(c.Player, value, common.DefenceTypeMAC, false)
		} else if target.GetRace() == common.ObjectTypeMonster {
			target.(*Monster).Attacked(c.Player, value, common.DefenceTypeMAC, false)
		}
		return
	case common.SpellFrostCrunch:
	case common.SpellVampirism:
	case common.SpellHealing:
		value := args[1].(int)
		target := args[2].(IMapObject)
		if target == nil || !target.IsFriendlyTarget(c.Player) {
			return
		}
		if target.GetRace() == common.ObjectTypePlayer {
			obj := target.(*Player)
			hp := int(obj.HP)
			maxHP := int(obj.MaxHP)
			if hp >= maxHP {
				return
			}
			obj.HP += uint16(value)
		} else if target.GetRace() == common.ObjectTypeMonster {
			obj := target.(*Monster)
			hp := int(obj.HP)
			maxHP := int(obj.MaxHP)
			if hp >= maxHP {
				return
			}
			obj.HP += uint32(value)
		}
		// LevelMagic(magic)
	case common.SpellElectricShock:
	case common.SpellPoisoning:
	case common.SpellStormEscape:
	case common.SpellTeleport:
	case common.SpellBlink:
	case common.SpellHiding:
	case common.SpellHaste:
	case common.SpellFury:
	case common.SpellImmortalSkin:
	case common.SpellLightBody:
	case common.SpellMagicShield:
	case common.SpellTurnUndead:
	case common.SpellMagicBooster:
	case common.SpellPurification:
	case common.SpellRevelation:
	case common.SpellReincarnation:
	case common.SpellEntrapment:
	case common.SpellHallucination:
	case common.SpellPetEnhancer:
	case common.SpellElementalBarrier:
	case common.SpellElementalShot:
	case common.SpellDelayedExplosion:
	}
}

func (c *Character) Fireball(target IMapObject, magic *common.UserMagic) bool {
	if target == nil || !target.IsAttackTarget(c.Player) {
		return false
	}
	damage := magic.GetDamage(c.GetAttackPower(int(c.MinMC), int(c.MaxMC)))
	action := NewDelayedAction(c.NewObjectID(), DelayedTypeMagic, NewTask(c.CompleteMagic, magic, damage, target))
	c.ActionList.Store(action.ID, action)
	return true
}

func (c *Character) Healing(target IMapObject, magic *common.UserMagic) {
	if target == nil || !target.IsFriendlyTarget(c.Player) {
		return
	}
	// int health = magic.GetDamage(GetAttackPower(MinSC, MaxSC) * 2) + Level;
	health := magic.GetDamage(c.GetAttackPower(int(c.MinSC), int(c.MaxSC))*2) + int(c.Level)
	action := NewDelayedAction(c.NewObjectID(), DelayedTypeMagic, NewTask(c.CompleteMagic, magic, health, target))
	c.ActionList.Store(action.ID, action)
}

func (c *Character) Repulsion(magic *common.UserMagic) {

}

func (c *Character) Poisoning(target IMapObject, magic *common.UserMagic) bool {
	return true
}

func (c *Character) HellFire(magic *common.UserMagic) {

}

func (c *Character) ThunderBolt(target IMapObject, magic *common.UserMagic) {

}

func (c *Character) SoulFireball(target IMapObject, magic *common.UserMagic) bool {
	return true
}

func (c *Character) SummonSkeleton(magic *common.UserMagic)                           {}
func (c *Character) Hiding(magic *common.UserMagic)                                   {}
func (c *Character) FurySpell(magic *common.UserMagic) bool                           { return true }
func (c *Character) ImmortalSkin(magic *common.UserMagic) bool                        { return true }
func (c *Character) FireBang(magic *common.UserMagic, location common.Point)          {}
func (c *Character) MassHiding(magic *common.UserMagic, location common.Point) bool   { return true }
func (c *Character) SoulShield(magic *common.UserMagic, location common.Point) bool   { return true }
func (c *Character) FireWall(magic *common.UserMagic, location common.Point)          {}
func (c *Character) Lightning(magic *common.UserMagic)                                {}
func (c *Character) HeavenlySword(magic *common.UserMagic)                            {}
func (c *Character) MassHealing(magic *common.UserMagic, location common.Point)       {}
func (c *Character) ShoulderDash(magic *common.UserMagic)                             {}
func (c *Character) FlameDisruptor(target IMapObject, magic *common.UserMagic)        {}
func (c *Character) TurnUndead(target IMapObject, magic *common.UserMagic)            {}
func (c *Character) MagicBooster(magic *common.UserMagic)                             {}
func (c *Character) Vampirism(target IMapObject, magic *common.UserMagic)             {}
func (c *Character) SummonShinsu(magic *common.UserMagic)                             {}
func (c *Character) Revelation(target IMapObject, magic *common.UserMagic)            {}
func (c *Character) PoisonCloud(magic *common.UserMagic, location common.Point) bool  { return true }
func (c *Character) Entrapment(target IMapObject, magic *common.UserMagic)            {}
func (c *Character) BladeAvalanche(magic *common.UserMagic)                           {}
func (c *Character) SlashingBurst(magic *common.UserMagic) bool                       { return true }
func (c *Character) Rage(magic *common.UserMagic)                                     {}
func (c *Character) Mirroring(magic *common.UserMagic)                                {}
func (c *Character) Blizzard(magic *common.UserMagic, location common.Point) bool     { return true }
func (c *Character) MeteorStrike(magic *common.UserMagic, location common.Point) bool { return true }
func (c *Character) IceThrust(magic *common.UserMagic)                                {}
func (c *Character) ProtectionField(magic *common.UserMagic)                          {}
func (c *Character) PetEnhancer(target IMapObject, magic *common.UserMagic) bool      { return true }
func (c *Character) TrapHexagon(magic *common.UserMagic, target IMapObject) bool      { return true }
func (c *Character) Reincarnation(magic *common.UserMagic, target IMapObject) bool    { return true }
func (c *Character) Curse(magic *common.UserMagic, location common.Point) bool        { return true }
func (c *Character) SummonHolyDeva(magic *common.UserMagic)                           {}
func (c *Character) Hallucination(target IMapObject, magic *common.UserMagic) bool    { return true }
func (c *Character) EnergyShield(target IMapObject, magic *common.UserMagic) bool     { return true }
func (c *Character) UltimateEnhancer(target IMapObject, magic *common.UserMagic) bool { return true }
func (c *Character) Plague(magic *common.UserMagic, location common.Point) bool       { return true }
