package main

import (
	"github.com/yenkeia/mirgo/common"
)

func (p *Player) GetMagic(spell common.Spell) *common.UserMagic {
	for i := range p.Magics {
		userMagic := p.Magics[i]
		if userMagic.Spell == spell {
			return &userMagic
		}
	}
	return nil
}

func (p *Player) GetClientMagics() []common.ClientMagic {
	gdb := p.Map.Env.GameDB
	res := make([]common.ClientMagic, 0)
	for i := range p.Magics {
		userMagic := p.Magics[i]
		info := gdb.GetMagicInfoByID(userMagic.MagicID)
		res = append(res, userMagic.GetClientMagic(info))
	}
	return res
}

func (p *Player) UseMagic(spell common.Spell, magic *common.UserMagic, target IMapObject) (cast bool, targetID uint32) {
	cast = true
	switch spell {
	case common.SpellFireBall, common.SpellGreatFireBall, common.SpellFrostCrunch:
		if ok := p.Fireball(target, magic); !ok {
			targetID = 0
		}
	case common.SpellHealing:
		if target == nil {
			target = p
			targetID = p.GetID()
		}
		p.Healing(target, magic)
	case common.SpellRepulsion, common.SpellEnergyRepulsor, common.SpellFireBurst:
		p.Repulsion(magic)
	case common.SpellElectricShock:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 500, magic, target as MonsterObject));
		action := NewDelayedAction(p.NewObjectID(), DelayedTypeMagic, NewTask(p.CompleteMagic, magic, target))
		p.ActionList.Store(action.ID, action)
	case common.SpellPoisoning:
		if !p.Poisoning(target, magic) {
			cast = false
		}
	case common.SpellHellFire:
		p.HellFire(magic)
	case common.SpellThunderBolt:
		p.ThunderBolt(target, magic)
	case common.SpellSoulFireBall:
		// if (!SoulFireball(target, magic, out cast)) targetID = 0;
		if !p.SoulFireball(target, magic) {
			targetID = 0
			cast = false
		}
	case common.SpellSummonSkeleton:
		p.SummonSkeleton(magic)
	case common.SpellTeleport, common.SpellBlink:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 200, magic, location));
		action := NewDelayedAction(p.NewObjectID(), DelayedTypeMagic, NewTask(p.CompleteMagic, magic, p.GetPoint()))
		p.ActionList.Store(action.ID, action)
	case common.SpellHiding:
		p.Hiding(magic)
	case common.SpellHaste, common.SpellLightBody:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 500, magic));
		action := NewDelayedAction(p.NewObjectID(), DelayedTypeMagic, NewTask(p.CompleteMagic, magic))
		p.ActionList.Store(action.ID, action)
	case common.SpellFury:
		cast = p.FurySpell(magic)
	case common.SpellImmortalSkin:
		cast = p.ImmortalSkin(magic)
	case common.SpellFireBang, common.SpellIceStorm:
		// FireBang(magic, target == null ? location : target.CurrentLocation);
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		p.FireBang(magic, location)
	case common.SpellMassHiding:
		// MassHiding(magic, target == null ? location : target.CurrentLocation, out cast);
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		cast = p.MassHiding(magic, location)
	case common.SpellSoulShield, common.SpellBlessedArmour:
		// SoulShield(magic, target == null ? location : target.CurrentLocation, out cast);
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		cast = p.SoulShield(magic, location)
	case common.SpellFireWall:
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		p.FireWall(magic, location)
	case common.SpellLightning:
		p.Lightning(magic)
	case common.SpellHeavenlySword:
		p.HeavenlySword(magic)
	case common.SpellMassHealing:
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		p.MassHealing(magic, location)
	case common.SpellShoulderDash:
		p.ShoulderDash(magic)
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
		action := NewDelayedAction(p.NewObjectID(), DelayedTypeMagic, NewTask(p.CompleteMagic, magic, magic.GetPower(p.GetAttackPower(int(p.MinMC), int(p.MaxMC))+15)))
		p.ActionList.Store(action.ID, action)
	case common.SpellFlameDisruptor:
		p.FlameDisruptor(target, magic)
	case common.SpellTurnUndead:
		p.TurnUndead(target, magic)
	case common.SpellMagicBooster:
		p.MagicBooster(magic)
	case common.SpellVampirism:
		p.Vampirism(target, magic)
	case common.SpellSummonShinsu:
		p.SummonShinsu(magic)
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
		p.Revelation(target, magic)
	case common.SpellPoisonCloud:
		cast = p.PoisonCloud(magic, p.GetPoint())
	case common.SpellEntrapment:
		p.Entrapment(target, magic)
	case common.SpellBladeAvalanche:
		p.BladeAvalanche(magic)
	case common.SpellSlashingBurst:
		cast = p.SlashingBurst(magic)
	case common.SpellRage:
		p.Rage(magic)
	case common.SpellMirroring:
		p.Mirroring(magic)
	case common.SpellBlizzard:
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		cast = p.Blizzard(magic, location)
	case common.SpellMeteorStrike:
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		cast = p.MeteorStrike(magic, location)
	case common.SpellIceThrust:
		p.IceThrust(magic)
	case common.SpellProtectionField:
		p.ProtectionField(magic)
	case common.SpellPetEnhancer:
		cast = p.PetEnhancer(target, magic)
	case common.SpellTrapHexagon:
		cast = p.TrapHexagon(magic, target)
	case common.SpellReincarnation:
		// Reincarnation(magic, target == null ? null : target as PlayerObject, out cast);
		if target != nil {
			target = p
		}
		cast = p.Reincarnation(magic, target)
	case common.SpellCurse:
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		cast = p.Curse(magic, location)
	case common.SpellSummonHolyDeva:
		p.SummonHolyDeva(magic)
	case common.SpellHallucination:
		p.Hallucination(target, magic)
	case common.SpellEnergyShield:
		cast = p.EnergyShield(target, magic)
	case common.SpellUltimateEnhancer:
		cast = p.UltimateEnhancer(target, magic)
	case common.SpellPlague:
		location := target.GetPoint()
		if target == nil {
			location = p.GetPoint()
		}
		cast = p.Plague(magic, location)
	default:
		cast = false
	}
	return
}

