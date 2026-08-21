package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// 用户命令族移植（Delphi ProcessUserLineMsg 用户段，ObjBase.pas:7362-7780）。
// 全部为 permMin=0 的普通玩家命令（ClearBag/ShowUseItemInfo/BindUseItem
// 为 10，M2Share.pas:2920-2922），权限不足统一走 permFailLow
//（g_sGameCommandPermissionTooLow，注册表 gmcmd_registry.go 处理）。

func init() {
	// --- 会员功能（ObjBase.pas:7551-7558，分发在 Date 之前）---
	registerGMCommand(&gmCmdDef{key: "MemberFunc", name: "MemberFunc", permMin: 0, handler: (*PlayObject).cmdMemberFunc})
	registerGMCommand(&gmCmdDef{key: "MemberFuncEx", name: "MemberFuncEx", permMin: 0, handler: (*PlayObject).cmdMemberFuncEx})
	// --- 用户段（ObjBase.pas:7613-7693）---
	registerGMCommand(&gmCmdDef{key: "Date", name: "Date", permMin: 0, handler: (*PlayObject).cmdDate})
	registerGMCommand(&gmCmdDef{key: "PrvMsg", name: "PrvMsg", permMin: 0, usage: "人物名称", handler: (*PlayObject).cmdPrvMsg})
	registerGMCommand(&gmCmdDef{key: "AllowMsg", name: "AllowMsg", permMin: 0, handler: (*PlayObject).cmdAllowMsg})
	registerGMCommand(&gmCmdDef{key: "LetShout", name: "LetShout", permMin: 0, handler: (*PlayObject).cmdLetShout})
	registerGMCommand(&gmCmdDef{key: "LetTrade", name: "LetTrade", permMin: 0, handler: (*PlayObject).cmdLetTrade})
	registerGMCommand(&gmCmdDef{key: "BanGuildChat", name: "BanGuildChat", permMin: 0, handler: (*PlayObject).cmdBanGuildChat})
	registerGMCommand(&gmCmdDef{key: "LetGuild", name: "LetGuild", permMin: 0, handler: (*PlayObject).cmdLetGuild})
	registerGMCommand(&gmCmdDef{key: "EndGuild", name: "EndGuild", permMin: 0, handler: (*PlayObject).cmdEndGuild})
	registerGMCommand(&gmCmdDef{key: "AuthAlly", name: "AuthAlly", permMin: 0, handler: (*PlayObject).cmdAuthAlly})
	// usage 留空：Delphi 在命令内部自带 '?' 帮助文案（ObjBase.pas:11357-11361）。
	registerGMCommand(&gmCmdDef{key: "AllowGroupCall", name: "AllowGroupRecall", permMin: 0, handler: (*PlayObject).cmdAllowGroupRecall})
	registerGMCommand(&gmCmdDef{key: "GroupCall", name: "GroupRecall", permMin: 0, handler: (*PlayObject).cmdGroupRecall})
	registerGMCommand(&gmCmdDef{key: "AllowGuildReCall", name: "AllowGuildRecall", permMin: 0, handler: (*PlayObject).cmdAllowGuildRecall})
	registerGMCommand(&gmCmdDef{key: "Auth", name: "联盟", permMin: 0, handler: (*PlayObject).cmdGuildAlly})
	registerGMCommand(&gmCmdDef{key: "AuthCancel", name: "取消联盟", permMin: 0, handler: (*PlayObject).cmdGuildBreakAlly})
	// --- GM 物品段（ObjBase.pas:7750-7762）---
	registerGMCommand(&gmCmdDef{key: "ClearBag", name: "ClearBag", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdClearBag})
	registerGMCommand(&gmCmdDef{key: "ShowUseItemInfo", name: "ShowUseItemInfo", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdShowUseItemInfo})
	registerGMCommand(&gmCmdDef{key: "BindUseItem", name: "BindUseItem", permMin: 10, usage: "人物名称 物品类型 绑定方法", handler: (*PlayObject).cmdBindUseItem})
	// --- 登录锁（ObjBase.pas:7780-7783）---
	registerGMCommand(&gmCmdDef{key: "LockLogin", name: "LockLogin", permMin: 0, handler: (*PlayObject).cmdLockLogin})
}

