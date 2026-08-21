package main

import (
	"fmt"
	"strings"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// GM 命令——封禁/发言过滤族移植。
// Delphi 参照：ObjBase.pas（分发 ProcessUserLineMsg 8128-8186,8236-8251，
// 实现 CmdDeny*/CmdDelDeny*/CmdShowDeny* 12466-12760、
// CmdDisableSendMsg/CmdDisableSendMsgList/CmdEnableSendMsg 12762-12840、
// CmdViewWhisper 15237-15268）+ M2Share.pas g_GameCommand 命令表
//（2838-2840,2903-2912）与文案常量（3187-3189,3287,3341-3343）。
//
// 登录拦截消费端见 main.go（CheckDenyLogon，usrengine.go）；
// 禁言/侦听私聊消费端见 chatsystem.go（HandleSay/HandleWhisper）。
//
// 偏差说明：
//   - Delphi 禁言名单（g_DisableSendMsgList）也落盘
//     Envir/DisableSendMsgList.txt（M2Share.pas:4491/4662）；本任务范围
//     仅三张封禁表持久化（deny_lists.json），禁言名单只在内存维护。
//   - Delphi @ShowDenyIPLogon 空表文案误用"禁止登录角色列表为空。"
//     （ObjBase.pas:12749 复用了角色列表常量），此处照抄该怪癖。

func init() {
	// --- 封禁九条（分发顺序 ObjBase.pas:8128-8186）---
	registerGMCommand(&gmCmdDef{key: "DenyIPaddrLogon", name: "DenyIPLogon", permMin: 10, usage: denyUsageIP, handler: (*PlayObject).cmdDenyIPLogon})
	registerGMCommand(&gmCmdDef{key: "DenyAccountLogon", name: "DenyAccountLogon", permMin: 10, usage: denyUsageAccount, handler: (*PlayObject).cmdDenyAccountLogon})
	registerGMCommand(&gmCmdDef{key: "DenyCharNameLogon", name: "DenyCharNameLogon", permMin: 10, usage: denyUsageChar, handler: (*PlayObject).cmdDenyCharNameLogon})
	registerGMCommand(&gmCmdDef{key: "DelDenyIPLogon", name: "DelDenyIPLogon", permMin: 10, usage: delDenyUsageIP, handler: (*PlayObject).cmdDelDenyIPLogon})
	registerGMCommand(&gmCmdDef{key: "DelDenyAccountLogon", name: "DelDenyAccountLogon", permMin: 10, usage: delDenyUsageAccount, handler: (*PlayObject).cmdDelDenyAccountLogon})
	registerGMCommand(&gmCmdDef{key: "DelDenyCharNameLogon", name: "DelDenyCharNameLogon", permMin: 10, usage: delDenyUsageChar, handler: (*PlayObject).cmdDelDenyCharNameLogon})
	registerGMCommand(&gmCmdDef{key: "ShowDenyIPLogon", name: "ShowDenyIPLogon", permMin: 10, handler: (*PlayObject).cmdShowDenyIPLogon})
	// Delphi 硬编码 m_btPermission<6 静默退出（ObjBase.pas:12692），
	// 不走 nPermissionMin=10 的常规"权限不够"提示。
	registerGMCommand(&gmCmdDef{key: "ShowDenyAccountLogon", name: "ShowDenyAccountLogon", permMin: 6, permFail: permFailSilent, handler: (*PlayObject).cmdShowDenyAccountLogon})
	registerGMCommand(&gmCmdDef{key: "ShowDenyCharNameLogon", name: "ShowDenyCharNameLogon", permMin: 10, handler: (*PlayObject).cmdShowDenyCharNameLogon})
	// --- 发言过滤三条（分发顺序 ObjBase.pas:8236-8251）---
	registerGMCommand(&gmCmdDef{key: "DisableSendMsg", name: "DisableSendMsg", permMin: 10, usage: muteUsage, handler: (*PlayObject).cmdDisableSendMsg})
	registerGMCommand(&gmCmdDef{key: "EnableSendMsg", name: "EnableSendMsg", permMin: 10, usage: muteUsage, handler: (*PlayObject).cmdEnableSendMsg})
	registerGMCommand(&gmCmdDef{key: "DisableSendMsgList", name: "DisableSendMsgList", permMin: 10, handler: (*PlayObject).cmdDisableSendMsgList})
	// --- 侦听私聊（分发 ObjBase.pas:8013）---
	registerGMCommand(&gmCmdDef{key: "ViewWhisper", name: "ViewWhisper", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdViewWhisper})
}

// 参数帮助文案（命令格式: @<命令> <usage>，g_sGameCommandParamUnKnow，
// M2Share.pas:3251；各 usage 为 Delphi 命令内原文）。
const (
	denyUsageIP         = "IP地址 是否永久封(0,1)" // ObjBase.pas:12533
	denyUsageAccount    = "登录帐号 是否永久封(0,1)" // ObjBase.pas:12475
	denyUsageChar       = "人物名称 是否永久封(0,1)" // ObjBase.pas:12504
	delDenyUsageIP      = "IP地址"            // ObjBase.pas:12662
	delDenyUsageAccount = "登录帐号"            // ObjBase.pas:12586
	delDenyUsageChar    = "人物名称"            // ObjBase.pas:12624
	muteUsage           = "人物名称"            // ObjBase.pas:12773/12820
)

// 封禁成功/删除文案（均为 c_Green；值为 Delphi 原文，前缀拼接命令参数）。
const (
	msgDenyIPAdded       = "已加入禁止登录IP列表"    // ObjBase.pas:12542
	msgDenyIPAddedTemp   = "已加入临时禁止登录IP列表"  // ObjBase.pas:12546
	msgDenyAccountAdded  = "已加入禁止登录帐号列表"    // ObjBase.pas:12484
	msgDenyAccountAddedT = "已加入临时禁止登录帐号列表"  // ObjBase.pas:12488
	msgDenyCharAdded     = "已加入禁止人物列表"      // ObjBase.pas:12513
	msgDenyCharAddedTemp = "已加入临时禁止人物列表"    // ObjBase.pas:12517
	msgDenyIPRemoved     = "已从禁止登录IP列表中删除。" // ObjBase.pas:12675
	msgDenyAccountRemove = "已从禁止登录帐号列表中删除。" // ObjBase.pas:12599
	msgDenyCharRemoved   = "已从禁止登录人物列表中删除。" // ObjBase.pas:12637
	msgDenyNotBanned     = "没有被禁止登录。"       // ObjBase.pas:12608/12646/12684
)

// 封禁列表为空文案（c_Green）。
const (
	msgDenyAccountEmpty = "禁止登录帐号列表为空。" // ObjBase.pas:12697
	msgDenyCharEmpty    = "禁止登录角色列表为空。" // ObjBase.pas:12723
	// Delphi @ShowDenyIPLogon 复用了角色列表常量（ObjBase.pas:12749），照抄。
	msgDenyIPEmpty = "禁止登录角色列表为空。"
)

// 禁言文案（ObjBase.pas:12783/12797/12801/12834/12838）。
const (
	msgMuteAdded     = " 已加入禁言列表。"   // c_Green，12783
	msgMuteRemoved   = " 已从禁言列表中删除。" // c_Green，12834
	msgMuteNotBanned = " 没有被禁言！！！"   // c_Red，12838
	msgMuteListEmpty = "禁言列表为空！！！"   // c_Red，12797
	msgMuteListHead  = "禁言列表:"       // c_Blue，12801
)

// 侦听私聊文案（M2Share.pas:3341-3343，c_Green；未在线走 g_sNowNotOnLine）。
const (
	msgWhisperListenStop = "已停止侦听%s的私聊信息..." // 3342
	msgWhisperListenOn   = "正在侦听%s的私聊信息..."  // 3343
)

// denyListSel 选择 UserEngine 上的一张封禁表。
type denyListSel func(e *UserEngine) *[]DenyEntry

func selDenyIP(e *UserEngine) *[]DenyEntry      { return &e.DenyIPList }
func selDenyAccount(e *UserEngine) *[]DenyEntry { return &e.DenyAccountList }
func selDenyChar(e *UserEngine) *[]DenyEntry    { return &e.DenyCharList }

// cmdDenyLogon 封禁通用实现（Delphi CmdDenyIPaddrLogon/CmdDenyAccountLogon/
// CmdDenyCharNameLogon，ObjBase.pas:12524/12466/12495）。
// 第二参以 '1' 开头为永久封禁（落盘），否则临时（仅内存）。
func (p *PlayObject) cmdDenyLogon(server *netserver.TCPServer, params []string,
	key, usage string, sel denyListSel, msgPerm, msgTemp string) bool {
	if p.Engine == nil {
		return true
	}
	if len(params) < 1 {
		p.gmParamHelp(server, key, usage)
		return true
	}
	value := params[0]
	permanent := len(params) >= 2 && len(params[1]) > 0 && params[1][0] == '1'
	e := p.Engine
	e.mu.Lock()
	*sel(e) = append(*sel(e), DenyEntry{Value: value, Permanent: permanent})
	e.mu.Unlock()
	if permanent {
		e.SaveDenyLists()
		p.sysMsgGreen(server, value+msgPerm)
	} else {
		p.sysMsgGreen(server, value+msgTemp)
	}
	return true
}

// cmdDelDenyLogon 解除封禁通用实现（Delphi CmdDelDenyIPaddrLogon/
// CmdDelDenyAccountLogon/CmdDelDenyCharNameLogon，ObjBase.pas:12649/12573/
// 12611）。CompareText 不敏感匹配第一条；删到永久条目时回写。
func (p *PlayObject) cmdDelDenyLogon(server *netserver.TCPServer, params []string,
	key, usage string, sel denyListSel, msgRemoved string) bool {
	if p.Engine == nil {
		return true
	}
	if len(params) < 1 {
		p.gmParamHelp(server, key, usage)
		return true
	}
	value := params[0]
	e := p.Engine
	removed, wasPerm := false, false
	e.mu.Lock()
	list := sel(e)
	for i := range *list {
		if strings.EqualFold((*list)[i].Value, value) {
			wasPerm = (*list)[i].Permanent
			*list = append((*list)[:i], (*list)[i+1:]...)
			removed = true
			break
		}
	}
	e.mu.Unlock()
	if removed {
		if wasPerm {
			e.SaveDenyLists()
		}
		p.sysMsgGreen(server, value+msgRemoved)
	} else {
		p.sysMsgGreen(server, value+msgDenyNotBanned)
	}
	return true
}

// cmdShowDenyLogon 列表打印通用实现（Delphi CmdShowDenyIPaddrLogon/
// CmdShowDenyAccountLogon/CmdShowDenyCharNameLogon，ObjBase.pas:12735/12687/
// 12709）。空表输出 msgEmpty，逐行输出各条目值（c_Green）。
func (p *PlayObject) cmdShowDenyLogon(server *netserver.TCPServer,
	sel denyListSel, msgEmpty string) bool {
	if p.Engine == nil {
		return true
	}
	e := p.Engine
	e.mu.RLock()
	list := append([]DenyEntry(nil), (*sel(e))...)
	e.mu.RUnlock()
	if len(list) == 0 {
		p.sysMsgGreen(server, msgEmpty)
		return true
	}
	for _, d := range list {
		p.sysMsgGreen(server, d.Value)
	}
	return true
}

// 封禁九条命令入口（参数帮助/权限由注册表统一处理）。

func (p *PlayObject) cmdDenyIPLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdDenyLogon(server, params, "DenyIPaddrLogon", denyUsageIP, selDenyIP, msgDenyIPAdded, msgDenyIPAddedTemp)
}
func (p *PlayObject) cmdDenyAccountLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdDenyLogon(server, params, "DenyAccountLogon", denyUsageAccount, selDenyAccount, msgDenyAccountAdded, msgDenyAccountAddedT)
}
func (p *PlayObject) cmdDenyCharNameLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdDenyLogon(server, params, "DenyCharNameLogon", denyUsageChar, selDenyChar, msgDenyCharAdded, msgDenyCharAddedTemp)
}
func (p *PlayObject) cmdDelDenyIPLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdDelDenyLogon(server, params, "DelDenyIPLogon", delDenyUsageIP, selDenyIP, msgDenyIPRemoved)
}
func (p *PlayObject) cmdDelDenyAccountLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdDelDenyLogon(server, params, "DelDenyAccountLogon", delDenyUsageAccount, selDenyAccount, msgDenyAccountRemove)
}
func (p *PlayObject) cmdDelDenyCharNameLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdDelDenyLogon(server, params, "DelDenyCharNameLogon", delDenyUsageChar, selDenyChar, msgDenyCharRemoved)
}
func (p *PlayObject) cmdShowDenyIPLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdShowDenyLogon(server, selDenyIP, msgDenyIPEmpty)
}
func (p *PlayObject) cmdShowDenyAccountLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdShowDenyLogon(server, selDenyAccount, msgDenyAccountEmpty)
}
func (p *PlayObject) cmdShowDenyCharNameLogon(server *netserver.TCPServer, params []string) bool {
	return p.cmdShowDenyLogon(server, selDenyChar, msgDenyCharEmpty)
}

