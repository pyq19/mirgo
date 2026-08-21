package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// GM 命令实现。注册表/分发机制见 gmcmd_registry.go。
// 注册顺序 = Delphi ProcessUserLineMsg 分发顺序（ObjBase.pas:7188-8600）：
// 用户命令段(7362-7780) → GM 命令段(7809-8253)。

func init() {
	// --- 用户命令段 ---
	registerGMCommand(&gmCmdDef{key: "UserMove", name: "Move", permMin: 0, usage: "X Y | 地图号", handler: (*PlayObject).cmdMove})
	registerGMCommand(&gmCmdDef{key: "Searching", name: "Searching", permMin: 0, usage: "人物名称", handler: (*PlayObject).cmdSearching})
	registerGMCommand(&gmCmdDef{key: "GuildReCall", name: "GuildRecall", permMin: 0, usage: "行会传送，行会掌门人可以将整个行会成员全部集中。", handler: (*PlayObject).cmdGuildRecall})
	registerGMCommand(&gmCmdDef{key: "Recall", name: "Recall", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdRecallHuman})
	registerGMCommand(&gmCmdDef{key: "Rest", name: "Rest", permMin: 0, handler: (*PlayObject).cmdRest})
	registerGMCommand(&gmCmdDef{key: "Slave", name: "slave", permMin: 0, usage: "[relax|recall]", handler: (*PlayObject).cmdSlave}) // Go 专属宠物管理
	registerGMCommand(&gmCmdDef{key: "DearRecall", name: "DearRecall", permMin: 0, handler: (*PlayObject).cmdDearRecall})             // Go 恢复：1.76 Delphi 已注释
	registerGMCommand(&gmCmdDef{key: "MasterRecall", name: "MasterRecall", permMin: 0, handler: (*PlayObject).cmdMasterRecall})       // Go 恢复：1.76 Delphi 已注释
	registerGMCommand(&gmCmdDef{key: "StorageSetPassword", name: "SetPassword", permMin: 0, usage: "密码(4-7位)", handler: (*PlayObject).cmdSetStoragePwd})
	registerGMCommand(&gmCmdDef{key: "StorageChgPassword", name: "ChgPassword", permMin: 0, usage: "旧密码 新密码", handler: (*PlayObject).cmdChgStoragePwd})
	registerGMCommand(&gmCmdDef{key: "StorageUnLock", name: "UnLockStorage", permMin: 0, usage: "密码", handler: (*PlayObject).cmdUnlockStorage})
	registerGMCommand(&gmCmdDef{key: "TakeOnHorse", name: "骑马", permMin: 0, handler: (*PlayObject).cmdTakeOnHorse})
	registerGMCommand(&gmCmdDef{key: "TakeOffHorse", name: "下马", permMin: 0, handler: (*PlayObject).cmdTakeOffHorse})
	registerGMCommand(&gmCmdDef{key: "AttackMode", name: "AttackMode", permMin: 0, usage: "[模式(0-4)]", handler: (*PlayObject).cmdAttackMode})

	// --- GM 命令段 ---
	registerGMCommand(&gmCmdDef{key: "Move", name: "MoveGM", permMin: 3, handler: nil}) // 无独立分发：并入 UserMove 按参数形态路由
	registerGMCommand(&gmCmdDef{key: "Info", name: "Info", permMin: 3, usage: "人物名称", handler: (*PlayObject).cmdInfo})
	registerGMCommand(&gmCmdDef{key: "Map", name: "Map", permMin: 3, handler: (*PlayObject).cmdMap})
	registerGMCommand(&gmCmdDef{key: "MapInfo", name: "MapInfo", permMin: 10, usage: "地图号 X Y", handler: (*PlayObject).cmdMapInfo})
	registerGMCommand(&gmCmdDef{key: "Make", name: "make", permMin: 10, usage: "物品名称 [数量(1-10)]", handler: (*PlayObject).cmdMake})
	registerGMCommand(&gmCmdDef{key: "Level", name: "Level", permMin: 10, usage: "等级(1-65535)", handler: (*PlayObject).cmdLevel})
	registerGMCommand(&gmCmdDef{key: "Mob", name: "Mob", permMin: 10, usage: "怪物名称 [数量(1-64)] [等级(0-10)]", handler: (*PlayObject).cmdMob})
	registerGMCommand(&gmCmdDef{key: "RecallMob", name: "RecallMob", permMin: 10, usage: "怪物名称 [数量(1-20)] [等级]", handler: (*PlayObject).cmdRecallMob})
	registerGMCommand(&gmCmdDef{key: "ClearMon", name: "ClearMon", permMin: 10, usage: "地图号 怪物名称|* [掉物品(0|1)]", handler: (*PlayObject).cmdClearMon})
	registerGMCommand(&gmCmdDef{key: "Gold", name: "gold", permMin: 10, usage: "数量", handler: (*PlayObject).cmdGold}) // Go 专属：给自己
	registerGMCommand(&gmCmdDef{key: "AddGold", name: "AddGold", permMin: 10, usage: "人物名称 数量", handler: (*PlayObject).cmdAddGold})
	registerGMCommand(&gmCmdDef{key: "DelGold", name: "DelGold", permMin: 10, usage: "人物名称 数量", handler: (*PlayObject).cmdDelGold})
	registerGMCommand(&gmCmdDef{key: "Heal", name: "heal", permMin: 10, handler: (*PlayObject).cmdHeal}) // Go 专属
	registerGMCommand(&gmCmdDef{key: "PkPoint", name: "PKpoint", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdPkPoint})
	registerGMCommand(&gmCmdDef{key: "IncPkPoint", name: "IncPkPoint", permMin: 10, usage: "人物名称 PK点数", handler: (*PlayObject).cmdIncPkPoint})
	registerGMCommand(&gmCmdDef{key: "FreePenalty", name: "FreePK", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdFreePK})
	registerGMCommand(&gmCmdDef{key: "Kick", name: "Kick", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdKick})
	registerGMCommand(&gmCmdDef{key: "ReAlive", name: "ReAlive", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdReAlive})
	registerGMCommand(&gmCmdDef{key: "Revive", name: "revive", permMin: 10, usage: "[人物名称]", handler: (*PlayObject).cmdRevive}) // Go 扩展：无参复活自己
	registerGMCommand(&gmCmdDef{key: "SuperMan", name: "Superman", permMin: 10, handler: (*PlayObject).cmdSuperman})
	registerGMCommand(&gmCmdDef{key: "GameMaster", name: "GameMaster", permMin: 10, handler: (*PlayObject).cmdGameMaster})
	registerGMCommand(&gmCmdDef{key: "ObServer", name: "Observer", permMin: 10, handler: (*PlayObject).cmdObserver})
	registerGMCommand(&gmCmdDef{key: "ReGoto", name: "ReGoto", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdReGoto})
	registerGMCommand(&gmCmdDef{key: "Shutup", name: "Shutup", permMin: 10, usage: "人物名称 [分钟]", handler: (*PlayObject).cmdShutup})
	registerGMCommand(&gmCmdDef{key: "ReleaseShutup", name: "ReleaseShutup", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdReleaseShutup})
	registerGMCommand(&gmCmdDef{key: "ShutupList", name: "ShutupList", permMin: 10, handler: (*PlayObject).cmdShutupList})
	registerGMCommand(&gmCmdDef{key: "MoveUser", name: "moveuser", permMin: 10, usage: "人物名称 地图号 [X Y]", handler: (*PlayObject).cmdMoveUser}) // Go 专属
	registerGMCommand(&gmCmdDef{key: "ReloadNpc", name: "ReloadNpc", permMin: 6, permFail: permFailSilent, usage: "[all]", handler: (*PlayObject).cmdReloadNpc})
	registerGMCommand(&gmCmdDef{key: "Status", name: "status", permMin: 10, handler: (*PlayObject).cmdStatus}) // Go 专属
	registerGMCommand(&gmCmdDef{key: "NoMob", name: "nomob", permMin: 10, usage: "[on|off]", handler: (*PlayObject).cmdNoMob}) // Go 专属：Delphi 无此命令
}