// gmNameByKey 按 Delphi Command.ini 键取已解析命令字面量
//（供命令内文案引用当前命令名，如 @LetGuild 提示）。
func gmNameByKey(key, fallback string) string {
	for _, def := range gmCmdDefs {
		if def.key == key && def.resolvedName != "" {
			return def.resolvedName
		}
	}
	return fallback
}

// gmParamHelp 输出 Delphi 参数帮助（g_sGameCommandParamUnKnow，
// M2Share.pas:3251）"命令格式: @%s %s"。
func (p *PlayObject) gmParamHelp(server *netserver.TCPServer, key, usage string) {
	p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey(key, key), usage))
}

// cmdDate 显示服务器日期时间（Delphi 内联，ObjBase.pas:7613-7617；
// g_sNowCurrDateTime='前系统日期时间: '，M2Share.pas:3167，c_Blue）。
// 偏差：Delphi FormatDateTime('dddddd,dddd,hh:mm:nn') 长格式，
// Go 用固定 2006-01-02 15:04:05 格式。
func (p *PlayObject) cmdDate(server *netserver.TCPServer, params []string) bool {
	p.sysMsgBlue(server, "前系统日期时间: "+time.Now().Format("2006-01-02 15:04:05"))
	return true
}

// cmdAllowMsg 私聊开关（Delphi 内联，ObjBase.pas:7624-7631；
// g_sEnableHearWhisper/g_sDisableHearWhisper，M2Share.pas:3168-3169）。
func (p *PlayObject) cmdAllowMsg(server *netserver.TCPServer, params []string) bool {
	p.HearWhisper = !p.HearWhisper
	if p.HearWhisper {
		p.sysMsgGreen(server, "[允许私聊]")
	} else {
		p.sysMsgGreen(server, "[禁止私聊]")
	}
	return true
}

// cmdLetShout 群聊开关（Delphi 内联，ObjBase.pas:7632-7638；
// g_sEnableShoutMsg/g_sDisableShoutMsg，M2Share.pas:3170-3171。
// 注意 BanShout 为反义命名，TRUE=允许）。
func (p *PlayObject) cmdLetShout(server *netserver.TCPServer, params []string) bool {
	p.BanShout = !p.BanShout
	if p.BanShout {
		p.sysMsgGreen(server, "[允许群聊]")
	} else {
		p.sysMsgGreen(server, "[禁止群聊]")
	}
	return true
}

// cmdLetTrade 交易开关（Delphi 内联，ObjBase.pas:7639-7645；
// g_sEnableDealMsg/g_sDisableDealMsg，M2Share.pas:3172-3173）。
func (p *PlayObject) cmdLetTrade(server *netserver.TCPServer, params []string) bool {
	p.AllowDeal = !p.AllowDeal
	if p.AllowDeal {
		p.sysMsgGreen(server, "[允许交易]")
	} else {
		p.sysMsgGreen(server, "[禁止交易]")
	}
	return true
}

// cmdBanGuildChat 行会聊天开关（Delphi 内联，ObjBase.pas:7646-7653；
// g_sEnableGuildChat/g_sDisableGuildChat，M2Share.pas:3174-3175。
// 注意 BanGuildChat 为反义命名，TRUE=允许）。
func (p *PlayObject) cmdBanGuildChat(server *netserver.TCPServer, params []string) bool {
	p.BanGuildChat = !p.BanGuildChat
	if p.BanGuildChat {
		p.sysMsgGreen(server, "[允许行会聊天]")
	} else {
		p.sysMsgGreen(server, "[禁止行会聊天]")
	}
	return true
}

// cmdLetGuild 加入行会开关（Delphi 内联，ObjBase.pas:7654-7660；
// g_sEnableJoinGuild/g_sDisableJoinGuild，M2Share.pas:3176-3177）。
func (p *PlayObject) cmdLetGuild(server *netserver.TCPServer, params []string) bool {
	p.AllowGuildInvite = !p.AllowGuildInvite
	if p.AllowGuildInvite {
		p.sysMsgGreen(server, "[允许加入行会]")
	} else {
		p.sysMsgGreen(server, "[禁止加入行会]")
	}
	return true
}

