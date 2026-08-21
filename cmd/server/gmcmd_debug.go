package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// GM 命令——调试杂项族移植。
// Delphi 参照：ObjBase.pas（分发全部位于 8254 m_btPermission>4 调试段，
// 内层 8387 >=5/boTestServer 门槛；实现 10604-15160）+ M2Share.pas
// g_GameCommand（2805-2924）与文案常量。
// 注册顺序 = Delphi 分发链出现顺序。
// 未移植：
//   - Reconnection（ObjBase.pas:8353/14086-14105）——语义为令客户端换网关
//     重连（SM_RECONNECT），Go 单端口架构无网关可换，客户端亦无消费。
//   - TestServerConfig（ObjBase.pas:8564/16803-16844）——发 SM_SERVERCONFIG，
//     Go 协议未定义该消息、客户端无消费端。

func init() {
	// 以下命令实现内另有 m_btPermission<6 静默退出 → perm 6 + permFailSilent。
	registerGMCommand(&gmCmdDef{key: "Backstep", name: "Backstep", permMin: 6, permFail: permFailSilent, usage: "[类型(0-8)] [次数]", handler: (*PlayObject).cmdBackstep})
	registerGMCommand(&gmCmdDef{key: "Hunger", name: "Hunger", permMin: 6, permFail: permFailSilent, usage: "人物名称 能量值", handler: (*PlayObject).cmdHunger})
	registerGMCommand(&gmCmdDef{key: "DisableFilter", name: "DisableFilter", permMin: 6, permFail: permFailSilent, handler: (*PlayObject).cmdDisableFilter})
	registerGMCommand(&gmCmdDef{key: "ChgUserFull", name: "CHGUSERFULL", permMin: 6, permFail: permFailSilent, usage: "人数", handler: (*PlayObject).cmdChgUserFull})
	registerGMCommand(&gmCmdDef{key: "ChgZenFastStep", name: "CHGZENFASTSTEP", permMin: 6, permFail: permFailSilent, usage: "速度", handler: (*PlayObject).cmdChgZenFastStep})
	registerGMCommand(&gmCmdDef{key: "TestStatus", name: "TestStatus", permMin: 6, permFail: permFailSilent, usage: "类型(0..11) 时长", handler: (*PlayObject).cmdTestStatus})
	registerGMCommand(&gmCmdDef{key: "Spirit", name: "祈祷生效", permMin: 6, permFail: permFailSilent, usage: "[秒数]", handler: (*PlayObject).cmdSpiritStart})
	registerGMCommand(&gmCmdDef{key: "SpiritStop", name: "停止叛变", permMin: 6, permFail: permFailSilent, handler: (*PlayObject).cmdSpiritStop})
	// ServerStatus 实现内 m_btPermission<10 静默退出（ObjBase.pas:16847）。
	registerGMCommand(&gmCmdDef{key: "ServerStatus", name: "ServerStatus", permMin: 10, permFail: permFailSilent, handler: (*PlayObject).cmdServerStatus})
	// 以下命令走常规 nPermissionMin 检查（'权限不够...'）。
	registerGMCommand(&gmCmdDef{key: "TestGetBagItem", name: "TestGetBagItem", permMin: 10, usage: "(用于测试升级武器方面参数)", handler: (*PlayObject).cmdTestGetBagItem})
	registerGMCommand(&gmCmdDef{key: "TestSpeedMode", name: "TestSpeedMode", permMin: 10, handler: (*PlayObject).cmdTestSpeedMode})
}

