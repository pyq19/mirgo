package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// GM 命令——行会/城堡族移植。
// Delphi 参照：ObjBase.pas（分发 7904-7907 常规段、8214-8232 争霸赛段、
// 8490-8508 调试内层段；实现 10483-14435）+ Guild.pas（AddGuild 114-127、
// DELGUILD 128-147、StartTeamFight/EndTeamFight 1272-1287）+
// Castle.pas（GetCastle 828-846、Start/StopWallconquestWar 848-889）+
// M2Share.pas g_GameCommand（2785,2831-2834,2873-2876）与文案常量（3290-3292）。
// 注册顺序 = Delphi 分发链出现顺序。
// 跨服消息（SS_204/205/206/211/212）在 Go 单服架构下无对应，全部跳过。

func init() {
	registerGMCommand(&gmCmdDef{key: "SabukWallGold", name: "SabukWallGold", permMin: 10, usage: "城堡名称 控制符(=、-、+) 金币数(1-100000000)", handler: (*PlayObject).cmdSabukWallGold})
	// --- 争霸赛段（ObjBase.pas:8214-8232，常规 nPermissionMin 检查）---
	registerGMCommand(&gmCmdDef{key: "ContestPoint", name: "ContestPoint", permMin: 10, handler: (*PlayObject).cmdContestPoint})
	registerGMCommand(&gmCmdDef{key: "StartContest", name: "StartContest", permMin: 10, handler: (*PlayObject).cmdStartContest})
	registerGMCommand(&gmCmdDef{key: "EndContest", name: "EndContest", permMin: 10, handler: (*PlayObject).cmdEndContest})
	registerGMCommand(&gmCmdDef{key: "Announcement", name: "Announcement", permMin: 10, handler: (*PlayObject).cmdAnnouncement})
	// --- 调试内层段（ObjBase.pas:8490-8508，>=5/boTestServer 门槛内，
	// 实现走常规 nPermissionMin 检查）---
	registerGMCommand(&gmCmdDef{key: "AddGuild", name: "AddGuild", permMin: 10, usage: "行会名称 掌门人名称", handler: (*PlayObject).cmdAddGuild})
	registerGMCommand(&gmCmdDef{key: "DelGuild", name: "DelGuild", permMin: 10, usage: "行会名称", handler: (*PlayObject).cmdDelGuild})
	registerGMCommand(&gmCmdDef{key: "ChangeSabukLord", name: "ChangeSabukLord", permMin: 10, usage: "城堡名称 行会名称", handler: (*PlayObject).cmdChangeSabukLord})
	registerGMCommand(&gmCmdDef{key: "ForcedWallconquestWar", name: "ForcedWallconquestWar", permMin: 10, usage: "城堡名称", handler: (*PlayObject).cmdForcedWallconquestWar})
}

// cryMap 地图范围喊话广播（Delphi UserEngine.CryCry，UsrEngn.pas:1500-1530：
// RM_CRY 同图广播，争霸赛公告用）。
func (p *PlayObject) cryMap(server *netserver.TCPServer, env *Environment, text string) {
	if env == nil || p.Engine == nil {
		return
	}
	cryMsg := protocol.MakeDefaultMsg(protocol.SMCry, p.ID, chatColorCry, 0, 0)
	body := protocol.EncodeString(p.Name + "/" + text)
	p.Engine.mu.RLock()
	for _, other := range p.Engine.PlayObjectList {
		if other.Ghost || other.envir != env {
			continue
		}
		server.Send(other.Session.ID, cryMsg, body)
	}
	p.Engine.mu.RUnlock()
}

// findGuildFold 行会名大小写不敏感查找（Delphi GuildManager.DELGUILD 用
// CompareText，Guild.pas:131；区别于 FindGuild 的 '=' 精确匹配）。
func (e *UserEngine) findGuildFold(name string) *Guild {
	for _, g := range e.Guilds {
		if strings.EqualFold(g.Name, name) {
			return g
		}
	}
	return nil
}