// cmdEndGuild 退出行会（Delphi CmdEndGuild，ObjBase.pas:12841-12863）。
// 掌门人不可退出（Delphi m_nGuildRankNo>1 检查）；成功后清理自身行会
// 状态并广播 SMChangeGuildName 空名（与 HandleGuildDelMember 移除语义
// 一致，guildsystem.go:427-437）。
// 偏差：跨服 SS_207 行会同步（ObjBase.pas:12849）在 Go 单服架构下无
// 对应，跳过；持久化依赖周期性 SaveGuilds（main.go:554-556，与
// HandleGuildDelMember 一致）。
func (p *PlayObject) cmdEndGuild(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	guild := p.Engine.PlayerGuild(p.Name)
	if guild == nil {
		p.sysMsg(server, "你都没加入行会！！！")
		return true
	}
	if guild.Leader == p.Name {
		p.sysMsg(server, "行会掌门人不能这样退出行会！！！")
		return true
	}
	var remaining []string
	for _, m := range guild.Members {
		if m != p.Name {
			remaining = append(remaining, m)
		}
	}
	guild.Members = remaining
	delete(guild.Ranks, p.Name)
	p.GuildName = ""
	p.GuildRank = ""
	server.Send(p.Session.ID, protocol.MakeDefaultMsg(protocol.SMChangeGuildName, 0, 0, 0, 0), protocol.EncodeString(""))
	// 行会类装备需求失效自动脱下（Delphi CheckItemsNeed 链路）。
	p.checkAutoTakeOff(server)
	p.sysMsgGreen(server, "你已经退出行会。")
	return true
}

// cmdAuthAlly 行会联盟授权（Delphi 内联，ObjBase.pas:7667-7676）。
// Delphi：掌门人切换 TGuild.m_boEnableAuthAlly（Guild.pas:29，默认
// False，Guild.pas:322），该标志在 ClientGuildAlly 中作为被联盟方的
// 准入门槛（ObjBase.pas:18232），提示 g_sEnableAuthAllyGuild/
// g_sDisableAuthAllyGuild='[允许行会联盟]'/'[禁止行会联盟]'
//（M2Share.pas:3178-3179）。
// 偏差：Go Guild 结构无此开关字段（不新建持久化字段），行会联盟由
// 客户端界面经 CMGuildAlly/CMGuildBreakAlly 直接管理
//（guildsystem.go HandleGuildAlly/HandleGuildBreakAlly），故实现为
// 掌门人查询当前联盟列表的最小闭环。
func (p *PlayObject) cmdAuthAlly(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	guild := p.Engine.PlayerGuild(p.Name)
	// Delphi: 非掌门人（含无行会）静默退出。
	if guild == nil || guild.Leader != p.Name {
		return true
	}
	if len(guild.AllyGuilds) > 0 {
		p.sysMsgGreen(server, "当前联盟行会: "+strings.Join(guild.AllyGuilds, " "))
	} else {
		p.sysMsgGreen(server, "当前没有联盟行会。")
	}
	p.sysMsg(server, "行会联盟的建立与解除请使用客户端行会界面。")
	return true
}

