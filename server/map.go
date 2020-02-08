package main

import (
	"fmt"

	"github.com/yenkeia/mirgo/common"
)

// Map ...
type Map struct {
	Env    *Environ
	Width  int
	Height int
	Info   *common.MapInfo
	AOI    *AOIManager
	cells  []*Cell
}

func NewMap(w, h int) *Map {
	m := &Map{
		Width:  w,
		Height: h,
		cells:  make([]*Cell, w*h),
	}
	m.AOI = newAOI(m, w, h)
	return m
}

func (m *Map) GetCell(p common.Point) *Cell {
	return m.GetCellXY(int(p.X), int(p.Y))
}
func (m *Map) GetCellXY(x, y int) *Cell {
	return m.cells[x+y*m.Width]
}
func (m *Map) SetCell(p common.Point, c *Cell) {
	m.SetCellXY(int(p.X), int(p.Y), c)
}
func (m *Map) SetCellXY(x, y int, c *Cell) {
	m.cells[x+y*m.Width] = c
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

func (m *Map) AddObject(obj IMapObject) (string, bool) {
	if obj == nil || obj.GetID() == 0 {
		return "", false
	}
	grid := m.AOI.GetGridByPoint(obj.GetPoint())
	grid.AddObject(obj)
	c := m.GetCell(obj.GetPoint())
	if c == nil {
		// FIXME
		return fmt.Sprintf("pos: %s is not walkable\n", obj.GetPoint()), false
	}
	c.AddObject(obj)
	return "", true
}

func (m *Map) DeleteObject(obj IMapObject) {
	if obj == nil || obj.GetID() == 0 {
		return
	}
	grid := m.AOI.GetGridByPoint(obj.GetPoint())
	grid.DeleteObject(obj)
	c := m.GetCell(obj.GetPoint())
	if c == nil {
		return
	}
	c.DeleteObject(obj)
}

// UpdateObject 更新对象在 Cells, AOI 中的数据, 如果更新成功返回 true
func (m *Map) UpdateObject(obj IMapObject, points ...common.Point) bool {
	for i := range points {
		c := m.GetCell(points[i])
		if c == nil || !c.CanWalk() {
			return false
		}

		blocking := false
		c.Objects.Range(func(k, v interface{}) bool {
			if v.(IMapObject).IsBlocking() {
				blocking = true
				return false
			}
			return true
		})

		if blocking {
			return false
		}
	}
	c1 := obj.GetCell()
	c1.DeleteObject(obj)
	c2 := m.GetCell(points[len(points)-1])
	c2.AddObject(obj)
	m.changeAOI(obj, c1, c2)
	return true
}

func (m *Map) changeAOI(obj IMapObject, c1 *Cell, c2 *Cell) {
	g1 := m.AOI.GetGridByPoint(c1.Point)
	g2 := m.AOI.GetGridByPoint(c2.Point)
	if g1.GID == g2.GID {
		return
	}
	g1.DeleteObject(obj)
	g2.AddObject(obj)
	switch obj.GetRace() {
	case common.ObjectTypePlayer:
		p := obj.(*Player)
		p.Broadcast(ServerMessage{}.ObjectPlayer(p))
		p.EnqueueAreaObjects(g1, g2)
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
			n := NewNPC(m, &ni)
			m.Env.AddNPC(n)
			m.AddObject(n)
		}
	}
	return nil
}

// InitMonsters 初始化地图上的怪物
func (m *Map) InitMonsters() error {
	for _, ri := range m.Env.GameDB.RespawnInfos {
		ri := ri
		if ri.MapID == m.Info.ID {
			cnt := ri.Count
			for i := 0; i < cnt; i++ {
				p, err := m.GetValidPoint(ri.LocationX, ri.LocationY, ri.Spread)
				if err != nil {
					continue
				}
				m.AddObject(NewMonster(m, p, m.Env.GameDB.GetMonsterInfoByID(ri.MonsterID)))
			}
		}
	}
	return nil
}

// GetValidPoint ...
func (m *Map) GetValidPoint(x int, y int, spread int) (common.Point, error) {
	if spread == 0 {
		c := m.GetCellXY(x, y)
		if c != nil && c.CanWalk() && !c.HasObject() {
			return c.Point, nil
		}
		return common.Point{}, fmt.Errorf("GetValidPoint: (x: %d, y: %d), spread: %d\n", x, y, spread)
	}

	for i := 0; i < 500; i++ {
		p := common.Point{
			X: uint32(AbsInt(x + RandomInt(-spread, spread))),
			Y: uint32(AbsInt(y + RandomInt(-spread, spread))),
		}
		c := m.GetCell(p)
		if c == nil || !c.CanWalk() {
			continue
		}
		return p, nil
	}
	return common.Point{}, fmt.Errorf("map(%v) no valid point in (%d,%d) spread: %d", m, x, y, spread)
}

func (m *Map) GetNextCell(cell *Cell, direction common.MirDirection, step uint32) *Cell {
	p := cell.Point.NextPoint(direction, step)
	return m.GetCell(p)
}

// GetAreaMapObjects 传入一个点，获取该点附近 9 个 AOI 区域内 MapObject
func (m *Map) GetAreaObjects(p common.Point) (objs []IMapObject) {
	grids := m.AOI.GetSurroundGrids(p)
	for i := range grids {
		g := grids[i]
		objs = append(objs, g.GetAllObjects()...)
	}
	return
}

// GetObjectInAreaByID 查找点 p 附近的区域中 ObjectID 为 id 的对象
func (m *Map) GetObjectInAreaByID(id uint32, p common.Point) IMapObject {
	areaObjects := m.GetAreaObjects(p)
	for i := range areaObjects {
		obj := areaObjects[i]
		if obj.GetID() == id {
			return obj
		}
	}
	return nil
}

// 从p点开始（包含P），由内至外向周围遍历cell。回调函数返回false，停止遍历
func (m *Map) RangeCell(p common.Point, depth int, fun func(c *Cell) bool) {

	px, py := int(p.X), int(p.Y)

	for d := 0; d <= depth; d++ {
		for y := py - d; y <= py+d; y++ {
			if y < 0 {
				continue
			}
			if y >= m.Height {
				break
			}

			for x := px - d; x <= px+d; {

				if x >= m.Width {
					break
				}

				if x >= 0 {
					if !fun(m.GetCellXY(x, y)) {
						return
					}
				}

				if y-py == d || y-py == -d {
					x++ // x += 1
				} else {
					x += d * 2
				}
			}
		}
	}
}

func (m *Map) RangeObject(p common.Point, depth int, fun func(IMapObject) bool) {
	var ret = true
	m.RangeCell(p, depth, func(c *Cell) bool {
		if c != nil && c.Objects != nil {
			c.Objects.Range(func(k, v interface{}) bool {
				ret = fun(v.(IMapObject))
				return ret
			})
		}

		return ret
	})
}

// CompleteMagic ...
func (m *Map) CompleteMagic(args ...interface{}) {
	magic := args[0].(*common.UserMagic)
	switch magic.Spell {
	case common.SpellSummonSkeleton, common.SpellSummonShinsu, common.SpellSummonHolyDeva, common.SpellSummonVampire, common.SpellSummonToad, common.SpellSummonSnakes:
		player := args[1].(*Player)
		monster := args[2].(*Monster)
		front := args[3].(common.Point)
		if monster.Master.IsDead() {
			return
		}
		cell := m.GetCell(front)
		if cell.IsValid() {
			monster.Spawn(m, front)
		} else {
			monster.Spawn(m, player.GetPoint())
		}
		pets := monster.Master.Pets
		pets = append(pets, monster)
	}
}
