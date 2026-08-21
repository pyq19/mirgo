package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// GM 命令——属性/经济/技能族移植。
// Delphi 参照：ObjBase.pas（分发 ProcessUserLineMsg 7188-8600）+
// M2Share.pas g_GameCommand 命令表（2724-2926）与文案常量（3250-3350）。
// 注册顺序 = Delphi 分发链出现顺序（ObjBase.pas:7963-8485）；
// 调试段（m_btPermission>4 门槛，ObjBase.pas:8253）内的命令统一注册为
// permMin 6 + permFailSilent，对齐 Delphi 硬编码静默退出语义。

func init() {
	registerGMCommand(&gmCmdDef{key: "LuckPoint", name: "LuckyPoint", permMin: 10, usage: "人物名称 控制符 幸运点数", handler: (*PlayObject).cmdLuckyPoint})
	registerGMCommand(&gmCmdDef{key: "ChangeJob", name: "ChangeJob", permMin: 10, usage: "人物名称 职业类型(Warr Wizard Taos)", handler: (*PlayObject).cmdChangeJob})
	registerGMCommand(&gmCmdDef{key: "ChangeGender", name: "ChangeGender", permMin: 10, usage: "人物名称 性别(男、女)", handler: (*PlayObject).cmdChangeGender})
	registerGMCommand(&gmCmdDef{key: "Hair", name: "hair", permMin: 10, usage: "人物名称 类型值", handler: (*PlayObject).cmdHair})
	registerGMCommand(&gmCmdDef{key: "BonusPoint", name: "BonusPoint", permMin: 10, usage: "人物名称 属性点数(不输入为查看点数)", handler: (*PlayObject).cmdBonusPoint})
	registerGMCommand(&gmCmdDef{key: "DelBonuPoint", name: "DelBonusPoint", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdDelBonusPoint})
	registerGMCommand(&gmCmdDef{key: "RestBonuPoint", name: "RestBonusPoint", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdRestBonusPoint})
	registerGMCommand(&gmCmdDef{key: "SetPermission", name: "SetPermission", permMin: 10, usage: "人物名称 权限等级(0 - 10)", handler: (*PlayObject).cmdSetPermission})
	// Delphi 实现内硬编码 m_btPermission<6 静默退出（ObjBase.pas:11558），
	// 命令表 nPermissionMin=10 未被实现消费，注册按实际行为 6+silent。
	registerGMCommand(&gmCmdDef{key: "ReNewLevel", name: "ReNewLevel", permMin: 6, permFail: permFailSilent, usage: "人物名称 点数(为空则查看)", handler: (*PlayObject).cmdReNewLevel})
	registerGMCommand(&gmCmdDef{key: "GameGold", name: "GameGold", permMin: 10, usage: "人物名称 控制符(+,-,=) 游戏币(1-200000000)", handler: (*PlayObject).cmdGameGold})
	registerGMCommand(&gmCmdDef{key: "GamePoint", name: "GamePoint", permMin: 10, usage: "人物名称 控制符(+,-,=) 游戏点数(1-100000000)", handler: (*PlayObject).cmdGamePoint})
	registerGMCommand(&gmCmdDef{key: "CreditPoint", name: "CreditPoint", permMin: 10, usage: "人物名称 控制符(+,-,=) 声望点数(0-255)", handler: (*PlayObject).cmdCreditPoint})
	registerGMCommand(&gmCmdDef{key: "Training", name: "Training", permMin: 10, usage: "人物名称 技能名称 修炼等级(0-3)", handler: (*PlayObject).cmdTraining})
	registerGMCommand(&gmCmdDef{key: "DeleteSkill", name: "DeleteSkill", permMin: 10, usage: "人物名称 技能名称", handler: (*PlayObject).cmdDeleteSkill})
	registerGMCommand(&gmCmdDef{key: "OpTraining", name: "TrainingSkill", permMin: 10, usage: "人物名称 技能名称 修炼等级(0-3)", handler: (*PlayObject).cmdTrainingSkill})
	// --- 调试段命令（ObjBase.pas:8253 m_btPermission>4 门槛内）---
	registerGMCommand(&gmCmdDef{key: "DelGameGold", name: "DelGamePoint", permMin: 6, permFail: permFailSilent, usage: "人物名称 游戏点数", handler: (*PlayObject).cmdDelGamePoint})
	registerGMCommand(&gmCmdDef{key: "AddGamePoint", name: "AddGamePoint", permMin: 6, permFail: permFailSilent, usage: "人物名称 金币数量", handler: (*PlayObject).cmdAddGamePoint})
	registerGMCommand(&gmCmdDef{key: "AdjuestTLevel", name: "AdjustLevel", permMin: 6, permFail: permFailSilent, usage: "人物名称 等级", handler: (*PlayObject).cmdAdjustLevel})
	registerGMCommand(&gmCmdDef{key: "AdjuestExp", name: "AdjustExp", permMin: 6, permFail: permFailSilent, usage: "人物名称 经验值", handler: (*PlayObject).cmdAdjustExp})
}

