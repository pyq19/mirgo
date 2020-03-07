package mir

import (
	"time"

	"github.com/yenkeia/mirgo/common"
)

type Poison struct {
	ObjectID   uint32
	Owner      IMapObject
	PoisonType common.PoisonType
	Value      int           // 效果总数
	NextTime   time.Time     // 下次生效时间
	Duration   time.Duration // 两次生效时间间隔
	TickNum    int           // 总共跳几次
	TickTime   int           // 当前第几跳
}

// NewPoison caster 释放者, value 总伤害, typ 毒类型, duration 两次间隔, tickNum 总共跳几次
func NewPoison(id uint32, caster IMapObject, value int, typ common.PoisonType, duration time.Duration, tickNum int) *Poison {
	return &Poison{
		ObjectID:   id,
		Owner:      caster,
		PoisonType: typ,
		Value:      value,
		NextTime:   time.Now().Add(duration),
		Duration:   duration,
		TickNum:    tickNum,
		TickTime:   0,
	}
}

type Buff struct {
	ObjectID   uint32
	BuffType   common.BuffType
	Visible    bool      // 是否可见
	Infinite   bool      // 是否永久
	Values     int       // public int[] Values
	ExpireTime time.Time // 过期时间️
}

func NewBuff(id uint32, typ common.BuffType, value int, expire time.Time) *Buff {
	return &Buff{
		ObjectID:   id,
		BuffType:   typ,
		Visible:    false,
		Infinite:   false,
		Values:     value,
		ExpireTime: expire,
	}
}

type BaseStats struct {
	MinAC    uint16
	MaxAC    uint16
	MinMAC   uint16
	MaxMAC   uint16
	MinDC    uint16
	MaxDC    uint16
	MinMC    uint16
	MaxMC    uint16
	MinSC    uint16
	MaxSC    uint16
	Accuracy uint8
	Agility  uint8
}

type IProcessObject interface {
	GetID() uint32
	Process(dt time.Duration)
}

type IMapObject interface {
	GetID() uint32
	GetMap() *Map
	GetName() string
	GetRace() common.ObjectType
	GetPoint() common.Point
	GetCell() *Cell
	Broadcast(interface{})
	GetDirection() common.MirDirection
	GetInfo() interface{}
	GetBaseStats() BaseStats
	IsAttackTarget(IMapObject) bool
	IsFriendlyTarget(IMapObject) bool
	IsDead() bool
	IsUndead() bool
	IsBlocking() bool
	AttackMode() common.AttackMode
	AddBuff(*Buff)
	ApplyPoison(*Poison, IMapObject)
	AddPlayerCount(n int)
	GetPlayerCount() int
	Attacked(attacker IMapObject, damage int, dtype common.DefenceType, damageWeapon bool) int
}

type MapObject struct {
	ID               uint32
	Name             string
	NameColor        common.Color
	Map              *Map
	CurrentLocation  common.Point
	CurrentDirection common.MirDirection
	Poisons          []*Poison
	Buffs            []*Buff
	Dead             bool
	PlayerCount      int // 记录在DataRange内有多少个玩家
}
