package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// GM 命令——NPC/怪物/地图模式族。
// Delphi 参照：ObjBase.pas（分发 ProcessUserLineMsg 7188-8600）+
// M2Share.pas g_GameCommand 命令表（2724-2926）与文案常量（3304-3306、3354-3355）。
// 注册顺序 = Delphi 分发链出现顺序（ObjBase.pas:7943-8582）。

func init() {
	// --- NPC/怪物攻城段（ObjBase.pas:7943-7955、8189-8206）---
	registerGMCommand(&gmCmdDef{key: "MobNpc", name: "MobNpc", permMin: 10, usage: gmMsgMobNpcUsage, handler: (*PlayObject).cmdMobNpc})
	registerGMCommand(&gmCmdDef{key: "NpcScript", name: "NpcScript", permMin: 10, usage: "？？？？", handler: (*PlayObject).cmdNpcScript})
	registerGMCommand(&gmCmdDef{key: "DelNpc", name: "DelNpc", permMin: 10, handler: (*PlayObject).cmdDelNpc})
	registerGMCommand(&gmCmdDef{key: "Mission", name: "Mission", permMin: 10, usage: "X  Y", handler: (*PlayObject).cmdMission})
	registerGMCommand(&gmCmdDef{key: "MobPlace", name: "MobPlace", permMin: 10, usage: "X  Y 怪物名称 怪物数量", handler: (*PlayObject).cmdMobPlace})
	// 地图模式：g_GameCommand 表权限 10（M2Share.pas:2915-2916），但实现内
	// m_btPermission<6 硬编码静默退出（ObjBase.pas:11867/11889）→ 注册 perm 6 + permFailSilent。
	registerGMCommand(&gmCmdDef{key: "SetMapMode", name: "SetMapMode", permMin: 6, permFail: permFailSilent, usage: "地图号 模式", handler: (*PlayObject).cmdSetMapMode})
	registerGMCommand(&gmCmdDef{key: "ShowMapMode", name: "ShowMapMode", permMin: 6, permFail: permFailSilent, usage: "地图号", handler: (*PlayObject).cmdShowMapMode})
	// --- 调试段（ObjBase.pas:8387 起 >=5||boTestServer 门槛内：8390/8395/8579）---
	// FireBurn 实现内另有 m_btPermission<6 静默退出（ObjBase.pas:12868）→ perm 6 + permFailSilent。
	registerGMCommand(&gmCmdDef{key: "FireBurn", name: "FireBurn", permMin: 6, permFail: permFailSilent, usage: "nInt nTime nN", handler: (*PlayObject).cmdFireBurn})
	registerGMCommand(&gmCmdDef{key: "TestFire", name: "TestFire", permMin: 6, permFail: permFailSilent, usage: "范围 类型 时间 点数", handler: (*PlayObject).cmdTestFire})
	registerGMCommand(&gmCmdDef{key: "MobFireBurn", name: "MobFireBurn", permMin: 6, permFail: permFailSilent, handler: (*PlayObject).cmdMobFireBurn})
}

// gmMsgMobNpcUsage @MobNpc 参数帮助（g_sGameCommandMobNpcHelpMsg，M2Share.pas:3304）。
const gmMsgMobNpcUsage = "NPC名称 脚本文件名 外形(数字) 属沙城(0,1)"

// getFacingNpc 返回面前一格的 NPC（Delphi GetPoseCreate：GetFrontPosition +
// GetMovingObject，ObjBase.pas:15875-15884；仅查面前格，不做双方面对面校验）。
func (p *PlayObject) getFacingNpc() *NpcObject {
	if p.envir == nil {
		return nil
	}
	dx, dy := dirToOffset(p.Dir)
	if npc, ok := p.envir.GetMovingObject(p.CurrX+dx, p.CurrY+dy).(*NpcObject); ok && !npc.Ghost {
		return npc
	}
	return nil
}