// cmdGuildAlly 联盟（Delphi 分发：IsGuildMaster 后 ClientGuildAlly，
// ObjBase.pas:7699-7703/18217-18260）：两位掌门人面对面结盟。
// 要求双方均为掌门人、两行会互不处于行会战；成功后双向记录联盟
// 并向两会全体成员广播。失败码为 Delphi 原版的负数码
//（-1 无有效面对面对象/-2 行会战中/-3 双方非掌门人）。
// 偏差：① Delphi 目标行会 m_boEnableAuthAlly 准入门槛（18232）
// 在 Go 无对应开关（见 cmdAuthAlly），跳过；② 已结盟防重复为 Go
// 扩展（Delphi 此路径未显式检查）；③ 跨服 SS_207 跳过；
// ④ Delphi SendGuildMsg 走 RM_GUILDMESSAGE 行会频道，Go 用行会
// 系统消息替代。
func (p *PlayObject) cmdGuildAlly(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	guild := p.Engine.PlayerGuild(p.Name)
	// Delphi: 非掌门人静默退出（ObjBase.pas:7700）。
	if guild == nil || guild.Leader != p.Name {
		return true
	}
	fail := func(code int) {
		server.Send(p.Session.ID, protocol.MakeDefaultMsg(protocol.SMGuildMakeAllyFail, int32(code), 0, 0, 0), "")
	}
	// Delphi GetPoseCreate 双方面对面检查（ObjBase.pas:18225-18229）。
	pose := p.getFacingPlayer()
	if pose == nil || pose.getFacingPlayer() != p {
		fail(-1)
		return true
	}
	targetGuild := p.Engine.PlayerGuild(pose.Name)
	if targetGuild == nil {
		fail(-1)
		return true
	}
	if targetGuild.Leader != pose.Name {
		fail(-3)
		return true
	}
	if guild.IsWarGuild(targetGuild.Name) || targetGuild.IsWarGuild(guild.Name) {
		fail(-2)
		return true
	}
	if guild.IsAllyGuild(targetGuild.Name) {
		fail(-1)
		return true
	}
	guild.AllyGuilds = append(guild.AllyGuilds, targetGuild.Name)
	targetGuild.AllyGuilds = append(targetGuild.AllyGuilds, guild.Name)
	// Delphi 两会广播 '%s行会已经和您的行会联盟成功。'（18244-18245）。
	guild.sendGuildSystemMsg(p.Engine, server, targetGuild.Name+"行会已经和您的行会联盟成功。")
	targetGuild.sendGuildSystemMsg(p.Engine, server, guild.Name+"行会已经和您的行会联盟成功。")
	server.Send(p.Session.ID, protocol.MakeDefaultMsg(protocol.SMGuildMakeAllyOK, 0, 0, 0, 0), "")
	return true
}

// cmdGuildBreakAlly 取消联盟（Delphi 分发：IsGuildMaster 后
// ClientGuildBreakAlly(sParam1)，ObjBase.pas:7704-7708/18270-18300）：
// 参数=目标行会名。行会不存在→失败、非联盟→失败（Delphi 原版
// 失败统一回 SM_GUILDMAKEALLY_FAIL，18296-18299，原文如此）。
// 偏差：跨服 SS_207 跳过；RefMemberName 称号刷新在 Go 无对应。
func (p *PlayObject) cmdGuildBreakAlly(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	guild := p.Engine.PlayerGuild(p.Name)
	// Delphi: 非掌门人静默退出（ObjBase.pas:18276）。
	if guild == nil || guild.Leader != p.Name {
		return true
	}
	name := ""
	if len(params) > 0 {
		name = params[0]
	}
	targetGuild := p.Engine.FindGuild(name)
	if targetGuild != nil && guild.IsAllyGuild(name) {
		// 双向移除（与 HandleGuildBreakAlly 一致，guildsystem.go:486-503）。
		var remaining []string
		for _, a := range guild.AllyGuilds {
			if a != name {
				remaining = append(remaining, a)
			}
		}
		guild.AllyGuilds = remaining
		var tRemaining []string
		for _, a := range targetGuild.AllyGuilds {
			if a != guild.Name {
				tRemaining = append(tRemaining, a)
			}
		}
		targetGuild.AllyGuilds = tRemaining
		// Delphi 两会广播（18282-18283，前导空格为原文）。
		guild.sendGuildSystemMsg(p.Engine, server, name+" 行会与您的行会解除联盟成功！！！")
		targetGuild.sendGuildSystemMsg(p.Engine, server, guild.Name+" 行会解除了与您行会的联盟！！！")
		server.Send(p.Session.ID, protocol.MakeDefaultMsg(protocol.SMGuildBreakAllyOK, 0, 0, 0, 0), "")
		return true
	}
	server.Send(p.Session.ID, protocol.MakeDefaultMsg(protocol.SMGuildMakeAllyFail, 0, 0, 0, 0), "")
	return true
}

