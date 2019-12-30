package main

import (
	"testing"
)

func TestGrid_GetPlayerID(t *testing.T) {
	grid := NewGrid(nil, 1, 1, 1, 1, 1)
	p1 := new(Player)
	p1.ID = 1
	p2 := new(Player)
	p2.ID = 22
	grid.AddPlayer(p1)
	grid.AddPlayer(p2)
	players := grid.GetAllPlayer()
	t.Log(players[0].ID, players[1].ID)
}
