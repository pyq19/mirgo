package main

import (
	"fmt"
	"sync"
)

// Grid 一个地图中的格子类
type Grid struct {
	GID     int       // 格子ID
	MinX    int       // 格子左边界坐标
	MaxX    int       // 格子右边界坐标
	MinY    int       // 格子上边界坐标
	MaxY    int       // 格子下边界坐标
	players *sync.Map // 当前格子内的玩家 key=playerID  value=*player
}

// NewGrid 初始化一个格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:     gID,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    minY,
		MaxY:    maxY,
		players: new(sync.Map),
	}
}

// Add 向当前格子中添加一个玩家
func (g *Grid) Add(p *Player) {
	g.players.Store(p.ID, p)
}

// Remove 从格子中删除一个玩家
func (g *Grid) Remove(p *Player) {
	v, ok := g.players.Load(p.ID)
	if !ok {
		return
	}
	g.players.Delete(v.(*Player).ID)
}

// GetPlayerID 得到当前格子中所有的玩家
func (g *Grid) GetAllPlayer() (players []*Player) {
	g.players.Range(func(k, v interface{}) bool {
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
