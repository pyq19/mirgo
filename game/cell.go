package game

import (
	"fmt"

	"github.com/yenkeia/mirgo/game/cm"
)

// Cell 地图格子
type Cell struct {
	Point     cm.Point
	Attribute cm.CellAttribute
	objects   map[uint32]IMapObject
}

func NewCell(attr cm.CellAttribute) *Cell {
	return &Cell{
		Attribute: attr,
		objects:   map[uint32]IMapObject{},
	}
}

// IsValid ...
func (c *Cell) IsValid() bool {
	return c.Attribute == cm.CellAttributeWalk
}

// IsEmpty ...
func (c *Cell) IsEmpty() bool {
	return len(c.objects) == 0
}

// HasItem 判断是否有游戏物品/装备/金币
func (c *Cell) HasItem() bool {

	for _, o := range c.objects {
		switch o.(type) {
		// 有 NPC 的 cell 也不能放置物品
		case *NPC, *Item:
			return true
		}
	}

	return false
}

// HasObject 判断是否有游戏对象 Player NPC Monster
func (c *Cell) HasObject() bool {

	for _, o := range c.objects {
		switch o.(type) {
		case *NPC, *Player, *Monster:
			return true
		}
	}

	return false
}

// CanWalk ...
func (c *Cell) CanWalk() bool {
	return c.Attribute == cm.CellAttributeWalk
}

func (c *Cell) String() string {
	return fmt.Sprintf("cell pos: %s, Objects: %v \n", c.Point, c.objects)
}

// AddObject ...
func (c *Cell) AddObject(obj IMapObject) {
	c.objects[obj.GetID()] = obj
}

// DeleteObject ...
func (c *Cell) DeleteObject(obj IMapObject) {
	delete(c.objects, obj.GetID())
}