func (p *Player) CompleteMagic(args ...interface{}) {
	userMagic := args[0].(*common.UserMagic)
	switch userMagic.Spell {
	case common.SpellFireBall, common.SpellGreatFireBall, common.SpellThunderBolt, common.SpellSoulFireBall, common.SpellFlameDisruptor, common.SpellStraightShot, common.SpellDoubleShot:
		value := args[1].(int)
		target := args[2].(IMapObject)
		if target == nil || !target.IsAttackTarget(p) {
			return
		}
		if target.GetRace() == common.ObjectTypePlayer {
			target.(*Player).Attacked(p, value, common.DefenceTypeMAC, false)
		} else if target.GetRace() == common.ObjectTypeMonster {
			target.(*Monster).Attacked(p, value, common.DefenceTypeMAC, false)
		}
		return
	case common.SpellFrostCrunch:
	case common.SpellVampirism:
	case common.SpellHealing:
		value := args[1].(int)
		target := args[2].(IMapObject)
		if target == nil || !target.IsFriendlyTarget(p) {
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

func (p *Player) Fireball(target IMapObject, magic *common.UserMagic) bool {
	if target == nil || !target.IsAttackTarget(p) {
		return false
	}
	damage := magic.GetDamage(p.GetAttackPower(int(p.MinMC), int(p.MaxMC)))
	action := NewDelayedAction(p.NewObjectID(), DelayedTypeMagic, NewTask(p.CompleteMagic, magic, damage, target))
	p.ActionList.Store(action.ID, action)
	return true
}

func (p *Player) Healing(target IMapObject, magic *common.UserMagic) {
	if target == nil || !target.IsFriendlyTarget(p) {
		return
	}
	// int health = magic.GetDamage(GetAttackPower(MinSC, MaxSC) * 2) + Level;
	health := magic.GetDamage(p.GetAttackPower(int(p.MinSC), int(p.MaxSC))*2) + int(p.Level)
	action := NewDelayedAction(p.NewObjectID(), DelayedTypeMagic, NewTask(p.CompleteMagic, magic, health, target))
	p.ActionList.Store(action.ID, action)
}

func (p *Player) Repulsion(magic *common.UserMagic) {

}

func (p *Player) Poisoning(target IMapObject, magic *common.UserMagic) bool {
	return true
}

func (p *Player) HellFire(magic *common.UserMagic) {

}

func (p *Player) ThunderBolt(target IMapObject, magic *common.UserMagic) {

}

func (p *Player) SoulFireball(target IMapObject, magic *common.UserMagic) bool {
	return true
}

func (p *Player) SummonSkeleton(magic *common.UserMagic)                           {}
func (p *Player) Hiding(magic *common.UserMagic)                                   {}
func (p *Player) FurySpell(magic *common.UserMagic) bool                           { return true }
func (p *Player) ImmortalSkin(magic *common.UserMagic) bool                        { return true }
func (p *Player) FireBang(magic *common.UserMagic, location common.Point)          {}
func (p *Player) MassHiding(magic *common.UserMagic, location common.Point) bool   { return true }
func (p *Player) SoulShield(magic *common.UserMagic, location common.Point) bool   { return true }
func (p *Player) FireWall(magic *common.UserMagic, location common.Point)          {}
func (p *Player) Lightning(magic *common.UserMagic)                                {}
func (p *Player) HeavenlySword(magic *common.UserMagic)                            {}
func (p *Player) MassHealing(magic *common.UserMagic, location common.Point)       {}
func (p *Player) ShoulderDash(magic *common.UserMagic)                             {}
func (p *Player) FlameDisruptor(target IMapObject, magic *common.UserMagic)        {}
func (p *Player) TurnUndead(target IMapObject, magic *common.UserMagic)            {}
func (p *Player) MagicBooster(magic *common.UserMagic)                             {}
func (p *Player) Vampirism(target IMapObject, magic *common.UserMagic)             {}
func (p *Player) SummonShinsu(magic *common.UserMagic)                             {}
func (p *Player) Revelation(target IMapObject, magic *common.UserMagic)            {}
func (p *Player) PoisonCloud(magic *common.UserMagic, location common.Point) bool  { return true }
func (p *Player) Entrapment(target IMapObject, magic *common.UserMagic)            {}
func (p *Player) BladeAvalanche(magic *common.UserMagic)                           {}
func (p *Player) SlashingBurst(magic *common.UserMagic) bool                       { return true }
func (p *Player) Rage(magic *common.UserMagic)                                     {}
func (p *Player) Mirroring(magic *common.UserMagic)                                {}
func (p *Player) Blizzard(magic *common.UserMagic, location common.Point) bool     { return true }
func (p *Player) MeteorStrike(magic *common.UserMagic, location common.Point) bool { return true }
func (p *Player) IceThrust(magic *common.UserMagic)                                {}
func (p *Player) ProtectionField(magic *common.UserMagic)                          {}
func (p *Player) PetEnhancer(target IMapObject, magic *common.UserMagic) bool      { return true }
func (p *Player) TrapHexagon(magic *common.UserMagic, target IMapObject) bool      { return true }
func (p *Player) Reincarnation(magic *common.UserMagic, target IMapObject) bool    { return true }
func (p *Player) Curse(magic *common.UserMagic, location common.Point) bool        { return true }
func (p *Player) SummonHolyDeva(magic *common.UserMagic)                           {}
func (p *Player) Hallucination(target IMapObject, magic *common.UserMagic) bool    { return true }
func (p *Player) EnergyShield(target IMapObject, magic *common.UserMagic) bool     { return true }
func (p *Player) UltimateEnhancer(target IMapObject, magic *common.UserMagic) bool { return true }
func (p *Player) Plague(magic *common.UserMagic, location common.Point) bool       { return true }