// g_Config.sGameGoldName 默认值（M2Share.pas:1634）；Go 配置无对应条目。
const gmGameGoldName = "元宝"

// cmdLuckyPoint 查询幸运点数（Delphi CmdLuckPoint，分发 ObjBase.pas:7963，
// 实现 13328-13357）。Delphi 原版只实现查看分支，修改分支（sCtr<>''）
// 为空实现，Go 同样只查看。
// 偏差：Go 无身体幸运系统（m_dBodyLuck/m_nBodyLuckLevel），输出 0 占位；
// 幸运值取 BaseObject.Luck（装备幸运合计，itemsystem.go RecalcAbilitys）。
func (p *PlayObject) cmdLuckyPoint(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	// g_sGameCommandLuckPointMsg（M2Share.pas:3310，c_Green）。
	p.sysMsgGreen(server, fmt.Sprintf("%s 的幸运点数为:%d/%g 幸运值为:%d", target.Name, 0, 0.0, target.Luck))
	return true
}

// cmdChangeJob 更改目标职业（Delphi CmdChangeJob，分发 ObjBase.pas:8031，
// 实现 10818-10847）。职业字面量 iconv 解码 sWarrior/sWizard/sTaos
//（M2Share.pas:178-180）为 'Warrior'/'Wizard'/'Taoist'；帮助文案写的
// Warr/Wizard/Taos（g_sGameCommandChangeJobHelpMsg，M2Share.pas:3346），
// Go 两种写法均接受。
// 原版 bug：Delphi 职业名不匹配时仍显示更改成功（ObjBase.pas:10838-10844），
// Go 改为不匹配时输出用法帮助。
func (p *PlayObject) cmdChangeJob(server *netserver.TCPServer, params []string) bool {
	jobUsage := func() {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("ChangeJob", "ChangeJob"), "人物名称 职业类型(Warr Wizard Taos)"))
	}
	if len(params) < 2 {
		jobUsage()
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	switch strings.ToLower(params[1]) {
	case "warrior", "warr":
		target.Job = 0 // jWarr
	case "wizard", "wiz":
		target.Job = 1 // jWizard
	case "taoist", "taos":
		target.Job = 2 // jTaos
	default:
		jobUsage()
		return true
	}
	// Delphi HasLevelUp(1)（ObjBase.pas:10840）：重算能力并刷新客户端。
	gmRefreshAbilities(target, server)
	// 职业不参与外观编码（Feature() 为发型/性别/武器/衣服），
	// 仍对齐 Delphi FeatureChanged 广播语义通知视野刷新。
	target.broadcastFeatureChanged(server)
	// g_sGameCommandChangeJobHumanMsg/g_sGameCommandChangeJobMsg
	//（M2Share.pas:3347-3348，c_Green）。
	target.sysMsgGreen(server, "职业更改成功。")
	p.sysMsgGreen(server, fmt.Sprintf("%s 的职业更改成功。", target.Name))
	return true
}