// cmdBackstep 把自己向后推（Delphi CmdBackStep，ObjBase.pas:10604-10617：
// nType=0 向身后推，否则随机 0..nType-1 方向；CharPushed 逐格推进，
// ObjBase.pas:2504-2537）。
// 偏差：Delphi 逐格 MoveToMovingObject + RM_PUSH 动画；Go 无 RM_PUSH
// 消费端，改为逐格 CanWalk 探测后一次性传送（RM_DISAPPEAR/RM_LOGON）。
func (p *PlayObject) cmdBackstep(server *netserver.TCPServer, params []string) bool {
	if p.envir == nil {
		return true
	}
	nType := 0
	if len(params) >= 1 {
		nType = atoiOr(params[0], 0) // Delphi Str_ToInt(sParam1, 0)
	}
	nCount := 1
	if len(params) >= 2 {
		nCount = atoiOr(params[1], 1) // Delphi Str_ToInt(sParam2, 1)
	}
	if nType > 8 {
		nType = 8 // Delphi _MIN(nType, 8)
	}
	pushDir := (p.Dir + 4) % 8 // Delphi GetBackDir：向身后推
	if nType != 0 {
		pushDir = rand.Intn(nType)
	}
	dx, dy := dirToOffset(pushDir)
	x, y := p.CurrX, p.CurrY
	for i := 0; i < nCount; i++ {
		if !p.envir.CanWalk(x+dx, y+dy) {
			break
		}
		x += dx
		y += dy
	}
	if x == p.CurrX && y == p.CurrY {
		return true
	}
	// Delphi CharPushed 末尾把朝向置为推力反方向（GetBackDir(nDir)，
	// ObjBase.pas:2533）；nType=0 时即保持原朝向。
	p.Dir = (pushDir + 4) % 8
	teleportPlayerTo(p, server, p.envir, x, y)
	return true
}

// cmdHunger 设置目标能量值（Delphi CmdHunger，ObjBase.pas:13167-13190）。
func (p *PlayObject) cmdHunger(server *netserver.TCPServer, params []string) bool {
	nHungerPoint := 0
	if len(params) >= 2 {
		nHungerPoint = atoiOr(params[1], 0) // Delphi Str_ToInt(sParam2, 0)
	}
	if len(params) < 1 || params[0] == "" || nHungerPoint < 0 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("Hunger", "Hunger")+" 人物名称 能量值")
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		// 原文（ObjBase.pas:13186）：与通用离线文案不同。
		p.sysMsg(server, params[0]+"没有在线！！！")
		return true
	}
	target.HungerStatus = nHungerPoint
	target.sendMyStatus(server) // Delphi RM_MYSTATUS + RefMyStatus
	p.sysMsgGreen(server, params[0]+" 的能量值已改变。")
	return true
}

// cmdDisableFilter 文字过滤总开关（Delphi CmdDisableFilter，
// ObjBase.pas:12554-12571：翻转全局 boFilterWord，M2Share.pas:1343）。
// Go 对应 wordfilter.go 过滤开关（filterChatText 消费）。
func (p *PlayObject) cmdDisableFilter(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.sysMsg(server, "启用/禁止文字过滤功能。")
		return true
	}
	setWordFilterEnabled(!wordFilterIsEnabled())
	if wordFilterIsEnabled() {
		p.sysMsgGreen(server, "已启用文字过滤。")
	} else {
		p.sysMsgGreen(server, "已禁止文字过滤。")
	}
	return true
}

// cmdChgUserFull 运行时改服务器上线人数（Delphi CmdChangeUserFull，
// ObjBase.pas:11015-11030：写 g_Config.nUserFull，仅内存不落盘）。
// Go 写 Config.Game.UserFull，mongen.go getZenTime 每轮刷怪消费。
func (p *PlayObject) cmdChgUserFull(server *netserver.TCPServer, params []string) bool {
	s := ""
	if len(params) > 0 {
		s = params[0]
	}
	nCount := atoiOr(s, -1) // Delphi Str_ToInt(sUserCount, -1)
	if s == "" || s[0] == '?' || nCount < 1 {
		p.sysMsg(server, "设置服务器最高上线人数。")
		p.sysMsg(server, "命令格式: @"+gmNameByKey("ChgUserFull", "CHGUSERFULL")+" 人数")
		return true
	}
	p.Engine.Config.Game.UserFull = nCount
	p.sysMsgGreen(server, fmt.Sprintf("服务器上线人数限制: %d", nCount))
	return true
}