// gmPermByKey 按 Delphi Command.ini 键取已解析权限（供跨命令引用，
// 如 UserMove 内部检查 MOVE 的权限）。
func gmPermByKey(key string, fallback int) int {
	for _, def := range gmCmdDefs {
		if def.key == key {
			return def.resolvedPerm
		}
	}
	return fallback
}

func (p *PlayObject) sysMsgBlue(server *netserver.TCPServer, text string) {
	msg := protocol.MakeDefaultMsg(protocol.SMSysMessage, 0, chatColorSysBlue, 0, 0)
	server.Send(p.Session.ID, msg, protocol.EncodeString(text))
}

func (p *PlayObject) sysMsgGreen(server *netserver.TCPServer, text string) {
	msg := protocol.MakeDefaultMsg(protocol.SMSysMessage, 0, chatColorSysGreen, 0, 0)
	server.Send(p.Session.ID, msg, protocol.EncodeString(text))
}

// teleportPlayerTo 把玩家移到目标地图坐标（跨图走 EnterAnotherMap，
// 同图走 RemoveObject/AddObject + RM_DISAPPEAR/RM_LOGON）。
func teleportPlayerTo(pl *PlayObject, server *netserver.TCPServer, env *Environment, x, y int) {
	if pl.envir != env || pl.MapName != env.Name {
		pl.EnterAnotherMap(server, env, x, y)
		return
	}
	pl.envir.RemoveObject(pl.CurrX, pl.CurrY, OS_MOVINGOBJECT, pl)
	pl.envir.broadcastRefMsg(pl.BaseObject, RM_DISAPPEAR, pl.ID, pl.CurrX, pl.CurrY, pl.Dir)
	pl.CurrX, pl.CurrY = x, y
	pl.envir.AddObject(x, y, OS_MOVINGOBJECT, pl)
	pl.envir.broadcastRefMsg(pl.BaseObject, RM_LOGON, pl.ID, x, y, pl.Dir)
}

// findWalkableNear 在 (cx,cy) 半径 radius 方框内找可走格子
// （Delphi sub_4C5370，ObjBase.pas:26709）。
func findWalkableNear(env *Environment, cx, cy, radius int) (int, int, bool) {
	if env == nil {
		return 0, 0, false
	}
	for r := 0; r <= radius; r++ {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				if r > 0 && dx != r && dx != -r && dy != r && dy != -r {
					continue // 按圈层由内向外
				}
				if env.CanWalk(cx+dx, cy+dy) {
					return cx + dx, cy + dy, true
				}
			}
		}
	}
	return 0, 0, false
}

// randomMapPosition 地图内随机可走落点（Delphi MapRandomMove 简化）。
func randomMapPosition(env *Environment, tries int) (int, int) {
	for i := 0; i < tries; i++ {
		x := rand.Intn(env.Width)
		y := rand.Intn(env.Height)
		if env.CanWalk(x, y) {
			return x, y
		}
	}
	return env.Width / 2, env.Height / 2
}

