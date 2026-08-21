package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// GM 命令——传送与人员管理族。
// Delphi 参照：ObjBase.pas（分发 ProcessUserLineMsg 7188-8600）+
// M2Share.pas g_GameCommand 命令表（2724-2926）与文案常量（3250-3350）。
// 注册顺序 = Delphi 分发链出现顺序（ObjBase.pas:7811-8023）。

func init() {
	registerGMCommand(&gmCmdDef{key: "HumanLocal", name: "HumanLocal", permMin: 3, usage: "人物名称", handler: (*PlayObject).cmdHumanLocal})
	registerGMCommand(&gmCmdDef{key: "PositionMove", name: "PositionMove", permMin: 3, usage: "地图号 座标X 座标Y", handler: (*PlayObject).cmdPositionMove})
	registerGMCommand(&gmCmdDef{key: "MobLevel", name: "MobLevel", permMin: 3, handler: (*PlayObject).cmdMobLevel})
	registerGMCommand(&gmCmdDef{key: "MobCount", name: "MobCount", permMin: 3, usage: "地图号", handler: (*PlayObject).cmdMobCount})
	registerGMCommand(&gmCmdDef{key: "HumanCount", name: "HumanCount", permMin: 3, usage: "地图号", handler: (*PlayObject).cmdHumanCount})
	registerGMCommand(&gmCmdDef{key: "Ting", name: "Ting", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdTing})
	registerGMCommand(&gmCmdDef{key: "SuperTing", name: "SuperTing", permMin: 10, usage: "人物名称 范围(0-10)", handler: (*PlayObject).cmdSuperTing})
	registerGMCommand(&gmCmdDef{key: "MapMove", name: "MapMove", permMin: 10, usage: "源地图 目标地图", handler: (*PlayObject).cmdMapMove})
	registerGMCommand(&gmCmdDef{key: "Kill", name: "Kill", permMin: 10, usage: "[人物名称]", handler: (*PlayObject).cmdKill})
}

// gmMsgNotOnLine '%s 现在不在线，或在其它服务器上！！！'
// （g_sNowNotOnLineOrOnOtherServer，M2Share.pas:3257）。
const gmMsgNotOnLine = "%s 现在不在线，或在其它服务器上！！！"

// cmdHumanLocal 查询玩家 IP 归属（Delphi CmdHumanLocal，分发 ObjBase.pas:7811，
// 实现 13143-13165）。文案 '%s来自:%s'（g_sGameCommandHumanLocalMsg，
// M2Share.pas:3331，c_Green）。
// 原版 bug：Delphi 输出的是 GM 自己的 m_sIPLocal（ObjBase.pas:13164），
// Go 修正为输出目标玩家的 IP；Go 无 IP 归属地库（Delphi GetIPLocal 亦被注释），
// 简化为输出目标登录 IP。
func (p *PlayObject) cmdHumanLocal(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	p.sysMsgGreen(server, fmt.Sprintf("%s来自:%s", target.Name, sessionRemoteIP(target.Session)))
	return true
}

