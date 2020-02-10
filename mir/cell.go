package mir

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/yenkeia/mirgo/common"
)

// Cell 地图格子
type Cell struct {
	Point     common.Point
	Map       *Map
	Attribute common.CellAttribute
	Objects   *sync.Map // map[IMapObject.ID]IMapObject
}

// IsValid ...
func (c *Cell) IsValid() bool {
	return c.Attribute == common.CellAttributeWalk
}

// IsEmpty ...
func (c *Cell) IsEmpty() bool {
	var cnt int32
	c.Objects.Range(func(k, v interface{}) bool {
		atomic.AddInt32(&cnt, 1)
		return true
	})
	return cnt == 0
}

// HasItem 判断是否有游戏物品/装备/金币
func (c *Cell) HasItem() bool {
	var cnt int32
	c.Objects.Range(func(k, v interface{}) bool {
		objectType := v.(IMapObject).GetRace()
		if objectType == common.ObjectTypeItem {
			atomic.AddInt32(&cnt, 1)
		}
		// 有 NPC 的 cell 也不能放置物品
		if objectType == common.ObjectTypeMerchant {
			atomic.AddInt32(&cnt, 1)
		}
		return true
	})
	return cnt > 0
}

// HasObject 判断是否有游戏对象 Player NPC Monster
func (c *Cell) HasObject() bool {
	var cnt int32
	c.Objects.Range(func(k, v interface{}) bool {
		objectType := v.(IMapObject).GetRace()
		if objectType == common.ObjectTypePlayer ||
			objectType == common.ObjectTypeMerchant ||
			objectType == common.ObjectTypeMonster {
			atomic.AddInt32(&cnt, 1)
		}
		return true
	})
	return cnt > 0
}

// CanWalk ...
func (c *Cell) CanWalk() bool {
	return c.Attribute == common.CellAttributeWalk
}

func (c *Cell) String() string {
	return fmt.Sprintf("cell pos: %s, Objects: %v \n", c.Point, c.Objects)
}

// AddObject ...
func (c *Cell) AddObject(obj IMapObject) {
	c.Objects.Store(obj.GetID(), obj)
}

// DeleteObject ...
func (c *Cell) DeleteObject(obj IMapObject) {
	c.Objects.Delete(obj.GetID())
}
