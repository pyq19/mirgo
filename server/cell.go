package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"sync"
	"sync/atomic"
)

type Cell struct {
	Map        *Map
	Coordinate string // 坐标 x,y
	Attribute  common.CellAttribute
	Objects    *sync.Map // map[IMapObject.ID]IMapObject
}

func (c *Cell) Point() common.Point {
	return common.NewPointByCoordinate(c.Coordinate)
}

func (c *Cell) IsEmpty() bool {
	var cnt int32
	c.Objects.Range(func(k, v interface{}) bool {
		atomic.AddInt32(&cnt, 1)
		return true
	})
	return cnt == 0
}

func (c *Cell) HasItemObject() bool {
	var cnt int32
	c.Objects.Range(func(k, v interface{}) bool {
		if v.(IMapObject).GetRace() == common.ObjectTypeItem {
			atomic.AddInt32(&cnt, 1)
		}
		return true
	})
	return cnt == 0
}

func (c *Cell) CanWalk() bool {
	return c.Attribute == common.CellAttributeWalk
}

func (c *Cell) CanWalkAndIsEmpty() bool {
	return c.CanWalk() && c.IsEmpty()
}

func (c *Cell) String() string {
	return fmt.Sprintf("Coordinate: %s, Objects: %v \n", c.Coordinate, c.Objects)
}

func (c *Cell) AddObject(obj IMapObject) {
	c.Objects.Store(obj.GetID(), obj)
}

func (c *Cell) DeleteObject(obj IMapObject) {
	c.Objects.Delete(obj.GetID())
}

/*
func (c *Cell) GetRace(obj IMapObject) common.ObjectType {
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
*/