// cmdAddGuild 为在线玩家创建行会（Delphi CmdAddGuild，
// ObjBase.pas:10483-10534 → GuildManager.AddGuild，Guild.pas:114-127）。
// Delphi 原版行为照搬：掌门人已入行会/行会重名时静默无输出；
// 原版 GetRankName(Self,...) 误传 GM 导致掌门人称号被置空的 bug 不继承，
// Go 直接赋"掌门人"（与 HandleBuildGuild 一致）。
func (p *PlayObject) cmdAddGuild(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("AddGuild", "AddGuild")+" 行会名称 掌门人名称")
		return true
	}
	e := p.Engine
	chief := e.GetPlayerByName(params[1])
	if chief == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[1]))
		return true
	}
	if e.PlayerGuild(chief.Name) != nil || e.FindGuild(params[0]) != nil {
		return true // Delphi 静默（ObjBase.pas:10511/10518 条件落空分支）
	}
	guild := &Guild{
		Name:    params[0],
		Leader:  chief.Name,
		Members: []string{chief.Name},
		Ranks:   map[string]string{chief.Name: "掌门人"},
	}
	e.Guilds = append(e.Guilds, guild)
	e.SaveGuilds()
	chief.GuildName = guild.Name
	chief.GuildRank = "掌门人"
	// Delphi RM_CHANGEGUILDNAME 刷新头顶行会名；Go 与 HandleGuildAddMember
	// 一致发 SMChangeGuildName（guildsystem.go:408-409）。
	server.Send(chief.Session.ID, protocol.MakeDefaultMsg(protocol.SMChangeGuildName, 0, 0, 0, 0), protocol.EncodeString(guild.Name))
	p.sysMsgGreen(server, "行会名称: "+guild.Name+" 掌门人: "+chief.Name)
	return true
}

// cmdDelGuild 删除行会（Delphi CmdDelGuild，ObjBase.pas:12343-12368
// → GuildManager.DELGUILD，Guild.pas:128-147 + BackupGuildFile 807-839）。
// 偏差：① Delphi 多职位行会（m_RankList.Count>1）不可删的限制在 Go 无
// 职位定义模型，不移植；② Delphi 行会文件备份（.bak）改为直接删 DB
// 记录（storage.DeleteGuild）；③ 成功无 SysMsg 为 Delphi 原版行为
// （仅发 SS_206 跨服消息），照搬。
func (p *PlayObject) cmdDelGuild(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("DelGuild", "DelGuild")+" 行会名称")
		return true
	}
	e := p.Engine
	guild := e.findGuildFold(params[0])
	if guild == nil {
		p.sysMsg(server, "没找到"+params[0]+"这个行会！！！")
		return true
	}
	// Delphi BackupGuildFile：清空全部在线成员的行会状态并刷新显示名。
	e.mu.Lock()
	members := make([]*PlayObject, 0, len(guild.Members))
	for _, pl := range e.PlayObjectList {
		if pl.GuildName == guild.Name {
			members = append(members, pl)
		}
	}
	e.mu.Unlock()
	for _, pl := range members {
		pl.GuildName = ""
		pl.GuildRank = ""
		server.Send(pl.Session.ID, protocol.MakeDefaultMsg(protocol.SMChangeGuildName, 0, 0, 0, 0), protocol.EncodeString(""))
		// 行会类装备需求失效自动脱下（与 cmdEndGuild 一致）。
		pl.checkAutoTakeOff(server)
	}
	for i, g := range e.Guilds {
		if g == guild {
			e.Guilds = append(e.Guilds[:i], e.Guilds[i+1:]...)
			break
		}
	}
	if e.db != nil {
		if err := e.db.DeleteGuild(guild.Name); err != nil {
			log.Logf(log.LevelError, "Guild", "failed to delete guild %s: %v", guild.Name, err)
		}
	}
	return true
}

// cmdChangeSabukLord 更改城堡所属行会（Delphi CmdChangeSabukLord，
// ObjBase.pas:10940-10984 → Castle.GetCastle，Castle.pas:828-846）。
// 偏差：原版 AddGameDataLog 第 5 字段误写字面量 'sGuildName' 的
// 复制粘贴 bug 不继承（Go 无 GameDataLog）。
func (p *PlayObject) cmdChangeSabukLord(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("ChangeSabukLord", "ChangeSabukLord")+" 城堡名称 行会名称")
		return true
	}
	e := p.Engine
	castle := e.Castle
	if castle == nil || !strings.EqualFold(params[0], castle.Config.Name) {
		// g_sGameCommandSbkGoldCastleNotFoundMsg（M2Share.pas:3291）。
		p.sysMsg(server, fmt.Sprintf("城堡%s未找到！！！", params[0]))
		return true
	}
	guild := e.FindGuild(params[1])
	if guild == nil {
		p.sysMsg(server, "行会 "+params[1]+"还没建立！！！")
		return true
	}
	castle.ChangeOwner(server, e, guild.Name)
	p.sysMsgGreen(server, castle.Config.Name+" 所属行会已经更改成 "+guild.Name)
	return true
}

