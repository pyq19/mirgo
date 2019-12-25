package main

import (
	"fmt"
	"sync"
)

// Grid 一个地图中的区域类
type Grid struct {
	GID     int       // 区域ID
	MinX    int       // 区域左边界坐标
	MaxX    int       // 区域右边界坐标
	MinY    int       // 区域上边界坐标
	MaxY    int       // 区域下边界坐标
	Players *sync.Map // 当前区域内的玩家 key=playerID  value=*player
	Cells   *sync.Map // key=Cell.Coordinate  value=*Cell
}

// NewGrid 初始化一个区域
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:     gID,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    minY,
		MaxY:    maxY,
		Players: new(sync.Map),
		Cells:   new(sync.Map),
	}
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
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d \n",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY)
}

func (g *Grid) AddCell(c *Cell) {
	g.Cells.Store(c.Coordinate, c)
}
