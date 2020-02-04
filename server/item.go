package server

import (
	"fmt"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type Item struct {
	MapObject
	Gold     uint64
	UserItem *common.UserItem
}

func (i *Item) GetID() uint32 {
	return i.ID
}

func (i *Item) GetRace() common.ObjectType {
	return common.ObjectTypeItem
}

func (i *Item) GetCoordinate() string {
	return i.CurrentLocation.Coordinate()
}

func (i *Item) GetPoint() common.Point {
	return i.CurrentLocation
}

func (i *Item) GetCell() *Cell {
	return i.Map.GetCell(i.GetCoordinate())
}

func (i *Item) Broadcast(msg interface{}) {
	i.Map.Submit(NewTask(func(args ...interface{}) {
		grids := i.Map.AOI.GetSurroundGridsByCoordinate(i.GetCoordinate())
		for i := range grids {
			areaPlayers := grids[i].GetAllPlayer()
			for i := range areaPlayers {
				areaPlayers[i].Enqueue(msg)
			}
		}
	}))
}

func (i *Item) GetDirection() common.MirDirection {
	return i.CurrentDirection
}

func (i *Item) GetInfo() interface{} {
	if i.UserItem == nil {
		res := &server.ObjectGold{
			ObjectID:  i.GetID(),
			Gold:      uint32(i.Gold),
			LocationX: int32(i.GetPoint().X),
			LocationY: int32(i.GetPoint().Y),
		}
		return res
	} else {
		res := &server.ObjectItem{
			ObjectID:  i.GetID(),
			Name:      i.Name,
			NameColor: i.NameColor.ToInt32(),
			LocationX: int32(i.GetPoint().X),
			LocationY: int32(i.GetPoint().Y),
			Image:     i.GetImage(),
			Grade:     common.ItemGradeNone, // TODO
		}
		return res
	}
}

func (i *Item) IsAttackTarget(attacker IMapObject) bool {
	return false
}

func (i *Item) IsFriendlyTarget(attacker IMapObject) bool {
	return true
}

func (i *Item) GetBaseStats() BaseStats {
	return BaseStats{}
}

func (i *Item) GetItemInfo() common.ItemInfo {
	return *i.Map.Env.GameDB.GetItemInfoByID(int(i.UserItem.ItemID))
}

func (i *Item) GetImage() uint16 {
	info := i.GetItemInfo()
	switch info.Type {
	case common.ItemTypeAmulet:
		if info.StackSize > 0 {
			switch info.Shape {
			case 0: //Amulet
				if i.UserItem.Count >= 300 {
					return 3662
				}
				if i.UserItem.Count >= 200 {
					return 3661
				}
				if i.UserItem.Count >= 100 {
					return 3660
				}
				return 3660
			case 1: //Grey Poison
				if i.UserItem.Count >= 150 {
					return 3675
				}
				if i.UserItem.Count >= 100 {
					return 2960
				}
				if i.UserItem.Count >= 50 {
					return 3674
				}
				return 3673
			case 2: //Yellow Poison
				if i.UserItem.Count >= 150 {
					return 3672
				}
				if i.UserItem.Count >= 100 {
					return 2961
				}
				if i.UserItem.Count >= 50 {
					return 3671
				}
				return 3670
			}
		}
	}
	return info.Image
}

// Drop 物品加入到地图上，传入中心点 center，范围 distance
func (i *Item) Drop(center common.Point, distance int) (string, bool) {
	// 以 center 为中心，向外获取点，放入集合
	x := int(center.X)
	y := int(center.Y)
	points := make([]common.Point, 0)
	if distance == 0 {
		points = append(points, center)
	} else {
		for k := 1; k <= distance; k++ {
			minX := x - k
			maxX := x + k
			minY := y - k
			maxY := y + k
			for n := minY; maxY >= n; n++ {
				for m := minX; maxX >= m; m++ {
					if m == minX || m == maxX || n == minY || n == maxY {
						points = append(points, common.Point{X: uint32(m), Y: uint32(n)})
					}
				}
			}
		}
	}
	for j := range points {
		p := points[j]
		c := i.Map.GetCell(p.Coordinate())
		if c == nil || c.HasItem() {
			continue
		}
		i.CurrentLocation = p
		i.Map.AddObject(i)
		i.Broadcast(i.GetInfo())
		return "", true
	}
	return fmt.Sprintf("坐标(%s)附近没有合适的点放置物品", center.Coordinate()), false
}
