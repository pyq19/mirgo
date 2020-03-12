package mir

import (
	"time"

	"github.com/yenkeia/mirgo/common"
)

func addMagic(sp common.Spell, c *MagicConfig) {

	c.Spell = sp
	if c.ItemCost == nil {
		c.ItemCost = Cost_None
	}
	if c.DelayAt == nil {
		c.DelayAt = DelayAt_Player
	}

	configsMap[sp] = c
}

func init() {
	add := addMagic

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
		Formula: func(ctx *MagicContext) int {
			return ctx.Magic.GetDamage(ctx.Player.GetAttackPower(int(ctx.Player.MinSC), int(ctx.Player.MaxSC))*2) + int(ctx.Player.Level)
		},
	})

	SummonMagic(common.SpellSummonSkeleton, "BoneFamiliar")
	SummonMagic(common.SpellSummonShinsu, "Shinsu")
}

type summonData struct {
	MonsterName string
	Pets        []*Monster
	SpellMap    *Map
}

func SummonMagic(sp common.Spell, sumname string) {
	addMagic(sp, &MagicConfig{
		SelectType:   Select_None,
		TargetType:   Target_None,
		ItemCost:     Cost_Amulet,
		BeforeAction: SummonMagic_BeforeAction,
		Action:       SummonMagic_Action,
		Data: &summonData{
			MonsterName: sumname,
		},
	})
}

func SummonMagic_BeforeAction(ctx *MagicContext) {
	sumdata := ctx.Config.Data.(*summonData)
	p := ctx.Player
	for i := range p.Pets {
		if p.Pets[i].GetName() == sumdata.MonsterName {
			m := p.Pets[i].(*Monster)
			p.ActionList.PushAction(DelayedTypeRecall, m.PetRecall)
			return
		}
	}

	if len(p.Pets) > 1 {
		return
	}

	if checkMagicItemCost(ctx) != nil {
		return
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
}

func SummonMagic_Action(ctx *MagicContext) {
	if ctx.Player.IsDead() {
		return
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
}
