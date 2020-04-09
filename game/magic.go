package game

import (
	"errors"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/proto/server"
	"github.com/yenkeia/mirgo/game/util"
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
type MagicItemCost func(ctx *MagicContext) *cm.UserItem

func Cost_Amulet(ctx *MagicContext) *cm.UserItem {
	return ctx.Player.GetAmulet(ctx.Config.ItemCostCount)
}

func Cost_Poison(ctx *MagicContext) *cm.UserItem {
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
type MagicBeforeAction func(ctx *MagicContext) (error, uint32)

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
	Spell         cm.Spell
	Formula       MagicDamageFormula
	SelectType    MagicSelectType
	TargetType    MagicTargetType
	DelayAt       MagicDelayAt
	ItemCost      MagicItemCost
	ItemCostCount uint32
	BeforeAction  MagicBeforeAction
	Action        MagicAction
	Data          interface{}
}

type MagicContext struct {
	Spell            cm.Spell
	Target           IMapObject
	Magic            *cm.UserMagic // 当前施法
	Player           *Player       // 施法玩家
	Item             *cm.UserItem  // 消耗的物品
	Map              *Map          // 施法的地图
	TargetPoint      cm.Point      // 施法目标位置
	PlayerSpellPoint cm.Point      // 施法时候玩家的位置

	Damage int
	Config *MagicConfig
}

var configsMap = map[cm.Spell]*MagicConfig{}

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

func startMagic(ctx *MagicContext) (err error, targetid uint32) {
	cfg := configsMap[ctx.Spell]

	if cfg.SelectType&Select_Point == Select_Point {

	} else {
		if cfg.SelectType&Select_Enemy == Select_Enemy {
			if ctx.Target == nil || !ctx.Target.IsAttackTarget(ctx.Player) {
				return errors.New("请选择一个目标"), 0
			}
			targetid = ctx.Target.GetID()
		}

		if cfg.SelectType&Select_Friend == Select_Friend {
			if ctx.Target == nil || !ctx.Target.IsFriendlyTarget(ctx.Player) {
				return errors.New("请选择一个目标"), 0
			}
			targetid = ctx.Target.GetID()
		}
	}

	ctx.Map = ctx.Player.Map
	ctx.Config = cfg

	if cfg.BeforeAction != nil {
		return cfg.BeforeAction(ctx)
	} else {

		if checkMagicItemCost(ctx) != nil {
			return errors.New("缺少释放技能的道具"), 0
		}

		if cfg.Formula != nil {
			ctx.Damage = cfg.Formula(ctx)
		}

		cfg.DelayAt(ctx, func() { completeMagic(ctx) })
		return nil, targetid
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

// GetAmulet 获取玩家身上装备的护身符
func (p *Player) GetAmulet(count uint32) *cm.UserItem {
	for _, userItem := range p.Equipment.Items {
		if userItem == nil {
			continue
		}
		if userItem.Info.Type == cm.ItemTypeAmulet && userItem.Count >= count {
			return userItem
		}
	}
	return nil
}

// GetPoison 获取玩家身上装备的毒药
func (p *Player) GetPoison(count uint32) *cm.UserItem {
	for _, userItem := range p.Equipment.Items {
		if userItem == nil {
			continue
		}
		if userItem.Info.Type == cm.ItemTypeAmulet && userItem.Count >= count {
			if userItem.Info.Shape == 1 || userItem.Info.Shape == 2 {
				return userItem
			}
		}
	}
	return nil
}

// GetMagic ...
func (p *Player) GetMagic(spell cm.Spell) *cm.UserMagic {
	for i := range p.Magics {
		userMagic := p.Magics[i]
		if userMagic.Spell == spell {
			return userMagic
		}
	}
	return nil
}

// LevelMagic ...
func (p *Player) LevelMagic(magic *cm.UserMagic) {
	exp := util.RandomNext(3) + 1

	magicLevel := 0
	magicNeed := 0
	oldLevel := magic.Level

	switch oldLevel {
	case 0:
		magicLevel = magic.Info.Level1
		magicNeed = magic.Info.Need1
	case 1:
		magicLevel = magic.Info.Level2
		magicNeed = magic.Info.Need2
	case 2:
		magicLevel = magic.Info.Level3
		magicNeed = magic.Info.Need3
	}

	if int(p.Level) < magicLevel {
		return
	}
	magic.Experience += exp
	if magic.Experience >= magicNeed {
		magic.Level++
		magic.Experience = (magic.Experience - magicNeed)
		p.RefreshStats()
	}

	if oldLevel != magic.Level {
		p.Enqueue(&server.MagicDelay{Spell: magic.Spell, Delay: int64(magic.GetDelay())})
	}

	p.Enqueue(&server.MagicLeveled{Spell: magic.Spell, Level: uint8(magic.Level), Experience: uint16(magic.Experience)})
}
