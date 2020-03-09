package mir

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/mir/script"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/ut"
)

type NPC struct {
	MapObject
	Image    int
	Light    uint8
	TurnTime time.Time
	Script   *script.Script
	Goods    []*common.UserItem
}

func NewNPC(m *Map, id uint32, ni *common.NpcInfo) *NPC {
	sc, err := script.LoadFile(ni.Filename + ".txt")
	if err != nil {
		log.Warnf("NPC [%s] [%s] 脚本加载失败: %s\n", ni.Name, ni.Filename, err.Error())
	}
	npc := &NPC{
		MapObject: MapObject{
			ID:               id,
			Name:             ni.Name,
			NameColor:        common.Color{R: 255, G: 255, B: 255},
			Map:              m,
			CurrentLocation:  common.NewPoint(ni.LocationX, ni.LocationY),
			CurrentDirection: common.MirDirection(ut.RandomInt(0, 1)),
		},
		Image:    ni.Image,
		Light:    0, // TODO
		TurnTime: time.Now(),
		Script:   sc,
		Goods:    []*common.UserItem{},
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
	p.Broadcast(p.GetInfo())
}

func (p *NPC) Spawned() {
	IMapObject_Spawned(p)
}

func (n *NPC) HasType(typ common.ItemType) bool {
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

func (n *NPC) AttackMode() common.AttackMode {
	return common.AttackModePeace
}

func (n *NPC) GetRace() common.ObjectType {
	return common.ObjectTypeMerchant
}

func (n *NPC) IsBlocking() bool {
	// return i.IsVisible()
	return false
}

func (n *NPC) GetPoint() common.Point {
	return n.CurrentLocation
}

func (n *NPC) GetCell() *Cell {
	return n.Map.GetCell(n.CurrentLocation)
}

func (n *NPC) GetDirection() common.MirDirection {
	return n.CurrentDirection
}

func (p *NPC) Attacked(attacker IMapObject, damageFinal int, defenceType common.DefenceType, damageWeapon bool) int {
	return 0
}

func (n *NPC) GetInfo() interface{} {
	res := &server.ObjectNPC{
		ObjectID:  n.ID,
		Name:      n.Name,
		NameColor: -16711936, // TODO
		Image:     uint16(n.Image),
		Color:     0, // TODO
		Location:  n.GetPoint(),
		Direction: n.GetDirection(),
		QuestIDs:  []int32{}, // TODO
	}
	return res
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
		n.TurnTime = time.Now().Add(time.Second * time.Duration(ut.RandomInt(20, 60)))
		n.CurrentDirection = common.MirDirection(ut.RandomInt(0, 1))
		n.Broadcast(ServerMessage{}.ObjectTurn(n))
	}
}

// GetUserItemByID 获取 NPC Goods
func (n *NPC) GetUserItemByID(id uint64) *common.UserItem {
	for _, v := range n.Goods {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// Buy 玩家向 NPC 购买物品
func (n *NPC) Buy(p *Player, userItemID uint64, count uint32) {
	userItem := n.GetUserItemByID(userItemID)
	if userItem == nil {
		return
	}
	itemInfo := data.GetItemInfoByID(int(userItem.ItemID))
	if itemInfo == nil {
		return
	}
	ui := env.NewUserItem(itemInfo)
	ui.Count = count
	p.GainItem(ui)
}

func (n *NPC) Sell(p *Player, item *common.UserItem) {
	// if (!BuyBack.ContainsKey(player.Name)) BuyBack[player.Name] = new List<UserItem>();

	// if (BuyBack[player.Name].Count >= Settings.GoodsBuyBackMaxStored)
	// 	BuyBack[player.Name].RemoveAt(0);

	// item.BuybackExpiryDate = Envir.Now;
	// BuyBack[player.Name].Add(item);
}

func (n *NPC) Craft(p *Player, index uint64, count uint32, slots []int) {
}