// cmdDisableSendMsg 禁言（Delphi CmdDisableSendMsg，ObjBase.pas:12762-12785）：
// 在线目标设 m_boFilterSendMsg=true，并把名字加入 g_DisableSendMsgList
// （无论是否在线都入表，供登录时恢复）。
// 偏差：Delphi 落盘 DisableSendMsgList.txt，Go 仅内存（见文件头说明）。
func (p *PlayObject) cmdDisableSendMsg(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	if len(params) < 1 {
		p.gmParamHelp(server, "DisableSendMsg", muteUsage)
		return true
	}
	name := params[0]
	e := p.Engine
	if target := e.GetPlayerByName(name); target != nil {
		target.FilterSendMsg = true
	}
	e.mu.Lock()
	e.DisableSendMsg = append(e.DisableSendMsg, name)
	e.mu.Unlock()
	p.sysMsgGreen(server, name+msgMuteAdded)
	return true
}

// cmdEnableSendMsg 解除禁言（Delphi CmdEnableSendMsg，ObjBase.pas:12808-12840）：
// CompareText 不敏感匹配，命中则清除在线目标 FilterSendMsg 并移出名单；
// 未命中文案 c_Red。
func (p *PlayObject) cmdEnableSendMsg(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	if len(params) < 1 {
		p.gmParamHelp(server, "EnableSendMsg", muteUsage)
		return true
	}
	name := params[0]
	e := p.Engine
	removed := false
	e.mu.Lock()
	for i := range e.DisableSendMsg {
		if strings.EqualFold(e.DisableSendMsg[i], name) {
			e.DisableSendMsg = append(e.DisableSendMsg[:i], e.DisableSendMsg[i+1:]...)
			removed = true
			break
		}
	}
	e.mu.Unlock()
	if !removed {
		p.sysMsg(server, name+msgMuteNotBanned)
		return true
	}
	if target := e.GetPlayerByName(name); target != nil {
		target.FilterSendMsg = false
	}
	p.sysMsgGreen(server, name+msgMuteRemoved)
	return true
}