// cmdChangeGender 更改目标性别（Delphi CmdChangeGender，分发 ObjBase.pas:8036，
// 实现 10749-10790）。性别字面量区分大小写（Delphi `=` 比较）。
func (p *PlayObject) cmdChangeGender(server *netserver.TCPServer, params []string) bool {
	nSex := -1
	if len(params) >= 2 {
		switch params[1] {
		case "Man", "Male", "0":
			nSex = 0
		case "Woman", "Female", "1":
			nSex = 1
		}
	}
	if len(params) < 1 || nSex == -1 {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("ChangeGender", "ChangeGender"), "人物名称 性别(男、女)"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		// 原文（ObjBase.pas:10786）：与通用离线文案不同。
		p.sysMsg(server, params[0]+"没有在线！！！")
		return true
	}
	if int(target.Gender) != nSex {
		target.Gender = byte(nSex)
		target.broadcastFeatureChanged(server) // Delphi FeatureChanged()
		p.sysMsgGreen(server, target.Name+" 的性别已改变。")
	} else {
		p.sysMsg(server, target.Name+" 的性别未改变！！！")
	}
	return true
}

// cmdHair 更改目标发型（Delphi CmdHair，分发 ObjBase.pas:8041，
// 实现 13092-13118）。分发 Str_ToInt(sParam2, 0)：不带类型值按 0。
func (p *PlayObject) cmdHair(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("Hair", "hair"), "人物名称 类型值"))
		return true
	}
	nHair := 0
	if len(params) >= 2 {
		nHair, _ = strconv.Atoi(params[1]) // Delphi Str_ToInt 解析失败按 0
	}
	if nHair < 0 {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("Hair", "hair"), "人物名称 类型值"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	target.Hair = byte(nHair)
	target.broadcastFeatureChanged(server) // Delphi FeatureChanged()
	p.sysMsgGreen(server, params[0]+" 的头发已改变。")
	return true
}

// cmdBonusPoint 查看/设置目标属性点（Delphi CmdBonuPoint，分发 ObjBase.pas:8048，
// 实现 10618-10660）。
// 偏差：Go 无 DC/MC/SC/AC/MAC/HP/MP/HIT/SPEED 分项分配模型（Delphi
// m_BonusAbil），查看输出的已分配明细恒为 0；设置后无 Delphi
// RM_ADJUST_BONUS 内部消息，改发 SMAbility（body 含 BonusPoint，
// playobject.go encodeAbilityBody）。
func (p *PlayObject) cmdBonusPoint(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("BonusPoint", "BonusPoint"), "人物名称 属性点数(不输入为查看点数)"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	nCount := 0
	if len(params) >= 2 {
		nCount, _ = strconv.Atoi(params[1])
	}
	if nCount > 0 {
		target.BonusPoint = nCount
		target.SendAbility(server)
		return true
	}
	// 查看：格式抄自 Delphi（ObjBase.pas:10647-10658，c_Red）。
	p.sysMsg(server, fmt.Sprintf("%s的属性点数为:%s", params[0],
		fmt.Sprintf("未分配点数:%d 已分配点数:(DC:%d MC:%d SC:%d AC:%d MAC:%d HP:%d MP:%d HIT:%d SPEED:%d)",
			target.BonusPoint, 0, 0, 0, 0, 0, 0, 0, 0, 0)))
	return true
}

// cmdDelBonusPoint 清除目标已分配属性点（Delphi CmdDelBonuPoint，
// 分发 ObjBase.pas:8051，实现 11526-11554）。
// 偏差：Go 无 DC/MC/... 分项分配模型，清除 = BonusPoint 清零。
func (p *PlayObject) cmdDelBonusPoint(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("DelBonuPoint", "DelBonusPoint"), "人物名称"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	target.BonusPoint = 0
	target.SendAbility(server)
	// 目标侧 c_Red、GM 侧 c_Green（ObjBase.pas:11549-11550）。
	target.sysMsg(server, "分配点数已清除！！！")
	p.sysMsgGreen(server, params[0]+" 的分配点数已清除.")
	return true
}

// cmdRestBonusPoint 复位目标已分配属性点（Delphi CmdRestBonuPoint，
// 分发 ObjBase.pas:8056，实现 11582-11623）。Delphi 把分项已分配点数
// 退还到未分配池；Go 无分项分配模型，无可退还点数，数值不变仅保留文案。
func (p *PlayObject) cmdRestBonusPoint(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("RestBonuPoint", "RestBonusPoint"), "人物名称"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	// 目标侧 c_Red、GM 侧 c_Green（ObjBase.pas:11618-11619）。
	target.sysMsg(server, "分配点数已复位！！！")
	p.sysMsgGreen(server, params[0]+" 的分配点数已复位.")
	return true
}

