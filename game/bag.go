package game

import (
	"fmt"

	"github.com/yenkeia/mirgo/game/cm"
)

type AnyMap map[string]interface{}

type Bag struct {
	Type   cm.UserItemType
	Player *Player
	Items  []*cm.UserItem
	// ItemCount int
}

func NewBag(p *Player, typ cm.UserItemType, n int) *Bag {
	b := &Bag{
		Player: p,
		Type:   typ,
	}
	b.Items = make([]*cm.UserItem, n)
	return b
}

func BagLoadFromDB(p *Player, typ cm.UserItemType, n int) *Bag {
	b := NewBag(p, typ, n)

	cui := []*cm.CharacterUserItem{}
	adb.Table("character_user_item").Where("character_id = ? AND type = ?", p.ID, typ).Find(&cui)
	ids := make([]int, n)

	userItemIDIndexMap := make(map[int]int)

	for i, item := range cui {
		ids[i] = item.UserItemID
		userItemIDIndexMap[item.UserItemID] = item.Index
	}

	items := []*cm.UserItem{}
	adb.Table("user_item").Where("id in (?)", ids).Find(&items)

	for _, item := range items {
		item.Info = data.GetItemInfoByID(int(item.ItemID))
		b.Items[userItemIDIndexMap[int(item.ID)]] = item
	}

	return b
}

func (b *Bag) ItemCount() int {
	cnt := 0
	for i := 0; i < len(b.Items); i++ {
		tmp := b.Items[i]
		if tmp == nil {
			continue
		}
		cnt++
	}
	return cnt
}

func (b *Bag) Length() int {
	return len(b.Items)
}

func (b *Bag) Move(from int, to int) error {
	return b.MoveTo(from, to, b)
}

func (b *Bag) Set(i int, item *cm.UserItem) {

	if item != nil {
		if b.Items[i] != nil {
			log.Errorln("该位置有物品了")
		}

		adb.Table("user_item").Create(item)
		adb.Table("character_user_item").Create(&cm.CharacterUserItem{
			CharacterID: int(b.Player.ID),
			UserItemID:  int(item.ID),
			Type:        int(b.Type),
			Index:       i,
		})
		b.Items[i] = item
	} else {
		item = b.Items[i]
		if item != nil {
			adb.Table("user_item").Where("id = ?", item.ID).Delete(&cm.UserItem{})
			adb.Table("character_user_item").Where("user_item_id = ?", item.ID).Delete(&cm.CharacterUserItem{})
		} else {
			log.Errorln("尝试删除空位置的物品")
		}
		b.Items[i] = nil
	}
}

func (b *Bag) Get(i int) *cm.UserItem {
	return b.Items[i]
}

func (b *Bag) SetCount(i int, c uint32) {
	if c == 0 {
		log.Infof("Delete UserItem %d \n", b.Items[i].ID)
		// adb.Table("user_item").Where("id = ?", b.Items[i].ID).Delete(&cm.UserItem{})
		// adb.Table("character_user_item").Where("user_item_id = ?", b.Items[i].ID).Delete(&cm.CharacterUserItem{})
		b.Set(i, nil)
	} else {
		adb.Table("user_item").Where("id = ?", b.Items[i].ID).Update("count", c)
		b.Items[i].Count = c
	}
}

func (b *Bag) UseCount(i int, c uint32) {
	b.SetCount(i, b.Items[i].Count-c)
}

func (b *Bag) MoveTo(from, to int, tobag *Bag) error {
	if from < 0 || to < 0 || from > len(b.Items) || to > len(tobag.Items) {
		return fmt.Errorf("Move: 位置不存在 from=%d to=%d", from, to)
	}

	item := b.Items[from]
	if item == nil {
		return fmt.Errorf("格子 %d 没有物品", from)
	}
	adb.Table("character_user_item").Where("user_item_id = ?", item.ID).Update(AnyMap{
		"type":  tobag.Type,
		"index": to,
	})

	toItem := tobag.Items[to]
	if toItem != nil {
		adb.Table("character_user_item").Where("user_item_id = ?", toItem.ID).Update(AnyMap{
			"type":  b.Type,
			"index": from,
		})
	}

	b.Items[from], tobag.Items[to] = tobag.Items[to], b.Items[from]

	return nil
}

func (b *Bag) EmptySlot(start int) int {
	for i := start; i < len(b.Items); i++ {
		tmp := b.Items[i]
		if tmp == nil {
			return i
		}
	}
	return -1
}
