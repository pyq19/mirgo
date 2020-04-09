package game

import (
	"fmt"
	"reflect"

	"github.com/yenkeia/mirgo/game/script"
	"github.com/yenkeia/mirgo/util"
)

func _gmKill(p *Player, playername string) {
	if playername != "" {
		o := env.Players.GetPlayerByName(playername)
		if o == nil {
			p.ReceiveChat(fmt.Sprintf("找不到玩家(%s)", playername), util.ChatTypeSystem)
			return
		}
		o.Die()
		return
	}

	c := p.Map.GetNextCell(p.GetCell(), p.GetDirection(), 1)
	if c == nil {
		return
	}
	for _, o := range c.objects {
		if m, ok := o.(*Monster); ok {
			m.Die()
		}
	}
}

func _gmMake(p *Player, itemname string, n int) string {
	info := data.GetItemInfoByName(itemname)
	if info == nil {
		return "无该物品"
	}
	if n > 100 {
		return "大于100"
	}
	count := uint32(n)
	for count > 0 {
		if info.StackSize >= count {
			userItem := env.NewUserItem(info)
			userItem.Count = count
			p.GainItem(userItem)
			return "超过StackSize"
		}
		userItem := env.NewUserItem(info)
		userItem.Count = count
		count -= info.StackSize
		p.GainItem(userItem)
	}
	return fmt.Sprintf("%s x %d 创建成功", info.Name, count)
}

func _gmMap(p *Player) string {
	return fmt.Sprintf("当前地图: %s, ID: %d version:%d", p.Map.Info.Title, p.Map.Info.ID, p.Map.Version)
}

func _gmMove(p *Player, x, y int) {
	p.Teleport(p.Map, util.NewPoint(x, y))
}

func _gmMob(p *Player, monstername string) string {

	c := p.Map.GetNextCell(p.GetCell(), p.GetDirection(), 1)
	if c == nil || c.HasObject() {
		return "生成怪物失败"
	}
	mi := data.GetMonsterInfoByName(monstername)
	if mi == nil {
		return fmt.Sprintf("生成怪物失败，找不到怪物 %s", monstername)
	}
	mo := NewMonster(p.Map, c.Point, mi)
	mo.Spawn()
	p.Map.AddObject(mo)

	return "生成怪物成功"
}

func _gmInfo(p *Player) {

	c := p.Map.GetNextCell(p.GetCell(), p.GetDirection(), 1)
	if c == nil || c.objects == nil {
		return
	}
	for _, o := range c.objects {
		if o.GetRace() == util.ObjectTypeMonster {
			mo := o.(*Monster)
			p.ReceiveChat("--Monster Info--", util.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("ID: %d, Name: %s, AI: %d", mo.ID, mo.Name, mo.AI), util.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("Level: %d, Pos: %s", mo.Level, mo.GetPoint()), util.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("HP: %d, MinDC: %d, MaxDC: %d", mo.HP, mo.MinDC, mo.MaxDC), util.ChatTypeSystem2)
		}
		if o.GetRace() == util.ObjectTypePlayer {
			po := o.(*Player)
			p.ReceiveChat("--Player Info--", util.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("Name: %s, Level: %d, Pos: %s", po.Name, po.Level, po.GetPoint()), util.ChatTypeSystem2)
		}
	}
}

func _gmGold(p *Player, gold int) {
	p.GainGold(uint64(gold))
}

func _gmExp(p *Player, exp int) {
	p.GainExp(uint32(exp))
}

func _gmGiveSkill(p *Player, name string, level int) {
	info := data.GetMagicInfoByName(name)
	p.GiveSkill(util.Spell(info.Spell), level)
}

func _gmAllowTrade(p *Player) {
	p.AllowTrade = !p.AllowTrade
	if p.AllowTrade {
		p.ReceiveChat("你现在允许交易。", util.ChatTypeSystem)
	} else {
		p.ReceiveChat("你现在拒绝交易。", util.ChatTypeSystem)
	}
}

var cmd = script.NewContext()

func init() {
	cmd.AddParser(reflect.TypeOf((*Player)(nil)), nil)
	cmd.Action("KILL", _gmKill, "")
	cmd.Action("MAKE", _gmMake)
	cmd.Action("MAP", _gmMap)
	cmd.Action("INFO", _gmInfo)
	cmd.Action("MOB", _gmMob)
	cmd.Action("MOVE", _gmMove)
	cmd.Action("GOLD", _gmGold)
	cmd.Action("EXP", _gmExp)
	cmd.Action("GIVESKILL", _gmGiveSkill)
	cmd.Action("ALLOWTRADE", _gmAllowTrade)
}
