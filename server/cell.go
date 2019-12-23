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

func (c *Cell) AddPlayer(p *Player) {

}

func (c *Cell) AddNPC(n *NPC) {

}

func (c *Cell) AddRespawn(r *Respawn) {

}
