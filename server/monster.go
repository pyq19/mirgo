package main

import "github.com/yenkeia/mirgo/common"

type Monster struct {
	ID               uint32
	Respawn          *Respawn
	Info             *common.MonsterInfo
	CurrentLocation  *common.Point
	CurrentDirection common.MirDirection
}

func NewMonster(r *Respawn) (m *Monster, err error) {
	m = new(Monster)
	m.Respawn = r
	m.Info = r.Map.Env.GameDB.GetMonsterInfoByID(r.Info.MonsterID)
	m.ID = r.Map.Env.NewObjectID()
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