// cmdChgZenFastStep 运行时改刷怪加速步长（Delphi CmdChangeZenFastStep，
// ObjBase.pas:11032-11047：写 g_Config.nZenFastStep，仅内存）。
// Go 写 Config.Game.ZenFastStep，mongen.go getZenTime 消费。
func (p *PlayObject) cmdChgZenFastStep(server *netserver.TCPServer, params []string) bool {
	s := ""
	if len(params) > 0 {
		s = params[0]
	}
	nFastStep := atoiOr(s, -1) // Delphi Str_ToInt(sFastStep, -1)
	if s == "" || s[0] == '?' || nFastStep < 1 {
		p.sysMsg(server, "设置怪物行动速度。")
		p.sysMsg(server, "命令格式: @"+gmNameByKey("ChgZenFastStep", "CHGZENFASTSTEP")+" 速度")
		return true
	}
	p.Engine.Config.Game.ZenFastStep = nFastStep
	p.sysMsgGreen(server, fmt.Sprintf("怪物行动速度: %d", nFastStep))
	return true
}

// cmdTestStatus 给自己附加状态时长（Delphi CmdTestStatus，
// ObjBase.pas:15146-15160：m_wStatusTimeArr[nType] := nTime*1000 后
// GetCharStatus + StatusChanged）。
// 偏差：Delphi 存毫秒按 tick 差递减；Go StatusTimeArr 每 tick 减 1，
// 按 tickInterval 换算（秒→tick 数，int16 上限钳制）；MakePoison 的
// Hidden 副作用不联动（调试命令直接写数组，与 Delphi 一致）。
func (p *PlayObject) cmdTestStatus(server *netserver.TCPServer, params []string) bool {
	nType := -1
	if len(params) >= 1 {
		nType = atoiOr(params[0], -1) // Delphi Str_ToInt(sParam1, -1)
	}
	nTime := 0
	if len(params) >= 2 {
		nTime = atoiOr(params[1], 0) // Delphi Str_ToInt(sParam2, 0)
	}
	if nType < 0 || nType > 11 || nTime < 0 {
		p.sysMsg(server, "命令格式: @"+gmNameByKey("TestStatus", "TestStatus")+" 类型(0..11) 时长")
		return true
	}
	tickMs := p.Engine.Config.GetTickInterval()
	if tickMs <= 0 {
		tickMs = 100
	}
	ticks := int64(nTime) * 1000 / tickMs
	if ticks > 32767 {
		ticks = 32767
	}
	p.StatusTimeArr[nType] = int16(ticks)
	p.BroadcastStatus() // Delphi StatusChanged（ObjBase.pas:20139-20142）
	p.sysMsgGreen(server, fmt.Sprintf("状态编号:%d 时间长度: %d 秒", nType, nTime))
	return true
}

// cmdSpiritStart 开始祈祷生效宝宝叛变（Delphi CmdSpirtStart，
// ObjBase.pas:14953-14975：写全局 g_dwSpiritMutinyTick =
// GetTickCount + dwTime，M2Share.pas:1509）。
// Go 消费端：monsterobject.go slaveFollow 生效期清除主从纽带。
func (p *PlayObject) cmdSpiritStart(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.sysMsg(server, "此命令用于开始祈祷生效宝宝叛变。")
		return true
	}
	nTime := -1
	if len(params) > 0 {
		nTime = atoiOr(params[0], -1) // Delphi Str_ToInt(sParam1, -1)
	}
	// Delphi g_Config.dwSpiritMutinyTime 默认 30 分钟（M2Share.pas:2046）。
	dwTime := int64(30 * 60 * 1000)
	if nTime > 0 {
		dwTime = int64(nTime) * 1000
	}
	e := p.Engine
	e.mu.Lock()
	e.SpiritMutinyTick = time.Now().UnixMilli() + dwTime
	e.mu.Unlock()
	p.sysMsgGreen(server, "祈祷叛变已开始。持续时长 "+strconv.FormatInt(dwTime/1000, 10)+" 秒。")
	return true
}