// cmdPositionMove 定点传送（Delphi CmdPositionMove，分发 ObjBase.pas:7821，
// 实现 13480-13523）。权限表 nPermissionMin=3、nPermissionMax=6
// （M2Share.pas:2768）：权限≥6 豁免禁传检查。
func (p *PlayObject) cmdPositionMove(server *netserver.TCPServer, params []string) bool {
	if len(params) < 3 {
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	if env == nil {
		return true // Delphi FindMap=nil 静默退出（ObjBase.pas:13498-13510）
	}
	// Delphi 用 CanMoveMap 禁传表 + g_sTheMapDisableMove
	//（'地图%s(%s)不允许传送...'，M2Share.pas:3185）；Go 无禁传表，
	// 以地图标志 NoPositionMove 替代，权限≥6 豁免（对齐 nPermissionMax=6）。
	// Go 地图格式无描述字段，第二参以地图号填充。
	if env.Flag.NoPositionMove && int(p.Permission) < 6 {
		p.sysMsg(server, fmt.Sprintf("地图%s(%s)不允许传送...", params[0], params[0]))
		return true
	}
	// Delphi Str_ToInt 解析失败按 0（ObjBase.pas:13501-13502）。
	x, _ := strconv.Atoi(params[1])
	y, _ := strconv.Atoi(params[2])
	// Delphi Envir.CanWalk(nX, nY, True)：boFlag=True 只查地形不查实体
	//（Envir.pas:392，等同 Go CanWalkEx ignoreEntities）。
	if env.CanWalkEx(x, y, true) {
		teleportPlayerTo(p, server, env, x, y)
	} else {
		// '无法移动到地图: %s X:%s Y:%s'
		//（g_sGameCommandPositionMoveCanotMoveToMap，M2Share.pas:3254，c_Green）。
		p.sysMsgGreen(server, fmt.Sprintf("无法移动到地图: %s X:%s Y:%s", params[0], params[1], params[2]))
	}
	return true
}

// cmdMobLevel 输出身边 2 格内所有对象基本信息（Delphi CmdMobLevel，
// 分发 ObjBase.pas:7831，实现 13751-13776：GetRangeBaseObject 半径 2，
// 逐个输出 GeTBaseObjectInfo，c_Green）。
// Go 简化为 @info 风格单行格式（Delphi GeTBaseObjectInfo 含攻防明细，
// ObjBase.pas:22564-22580）。
func (p *PlayObject) cmdMobLevel(server *netserver.TCPServer, params []string) bool {
	if p.envir == nil {
		return true
	}
	for _, obj := range p.envir.GetRangeObjects(p.CurrX, p.CurrY, 2) {
		switch o := obj.(type) {
		case *PlayObject:
			if o.Ghost {
				continue
			}
			p.sysMsgGreen(server, fmt.Sprintf("[%s] Lv%d %s HP:%d/%d MP:%d/%d Gold:%d Map:%s(%d,%d) PK:%d",
				o.Name, o.WAbil.Level, jobName(o.Job),
				o.WAbil.HP, o.WAbil.MaxHP, o.WAbil.MP, o.WAbil.MaxMP,
				o.Gold, o.MapName, o.CurrX, o.CurrY, o.PkPoint))
		case *MonsterObject:
			if o.Ghost {
				continue
			}
			p.sysMsgGreen(server, fmt.Sprintf("[%s] Lv%d 怪物 HP:%d/%d MP:%d/%d Map:%s(%d,%d)",
				o.Name, o.WAbil.Level,
				o.WAbil.HP, o.WAbil.MaxHP, o.WAbil.MP, o.WAbil.MaxMP,
				o.MapName, o.CurrX, o.CurrY))
		}
	}
	return true
}

// cmdMobCount 查询地图存活怪物数量（Delphi CmdMobCount，分发 ObjBase.pas:7836，
// 实现 13659-13680）。Delphi 统计 MonGen CertList 中 !Death && !Ghost 且
// 位于目标地图的怪物（UsrEngn.pas:2879-2901）。
func (p *PlayObject) cmdMobCount(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	if env == nil {
		// '指定的地图不存在！！！'（g_sGameCommandMobCountMapNotFound，M2Share.pas:3259）。
		p.sysMsg(server, "指定的地图不存在！！！")
		return true
	}
	// '怪物数量：%d'（g_sGameCommandMobCountMonsterCount，M2Share.pas:3260，c_Green）。
	p.sysMsgGreen(server, fmt.Sprintf("怪物数量：%d", countMapMonsters(p.Engine, env)))
	return true
}

// countMapMonsters 目标地图存活怪物数（Delphi UserEngine.GetMapMonster，
// UsrEngn.pas:2879-2901：!Death && !Ghost 且 m_PEnvir = Envir）。
func countMapMonsters(e *UserEngine, env *Environment) int {
	count := 0
	for _, mon := range e.Monsters {
		if !mon.Death && !mon.Ghost && mon.envir == env {
			count++
		}
	}
	return count
}

// cmdHumanCount 查询地图在线人数（Delphi CmdHumanCount，分发 ObjBase.pas:7841，
// 实现 13683-13705）。统计 !Death && !Ghost 且位于目标地图的玩家
// （GetMapHuman，UsrEngn.pas:2920-2936）。
// 原版 bug：Delphi 复用 g_sGameCommandMobCountMonsterCount 输出
// '怪物数量：人数' 并追加一行 Envir.HumCount（ObjBase.pas:13703-13704）；
// M2Share 无 HumanCount 专用文案常量（仅 HelpMsg，M2Share.pas:3261），
// Go 修正为单行 '人数：%d'。
func (p *PlayObject) cmdHumanCount(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	env := p.MapMgr.FindMap(params[0])
	if env == nil {
		p.sysMsg(server, "指定的地图不存在！！！")
		return true
	}
	p.sysMsgGreen(server, fmt.Sprintf("人数：%d", countMapHumans(p.Engine, env)))
	return true
}

// countMapHumans 目标地图在线人数（Delphi UserEngine.GetMapHuman，
// UsrEngn.pas:2920-2936：!Death && !Ghost 且 m_PEnvir = Envir）。
func countMapHumans(e *UserEngine, env *Environment) int {
	count := 0
	e.mu.Lock()
	for _, pl := range e.PlayObjectList {
		if !pl.Death && !pl.Ghost && pl.envir == env {
			count++
		}
	}
	e.mu.Unlock()
	return count
}

// cmdTing 把目标随机传送回主城（Delphi CmdTing，分发 ObjBase.pas:7851，
// 实现 15161-15184：PlayObject.MapRandomMove(m_sHomeMap, 0)）。
// 原版行为：用的是 GM 自己的 m_sHomeMap；Go 无各玩家 hometown 字段，
// 统一用 config 回城点（GetHomeMap/GetHomeX/GetHomeY）。
func (p *PlayObject) cmdTing(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	p.gmSendHome(target, server)
	return true
}

// gmSendHome 随机传送回 config 回城点（Delphi MapRandomMove，ObjBase.pas:
// TBaseObject.MapRandomMove：随机落点 + SpaceMove；地图不存在时静默）。
func (p *PlayObject) gmSendHome(target *PlayObject, server *netserver.TCPServer) {
	cfg := p.Engine.Config
	homeEnv := p.MapMgr.FindMap(cfg.GetHomeMap())
	if homeEnv == nil {
		return // Delphi MapRandomMove 对不存在地图静默退出
	}
	x, y := randomMapPosition(homeEnv, 20)
	teleportPlayerTo(target, server, homeEnv, x, y)
}

// cmdSuperTing 把目标周围玩家全部遣送回城（Delphi CmdSuperTing，
// 分发 ObjBase.pas:7856，实现 15006-15043）。
// 原版 bug：nRange := _MAX(10, Str_ToInt(sRange, 2)) 实为 max(10, x)，
// 范围恒 ≥10（ObjBase.pas:15025，_MAX 定义 Common/HUtil32.pas:1959）；
// Go 修正为默认 2、上限 10。Delphi 传送各人自己的 m_sHomeMap，
// Go 统一 config 回城点。排除 GM 自己（ObjBase.pas:15034）。
func (p *PlayObject) cmdSuperTing(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	nRange := 2 // Delphi Str_ToInt(sRange, 2) 默认值
	if len(params) >= 2 {
		if v, err := strconv.Atoi(params[1]); err == nil {
			nRange = v
		}
	}
	if nRange < 0 {
		nRange = 0
	}
	if nRange > 10 {
		nRange = 10
	}
	cfg := p.Engine.Config
	homeEnv := p.MapMgr.FindMap(cfg.GetHomeMap())
	if homeEnv == nil {
		return true
	}
	// Delphi GetMapRageHuman：同图、!Death && !Ghost、
	// abs(dx)<=nRange && abs(dy)<=nRange（UsrEngn.pas:2938-2958）。
	p.Engine.mu.Lock()
	victims := make([]*PlayObject, 0, 8)
	for _, pl := range p.Engine.PlayObjectList {
		if pl == p || pl.Death || pl.Ghost {
			continue
		}
		if pl.envir != target.envir || pl.MapName != target.MapName {
			continue
		}
		if abs(pl.CurrX-target.CurrX) <= nRange && abs(pl.CurrY-target.CurrY) <= nRange {
			victims = append(victims, pl)
		}
	}
	p.Engine.mu.Unlock()
	for _, pl := range victims {
		p.gmSendHome(pl, server)
	}
	return true
}

// cmdMapMove 整图玩家迁移（Delphi CmdMapMoveHuman，分发 ObjBase.pas:7861，
// 实现 13524-13564）。源/目标地图不存在均报
// '地图%s不存在！！！'（g_sGameCommandMapMoveMapNotFound，M2Share.pas:3266）。
// Delphi 以源图中心半径 1000 收集（等效全图），除 GM 自己外
// MapRandomMove 到目标图；Go 用 randomMapPosition + teleportPlayerTo。
func (p *PlayObject) cmdMapMove(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		return true
	}
	srcEnv := p.MapMgr.FindMap(params[0])
	if srcEnv == nil {
		p.sysMsg(server, fmt.Sprintf("地图%s不存在！！！", params[0]))
		return true
	}
	dstEnv := p.MapMgr.FindMap(params[1])
	if dstEnv == nil {
		p.sysMsg(server, fmt.Sprintf("地图%s不存在！！！", params[1]))
		return true
	}
	p.Engine.mu.Lock()
	moved := make([]*PlayObject, 0, 16)
	for _, pl := range p.Engine.PlayObjectList {
		if pl == p || pl.Death || pl.Ghost {
			continue
		}
		if pl.envir == srcEnv {
			moved = append(moved, pl)
		}
	}
	p.Engine.mu.Unlock()
	for _, pl := range moved {
		x, y := randomMapPosition(dstEnv, 20)
		teleportPlayerTo(pl, server, dstEnv, x, y)
	}
	return true
}