// cmdAllowGroupRecall 编组传送（天地合一）许可开关（Delphi
// CmdAllowGroupReCall，ObjBase.pas:11355-11368；g_sEnableGroupRecall/
// g_sDisableGroupRecall，M2Share.pas:3180-3181）。
// '?' 帮助为命令内自定义文案（ObjBase.pas:11357-11361），故注册表
// usage 留空、由此处分发。
func (p *PlayObject) cmdAllowGroupRecall(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && strings.HasPrefix(params[0], "?") {
		p.sysMsg(server, "此命令用于允许或禁止编组传送功能。")
		return true
	}
	p.AllowGroupReCall = !p.AllowGroupReCall
	if p.AllowGroupReCall {
		p.sysMsgGreen(server, "[允许天地合一]")
	} else {
		p.sysMsgGreen(server, "[禁止天地合一]")
	}
	return true
}

// cmdGroupRecall 编组传送（天地合一）执行（Delphi CmdGroupRecall，
// ObjBase.pas:12942-12996）：队长把已用 @AllowGroupRecall 许可的成员
// 全部召到自己身边。需召回套装或 GM≥6；地图 NORECALL 禁用；
// 冷却 nGroupRecallTime=180 秒（M2Share.pas:2010），GM≥6 免冷却。
// 文案为 Delphi 原文（中英混合即原版如此）。
// 偏差：Delphi m_wGroupRcallTime 秒数递减状态机简化为到期时间戳
//（GroupRecallTick），行为等价；跨服 SS_207 无对应跳过。
func (p *PlayObject) cmdGroupRecall(server *netserver.TCPServer, params []string) bool {
	if !p.HasRecallSuite && int(p.Permission) < 6 {
		p.sysMsg(server, "您现在还无法使用此功能！！！")
		return true
	}
	if p.envir != nil && p.envir.Flag.NoRecall {
		p.sysMsg(server, "此地图禁止使用此命令！！！")
		return true
	}
	now := time.Now().UnixMilli()
	if int(p.Permission) < 6 && now < p.GroupRecallTick {
		left := (p.GroupRecallTick-now)/1000 + 1
		p.sysMsg(server, fmt.Sprintf("%d seconds, You can use it again.", left))
		return true
	}
	party := p.partyOf()
	// Delphi: 非队长（m_GroupOwner<>Self）静默退出。
	if party == nil || party.Leader != p.ID || p.Engine == nil {
		return true
	}
	for _, id := range party.Members {
		if id == p.ID {
			continue
		}
		member := p.Engine.GetPlayer(id)
		if member == nil {
			continue
		}
		if !member.AllowGroupReCall {
			p.sysMsg(server, fmt.Sprintf("%s is now rejecting GroupRecall.", member.Name))
			continue
		}
		if member.envir != nil && member.envir.Flag.NoRecall {
			p.sysMsg(server, fmt.Sprintf("%s map does not allow recall.", member.Name))
			continue
		}
		x, y, ok := findWalkableNear(p.envir, p.CurrX, p.CurrY, 3)
		if !ok {
			x, y = p.CurrX, p.CurrY
		}
		teleportPlayerTo(member, server, p.envir, x, y)
	}
	if int(p.Permission) < 6 {
		p.GroupRecallTick = now + 180000
	}
	return true
}

// cmdAllowGuildRecall 行会合一许可开关（Delphi 内联，
// ObjBase.pas:7687-7693；g_sEnableGuildRecall/g_sDisableGuildRecall，
// M2Share.pas:3182-3183）。
func (p *PlayObject) cmdAllowGuildRecall(server *netserver.TCPServer, params []string) bool {
	p.AllowGuildReCall = !p.AllowGuildReCall
	if p.AllowGuildReCall {
		p.sysMsgGreen(server, "[允许行会合一]")
	} else {
		p.sysMsgGreen(server, "[禁止行会合一]")
	}
	return true
}