// resolveNpcScriptPath 把 @MobNpc 的"脚本名"参数解析为 Go 可读的脚本路径。
// Delphi 存相对名，读取时拼接 EnvirDir + 'Market_Def/' + m_sScript + '-' +
// m_sMapName + '.txt'（ObjBase.pas:13911-13916）；Go 脚本平铺在
// npcConfigDir 的 merchant_scripts/npc_scripts 下（mongen.go LoadNpcs 命名
// <ID>-<地图名>.txt）。优先原路径，再试两个脚本目录的 <名>-<地图>.txt 与
// <名>.txt；都找不到时原样返回（GetScript 加载失败静默无脚本，等同 Delphi
// 脚本文件不存在）。
func resolveNpcScriptPath(e *UserEngine, name, mapName string) string {
	if _, err := os.Stat(name); err == nil {
		return name
	}
	if dir := e.npcConfigDir; dir != "" {
		for _, sub := range []string{"merchant_scripts", "npc_scripts"} {
			for _, c := range []string{
				filepath.Join(dir, sub, name+"-"+mapName+".txt"),
				filepath.Join(dir, sub, name+".txt"),
			} {
				if _, err := os.Stat(c); err == nil {
					return c
				}
			}
		}
	}
	return name
}

// cmdMobNpc 在面前生成 NPC（Delphi CmdMobNpc，分发 ObjBase.pas:7943-7945，
// 实现 13777-13815）。Delphi 创建 TMerchant（Name/MapName/PEnvir/Appr/
// m_nFlag=0/Castle/Script），置于 GetFrontPosition 后 Initialize +
// UserEngine.AddMerchant；成功无文案（静默）。
// Go 运行时创建等价闭环：NewNpcObject 设 Name/Appr/Castle/Script，
// 插入 Engine.Npcs 与当前地图格子，再按脚本初始化商人能力/商品
//（InitGoodsFromScript ≡ Merchant.Initialize 的脚本部分）。
func (p *PlayObject) cmdMobNpc(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		p.gmParamHelp(server, "MobNpc", gmMsgMobNpcUsage)
		return true
	}
	if p.envir == nil || p.Engine == nil {
		return true
	}
	// Delphi Str_ToInt(sParam3, 0) / Str_ToInt(sParam4, 0)=1（ObjBase.pas:13795-13796）。
	appr, _ := strconv.Atoi(params[2])
	isCastle := false
	if len(params) >= 4 {
		if v, err := strconv.Atoi(params[3]); err == nil && v == 1 {
			isCastle = true
		}
	}
	e := p.Engine
	e.mu.Lock()
	id := e.nextMonsterID
	e.nextMonsterID++
	e.mu.Unlock()

	npc := NewNpcObject(params[0], id, uint16(appr))
	npc.MapName = p.MapName
	npc.Castle = isCastle
	npc.IsMerchant = true // Delphi TMerchant.Create（ObjBase.pas:13802）
	npc.envir = p.envir
	npc.Script = resolveNpcScriptPath(e, params[1], p.MapName)
	// Delphi GetFrontPosition（ObjBase.pas:13810-13812）；朝向取面向 GM。
	dx, dy := dirToOffset(p.Dir)
	npc.CurrX, npc.CurrY = p.CurrX+dx, p.CurrY+dy
	npc.Dir = (p.Dir + 4) % 8

	if !p.envir.AddObject(npc.CurrX, npc.CurrY, OS_MOVINGOBJECT, npc) {
		return true // 越界格：Delphi 无此检查，Go 静默放弃（无文案）
	}
	e.mu.Lock()
	e.Npcs = append(e.Npcs, npc)
	e.mu.Unlock()
	// Merchant.Initialize 的脚本部分：加载脚本并应用价格比率/能力/商品。
	if s := npc.GetScript(); s != nil {
		npc.InitGoodsFromScript(s, e.ItemDB, e)
	}
	// 近旁玩家即时可见（与 SpawnMonsterByName 同款 RM_TURN 登场广播）。
	npc.SendRefMsg(RM_TURN, npc.Dir, npc.CurrX, npc.CurrY, npc.Name)
	return true
}