// cmdKill 杀死目标（Delphi CmdKill，分发 ObjBase.pas:8023，实现 13243-13270）。
// 无参=杀面前一格对象（GetPoseCreate：GetFrontPosition + GetMovingObject(x,y,True)，
// ObjBase.pas:15875-15884；True 跳过 ghost/dead，Envir.pas GetMovingObject）；
// 有参=按名字找玩家。Delphi 调 BaseObject.Die（内部豁免 m_boSuperMan，
// ObjBase.pas:20765-20782）；Go 走正常死亡结算路径（applyDamage 致死量）。
func (p *PlayObject) cmdKill(server *netserver.TCPServer, params []string) bool {
	if p.envir == nil {
		return true
	}
	if len(params) >= 1 {
		target := p.Engine.GetPlayerByName(params[0])
		if target == nil {
			p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
			return true
		}
		if target.MapName == p.MapName {
			p.applyDamage(server, target.BaseObject, gmLethalDamage(target.BaseObject), p.Dir, false)
		} else {
			// applyDamage 的死亡广播/结算基于攻击者所在地图，
			// 跨图目标按目标自身地图执行同款结算。
			p.gmKillRemotePlayer(server, target)
		}
		return true
	}
	dx, dy := dirToOffset(p.Dir)
	var base *BaseObject
	switch o := p.envir.GetMovingObject(p.CurrX+dx, p.CurrY+dy).(type) {
	case *PlayObject:
		base = o.BaseObject
	case *MonsterObject:
		base = o.BaseObject
	}
	// Delphi GetMovingObject(x, y, True) 排除 ghost/dead。
	if base == nil || base.Ghost || base.Death {
		// 原文如此（ObjBase.pas:13265）。
		p.sysMsg(server, "命令使用方法不正确，必须与角色面对面站好！！！")
		return true
	}
	p.applyDamage(server, base, gmLethalDamage(base), p.Dir, false)
	return true
}

