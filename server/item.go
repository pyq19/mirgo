package main

import (
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type ItemObject struct {
	MapObject
	Gold     uint64
	UserItem *common.UserItem
}

func (i *ItemObject) GetID() uint32 {
	return i.ID
}

func (i *ItemObject) GetRace() common.ObjectType {
	return common.ObjectTypeItem
}

func (i *ItemObject) GetCoordinate() string {
	return i.CurrentLocation.Coordinate()
}

func (i *ItemObject) GetPoint() common.Point {
	return i.CurrentLocation
}

func (i *ItemObject) GetCell() *Cell {
	return i.Map.GetCell(i.GetCoordinate())
}

func (i *ItemObject) Broadcast(interface{}) {
	return
}

func (i *ItemObject) GetDirection() common.MirDirection {
	return i.CurrentDirection
}

func (i *ItemObject) GetInfo() interface{} {
	if i.UserItem == nil {
		res := &server.ObjectGold{
			ObjectID: i.GetID(),
			Gold:     uint32(i.Gold),
			Location: i.GetPoint(),
		}
		return res
	} else {
		res := &server.ObjectItem{
			ObjectID:  i.GetID(),
			Name:      i.Name,
			NameColor: i.NameColor,
			Location:  i.GetPoint(),
			Image:     0,                    // TODO
			Grade:     common.ItemGradeNone, // TODO
		}
		return res
	}
}

// TODO
// Drop 物品加入到地图上，传入中心点 p，范围 r
func (i *ItemObject) Drop(p common.Point, r int) (string, bool) {
	i.Map.AddObject(i)
	i.Broadcast(i.GetInfo())
	return "", true
}