// cmdDelNpc 删除面对的 NPC（Delphi CmdDelNpc，分发 ObjBase.pas:7953-7955，
// 实现 12368-12414）。成功文案 sDelOK='删除NPC成功...'（ObjBase.pas:12373）。
// Delphi 置 Ghost 后留待清理循环；Go 无 NPC ghost 清理，直接从
// Engine.Npcs 与地图格子移除。
func (p *PlayObject) cmdDelNpc(server *netserver.TCPServer, params []string) bool {
	// Delphi '?' 帮助用空参数格式（ObjBase.pas:12380-12382）；注册表 usage
	// 为空不拦截，由此处按原文输出。
	if len(params) > 0 && strings.HasPrefix(params[0], "?") {
		p.sysMsg(server, fmt.Sprintf(gmMsgParamUnKnow, gmNameByKey("DelNpc", "DelNpc"), ""))
		return true
	}
	if p.envir == nil || p.Engine == nil {
		return true
	}
	npc := p.getFacingNpc()
	if npc == nil {
		// g_sGameCommandDelNpcMsg（M2Share.pas:3306）。
		p.sysMsg(server, "命令使用方法不正确，必须与NPC面对面，才能使用此命令！！！")
		return true
	}
	// Delphi: Ghost + RM_DISAPPEAR（ObjBase.pas:12392-12394/12403-12405）。
	npc.SendRefMsg(RM_DISAPPEAR, 0, 0, 0, "")
	npc.Ghost = true
	p.envir.RemoveObject(npc.CurrX, npc.CurrY, OS_MOVINGOBJECT, npc)
	e := p.Engine
	e.mu.Lock()
	for i, n := range e.Npcs {
		if n == npc {
			e.Npcs = append(e.Npcs[:i], e.Npcs[i+1:]...)
			break
		}
	}
	e.mu.Unlock()
	p.sysMsg(server, "删除NPC成功...")
	return true
}

// cmdNpcScript 打印面对 NPC 的脚本内容（Delphi CmdNpcScript，分发
// ObjBase.pas:7948-7950，实现 13857-13935）。帮助文案
// g_sGameCommandNpcScriptHelpMsg='？？？？'（M2Share.pas:3305）。
// 原版 bug：实现体位于 `if sParam1 = ''` 分支内，而空参在前面帮助分支已
// Exit（ObjBase.pas:13873-13877 vs 13906），原版恒只出帮助；Go 按其死代码
// 的意图实装：读取 npc.Script 文件逐行输出 '%d,%s'（空格替换为逗号）。
func (p *PlayObject) cmdNpcScript(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.gmParamHelp(server, "NpcScript", "？？？？")
		return true
	}
	if p.envir == nil {
		return true
	}
	npc := p.getFacingNpc()
	if npc == nil {
		// 原文与 g_sGameCommandDelNpcMsg 相同（ObjBase.pas:13902）。
		p.sysMsg(server, "命令使用方法不正确，必须与NPC面对面，才能使用此命令！！！")
		return true
	}
	if npc.Script == "" {
		return true // Delphi FileExists=False 分支静默（ObjBase.pas:13918）
	}
	data, err := os.ReadFile(npc.Script)
	if err != nil {
		// '读取脚本文件错误: '（ObjBase.pas:13924）。
		p.sysMsg(server, "读取脚本文件错误: "+npc.Script)
		return true
	}
	// Delphi 全量逐行 SysMsg（c_Blue，ObjBase.pas:13926-13931）；
	// Go 聊天窗口保护：至多前 100 行 + 省略提示。
	lines := npcScriptLines(string(data), 100)
	for _, l := range lines {
		p.sysMsgBlue(server, l)
	}
	return true
}

// npcScriptLines 按 Delphi 格式整理脚本行：Trim 后空格替换为逗号，
// 输出 '%d,%s'（行号 0 起，ObjBase.pas:13928-13930）；超过 max 行时
// 追加省略提示。
func npcScriptLines(content string, max int) []string {
	// Delphi TStringList 不为结尾换行产生空行（ObjBase.pas:13922 LoadFromFile）。
	content = strings.TrimSuffix(strings.TrimSuffix(content, "\n"), "\r")
	raw := strings.Split(content, "\n")
	out := make([]string, 0, len(raw))
	for i, line := range raw {
		if i >= max {
			out = append(out, fmt.Sprintf("...(脚本共 %d 行，仅显示前 %d 行)", len(raw), max))
			break
		}
		s := strings.ReplaceAll(strings.TrimSpace(line), " ", ",")
		out = append(out, fmt.Sprintf("%d,%s", i, s))
	}
	return out
}

