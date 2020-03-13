package mir

import (
	"errors"

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
type MagicItemCost func(ctx *MagicContext) *common.UserItem

func Cost_Amulet(ctx *MagicContext) *common.UserItem {
	return ctx.Player.GetAmulet(ctx.Config.ItemCostCount)
}

func Cost_Poison(ctx *MagicContext) *common.UserItem {
	return ctx.Player.GetPoison(ctx.Config.ItemCostCount)
}

// 技能延迟在哪里执行
type MagicDelayAt func(ctx *MagicContext, f func())

func DelayAt_Map(ctx *MagicContext, f func()) {
	ctx.Player.Map.ActionList.PushAction(DelayedTypeMagic, f)
}

func DelayAt_Player(ctx *MagicContext, f func()) {
	ctx.Player.ActionList.PushAction(DelayedTypeMagic, f)
}

//
type MagicBeforeAction func(ctx *MagicContext) (bool, uint32)

// 技能最终作用函数
type MagicAction func(ctx *MagicContext) bool

// 伤害值计算公式
type MagicDamageFormula func(ctx *MagicContext) int

func Formula_SC(ctx *MagicContext) int {
	return ctx.Magic.GetDamage(ctx.Player.GetAttackPower(int(ctx.Player.MinSC), int(ctx.Player.MaxSC)))
}

func Formula_MC(ctx *MagicContext) int {
	return ctx.Magic.GetDamage(ctx.Player.GetAttackPower(int(ctx.Player.MinMC), int(ctx.Player.MaxMC)))
}

type MagicConfig struct {
	Spell         common.Spell
	Formula       MagicDamageFormula
	SelectType    MagicSelectType
	TargetType    MagicTargetType
	DelayAt       MagicDelayAt
	ItemCost      MagicItemCost
	ItemCostCount int
	BeforeAction  MagicBeforeAction
	Action        MagicAction
	Data          interface{}
}

type MagicContext struct {
	Spell            common.Spell
	Target           IMapObject
	Magic            *common.UserMagic // 当前施法
	Player           *Player           // 施法玩家
	Item             *common.UserItem  // 消耗的物品
	Map              *Map              // 施法的地图
	TargetPoint      common.Point      // 施法目标位置
	PlayerSpellPoint common.Point      // 施法时候玩家的位置

	Damage int
	Config *MagicConfig
}

var configsMap = map[common.Spell]*MagicConfig{}

func checkMagicItemCost(ctx *MagicContext) error {
	if ctx.Config.ItemCost == nil || ctx.Config.ItemCostCount == 0 {
		return nil
	}

	item := ctx.Config.ItemCost(ctx)

	if item == nil {
		return errors.New("没有施法道具")
	} else {
		ctx.Player.ConsumeItem(item, ctx.Config.ItemCostCount)
	}

	ctx.Item = item

	return nil
}

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

	ctx.Map = ctx.Player.Map
	ctx.Config = cfg

	if cfg.BeforeAction != nil {
		return cfg.BeforeAction(ctx)
	} else {

		if checkMagicItemCost(ctx) != nil {
			return false, 0
		}

		if cfg.Formula != nil {
			ctx.Damage = cfg.Formula(ctx)
		}

		cfg.DelayAt(ctx, func() { completeMagic(ctx) })
		return true, targetid
	}
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

	if cfg.Action(ctx) {
		ctx.Player.LevelMagic(ctx.Magic)
	}
}

func Action_DamageTarget(ctx *MagicContext) bool {
	return ctx.Target.Attacked(ctx.Player, ctx.Damage, common.DefenceTypeMAC, false) > 0
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