// cmdSetPermission 设置目标权限等级（Delphi CmdSetPermission，
// 分发 ObjBase.pas:8061，实现 11718-11748）。
// Delphi Str_ToInt(sPermission, 0)：解析失败按 0（合法值，设为 0）。
func (p *PlayObject) cmdSetPermission(server *netserver.TCPServer, params []string) bool {
	nPermission := 0
	if len(params) >= 2 {
		nPermission, _ = strconv.Atoi(params[1])
	}
	if len(params) < 1 || nPermission < 0 || nPermission > 10 {
		// 原文（ObjBase.pas:11734）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("SetPermission", "SetPermission"), "人物名称 权限等级(0 - 10)"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	target.Permission = byte(nPermission)
	// 原文 c_Red（ObjBase.pas:11747）。
	p.sysMsg(server, params[0]+" 当前权限为: "+strconv.Itoa(int(target.Permission)))
	return true
}

// cmdReNewLevel 查看/设置转生等级（Delphi CmdReNewLevel，分发 ObjBase.pas:8066，
// 实现 11555-11581；实现内硬编码 m_btPermission<6 静默退出，注册 6+silent）。
// 偏差：Delphi 设置后 RefShowName() 按职业加转生名缀（*<神>*/*<尊>*），
// Go 显示名无转生修饰，仅改数值。
func (p *PlayObject) cmdReNewLevel(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		// 原文（ObjBase.pas:11562）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("ReNewLevel", "ReNewLevel"), "人物名称 点数(为空则查看)"))
		return true
	}
	nLevel := -1
	if len(params) >= 2 {
		nLevel, _ = strconv.Atoi(params[1]) // Delphi Str_ToInt(sLevel, -1)
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		// 原文（ObjBase.pas:11578）：与通用离线文案不同。
		p.sysMsg(server, params[0]+" 没在线上！！！")
		return true
	}
	if nLevel >= 0 && nLevel <= 255 {
		target.ReNewLevel = nLevel
	}
	p.sysMsgGreen(server, fmt.Sprintf("%s 的转生等级为 %d", params[0], target.ReNewLevel))
	return true
}

// gmPointCmdCheck 经济三命令（GameGold/GamePoint/CreditPoint）共用的
// 参数校验（Delphi 同构实现 ObjBase.pas:10280-10435）：
// 人物名非空、控制符 ∈ =/+/-（缺省 '1' 触发用法）、数量在 [0, max]。
func gmPointCmdCheck(params []string, max int) (ctr byte, amount int, ok bool) {
	ctr = '1'
	if len(params) >= 2 {
		ctr = params[1][0] // Delphi Ctr := sCtr[1] 取首字符
	}
	if len(params) >= 3 {
		amount, _ = strconv.Atoi(params[2]) // Delphi Str_ToInt(sParam3, 0)
	}
	if len(params) < 1 || (ctr != '=' && ctr != '+' && ctr != '-') || amount < 0 || amount > max {
		return ctr, amount, false
	}
	return ctr, amount, true
}

// cmdGameGold 设置目标游戏币（Delphi CmdGameGold，分发 ObjBase.pas:8081，
// 实现 10280-10331）。输入上限 2 亿；+/- 结果不钳制（原版行为）。
// 原版 bug：'=' 分支误写 m_nGamePoint（ObjBase.pas:10303-10306，
// 复制 CmdGamePoint 遗留），Go 修正为写 GameGold。
// 偏差：Delphi GameGoldChanged()（ObjBase.pas:1655）经 RM_GAMEGOLDCHANGED →
// SendGoldInfo(SM_GAMEGOLDNAME) 同步客户端；Go 协议无游戏币同步消息，
// 数值变更随下次登录刷新。
func (p *PlayObject) cmdGameGold(server *netserver.TCPServer, params []string) bool {
	ctr, nGold, ok := gmPointCmdCheck(params, 200000000)
	if !ok {
		// g_sGameCommandGameGoldHelpMsg（M2Share.pas:3277）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("GameGold", "GameGold"), "人物名称 控制符(+,-,=) 游戏币(1-200000000)"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	switch ctr {
	case '=':
		target.GameGold = nGold
	case '+':
		target.GameGold += nGold
	case '-':
		target.GameGold -= nGold
	}
	// g_sGameCommandGameGoldHumanMsg/GMMsg（M2Share.pas:3278-3279，c_Green）：
	// 原文对 '='/'-' 同样显示"已增加"。
	target.sysMsgGreen(server, fmt.Sprintf("你的%s已增加%d，当前拥有%d%s。", gmGameGoldName, nGold, target.GameGold, gmGameGoldName))
	p.sysMsgGreen(server, fmt.Sprintf("%s的%s已增加%d，当前拥有%d%s。", target.Name, gmGameGoldName, nGold, target.GameGold, gmGameGoldName))
	return true
}

