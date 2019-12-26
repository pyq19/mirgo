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
	Players  *sync.Map // 当前区域内的玩家  {Player.Character.Id: *Player}
	Monsters *sync.Map // 当前区域内的怪物
	NPCs     *sync.Map // 当前区域内的 NPC  {NPC.Info.Id: *NPC}
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
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d \n",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY)
}

// AddPlayer 向当前区域中添加一个玩家
func (g *Grid) AddPlayer(p *Player) {
	g.Players.Store(p.Character.Id, p)
}

// DeletePlayer 从区域中删除一个玩家
func (g *Grid) DeletePlayer(p *Player) {
	v, ok := g.Players.Load(p.Character.Id)
	if !ok {
		return
	}
	g.Players.Delete(v.(*Player).Character.Id)
}

func (g *Grid) AddNPC(n *NPC) {

}
