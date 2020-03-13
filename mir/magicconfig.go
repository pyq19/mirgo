package mir

import (
	"errors"
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/ut"
)

func addMagic(sp common.Spell, c *MagicConfig) {

	c.Spell = sp

	if c.ItemCost != nil {
		if c.ItemCostCount == 0 {
			c.ItemCostCount = 1
		}
	}

	if c.DelayAt == nil {
		c.DelayAt = DelayAt_Player
	}

	configsMap[sp] = c
}

func cloneMagic(spDest, spSrc common.Spell) {
	configsMap[spDest] = configsMap[spSrc]
}

func init() {
	add := addMagic

	// 火球术
	add(common.SpellFireBall, &MagicConfig{
		SelectType: Select_Point | Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
	})

	// 大火球
	cloneMagic(common.SpellGreatFireBall, common.SpellFireBall)
	// 雷电术
	cloneMagic(common.SpellThunderBolt, common.SpellFireBall)

	// 寒冰掌
	add(common.SpellFrostCrunch, &MagicConfig{
		SelectType: Select_Point | Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_FrostCrunch,
		Formula:    Formula_MC,
	})

	// 灵魂火符
	add(common.SpellSoulFireBall, &MagicConfig{
		SelectType: Select_Enemy | Select_Point,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_SC,
		ItemCost:   Cost_Amulet,
	})

	// 治愈术
	add(common.SpellHealing, &MagicConfig{
		SelectType: Select_Self | Select_Friend,
		TargetType: Target_Friend | Target_Self,
		Action:     Action_HealingTarget,
		Formula: func(ctx *MagicContext) int {
			return ctx.Magic.GetDamage(ctx.Player.GetAttackPower(int(ctx.Player.MinSC), int(ctx.Player.MaxSC))*2) + int(ctx.Player.Level)
		},
	})

	// 召唤骷髅
	SummonMagic(common.SpellSummonSkeleton, "BoneFamiliar", 1)
	// 召唤神兽
	SummonMagic(common.SpellSummonShinsu, "Shinsu", 5)
	// 召唤月灵
	SummonMagic(common.SpellSummonHolyDeva, "HolyDeva", 2)
	// TODO: AliveTime，SetTarget
	//
	SummonMagic1(common.SpellSummonVampire, "VampireSpider", Select_Enemy, 0)
	//
	SummonMagic1(common.SpellSummonToad, "SpittingToad", Select_Enemy, 0)
	//
	SummonMagic1(common.SpellSummonSnakes, "SnakeTotem", Select_Enemy, 0)

	// 施毒术
	add(common.SpellPoisoning, &MagicConfig{
		SelectType: Select_Enemy,
		TargetType: Select_Enemy,
		Formula:    Formula_SC,
		ItemCost:   Cost_Poison,
		Action:     Action_Poisoning,
	})

	// 隐身术
	add(common.SpellHiding, &MagicConfig{
		Formula: func(ctx *MagicContext) int {
			return ctx.Player.GetAttackPower(int(ctx.Player.MinSC), int(ctx.Player.MaxSC)) + (ctx.Magic.Level+1)*5
		},
		ItemCost: Cost_Amulet,
		Action:   Action_Hidding,
	})

	// 龙血剑法
	add(common.SpellFury, &MagicConfig{
		Action: Action_Fury,
	})

	// 爆裂火焰
	add(common.SpellFireBang, &MagicConfig{
		SelectType: Select_Point,
		TargetType: Target_Point,
		Action:     Action_RangeDamage1,
		Formula:    Formula_MC,
	})
	// 冰咆哮
	cloneMagic(common.SpellIceStorm, common.SpellFireBang)

	// 火龙术
	add(common.SpellFlameDisruptor, &MagicConfig{
		SelectType: Select_Enemy,
		TargetType: Target_Enemy,
		Formula:    Formula_MC,
		Action:     Action_DamageTarget,
		// if (!target.Undead) damage = (int)(damage * 1.5F);
	})

	// 净化术
	add(common.SpellPurification, &MagicConfig{
		SelectType: Select_Self | Select_Friend,
		TargetType: Target_Friend | Target_Self,
		Action:     Action_Purification,
	})

	// 幽灵盾
	add(common.SpellSoulShield, &MagicConfig{
		SelectType: Select_Point,
		TargetType: Target_Point,
		ItemCost:   Cost_Amulet,
		DelayAt:    DelayAt_Map,
		Action:     Action_SoulShield,
		Formula: func(ctx *MagicContext) int {
			return ctx.Player.GetAttackPower(int(ctx.Player.MinSC), int(ctx.Player.MaxSC))*2 + (ctx.Magic.Level+1)*10
		},
	})

	// 神圣战甲术
	cloneMagic(common.SpellBlessedArmour, common.SpellSoulShield)

	// 诱惑之光
	add(common.SpellElectricShock, &MagicConfig{
		SelectType: Select_Enemy,
		TargetType: Target_Enemy,
		Formula:    Formula_MC,
		Action:     Action_DamageTarget,
	})
}

