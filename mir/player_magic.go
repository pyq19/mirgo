package mir

import (
	"github.com/yenkeia/mirgo/common"
)

// GetMagic ...
func (p *Player) GetMagic(spell common.Spell) *common.UserMagic {
	for i := range p.Magics {
		userMagic := p.Magics[i]
		if userMagic.Spell == spell {
			return userMagic
		}
	}
	return nil
}

// LevelMagic ...
func (p *Player) LevelMagic(userMagic *common.UserMagic) {

}

// UseMagic ...
func (p *Player) UseMagic(spell common.Spell, magic *common.UserMagic, target IMapObject) (cast bool, targetID uint32) {

	cast = true
	if target != nil {
		targetID = target.GetID()
	}

	if _, ok := configsMap[spell]; ok {
		ctx := &MagicContext{Spell: spell, Magic: magic, Target: target, Player: p}
		return startMagic(ctx)
	}

	switch spell {
	case common.SpellRepulsion, common.SpellEnergyRepulsor, common.SpellFireBurst:
		p.Repulsion(magic)
	case common.SpellElectricShock:
		p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, target)
	case common.SpellHellFire:
		p.HellFire(magic)
	case common.SpellTeleport, common.SpellBlink:
		p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, p.GetPoint())
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 200, magic, location));
	case common.SpellHiding:
		p.Hiding(magic)
	case common.SpellHaste, common.SpellLightBody:
		// ActionList.Add(new DelayedAction(DelayedType.Magic, Envir.Time + 500, magic));
		p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic)
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
		p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, magic.GetPower1(p.GetAttackPower(int(p.MinMC), int(p.MaxMC))+15))
	case common.SpellFlameDisruptor:
		p.FlameDisruptor(target, magic)
	case common.SpellTurnUndead:
		p.TurnUndead(target, magic)
	case common.SpellMagicBooster:
		p.MagicBooster(magic)
	case common.SpellVampirism:
		p.Vampirism(target, magic)
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

// CompleteMagic ...
func (p *Player) CompleteMagic(args ...interface{}) {
	magic := args[0].(*common.UserMagic)
	switch magic.Spell {
	case common.SpellFlameDisruptor, common.SpellStraightShot, common.SpellDoubleShot:
		value := args[1].(int)
		target := args[2].(IMapObject)
		if target == nil || !target.IsAttackTarget(p) {
			return
		}
		if target.Attacked(p, value, common.DefenceTypeMAC, false) > 0 {
			p.LevelMagic(magic)
		}

	case common.SpellVampirism:

	case common.SpellElectricShock:
		monster, ok := args[1].(*Monster)
		if !ok {
			return
		}
		if monster == nil || !monster.IsAttackTarget(p) { // || monster.CurrentMap != CurrentMap || monster.Node == null) return;
			return
		}
		p.ElectricShock(monster, magic)

	case common.SpellMagicBooster:
		value := args[1].(int)
		buff := NewBuff(common.BuffTypeMagicBooster, p, 60000, []int{value, 6 + magic.Level})
		buff.Visible = true
		p.AddBuff(buff)
		p.LevelMagic(magic)
	case common.SpellPurification:
		target := args[1].(IMapObject)
		if target == nil || !target.IsFriendlyTarget(p) { // || target.CurrentMap != CurrentMap || target.Node == null) return;
			return
		}
		// if (Envir.Random.Next(4) > magic.Level || target.PoisonList.Count == 0) return;
		// target.ExplosionInflictedTime = 0;
		// target.ExplosionInflictedStage = 0;

		// for (int i = 0; i < target.Buffs.Count; i++)
		// {
		//     if (target.Buffs[i].Type == BuffType.Curse)
		//     {
		//         target.Buffs.RemoveAt(i);
		//         break;
		//     }
		// }

		// target.PoisonList.Clear();
		// target.OperateTime = 0;

		// if (target.ObjectID == ObjectID)
		//     Enqueue(new S.RemoveDelayedExplosion { ObjectID = target.ObjectID });
		// target.Broadcast(new S.RemoveDelayedExplosion { ObjectID = target.ObjectID });

		// LevelMagic(magic);
	case common.SpellRevelation:
		// value := args[1].(int)
		target := args[2].(IMapObject)
		if target == nil { // || target.CurrentMap != CurrentMap || target.Node == null) return;
			return
		}
		if target.GetRace() != common.ObjectTypePlayer && target.GetRace() != common.ObjectTypeMonster {
			return
		}
		// if (Envir.Random.Next(4) > magic.Level || Envir.Time < target.RevTime) return;
		// target.RevTime = Envir.Time + value * 1000;
		// target.OperateTime = 0;
		// target.BroadcastHealthChange()
		p.LevelMagic(magic)
	case common.SpellReincarnation:
	case common.SpellEntrapment:
	case common.SpellHallucination:
	case common.SpellPetEnhancer:
	case common.SpellElementalBarrier:
	case common.SpellElementalShot:
	case common.SpellDelayedExplosion:
	}
}

