package main

import "testing"

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

func TestAOIManager_GetGridByCoordinate(t *testing.T) {
	var g *Grid
	m := NewAOIManager(nil, 0, 2000, 4, 0, 1200, 3)
	g = m.GetGridByCoordinate("550,880")
	t.Log(g.String())
	g = m.GetGridByCoordinate("0,0")
	t.Log(g.String())
	g = m.GetGridByCoordinate("1200,550")
	t.Log(g.String())
	g = m.GetGridByCoordinate("300,700")
	t.Log(g.String())
	g = m.GetGridByCoordinate("300,1100")
	t.Log(g.String())
}

func TestAOIManager_GetSurroundGridsByCoordinate(t *testing.T) {
	m := NewAOIManager(nil, 0, 2000, 4, 0, 1200, 3)
	var gs []*Grid
	gs = m.GetSurroundGridsByCoordinate("0,0")
	t.Log(gs)
	gs = m.GetSurroundGridsByCoordinate("1700,0")
	t.Log(gs)
	gs = m.GetSurroundGridsByCoordinate("1100,500")
	t.Log(gs)
}

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