// cmdGamePoint 设置目标游戏点（Delphi CmdGamePoint，分发 ObjBase.pas:8086，
// 实现 10332-10384）。输入上限 1 亿；+/- 结果不钳制（原版行为）。
func (p *PlayObject) cmdGamePoint(server *netserver.TCPServer, params []string) bool {
	ctr, nPoint, ok := gmPointCmdCheck(params, 100000000)
	if !ok {
		// g_sGameCommandGamePointHelpMsg（M2Share.pas:3269）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("GamePoint", "GamePoint"), "人物名称 控制符(+,-,=) 游戏点数(1-100000000)"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	switch ctr {
	case '=':
		target.GamePoint = nPoint
	case '+':
		target.GamePoint += nPoint
	case '-':
		target.GamePoint -= nPoint
	}
	// g_sGameCommandGamePointHumanMsg/GMMsg（M2Share.pas:3270-3271，c_Green）。
	target.sysMsgGreen(server, fmt.Sprintf("你的游戏点已增加%d点，当前总点数为%d点。", nPoint, target.GamePoint))
	p.sysMsgGreen(server, fmt.Sprintf("%s的游戏点已增加%d点，当前总点数为%d点。", target.Name, nPoint, target.GamePoint))
	return true
}

// cmdCreditPoint 设置目标声望（Delphi CmdCreditPoint，分发 ObjBase.pas:8091，
// 实现 10385-10435）。输入范围 0-255（High(Byte)）。
// 原版 bug：+/- 分支校验的是入参 nPoint（已通过校验）而非结果
//（ObjBase.pas:10417-10431），结果可越界；Go 修正为结果钳制 0-255。
func (p *PlayObject) cmdCreditPoint(server *netserver.TCPServer, params []string) bool {
	ctr, nPoint, ok := gmPointCmdCheck(params, 255)
	if !ok {
		// g_sGameCommandCreditPointHelpMsg（M2Share.pas:3273）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("CreditPoint", "CreditPoint"), "人物名称 控制符(+,-,=) 声望点数(0-255)"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	switch ctr {
	case '=':
		target.CreditPoint = nPoint
	case '+':
		target.CreditPoint += nPoint
	case '-':
		target.CreditPoint -= nPoint
	}
	if target.CreditPoint < 0 {
		target.CreditPoint = 0
	}
	if target.CreditPoint > 255 {
		target.CreditPoint = 255
	}
	// g_sGameCommandCreditPointHumanMsg/GMMsg（M2Share.pas:3274-3275，c_Green）：
	// 原文对 '='/'-' 同样显示"已增加"。
	target.sysMsgGreen(server, fmt.Sprintf("你的声望点已增加%d点，当前总声望点数为%d点。", nPoint, target.CreditPoint))
	p.sysMsgGreen(server, fmt.Sprintf("%s的声望点已增加%d点，当前总声望点数为%d点。", target.Name, nPoint, target.CreditPoint))
	return true
}