// Fireball 火球术
func (p *Player) Fireball(target IMapObject, magic *common.UserMagic) bool {
	if target == nil || !target.IsAttackTarget(p) {
		return false
	}
	damage := magic.GetDamage(p.GetAttackPower(int(p.MinMC), int(p.MaxMC)))
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, damage, target)
	return true
}

// Repulsion 抗拒火环
func (p *Player) Repulsion(magic *common.UserMagic) {

}

// HellFire 地狱火
func (p *Player) HellFire(magic *common.UserMagic) {

}

// GetAmulet 获取玩家身上装备的护身符
func (p *Player) GetAmulet(count int) *common.UserItem {
	for _, userItem := range p.Equipment.Items {
		if userItem == nil {
			continue
		}
		itemInfo := data.GetItemInfoByID(int(userItem.ItemID))
		if itemInfo != nil && itemInfo.Type == common.ItemTypeAmulet && int(userItem.Count) > count {
			return userItem
		}
	}
	return nil
}

// GetPoison 获取玩家身上装备的毒药
func (p *Player) GetPoison(count int) *common.UserItem {
	for _, userItem := range p.Equipment.Items {
		if userItem == nil {
			continue
		}
		itemInfo := data.GetItemInfoByID(int(userItem.ItemID))
		if itemInfo != nil && itemInfo.Type == common.ItemTypeAmulet && int(userItem.Count) > count {
			if itemInfo.Shape == 1 || itemInfo.Shape == 2 {
				return userItem
			}
		}
	}
	return nil
}

// Hiding 隐身术
func (p *Player) Hiding(magic *common.UserMagic) {
	userItem := p.GetAmulet(1)
	if userItem == nil {
		return
	}
	p.ConsumeItem(userItem, 1)
	damage := p.GetAttackPower(int(p.MinSC), int(p.MaxSC)) + (magic.Level+1)*5
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, damage)
}

// ImmortalSkin ...
func (p *Player) ImmortalSkin(magic *common.UserMagic) bool { return true }

// FireBang 爆裂火焰
func (p *Player) FireBang(magic *common.UserMagic, location common.Point) {}

// MassHiding 集体隐身术
func (p *Player) MassHiding(magic *common.UserMagic, location common.Point) bool {
	userItem := p.GetAmulet(1)
	if userItem == nil {
		return false
	}
	// int delay = Functions.MaxDistance(CurrentLocation, location) * 50 + 500; //50 MS per Step
	// DelayedAction action = new DelayedAction(DelayedType.Magic, Envir.Time + delay, this, magic, GetAttackPower(MinSC, MaxSC) / 2 + (magic.Level + 1) * 2, location);
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.Map.CompleteMagic, magic, p.GetAttackPower(int(p.MinSC), int(p.MaxSC))/2+(magic.Level+1)*2, location, p)
	return true
}

// SoulShield 幽灵盾
func (p *Player) SoulShield(magic *common.UserMagic, location common.Point) bool {
	userItem := p.GetAmulet(1)
	if userItem == nil {
		return false
	}
	// int delay = Functions.MaxDistance(CurrentLocation, location) * 50 + 500; //50 MS per Step
	// DelayedAction action = new DelayedAction(DelayedType.Magic, Envir.Time + delay, this, magic, GetAttackPower(MinSC, MaxSC) * 2 + (magic.Level + 1) * 10, location);
	delay := 500
	p.Map.ActionList.PushDelayAction(DelayedTypeMagic, delay, func() {
		p.Map.CompleteMagic(p, magic, p.GetAttackPower(int(p.MinSC), int(p.MaxSC))*2+(magic.Level+1)*10, location)
	})
	p.ConsumeItem(userItem, 1)
	return true
}

// FireWall 火墙
func (p *Player) FireWall(magic *common.UserMagic, location common.Point) {
	damage := magic.GetDamage(p.GetAttackPower(int(p.MinMC), int(p.MaxMC)))
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.Map.CompleteMagic, magic, p, damage, location)
}

// Lightning 疾光电影
func (p *Player) Lightning(magic *common.UserMagic) {
	damage := magic.GetDamage(p.GetAttackPower(int(p.MinMC), int(p.MaxMC)))
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.Map.CompleteMagic, magic, p, damage, p.CurrentLocation, p.CurrentDirection)
}

// HeavenlySword ..
func (p *Player) HeavenlySword(magic *common.UserMagic) {}

// MassHealing 群体治疗术
func (p *Player) MassHealing(magic *common.UserMagic, location common.Point) {
	value := magic.GetDamage(p.GetAttackPower(int(p.MinSC), int(p.MaxSC)))
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.Map.CompleteMagic, magic, p, value, location)
}

// ShoulderDash 野蛮冲撞
func (p *Player) ShoulderDash(magic *common.UserMagic) {}

