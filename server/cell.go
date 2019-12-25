package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
)

type Cell struct {
	Map        *Map
	Coordinate string // 坐标 x,y
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

func (c *Cell) String() string {
	return fmt.Sprintf("Coordinate: %s, Respawn: %v, Player: %v, NPC: %v\n", c.Coordinate, c.Respawn, c.Player, c.NPC)
}
