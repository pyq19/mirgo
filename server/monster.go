package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type Monster struct {
	MapObject
	RespawnID   int
	Image       common.Monster
	AI          int
	Effect      int
	Poison      common.PoisonType
	Light       uint8
	MaxHP       uint32
	MinAC       uint16
	MaxAC       uint16
	MinMAC      uint16
	MaxMAC      uint16
	MinDC       uint16
	MaxDC       uint16
	MinMC       uint16
	MaxMC       uint16
	MinSC       uint16
	MaxSC       uint16
	Accuracy    uint8
	Agility     uint8
	MoveSpeed   uint16
	AttackSpeed int32
}

func (m *Monster) String() string {
	return fmt.Sprintf("Monster: %s, (%v), ID: %d, ptr: %p\n", m.Name, m.CurrentLocation, m.ID, m)
}

func NewMonster(mp *Map, p common.Point, mi *common.MonsterInfo, ri int) (m *Monster) {
	m = new(Monster)
	m.RespawnID = ri
	m.ID = mp.Env.NewObjectID()
	m.Name = mi.Name
	m.NameColor = common.Color{R: 255, G: 255, B: 255}
	m.Image = common.Monster(mi.Image)
	m.AI = mi.AI
	m.Effect = mi.Effect
	m.Light = uint8(mi.Light)
	m.Poison = common.PoisonTypeNone
	m.CurrentLocation = p
	m.CurrentDirection = common.MirDirection(G_Rand.RandInt(0, 7))
	m.MaxHP = 0
	m.MinAC = 0
	m.MaxAC = 0
	m.MinMAC = 0
	m.MaxMAC = 0
	m.MinDC = 0
	m.MaxDC = 0
	m.MinMC = 0
	m.MaxMC = 0
	m.MinSC = 0
	m.MaxSC = 0
	m.Accuracy = 0
	m.Agility = 0
	m.MoveSpeed = 0
	m.AttackSpeed = 0
	return m
}

func (m *Monster) GetID() uint32 {
	return m.ID
}

func (m *Monster) GetRace() common.ObjectType {
	return common.ObjectTypeMonster
}

func (m *Monster) GetCoordinate() string {
	return m.GetPoint().Coordinate()
}

func (m *Monster) GetPoint() common.Point {
	return m.CurrentLocation
}

func (m *Monster) GetCell() *Cell {
	return m.Map.GetCell(m.GetCoordinate())
}

func (m *Monster) GetDirection() common.MirDirection {
	return m.CurrentDirection
}

func (m *Monster) GetInfo() interface{} {
	res := &server.ObjectMonster{
		ObjectID:          m.ID,
		Name:              m.Name,
		NameColor:         m.NameColor.ToInt32(),
		Location:          m.GetPoint(),
		Image:             m.Image,
		Direction:         m.GetDirection(),
		Effect:            uint8(m.Effect),
		AI:                uint8(m.AI),
		Light:             m.Light,
		Dead:              m.IsDead(),
		Skeleton:          m.IsSkeleton(),
		Poison:            m.Poison,
		Hidden:            m.IsHidden(),
		ShockTime:         0,     // TODO
		BindingShotCenter: false, // TODO
		Extra:             false, // TODO
		ExtraByte:         0,     // TODO
	}
	return res
}

func (m *Monster) Broadcast(msg interface{}) {

}

func (m *Monster) IsDead() bool {
	return false
}

func (m *Monster) IsSkeleton() bool {
	return false
}

func (m *Monster) IsHidden() bool {
	return false
}

func (m *Monster) Process() {

}

func (m *Monster) isAttackTarget(attacker *Player) bool {
	return true
}

func (m *Monster) attacked(attacker *Player, finalDamage int, defenceType common.DefenceType, damageWeapon bool) {

}