// cmdForcedWallconquestWar 强制切换攻城战状态（Delphi
// CmdForcedWallconquestWar，ObjBase.pas:12877-12916：翻转 m_boUnderWar，
// 开战广播 '[X攻城战已经开始]' + 关城门，停战广播 '[X攻城战已经结束]'）。
// 对执行者无 SysMsg 反馈为 Delphi 原版行为，照搬。
// 偏差：Delphi 停战分支 GetMapOfRangeHumanCount 未填充清理列表的 bug
// 不继承（Go EndWar 完整清理状态）。
func (p *PlayObject) cmdForcedWallconquestWar(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("ForcedWallconquestWar", "ForcedWallconquestWar")+" 城堡名称")
		return true
	}
	e := p.Engine
	castle := e.Castle
	if castle == nil || !strings.EqualFold(params[0], castle.Config.Name) {
		p.sysMsg(server, fmt.Sprintf("城堡%s未找到！！！", params[0]))
		return true
	}
	if castle.IsAtWar() {
		castle.EndWar(server, e)
	} else {
		castle.ForceStartWar(server, e)
	}
	return true
}

// cmdSabukWallGold 查询/调整城堡金库（Delphi CmdShowSbkGold，
// ObjBase.pas:14383-14435；g_sGameCommandSbkGoldShowMsg/HelpMsg，
// M2Share.pas:3290-3292）。
// 无参=查询（Go 单城堡 ≡ Delphi 全部城堡列表）；带参=修改。
// 原版 bug 修正：'-' 分支 Dec 不带步长只减 1（ObjBase.pas:14424），
// Go 按语义减 nGold。修改后无反馈消息为 Delphi 原版行为，照搬。
func (p *PlayObject) cmdSabukWallGold(server *netserver.TCPServer, params []string) bool {
	castle := p.Engine.Castle
	if castle == nil {
		return true
	}
	if len(params) == 0 {
		castle.mu.Lock()
		gold, income := castle.Gold, castle.TodayIncome
		castle.mu.Unlock()
		p.sysMsgGreen(server, fmt.Sprintf("%s的金币数为: %d 今天收入: %d", castle.Config.Name, gold, income))
		return true
	}
	if !strings.EqualFold(params[0], castle.Config.Name) {
		p.sysMsg(server, fmt.Sprintf("城堡%s未找到！！！", params[0]))
		return true
	}
	ctr := byte(0) // Delphi Ctr := sCtr[1]，空串时 #0 触发用法分支
	if len(params) >= 2 && params[1] != "" {
		ctr = params[1][0]
	}
	nGold := -1 // Delphi Str_ToInt(sGold, -1)
	if len(params) >= 3 {
		nGold = atoiOr(params[2], -1)
	}
	if len(params) < 3 || (ctr != '=' && ctr != '-' && ctr != '+') || nGold < 0 || nGold > 100000000 {
		p.gmParamHelp(server, "SabukWallGold", "城堡名称 控制符(=、-、+) 金币数(1-100000000)")
		return true
	}
	castle.mu.Lock()
	switch ctr {
	case '=':
		castle.Gold = int64(nGold)
	case '+':
		castle.Gold += int64(nGold)
	case '-':
		castle.Gold -= int64(nGold)
	}
	if castle.Gold < 0 {
		castle.Gold = 0
	}
	castle.mu.Unlock()
	castle.SaveState(p.Engine.DB())
	return true
}

// contestPlayers 收集当前 FIGHT3 地图的参赛玩家（Delphi
// GetMapRageHuman 半径 1000 ≡ 全图，UsrEngn.pas:2938-2958）。
// 原版 bug 修正：Delphi 观察/管理员过滤条件误用 OR（ObjBase.pas:
// 11246/11315），Go 按意图用 AND 排除。
func (p *PlayObject) contestPlayers() []*PlayObject {
	e := p.Engine
	e.mu.Lock()
	players := make([]*PlayObject, 0, 16)
	for _, pl := range e.PlayObjectList {
		if pl.Death || pl.Ghost || pl.envir != p.envir {
			continue
		}
		if pl.ObMode || pl.AdminMode {
			continue
		}
		players = append(players, pl)
	}
	e.mu.Unlock()
	return players
}

// contestGuilds 参赛行会去重列表（按首次出现顺序）。
func (p *PlayObject) contestGuilds(players []*PlayObject) []*Guild {
	e := p.Engine
	var guilds []*Guild
	seen := map[string]bool{}
	for _, pl := range players {
		g := e.PlayerGuild(pl.Name)
		if g == nil || seen[g.Name] {
			continue
		}
		seen[g.Name] = true
		guilds = append(guilds, g)
	}
	return guilds
}