// cmdMission 设置怪物攻城集中点为当前地图 (X,Y)（Delphi CmdMission，
// 分发 ObjBase.pas:8189-8191，实现 13593-13614：写 g_boMission/
// g_sMissionMap/g_nMissionX/g_nMissionY）。Go 存于 UserEngine.MissionMap/
// MissionX/MissionY（受 e.mu 保护）；MissionMap 非空 ≡ g_boMission=True。
func (p *PlayObject) cmdMission(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		// Delphi 内联帮助 '命令格式: @Mission X  Y'（ObjBase.pas:13604）。
		p.gmParamHelp(server, "Mission", "X  Y")
		return true
	}
	// Delphi Str_ToInt(sX/sY, 0)（ObjBase.pas:13607-13608）。
	x, _ := strconv.Atoi(params[0])
	y, _ := strconv.Atoi(params[1])
	e := p.Engine
	e.mu.Lock()
	e.MissionMap = p.MapName
	e.MissionX, e.MissionY = x, y
	e.mu.Unlock()
	// '怪物集中目标已设定为: %s(%d:%d)'（c_Green，ObjBase.pas:13613）。
	p.sysMsgGreen(server, fmt.Sprintf("怪物集中目标已设定为: %s(%d:%d)", p.MapName, x, y))
	return true
}

// cmdMobPlace 在怪物攻城集中点地图刷怪（Delphi CmdMobPlace，分发
// ObjBase.pas:8194-8196，实现 13816-13856）。
// Delphi 语义：sX/sY 是集中点地图上的绝对刷怪坐标（非相对集中点偏移），
// RegenMonsterByName(g_sMissionMap, X, Y) 后设 m_boMission 让怪物行军到
// 集中点（ObjBase.pas:13846-13851）。Go 未实装怪物攻城 AI，改为直接散布
// 在集中点周围小范围（"行军集中"行为缺失）；X/Y 仍按 Delphi 做形式校验。
func (p *PlayObject) cmdMobPlace(server *netserver.TCPServer, params []string) bool {
	if p.Engine == nil {
		return true
	}
	if len(params) < 4 {
		// Delphi 内联帮助 '命令格式: @MobPlace X  Y 怪物名称 怪物数量'（ObjBase.pas:13834）。
		p.gmParamHelp(server, "MobPlace", "X  Y 怪物名称 怪物数量")
		return true
	}
	// Delphi _MIN(500, Str_ToInt(sCount, 0))（ObjBase.pas:13828）。
	count, _ := strconv.Atoi(params[3])
	if count > 500 {
		count = 500
	}
	x, _ := strconv.Atoi(params[0])
	y, _ := strconv.Atoi(params[1])
	if x <= 0 || y <= 0 || params[2] == "" || count <= 0 {
		p.gmParamHelp(server, "MobPlace", "X  Y 怪物名称 怪物数量")
		return true
	}
	e := p.Engine
	e.mu.RLock()
	missionMap, mx, my := e.MissionMap, e.MissionX, e.MissionY
	e.mu.RUnlock()
	env := p.MapMgr.FindMap(missionMap)
	if missionMap == "" || env == nil {
		// '还没有设定怪物集中点！！！' + '请先用命令Mission设置怪物的集中点。'
		//（ObjBase.pas:13839-13840，Mission 字面量取注册表解析名）。
		p.sysMsg(server, "还没有设定怪物集中点！！！")
		p.sysMsg(server, "请先用命令"+gmNameByKey("Mission", "Mission")+"设置怪物的集中点。")
		return true
	}
	now := time.Now().UnixMilli()
	for i := 0; i < count; i++ {
		sx, sy, ok := findWalkableNear(env, mx, my, 8)
		if !ok {
			sx, sy = mx, my
		}
		if e.SpawnMonsterByName(missionMap, sx, sy, params[2], now) == nil {
			break // Delphi: mon = nil then Break（ObjBase.pas:13852）
		}
	}
	// '%d 只 %s 已正在往地图 %s %d:%d 集中。'（c_Green，数量按 Delphi 恒用
	// 请求数 nCount，ObjBase.pas:13854）。
	p.sysMsgGreen(server, fmt.Sprintf("%d 只 %s 已正在往地图 %s %d:%d 集中。", count, params[2], missionMap, mx, my))
	return true
}