func Action_ElectricShock(ctx *MagicContext) bool {
	if ut.RandomNext(4-ctx.Magic.Level) > 0 {
		if ut.RandomNext(2) == 0 {
			return true
		}
		return false
	}
	/*
		if (target.Master == this)
		{
			target.ShockTime = Envir.Time + (magic.Level * 5 + 10) * 1000;
			target.Target = null;
			return;
		}

		if (Envir.Random.Next(2) > 0)
		{
			target.ShockTime = Envir.Time + (magic.Level * 5 + 10) * 1000;
			target.Target = null;
			return;
		}

		if (target.Level > Level + 2 || !target.Info.CanTame) return;

		if (Envir.Random.Next(Level + 20 + magic.Level * 5) <= target.Level + 10)
		{
			if (Envir.Random.Next(5) > 0 && target.Master == null)
			{
				target.RageTime = Envir.Time + (Envir.Random.Next(20) + 10) * 1000;
				target.Target = null;
			}
			return;
		}

		if (Pets.Count(t => !t.Dead) >= magic.Level + 2) return;
		int rate = (int)(target.MaxHP / 100);
		if (rate <= 2) rate = 2;
		else rate *= 2;

		if (Envir.Random.Next(rate) != 0) return;
		//else if (Envir.Random.Next(20) == 0) target.Die();

		if (target.Master != null)
		{
			target.SetHP(target.MaxHP / 10);
			target.Master.Pets.Remove(target);
		}
		else if (target.Respawn != null)
		{
			target.Respawn.Count--;
			Envir.MonsterCount--;
			CurrentMap.MonsterCount--;
			target.Respawn = null;
		}

		target.Master = this;
		//target.HealthChanged = true;
		target.BroadcastHealthChange();
		Pets.Add(target);
		target.Target = null;
		target.RageTime = 0;
		target.ShockTime = 0;
		target.OperateTime = 0;
		target.MaxPetLevel = (byte)(1 + magic.Level * 2);
		//target.TameTime = Envir.Time + (Settings.Minute * 60);

		target.Broadcast(new S.ObjectName { ObjectID = target.ObjectID, Name = target.Name });
	*/
	return true
}

func Action_SoulShield(ctx *MagicContext) bool {
	const damageRange = 3
	var loc common.Point
	if ctx.Target == nil {
		loc = ctx.TargetPoint
	} else {
		loc = ctx.Target.GetPoint()
	}

	buffType := common.BuffTypeSoulShield
	if ctx.Magic.Spell == common.SpellBlessedArmour {
		buffType = common.BuffTypeBlessedArmour
	}

	ctx.Map.RangeObject(loc, damageRange, func(o IMapObject) bool {
		switch o.GetRace() {
		case common.ObjectTypePlayer, common.ObjectTypeMonster:
			if o.IsFriendlyTarget(ctx.Player) {
				buff := NewBuff(buffType, ctx.Player, ctx.Damage*1000, []int32{int32(o.GetLevel()/7) + 4})
				o.AddBuff(buff)
			}
		}
		return true
	})

	return true
}

