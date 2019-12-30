package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
)

type Monster struct {
	MapObject
	Respawn *Respawn
}

func (m *Monster) String() string {
	return fmt.Sprintf("Monster: %s, (%v), ID: %d, ptr: %p\n", m.Name, m.CurrentLocation, m.ID, m)
}

func NewMonster(r *Respawn) (m *Monster, err error) {
	m = new(Monster)
	m.Respawn = r
	mi := r.Map.Env.GameDB.GetMonsterInfoByID(r.Info.MonsterID)
	m.ID = r.Map.Env.NewObjectID()
	m.Name = mi.Name
	p, err := r.Map.GetValidPoint(r.Info.LocationX, r.Info.LocationY, r.Info.Spread)
	if err != nil {
		return nil, err
	}
	m.CurrentLocation = p
	m.CurrentDirection = common.MirDirection(G_Rand.RandInt(0, 7))
	c := r.Map.GetCell(p.Coordinate())
	c.SetObject(m)
	return m, nil
}

func (m *Monster) Point() common.Point {
	return *m.CurrentLocation
}