// cmdTraining 调整目标已有技能修炼等级（Delphi CmdTrainingSkill，
// 分发 ObjBase.pas:8096，实现 10177-10221）。
// 偏差：Delphi 以 RM_MAGIC_LVEXP 通知客户端单技能更新，Go 无该链路，
// 改发 SendMyMagicFull 全量技能列表。
func (p *PlayObject) cmdTraining(server *netserver.TCPServer, params []string) bool {
	nLevel := 0
	if len(params) >= 3 {
		nLevel, _ = strconv.Atoi(params[2]) // Delphi Str_ToInt(sParam3, 0)
	}
	if len(params) < 3 || nLevel <= 0 {
		// 原文（ObjBase.pas:10189）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("Training", "Training"), "人物名称 技能名称 修炼等级(0-3)"))
		return true
	}
	if nLevel > 3 {
		nLevel = 3 // Delphi _MIN(3, nLevel)
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		// 原文（ObjBase.pas:10196）：与通用离线文案不同。
		p.sysMsg(server, fmt.Sprintf("%s不在线，或在其它服务器上！！", params[0]))
		return true
	}
	pm := findLearnedMagicByName(target, params[1])
	if pm == nil {
		return true // Delphi 未找到已学技能时静默退出（ObjBase.pas:10203-10220）
	}
	pm.Level = nLevel
	target.SendMyMagicFull(server)
	// 目标侧/GM 侧文案原文如此（ObjBase.pas:10214-10215，c_Green）。
	target.sysMsgGreen(server, fmt.Sprintf("%s的修改炼等级为%d", params[1], nLevel))
	p.sysMsgGreen(server, fmt.Sprintf("%s的技能%s修炼等级为%d", target.Name, params[1], nLevel))
	return true
}

// cmdDeleteSkill 删除目标技能（Delphi CmdDelSkill，分发 ObjBase.pas:8106，
// 实现 12415-12465）。All = 全部删除（Delphi 全删不发文案）。
// 原版 bug：Delphi 用法文案带多余 ')'（ObjBase.pas:12432），不继承。
// 偏差：Delphi 逐个 SendDelMagic(SM_DELMAGIC)，Go 改发 SendMyMagicFull。
func (p *PlayObject) cmdDeleteSkill(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("DeleteSkill", "DeleteSkill"), "人物名称 技能名称"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	if strings.EqualFold(params[1], "All") {
		target.LearnedMagics = nil
		target.SendMyMagicFull(server)
		return true
	}
	pm := findLearnedMagicByName(target, params[1])
	if pm == nil {
		return true // Delphi 未找到时静默退出（ObjBase.pas:12440-12464）
	}
	target.removeMagic(pm.MagID)
	target.SendMyMagicFull(server)
	// 文案原文（ObjBase.pas:12456-12457，c_Green）。
	target.sysMsgGreen(server, fmt.Sprintf("技能%s已删除。", params[1]))
	p.sysMsgGreen(server, fmt.Sprintf("%s的技能%s已删除。", target.Name, params[1]))
	return true
}