func Action_Purification(ctx *MagicContext) bool {
	if ctx.Target.GetMap() != ctx.Player.Map {
		return false
	}

	// target.ExplosionInflictedTime = 0;
	// target.ExplosionInflictedStage = 0;

	// for (int i = 0; i < target.Buffs.Count; i++)
	// {
	// 	if (target.Buffs[i].Type == BuffType.Curse)
	// 	{
	// 		target.Buffs.RemoveAt(i);
	// 		break;
	// 	}
	// }

	// target.PoisonList.Clear();
	// target.OperateTime = 0;

	// if (target.ObjectID == ObjectID)
	// 	Enqueue(new S.RemoveDelayedExplosion { ObjectID = target.ObjectID });
	// target.Broadcast(new S.RemoveDelayedExplosion { ObjectID = target.ObjectID });

	return true
}

func Action_RangeDamage1(ctx *MagicContext) bool {
	const damageRange = 1
	var loc common.Point
	if ctx.Target == nil {
		loc = ctx.TargetPoint
	} else {
		loc = ctx.Target.GetPoint()
	}

	ctx.Map.RangeObject(loc, damageRange, func(o IMapObject) bool {
		switch o.GetRace() {
		case common.ObjectTypePlayer, common.ObjectTypeMonster:
			if o.IsAttackTarget(ctx.Player) {
				o.Attacked(ctx.Player, ctx.Damage, common.DefenceTypeMAC, false)
			}
		}
		return true
	})

	return true
}

func Action_Fury(ctx *MagicContext) bool {
	buff := NewBuff(common.BuffTypeFury, ctx.Player, 60000+ctx.Magic.Level*10000, []int32{4})
	buff.Visible = true
	ctx.Player.AddBuff(buff)
	return true
}

func Action_Hidding(ctx *MagicContext) bool {
	p := ctx.Player

	for e := p.BuffList.List.Front(); e != nil; e = e.Next() {
		if e.Value.(*Buff).BuffType == common.BuffTypeHiding {
			return false
		}
	}
	buff := NewBuff(common.BuffTypeHiding, p, 1000*ctx.Damage, []int32{})
	p.AddBuff(buff)
	return true
}

func Action_HealingTarget(ctx *MagicContext) bool {
	// target.HealAmount = (ushort)Math.Min(ushort.MaxValue, target.HealAmount + value);

	if ctx.Target == nil {
		ctx.Player.ChangeHP(ctx.Damage)
	} else {
		target := ctx.Target.(ILifeObject)
		target.ChangeHP(ctx.Damage)
	}
	return true
}

func Action_FrostCrunch(ctx *MagicContext) bool {
	target := ctx.Target
	p := ctx.Player

	if target.Attacked(p, ctx.Damage, common.DefenceTypeMAC, false) > 0 {
		var tmp1 int
		var tmp2 int
		var duration int
		if target.GetRace() == common.ObjectTypePlayer {
			tmp1 = 2
			tmp2 = 100
			duration = 4
		} else {
			tmp1 = 10
			tmp2 = 20
			duration = 5 + ut.RandomNext(5)
		}
		if int(p.Level)+tmp1 >= target.GetLevel() && ut.RandomNext(tmp2) <= ctx.Magic.Level {
			target.ApplyPoison(NewPoison(duration, p, common.PoisonTypeSlow, 1000, 0), p)
			// TODO // target.OperateTime = 0;
		}
		if target.GetRace() == common.ObjectTypePlayer {
			tmp1 = 2
			tmp2 = 100
			duration = 2
		} else {
			tmp1 = 10
			tmp2 = 40
			duration = 5 + int(p.Freezing)
		}
		if int(p.Level)+tmp1 >= target.GetLevel() && ut.RandomNext(tmp2) <= ctx.Magic.Level {
			target.ApplyPoison(NewPoison(duration, p, common.PoisonTypeFrozen, 1000, 0), p)
			// TODO // target.OperateTime = 0;
		}
		return true
	}
	return false
}

