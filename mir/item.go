package mir

import (
	"fmt"
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type Item struct {
	MapObject
	Gold     uint64
	UserItem *common.UserItem
}

func NewGold(dropper IMapObject, gold uint64) *Item {
	item := &Item{}
	item.ID = env.NewObjectID()
	item.Map = dropper.GetMap()
	item.Gold = gold

	return item
}

func NewItem(dropper IMapObject, ui *common.UserItem) *Item {
	item := &Item{UserItem: ui}
	item.Name = ui.Info.Name
	item.ID = env.NewObjectID()
	item.Map = dropper.GetMap()

	// if ui.IsAdded {
	// 	item.NameColor = Color.Cyan
	// } else {
	if ui.Info.Grade == common.ItemGradeNone {
		item.NameColor = common.ColorWhite
	}
	if ui.Info.Grade == common.ItemGradeCommon {
		item.NameColor = common.ColorWhite
	}
	if ui.Info.Grade == common.ItemGradeRare {
		item.NameColor = common.ColorDeepSkyBlue
	}
	if ui.Info.Grade == common.ItemGradeLegendary {
		item.NameColor = common.ColorDarkOrange
	}
	if ui.Info.Grade == common.ItemGradeMythical {
		item.NameColor = common.ColorPlum
	}
	// }

	return item
}

func (i *Item) GetMap() *Map {
	return i.Map
}

func (i *Item) GetID() uint32 {
	return i.ID
}

func (i *Item) GetName() string {
	return i.Name
}

func (m *Item) AddPlayerCount(n int) {
	m.PlayerCount += n
	switch m.PlayerCount {
	case 1:
		m.Map.AddActiveObj(m)
	case 0:
		m.Map.DelActiveObj(m)
	}
}

func (m *Item) Attacked(attacker IMapObject, damage int, dtype common.DefenceType, damageWeapon bool) int {
	return 0
}

func (m *Item) GetPlayerCount() int {
	return m.PlayerCount
}

func (i *Item) AttackMode() common.AttackMode {
	return common.AttackModePeace
}

func (i *Item) IsDead() bool { return i.Dead }

func (i *Item) IsUndead() bool {
	return false
}

func (i *Item) Process(dt time.Duration) {

}

func (i *Item) GetRace() common.ObjectType {
	return common.ObjectTypeItem
}

func (i *Item) IsBlocking() bool {
	return false
}

func (i *Item) GetPoint() common.Point {
	return i.CurrentLocation
}

func (i *Item) GetCell() *Cell {
	return i.Map.GetCell(i.CurrentLocation)
}

func (i *Item) Broadcast(msg interface{}) {
	i.Map.BroadcastP(i.CurrentLocation, msg, nil)
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

func (i *Item) AddBuff(buff *Buff) {}

func (i *Item) ApplyPoison(poison *Poison, caster IMapObject) {}

func (i *Item) GetItemInfo() *common.ItemInfo {
	return data.GetItemInfoByID(int(i.UserItem.ItemID))
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

	ok := false

	i.Map.RangeCell(center, distance, func(c *Cell, x, y int) bool {
		if c == nil || c.HasItem() {
			return true
		}

		ok = true
		i.CurrentLocation = common.NewPoint(x, y)
		i.Map.AddObject(i)
		i.Broadcast(i.GetInfo())

		return false
	})

	if !ok {
		return fmt.Sprintf("坐标(%s)附近没有合适的点放置物品", center), false
	}

	return "", true
}
