package main

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/yenkeia/mirgo/server/script"
)

func _CHECKPKPOINT(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return false
}

func _LEVEL(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return true
}

func _CHECKGOLD(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return true
}

func _GIVEBUFF(npc *NPC, plr *Player, bufname string, time int) {
}

func _MOVE(npc *NPC, plr *Player, mapname string, x, y int) {

}

func _TAKEGOLD(npc *NPC, plr *Player, gold int) {

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

func init() {
	script.AddParser(reflect.TypeOf(Flag(0)), parseFlag)
	script.AddParser(reflect.TypeOf(QuestStatus(0)), parseQuestStatus)

	// if
	script.Check("CHECKPKPOINT", _CHECKPKPOINT)
	script.Check("LEVEL", _CHECKPKPOINT)
	script.Check("CHECKGOLD", _CHECKGOLD)
	script.Check("CHECKITEM", _CHECKITEM, 1) // TODO: 验证默认是否为1
	script.Check("CHECKQUEST", _CHECKQUEST)
	script.Check("INGUILD", _INGUILD)
	script.Check("CHECK", _CHECK)

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
