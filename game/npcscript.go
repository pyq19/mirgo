package game

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/script"
)

func _CHECKPKPOINT(npc *NPC, p *Player, op script.CompareOp, v int) bool {
	return script.CompareInt(op, p.PKPoints, v)
}

func _LEVEL(npc *NPC, p *Player, op script.CompareOp, v int) bool {
	return script.CompareInt(op, int(p.Level), v)
}

func _CHECKGOLD(npc *NPC, p *Player, op script.CompareOp, v int) bool {
	// FIX int64->int
	return script.CompareInt(op, int(p.Gold), v)
}

func _GIVEBUFF(npc *NPC, p *Player, bufname string, time int) {
}
func _REMOVEBUFF(npc *NPC, p *Player, bufname string) {
}

func _MOVE(npc *NPC, p *Player, mapname string, x, y int) {
	m := env.GetMapByName(strings.ToUpper(mapname))
	if x < 0 || y < 0 {
		// random teleport
	} else {
		p.Teleport(m, cm.NewPoint(x, y))
	}
}

func _TAKEGOLD(npc *NPC, p *Player, gold int) {
	if gold < 0 {
		log.Warnf("npc take gold error")
		return
	}
	p.TakeGold(uint64(gold))
}

func _CHECKITEM(npc *NPC, p *Player, itemname string, n int) bool {

	info := data.GetItemInfoByName(itemname)

	for _, item := range p.Inventory.Items {
		if item != nil && item.ItemID == info.ID {
			n -= int(item.Count)
			if n <= 0 {
				return true
			}
		}
	}

	return false
}

func _GIVEITEM(npc *NPC, p *Player, itemname string, n int) {
	count := uint32(n)
	info := data.GetItemInfoByName(itemname)
	if info == nil {
		return
	}
	for count > 0 {
		if info.StackSize >= count {
			userItem := env.NewUserItem(info)
			userItem.Count = count
			p.GainItem(userItem)
			return
		}
		userItem := env.NewUserItem(info)
		userItem.Count = count
		count -= info.StackSize
		p.GainItem(userItem)
	}
}

func _LINEMESSAGE(npc *NPC, p *Player, msg string, t string) {
	p.ReceiveChat(msg, cm.ChatTypeHint)
}
func _LOCALMESSAGE(npc *NPC, p *Player, msg string, typ string) {
	p.ReceiveChat(msg, cm.ChatTypeHint)
}

func _GIVEGOLD(npc *NPC, p *Player, v int) {
	p.GainGold(uint64(v))
}

func _TAKEITEM(npc *NPC, p *Player, itemname string, n int) {
	p.TakeItem(itemname, n)
}

func _CHECKQUEST(npc *NPC, p *Player, quest int, stat QuestStatus) bool {
	return true
}

func _SET(npc *NPC, p *Player, flag Flag, v int) {
}

func _CHECK(npc *NPC, p *Player, flag Flag, v int) bool {
	return false
}

func _ADDTOGUILD(npc *NPC, p *Player, message string) {

}
func _ADDNAMELIST(npc *NPC, p *Player, message string) {

}

func _INGUILD(npc *NPC, p *Player) bool {
	return true
}

func _CLOSE(npc *NPC, p *Player) {

}

func _REDUCEPKPOINT(npc *NPC, p *Player, n int) {

}

func _CHECKMAP(npc *NPC, p *Player, mapname string) bool {
	return false
}

func _ENTERMAP(npc *NPC, p *Player) {
}

func _PETCOUNT(npc *NPC, p *Player, op script.CompareOp, n int) bool {
	return false
}

func _CHECKGENDER(npc *NPC, p *Player, g string) bool {
	return false
}
func _CHANGEGENDER(npc *NPC, p *Player) {
}

func _CLEARPETS(npc *NPC, p *Player) {
}
func _GIVEPET(npc *NPC, p *Player, petname string) {
}

func _CHECKNAMELIST(npc *NPC, p *Player, g string) bool {
	return false
}

func _REMOVENAMELIST(npc *NPC, p *Player, g string) bool {
	return false
}

func _ISADMIN(npc *NPC, p *Player) bool {
	return false
}

func _REMOVEFROMGUILD(npc *NPC, p *Player, g string) {
}

func _CHECKHUM(npc *NPC, p *Player, g1 string, g2 string) bool {
	return false
}

func _MONCLEAR(npc *NPC, p *Player, g1 string) {
}

func _PARAM1(npc *NPC, p *Player, g1 string) {
}
func _PARAM2(npc *NPC, p *Player, g1 string) {
}
func _PARAM3(npc *NPC, p *Player, g1 string) {
}

func _MONGEN(npc *NPC, p *Player, g1 string, n int) {
}

func _GIVESKILL(npc *NPC, p *Player, name string, level int) {
	info := data.GetMagicInfoByName(name)
	p.GiveSkill(cm.Spell(info.Spell), level)
}

func _CHANGELEVEL(npc *NPC, p *Player, lv int) {
}

func _CHECKBUFF(npc *NPC, p *Player, bufname string) bool {
	return false
}

func _CHECKCLASS(npc *NPC, p *Player, classname string) bool {
	return false
}