// cmdPrvMsg 私聊拉黑切换（Delphi CmdPrvMsg，ObjBase.pas:13969-13996；
// 文案 g_sGameCommandPrvMsgUnLimitMsg/g_sGameCommandPrvMsgLimitMsg，
// M2Share.pas:3333-3334，c_Green）。
func (p *PlayObject) cmdPrvMsg(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.gmParamHelp(server, "PrvMsg", "人物名称")
		return true
	}
	name := params[0]
	for i, n := range p.BlockWhisper {
		if strings.EqualFold(n, name) {
			p.BlockWhisper = append(p.BlockWhisper[:i], p.BlockWhisper[i+1:]...)
			p.sysMsgGreen(server, fmt.Sprintf("%s 已从禁止私聊列表中删除...", name))
			return true
		}
	}
	p.BlockWhisper = append(p.BlockWhisper, name)
	p.sysMsgGreen(server, fmt.Sprintf("%s 已被加入禁止私聊列表...", name))
	return true
}

// isBlockWhisper 私聊拉黑名单查询（Delphi IsBlockWhisper，
// ObjBase.pas:2217-2230；CompareText 大小写不敏感）。
func (p *PlayObject) isBlockWhisper(name string) bool {
	for _, n := range p.BlockWhisper {
		if strings.EqualFold(n, name) {
			return true
		}
	}
	return false
}

// cmdMemberFunc 会员功能窗口（Delphi CmdMemberFunction，
// ObjBase.pas:13566-13578：g_ManageNPC.GotoLable(Self,'@Member',False)）。
// '?' 帮助为命令内自定义文案（ObjBase.pas:13568-13572），故注册表
// usage 留空、由此处分发。
func (p *PlayObject) cmdMemberFunc(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && strings.HasPrefix(params[0], "?") {
		p.sysMsg(server, "打开会员功能窗口.")
		return true
	}
	p.gotoMemberLabel(server)
	return true
}

// cmdMemberFuncEx 会员功能窗口（功能 NPC 版）（Delphi
// CmdMemberFunctionEx，ObjBase.pas:13580-13592：
// g_FunctionNPC.GotoLable(Self,'@Member',False)）。
func (p *PlayObject) cmdMemberFuncEx(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && strings.HasPrefix(params[0], "?") {
		p.sysMsg(server, "打开会员功能窗口.")
		return true
	}
	p.gotoMemberLabel(server)
	return true
}

// gotoMemberLabel 触发同地图最近管理型 NPC 的 @Member 标签
//（Delphi g_ManageNPC/g_FunctionNPC.GotoLable(Self,'@Member',False)）。
func (p *PlayObject) gotoMemberLabel(server *netserver.TCPServer) {
	p.gotoNpcLabel(server, "Member")
}

// gotoNpcLabel 查找同地图最近的带指定标签的管理型 NPC（IsMerchant=false）
// 并执行其脚本标签。
// 偏差：Delphi 走全局单例 g_ManageNPC/g_FunctionNPC（M2Share.pas
// 变量，启动时由 Envir 配置指定）；Go 未实装管理/功能 NPC 单例
//（itemsystem.go:566 同样注明"Go 无 FunctionNPC"），改为按地图就近
// 扫描，找不到时静默（与 Delphi g_ManageNPC=nil 时一致）。
func (p *PlayObject) gotoNpcLabel(server *netserver.TCPServer, label string) {
	if p.Engine == nil || p.envir == nil {
		return
	}
	var target *NpcObject
	var script *NpcScript
	bestDist := -1
	for _, npc := range p.Engine.Npcs {
		if npc.IsMerchant || npc.envir != p.envir {
			continue
		}
		s := npc.GetScript()
		if s == nil {
			continue
		}
		// 脚本标签存储时已剥离 '@' 前缀（npcscript.go:98-100）。
		if _, ok := s.Labels[label]; !ok {
			continue
		}
		d := abs(npc.CurrX-p.CurrX) + abs(npc.CurrY-p.CurrY)
		if target == nil || d < bestDist {
			target, script, bestDist = npc, s, d
		}
	}
	if target == nil {
		return
	}
	script.Execute(label, p, target, server)
}

// cmdLockLogin 登录锁开关（Delphi CmdLockLogin，ObjBase.pas:13272-13300）。
// 偏差：Go 未实装登录锁（Delphi 依赖 g_Config.boLockHumanLogin 与
// CM_PASSWORD 聊天状态机 m_boLockLogon/m_boLockLogoned），恒走
// "未启用"分支（ObjBase.pas:13278-13282）。
func (p *PlayObject) cmdLockLogin(server *netserver.TCPServer, params []string) bool {
	p.sysMsg(server, "本服务器还没有启用登录锁功能！！！")
	return true
}

