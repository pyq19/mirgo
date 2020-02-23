package mir

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/yenkeia/mirgo/mir/script"
	"github.com/yenkeia/mirgo/proto/server"
)

func _CHECKPKPOINT(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return script.CompareInt(op, plr.PKPoints, v)
}

func _LEVEL(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return script.CompareInt(op, int(plr.Level), v)
}

func _CHECKGOLD(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	// FIX int64->int
	return script.CompareInt(op, int(plr.Gold), v)
}

func _GIVEBUFF(npc *NPC, plr *Player, bufname string, time int) {
}
func _REMOVEBUFF(npc *NPC, plr *Player, bufname string) {
}

func _MOVE(npc *NPC, plr *Player, mapname string, x, y int) {

}

func _TAKEGOLD(npc *NPC, plr *Player, gold int) {
	if gold < 0 || uint64(gold) > plr.Gold {
		log.Warnf("gold error")
	}
	plr.Gold -= uint64(gold)
	plr.Enqueue(&server.LoseGold{Gold: uint32(gold)})
}

func _INGUILD(npc *NPC, plr *Player) bool {
	return true
}

func _CHECKITEM(npc *NPC, plr *Player, itemname string, n int) bool {
	return true
}

func _CHECKQUEST(npc *NPC, plr *Player, quest int, stat QuestStatus) bool {
	return true
}

func _LOCALMESSAGE(npc *NPC, plr *Player, message string, typ string) {

}

func _SET(npc *NPC, plr *Player, flag Flag, v int) {
}

func _CHECK(npc *NPC, plr *Player, flag Flag, v int) bool {
	return false
}

func _ADDTOGUILD(npc *NPC, plr *Player, message string) {

}
func _ADDNAMELIST(npc *NPC, plr *Player, message string) {

}

func _GIVEITEM(npc *NPC, plr *Player, itemname string, n int) {

}

func _CLOSE(npc *NPC, plr *Player) {

}

func _REDUCEPKPOINT(npc *NPC, plr *Player, n int) {

}

func _CHECKMAP(npc *NPC, plr *Player, mapname string) bool {
	return false
}

func _ENTERMAP(npc *NPC, plr *Player) {
}

func _PETCOUNT(npc *NPC, plr *Player, op script.CompareOp, n int) bool {
	return false
}

func _CHECKGENDER(npc *NPC, plr *Player, g string) bool {
	return false
}
func _CHANGEGENDER(npc *NPC, plr *Player) {
}

func _CLEARPETS(npc *NPC, plr *Player) {
}
func _GIVEPET(npc *NPC, plr *Player, petname string) {
}

func _CHECKNAMELIST(npc *NPC, plr *Player, g string) bool {
	return false
}

func _REMOVENAMELIST(npc *NPC, plr *Player, g string) bool {
	return false
}

func _ISADMIN(npc *NPC, plr *Player) bool {
	return false
}

func _REMOVEFROMGUILD(npc *NPC, plr *Player, g string) {
}

func _CHECKHUM(npc *NPC, plr *Player, g1 string, g2 string) bool {
	return false
}

func _MONCLEAR(npc *NPC, plr *Player, g1 string) {
}

func _PARAM1(npc *NPC, plr *Player, g1 string) {
}
func _PARAM2(npc *NPC, plr *Player, g1 string) {
}
func _PARAM3(npc *NPC, plr *Player, g1 string) {
}

func _MONGEN(npc *NPC, plr *Player, g1 string, n int) {
}

func init() {
	script.AddParser(reflect.TypeOf(Flag(0)), parseFlag)
	script.AddParser(reflect.TypeOf(QuestStatus(0)), parseQuestStatus)

	// if
	script.Check("CHECKPKPOINT", _CHECKPKPOINT)
	script.Check("LEVEL", _LEVEL)
	script.Check("CHECKLEVEL", _LEVEL)
	script.Check("CHECKGOLD", _CHECKGOLD)
	script.Check("CHECKITEM", _CHECKITEM, 1) // TODO: 验证默认是否为1
	script.Check("CHECKQUEST", _CHECKQUEST)
	script.Check("INGUILD", _INGUILD)
	script.Check("CHECK", _CHECK)
	script.Check("CHECKMAP", _CHECKMAP)
	script.Check("PETCOUNT", _PETCOUNT)
	script.Check("CHECKGENDER", _CHECKGENDER)
	script.Check("CHECKNAMELIST", _CHECKNAMELIST)
	script.Check("ISADMIN", _ISADMIN)
	script.Check("CHECKHUM", _CHECKHUM)

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
	switch s {
	case "<$USERNAME>":
		return player.Name
	case "<$NPCNAME>":
		return npc.Name
	}

	log.Warnf("NPC 脚本缺少替换文本: %s %s\n", npc.Name, s)

	return s
}
