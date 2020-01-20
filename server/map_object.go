package main

import (
	"github.com/yenkeia/mirgo/common"
	"time"
)

type MapObject struct {
	ID               uint32
	Name             string
	NameColor        common.Color
	Map              *Map
	CurrentLocation  common.Point
	CurrentDirection common.MirDirection
	Poisons          []*Poison
	Buffs            []*Buff
}

type Poison struct {
	ObjectID   uint32
	PoisonType common.PoisonType
	Value      int       // 效果总数
	NextTime   time.Time // 下次生效时间
	Duration   time.Time // 两次生效时间间隔
	ExpireTime time.Time // 结束时间
	TickNum    int       // 总共跳几次
	TickTime   int       // 当前第几跳
}

type Buff struct {
	ObjectID   uint32
	BuffType   common.BuffType
	Visible    bool      // 是否可见
	Infinite   bool      // 是否永久
	Value      int       // public int[] Values
	NextTime   time.Time // 下次生效时间
	Duration   time.Time // 两次生效时间间隔
	ExpireTime time.Time // 过期时间️
	TickNum    int       // 总共跳几次
	TickTime   int       // 当前第几跳
}

type IMapObject interface {
	GetID() uint32
	GetRace() common.ObjectType
	GetCoordinate() string
	GetPoint() common.Point
	GetCell() *Cell
	Broadcast(interface{})
	GetDirection() common.MirDirection
	GetInfo() interface{}
}