// ThunderStorm 地狱雷光
func (p *Player) ThunderStorm(magic *common.UserMagic) {
	damage := magic.GetDamage(p.GetAttackPower(int(p.MinMC), int(p.MaxMC)))
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.Map.CompleteMagic, magic, p, damage, p.CurrentLocation)
}

// FlameDisruptor 火龙术
func (p *Player) FlameDisruptor(target IMapObject, magic *common.UserMagic) {
	if target == nil || (target.GetRace() != common.ObjectTypePlayer && target.GetRace() != common.ObjectTypeMonster) || !target.IsAttackTarget(p) {
		return
	}
	damage := magic.GetDamage(p.GetAttackPower(int(p.MinMC), int(p.MaxMC)))
	// if (!target.Undead) damage = (int)(damage * 1.5F);
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, damage, target)
}

// TurnUndead 圣言术
func (p *Player) TurnUndead(target IMapObject, magic *common.UserMagic) {
	if target == nil || target.GetRace() != common.ObjectTypeMonster || !target.IsUndead() || !target.IsAttackTarget(p) {
		return
	}
	// if (Envir.Random.Next(2) + Level - 1 <= target.Level)
	// {
	// 	target.Target = this;
	// 	return;
	// }
	// int dif = Level - target.Level + 15;
	// if (Envir.Random.Next(100) >= (magic.Level + 1 << 3) + dif)
	// {
	// 	target.Target = this;
	// 	return;
	// }
	// DelayedAction action = new DelayedAction(DelayedType.Magic, Envir.Time + 500, magic, target);
	// ActionList.Add(action);
}

// MagicBooster 深延术
func (p *Player) MagicBooster(magic *common.UserMagic) {
	bonus := 6 + magic.Level*6
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, bonus)
}

// Vampirism 嗜血术
func (p *Player) Vampirism(target IMapObject, magic *common.UserMagic) {
	if target == nil || !target.IsAttackTarget(p) {
		return
	}
	damage := magic.GetDamage(p.GetAttackPower(int(p.MinMC), int(p.MaxMC)))
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, damage, target)
}

// Purification 净化术
func (p *Player) Purification(target IMapObject, magic *common.UserMagic) {
	if target == nil || !target.IsFriendlyTarget(p) {
		return
	}
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, target)
}

// Revelation 心灵启示
func (p *Player) Revelation(target IMapObject, magic *common.UserMagic) {
	if target == nil {
		return
	}
	value := p.GetAttackPower(int(p.MinSC), int(p.MaxSC)) + magic.GetPower()
	p.ActionList.PushActionSuper(DelayedTypeMagic, p.CompleteMagic, magic, value, target)
}

// PoisonCloud 毒云
func (p *Player) PoisonCloud(magic *common.UserMagic, location common.Point) bool { return true }

// Entrapment 捕绳剑
func (p *Player) Entrapment(target IMapObject, magic *common.UserMagic) {}

// BladeAvalanche 攻破斩
func (p *Player) BladeAvalanche(magic *common.UserMagic) {}

// SlashingBurst 日闪
func (p *Player) SlashingBurst(magic *common.UserMagic) bool { return true }

// Rage 剑气爆
func (p *Player) Rage(magic *common.UserMagic) {}

// Mirroring 分身术
func (p *Player) Mirroring(magic *common.UserMagic) {}

// Blizzard 天霜冰环
func (p *Player) Blizzard(magic *common.UserMagic, location common.Point) bool { return true }

// MeteorStrike 流星火雨
func (p *Player) MeteorStrike(magic *common.UserMagic, location common.Point) bool { return true }

// IceThrust 冰焰术
func (p *Player) IceThrust(magic *common.UserMagic) {}

// ProtectionField 护身气幕
func (p *Player) ProtectionField(magic *common.UserMagic) {}

// PetEnhancer 血龙水
func (p *Player) PetEnhancer(target IMapObject, magic *common.UserMagic) bool { return true }

// TrapHexagon 困魔咒
func (p *Player) TrapHexagon(magic *common.UserMagic, target IMapObject) bool { return true }

// Reincarnation 复活术
func (p *Player) Reincarnation(magic *common.UserMagic, target IMapObject) bool { return true }

// Curse 诅咒术
func (p *Player) Curse(magic *common.UserMagic, location common.Point) bool { return true }

// SummonHolyDeva 召唤月灵
func (p *Player) SummonHolyDeva(magic *common.UserMagic) {}

// Hallucination 迷魂术
func (p *Player) Hallucination(target IMapObject, magic *common.UserMagic) bool { return true }

// EnergyShield 阴阳盾
func (p *Player) EnergyShield(target IMapObject, magic *common.UserMagic) bool { return true }

// UltimateEnhancer 无极真气
func (p *Player) UltimateEnhancer(target IMapObject, magic *common.UserMagic) bool { return true }

// Plague 瘟疫
func (p *Player) Plague(magic *common.UserMagic, location common.Point) bool { return true }

// ElectricShock 诱惑之光
func (p *Player) ElectricShock(target *Monster, magic *common.UserMagic) {}
