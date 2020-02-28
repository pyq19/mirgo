package mir

import (
	"fmt"
	"reflect"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/mir/script"
)

func _gmKill(p *Player, playername string) {
	if playername != "" {
		o := env.GetPlayerByName(playername)
		if o == nil {
			p.ReceiveChat(fmt.Sprintf("找不到玩家(%s)", playername), common.ChatTypeSystem)
			return
		}
		o.Die()
		return
	}

	c := p.Map.GetNextCell(p.GetCell(), p.GetDirection(), 1)
	if c == nil {
		return
	}
	c.Objects.Range(func(k, v interface{}) bool {
		if o, ok := v.(*Monster); ok {
			o.Die()
		}
		return true
	})
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
	return fmt.Sprintf("当前地图: %s, ID: %d", p.Map.Info.Title, p.Map.Info.ID)
}

func _gmMove(p *Player, x, y int) {
	p.Teleport(p.Map, common.NewPoint(x, y))
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
	p.Map.AddObject(NewMonster(p.Map, c.Point, mi))

	return ""
}

func _gmInfo(p *Player) {

	c := p.Map.GetNextCell(p.GetCell(), p.GetDirection(), 1)
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
			po := o.(*Player)
			p.ReceiveChat("--Player Info--", common.ChatTypeSystem2)
			p.ReceiveChat(fmt.Sprintf("Name: %s, Level: %d, Pos: %s", po.Name, po.Level, po.GetPoint()), common.ChatTypeSystem2)
		}
		return true
	})
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
}