// cmdStartContest 开始行会争霸赛（Delphi CmdStartContest，
// ObjBase.pas:11226-11295）：FIGHT3 地图内收集参赛行会，
// StartTeamFight + 成员登记，公告开战与参赛门派。
// 偏差：Delphi 内层循环误用外层下标 i 致成员重复登记/漏登记的 bug
// （ObjBase.pas:11275）不继承，Go 逐玩家正确登记；
// Delphi 逐玩家清零 m_nFightZoneDieCount 在 Go 无对应字段
// （团队战死亡计数走 TeamFightRecord），跳过。
func (p *PlayObject) cmdStartContest(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.sysMsg(server, "开始行会争霸赛。")
		p.sysMsg(server, "命令格式: @"+gmNameByKey("StartContest", "StartContest"))
		return true
	}
	if p.envir == nil || !p.envir.Flag.Fight3 {
		p.sysMsg(server, "此命令不能在当前地图中使用！！！")
		return true
	}
	players := p.contestPlayers()
	guilds := p.contestGuilds(players)
	p.sysMsgGreen(server, "行会争霸赛已经开始。")
	p.cryMap(server, p.envir, "- 行会战争已爆发。")
	for _, g := range guilds {
		g.StartTeamFight()
		for _, pl := range players {
			if p.Engine.PlayerGuild(pl.Name) == g {
				g.AddTeamFightMember(pl.Name)
			}
		}
	}
	names := make([]string, 0, len(guilds))
	for _, g := range guilds {
		names = append(names, g.Name)
	}
	p.cryMap(server, p.envir, " -参加的门派:"+strings.Join(names, " "))
	return true
}

// cmdEndContest 结束行会争霸赛（Delphi CmdEndContest，
// ObjBase.pas:11297-11357）：逐行会 EndTeamFight 并公告。
// 无 SysMsg 反馈为 Delphi 原版行为，照搬。
func (p *PlayObject) cmdEndContest(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.sysMsg(server, "结束行会争霸赛。")
		p.sysMsg(server, "命令格式: @"+gmNameByKey("EndContest", "EndContest"))
		return true
	}
	if p.envir == nil || !p.envir.Flag.Fight3 {
		p.sysMsg(server, "此命令不能在当前地图中使用！！！")
		return true
	}
	guilds := p.contestGuilds(p.contestPlayers())
	for _, g := range guilds {
		g.EndTeamFight()
		p.cryMap(server, p.envir, fmt.Sprintf(" - %s 行会争霸赛已结束。", g.Name))
	}
	return true
}

// cmdContestPoint 查询行会争霸赛得分（Delphi CmdContestPoint，
// ObjBase.pas:11200-11224）。
// 原版怪癖照抄：行会不存在的错误提示也用 c_Green（ObjBase.pas:11219）。
func (p *PlayObject) cmdContestPoint(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 || (params[0] != "" && params[0][0] == '?') {
		p.sysMsg(server, "查看行会战的得分数。")
		p.sysMsg(server, fmt.Sprintf("命令格式: @%s 行会名称", gmNameByKey("ContestPoint", "ContestPoint")))
		return true
	}
	guild := p.Engine.FindGuild(params[0])
	if guild == nil {
		p.sysMsgGreen(server, fmt.Sprintf("行会: %s 不存在！！！", params[0]))
		return true
	}
	p.sysMsgGreen(server, fmt.Sprintf("%s 的得分为: %d", guild.Name, guild.ContestPoint))
	return true
}

// cmdAnnouncement 公告行会争霸赛结果（Delphi CmdAnnouncement，
// ObjBase.pas:11369-11429：CryCry 逐成员 '得分/死亡次数'）。
// 原版 bug 修正：尾部两条 CryCry 在 Guild<>nil 检查之外，行会不存在时
// 访问违规（ObjBase.pas:11421-11426）；Go 行会不存在时输出错误提示。
// 偏差：Delphi TeamFightDeadList 为有序列表，Go map 遍历需排序保证
// 输出稳定。
func (p *PlayObject) cmdAnnouncement(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 || (params[0] != "" && params[0][0] == '?') {
		p.sysMsg(server, "查看行会争霸赛结果。")
		p.sysMsg(server, fmt.Sprintf("命令格式: @%s 行会名称", gmNameByKey("Announcement", "Announcement")))
		return true
	}
	if p.envir == nil || !p.envir.Flag.Fight3 {
		p.sysMsg(server, "此命令不能在当前地图中使用！！！")
		return true
	}
	guild := p.Engine.FindGuild(params[0])
	if guild == nil {
		p.sysMsgGreen(server, fmt.Sprintf("行会: %s 不存在！！！", params[0]))
		return true
	}
	p.cryMap(server, p.envir, fmt.Sprintf(" - %s 行会争霸赛结果: ", guild.Name))
	names := make([]string, 0, len(guild.TeamFightMembers))
	for name := range guild.TeamFightMembers {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		rec := guild.TeamFightMembers[name]
		p.cryMap(server, p.envir, fmt.Sprintf(" - %s  : %d 分/死亡%d次。 ", name, rec.Score, rec.Deaths))
	}
	p.cryMap(server, p.envir, fmt.Sprintf(" - [%s] : %d 分。", guild.Name, guild.ContestPoint))
	p.cryMap(server, p.envir, "------------------------------------")
	return true
}
