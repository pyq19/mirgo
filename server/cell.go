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
}

func (c *Cell) SetNPC(n *NPC) {
	c.NPC = n
	n.Cell = c
}

func (c *Cell) SetRespawn(r *Respawn) {
	c.Respawn = r
	r.Cell = c
}