// cmdMove 统一 @Move：数字参数=传送戒指坐标版（USERMOVE，ObjBase.pas:7770→15190），
// 地图名参数=GM 地图版（MOVE，ObjBase.pas:7814→13450）。Delphi 默认配置两个字面量
// 同为 'Move' 且 USERMOVE 遮蔽 MOVE；Go 按参数形态路由保持两者可用。
func (p *PlayObject) cmdMove(server *netserver.TCPServer, params []string) bool {
	if len(params) >= 2 {
		if !p.HasTeleport {
			p.sysMsg(server, "您现在还无法使用此功能！！！")
			return true
		}
		p.ringTeleport(server, params[0], params[1])
		return true
	}
	if len(params) == 1 {
		if int(p.Permission) < gmPermByKey("Move", 3) {
			p.sysMsg(server, gmMsgPermTooLow)
			return true
		}
		mapName := params[0]
		env := p.MapMgr.FindMap(mapName)
		if env == nil {
			p.sysMsg(server, mapName+" 此地图号不存在.")
			return true
		}
		x, y := randomMapPosition(env, 20)
		p.EnterAnotherMap(server, env, x, y)
		p.sysMsg(server, "传送到: "+mapName)
	}
	return true
}

// ringTeleport 传送戒指坐标传送（Delphi CmdUserMoveXY，ObjBase.pas:15190-15231）。
func (p *PlayObject) ringTeleport(server *netserver.TCPServer, sx, sy string) {
	if p.envir != nil && p.envir.Flag.NoPositionMove {
		p.sysMsg(server, "此地图禁止使用此命令！！！")
		return
	}
	now := time.Now().UnixMilli()
	// Delphi 每玩家 m_dwTeleportTick，dwUserMoveTime 默认 10 秒。
	if now-p.TeleportTick < 10000 {
		left := 10 - (now-p.TeleportTick)/1000
		p.sysMsg(server, strconv.FormatInt(left, 10)+"秒之后才可以再使用此功能！！！")
		return
	}
	x, err1 := strconv.Atoi(sx)
	y, err2 := strconv.Atoi(sy)
	if err1 != nil || err2 != nil || x < 0 || y < 0 {
		p.sysMsg(server, "格式: @move X Y")
		return
	}
	if p.envir != nil && p.envir.CanWalk(x, y) {
		p.TeleportTick = now
		teleportPlayerTo(p, server, p.envir, x, y)
		p.sysMsg(server, "已传送")
	} else {
		p.sysMsg(server, fmt.Sprintf("无法移动到地图: %s X:%s Y:%s", p.MapName, sx, sy))
	}
}

// cmdSearching 探测项链定位玩家（Delphi CmdSearchHuman，ObjBase.pas:14352-14382）。
func (p *PlayObject) cmdSearching(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	if !p.HasProbe && int(p.Permission) < 6 {
		p.sysMsg(server, "您现在还无法使用此功能！！！")
		return true
	}
	now := time.Now().UnixMilli()
	if now-p.SearchTick < 10000 { // Delphi 10 秒冷却
		return true
	}
	p.SearchTick = now
	if p.Engine == nil {
		return true
	}
	p.Engine.mu.Lock()
	var found *PlayObject
	for _, pl := range p.Engine.PlayObjectList {
		if pl.Name == params[0] {
			found = pl
			break
		}
	}
	p.Engine.mu.Unlock()
	if found != nil {
		p.sysMsgBlue(server, fmt.Sprintf("%s 现在位于 %s %d:%d", found.Name, found.MapName, found.CurrX, found.CurrY))
	} else {
		p.sysMsg(server, fmt.Sprintf("%s 现在不在线，或在其它服务器上！！！", params[0]))
	}
	return true
}

// cmdGuildRecall 行会合一（Delphi CmdGuildRecall，ObjBase.pas:12997-13087）。
func (p *PlayObject) cmdGuildRecall(server *netserver.TCPServer, params []string) bool {
	// Delphi: 行会传送装备(m_boGuildMove)或 GM≥6。
	if !p.HasGuildMove && int(p.Permission) < 6 {
		p.sysMsg(server, "您现在还无法使用此功能！！！")
		return true
	}
	guild := p.Engine.PlayerGuild(p.Name)
	if guild == nil || guild.Leader != p.Name {
		p.sysMsg(server, "行会掌门人才可以使用此功能！！！")
		return true
	}
	if p.envir != nil && p.envir.Flag.NoGuildRecall {
		p.sysMsg(server, "本地图不允许使用此功能！！！")
		return true
	}
	now := time.Now().UnixMilli()
	// Delphi nGuildRecallTime 默认 180 秒（M2Share.pas:2009）。
	if now-p.GuildRecallTick < 180000 {
		left := 180 - (now-p.GuildRecallTick)/1000
		p.sysMsg(server, fmt.Sprintf("%d 秒之后才可以再使用此功能！！！", left))
		return true
	}
	p.GuildRecallTick = now

	p.Engine.mu.Lock()
	recalled, refused := 0, 0
	idx := 0
	for _, pl := range p.Engine.PlayObjectList {
		if pl.GuildName != p.GuildName || pl.Name == p.Name || pl.envir == nil {
			continue
		}
		if pl.envir.Flag.NoGuildRecall {
			refused++
			continue
		}
		x, y, ok := findWalkableNear(p.envir, p.CurrX, p.CurrY, 3)
		if !ok {
			x = p.CurrX + idx%3 - 1
			y = p.CurrY + idx/3 - 1
		}
		teleportPlayerTo(pl, server, p.envir, x, y)
		recalled++
		idx++
	}
	p.Engine.mu.Unlock()
	p.sysMsg(server, fmt.Sprintf("已传送%d个成员，%d个成员未被传送。", recalled, refused))
	return true
}