func init() {
	script.AddParser(reflect.TypeOf((*Player)(nil)), nil)
	script.AddParser(reflect.TypeOf((*NPC)(nil)), nil)
	script.AddParser(reflect.TypeOf(Flag(0)), parseFlag)
	script.AddParser(reflect.TypeOf(QuestStatus(0)), parseQuestStatus)

	// if
	script.Check("CHECKPKPOINT", _CHECKPKPOINT)
	script.Check("LEVEL", _LEVEL)
	script.Check("CHECKLEVEL", _LEVEL)
	script.Check("CHECKGOLD", _CHECKGOLD)
	script.Check("CHECKITEM", _CHECKITEM, 1)
	script.Check("CHECKQUEST", _CHECKQUEST)
	script.Check("INGUILD", _INGUILD)
	script.Check("CHECK", _CHECK)
	script.Check("CHECKMAP", _CHECKMAP)
	script.Check("PETCOUNT", _PETCOUNT)
	script.Check("CHECKGENDER", _CHECKGENDER)
	script.Check("CHECKNAMELIST", _CHECKNAMELIST)
	script.Check("ISADMIN", _ISADMIN)
	script.Check("CHECKHUM", _CHECKHUM)
	script.Check("CHECKBUFF", _CHECKBUFF)
	script.Check("CHECKCLASS", _CHECKCLASS)

	script.Action("TAKEITEM", _TAKEITEM)
	script.Action("CHECKITEM", _CHECKITEM, 1) // GM-Manager.txt 32行：可能是配置写错了。
	script.Action("CHANGELEVEL", _CHANGELEVEL)
	script.Action("LINEMESSAGE", _LINEMESSAGE)
	script.Action("GIVEGOLD", _GIVEGOLD)
	script.Action("GIVESKILL", _GIVESKILL, 1)
	script.Action("PARAM1", _PARAM1)
	script.Action("PARAM2", _PARAM2)
	script.Action("PARAM3", _PARAM3)
	script.Action("MONGEN", _MONGEN)
	script.Action("REMOVENAMELIST", _REMOVENAMELIST)
	script.Action("MONCLEAR", _MONCLEAR)
	script.Action("REMOVEFROMGUILD", _REMOVEFROMGUILD)
	script.Action("GIVEPET", _GIVEPET)
	script.Action("CHANGEGENDER", _CHANGEGENDER)
	script.Action("REMOVEBUFF", _REMOVEBUFF)
	script.Action("CLEARPETS", _CLEARPETS)
	script.Action("ENTERMAP", _ENTERMAP)
	script.Action("REDUCEPKPOINT", _REDUCEPKPOINT)
	script.Action("CLOSE", _CLOSE)
	script.Action("GIVEBUFF", _GIVEBUFF)
	script.Action("SET", _SET)
	script.Action("MOVE", _MOVE, -1, -1)
	script.Action("TAKEGOLD", _TAKEGOLD)
	script.Action("LOCALMESSAGE", _LOCALMESSAGE)
	script.Action("ADDTOGUILD", _ADDTOGUILD)
	script.Action("ADDNAMELIST", _ADDNAMELIST)
	script.Action("GIVEITEM", _GIVEITEM)
}

//
type Flag int

// [888]
func parseFlag(s string) (reflect.Value, error) {
	if len(s) < 3 {
		return reflect.Value{}, errors.New("invalid flag:" + s)
	}

	return script.ParseInt(s[1 : len(s)-1])
}

type QuestStatus int

// 0,1,COMPLETE
func parseQuestStatus(s string) (reflect.Value, error) {
	if s == "COMPLETE" {
		s = "1"
	}

	return script.ParseInt(s)
}

var regNPCHotkey = regexp.MustCompile(`\<\$\w+\>`)

func replaceTemplates(npc *NPC, player *Player, say []string) []string {
	ret := make([]string, len(say))
	for i, v := range say {
		ret[i] = regNPCHotkey.ReplaceAllStringFunc(v, func(s string) string {
			return replaceTemplateName(npc, player, s)
		})
	}
	return ret
}

func replaceTemplateName(npc *NPC, player *Player, s string) string {

	// <$USERNAME>
	switch s[2 : len(s)-1] {
	case "USERNAME":
		return player.Name
	case "NPCNAME":
		return npc.Name
	case "PKPOINT":
		return fmt.Sprintf("%d", player.PKPoints)
	case "ARMOUR":
		return getEquipmentName(player, EquipmentSlotArmour)
	case "WEAPON":
		return getEquipmentName(player, EquipmentSlotWeapon)
	case "RING_L":
		return getEquipmentName(player, EquipmentSlotRingL)
	case "RING_R":
		return getEquipmentName(player, EquipmentSlotRingR)
	case "NECKLACE":
		return getEquipmentName(player, EquipmentSlotNecklace)
	case "BELT":
		return getEquipmentName(player, EquipmentSlotBelt)
	case "BOOTS":
		return getEquipmentName(player, EquipmentSlotBoots)
	case "STONE":
		return getEquipmentName(player, EquipmentSlotStone)
	case "HELMET":
		return getEquipmentName(player, EquipmentSlotHelmet)
	case "GAMEGOLD":
		return fmt.Sprintf("%d", player.Gold)
	case "HP":
		return fmt.Sprintf("%d", player.HP)
	case "MP":
		return fmt.Sprintf("%d", player.MP)
	case "MAXHP":
		return fmt.Sprintf("%d", player.MaxHP)
	case "MAXMP":
		return fmt.Sprintf("%d", player.MaxMP)
	case "LEVEL":
		return fmt.Sprintf("%d", player.Level)
	case "DATE":
		return time.Now().Format("2006-01-02")
	}

	log.Warnf("NPC 脚本缺少替换文本: %s %s\n", npc.Name, s)

	return s
}

func getEquipmentName(p *Player, slot int) string {
	item := p.Equipment.Items[slot]
	if item == nil {
		return "无"
	}

	return item.Info.Name
}
