package main

import (
	"github.com/yenkeia/mirgo/common"
	"testing"
)

func TestGrid_GetPlayerID(t *testing.T) {
	grid := NewGrid(nil, 1, 1, 1, 1, 1)
	p1 := new(Player)
	p1.Character = new(common.Character)
	p1.Character.ID = 1
	p2 := new(Player)
	p2.Character = new(common.Character)
	p2.Character.ID = 22
	grid.AddPlayer(p1)
	grid.AddPlayer(p2)
	players := grid.GetAllPlayer()
	t.Log(players[0].Character.ID, players[1].Character.ID)
}
