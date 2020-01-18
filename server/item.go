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
	res := &server.ObjectItem{
		ObjectID:  i.GetID(),
		Name:      i.Name,
		NameColor: i.NameColor,
		Location:  i.CurrentLocation,
		Image:     nil,                  // TODO
		Grade:     common.ItemGradeNone, // TODO
	}
	return res
}
