package game

import (
	"fmt"
	"reflect"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/script"
)

func _gmKill(p *Player, playername string) {
	if playername != "" {
		o := env.Players.GetPlayerByName(playername)
		if o == nil {
			p.ReceiveChat(fmt.Sprintf("找不到玩家(%s)", playername), cm.ChatTypeSystem)
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
	p.Teleport(p.Map, cm.NewPoint(x, y))
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
		if o.GetRace() == cm.ObjectTypeMonster {
			mo := o.(*Monster)
			p.ReceiveChat("--Monster Info--", cm.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("ID: %d, Name: %s, AI: %d", mo.ID, mo.Name, mo.AI), cm.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("Level: %d, Pos: %s", mo.Level, mo.GetPoint()), cm.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("HP: %d, MinDC: %d, MaxDC: %d", mo.HP, mo.MinDC, mo.MaxDC), cm.ChatTypeSystem2)
		}
		if o.GetRace() == cm.ObjectTypePlayer {
			po := o.(*Player)
			p.ReceiveChat("--Player Info--", cm.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("Name: %s, Level: %d, Pos: %s", po.Name, po.Level, po.GetPoint()), cm.ChatTypeSystem2)
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
	p.GiveSkill(cm.Spell(info.Spell), level)
}

func _gmAllowTrade(p *Player) {
	p.AllowTrade = !p.AllowTrade
	if p.AllowTrade {
		p.ReceiveChat("你现在允许交易。", cm.ChatTypeSystem)
	} else {
		p.ReceiveChat("你现在拒绝交易。", cm.ChatTypeSystem)
	}
}

func _gmDie(p *Player) {
	p.Die()
}

func _gmLeaveGuild(p *Player) {
	if p.MyGuild == nil {
		return
	}
	if p.MyGuildRank == nil {
		return
	}
	if p.MyGuild.IsAtWar() {
		p.ReceiveChat("在战争中不能离开行会。", cm.ChatTypeSystem)
		return
	}
	p.MyGuild.DeleteMember(p, p.Name)
}

func _gmCreateGuild(p *Player, guildname string) {
	player := p
	gName := guildname
	if player.MyGuild != nil {
		p.ReceiveChat(fmt.Sprintf("玩家 %s 已经在一个行会中了。", player.Name), cm.ChatTypeSystem)
		return
	}
	if (len(gName) < 3) || (len(gName) > 20) {
		p.ReceiveChat("行会名字必须在3-20字符之间。", cm.ChatTypeSystem)
		return
	}
	guild := env.GetGuild(gName)
	if guild != nil {
		p.ReceiveChat(fmt.Sprintf("行会 %s 已存在。", gName), cm.ChatTypeSystem)
		return
	}
	player.CanCreateGuild = true
	if player.CreateGuild(gName) {
		p.ReceiveChat(fmt.Sprintf("成功创建行会 %s 。", gName), cm.ChatTypeSystem)
	} else {
		p.ReceiveChat("创建行会失败。", cm.ChatTypeSystem)
	}
	player.CanCreateGuild = false
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
	cmd.Action("DIE", _gmDie)
	cmd.Action("LEAVEGUILD", _gmLeaveGuild)
	cmd.Action("CREATEGUILD", _gmCreateGuild)
}
