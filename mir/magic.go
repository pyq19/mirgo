package mir

import (
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/ut"
)

// 技能选择类型
// 玩家在按施法键之后，选择一个目标或者目标点。或者没有目标
type MagicSelectType int

const (
	Select_None   MagicSelectType = 0
	Select_Enemy                  = 1 << 0
	Select_Point                  = 1 << 1
	Select_Friend                 = 1 << 2
	Select_Self                   = 1 << 3
)

// 技能作用类型
// 比如 灵魂火符可以对一个点施法，但最终只能作用在一个敌人身上
type MagicTargetType int

const (
	Target_None   MagicTargetType = 0
	Target_Enemy                  = 1 << 0
	Target_Point                  = 1 << 1
	Target_Friend                 = 1 << 2
	Target_Self                   = 1 << 3
)

// 技能额外消耗类型
type MagicItemCost func(ctx *MagicContext) (*common.UserItem, int)

func Cost_None(ctx *MagicContext) (*common.UserItem, int) {
	return nil, 0
}

func Cost_Amulet(ctx *MagicContext) (*common.UserItem, int) {
	return ctx.Player.GetAmulet(1), 1
}

func Cost_Poison(ctx *MagicContext) (*common.UserItem, int) {
	return ctx.Player.GetPoison(1), 1
}

// 技能延迟在哪里执行
type MagicDelayAt func(ctx *MagicContext, f func())

func DelayAt_Map(ctx *MagicContext, f func()) {
	ctx.Player.Map.ActionList.PushAction(DelayedTypeMagic, f)
}

func DelayAt_Player(ctx *MagicContext, f func()) {
	ctx.Player.ActionList.PushAction(DelayedTypeMagic, f)
}

// 技能最终作用函数
type MagicAction func(ctx *MagicContext)

// 伤害值计算公式
type MagicDamageFormula func(ctx *MagicContext) int

func Formula_SC(ctx *MagicContext) int {
	return ctx.Magic.GetDamage(ctx.Player.GetAttackPower(int(ctx.Player.MinSC), int(ctx.Player.MaxSC)))
}

func Formula_MC(ctx *MagicContext) int {
	return ctx.Magic.GetDamage(ctx.Player.GetAttackPower(int(ctx.Player.MinMC), int(ctx.Player.MaxMC)))
}

type MagicConfig struct {
	Spell      common.Spell
	Formula    MagicDamageFormula
	SelectType MagicSelectType
	TargetType MagicTargetType
	DelayAt    MagicDelayAt
	ItemCost   MagicItemCost
	Action     MagicAction
}

type MagicContext struct {
	Spell  common.Spell
	Target IMapObject
	Magic  *common.UserMagic
	Player *Player

	Damage int
	Config *MagicConfig
}

var configsMap = map[common.Spell]*MagicConfig{}

func startMagic(ctx *MagicContext) (cast bool, targetid uint32) {
	cfg := configsMap[ctx.Spell]

	if cfg.SelectType&Select_Enemy == Select_Enemy {
		if ctx.Target == nil || !ctx.Target.IsAttackTarget(ctx.Player) {
			return false, 0
		}
		targetid = ctx.Target.GetID()
	}

	if cfg.SelectType&Select_Friend == Select_Friend {
		if ctx.Target == nil || !ctx.Target.IsFriendlyTarget(ctx.Player) {
			return false, 0
		}
		targetid = ctx.Target.GetID()
	}

	if cfg.SelectType&Select_Point == Select_Point {
		// TODO
	}

	item, count := cfg.ItemCost(ctx)

	if item != nil {
		ctx.Player.ConsumeItem(item, count)
	}

	ctx.Config = cfg
	ctx.Damage = cfg.Formula(ctx)

	ctx.Config.DelayAt(ctx, func() { completeMagic(ctx) })

	return true, targetid
}

func completeMagic(ctx *MagicContext) {
	cfg := ctx.Config

	if cfg.TargetType&Target_Enemy == Target_Enemy {
		if ctx.Target == nil || !ctx.Target.IsAttackTarget(ctx.Player) {
			return
		}
	}

	if cfg.SelectType&Target_Friend == Target_Friend {
		if ctx.Target == nil || !ctx.Target.IsFriendlyTarget(ctx.Player) {
			return
		}
	}

	if cfg.SelectType&Target_Point == Target_Point {
		// TODO
	}

	cfg.Action(ctx)

	ctx.Player.LevelMagic(ctx.Magic)
}

func Action_DamageTarget(ctx *MagicContext) {
	ctx.Target.Attacked(ctx.Player, ctx.Damage, common.DefenceTypeMAC, false)
}

func Action_HealingTarget(ctx *MagicContext) {
	// target.HealAmount = (ushort)Math.Min(ushort.MaxValue, target.HealAmount + value);

	if ctx.Target == nil {
		ctx.Player.ChangeHP(ctx.Damage)
	} else {
		target := ctx.Target.(ILifeObject)
		target.ChangeHP(ctx.Damage)
	}
}

func Action_FrostCrunch(ctx *MagicContext) {
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
	}
}
