package mir

import (
	"fmt"
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/mir/script"
	"github.com/yenkeia/mirgo/proto/server"
	"github.com/yenkeia/mirgo/setting"
)

type NPC struct {
	MapObject
	Image    int
	Light    uint8
	TurnTime time.Time
	Script   *script.Script
	Goods    []common.UserItem
}

func NewNPC(m *Map, ni *common.NpcInfo) *NPC {
	sc, err := script.LoadFile(setting.Conf.NPCDirPath + ni.Filename + ".txt")
	if err != nil {
		log.Warnf("NPC %s 脚本加载失败: %s\n", ni.Name, err.Error())
	}
	return &NPC{
		MapObject: MapObject{
			ID:               m.Env.NewObjectID(),
			Name:             ni.Name,
			NameColor:        common.Color{R: 255, G: 255, B: 255},
			Map:              m,
			CurrentLocation:  common.NewPoint(ni.LocationX, ni.LocationY),
			CurrentDirection: common.MirDirection(RandomInt(0, 1)),
		},
		Image:    ni.Image,
		Light:    0, // TODO
		TurnTime: time.Now(),
		Script:   sc,
		Goods:    make([]common.UserItem, 0),
	}
}

func (n *NPC) CallScript(p *Player, key string) ([]string, error) {
	say, err := n.Script.Call(n, p, key)
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

func (n *NPC) Process() {
	if n.TurnTime.Before(time.Now()) {
		n.TurnTime = time.Now().Add(time.Second * time.Duration(RandomInt(20, 60)))
		n.CurrentDirection = common.MirDirection(RandomInt(0, 1))
		n.Broadcast(ServerMessage{}.ObjectTurn(n))
	}
}

// GetUserItemByID 获取 NPC Goods
func (n *NPC) GetUserItemByID(id uint64) (item *common.UserItem) {
	for i := range n.Goods {
		if n.Goods[i].ID == id {
			return &n.Goods[i]
		}
	}
	return nil
}

// Buy 玩家向 NPC 购买物品
func (n *NPC) Buy(p *Player, userItemID uint64, count uint32) {
	env := n.Map.Env
	userItem := n.GetUserItemByID(userItemID)
	if userItem == nil {
		return
	}
	itemInfo := env.GameDB.GetItemInfoByID(int(userItem.ItemID))
	if itemInfo == nil {
		return
	}
	ui := env.NewUserItem(itemInfo)
	ui.Count = count
	p.GainItem(ui)
}
