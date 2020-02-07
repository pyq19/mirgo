package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

const (
	LOGIN = iota
	SELECT
	GAME
	DISCONNECTED
)

// Player ...
type Player struct {
	AccountID int
	GameStage int
	Session   *cellnet.Session
	Char      *Character
	// MapObject
}

func (p *Player) Enqueue(msg interface{}) {
	if msg == nil {
		log.Errorln("warning: enqueue nil message")
		return
	}
	(*p.Session).Send(msg)
}

func (p *Player) ReceiveChat(text string, ct common.ChatType) {
	p.Enqueue(&server.Chat{
		Message: text,
		Type:    ct,
	})
}

func (p *Player) StartGame() {
	p.ReceiveChat("这是一个以学习为目的传奇服务端", common.ChatTypeSystem)
	p.ReceiveChat("如有任何建议、疑问欢迎交流", common.ChatTypeSystem)
	p.ReceiveChat("源码地址 https://github.com/yenkeia/mirgo", common.ChatTypeSystem)
	p.Char.EnqueueItemInfos()
	p.Char.RefreshStats()
	p.Char.EnqueueQuestInfo()
	p.Enqueue(ServerMessage{}.MapInformation(p.Char.Map.Info))
	p.Enqueue(ServerMessage{}.UserInformation(p.Char))
	p.Enqueue(ServerMessage{}.TimeOfDay(common.LightSettingDay))
	p.Char.EnqueueAreaObjects(nil, p.Char.Map.AOI.GetGridByPoint(p.Char.GetPoint()))
	p.Enqueue(ServerMessage{}.NPCResponse([]string{}))
	p.Char.Broadcast(ServerMessage{}.ObjectPlayer(p.Char))
}

func (p *Player) StopGame(reason int) {
	p.Char.Broadcast(ServerMessage{}.ObjectRemove(p.Char))
}

func (p *Player) Turn(direction common.MirDirection) {
	p.Char.Turn(direction)
}

func (p *Player) Walk(direction common.MirDirection) {
	p.Char.Walk(direction)
}

func (p *Player) Run(direction common.MirDirection) {
	p.Char.Run(direction)
}

