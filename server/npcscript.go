package main

import "github.com/yenkeia/mirgo/server/script"

func _CHECKPKPOINT(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return false
}

func _LEVEL(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return true
}

func _CHECKGOLD(npc *NPC, plr *Player, op script.CompareOp, v int) bool {
	return true
}

func _SET(npc *NPC, plr *Player, stat int, v int) {
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

func _CHECKQUEST(npc *NPC, plr *Player, quest int, stat int) bool {
	return true
}

func _LOCALMESSAGE(npc *NPC, plr *Player, message string, typ string) {

}

func _CHECK(npc *NPC, plr *Player, message string, typ string) bool {
	return false
}

func _ADDTOGUILD(npc *NPC, plr *Player, message string) {

}
func _ADDNAMELIST(npc *NPC, plr *Player, message string) {

}

func init() {
	// if
	script.Check("CHECKPKPOINT", _CHECKPKPOINT)
	script.Check("LEVEL", _CHECKPKPOINT)
	script.Check("CHECKGOLD", _CHECKGOLD)
	script.Check("CHECKITEM", _CHECKITEM)
	script.Check("CHECKQUEST", _CHECKQUEST)
	script.Check("INGUILD", _INGUILD)
	script.Check("CHECK", _CHECK)

	script.Action("GIVEBUFF", _GIVEBUFF)
	script.Action("SET", _SET)
	script.Action("MOVE", _MOVE)
	script.Action("TAKEGOLD", _TAKEGOLD)
	script.Action("LOCALMESSAGE", _LOCALMESSAGE)
	script.Action("ADDTOGUILD", _ADDTOGUILD)
	script.Action("ADDNAMELIST", _ADDNAMELIST)
}
