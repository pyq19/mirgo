package main

import (
	"fmt"
	"sync"
)

// Grid 一个地图中的区域类
type Grid struct {
	AOI      *AOIManager
	GID      int       // 区域ID
	MinX     int       // 区域左边界坐标
	MaxX     int       // 区域右边界坐标
	MinY     int       // 区域上边界坐标
	MaxY     int       // 区域下边界坐标
	Players  *sync.Map // 当前区域内的玩家  {Player.ID: *Player}
	Monsters *sync.Map // 当前区域内的怪物
	NPCs     *sync.Map // 当前区域内的 NPC  {NPC.Info.ID: *NPC}
}

// NewGrid 初始化一个区域
func NewGrid(aoi *AOIManager, gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		AOI:      aoi,
		GID:      gID,
		MinX:     minX,
		MaxX:     maxX,
		MinY:     minY,
		MaxY:     maxY,
		Players:  new(sync.Map),
		Monsters: new(sync.Map),
		NPCs:     new(sync.Map),
	}
}

// GetPlayerID 得到当前区域中所有的玩家
func (g *Grid) GetAllPlayer() (players []*Player) {
	g.Players.Range(func(k, v interface{}) bool {
		players = append(players, v.(*Player))
		return true
	})
	return
}

// String 打印信息方法
func (g *Grid) String() string {
	res1 := fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d \n",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY)
	//res2 := fmt.Sprintf("Players: %v, NPCs: %v, Monsters: %v", g.Players, g.NPCs, g.Monsters)
	res2 := ""
	g.Monsters.Range(func(k, v interface{}) bool {
		m := v.(*Monster)
		res2 += fmt.Sprintf("Coordinate: %s, MonsterID: %d, ptr: %p\n", m.CurrentLocation.Coordinate(), m.ID, m)
		return true
	})
	res3 := ""
	g.NPCs.Range(func(k, v interface{}) bool {
		n := v.(*NPC)
		res3 += fmt.Sprintf("NPC Coordinate: %s, ID: %d, name: %s\n", n.Point().Coordinate(), n.Info.ID, n.Info.Name)
		return true
	})
	return res1 + "Monsters: \n" + res2 + "NPCs: \n" + res3
}

// AddPlayer 向当前区域中添加一个玩家
func (g *Grid) AddPlayer(p *Player) {
	g.Players.Store(p.ID, p)
}

// DeletePlayer 从区域中删除一个玩家
func (g *Grid) DeletePlayer(p *Player) {
	v, ok := g.Players.Load(p.ID)
	if !ok {
		return
	}
	g.Players.Delete(v.(*Player).ID)
}

func (g *Grid) AddNPC(n *NPC) {
	g.NPCs.Store(n.Info.ID, n)
}

func (g *Grid) DeleteNPC(n *NPC) {
	v, ok := g.NPCs.Load(n.Info.ID)
	if !ok {
		return
	}
	g.NPCs.Delete(v.(*NPC).Info.ID)
}

func (g *Grid) AddMonster(m *Monster) {
	g.Monsters.Store(m.ID, m)
}

func (g *Grid) DeleteMonster(m *Monster) {
	v, ok := g.Monsters.Load(m.ID)
	if !ok {
		return
	}
	g.Monsters.Delete(v.(*Monster).ID)
}
