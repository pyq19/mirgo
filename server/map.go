package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// Map ...
type Map struct {
	Env    *Environ
	Width  uint16 // 测试用
	Height uint16 // 测试用
	Info   *common.MapInfo
	AOI    *AOIManager
	Cells  *sync.Map // key=Cell.Coordinate  value=*Cell
}

func (m *Map) String() string {
	return fmt.Sprintf("Map(%d, Filename: %s, Title: %s, Width: %d, Height: %d)", m.Info.ID, m.Info.Filename, m.Info.Title, m.Width, m.Height)
}

func (m *Map) Submit(t *Task) {
	m.Env.Game.Pool.EntryChan <- t
}

func (m *Map) GetAllPlayers() []*Player {
	players := make([]*Player, 0)
	m.AOI.grids.Range(func(k, v interface{}) bool {
		g := v.(*Grid)
		players = append(players, g.GetAllPlayer()...)
		return true
	})
	return players
}

// Broadcast send message to all players in this map
func (m *Map) Broadcast(msg interface{}) {
	players := m.GetAllPlayers()
	for i := range players {
		players[i].Enqueue(msg)
	}
}

func (m *Map) GetCell(coordinate string) *Cell {
	v, ok := m.Cells.Load(coordinate)
	if !ok {
		return nil
	}
	return v.(*Cell)
}

func (m *Map) AddObject(obj IMapObject) {
	coordinate := obj.GetCoordinate()
	grid := m.AOI.GetGridByCoordinate(coordinate)
	grid.AddObject(obj)
	c := m.GetCell(coordinate)
	if c == nil {
		// FIXME
		//log.Warnf("coordinate: %s is not walkable\n", coordinate)
		return
	}
	c.AddObject(obj)
}

func (m *Map) DeleteObject(obj IMapObject) {
	coordinate := obj.GetCoordinate()
	grid := m.AOI.GetGridByCoordinate(coordinate)
	grid.DeleteObject(obj)
	c := m.GetCell(coordinate)
	if c == nil {
		return
	}
	c.DeleteObject(obj)
}

// UpdateObject 更新对象在 Cells, AOI 中的数据, 如果更新成功返回 true
func (m *Map) UpdateObject(obj IMapObject, points ...common.Point) bool {
	for i := range points {
		c := m.GetCell(points[i].Coordinate())
		if c == nil || !c.CanWalkAndIsEmpty() {
			return false
		}
	}
	c1 := obj.GetCell()
	c1.DeleteObject(obj)
	c2 := m.GetCell(points[len(points)-1].Coordinate())
	c2.AddObject(obj)
	m.changeAOI(obj, c1, c2)
	return true
}

func (m *Map) changeAOI(obj IMapObject, c1 *Cell, c2 *Cell) {
	g1 := m.AOI.GetGridByPoint(c1.Point())
	g2 := m.AOI.GetGridByPoint(c2.Point())
	if g1.GID == g2.GID {
		return
	}
	g1.DeleteObject(obj)
	g2.AddObject(obj)
	switch obj.GetRace() {
	case common.ObjectTypePlayer:
		p := obj.(*Player)
		p.Broadcast(ServerMessage{}.ObjectPlayer(p))
	case common.ObjectTypeMonster:
		m := obj.(*Monster)
		m.Broadcast(ServerMessage{}.ObjectMonster(m))
	}
}

// InitNPCs 初始化地图上的 NPC
func (m *Map) InitNPCs() error {
	for _, ni := range m.Env.GameDB.NpcInfos {
		ni := ni
		if ni.MapID == m.Info.ID {
			m.AddObject(NewNPC(m, &ni))
		}
	}
	return nil
}

// InitMonsters 初始化地图上的怪物
func (m *Map) InitMonsters() error {
	for _, ri := range m.Env.GameDB.RespawnInfos {
		if ri.MapID == m.Info.ID {
			r, err := NewRespawn(m, ri)
			if err != nil {
				return err
			}
			for _, a := range r.AliveMonster {
				a := a
				m.AddObject(a)
			}
		}
	}
	return nil
}

// GetValidPoint
func (m *Map) GetValidPoint(x int, y int, spread int) (common.Point, error) {
	if spread == 0 {
		c := m.GetCell(common.Point{X: uint32(x), Y: uint32(y)}.Coordinate())
		if c != nil && c.CanWalkAndIsEmpty() {
			return common.NewPointByCoordinate(c.Coordinate), nil
		}
		return common.Point{}, fmt.Errorf("GetValidPoint: (x: %d, y: %d), spread: %d\n", x, y, spread)
	}

	for i := 0; i < 500; i++ {
		p := common.Point{
			X: uint32(x + G_Rand.RandInt(-spread, spread+1)),
			Y: uint32(y + G_Rand.RandInt(-spread, spread+1)),
		}
		c := m.GetCell(p.Coordinate())
		if c == nil || !c.CanWalk() {
			continue
		}
		return p, nil
	}
	return common.Point{}, fmt.Errorf("map(%v) no valid point in (%d,%d) spread: %d", m, x, y, spread)
}

func (m *Map) GetNextCell(cell *Cell, direction common.MirDirection, step uint32) *Cell {
	p := cell.Point().NextPoint(direction, step)
	return m.GetCell(p.Coordinate())
}

// GetAreaMapObjects 传入一个点，获取该点附近 9 个 AOI 区域内 MapObject
func (m *Map) GetAreaObjects(p common.Point) (objs []IMapObject) {
	grids := m.AOI.GetSurroundGridsByCoordinate(p.Coordinate())
	for i := range grids {
		g := grids[i]
		objs = append(objs, g.GetAllObjects()...)
	}
	return
}
