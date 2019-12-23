package main

import "github.com/yenkeia/mirgo/common"

type Cell struct {
	Map        *Map
	Coordinate string       // 坐标 x,y
	Point      common.Point // 坐标点
	Attribute  common.CellAttribute
	Respawn    *Respawn
	NPC        *NPC
	Player     *Player
}

func (c *Cell) Empty() bool {
	if c.Respawn != nil {
		return false
	}
	if c.NPC != nil {
		return false
	}
	if c.Player != nil {
		return false
	}
	return true
}

func (c *Cell) SetPlayer(p *Player) {
	c.Player = p
	p.Cell = c
	p.Character.CurrentMapId = int32(c.Map.Id)
	p.Character.CurrentLocationX = int32(c.Point.X)
	p.Character.CurrentLocationY = int32(c.Point.Y)
}

func (c *Cell) SetNPC(n *NPC) {
	c.NPC = n
	n.Cell = c
	// NPC 不会动，设置位置也没用
	//n.Info.LocationX = int(c.Point.X)
	//n.Info.LocationY = int(c.Point.Y)
	//n.Info.MapId = c.Map.Id
}

func (c *Cell) SetRespawn(r *Respawn) {
	c.Respawn = r
	r.Cell = c
	r.Info.LocationX = int(c.Point.X)
	r.Info.LocationY = int(c.Point.Y)
}