// cmdRecallHuman GM 召唤玩家到身边（Delphi CmdRecallHuman→RecallHuman，
// ObjBase.pas:7909/26697-26720）。
func (p *PlayObject) cmdRecallHuman(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf("%s 现在不在线，或在其它服务器上！！！", params[0]))
		return true
	}
	dx, dy := dirToOffset(p.Dir)
	x, y, ok := findWalkableNear(p.envir, p.CurrX+dx, p.CurrY+dy, 3)
	if !ok {
		p.sysMsg(server, "召唤失败！！！")
		return true
	}
	teleportPlayerTo(target, server, p.envir, x, y)
	return true
}

// cmdRest 宝宝休息/攻击切换（Delphi CmdChangeSalveStatus，ObjBase.pas:10987-10996）。
func (p *PlayObject) cmdRest(server *netserver.TCPServer, params []string) bool {
	p.toggleSlaveRelax(server)
	return true
}

// cmdSlave Go 专属宠物管理（Delphi 无此命令）。
func (p *PlayObject) cmdSlave(server *netserver.TCPServer, params []string) bool {
	if len(params) >= 1 {
		switch strings.ToLower(params[0]) {
		case "relax":
			p.toggleSlaveRelax(server)
		case "recall":
			p.recallSlaves(server)
		default:
			p.sysMsg(server, "用法: @slave relax|recall")
		}
	} else {
		p.cleanSlaveList()
		p.sysMsg(server, "宠物: "+strconv.Itoa(len(p.SlaveIDs))+"/"+strconv.Itoa(MaxSlaveCount)+" 等级: "+strconv.Itoa(p.SlaveLevel))
	}
	return true
}

func (p *PlayObject) cmdDearRecall(server *netserver.TCPServer, params []string) bool {
	p.DearRecall(server)
	return true
}

func (p *PlayObject) cmdMasterRecall(server *netserver.TCPServer, params []string) bool {
	p.MasterRecall(server)
	return true
}

// 仓库密码：Go 采用一行带参形态（Delphi 为多步聊天状态机，
// ObjBase.pas:7229-7323），命令名对齐 Command.ini 转译表。
func (p *PlayObject) cmdSetStoragePwd(server *netserver.TCPServer, params []string) bool {
	if len(params) >= 1 {
		p.setStoragePassword(server, params[0])
	} else {
		p.sysMsg(server, "用法: @setstoragepwd 密码(4-7位)")
	}
	return true
}

func (p *PlayObject) cmdChgStoragePwd(server *netserver.TCPServer, params []string) bool {
	if len(params) >= 2 {
		p.changeStoragePassword(server, params[0], params[1])
	} else {
		p.sysMsg(server, "用法: @chgstoragepwd 旧密码 新密码")
	}
	return true
}

func (p *PlayObject) cmdUnlockStorage(server *netserver.TCPServer, params []string) bool {
	if len(params) >= 1 {
		p.unlockStorage(server, params[0])
	} else {
		p.sysMsg(server, "用法: @unlockstorage 密码")
	}
	return true
}

// cmdTakeOnHorse/cmdTakeOffHorse Go 坐骑闭环扩展。Delphi 命令名为
// @骑马/@下马（Command.ini 转译保留），且要求佩戴马牌；Go 无马牌物品系统，
// 不检查前置装备（偏差已在审查报告标注）。
func (p *PlayObject) cmdTakeOnHorse(server *netserver.TCPServer, params []string) bool {
	if !p.OnHorse {
		p.OnHorse = true
		p.broadcastFeatureChanged(server)
	}
	p.sysMsg(server, "已上马")
	return true
}

func (p *PlayObject) cmdTakeOffHorse(server *netserver.TCPServer, params []string) bool {
	if p.OnHorse {
		p.OnHorse = false
		p.broadcastFeatureChanged(server)
	}
	p.sysMsg(server, "已下马")
	return true
}

// cmdAttackMode 切换攻击模式（Delphi CmdChangeAttackMode，ObjBase.pas:10679-10708）。
// 用户命令：无参循环切换，有参指定；文案按 Go 攻击模式编号映射 Delphi 原文。
func (p *PlayObject) cmdAttackMode(server *netserver.TCPServer, params []string) bool {
	if len(params) >= 1 {
		if m, err := strconv.Atoi(params[0]); err == nil && m >= AttackModePeace && m <= AttackModePK {
			p.AttackMode = byte(m)
		} else {
			p.cycleAttackMode()
		}
	} else {
		p.cycleAttackMode()
	}
	modeText := [...]string{
		"[攻击模式: 和平攻击]", // AttackModePeace
		"[攻击模式: 编组攻击]", // AttackModeGroup
		"[攻击模式: 行会攻击]", // AttackModeGuild
		"[攻击模式: 全体攻击]", // AttackModeAll
		"[攻击模式: 红名攻击]", // AttackModePK
	}
	p.sysMsgBlue(server, modeText[p.AttackMode])
	return true
}

func (p *PlayObject) cycleAttackMode() {
	p.AttackMode++
	if p.AttackMode > AttackModePK {
		p.AttackMode = AttackModePeace
	}
}