func Action_Poisoning(ctx *MagicContext) bool {

	duration := (ctx.Damage * 2) + ((ctx.Magic.Level + 1) * 7)
	value := ctx.Damage/15 + ctx.Magic.Level + 1 + ut.RandomNext(int(ctx.Player.PoisonAttack))
	switch ctx.Item.Info.Shape {
	case 1:
		ctx.Target.ApplyPoison(NewPoison(duration, ctx.Player, common.PoisonTypeGreen, 2000, value), ctx.Player)
	case 2:
		ctx.Target.ApplyPoison(NewPoison(duration, ctx.Player, common.PoisonTypeRed, 2000, 0), ctx.Player)
	}
	return true
}

type summonData struct {
	MonsterName string
	Pets        []*Monster
	SpellMap    *Map
}

func SummonMagic1(sp common.Spell, sumname string, sel MagicSelectType, amuletCount uint32) {
	cfg := &MagicConfig{
		SelectType:    sel,
		TargetType:    Target_None,
		ItemCost:      Cost_Amulet,
		ItemCostCount: amuletCount,
		BeforeAction:  SummonMagic_BeforeAction,
		Action:        SummonMagic_Action,
		Data: &summonData{
			MonsterName: sumname,
		},
	}

	if amuletCount == 0 {
		cfg.ItemCost = nil
	}

	addMagic(sp, cfg)
}

func SummonMagic(sp common.Spell, sumname string, amuletCount int) {
	SummonMagic1(sp, sumname, Select_Enemy, 0)
}

func SummonMagic_BeforeAction(ctx *MagicContext) (error, uint32) {
	sumdata := ctx.Config.Data.(*summonData)
	p := ctx.Player
	for i := range p.Pets {
		if p.Pets[i].GetName() == sumdata.MonsterName {
			m := p.Pets[i].(*Monster)
			p.ActionList.PushAction(DelayedTypeRecall, m.PetRecall)
			return nil, 0
		}
	}

	if len(p.Pets) > 1 {
		return nil, 0
	}

	if checkMagicItemCost(ctx) != nil {
		return errors.New("缺少释放技能的道具"), 0
	}

	sumdata.Pets = []*Monster{}
	sumdata.SpellMap = p.Map

	monsterInfo := data.GetMonsterInfoByName(sumdata.MonsterName)

	dir := int(p.CurrentDirection) + 4
	if dir > 8 {
		dir -= 8
	}
	monster := NewMonster(p.Map, p.GetPoint().NextPoint(common.MirDirection(dir), 1), monsterInfo)
	monster.PetLevel = uint16(ctx.Magic.Level)
	monster.Master = p
	monster.ActionTime = time.Now().Add(time.Duration(1000) * time.Millisecond)
	monster.CurrentLocation = p.GetFrontPoint()

	sumdata.Pets = append(sumdata.Pets, monster)

	DelayAt_Map(ctx, func() { completeMagic(ctx) })

	return nil, 0
}

func SummonMagic_Action(ctx *MagicContext) bool {
	if ctx.Player.IsDead() {
		return false
	}

	sumdata := ctx.Config.Data.(*summonData)

	for _, v := range sumdata.Pets {
		cell := sumdata.SpellMap.GetCell(v.CurrentLocation)
		if cell.IsValid() {
			v.Master = ctx.Player
			ctx.Player.Pets = append(ctx.Player.Pets, v)
			v.Spawn()
		}
	}

	return true
}