// cmdClearBag 清空目标背包（Delphi CmdClearBagItem，
// ObjBase.pas:11048-11087，权限 10）。
// 偏差：Delphi 构造删除清单经 RM_SENDDELITEMLIST 通知客户端逐条删除；
// Go 用 SendBagItemsFull 全量刷新背包替代（等价闭环）。
func (p *PlayObject) cmdClearBag(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.gmParamHelp(server, "ClearBag", "人物名称")
		return true
	}
	if p.Engine == nil {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf("%s 现在不在线，或在其它服务器上！！！", params[0]))
		return true
	}
	target.ItemList = nil
	target.SendBagItemsFull(server)
	return true
}

// useItemSlotName 穿戴槽位名（Delphi GetUseItemName，M2Share.pas:4469-4484；
// 1.76 启用英文字面量 M2Share.pas:3373-3382，腰带/靴子/宝石被注释，
// 10-12 槽原版返回空串）。
func useItemSlotName(slot int) string {
	switch slot {
	case protocol.UDress:
		return "Dress"
	case protocol.UWeapon:
		return "Weapon"
	case protocol.URightHand:
		return "RightHand"
	case protocol.UNecklace:
		return "Necklace"
	case protocol.UHelmet:
		return "Helmet"
	case protocol.UArmRingL:
		return "BraceL"
	case protocol.UArmRingR:
		return "BraceR"
	case protocol.URingL:
		return "RingL"
	case protocol.URingR:
		return "RingR"
	case protocol.UBujuk:
		return "Bujuk"
	}
	return ""
}

// useItemSlotIdx 槽位名→索引（Delphi GetUseItemIdx，M2Share.pas:4407-4468；
// CompareText 大小写不敏感）。
func useItemSlotIdx(name string) int {
	switch {
	case strings.EqualFold(name, "Dress"):
		return protocol.UDress
	case strings.EqualFold(name, "Weapon"):
		return protocol.UWeapon
	case strings.EqualFold(name, "RightHand"):
		return protocol.URightHand
	case strings.EqualFold(name, "Necklace"):
		return protocol.UNecklace
	case strings.EqualFold(name, "Helmet"):
		return protocol.UHelmet
	case strings.EqualFold(name, "BraceL"):
		return protocol.UArmRingL
	case strings.EqualFold(name, "BraceR"):
		return protocol.UArmRingR
	case strings.EqualFold(name, "RingL"):
		return protocol.URingL
	case strings.EqualFold(name, "RingR"):
		return protocol.URingR
	case strings.EqualFold(name, "Bujuk"):
		return protocol.UBujuk
	}
	return -1
}

// itemNameOf 物品实例的定义名（Delphi UserEngine.GetStdItemName）。
func itemNameOf(p *PlayObject, wIndex uint16) string {
	if p != nil && p.ItemDB != nil {
		if def := p.ItemDB.GetByIdx(int(wIndex)); def != nil {
			return def.Name
		}
	}
	return "?"
}

// cmdShowUseItemInfo 显示目标穿戴物品明细（Delphi CmdShowUseItemInfo，
// ObjBase.pas:14437-14474，权限 10；格式串为原文内联，c_Blue）。
func (p *PlayObject) cmdShowUseItemInfo(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.gmParamHelp(server, "ShowUseItemInfo", "人物名称")
		return true
	}
	if p.Engine == nil {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf("%s 现在不在线，或在其它服务器上！！！", params[0]))
		return true
	}
	for i := range target.UseItems {
		item := target.UseItems[i]
		// Delphi: UserItem.wIndex = 0 then Continue。
		if item == nil || item.WIndex == 0 {
			continue
		}
		p.sysMsgBlue(server, fmt.Sprintf("%s[%s]IDX[%d]系列号[%d]持久[%d-%d]",
			useItemSlotName(i), itemNameOf(target, item.WIndex),
			item.WIndex, item.MakeIndex, item.Dura, item.DuraMax))
	}
	return true
}

