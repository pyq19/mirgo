package main

import (
	"errors"
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type Monster struct {
	Respawn *Respawn
	MapObject
	Image  common.Monster
	AI     int
	Effect int
	Poison common.PoisonType
}

func (m *Monster) String() string {
	return fmt.Sprintf("Monster: %s, (%v), ID: %d, ptr: %p\n", m.Name, m.CurrentLocation, m.ID, m)
}

func NewMonster(r *Respawn) (m *Monster, err error) {
	m = new(Monster)
	m.Respawn = r
	mi := r.Map.Env.GameDB.GetMonsterInfoByID(r.Info.MonsterID)
	if mi == nil {
		return nil, errors.New("new monster error")
	}
	m.ID = r.Map.Env.NewObjectID()
	m.Name = mi.Name
	m.NameColor = common.Color{R: 255, G: 255, B: 255}
	m.Image = common.Monster(mi.Image)
	m.AI = mi.AI
	m.Effect = mi.Effect
	m.Light = uint8(mi.Light)
	m.Poison = common.PoisonTypeNone
	p, err := r.Map.GetValidPoint(r.Info.LocationX, r.Info.LocationY, r.Info.Spread)
	if err != nil {
		return nil, err
	}
	m.CurrentLocation = p
	m.CurrentDirection = common.MirDirection(G_Rand.RandInt(0, 7))
	c := r.Map.GetCell(p.Coordinate())
	c.AddObject(m)
	return m, nil
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
