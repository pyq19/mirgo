package mir

import (
	"fmt"

	"github.com/yenkeia/mirgo/common"
)

type AnyMap map[string]interface{}

type Bag struct {
	Type   common.UserItemType
	Player *Player
	Items  []*common.UserItem
}

func BagLoadFromDB(p *Player, typ common.UserItemType, n int) *Bag {
	b := &Bag{
		Player: p,
		Type:   typ,
	}
	b.Items = make([]*common.UserItem, n)

	cui := []*common.CharacterUserItem{}
	adb.Table("character_user_item").Where("character_id = ? AND type = ?", p.ID, typ).Find(&cui)
	ids := make([]int, n)

	userItemIDIndexMap := make(map[int]int)

	for i, item := range cui {
		ids[i] = item.UserItemID
		userItemIDIndexMap[item.UserItemID] = item.Index
	}

	items := []*common.UserItem{}
	adb.Table("user_item").Where("id in (?)", ids).Find(&items)

	for _, item := range items {
		item.Info = data.GetItemInfoByID(int(item.ItemID))
		b.Items[userItemIDIndexMap[int(item.ID)]] = item
	}

	return b
}

func (b *Bag) Move(from int, to int) error {
	if from < 0 || to < 0 || from > len(b.Items) || to > len(b.Items) {
		return fmt.Errorf("Move: 位置不存在 from=%d to=%d", from, to)
	}
	fromItem := b.Items[from]
	if fromItem == nil {
		return fmt.Errorf("格子 %d 没有物品", from)
	}
	toItem := b.Items[to]
	if toItem != nil {
		adb.Table("character_user_item").Where("user_item_id = ?", toItem.ID).Update("index", from)
	}
	adb.Table("character_user_item").Where("user_item_id = ?", fromItem.ID).Update("index", to)

	b.Items[from], b.Items[to] = b.Items[to], b.Items[from]
	return nil
}

func (b *Bag) Set(i int, item *common.UserItem) {

	if item != nil {
		if b.Items[i] != nil {
			log.Errorln("该位置有物品了")
		}

		adb.Table("user_item").Create(item)
		adb.Table("character_user_item").Create(&common.CharacterUserItem{
			CharacterID: int(b.Player.ID),
			UserItemID:  int(item.ID),
			Type:        int(b.Type),
			Index:       i,
		})
		b.Items[i] = item
	} else {
		item = b.Items[i]
		if item != nil {
			adb.Table("user_item").Where("id = ?", item.ID).Delete(&common.UserItem{})
			adb.Table("character_user_item").Where("user_item_id = ?", item.ID).Delete(&common.CharacterUserItem{})
		} else {
			log.Errorln("尝试删除空位置的物品")
		}
		b.Items[i] = nil
	}
}

func (b *Bag) Get(i int) *common.UserItem {
	return b.Items[i]
}

func (b *Bag) SetCount(i int, c uint32) {
	if c == 0 {
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

	item := b.Items[from]
	adb.Table("character_user_item").Where("user_item_id = ?", item.ID).Update(AnyMap{
		"type":  tobag.Type,
		"index": to,
	})
	b.Items[from] = nil

	toItem := tobag.Items[to]
	if toItem != nil {
		adb.Table("character_user_item").Where("user_item_id = ?", toItem.ID).Update(AnyMap{
			"type":  b.Type,
			"index": from,
		})
	}
	tobag.Items[to] = item

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