// cmdSpiritStop 停止祈祷叛变（Delphi CmdSpirtStop，
// ObjBase.pas:14976-14989：g_dwSpiritMutinyTick := 0）。
func (p *PlayObject) cmdSpiritStop(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.sysMsg(server, "此命令用于停止祈祷生效导致宝宝叛变。")
		return true
	}
	e := p.Engine
	e.mu.Lock()
	e.SpiritMutinyTick = 0
	e.mu.Unlock()
	p.sysMsgGreen(server, "祈祷叛变已停止。")
	return true
}

// cmdServerStatus 输出服务端可执行文件 CRC（Delphi SendServerStatus，
// ObjBase.pas:16846-16850：CalcFileCRC(Application.ExeName) 防篡改自检，
// 实现内 m_btPermission<10 静默退出）。Go 对当前进程二进制做 CRC32。
func (p *PlayObject) cmdServerStatus(server *netserver.TCPServer, params []string) bool {
	exe, err := os.Executable()
	if err != nil {
		return true
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		return true
	}
	p.sysMsg(server, strconv.FormatUint(uint64(crc32.ChecksumIEEE(data)), 10))
	return true
}

// cmdTestGetBagItem 升级武器材料测试（Delphi CmdTestGetBagItems，
// ObjBase.pas:15105-15126 → GetBagUseItems，15630 起）。
// 破坏性命令（Delphi 原版如此）：消耗背包全部黑铁矿（g_Config.sBlackStone
// ≡ Go GetUpgradeMaterial）并删除全部首饰（StdMode 19-26，与 upgrade.go
// 材料分析口径一致），回显累计 DC/SC/MC/DURA。
func (p *PlayObject) cmdTestGetBagItem(server *netserver.TCPServer, params []string) bool {
	var dc, sc, mc, dura int
	if p.ItemDB != nil {
		matName := p.Engine.Config.GetUpgradeMaterial()
		remaining := make([]*protocol.UserItem, 0, len(p.ItemList))
		for _, item := range p.ItemList {
			if item == nil {
				continue
			}
			def := p.ItemDB.GetByIdx(int(item.WIndex))
			if def == nil {
				remaining = append(remaining, item)
				continue
			}
			if def.Name == matName {
				dura += int(item.Dura) / 1000
				continue // Delphi DuraList 收集后删除
			}
			if def.StdMode >= 19 && def.StdMode <= 26 {
				dc += int(def.DC + def.DCMax)
				mc += int(def.MC + def.MCMax)
				sc += int(def.SC + def.SCMax)
				continue // Delphi DelItemList 收集后删除
			}
			remaining = append(remaining, item)
		}
		p.ItemList = remaining
		p.SendBagItemsFull(server)
	}
	// g_sGameCommandTestGetBagItemsHelpMsg='(用于测试升级武器方面参数)'
	//（M2Share.pas:3349）作 '?' 帮助由注册表输出；结果蓝字（ObjBase.pas:15124）。
	p.sysMsgBlue(server, fmt.Sprintf("DC:%d SC:%d MC:%d DURA:%d", dc, sc, mc, dura))
	return true
}

// cmdTestSpeedMode 速度测试模式开关（Delphi CmdTestSpeedMode，
// ObjBase.pas:15128-15144：翻转 m_boTestSpeedMode，动作速度校验处豁免，
// ObjBase.pas:4735/4949/4993）。Go 消费端：checkActionSpeed。
func (p *PlayObject) cmdTestSpeedMode(server *netserver.TCPServer, params []string) bool {
	p.TestSpeedMode = !p.TestSpeedMode
	if p.TestSpeedMode {
		p.sysMsg(server, "开启速度测试模式")
	} else {
		p.sysMsg(server, "关闭速度测试模式")
	}
	return true
}
