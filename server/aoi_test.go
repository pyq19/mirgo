package main

import (
	"testing"
)

/*
m := NewAOIManager(0, 2000, 4, 0, 1200, 3)

	 w = 2000 / 4   h = 1200 / 3

	-----------------	x
	| 0 | 1 | 2 | 3 |
	-----------------
	| 4 | 5 | 6 | 7 |
	-----------------
	| 8 | 9 | 10| 11|
	-----------------
	y
*/

func TestAOIManager_GetSurroundGridsByGridID(t *testing.T) {
	m := NewAOIManager(nil, 0, 2000, 4, 0, 1200, 3)
	var gs []*Grid
	gs = m.GetSurroundGridsByGridID(6)
	t.Log(gs)
}

func TestAOIManager_GetAllGrid(t *testing.T) {
	m := NewAOIManager(nil, 0, 2000, 4, 0, 1200, 3)
	m.grids.Range(func(key, value interface{}) bool {
		t.Log(key.(int))
		t.Log(value.(*Grid).String())
		t.Log()
		return true
	})
}

func TestGrid_GetPlayerID(t *testing.T) {
	grid := NewGrid(nil, 1, 1, 1, 1, 1)
	p1 := new(Character)
	p1.ID = 1
	p2 := new(Character)
	p2.ID = 22
	grid.AddObject(p1)
	grid.AddObject(p2)
	players := grid.GetAllPlayer()
	t.Log(players[0].ID, players[1].ID)
}