// gmLethalDamage 返回保证致死的伤害值：需大于魔法盾完全吸收上限
// （mpCost = damage×(1+ratio%) 必须超过 MaxMP 才走部分吸收分支，
// 部分吸收至多扣 MaxMP×2/3），故取 MaxHP + 2×MaxMP + 10000。
func gmLethalDamage(base *BaseObject) int {
	return int(base.WAbil.MaxHP) + int(base.WAbil.MaxMP)*2 + 10000
}

// gmKillRemotePlayer 跨图 @Kill：在目标所在地图复刻 applyDamage 致死分支的
// 结算（playobject.go applyDamage hp<=0 段）——无敌豁免、复活戒指、
// 死亡广播、PK 结算、死亡掉落。
func (p *PlayObject) gmKillRemotePlayer(server *netserver.TCPServer, target *PlayObject) {
	if target.Death || target.Ghost || target.SuperManMode {
		return // Delphi Die 豁免 m_boSuperMan（ObjBase.pas:20782）
	}
	cfg := p.Engine.Config
	// 复活戒指判定与 applyDamage 一致（攻击者佩戴禁复活则不生效）。
	if target.HasRevival && !p.HasUnRevival {
		target.HasRevival = false
		target.WAbil.HP = target.WAbil.MaxHP / uint16(cfg.GetReviveHPRatio())
		target.damageRevivalRings(server)
		target.sendHealthSpell(server)
		return
	}
	target.WAbil.HP = 0
	target.Death = true
	target.deathTick = time.Now().UnixMilli()
	if target.envir != nil {
		target.envir.broadcastDeathMsg(target.BaseObject, target.ID, target.CurrX, target.CurrY, target.Dir, true)
	}
	p.OnPlayerKilled(server, target)
	// 被玩家击杀：默认不掉装备（Delphi boKillByHumanDropUseItem=False）。
	target.DropDeathItems(server, false)
}