func (p *PlayObject) cmdInfo(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target != nil {
		p.sysMsg(server, fmt.Sprintf("[%s] Lv%d %s HP:%d/%d MP:%d/%d Gold:%d Map:%s(%d,%d) PK:%d",
			target.Name, target.WAbil.Level, jobName(target.Job),
			target.WAbil.HP, target.WAbil.MaxHP, target.WAbil.MP, target.WAbil.MaxMP,
			target.Gold, target.MapName, target.CurrX, target.CurrY, target.PkPoint))
	} else {
		p.sysMsg(server, "玩家不在线: "+params[0])
	}
	return true
}

// cmdMap 当前地图信息（Delphi CmdShowMapInfo，ObjBase.pas:11846-11861）。
func (p *PlayObject) cmdMap(server *netserver.TCPServer, params []string) bool {
	if p.envir == nil {
		return true
	}
	// Delphi 第二参数为 sMapDesc（地图描述），Go 地图格式未解析该字段。
	p.sysMsgGreen(server, fmt.Sprintf("地图名称: %s(%s)", p.envir.Name, p.envir.Name))
	p.sysMsgGreen(server, fmt.Sprintf("地图大小: X(%d) Y(%d)", p.envir.Width, p.envir.Height))
	return true
}

// cmdMapInfo 查询指定地图格子信息（Delphi ShowMapInfo，ObjBase.pas:19262-19290）。
func (p *PlayObject) cmdMapInfo(server *netserver.TCPServer, params []string) bool {
	if len(params) < 3 {
		p.sysMsgGreen(server, "请按正确格式输入: MapInfo 地图号 X Y")
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	x, errX := strconv.Atoi(params[1])
	y, errY := strconv.Atoi(params[2])
	if env == nil || errX != nil || errY != nil || x < 0 || y < 0 || x >= env.Width || y >= env.Height {
		p.sysMsg(server, "Failed to get cell information for map: "+params[0])
		return true
	}
	cell := env.Cells[y*env.Width+x]
	// Delphi 输出 chFlag 与 ObjList.Count；Go 格子无标志字节，Flag 固定 0。
	p.sysMsgGreen(server, "Cell Flag: 0")
	p.sysMsgGreen(server, "Cell Count: "+strconv.Itoa(len(cell.ObjList)))
	return true
}

// cmdMake 制造物品（Delphi CmdMakeItem，ObjBase.pas:13358-13449；数量钳制 1-10）。
func (p *PlayObject) cmdMake(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 || p.ItemDB == nil {
		return true
	}
	count := 1
	if len(params) >= 2 {
		count, _ = strconv.Atoi(params[1])
	}
	if count < 1 {
		count = 1
	}
	if count > 10 {
		count = 10
	}
	def := p.ItemDB.GetByName(params[0])
	if def == nil {
		p.sysMsg(server, fmt.Sprintf("%s 物品名称不正确！！！", params[0]))
		return true
	}
	for i := 0; i < count; i++ {
		p.GiveItem(def.Idx)
	}
	p.SendBagItemsFull(server)
	p.sysMsg(server, "获得: "+def.Name)
	return true
}

// cmdLevel 调整自己等级（Delphi CmdChangeLevel，ObjBase.pas:10848-10877；
// 上限 MAXUPLEVEL=65535，M2Share.pas:91；Delphi 无玩家侧文案，仅控制台）。
func (p *PlayObject) cmdLevel(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	lvl, _ := strconv.Atoi(params[0])
	if lvl < 1 {
		lvl = 1
	}
	if lvl > 65535 {
		lvl = 65535
	}
	old := p.WAbil.Level
	p.WAbil.Level = uint16(lvl)
	p.RecalcAbilitys()
	p.WAbil.HP = p.WAbil.MaxHP
	p.WAbil.MP = p.WAbil.MaxMP
	p.sendHealthSpell(server)
	log.Logf(log.LevelInfo, "GM", "[等级调整] %s (%d -> %d)", p.Name, old, lvl)
	return true
}

// cmdMob 面前刷怪（Delphi CmdMob，ObjBase.pas:13615-13658；数量≤64、等级 0-10）。
func (p *PlayObject) cmdMob(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	count := 1
	if len(params) >= 2 {
		count, _ = strconv.Atoi(params[1])
	}
	if count < 1 {
		count = 1
	}
	if count > 64 {
		count = 64
	}
	level := 0
	if len(params) >= 3 {
		level, _ = strconv.Atoi(params[2])
		if level < 0 || level > 10 {
			level = 0
		}
	}
	dx, dy := dirToOffset(p.Dir)
	now := time.Now().UnixMilli()
	spawned := 0
	for i := 0; i < count; i++ {
		mx := p.CurrX + dx*2 + (i%3 - 1)
		my := p.CurrY + dy*2 + (i/3%3 - 1)
		mon := p.Engine.SpawnMonsterByName(p.MapName, mx, my, params[0], now)
		if mon != nil {
			if level > 0 {
				mon.WAbil.Level = uint16(level)
			}
			spawned++
		}
	}
	if spawned == 0 {
		// 原文如此（Delphi g_sGameCommandMobGenFail，M2Share.pas:3303）。
		p.sysMsg(server, "怪物名称不正确或其它未问题！！！")
	} else {
		p.sysMsg(server, fmt.Sprintf("召唤: %s x%d", params[0], spawned))
	}
	return true
}

// cmdRecallMob 召唤宠物怪（Delphi CmdRecallMob，ObjBase.pas:14040-14085；上限 20）。
func (p *PlayObject) cmdRecallMob(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	count := 1
	if len(params) >= 2 {
		count, _ = strconv.Atoi(params[1])
	}
	if count < 1 {
		count = 1
	}
	if count > 20 {
		count = 20
	}
	level := 0
	if len(params) >= 3 {
		level, _ = strconv.Atoi(params[2])
	}
	dx, dy := dirToOffset(p.Dir)
	now := time.Now().UnixMilli()
	spawned := 0
	for i := 0; i < count; i++ {
		mx := p.CurrX + dx*(2+i%3) + (i / 3)
		my := p.CurrY + dy*(2+i%3) + (i / 3)
		if mon := p.Engine.SpawnMonsterByName(p.MapName, mx, my, params[0], now); mon != nil {
			if level > 0 {
				mon.WAbil.Level = uint16(level)
			}
			mon.PlayerMasterID = p.ID
			p.SlaveIDs = append(p.SlaveIDs, mon.ID)
			spawned++
		}
	}
	p.sysMsg(server, fmt.Sprintf("召唤 %d 个 %s", spawned, params[0]))
	return true
}

// cmdClearMon 清除指定地图怪物（Delphi CmdClearMapMonster，ObjBase.pas:11111-11169）。
func (p *PlayObject) cmdClearMon(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		p.sysMsg(server, "用法: @clearmon 地图号 怪物名称|* [掉物品(0|1)]")
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	if env == nil {
		p.sysMsg(server, "输入的地图不存在！！！")
		return true
	}
	monName := params[1]
	// Delphi 第三参控制 m_boNoItem（是否爆物品）；Go 死亡掉落走统一管线，
	// 暂不支持按命令抑制（偏差已在审查报告标注）。
	now := time.Now().UnixMilli()
	count := 0
	for i := range env.Cells {
		for _, o := range env.Cells[i].ObjList {
			if o.Type != OS_MOVINGOBJECT {
				continue
			}
			mon, ok := o.Obj.(*MonsterObject)
			if !ok || mon.Death {
				continue
			}
			if monName != "*" && mon.Name != monName {
				continue
			}
			mon.Death = true
			mon.WAbil.HP = 0
			mon.DeathTick = now
			env.broadcastDeathMsg(mon.BaseObject, mon.ID, mon.CurrX, mon.CurrY, mon.Dir, true)
			count++
		}
	}
	p.sysMsg(server, fmt.Sprintf("已清除怪物数: %d", count))
	return true
}

func clampGold(p *PlayObject) {
	maxGold := 10000000
	if p.Engine != nil && p.Engine.Config != nil {
		maxGold = p.Engine.Config.GetMaxGold()
	}
	if p.Gold < 0 {
		p.Gold = 0
	}
	if p.Gold > maxGold {
		p.Gold = maxGold
	}
}

// cmdGold Go 专属：给自己加金币（Delphi 无自我金币命令）。
func (p *PlayObject) cmdGold(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	gold, _ := strconv.Atoi(params[0])
	p.Gold += gold
	clampGold(p)
	server.Send(p.Session.ID, protocol.MakeDefaultMsg(protocol.SMGoldChanged, int32(p.Gold), 0, 0, 0), "")
	p.sysMsg(server, "金币: "+strconv.Itoa(p.Gold))
	return true
}

// cmdAddGold/cmdDelGold 对目标加/扣金币（Delphi CmdAddGold/CmdDelGold，
// ObjBase.pas:10436-10482/12296-12342）。Delphi 支持离线记账，Go 仅在线
//（偏差已在审查报告标注）。
func (p *PlayObject) cmdAddGold(server *netserver.TCPServer, params []string) bool {
	return p.adjustTargetGold(server, params, true)
}

func (p *PlayObject) cmdDelGold(server *netserver.TCPServer, params []string) bool {
	return p.adjustTargetGold(server, params, false)
}

func (p *PlayObject) adjustTargetGold(server *netserver.TCPServer, params []string, add bool) bool {
	if len(params) < 2 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, "玩家不在线: "+params[0])
		return true
	}
	gold, _ := strconv.Atoi(params[1])
	if gold < 0 {
		gold = -gold
	}
	if add {
		target.Gold += gold
	} else {
		target.Gold -= gold
	}
	clampGold(target)
	server.Send(target.Session.ID, protocol.MakeDefaultMsg(protocol.SMGoldChanged, int32(target.Gold), 0, 0, 0), "")
	if add {
		p.sysMsg(server, fmt.Sprintf("%s的金币已增加%d.", params[0], gold))
	} else {
		p.sysMsg(server, fmt.Sprintf("%s的金币已减少%d.", params[0], gold))
	}
	return true
}