// getEnvirInfo 输出地图已置位标志列表（Delphi TEnvirnoment.GetEnvirInfo，
// Envir.pas:1499-1550）。Delphi 以 BoolToStr 全量输出所有标志；Go 简化为
// 仅列已置位标志（带值标志沿用 Delphi NAME(值) 形态），字段按 MapFlag 枚举。
// Go 无 PK 奖惩/GameGold/GamePoint/RUN* 对应字段，不输出（同 SetMapMode 不支持项）。
func getEnvirInfo(env *Environment) string {
	f := env.Flag
	var parts []string
	add := func(s string) { parts = append(parts, s) }
	// 顺序对齐 Delphi Format 串（Envir.pas:1499-1512）。
	if f.DayLight {
		add("DAY")
	}
	if f.Dark {
		add("DARK")
	}
	if f.Safe {
		add("SAFE")
	}
	if f.Fight {
		add("FIGHT")
	}
	if f.Fight3 {
		add("FIGHT3")
	}
	if f.Quiz {
		add("QUIZ")
	}
	if f.NoReconnect {
		add(fmt.Sprintf("NORECONNECT(%s)", f.ReconnectMap))
	}
	if f.MusicID >= 0 {
		add(fmt.Sprintf("MUSIC(%d)", f.MusicID))
	}
	if f.ExpRate != 0 {
		add(fmt.Sprintf("EXPRATE(%d)", f.ExpRate))
	}
	if f.DecHPPoint > 0 || f.DecHPTime > 0 {
		add(fmt.Sprintf("DECHP(%d/%d)", f.DecHPPoint, f.DecHPTime))
	}
	if f.NeedHole {
		add("NEEDHOLE")
	}
	if f.NoRecall {
		add("NORECALL")
	}
	if f.NoGuildRecall {
		add("NOGUILDRECALL")
	}
	if f.NoRandomMove {
		add("NORANDOMMOVE")
	}
	if f.NoPositionMove {
		add("NOPOSITIONMOVE")
	}
	if f.NoDrug {
		add("NODRUG")
	}
	if f.Mine {
		add("MINE")
	}
	if f.Mine2 {
		add("MINE2")
	}
	if f.NoDropItem {
		add("NODROPITEM")
	}
	if f.NoThrowItem {
		add("NOTHROWITEM")
	}
	if f.NoHorse {
		add("NOHORSE")
	}
	if f.NoChat {
		add("NOCHAT")
	}
	return strings.Join(parts, " ")
}

// cmdShowMapMode 显示地图模式（Delphi CmdShowMapMode，分发
// ObjBase.pas:8204-8206，实现 11862-11881）。'地图模式: ' 蓝字（11879-11880）。
func (p *PlayObject) cmdShowMapMode(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		// Delphi 内联帮助 '命令格式: @ShowMapMode 地图号'（ObjBase.pas:11870）。
		p.gmParamHelp(server, "ShowMapMode", "地图号")
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	if env == nil {
		// '%s 不存在！！！'（ObjBase.pas:11876）。
		p.sysMsg(server, params[0]+" 不存在！！！")
		return true
	}
	p.sysMsgBlue(server, "地图模式: "+getEnvirInfo(env))
	return true
}