func (p *Player) Chat(message string) {
	// private message
	if strings.HasPrefix(message, "/") {
		return
	}
	// group
	if strings.HasPrefix(message, "!!") {
		return
	}

	curMap := p.Char.Map

	if strings.HasPrefix(message, "@") {
		parts := strings.Split(message[1:], " ")
		switch strings.ToUpper(parts[0]) {
		case "LOGIN":
		case "KILL": // @kill 杀死面前的怪物，@kill name 杀死名字为 name 的玩家
			if len(parts) == 2 {
				o := curMap.Env.GetPlayerByName(parts[1])
				if o == nil {
					p.ReceiveChat(fmt.Sprintf("找不到玩家(%s)", parts[1]), common.ChatTypeSystem)
					return
				}
				o.Die()
				return
			}
			c := curMap.GetNextCell(p.Char.GetCell(), p.Char.GetDirection(), 1)
			if c == nil {
				return
			}
			c.Objects.Range(func(k, v interface{}) bool {
				if o, ok := v.(*Monster); ok {
					o.Die()
				}
				return true
			})
		case "RESTORE":
		case "CHANGEGENDER":
		case "LEVEL":
		case "MAKE": // @make 物品名 数量
			if len(parts) != 3 {
				return
			}
			info := curMap.Env.GameDB.GetItemInfoByName(parts[1])
			if info == nil {
				return
			}
			tmp, err := strconv.Atoi(parts[2])
			if err != nil || tmp > 100 {
				return
			}
			count := uint32(tmp)
			for count > 0 {
				if info.StackSize >= count {
					userItem := curMap.Env.NewUserItem(info)
					userItem.Count = count
					p.Char.GainItem(userItem)
					return
				}
				userItem := curMap.Env.NewUserItem(info)
				userItem.Count = count
				count -= info.StackSize
				p.Char.GainItem(userItem)
			}
			p.ReceiveChat(fmt.Sprintf("%s x %d 创建成功", info.Name, count), common.ChatTypeSystem)
		case "CLEARBUFFS":
		case "CLEARBAG":
		case "SUPERMAN":
		case "GAMEMASTER":
		case "OBSERVER":
		case "ALLOWGUILD":
		case "RECALL":
		case "ENABLEGROUPRECALL":
		case "GROUPRECALL":
		case "RECALLMEMBER":
		case "RECALLLOVER":
		case "TIME":
		case "ROLL":
		case "MAP":
			p.ReceiveChat(fmt.Sprintf("当前地图: %s, ID: %d", curMap.Info.Title, curMap.Info.ID), common.ChatTypeSystem)
		case "MOVE": // @move x y
			if len(parts) != 3 {
				p.ReceiveChat(fmt.Sprintf("移动失败，正确命令格式: @move 123 456"), common.ChatTypeSystem)
				return
			}
			x, err := strconv.Atoi(parts[1])
			if err != nil {
				p.ReceiveChat(fmt.Sprintf("移动失败，正确命令格式: @move 123 456"), common.ChatTypeSystem)
				return
			}
			y, err := strconv.Atoi(parts[2])
			if err != nil {
				p.ReceiveChat(fmt.Sprintf("移动失败，正确命令格式: @move 123 456"), common.ChatTypeSystem)
				return
			}
			p.Char.Teleport(curMap, common.NewPoint(x, y))
		case "MAPMOVE":
		case "GOTO":
		case "MOB": // @mob 怪物名称		在玩家周围生成 1 个怪物
			if len(parts) != 2 {
				p.ReceiveChat(fmt.Sprintf("生成怪物失败，正确命令格式: @mob 怪物名"), common.ChatTypeSystem)
				return
			}
			c := curMap.GetNextCell(p.Char.GetCell(), p.Char.GetDirection(), 1)
			if c == nil || c.HasObject() {
				p.ReceiveChat(fmt.Sprintf("生成怪物失败"), common.ChatTypeSystem)
				return
			}
			mi := curMap.Env.GameDB.GetMonsterInfoByName(parts[1])
			if mi == nil {
				p.ReceiveChat(fmt.Sprintf("生成怪物失败，找不到怪物 %s", parts[1]), common.ChatTypeSystem)
				return
			}
			curMap.AddObject(NewMonster(curMap, c.Point, mi))
		case "RECALLMOB":
		case "RELOADDROPS":
		case "RELOADNPCS":
		case "GIVEGOLD":
		case "GIVEPEARLS":
		case "GIVECREDIT":
		case "GIVESKILL":
		case "FIND":
		case "LEAVEGUILD":
		case "CREATEGUILD":
		case "ALLOWTRADE":
		case "TRIGGER":
		case "RIDE":
		case "SETFLAG":
		case "LISTFLAGS":
		case "CLEARFLAGS":
		case "CLEARMOB":
		case "CHANGECLASS": //@changeclass [Player] [Class]
		case "DIE":
		case "HAIR":
		case "DECO": //TEST CODE
		case "ADJUSTPKPOINT":
		case "ADDINVENTORY":
		case "ADDSTORAGE":
		case "INFO": // @info
			if len(parts) != 1 {
				return
			}
			c := curMap.GetNextCell(p.Char.GetCell(), p.Char.GetDirection(), 1)
			if c == nil {
				return
			}
			c.Objects.Range(func(k, v interface{}) bool {
				o := v.(IMapObject)
				if o.GetRace() == common.ObjectTypeMonster {
					mo := o.(*Monster)
					p.ReceiveChat("--Monster Info--", common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("ID: %d, Name: %s", mo.ID, mo.Name), common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("Level: %d, Pos: %s", mo.Level, mo.GetPoint()), common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("HP: %d, MinDC: %d, MaxDC: %d", mo.HP, mo.MinDC, mo.MaxDC), common.ChatTypeSystem2)
				}
				if o.GetRace() == common.ObjectTypePlayer {
					po := o.(*Character)
					p.ReceiveChat("--Player Info--", common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("Name: %s, Level: %d, Pos: %s", po.Name, po.Level, po.GetPoint()), common.ChatTypeSystem2)
				}
				return true
			})
		case "CLEARQUESTS":
		case "SETQUEST":
		case "TOGGLETRANSFORM":
		case "CREATEMAPINSTANCE": //TEST CODE
		case "STARTCONQUEST":
		case "RESETCONQUEST":
		case "GATES":
		case "CHANGEFLAG":
		case "CHANGEFLAGCOLOUR":
		case "REVIVE":
		case "DELETESKILL":
		}
		return
	}
	msg := ServerMessage{}.ObjectChat(p.Char, message, common.ChatTypeNormal)
	p.Enqueue(msg)
	p.Char.Broadcast(msg)
}

func (p *Player) MoveItem(mirGridType common.MirGridType, from int32, to int32) {
	p.Char.MoveItem(mirGridType, from, to)
}

func (p *Player) StoreItem(from int32, to int32) {
	p.Char.StoreItem(from, to)
}

func (p *Player) DepositRefineItem(from int32, to int32) {

}

func (p *Player) RetrieveRefineItem(from int32, to int32) {

}

func (p *Player) RefineCancel() {

}

func (p *Player) RefineItem(id uint64) {

}

func (p *Player) CheckRefine(id uint64) {

}

func (p *Player) ReplaceWeddingRing(id uint64) {

}