func (p *PlayObject) cmdHeal(server *netserver.TCPServer, params []string) bool {
	p.WAbil.HP = p.WAbil.MaxHP
	p.WAbil.MP = p.WAbil.MaxMP
	p.sendHealthSpell(server)
	p.sysMsg(server, "已回复")
	return true
}

// cmdPkPoint 查询目标 PK 值（Delphi CmdPKPoint，ObjBase.pas:13945-13968）。
func (p *PlayObject) cmdPkPoint(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf("%s 现在不在线，或在其它服务器上！！！", params[0]))
		return true
	}
	p.sysMsg(server, fmt.Sprintf("%s的PK点数为:%d", target.Name, target.PkPoint))
	return true
}

// cmdIncPkPoint 增减目标 PK 值（Delphi CmdIncPkPoint，ObjBase.pas:13190-13217）。
func (p *PlayObject) cmdIncPkPoint(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, "玩家不在线: "+params[0])
		return true
	}
	pts, _ := strconv.Atoi(params[1])
	target.PkPoint += pts
	if target.PkPoint < 0 {
		target.PkPoint = 0
	}
	target.BroadcastNameColor(server) // Delphi RefNameColor()（ObjBase.pas:13212）
	if pts >= 0 {
		p.sysMsg(server, fmt.Sprintf("%s的PK值已增加%d点...", target.Name, pts))
	} else {
		p.sysMsg(server, fmt.Sprintf("%s的PK值已减少%d点...", target.Name, -pts))
	}
	return true
}

