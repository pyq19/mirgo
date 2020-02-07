package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// Grid 一个地图中的区域类
type Grid struct {
	AOI     *AOIManager
	GID     int       // 区域ID
	MinX    int       // 区域左边界坐标
	MaxX    int       // 区域右边界坐标
	MinY    int       // 区域上边界坐标
	MaxY    int       // 区域下边界坐标
	Objects *sync.Map // 当前区域内的玩家
}

// NewGrid 初始化一个区域
func NewGrid(aoi *AOIManager, gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		AOI:     aoi,
		GID:     gID,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    minY,
		MaxY:    maxY,
		Objects: new(sync.Map),
	}
}

// String 打印信息方法
func (g *Grid) String() string {
	var (
		players  []*Character
		monsters []*Monster
		npcs     []*NPC
	)
	gridInfo := fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d \n",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY)
	g.Objects.Range(func(k, v interface{}) bool {
		o := v.(IMapObject)
		switch o.GetRace() {
		case common.ObjectTypePlayer:
			players = append(players, o.(*Character))
		case common.ObjectTypeMonster:
			monsters = append(monsters, o.(*Monster))
		case common.ObjectTypeMerchant:
			npcs = append(npcs, o.(*NPC))
		}
		return true
	})
	return fmt.Sprintf("%s\nPlayers: %v\nMonsters: %v\nNPCs: %v\n", gridInfo, players, monsters, npcs)
}

func (g *Grid) AddObject(obj IMapObject) {
	g.Objects.Store(obj.GetID(), obj)
}

func (g *Grid) DeleteObject(obj IMapObject) {
	v, ok := g.Objects.Load(obj.GetID())
	if !ok {
		return
	}
	g.Objects.Delete(v.(IMapObject).GetID())
}

// GetAllPlayer 得到当前区域中所有的玩家
func (g *Grid) GetAllPlayer() (players []*Character) {
	g.Objects.Range(func(k, v interface{}) bool {
		o := v.(IMapObject)
		if o.GetRace() == common.ObjectTypePlayer {
			players = append(players, o.(*Character))
		}
		return true
	})
	return
}

func (g *Grid) GetAllObjects() (objs []IMapObject) {
	g.Objects.Range(func(k, v interface{}) bool {
		o := v.(IMapObject)
		objs = append(objs, o)
		return true
	})
	return
}