// cmdSetMapMode 运行时修改地图标志（Delphi CmdSetMapMode，分发
// ObjBase.pas:8199-8201，实现 11883-12253）。
// 语义：带值参数=置位（部分模式吃参数值），无参数=清除（非翻转）。
// Go MapFlag 无对应字段的 PK 奖惩/元宝/跑动模式（Delphi 11994-12105）
// 输出 '暂不支持: 模式名'；未识别模式按 Delphi if 链落空静默（12251 前落空）。
func (p *PlayObject) cmdSetMapMode(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		// Delphi 内联帮助 '命令格式: @SetMapMode 地图号 模式'（ObjBase.pas:11892）。
		p.gmParamHelp(server, "SetMapMode", "地图号 模式")
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	if env == nil {
		p.sysMsg(server, params[0]+" 不存在！！！")
		return true
	}
	val, val2 := "", ""
	if len(params) >= 3 {
		val = params[2]
	}
	if len(params) >= 4 {
		val2 = params[3]
	}
	on := val != ""
	// CompareText 不敏感（ObjBase.pas:11901 起整条 if 链）。
	switch strings.ToUpper(params[1]) {
	case "SAFE": // ObjBase.pas:11901
		env.Flag.Safe = on
	case "DARK": // 11911
		env.Flag.Dark = on
	case "FIGHT": // 11921
		env.Flag.Fight = on
	case "FIGHT3": // 11931
		env.Flag.Fight3 = on
	case "DAY": // 11941
		env.Flag.DayLight = on
	case "QUIZ": // 11951
		env.Flag.Quiz = on
	case "NORECONNECT": // 11961-11966：带目标地图名
		env.Flag.NoReconnect = on
		if on {
			env.Flag.ReconnectMap = val
		}
	case "MUSIC": // 11972：Str_ToInt(sParam1, -1)
		if on {
			env.Flag.MusicID = atoiOr(val, -1)
		} else {
			env.Flag.MusicID = -1 // Go 以 -1 表示无音乐（mapmanager.go parseMapFlag）
		}
	case "EXPRATE": // 11983
		if on {
			env.Flag.ExpRate = atoiOr(val, -1)
		} else {
			env.Flag.ExpRate = 0
		}
	case "DECHP": // 12038-12044：sParam1=时间 sParam2=点数，缺一即清除
		if val != "" && val2 != "" {
			env.Flag.DecHPTime = atoiOr(val, -1)
			env.Flag.DecHPPoint = atoiOr(val2, -1)
		} else {
			env.Flag.DecHPTime, env.Flag.DecHPPoint = 0, 0
		}
	case "NEEDHOLE": // 12106
		env.Flag.NeedHole = on
	case "NORECALL": // 12116
		env.Flag.NoRecall = on
	case "NORANDOMMOVE": // 12161（NOGUILDRECALL/NODEARRECALL/NOMASTERRECALL
		env.Flag.NoRandomMove = on // 在 ObjBase.pas:12127-12159 被注释，不实装）
	case "NODRUG": // 12171
		env.Flag.NoDrug = on
	case "MINE": // 12181
		env.Flag.Mine = on
	case "MINE2": // 12191
		env.Flag.Mine2 = on
	case "NOTHROWITEM": // 12201
		env.Flag.NoThrowItem = on
	case "NODROPITEM": // 12211
		env.Flag.NoDropItem = on
	case "NOPOSITIONMOVE": // 12221
		env.Flag.NoPositionMove = on
	case "NOHORSE": // 12231
		env.Flag.NoHorse = on
	case "NOCHAT": // 12241
		env.Flag.NoChat = on
	case "PKWINLEVEL", "PKWINEXP", "PKLOSTLEVEL", "PKLOSTEXP",
		"DECGAMEGOLD", "INCGAMEGOLD", "INCGAMEPOINT", "RUNHUMAN", "RUNMON":
		// 杀怪奖惩/元宝/跑动参数（ObjBase.pas:11994-12105），Go MapFlag 无对应字段。
		p.sysMsg(server, "暂不支持: "+strings.ToUpper(params[1]))
	}
	// 'Environment: ' + GetEnvirInfo（c_Blue，ObjBase.pas:12251-12252，
	// 无论模式是否识别都输出）。
	p.sysMsgBlue(server, "Environment: "+getEnvirInfo(env))
	return true
}

// cmdFireBurn 身边放火事件（Delphi CmdFireBurn，分发 ObjBase.pas:8390-8393
// 调试段，实现 12864-12876：TFireBurnEvent.Create(Self, X, Y, nInt, nTime, nN)）。
// Go 复用魔法火墙同款地图火事件（envir.AddFireEvent，mapevent.go；
// 火墙施法 magicsystem.go case 22）。
// 偏差：① Delphi TFireBurnEvent 每 3 秒对格内全部对象按 IsProperTarget
// 结算 RM_MAGSTRUCK_MINE（Event.pas:239-264），Go 火事件每 tick 对格内
// 对象扣血、玩家钳制 1HP（mapevent.go processFireDamage）；② nInt 事件
// 类型在 Go 固定为 ETFire（火渲染）；③ nTime 单位毫秒（Delphi 原样传入，
// TestFire/MobFireBurn 则是秒×1000）。
func (p *PlayObject) cmdFireBurn(server *netserver.TCPServer, params []string) bool {
	nInt, nTime, nN := 0, 0, 0
	if len(params) >= 1 {
		nInt, _ = strconv.Atoi(params[0])
	}
	if len(params) >= 2 {
		nTime, _ = strconv.Atoi(params[1])
	}
	if len(params) >= 3 {
		nN, _ = strconv.Atoi(params[2])
	}
	if nInt == 0 || nTime == 0 || nN == 0 {
		// '命令格式: @FireBurn nInt nTime nN'（ObjBase.pas:12871）。
		p.gmParamHelp(server, "FireBurn", "nInt nTime nN")
		return true
	}
	if p.envir == nil {
		return true
	}
	p.envir.AddFireEvent(server, p.CurrX, p.CurrY, nN, int64(nTime), p.ID)
	return true
}

