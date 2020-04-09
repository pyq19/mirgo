package game

import (
	"fmt"
	"time"

	"github.com/yenkeia/mirgo/game/cm"
)

type Item struct {
	MapObject
	Gold     uint64
	UserItem *cm.UserItem
}

func NewGold(dropper IMapObject, gold uint64) *Item {
	item := &Item{}
	item.ID = env.NewObjectID()
	item.Map = dropper.GetMap()
	item.Gold = gold

	return item
}

func NewItem(dropper IMapObject, ui *cm.UserItem) *Item {
	item := &Item{UserItem: ui}
	item.Name = ui.Info.Name
	item.ID = env.NewObjectID()
	item.Map = dropper.GetMap()

	// if ui.IsAdded {
	// 	item.NameColor = Color.Cyan
	// } else {
	if ui.Info.Grade == cm.ItemGradeNone {
		item.NameColor = cm.ColorWhite
	}
	if ui.Info.Grade == cm.ItemGradeCommon {
		item.NameColor = cm.ColorWhite
	}
	if ui.Info.Grade == cm.ItemGradeRare {
		item.NameColor = cm.ColorDeepSkyBlue
	}
	if ui.Info.Grade == cm.ItemGradeLegendary {
		item.NameColor = cm.ColorDarkOrange
	}
	if ui.Info.Grade == cm.ItemGradeMythical {
		item.NameColor = cm.ColorPlum
	}
	// }

	return item
}

func (p *Item) Spawned() {
	IMapObject_Spawned(p)
}

func (p *Item) BroadcastHealthChange() {

}

func (i *Item) BroadcastInfo() {
	if i.UserItem == nil {
		i.Broadcast(ServerMessage{}.ObjectGold(i))
	} else {
		i.Broadcast(ServerMessage{}.ObjectItem(i))
	}
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

func (i *Item) GetLevel() int {
	return 0
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

func (m *Item) Attacked(attacker IMapObject, damage int, dtype cm.DefenceType, damageWeapon bool) int {
	return 0
}

func (m *Item) GetPlayerCount() int {
	return m.PlayerCount
}

func (i *Item) IsDead() bool { return i.Dead }

func (i *Item) IsUndead() bool {
	return false
}

func (i *Item) Process(dt time.Duration) {

}

func (i *Item) GetRace() cm.ObjectType {
	return cm.ObjectTypeItem
}

func (i *Item) IsBlocking() bool {
	return false
}

func (i *Item) GetPoint() cm.Point {
	return i.CurrentLocation
}

func (i *Item) GetCell() *Cell {
	return i.Map.GetCell(i.CurrentLocation)
}

func (i *Item) Broadcast(msg interface{}) {
	i.Map.BroadcastP(i.CurrentLocation, msg, nil)
}

func (i *Item) GetDirection() cm.MirDirection {
	return i.Direction
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

func (i *Item) GetItemInfo() *cm.ItemInfo {
	return data.GetItemInfoByID(int(i.UserItem.ItemID))
}

func (i *Item) GetImage() uint16 {
	info := i.GetItemInfo()
	switch info.Type {
	case cm.ItemTypeAmulet:
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
func (i *Item) Drop(center cm.Point, distance int) (string, bool) {

	ok := false

	i.Map.RangeCell(center, distance, func(c *Cell, x, y int) bool {
		if c == nil || c.HasItem() {
			return true
		}

		ok = true
		i.CurrentLocation = cm.NewPoint(x, y)
		i.Map.AddObject(i)
		i.BroadcastInfo()

		return false
	})

	if !ok {
		return fmt.Sprintf("坐标(%s)附近没有合适的点放置物品", center), false
	}

	return "", true
}
