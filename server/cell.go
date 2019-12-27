package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"sync"
)

type Cell struct {
	Map        *Map
	Coordinate string // 坐标 x,y
	Attribute  common.CellAttribute
	Object     interface{}
	lock       sync.RWMutex
}

func (c *Cell) isEmpty() bool {
	if c.Object == nil {
		return true
	}
	return false
}

func (c *Cell) canWalk() bool {
	return c.Attribute == common.CellAttributeWalk
}

func (c *Cell) IsValid() bool {
	return c.canWalk() && c.isEmpty()
}

func (c *Cell) String() string {
	return fmt.Sprintf("Coordinate: %s, Object: %v \n", c.Coordinate, c.Object)
}

func (c *Cell) SetObject(obj interface{}) {
	c.lock.Lock()
	c.Object = obj
	c.lock.Unlock()
}
