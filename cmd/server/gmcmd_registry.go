package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// GM 命令注册表：Delphi ProcessUserLineMsg 顺序 if 链
// (ObjBase.pas:7188-8600) + g_GameCommand 命令表 (M2Share.pas:2724-2926)
// 的表驱动移植。命令字面量与权限等级可被 serverconfig/server.jsonc
// commands.names / commands.permissions（Command.ini 转译）覆盖。

// gmPermFailMode 权限不足时的反馈方式（Delphi 三套并存）：
// permFailLow      — "权限不够..."（g_sGameCommandPermissionTooLow，M2Share.pas:3250）
// permFailSilent   — 静默退出（Delphi 硬编码 m_btPermission<6 Exit 的命令）
// permFailFallback — 与未知命令同文案（调试命令段门槛，ObjBase.pas:8254）
type gmPermFailMode int

const (
	permFailLow gmPermFailMode = iota
	permFailSilent
	permFailFallback
)

type gmCmdHandler func(p *PlayObject, server *netserver.TCPServer, params []string) bool

type gmCmdDef struct {
	key      string // Delphi Command.ini 键名（server.jsonc names/permissions 查找键）
	name     string // Delphi 默认字面量
	permMin  int    // Delphi 默认 nPermissionMin
	permFail gmPermFailMode
	usage    string // "?" 帮助显示的参数格式
	handler  gmCmdHandler

	resolvedName string
	resolvedPerm int
}

var (
	gmCmdDefs  []*gmCmdDef
	gmCmdTable map[string]*gmCmdDef
)

func registerGMCommand(def *gmCmdDef) {
	gmCmdDefs = append(gmCmdDefs, def)
}

func (p *PlayObject) sysMsg(server *netserver.TCPServer, text string) {
	msg := protocol.MakeDefaultMsg(protocol.SMSysMessage, 0, chatColorSys, 0, 0)
	server.Send(p.Session.ID, msg, protocol.EncodeString(text))
}

// buildGMCommandTable 用配置解析每条命令的字面量与权限，生成小写化分发表。
// 同名字面量先注册者优先（Delphi 顺序匹配语义，如 USERMOVE 遮蔽 MOVE）。
func buildGMCommandTable(cfg *ServerConfig) {
	gmCmdTable = make(map[string]*gmCmdDef, len(gmCmdDefs))
	for _, def := range gmCmdDefs {
		def.resolvedName = def.name
		def.resolvedPerm = def.permMin
		if cfg != nil {
			if v, ok := cfg.Commands.Names[def.key]; ok && v != "" {
				def.resolvedName = v
			}
			if v, ok := cfg.Commands.Permissions[def.key]; ok {
				def.resolvedPerm = v
			}
		}
		lk := strings.ToLower(def.resolvedName)
		if _, exists := gmCmdTable[lk]; !exists {
			gmCmdTable[lk] = def
		}
	}
}

// splitGMParams 复刻 GetValidStr3（HUtil32.pas:695）：按空格/冒号/逗号/TAB
// 拆分，连续分隔符合并（不产生空段），上限 8 段（命令 + 7 个参数）。
func splitGMParams(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == ':' || r == ',' || r == '\t'
	})
	if len(fields) > 8 {
		fields = fields[:8]
	}
	return fields
}

const (
	gmMsgPermTooLow  = "权限不够..."                             // g_sGameCommandPermissionTooLow（M2Share.pas:3250）
	gmMsgParamUnKnow = "命令格式: @%s %s"                          // g_sGameCommandParamUnKnow（M2Share.pas:3251）
	gmMsgBadCmd      = "@%s 此命令不正确，或没有足够的权限！！！" // 兜底（ObjBase.pas:8591）
)

// HandleGMCommand 是聊天行 "@..." 的主分发器（Delphi ProcessUserLineMsg）。
func (p *PlayObject) HandleGMCommand(cmd string, server *netserver.TCPServer) bool {
	if len(cmd) > 1 && cmd[0] == '!' {
		return p.gmRedBroadcast(cmd[1:], server)
	}
	parts := splitGMParams(cmd)
	if len(parts) == 0 {
		return false
	}
	command := strings.ToLower(parts[0])
	params := parts[1:]

	def := gmCmdTable[command]
	if def == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgBadCmd, parts[0]))
		return true
	}
	if len(params) > 0 && strings.HasPrefix(params[0], "?") && def.usage != "" {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, def.resolvedName, def.usage))
		return true
	}
	if int(p.Permission) < def.resolvedPerm {
		switch def.permFail {
		case permFailLow:
			p.sysMsg(server, gmMsgPermTooLow)
		case permFailSilent:
		case permFailFallback:
			p.sysMsg(server, fmt.Sprintf(gmMsgBadCmd, parts[0]))
		}
		return true
	}
	return def.handler(p, server, params)
}

// gmRedBroadcast GM 红字全服广播 "@!内容"（ObjBase.pas:7785-7798）。
// Delphi：权限≥6（g_GMRedMsgCmd 默认 '!'，M2Share.pas:1495），2 秒冷却，
// 截断到 nSayRedMsgMaxLen，UserEngine.SendBroadCastMsg(SC, t_GM)。
func (p *PlayObject) gmRedBroadcast(text string, server *netserver.TCPServer) bool {
	perm := 6
	if p.Engine != nil && p.Engine.Config != nil {
		if v, ok := p.Engine.Config.Commands.Permissions["GMRedMsgCmd"]; ok {
			perm = v
		}
	}
	if int(p.Permission) < perm {
		p.sysMsg(server, fmt.Sprintf(gmMsgBadCmd, "!"+text))
		return true
	}
	if p.Engine == nil {
		return true
	}
	now := time.Now().UnixMilli()
	if now-p.LastRedMsgTick < 2000 {
		return true
	}
	p.LastRedMsgTick = now
	if r := []rune(text); len(r) > 255 { // nSayRedMsgMaxLen 默认 255（M2Share.pas:1774）
		text = string(r[:255])
	}
	p.Engine.mu.Lock()
	targets := make([]*PlayObject, 0, len(p.Engine.PlayObjectList))
	for _, pl := range p.Engine.PlayObjectList {
		if pl.Session != nil {
			targets = append(targets, pl)
		}
	}
	p.Engine.mu.Unlock()
	// sGMRedMsgpreFix 默认 '〖ＧＭ〗'（M2Share.pas:2148）；
	// boShutRedMsgShowGMName 默认 False 不带角色名（M2Share.pas:1769）。
	msg := "〖ＧＭ〗" + text
	for _, pl := range targets {
		pl.sysMsg(server, msg)
	}
	return true
}
