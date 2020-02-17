package mir

import (
	"fmt"
	"time"

	"github.com/yenkeia/mirgo/common"
)

// Map ...
type Map struct {
	Env    *Environ
	Width  int
	Height int
	Info   *common.MapInfo
	cells  []*Cell
	// AOI    *AOIManager

	players  map[uint32]*Player
	monsters map[uint32]*Monster
	npcs     map[uint32]*NPC
}

func NewMap(w, h int) *Map {
	m := &Map{
		Width:    w,
		Height:   h,
		cells:    make([]*Cell, w*h),
		players:  map[uint32]*Player{},
		monsters: map[uint32]*Monster{},
		npcs:     map[uint32]*NPC{},
	}
	// m.AOI = newAOI(m, w, h)
	return m
}

func (m *Map) GetCell(p common.Point) *Cell {
	return m.GetCellXY(int(p.X), int(p.Y))
}
func (m *Map) GetCellXY(x, y int) *Cell {
	if m.InMap(x, y) {
		return m.cells[x+y*m.Width]
	} else {
		return nil
	}
}

func (m *Map) InMap(x, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
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

func (m *Map) GetAllPlayers() map[uint32]*Player {
	return m.players
}
func (m *Map) GetNPC(id uint32) *NPC {
	return m.npcs[id]
}

// Broadcast send message to all players in this map
func (m *Map) Broadcast(msg interface{}) {
	for _, p := range m.players {
		p.Enqueue(msg)
	}
}

const DataRange = 6

// 位置，消息，跳过玩家
func (m *Map) BroadcastP(pos common.Point, msg interface{}, me *Player) {
	m.Submit(NewTask(func(args ...interface{}) {
		for _, plr := range m.players {
			if InRange(pos, plr.CurrentLocation, DataRange) {
				if plr != me {
					plr.Enqueue(msg)
				}
			}
		}
		// grids := m.AOI.GetSurroundGrids(p)
		// for i := range grids {
		// 	areaPlayers := grids[i].GetAllPlayer()
		// 	for _, p := range areaPlayers {
		// 		if p != me {
		// 			areaPlayers[i].Enqueue(msg)
		// 		}
		// 	}
		// }
	}))
}

func (m *Map) AddObject(obj IMapObject) (string, bool) {
	if obj == nil || obj.GetID() == 0 {
		return "", false
	}
	// grid := m.AOI.GetGridByPoint(obj.GetPoint())
	// grid.AddObject(obj)
	c := m.GetCell(obj.GetPoint())
	if c == nil {
		// FIXME
		return fmt.Sprintf("pos: %s is not walkable\n", obj.GetPoint()), false
	}
	c.AddObject(obj)

	switch obj.(type) {
	case *Player:
		m.players[obj.GetID()] = obj.(*Player)
	case *NPC:
		m.npcs[obj.GetID()] = obj.(*NPC)
	case *Monster:
		m.monsters[obj.GetID()] = obj.(*Monster)
	}

	return "", true
}

func (m *Map) DeleteObject(obj IMapObject) {
	if obj == nil || obj.GetID() == 0 {
		return
	}
	// grid := m.AOI.GetGridByPoint(obj.GetPoint())
	// grid.DeleteObject(obj)
	c := m.GetCell(obj.GetPoint())
	if c == nil {
		return
	}
	c.DeleteObject(obj)

	switch obj.(type) {
	case *Player:
		delete(m.players, obj.GetID())
	case *Monster:
		delete(m.monsters, obj.GetID())
	case *NPC:
		delete(m.npcs, obj.GetID())
	}
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
	// g1 := m.AOI.GetGridByPoint(c1.Point)
	// g2 := m.AOI.GetGridByPoint(c2.Point)
	// if g1.GID == g2.GID {
	// 	return
	// }
	// g1.DeleteObject(obj)
	// g2.AddObject(obj)
	switch obj.GetRace() {
	case common.ObjectTypePlayer:
		p := obj.(*Player)
		p.Broadcast(ServerMessage{}.ObjectPlayer(p))
		p.EnqueueAreaObjects(c1, c2)
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
			// m.Env.AddNPC(n)
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
// func (m *Map) GetAreaObjects(p common.Point) (objs []IMapObject) {
// 	grids := m.AOI.GetSurroundGrids(p)
// 	for i := range grids {
// 		g := grids[i]
// 		objs = append(objs, g.GetAllObjects()...)
// 	}
// 	return
// }

// GetObjectInAreaByID 查找点 p 附近的区域中 ObjectID 为 id 的对象
// func (m *Map) GetObjectInAreaByID(id uint32, p common.Point) IMapObject {
// 	areaObjects := m.GetAreaObjects(p)
// 	for i := range areaObjects {
// 		obj := areaObjects[i]
// 		if obj.GetID() == id {
// 			return obj
// 		}
// 	}
// 	return nil
// }

func (m *Map) GetObjectInAreaByID(id uint32, p common.Point) IMapObject {
	var ret IMapObject
	m.RangeObject(p, 1, func(o IMapObject) bool {
		if o.GetID() == id {
			ret = o
			return false
		}
		return true
	})

	return ret
}

// 从p点开始（包含P），由内至外向周围遍历cell。回调函数返回false，停止遍历
func (m *Map) RangeCell(p common.Point, depth int, fun func(c *Cell, x, y int) bool) {

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
					if !fun(m.GetCellXY(x, y), x, y) {
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
	m.RangeCell(p, depth, func(c *Cell, _, _ int) bool {
		if c != nil && c.Objects != nil {
			c.Objects.Range(func(k, v interface{}) bool {
				ret = fun(v.(IMapObject))
				return ret
			})
		}

		return ret
	})
}

// Cell集合
type CellSet struct {
	M map[*Cell]bool
}

func NewCellSet() *CellSet {
	return &CellSet{M: map[*Cell]bool{}}
}

func (c *CellSet) Add(m *Map, x, y int, b bool) {
	cell := m.GetCellXY(x, y)
	if cell != nil {
		c.M[cell] = b
	}
}

// 根据两个点，求出 远离datarange内的cell，和新加的cell
func (m *Map) CalcDiff(from, to common.Point, datarange int) *CellSet {
	fx, fy, tx, ty := int(from.X), int(from.Y), int(to.X), int(to.Y)

	xChange, yChange := tx-fx, ty-fy
	set := NewCellSet()

	if xChange > 0 { // 右移
		for x := 0; x < xChange; x++ {
			for y := fy - datarange; y <= fy+datarange; y++ {
				set.Add(m, fx-datarange+x, y, false) // 左
			}
			for y := ty - datarange; y <= ty+datarange; y++ {
				set.Add(m, tx+datarange-x, y, true) // 右
			}
		}
	} else { // 左移
		for x := 0; x > xChange; x-- {
			for y := ty - datarange; y <= ty+datarange; y++ {
				set.Add(m, tx-datarange-x, y, true) // 左
			}
			for y := fy - datarange; y <= fy+datarange; y++ {
				set.Add(m, fx+datarange+x, y, false) // 右
			}
		}
	}
	if yChange < 0 { // 上移
		for y := 0; y > yChange; y-- {
			for x := tx - datarange; x <= tx+datarange; x++ {
				set.Add(m, x, ty-datarange-y, true) // 上
			}
			for x := fx - datarange; x <= fx+datarange; x++ {
				set.Add(m, x, fy+datarange+y, false) // 下
			}
		}
	} else { // 下移
		for y := 0; y < yChange; y++ {
			for x := fx - datarange; x <= fx+datarange; x++ {
				set.Add(m, x, fy-datarange+y, false) // 上
			}
			for x := tx - datarange; x <= tx+datarange; x++ {
				set.Add(m, x, ty+datarange-y, true) // 下
			}
		}
	}

	return set
}

// for test CalcDiff.
// func (m *Map) CalcDiff1(from, to common.Point, datarange int) *CellSet {
// 	fx, fy, tx, ty := int(from.X), int(from.Y), int(to.X), int(to.Y)

// 	oldcells := map[int]bool{}
// 	m.RangeCell(common.NewPoint(fx, fy), datarange, func(c *Cell, x, y int) bool {
// 		oldcells[x*10000+y] = true
// 		return true
// 	})

// 	newcells := map[int]bool{}
// 	m.RangeCell(common.NewPoint(tx, ty), datarange, func(c *Cell, x, y int) bool {
// 		newcells[x*10000+y] = true
// 		return true
// 	})

// 	added := map[int]bool{}
// 	for c := range newcells {
// 		if _, ok := oldcells[c]; ok {
// 			delete(oldcells, c)
// 		} else {
// 			added[c] = true
// 		}
// 	}

// 	cs := NewCellSet()
// 	for c := range oldcells {
// 		cs.Add(m, c/10000, c%10000, false)
// 	}

// 	for c := range added {
// 		cs.Add(m, c/10000, c%10000, true)
// 	}
// 	return cs
// }

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
	case common.SpellMassHealing:
		value := args[1].(int)
		location := args[2].(common.Point)
		player := args[3].(*Player)
		m.RangeObject(location, 1, func(o IMapObject) bool {
			if o.GetRace() == common.ObjectTypePlayer && o.IsFriendlyTarget(player) {
				target := o.(*Player)
				for i := range target.Buffs {
					if target.Buffs[i].BuffType == common.BuffTypeHiding {
						return true
					}
				}
				target.AddBuff(NewBuff(player.NewObjectID(), common.BuffTypeHiding, 0, time.Now().Add(time.Duration(value*1000)*time.Millisecond)))
			}
			return true
		})
	case common.SpellSoulShield, common.SpellBlessedArmour:
		value := args[1].(int)
		location := args[2].(common.Point)
		player := args[3].(*Player)
		buffType := common.BuffTypeSoulShield
		if magic.Spell == common.SpellBlessedArmour {
			buffType = common.BuffTypeBlessedArmour
		}
		m.RangeObject(location, 1, func(o IMapObject) bool {
			if o.GetRace() == common.ObjectTypePlayer {
				target := o.(*Player)
				target.AddBuff(NewBuff(player.NewObjectID(), buffType, int(target.Level)/7+4, time.Now().Add(time.Duration(value*1000)*time.Millisecond)))
			}
			return true
		})
	case common.SpellFireWall:
		// player := args[1].(*Player)
		// value := args[2].(int)
		// location := args[3].(common.Point)
		// player.LevelMagic(magic)
		// TODO SpellObject
	case common.SpellLightning:
		// player := args[1].(*Player)
		// value := args[2].(int)
		// location := args[3].(common.Point)
		// direction := args[4].(common.MirDirection)
	case common.SpellThunderStorm, common.SpellFlameField, common.SpellNapalmShot, common.SpellStormEscape:
		// player := args[1].(*Player)
		// value := args[2].(int)
		// location := args[3].(common.Point)
	}
}
