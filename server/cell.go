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
	Objects    []interface{}
	lock       sync.RWMutex
}

func (c *Cell) Point() *common.Point {
	return common.NewPointByCoordinate(c.Coordinate)
}

func (c *Cell) IsEmpty() bool {
	if len(c.Objects) == 0 {
		return true
	}
	return false
}

func (c *Cell) CanWalk() bool {
	return c.Attribute == common.CellAttributeWalk
}

func (c *Cell) IsValid() bool {
	return c.CanWalk() && c.IsEmpty()
}

func (c *Cell) String() string {
	return fmt.Sprintf("Coordinate: %s, Objects: %v \n", c.Coordinate, c.Objects)
}

func (c *Cell) SetObject(obj interface{}) {
	c.lock.Lock()
	c.Objects = append(c.Objects, obj)
	c.lock.Unlock()
}

func (c *Cell) GetRace(obj interface{}) common.ObjectType {
	switch obj.(type) {
	case *Player:
		return common.ObjectTypePlayer
	case *NPC:
		return common.ObjectTypeMerchant
	case *Monster:
		return common.ObjectTypeMonster
	}
	return common.ObjectTypeNone
}
