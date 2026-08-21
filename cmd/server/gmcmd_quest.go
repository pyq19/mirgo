package main

import (
	"fmt"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// GM 命令——任务标志族移植。
// Delphi 参照：ObjBase.pas（分发 7919-7932/8116-8124 常规段、
// 8290-8351 调试段内联）+ M2Share.pas g_GameCommand（2788-2790,
// 2823-2826,2899）与文案常量（3295-3301）。
// 标志载体 = PlayObject.QuestFlag/QuestUnitOpen/QuestUnit（各 1024 位，
// playobject.go:192-195，位操作 questsystem.go），与 Delphi
// m_QuestFlag/m_QuestUnitOpen/m_QuestUnit 对应（ObjBase.pas:110-112）。
// 注册顺序 = Delphi 分发链出现顺序。

func init() {
	// --- 常规 GM 段（ObjBase.pas:7919-7932,8116-8124，nPermissionMin 检查）---
	registerGMCommand(&gmCmdDef{key: "ShowFlag", name: "showflag", permMin: 10, usage: "人物名称 标识号", handler: (*PlayObject).cmdShowFlag})
	registerGMCommand(&gmCmdDef{key: "ShowOpen", name: "showopen", permMin: 10, usage: "人物名称 单元号", handler: (*PlayObject).cmdShowOpen})
	registerGMCommand(&gmCmdDef{key: "ShowUnit", name: "showunit", permMin: 10, usage: "人物名称 单元号", handler: (*PlayObject).cmdShowUnit})
	registerGMCommand(&gmCmdDef{key: "ClearMission", name: "ClearMission", permMin: 10, handler: (*PlayObject).cmdClearMission})
	registerGMCommand(&gmCmdDef{key: "StartQuest", name: "StartQuest", permMin: 10, handler: (*PlayObject).cmdStartQuest})
	// --- 调试段内联（ObjBase.pas:8290-8351）：分发仅处 m_btPermission>4
	// 门槛内，不检查 nPermissionMin（表值 10 实际不生效），权限≤4 落
	// 8591 兜底 → permMin 5 + permFailFallback。
	registerGMCommand(&gmCmdDef{key: "SetFlag", name: "setflag", permMin: 5, permFail: permFailFallback, handler: (*PlayObject).cmdSetFlag})
	registerGMCommand(&gmCmdDef{key: "SetOpen", name: "setopen", permMin: 5, permFail: permFailFallback, handler: (*PlayObject).cmdSetOpen})
	registerGMCommand(&gmCmdDef{key: "SetUnit", name: "setunit", permMin: 5, permFail: permFailFallback, handler: (*PlayObject).cmdSetUnit})
}

// cmdShowFlag 查询目标任务标志（Delphi CmdShowHumanFlag，
// ObjBase.pas:11749-11780；g_sGameCommandShowHumanFlagONMsg/OFFMsg，
// M2Share.pas:3296-3297，c_Green）。
func (p *PlayObject) cmdShowFlag(server *netserver.TCPServer, params []string) bool {
	return p.showQuestBit(server, params, "ShowFlag", "人物名称 标识号", func(t *PlayObject, idx int) bool {
		return questGetBit(&t.QuestFlag, idx)
	})
}

// cmdShowOpen 查询目标任务单元开启状态（Delphi CmdShowHumanUnitOpen，
// ObjBase.pas:11814-11845；帮助串复用单元文案
// g_sGameCommandShowHumanUnitHelpMsg，M2Share.pas:3299）。
func (p *PlayObject) cmdShowOpen(server *netserver.TCPServer, params []string) bool {
	return p.showQuestBit(server, params, "ShowOpen", "人物名称 单元号", func(t *PlayObject, idx int) bool {
		return questGetBit(&t.QuestUnitOpen, idx)
	})
}

// cmdShowUnit 查询目标任务单元状态（Delphi CmdShowHumanUnit，
// ObjBase.pas:11782-11813）。
func (p *PlayObject) cmdShowUnit(server *netserver.TCPServer, params []string) bool {
	return p.showQuestBit(server, params, "ShowUnit", "人物名称 单元号", func(t *PlayObject, idx int) bool {
		return questGetBit(&t.QuestUnit, idx)
	})
}

func (p *PlayObject) showQuestBit(server *netserver.TCPServer, params []string,
	key, usage string, get func(*PlayObject, int) bool) bool {
	if len(params) < 2 {
		p.gmParamHelp(server, key, usage)
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	idx := atoiOr(params[1], 0) // Delphi Str_ToInt(sFlag, 0)
	on := "OFF"
	if get(target, idx) {
		on = "ON"
	}
	p.sysMsgGreen(server, fmt.Sprintf("%s: [%d] = %s", target.Name, idx, on))
	return true
}

// cmdSetQuestBit SetFlag/SetOpen/SetUnit 共用实现（Delphi 内联分发，
// ObjBase.pas:8290-8309/8311-8330/8332-8351）：目标不存在时输出
// '@<命令字面量> 人物名称 标志号 数字(0 - 1)'（原文无"命令格式:"前缀）。
func (p *PlayObject) cmdSetQuestBit(server *netserver.TCPServer, params []string,
	key string, arr func(*PlayObject) *[128]byte) bool {
	if len(params) < 1 {
		p.sysMsg(server, "@"+gmNameByKey(key, key)+" 人物名称 标志号 数字(0 - 1)")
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, "@"+gmNameByKey(key, key)+" 人物名称 标志号 数字(0 - 1)")
		return true
	}
	idx, val := 0, 0 // Delphi Str_ToInt(sParam2/sParam3, 0)；缺参按 0
	if len(params) >= 2 {
		idx = atoiOr(params[1], 0)
	}
	if len(params) >= 3 {
		val = atoiOr(params[2], 0)
	}
	if val != 0 {
		questSetBit(arr(target), idx)
	} else {
		questClearBit(arr(target), idx)
	}
	on := "OFF"
	if questGetBit(arr(target), idx) {
		on = "ON"
	}
	p.sysMsgGreen(server, fmt.Sprintf("%s: [%d] = %s", target.Name, idx, on))
	return true
}

func (p *PlayObject) cmdSetFlag(server *netserver.TCPServer, params []string) bool {
	return p.cmdSetQuestBit(server, params, "SetFlag", func(t *PlayObject) *[128]byte { return &t.QuestFlag })
}

func (p *PlayObject) cmdSetOpen(server *netserver.TCPServer, params []string) bool {
	return p.cmdSetQuestBit(server, params, "SetOpen", func(t *PlayObject) *[128]byte { return &t.QuestUnitOpen })
}

func (p *PlayObject) cmdSetUnit(server *netserver.TCPServer, params []string) bool {
	return p.cmdSetQuestBit(server, params, "SetUnit", func(t *PlayObject) *[128]byte { return &t.QuestUnit })
}

// cmdClearMission 清除目标任务标志（Delphi CmdClearMission，
// ObjBase.pas:11170-11197）。
// 原版怪癖照抄：无参帮助串尾部多余 ')'（ObjBase.pas:11178）；
// 离线文案为硬编码变体（11190，未用公共常量）。
// Delphi 只清 m_QuestFlag 不清 QuestUnit/QuestUnitOpen（11193 FillChar）。
func (p *PlayObject) cmdClearMission(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("ClearMission", "ClearMission")+" 人物名称)")
		return true
	}
	if params[0] != "" && params[0][0] == '?' {
		// '?' 帮助为命令内自定义文案（ObjBase.pas:11181-11184，c_Blue），
		// 故注册表 usage 留空、由此处分发。
		p.sysMsgBlue(server, "此命令用于清除人物的任务标志。")
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf("%s不在线，或在其它服务器上！！", params[0]))
		return true
	}
	target.QuestFlag = [128]byte{}
	p.sysMsgGreen(server, fmt.Sprintf("%s的任务标志已经全部清零。", params[0]))
	return true
}

// cmdStartQuest 全员触发任务脚本标签（Delphi CmdStartQuest，
// ObjBase.pas:14991-15006 → UserEngine.SendQuestMsg，UsrEngn.pas:3118-3129：
// 遍历全部在线未死亡玩家，逐一 g_ManageNPC.GotoLable(Self, sQuestName, False)）。
// 偏差：Go 无 QManage 单例（同 cmdMemberFunc/gotoMemberLabel 的就近扫描
// 策略），改为每个玩家各自触发其同地图最近的含该标签 NPC；无脚本命中
// 时静默（与 Delphi GotoLable 标签缺失行为一致）。无任何反馈消息。
func (p *PlayObject) cmdStartQuest(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("StartQuest", "StartQuest")+" 问答名称")
		return true
	}
	if p.Engine == nil {
		return true
	}
	p.Engine.mu.Lock()
	targets := make([]*PlayObject, 0, len(p.Engine.PlayObjectList))
	for _, pl := range p.Engine.PlayObjectList {
		if !pl.Death && !pl.Ghost {
			targets = append(targets, pl)
		}
	}
	p.Engine.mu.Unlock()
	for _, pl := range targets {
		pl.gotoNpcLabel(server, params[0])
	}
	return true
}
