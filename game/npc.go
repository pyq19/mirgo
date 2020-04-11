package game

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/script"
	"github.com/yenkeia/mirgo/game/util"
)

type NPC struct {
	MapObject
	Image    int
	Light    uint8
	TurnTime time.Time
	Script   *script.Script
	Goods    []*cm.UserItem
	BuyBack  map[uint32]*list.List
	Info     *cm.NpcInfo
}

func NewNPC(m *Map, id uint32, ni *cm.NpcInfo) *NPC {
	sc, err := script.LoadFile(ni.Filename + ".txt")
	if err != nil {
		log.Warnf("NPC [%s] [%s] 脚本加载失败: %s\n", ni.Name, ni.Filename, err.Error())
	}
	npc := &NPC{
		MapObject: MapObject{
			ID:              id,
			Name:            ni.ChineseName,
			NameColor:       cm.ColorLime,
			Map:             m,
			CurrentLocation: cm.NewPoint(ni.LocationX, ni.LocationY),
			Direction:       cm.MirDirection(util.RandomInt(0, 1)),
		},
		Image:    ni.Image,
		Light:    0, // TODO
		TurnTime: time.Now(),
		Script:   sc,
		Goods:    []*cm.UserItem{},
		BuyBack:  map[uint32]*list.List{},
	}

	for _, name := range npc.Script.Goods {
		res := strings.Split(name, " ")
		name := res[0]
		count := 1
		if len(res) == 2 {
			c, _ := strconv.Atoi(res[1])
			count = c
		}
		item := data.GetItemInfoByName(name)
		if item == nil {
			// FIXME 在 mir.sqlite item 表加上物品
			// log.Warnf("npc: %s 找不到 ItemInfo.Name = %s\n", npc.Name, name)
			continue
		}
		g := env.NewUserItem(item)
		g.Count = uint32(count)
		npc.Goods = append(npc.Goods, g)
	}

	return npc
}

func (p *NPC) BroadcastHealthChange() {

}

func (p *NPC) BroadcastInfo() {
	p.Broadcast(ServerMessage{}.ObjectNPC(p))
}

func (p *NPC) Spawned() {
	IMapObject_Spawned(p)
}

func (n *NPC) HasType(typ cm.ItemType) bool {
	if n.Script.Types != nil {
		for _, v := range n.Script.Types {
			if v == int(typ) {
				return true
			}
		}
	}
	return false
}

func (n *NPC) CallScript(p *Player, key string) ([]string, error) {
	say, err := n.Script.Call(key, n, p)
	if err != nil {
		return nil, err
	}
	return say, nil
}

func (n *NPC) IsDead() bool {
	return n.Dead
}

func (n *NPC) IsUndead() bool {
	return false
}

func (m *NPC) AddPlayerCount(n int) {
	m.PlayerCount += n
	switch m.PlayerCount {
	case 1:
		m.Map.AddActiveObj(m)
	case 0:
		m.Map.DelActiveObj(m)
	}
}

func (m *NPC) GetPlayerCount() int {
	return m.PlayerCount
}

func (i *NPC) GetMap() *Map {
	return i.Map
}

func (n *NPC) GetID() uint32 {
	return n.ID
}

func (n *NPC) GetName() string {
	return n.Name
}

func (n *NPC) GetLevel() int {
	return 0
}

func (n *NPC) GetRace() cm.ObjectType {
	return cm.ObjectTypeMerchant
}

func (n *NPC) IsBlocking() bool {
	// return i.IsVisible()
	return false
}

func (n *NPC) GetPoint() cm.Point {
	return n.CurrentLocation
}

func (n *NPC) GetCell() *Cell {
	return n.Map.GetCell(n.CurrentLocation)
}

func (n *NPC) GetDirection() cm.MirDirection {
	return n.Direction
}

func (p *NPC) Attacked(attacker IMapObject, damageFinal int, defenceType cm.DefenceType, damageWeapon bool) int {
	return 0
}

func (n *NPC) IsAttackTarget(IMapObject) bool {
	return false
}

func (n *NPC) IsFriendlyTarget(attacker IMapObject) bool {
	return true
}