// itemBindRegistry 装备绑定登记表（Go 最小闭环）。
// Delphi 按绑定方式分三张全局表 g_ItemBindAccount/g_ItemBindCharName/
// g_ItemBindIPaddr，以 ItemIdx+MakeIndex 查重并 SaveItemBind* 落盘
//（ObjBase.pas:14516-14634）；Go UserItem 协议结构无绑定字段，且未
// 实装落盘与登录/交易校验协议，此处用 MakeIndex（实例唯一 ID）为键
// 的内存表作为最小闭环。偏差：不落盘、重启丢失、不参与任何校验、
// 不区分绑定方式。
var (
	itemBindMu       sync.Mutex
	itemBindRegistry = map[int32]string{}
)

// cmdBindUseItem 绑定目标穿戴物品（Delphi CmdBindUseItem，
// ObjBase.pas:14476-14637，权限 10）。
// 参数：人物名称 物品类型(Dress/Weapon/...) 绑定方法(帐号/人物/IP)。
func (p *PlayObject) cmdBindUseItem(server *netserver.TCPServer, params []string) bool {
	if len(params) < 3 {
		p.gmParamHelp(server, "BindUseItem", "人物名称 物品类型 绑定方法")
		return true
	}
	humanName, itemName, bindType := params[0], params[1], params[2]
	slot := useItemSlotIdx(itemName)
	bind := -1
	// Delphi CompareText(sType, '帐号'/'人物'/'IP')（ObjBase.pas:14494-14496）。
	switch {
	case strings.EqualFold(bindType, "帐号"):
		bind = 0
	case strings.EqualFold(bindType, "人物"):
		bind = 1
	case strings.EqualFold(bindType, "IP"):
		bind = 2
	}
	if slot < 0 || bind < 0 {
		p.gmParamHelp(server, "BindUseItem", "人物名称 物品类型 绑定方法")
		return true
	}
	if p.Engine == nil {
		return true
	}
	target := p.Engine.GetPlayerByName(humanName)
	if target == nil {
		p.sysMsg(server, fmt.Sprintf("%s 现在不在线，或在其它服务器上！！！", humanName))
		return true
	}
	item := target.UseItems[slot]
	if item == nil || item.WIndex == 0 {
		// g_sGameCommandBindUseItemNoItemMsg（M2Share.pas:3352）。
		p.sysMsg(server, fmt.Sprintf("%s的%s没有戴物品！！！", humanName, itemName))
		return true
	}
	// 绑定名：帐号→m_sUserID，人物→m_sCharName，IP→m_sIPaddr
	//（ObjBase.pas:14521/14560/14599）。
	var bindName string
	switch bind {
	case 0:
		bindName = target.AccountName
	case 1:
		bindName = target.Name
	default:
		bindName = sessionRemoteIP(target.Session)
	}
	itemBindMu.Lock()
	if _, exists := itemBindRegistry[item.MakeIndex]; exists {
		itemBindMu.Unlock()
		// g_sGameCommandBindUseItemAlreadBindMsg（M2Share.pas:3353）。
		p.sysMsg(server, fmt.Sprintf("%s的%s上的物品早已绑定过了！！！", humanName, itemName))
		return true
	}
	itemBindRegistry[item.MakeIndex] = fmt.Sprintf("%s:%s", bindType, bindName)
	itemBindMu.Unlock()
	slotName := useItemSlotName(slot)
	stdName := itemNameOf(target, item.WIndex)
	// Delphi 双发蓝字：GM 侧成功明细 + 目标侧通知
	//（ObjBase.pas:14541-14554 等三段相同结构）。
	p.sysMsgBlue(server, fmt.Sprintf("%s[%s]IDX[%d]系列号[%d]持久[%d-%d]，绑定到%s成功。",
		slotName, stdName, item.WIndex, item.MakeIndex, item.Dura, item.DuraMax, bindName))
	target.sysMsgBlue(server, fmt.Sprintf("你的%s[%s]已经绑定到%s[%s]上了。",
		slotName, stdName, bindType, bindName))
	return true
}