func (p *Player) DepositTradeItem(from int32, to int32) {

}

func (p *Player) RetrieveTradeItem(from int32, to int32) {

}

func (p *Player) TakeBackItem(from int32, to int32) {

}

func (p *Player) MergeItem(from common.MirGridType, to common.MirGridType, from2 uint64, to2 uint64) {

}

func (p *Player) EquipItem(mirGridType common.MirGridType, id uint64, to int32) {
	p.Char.EquipItem(mirGridType, id, to)
}

func (p *Player) RemoveItem(mirGridType common.MirGridType, id uint64, to int32) {
	p.Char.RemoveItem(mirGridType, id, to)
}

func (p *Player) RemoveSlotItem(grid common.MirGridType, id uint64, to int32, to2 common.MirGridType) {

}

func (p *Player) SplitItem(grid common.MirGridType, id uint64, count uint32) {

}

func (p *Player) UseItem(id uint64) {
	p.Char.UseItem(id)
}

func (p *Player) DropItem(id uint64, count uint32) {
	p.Char.DropItem(id, count)
}

func (p *Player) DropGold(gold uint64) {
	p.Char.DropGold(gold)
}

func (p *Player) PickUp() {
	p.Char.PickUp()
}

func (p *Player) Inspect(id uint32) {
}

func (p *Player) ChangeAMode(mode common.AttackMode) {

}

func (p *Player) ChangePMode(mode common.AttackMode) {

}

func (p *Player) ChangeTrade(trade bool) {

}

func (p *Player) Attack(direction common.MirDirection, spell common.Spell) {
	p.Char.Attack(direction, spell)
}

func (p *Player) RangeAttack(direction common.MirDirection, location common.Point, id uint32) {

}

func (p *Player) Harvest(direction common.MirDirection) {

}

func (p *Player) CallNPC(id uint32, key string) {
	npc := p.Char.Map.Env.GetNPC(id)
	if npc == nil {
		return
	}
	say, err := npc.CallScript(p, key)
	if err != nil {
		log.Warnf("NPC 脚本执行失败: %d %s %s\n", id, key, err.Error())
	}

	p.Enqueue(ServerMessage{}.NPCResponse(replaceTemplates(npc, p, say)))

	// ProcessSpecial
	switch strings.ToUpper(key) {
	case "[@BUY]":
		sendBuyKey(p, npc)
	case "[@SELL]":
		p.Enqueue(&server.NPCSell{})
	case "[@BUYSELL]":
		sendBuyKey(p, npc)
		p.Enqueue(&server.NPCSell{})
	default:
		// TODO
	}
}

func sendBuyKey(p *Player, npc *NPC) {

	goods := []*common.UserItem{}

	// TODO: fix..
	// for _, name := range npc.Script.Goods {
	// 	item := p.Map.Env.GameDB.GetItemInfoByName(name)
	// 	if item != nil {
	// 		p.EnqueueItemInfo(item.ID)
	// 		goods = append(goods, p.Map.Env.NewUserItem(item))
	// 	}
	// }

	p.Enqueue(&server.NPCGoods{
		Goods: goods,
		Rate:  1.0,
		Type:  common.PanelTypeBuy,
	})
}

func (p *Player) TalkMonsterNPC(id uint32) {

}

func (p *Player) BuyItem(index uint64, count uint32, panelType common.PanelType) {

}

func (p *Player) CraftItem() {

}

func (p *Player) SellItem(id uint64, count uint32) {

}

func (p *Player) RepairItem(id uint64) {

}

func (p *Player) BuyItemBack(id uint64, count uint32) {

}

func (p *Player) SRepairItem(id uint64) {

}

func (p *Player) Magic(spell common.Spell, direction common.MirDirection, targetID uint32, targetLocation common.Point) {
	p.Char.Magic(spell, direction, targetID, targetLocation)
}

func (p *Player) MagicKey(spell common.Spell, key uint8) {

}

func (p *Player) SwitchGroup(group bool) {

}

func (p *Player) AddMember(name string) {

}

func (p *Player) DelMember(name string) {

}

func (p *Player) GroupInvite(invite bool) {

}

func (p *Player) TownRevive() {

}

func (p *Player) SpellToggle(spell common.Spell, use bool) {

}

func (p *Player) ConsignItem(id uint64, price uint32) {

}

func (p *Player) MarketSearch(match string) {

}

func (p *Player) MarketRefresh() {

}

func (p *Player) MarketPage(page int32) {

}

func (p *Player) MarketBuy(id uint64) {

}

func (p *Player) MarketGetBack(id uint64) {

}

func (p *Player) RequestUserName(id uint32) {

}

func (p *Player) RequestChatItem(id uint64) {

}

func (p *Player) EditGuildMember(name string, name2 string, index uint8, changeType uint8) {

}