func (n *NPC) GetBaseStats() BaseStats {
	return BaseStats{}
}

func (n *NPC) AddBuff(buff *Buff) {}

func (n *NPC) ApplyPoison(poison *Poison, caster IMapObject) {}

func (n *NPC) String() string {
	return fmt.Sprintf("NPC pos: %s, ID: %d, name: %s\n", n.GetPoint(), n.ID, n.Name)
}

func (n *NPC) Broadcast(msg interface{}) {
	n.Map.BroadcastP(n.CurrentLocation, msg, nil)
}

func (n *NPC) Process(dt time.Duration) {
	if n.TurnTime.Before(time.Now()) {
		n.TurnTime = time.Now().Add(time.Second * time.Duration(util.RandomInt(20, 60)))
		n.Direction = cm.MirDirection(util.RandomInt(0, 1))
		n.Broadcast(ServerMessage{}.ObjectTurn(n))
	}

	// TODO: 过期的buyback 放入商店供所有人购买
	for _, items := range n.BuyBack {
		for it := items.Front(); it != nil; {
			ele := it.Value.(*BuyBackItem)
			if ele.Expire <= dt {
				old := it
				it = it.Next()
				items.Remove(old)
			} else {
				ele.Expire -= dt
				it = it.Next()
			}
		}
	}
}

// GetUserItemByID 获取 NPC Goods
func (n *NPC) GetUserItemByID(id uint64) *cm.UserItem {
	for _, v := range n.Goods {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// Buy 玩家向 NPC 购买物品
func (n *NPC) Buy(p *Player, userItemID uint64, count uint32) {

	var userItem *cm.UserItem
	var iter *list.Element
	var isBuyBack bool

	items, has := n.BuyBack[p.ID]
	if has {
		for iter = items.Front(); iter != nil; iter = iter.Next() {
			if iter.Value.(*BuyBackItem).Item.ID == userItemID {
				userItem = iter.Value.(*BuyBackItem).Item
				isBuyBack = true
				break
			}
		}
		isBuyBack = false
	}

	if !isBuyBack {
		userItem = n.GetUserItemByID(userItemID)
	}

	if userItem == nil || count == 0 || count > userItem.Info.StackSize {
		return
	}

	price := userItem.Price()
	if price > p.Gold {
		return
	}

	if isBuyBack {
		count = userItem.Count
		items.Remove(iter)
		sendBuyBackGoods(p, n, false)
	} else {
		userItem = env.NewUserItem(userItem.Info)
		userItem.Count = count
	}
	if p.GainItem(userItem) {
		p.TakeGold(price)
	}
}

type BuyBackItem struct {
	Expire time.Duration
	Item   *cm.UserItem
}

func (n *NPC) GetPlayerBuyBack(p *Player) (ret []*cm.UserItem) {

	items, has := n.BuyBack[p.ID]
	if !has {
		return
	}
	for it := items.Front(); it != nil; it = it.Next() {
		ret = append(ret, it.Value.(*BuyBackItem).Item)
	}
	return
}

func (n *NPC) Sell(p *Player, item *cm.UserItem) {

	// TODO: config
	const GoodsBuyBackMaxStored = 20
	const GoodsBuyBackTime = 1 * time.Hour

	items, has := n.BuyBack[p.ID]
	if !has {
		items = list.New()
		n.BuyBack[p.ID] = items
	}

	if items.Len() >= GoodsBuyBackMaxStored {
		items.Remove(items.Front())
	}

	items.PushBack(&BuyBackItem{
		Item:   item,
		Expire: GoodsBuyBackTime,
	})
}

func (n *NPC) Craft(p *Player, index uint64, count uint32, slots []int) {
}

func (n *NPC) PriceRate(player *Player, baseRate bool) float32 {
	/* FIXME
	if n.Conq == nil || baseRate {
		return n.Info.Rate / 100.0
	}
	if player.MyGuild != nil && player.MyGuild.Guildindex == n.Conq.Owner {
		return n.Info.Rate / 100.0
	} else {
		return (((Info.Rate / 100.0) * Conq.npcRate) + Info.Rate) / 100.0
	}
	*/
	return float32(n.Info.Rate) / 100
}
