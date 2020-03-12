package mir

import (
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
	add(common.SpellGreatFireBall, &MagicConfig{
		SelectType: Select_Point | Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
	})

	// 寒冰掌
	add(common.SpellFrostCrunch, &MagicConfig{
		SelectType: Select_Point | Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_FrostCrunch,
		Formula:    Formula_MC,
	})

	// 雷电术
	add(common.SpellThunderBolt, &MagicConfig{
		SelectType: Select_Enemy,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
	})

	// 灵魂火符
	add(common.SpellSoulFireBall, &MagicConfig{
		SelectType: Select_Enemy | Select_Point,
		TargetType: Target_Enemy,
		Action:     Action_DamageTarget,
		Formula:    Formula_MC,
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

	// 施毒术
	add(common.SpellPoisoning, &MagicConfig{
		SelectType: Select_Enemy,
		TargetType: Select_Enemy,
		Formula:    Formula_SC,
		ItemCost:   Cost_Poison,
		Action:     Action_Poisoning,
	})

	// 隐身术
	add(common.SpellPoisoning, &MagicConfig{
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
}

func Action_Fury(ctx *MagicContext) bool {
	buff := NewBuff(common.BuffTypeFury, ctx.Player, 60000+ctx.Magic.Level*10000, []int{4})
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
	buff := NewBuff(common.BuffTypeHiding, p, 1000*ctx.Damage, []int{})
	p.AddBuff(buff)
	return true
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

func SummonMagic(sp common.Spell, sumname string, amuletCount int) {
	addMagic(sp, &MagicConfig{
		SelectType:    Select_None,
		TargetType:    Target_None,
		ItemCost:      Cost_Amulet,
		ItemCostCount: amuletCount,
		BeforeAction:  SummonMagic_BeforeAction,
		Action:        SummonMagic_Action,
		Data: &summonData{
			MonsterName: sumname,
		},
	})
}

func SummonMagic_BeforeAction(ctx *MagicContext) (bool, uint32) {
	sumdata := ctx.Config.Data.(*summonData)
	p := ctx.Player
	for i := range p.Pets {
		if p.Pets[i].GetName() == sumdata.MonsterName {
			m := p.Pets[i].(*Monster)
			p.ActionList.PushAction(DelayedTypeRecall, m.PetRecall)
			return false, 0
		}
	}

	if len(p.Pets) > 1 {
		return false, 0
	}

	if checkMagicItemCost(ctx) != nil {
		return false, 0
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

	return true, 0
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