// cmdTrainingSkill 给目标新增技能（Delphi CmdTrainingMagic，
// 分发 ObjBase.pas:8111，实现 10130-10176）。
// 偏差：Delphi SendAddMagic(SM_ADDMAGIC) 单技能通知 + RecalcAbilitys；
// Go 客户端技能列表以 SendMyMagicFull 全量刷新。
func (p *PlayObject) cmdTrainingSkill(server *netserver.TCPServer, params []string) bool {
	nLevel := 0
	if len(params) >= 3 {
		nLevel, _ = strconv.Atoi(params[2])
	}
	if len(params) < 3 || nLevel < 0 || nLevel > 3 {
		// 原文（ObjBase.pas:10143）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("OpTraining", "TrainingSkill"), "人物名称 技能名称 修炼等级(0-3)"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	// Delphi UserEngine.FindMagic 按名字查技能库（不区分大小写）。
	def := lookupMagicByName(p.MagicDB, params[1])
	if def == nil {
		// 原文（ObjBase.pas:10157，c_Red）。
		p.sysMsg(server, fmt.Sprintf("%s 技能名称不正确！！！", params[1]))
		return true
	}
	if target.findMagic(def.MagID) != nil {
		// Delphi IsTrainingSkill 按 wMagIdx 判重（ObjBase.pas:10160-10164，c_Red）。
		p.sysMsg(server, fmt.Sprintf("%s 技能已修炼过了！！！", params[1]))
		return true
	}
	target.learnMagic(def.MagID, nLevel, 0) // Delphi btKey=0、nTranPoint=0
	target.SendMyMagicFull(server)
	target.RecalcAbilitys() // Delphi PlayObject.RecalcAbilitys
	// 原文（ObjBase.pas:10173，c_Green）。
	p.sysMsgGreen(server, fmt.Sprintf("%s 的 %s 技能修炼成功！！！", target.Name, params[1]))
	return true
}

// cmdAddGamePoint 给目标加游戏点（Delphi CmdAddGameGold，分发 ObjBase.pas:8410
// 调试段内，实现 10222-10252；实现内硬编码 m_btPermission<6 静默退出，
// 注册 6+silent）。上限 2000000。
// 偏差：Delphi 操作 m_nGameGold 但文案称"游戏点"（ObjBase.pas:10233-10246），
// Go 统一操作 GamePoint 字段。
func (p *PlayObject) cmdAddGamePoint(server *netserver.TCPServer, params []string) bool {
	nPoint := 0
	if len(params) >= 2 {
		nPoint, _ = strconv.Atoi(params[1])
	}
	if len(params) < 1 || nPoint <= 0 {
		// 原文（ObjBase.pas:10229）：参数说明写的"金币数量"，原文如此。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("AddGamePoint", "AddGamePoint"), "人物名称 金币数量"))
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	// Delphi 钳制到 2000000，文案报告实际增加量（ObjBase.pas:10236-10242）。
	if target.GamePoint+nPoint < 2000000 {
		target.GamePoint += nPoint
	} else {
		nPoint = 2000000 - target.GamePoint
		target.GamePoint = 2000000
	}
	// 文案原文（ObjBase.pas:10244-10245，c_Green）。
	p.sysMsgGreen(server, fmt.Sprintf("%s的游戏点已增加%d.", target.Name, nPoint))
	target.sysMsgGreen(server, fmt.Sprintf("游戏点已增加%d.", nPoint))
	return true
}

// cmdDelGamePoint 扣减目标游戏点（Delphi CmdDelGameGold，分发 ObjBase.pas:8405
// 调试段内，实现 10253-10279；实现内硬编码 m_btPermission<6 静默退出，
// 注册 6+silent）。
// 偏差：字段同 cmdAddGamePoint 统一为 GamePoint；Delphi 参数非法静默退出，
// Go 保留静默（usage 仅注册表 "?" 触发）。
func (p *PlayObject) cmdDelGamePoint(server *netserver.TCPServer, params []string) bool {
	nPoint := 0
	if len(params) >= 2 {
		nPoint, _ = strconv.Atoi(params[1])
	}
	if len(params) < 1 || nPoint <= 0 {
		return true // Delphi 静默退出（ObjBase.pas:10259）
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	// Delphi 最多扣到 0，文案报告实际扣减量（ObjBase.pas:10266-10272）。
	if target.GamePoint > nPoint {
		target.GamePoint -= nPoint
	} else {
		nPoint = target.GamePoint
		target.GamePoint = 0
	}
	// 文案原文（ObjBase.pas:10274-10275，c_Green）。
	p.sysMsgGreen(server, fmt.Sprintf("%s的游戏点已减少%d.", target.Name, nPoint))
	target.sysMsgGreen(server, fmt.Sprintf("游戏点已减少%d.", nPoint))
	return true
}

// cmdAdjustLevel 调整目标等级（Delphi CmdAdjuestLevel，分发 ObjBase.pas:8480
// 调试段内，实现 10568-10598）。上限 MAXUPLEVEL=65535（M2Share.pas:91）。
func (p *PlayObject) cmdAdjustLevel(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		// 原文（ObjBase.pas:10580）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("AdjuestTLevel", "AdjustLevel"), "人物名称 等级"))
		return true
	}
	nLevel := 1 // Delphi Str_ToInt(sParam2, 1)
	if len(params) >= 2 {
		nLevel, _ = strconv.Atoi(params[1])
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	// Delphi _MAX(1, _MIN(MAXUPLEVEL, nLevel))（ObjBase.pas:10587）。
	if nLevel < 1 {
		nLevel = 1
	}
	if nLevel > 65535 {
		nLevel = 65535
	}
	target.WAbil.Level = uint16(nLevel)
	// Delphi HasLevelUp(1)（ObjBase.pas:10588）：重算能力并通知客户端。
	gmRefreshAbilities(target, server)
	p.sysMsgGreen(server, params[0]+" 等级调整完成。")
	return true
}

// cmdAdjustExp 调整目标经验（Delphi CmdAdjuestExp，分发 ObjBase.pas:8485
// 调试段内，实现 10535-10567）。Delphi 直接写 m_Abil.Exp 后 HasLevelUp(1)
// 检查升级；Go 写入后走与 addExp 相同的升级循环（playobject.go:1926-1948）。
func (p *PlayObject) cmdAdjustExp(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		// 原文（ObjBase.pas:10548）。
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("AdjuestExp", "AdjustExp"), "人物名称 经验值"))
		return true
	}
	nExp := 0 // Delphi Str_ToInt(sExp, 0)
	if len(params) >= 2 {
		nExp, _ = strconv.Atoi(params[1])
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	if nExp < 0 {
		nExp = 0
	}
	target.WAbil.Exp = uint32(nExp)
	gmCheckLevelUps(target, server)
	p.sysMsgGreen(server, params[0]+" 经验调整完成。")
	return true
}

// gmRefreshAbilities Delphi HasLevelUp(1) 的客户端刷新等效（ObjBase.pas:
// 1943-1950：RecalcLevelAbilitys + RecalcAbilitys + RM_LEVELUP → SM_LEVELUP）。
func gmRefreshAbilities(target *PlayObject, server *netserver.TCPServer) {
	target.RecalcAbilitys()
	target.WAbil.HP = target.WAbil.MaxHP
	target.WAbil.MP = target.WAbil.MaxMP
	levelMsg := protocol.MakeDefaultMsg(protocol.SMLevelUp, int32(target.WAbil.Exp), uint16(target.WAbil.Level), 0, 0)
	server.Send(target.Session.ID, levelMsg, "")
	target.SendAbility(server)
	target.sendHealthSpell(server)
}

// gmCheckLevelUps 绝对经验写入后的循环升级检查（Delphi HasLevelUp，
// ObjBase.pas:1943；单级逻辑复刻 Go addExp 升级段 playobject.go:1928-1943：
// 扣经验、等级+1、重算属性、满血蓝、奖励属性点）。
func gmCheckLevelUps(target *PlayObject, server *netserver.TCPServer) {
	leveled := false
	for target.WAbil.Exp >= target.GetMaxExp() {
		target.WAbil.Exp -= target.GetMaxExp()
		target.WAbil.Level++
		target.RecalcAbilitys()
		target.WAbil.HP = target.WAbil.MaxHP
		target.WAbil.MP = target.WAbil.MaxMP
		target.BonusPoint += target.Engine.Config.GetBonusPerLevel()
		leveled = true
	}
	if leveled {
		levelMsg := protocol.MakeDefaultMsg(protocol.SMLevelUp, int32(target.WAbil.Exp), uint16(target.WAbil.Level), 0, 0)
		server.Send(target.Session.ID, levelMsg, "")
		target.sendWeightChanged(server)
	}
	target.SendAbility(server)
	target.sendHealthSpell(server)
}

// lookupMagicByName 技能名查库（Delphi UserEngine.FindMagic 按名字
// CompareText；Go MagicDB.GetByName 精确匹配，补 EqualFold 兜底）。
func lookupMagicByName(db *MagicDB, name string) *MagicDef {
	if db == nil {
		return nil
	}
	if def := db.GetByName(name); def != nil {
		return def
	}
	for i := range db.Magics {
		if strings.EqualFold(db.Magics[i].MagName, name) {
			return &db.Magics[i]
		}
	}
	return nil
}

// findLearnedMagicByName 在已学技能中按名字查找（Delphi
// CompareText(UserMagic.MagicInfo.sMagicName, sSkillName)，
// ObjBase.pas:10205-10219）。
func findLearnedMagicByName(p *PlayObject, name string) *PlayerMagic {
	db := p.MagicDB
	if db == nil {
		return nil
	}
	for _, pm := range p.LearnedMagics {
		if def := db.GetByID(pm.MagID); def != nil && strings.EqualFold(def.MagName, name) {
			return pm
		}
	}
	return nil
}