// cmdDisableSendMsgList 禁言列表（Delphi CmdDisableSendMsgList，
// ObjBase.pas:12786-12807）：空表 c_Red 提示，否则标题 c_Blue+逐行 c_Green。
func (p *PlayObject) cmdDisableSendMsgList(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	e := p.Engine
	e.mu.RLock()
	names := append([]string(nil), e.DisableSendMsg...)
	e.mu.RUnlock()
	if len(names) == 0 {
		p.sysMsg(server, msgMuteListEmpty)
		return true
	}
	p.sysMsgBlue(server, msgMuteListHead)
	for _, n := range names {
		p.sysMsgGreen(server, n)
	}
	return true
}

// cmdViewWhisper 侦听私聊（Delphi CmdViewWhisper，ObjBase.pas:15237-15268）。
// Delphi 把侦听者记在目标侧（target.m_GetWhisperHuman := Self）；Go 反向
// 记在 GM 侧（p.ViewWhisperName = 目标名），消费端 chatsystem.go 据此转发。
// 切换语义：再输同一目标名取消；空参在侦听中=取消、否则出参数帮助。
// 目标须在线（Delphi GetPlayObject 判空 → g_sNowNotOnLine）。
func (p *PlayObject) cmdViewWhisper(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	if len(params) < 1 || params[0] == "" {
		if p.ViewWhisperName != "" {
			old := p.ViewWhisperName
			p.ViewWhisperName = ""
			p.sysMsgGreen(server, fmt.Sprintf(msgWhisperListenStop, old))
		} else {
			p.gmParamHelp(server, "ViewWhisper", "人物名称")
		}
		return true
	}
	name := params[0]
	if strings.EqualFold(p.ViewWhisperName, name) {
		p.ViewWhisperName = ""
		p.sysMsgGreen(server, fmt.Sprintf(msgWhisperListenStop, name))
		return true
	}
	if p.Engine.GetPlayerByName(name) == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, name))
		return true
	}
	p.ViewWhisperName = name
	p.sysMsgGreen(server, fmt.Sprintf(msgWhisperListenOn, name))
	return true
}

// forwardWhisperWatchers 将私聊转发给在线侦听者（Delphi m_GetWhisperHuman
// 转发，ObjBase.pas:2190-2202：发起者与接收者两侧的侦听者都收到副本）。
// Go 扩展格式：'[侦听] A→B: 内容'（Delphi 原文为 'A=>B 内容'，走
// RM_WHISPER 通道；Go 以系统消息呈现，格式自拟）。
func (p *PlayObject) forwardWhisperWatchers(server *netserver.TCPServer, target *PlayObject, text string) {
	if p.Engine == nil || target == nil {
		return
	}
	e := p.Engine
	e.mu.RLock()
	var watchers []*PlayObject
	for _, pl := range e.PlayObjectList {
		if pl == p || pl == target || pl.Ghost || pl.ViewWhisperName == "" {
			continue
		}
		if strings.EqualFold(pl.ViewWhisperName, p.Name) || strings.EqualFold(pl.ViewWhisperName, target.Name) {
			watchers = append(watchers, pl)
		}
	}
	e.mu.RUnlock()
	for _, w := range watchers {
		w.sysMsg(server, fmt.Sprintf("[侦听] %s→%s: %s", p.Name, target.Name, text))
	}
}
