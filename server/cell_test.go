package main

import "testing"

func TestCell_SetObject(t *testing.T) {
	c := new(Cell)
	c.SetObject(nil)
	t.Log(c.Object == nil) // true
	c.SetObject(1)
	t.Log(c.Object == nil) // false
	c.SetObject(nil)
	t.Log(c.Object == nil) // true
}