// cmdTestFire 周围一圈火墙（Delphi CmdTestFire，分发 ObjBase.pas:8395-8398
// 调试段，实现 15080-15104）。nRange 方形环（边长 2n+1，仅边缘格），
// 每格 TFireBurnEvent.Create(Self, nX, nY, nType, nTime*1000, nPoint)。
// 实现内无参数校验/权限检查（门槛在调试段分发处）；nType 在 Go 固定 ETFire。
func (p *PlayObject) cmdTestFire(server *netserver.TCPServer, params []string) bool {
	nRange, nTime, nPoint := 0, 0, 0
	if len(params) >= 1 {
		nRange, _ = strconv.Atoi(params[0])
	}
	if len(params) >= 3 {
		nTime, _ = strconv.Atoi(params[2])
	}
	if len(params) >= 4 {
		nPoint, _ = strconv.Atoi(params[3])
	}
	if p.envir == nil {
		return true
	}
	durationMs := int64(nTime) * 1000
	for dx := -nRange; dx <= nRange; dx++ {
		for dy := -nRange; dy <= nRange; dy++ {
			// Delphi 环形条件（ObjBase.pas:15094-15096）：
			// ((nX<nMaxX) and (nY=nMinY)) or ((nY<nMaxY) and (nX=nMinX))
			// or (nX=nMaxX) or (nY=nMaxY)。
			if !((dx < nRange && dy == -nRange) ||
				(dy < nRange && dx == -nRange) ||
				dx == nRange || dy == nRange) {
				continue
			}
			p.envir.AddFireEvent(server, p.CurrX+dx, p.CurrY+dy, nPoint, durationMs, p.ID)
		}
	}
	return true
}

// cmdMobFireBurn 指定地图放火（Delphi CmdMobFireBurn，分发
// ObjBase.pas:8579-8582 调试段，实现 13707-13750：临时切换 m_PEnvir 后
// TFireBurnEvent.Create(Self, nX, nY, nType, nTime*1000, nPoint)）。
// Go 直接对目标地图 AddFireEvent。nPoint 解析失败/负数按 1
//（ObjBase.pas:13731）；nTime 单位秒。
func (p *PlayObject) cmdMobFireBurn(server *netserver.TCPServer, params []string) bool {
	name := gmNameByKey("MobFireBurn", "MobFireBurn")
	// g_sGameCommandMobFireBurnHelpMsg='命令格式: %s %s %s %s %s %s %s'
	//（M2Share.pas:3354，参数原样回填）。
	help := func() {
		args := [7]string{name}
		for i := 0; i < len(params) && i < 6; i++ {
			args[i+1] = params[i]
		}
		p.sysMsg(server, fmt.Sprintf("命令格式: %s %s %s %s %s %s %s",
			args[0], args[1], args[2], args[3], args[4], args[5], args[6]))
	}
	if len(params) < 1 || strings.HasPrefix(params[0], "?") {
		help()
		return true
	}
	// Delphi Str_ToInt(s, -1)（ObjBase.pas:13726-13730）。
	atoiDef := func(s string) int {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
		return -1
	}
	nX, nY, nType, nTime, nPoint := -1, -1, -1, -1, -1
	if len(params) >= 2 {
		nX = atoiDef(params[1])
	}
	if len(params) >= 3 {
		nY = atoiDef(params[2])
	}
	if len(params) >= 4 {
		nType = atoiDef(params[3])
	}
	if len(params) >= 5 {
		nTime = atoiDef(params[4])
	}
	if len(params) >= 6 {
		nPoint = atoiDef(params[5])
	}
	if nPoint < 0 {
		nPoint = 1
	}
	if nX < 0 || nY < 0 || nType < 0 || nTime < 0 || nPoint < 0 {
		help()
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	if env == nil {
		// g_sGameCommandMobFireBurnMapNotFountMsg='地图%s 不存在'（M2Share.pas:3355）。
		p.sysMsg(server, fmt.Sprintf("地图%s 不存在", params[0]))
		return true
	}
	env.AddFireEvent(server, nX, nY, nPoint, int64(nTime)*1000, p.ID)
	return true
}