// cmdFreePK 清除目标 PK 值（Delphi CmdFreePenalty，ObjBase.pas:12917-12941）。
func (p *PlayObject) cmdFreePK(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, "玩家不在线: "+params[0])
		return true
	}
	target.PkPoint = 0
	target.BroadcastNameColor(server) // Delphi RefNameColor()（ObjBase.pas:12938）
	target.sysMsg(server, "你的PK值已经被清除...")
	p.sysMsg(server, fmt.Sprintf("%s的PK值已经被清除...", target.Name))
	return true
}

func (p *PlayObject) cmdKick(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target != nil && target.Session != nil {
		server.CloseSession(target.Session.ID)
		p.sysMsg(server, "已踢出: "+params[0])
	} else {
		p.sysMsg(server, "玩家不在线: "+params[0])
	}
	return true
}

// revivePlayer 复活玩家：满血蓝 + 解除死亡 + RM_ALIVE 广播（客户端起身）。
func revivePlayer(target *PlayObject, server *netserver.TCPServer) {
	target.WAbil.HP = target.WAbil.MaxHP
	target.WAbil.MP = target.WAbil.MaxMP
	target.Death = false
	target.sendHealthSpell(server)
	if target.envir != nil {
		target.envir.broadcastRefMsg(target.BaseObject, RM_ALIVE, target.ID, target.CurrX, target.CurrY, target.Dir)
	}
}

// cmdReAlive 复活目标（Delphi CmdReAlive，ObjBase.pas:13998-14025，双发文案为原版行为）。
func (p *PlayObject) cmdReAlive(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, "玩家不在线: "+params[0])
		return true
	}
	revivePlayer(target, server)
	p.sysMsg(server, fmt.Sprintf("%s 已获重生.", target.Name))
	p.sysMsg(server, fmt.Sprintf("%s 已获重生。", target.Name))
	return true
}

// cmdRevive Go 扩展：无参复活自己，有参复活他人。
func (p *PlayObject) cmdRevive(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		revivePlayer(p, server)
		p.sysMsg(server, "已复活")
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target != nil {
		revivePlayer(target, server)
		p.sysMsg(server, "已复活: "+params[0])
	} else {
		p.sysMsg(server, "玩家不在线: "+params[0])
	}
	return true
}

// cmdSuperman 无敌模式（Delphi CmdChangeSuperManMode，ObjBase.pas:10997-11014）。
// Go 以 SuperManMode 模式位实现（伤害入口 applyDamage 豁免），
// 不再用 MaxHP=65535 hack。
func (p *PlayObject) cmdSuperman(server *netserver.TCPServer, params []string) bool {
	p.SuperManMode = !p.SuperManMode
	if p.SuperManMode {
		p.sysMsgGreen(server, "已进入无敌模式.")
	} else {
		p.sysMsgGreen(server, "已退出无敌模式.")
	}
	return true
}

// cmdGameMaster 管理员模式（Delphi CmdChangeAdminMode，ObjBase.pas:10661-10678）。
func (p *PlayObject) cmdGameMaster(server *netserver.TCPServer, params []string) bool {
	p.AdminMode = !p.AdminMode
	if p.AdminMode {
		p.sysMsgGreen(server, "已进入管理员模式.")
	} else {
		p.sysMsgGreen(server, "已退出管理员模式.")
	}
	return true
}

// cmdObserver 隐身模式（Delphi CmdChangeObMode，ObjBase.pas:10918-10935）。
func (p *PlayObject) cmdObserver(server *netserver.TCPServer, params []string) bool {
	p.ObMode = !p.ObMode
	if p.ObMode {
		p.sysMsgGreen(server, "已进入隐身模式.")
	} else {
		p.sysMsgGreen(server, "已退出隐身模式.")
	}
	return true
}

// cmdReGoto 跟踪：传送到目标身边（Delphi CmdReGotoHuman，ObjBase.pas:14125-14150）。
func (p *PlayObject) cmdReGoto(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target != nil && target.envir != nil {
		teleportPlayerTo(p, server, target.envir, target.CurrX, target.CurrY)
		p.sysMsg(server, "已传送到 "+params[0]+" 身边")
	} else {
		p.sysMsg(server, "玩家不在线: "+params[0])
	}
	return true
}

// cmdShutup 禁言（Delphi CmdShutup，ObjBase.pas:14800-14831；默认 5 分钟）。
func (p *PlayObject) cmdShutup(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target != nil {
		minutes := 5
		if len(params) >= 2 {
			if m, _ := strconv.Atoi(params[1]); m > 0 {
				minutes = m
			}
		}
		target.ShutupTick = time.Now().UnixMilli() + int64(minutes)*60000
		p.sysMsg(server, fmt.Sprintf("%s 已被禁言%d分钟", params[0], minutes))
	} else {
		p.sysMsg(server, "玩家不在线: "+params[0])
	}
	return true
}

// cmdReleaseShutup 解除禁言（Delphi CmdReleaseShutup，ObjBase.pas:14863-14898）。
func (p *PlayObject) cmdReleaseShutup(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, "玩家不在线: "+params[0])
		return true
	}
	target.ShutupTick = 0
	target.sysMsg(server, "你已经恢复聊天功能！！！")
	p.sysMsg(server, fmt.Sprintf("%s 已经恢复聊天。", target.Name))
	return true
}

