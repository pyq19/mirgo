package main

import "testing"

func TestCell_SetObject(t *testing.T) {
	c := new(Cell)
	c.SetObject(nil)
	t.Log(c.Objects == nil) // true
	c.SetObject(1)
	t.Log(c.Objects == nil) // false
	c.SetObject(nil)
	t.Log(c.Objects == nil) // true
}
