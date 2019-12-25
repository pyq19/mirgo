package common

import "testing"

func TestNewPointByCoordinate(t *testing.T) {
	p := NewPointByCoordinate("100,200")
	t.Log(p.Coordinate())
}
