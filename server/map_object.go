package main

import (
	"github.com/yenkeia/mirgo/common"
)

type MapObject struct {
	ID               uint32
	Name             string
	Map              *Map
	CurrentLocation  *common.Point
	CurrentDirection common.MirDirection
	Level            uint16
	Health           uint32 // 当前生命值
	MaxMaxHealth     uint32 // 最大生命值
	MinAC            uint16 // 物理防御力
	MaxAC            uint16
	MinMAC           uint16 // 魔法防御力
	MaxMAC           uint16
	MinDC            uint16 // 攻击力
	MaxDC            uint16
	MinMC            uint16 // 魔法力
	MaxMC            uint16
	MinSC            uint16 // 道术力
	MaxSC            uint16
}

type Poison struct {
	ObjectID   uint32
	PoisonType common.PoisonType
	Value      int
	Duration
	Time
	TickTime
	TickSpeed
}

type Buff struct {
	ObjectID uint32
}