// cmdShutupList 禁言列表（Delphi CmdShutupList，ObjBase.pas:14832-14862）。
// Delphi 为全局名单（离线可查）；Go 按在线玩家禁言倒计时输出。
func (p *PlayObject) cmdShutupList(server *netserver.TCPServer, params []string) bool {
	now := time.Now().UnixMilli()
	p.Engine.mu.Lock()
	var lines []string
	for _, pl := range p.Engine.PlayObjectList {
		if pl.ShutupTick > now {
			left := (pl.ShutupTick-now)/60000 + 1
			lines = append(lines, fmt.Sprintf("%s 剩余%d分钟", pl.Name, left))
		}
	}
	p.Engine.mu.Unlock()
	if len(lines) == 0 {
		p.sysMsg(server, "禁言列表为空！！！")
		return true
	}
	p.sysMsg(server, "禁言列表:")
	for _, l := range lines {
		p.sysMsg(server, l)
	}
	return true
}

// cmdMoveUser Go 专属：单体跨图传送（Delphi @MapMove 为整图迁移）。
func (p *PlayObject) cmdMoveUser(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target != nil {
		newEnvir := p.MapMgr.FindMap(params[1])
		if newEnvir != nil {
			tx, ty := newEnvir.Width/2, newEnvir.Height/2
			if len(params) >= 4 {
				tx, _ = strconv.Atoi(params[2])
				ty, _ = strconv.Atoi(params[3])
			}
			teleportPlayerTo(target, server, newEnvir, tx, ty)
			p.sysMsg(server, "已传送 "+params[0]+" 到 "+params[1])
		} else {
			p.sysMsg(server, "地图不存在: "+params[1])
		}
	} else {
		p.sysMsg(server, "玩家不在线: "+params[0])
	}
	return true
}

// cmdReloadNpc 重载 NPC 脚本（Delphi CmdReloadNpc，ObjBase.pas:14305-14351：
// 无参=附近 9 格，all=全部）。
func (p *PlayObject) cmdReloadNpc(server *netserver.TCPServer, params []string) bool {
	if len(params) >= 1 && strings.EqualFold(params[0], "all") {
		p.Engine.InvalidateAllNpcScripts()
		p.sysMsg(server, "交易NPC重新加载完成！！！")
		p.sysMsg(server, "管理NPC重新加载！！！")
		return true
	}
	count := 0
	for _, npc := range p.Engine.Npcs {
		if npc.envir == p.envir && abs(npc.CurrX-p.CurrX) <= 9 && abs(npc.CurrY-p.CurrY) <= 9 {
			npc.InvalidateScript()
			p.sysMsg(server, npc.Name+"重新加载成功...")
			count++
		}
	}
	if count == 0 {
		p.sysMsg(server, "附近未发现任何交易NPC！！！")
	}
	return true
}

// cmdStatus Go 专属：服务器状态报表（替代 Delphi ViewKernelInfo/ViewOnlineHuman 监控窗体）。
func (p *PlayObject) cmdStatus(server *netserver.TCPServer, params []string) bool {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	p.sysMsg(server, fmt.Sprintf("在线: %d  会话: %d", p.Engine.GetPlayerCount(), server.GetSessionCount()))
	p.sysMsg(server, fmt.Sprintf("怪物: %d  NPC: %d  行会: %d", len(p.Engine.Monsters), len(p.Engine.Npcs), len(p.Engine.Guilds)))
	p.sysMsg(server, fmt.Sprintf("堆内存: %dMB  GC: %d次  Goroutine: %d", ms.HeapAlloc/1024/1024, ms.NumGC, runtime.NumGoroutine()))
	monByMap := make(map[string]int)
	for _, m := range p.Engine.Monsters {
		if !m.Ghost {
			monByMap[m.MapName]++
		}
	}
	if len(monByMap) > 0 {
		names := make([]string, 0, len(monByMap))
		for name := range monByMap {
			names = append(names, name)
		}
		sort.Strings(names)
		var sb strings.Builder
		sb.WriteString("怪物分布: ")
		for i, name := range names {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(&sb, "%s=%d", name, monByMap[name])
		}
		p.sysMsg(server, sb.String())
	}
	return true
}

// cmdNoMob Go 专属停怪开关（Delphi 无对应命令；noarg=切换，on/off 幂等设定
// 供客户端调试指令避免误触）。
func (p *PlayObject) cmdNoMob(server *netserver.TCPServer, params []string) bool {
	mode := ""
	if len(params) >= 1 {
		mode = strings.ToLower(params[0])
	}
	var turnOn bool
	switch mode {
	case "on":
		turnOn = true
	case "off":
		turnOn = false
	default:
		turnOn = !p.Engine.NoMonGen
	}
	p.Engine.NoMonGen = turnOn
	if turnOn {
		now := time.Now().UnixMilli()
		killed := 0
		for _, m := range p.Engine.Monsters {
			if !m.Death && !m.Ghost {
				m.Death = true
				m.DeathTick = now
				if m.envir != nil {
					m.envir.broadcastDeathMsg(m.BaseObject, m.ID, m.CurrX, m.CurrY, m.Dir, true)
				}
				killed++
			}
		}
		for i := range p.Engine.MonGenList {
			p.Engine.MonGenList[i].LiveList = nil
		}
		p.sysMsg(server, fmt.Sprintf("怪物生成已停止，已清除 %d 只", killed))
	} else {
		p.sysMsg(server, "怪物生成已恢复")
	}
	return true
}
