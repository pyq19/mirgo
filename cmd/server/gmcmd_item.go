package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// GM 命令——物品族。
// Delphi 参照：ObjBase.pas（分发 ProcessUserLineMsg 7188-8600）+
// M2Share.pas g_GameCommand 命令表（2724-2926）与文案常量。
// 注册顺序 = Delphi 分发链出现顺序（ObjBase.pas:8026-8384）。
// 已实现跳过：@make（含数量）、@ClearBag、@ShowUseItemInfo、@BindUseItem（T1.3）。

func init() {
	registerGMCommand(&gmCmdDef{key: "SuperMake", name: "Supermake", permMin: 10, usage: "部位(0-12) 属性位(0-15) 值(0-255)", handler: (*PlayObject).cmdSuperMake})
	registerGMCommand(&gmCmdDef{key: "DeleteItem", name: "DeleteItem", permMin: 10, usage: "人物名称 物品名 数量", handler: (*PlayObject).cmdDeleteItem})
	registerGMCommand(&gmCmdDef{key: "StorageClearPassword", name: "ClrPassword", permMin: 10, usage: "人物名称", handler: (*PlayObject).cmdClrPassword})
	registerGMCommand(&gmCmdDef{key: "RefineWeapon", name: "RefineWeapon", permMin: 10, usage: "DC MC SC HIT", handler: (*PlayObject).cmdRefineWeapon})
	// Delphi 分发位于调试段（门槛 >4，ObjBase.pas:8254），实现内硬编码
	// m_btPermission<6 静默退出（ObjBase.pas:10795）→ 等价 perm 6 + 静默失败。
	registerGMCommand(&gmCmdDef{key: "ChangeItemName", name: "ChangeItemName", permMin: 6, permFail: permFailSilent, usage: "物品编号 物品ID号 物品名称", handler: (*PlayObject).cmdChangeItemName})
}

// cmdSuperMake 调整自己穿戴物品属性（Delphi CmdSmakeItem，
// 分发 ObjBase.pas:8026-8028，实现 14899-14952；命令表 SMAKE
// 'Supermake' 权限 10，M2Share.pas:2843，ini 键 SuperMake 6223）。
//
// 属性位映射：Delphi TUserItem.btValue[0..13] 与 Go UserItem.BtValue[0..13]
// 下标一一对应（布局一致，无需换算）：武器 0=DC加点 1=MC加点 2=SC加点
// 3=祝福 4=诅咒 5=HIT 6=攻速 7=强化（ItmUnit.pas:120-133，
// Go 消费端 itemsystem.go RecalcAbilitys 武器分支）；
// 属性位 14/15 为持久/持久上限（Delphi 上限 65，×1000 存储）。
func (p *PlayObject) cmdSuperMake(server *netserver.TCPServer, params []string) bool {
	if len(params) < 3 {
		return true
	}
	// Delphi Str_ToInt(sParamN, 0)（ObjBase.pas:8028）：解析失败按 0。
	nWhere, _ := strconv.Atoi(params[0])
	nValueType, _ := strconv.Atoi(params[1])
	nValue, _ := strconv.Atoi(params[2])
	// Delphi 越界静默退出（ObjBase.pas:14910）。
	if nWhere < 0 || nWhere > 12 || nValueType < 0 || nValueType > 15 || nValue < 0 || nValue > 255 {
		return true
	}
	// Delphi m_UseItems[nWhere].wIndex <= 0 →
	// g_sGamecommandSuperMakeHelpMsg（M2Share.pas:3340）。
	if p.UseItems[nWhere] == nil || p.UseItems[nWhere].WIndex <= 0 {
		p.sysMsg(server, "身上没指定物品！！！")
		return true
	}
	item := p.UseItems[nWhere]
	var def *ItemDef
	if p.ItemDB != nil {
		def = p.ItemDB.GetByIdx(int(item.WIndex))
	}
	if def == nil {
		return true // Delphi GetStdItem=nil 静默退出（ObjBase.pas:14915）
	}
	if nValueType > 13 {
		// Delphi nValue := _MIN(65, nValue)（ObjBase.pas:14919-14921）。
		if nValue > 65 {
			nValue = 65
		}
		if nValueType == 14 {
			item.Dura = uint16(nValue * 1000)
		} else {
			item.DuraMax = uint16(nValue * 1000)
		}
	} else {
		item.BtValue[nValueType] = byte(nValue)
	}
	p.RecalcAbilitys()
	p.SendUseItemsFull(server) // Delphi SendUpdateItem（ObjBase.pas:14927）
	// 完整属性串：wIndex-MakeIndex Dura/DuraMax bt0/bt1/.../bt13
	//（ObjBase.pas:14929-14944，蓝字 c_Blue）。
	bts := make([]string, 0, 14)
	for i := 0; i < 14; i++ {
		bts = append(bts, strconv.Itoa(int(item.BtValue[i])))
	}
	sShowMsg := fmt.Sprintf("%d-%d %d/%d %s",
		item.WIndex, item.MakeIndex, item.Dura, item.DuraMax, strings.Join(bts, "/"))
	p.sysMsgBlue(server, sShowMsg)
	// Delphi 控制台输出受 g_Config.boShowMakeItemMsg 控制（默认 False，
	// M2Share.pas:1768）；Go 无此开关，与 @level 一致恒输出（ObjBase.pas:14945-14946）。
	log.Logf(log.LevelInfo, "GM", "[物品调整] %s(%s -> %s)", p.Name, def.Name, sShowMsg)
	return true
}

// cmdDeleteItem 按名称从目标背包删除指定数量物品（Delphi CmdDeleteItem，
// 分发 ObjBase.pas:8101-8103，实现 12255-12295；命令表 DELETEITEM
// 'DeleteItem' 权限 10，M2Share.pas:2821，ini 键 DeleteItem 6019）。
// 倒序遍历、名称大小写不敏感（CompareText）、最多删 nCount 件（默认 1，
// Str_ToInt(sParam3, 1)）。Delphi 逐件 SendDelItems 通知客户端，
// Go 无单件删除消息，改发整包 SendBagItemsFull；Delphi 无成功提示（照做）。
func (p *PlayObject) cmdDeleteItem(server *netserver.TCPServer, params []string) bool {
	if len(params) < 2 {
		// 尾部多余的 ')' 为 Delphi 原文笔误（ObjBase.pas:12271），照抄。
		p.sysMsg(server, "命令格式: @DeleteItem 人物名称 物品名称 数量)")
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		p.sysMsg(server, fmt.Sprintf(gmMsgNotOnLine, params[0]))
		return true
	}
	count := 1
	if len(params) >= 3 {
		if v, err := strconv.Atoi(params[2]); err == nil {
			count = v
		}
	}
	deleted := 0
	for i := len(target.ItemList) - 1; i >= 0 && deleted < count; i-- {
		item := target.ItemList[i]
		if item == nil || target.ItemDB == nil {
			continue
		}
		def := target.ItemDB.GetByIdx(int(item.WIndex))
		if def != nil && strings.EqualFold(params[1], def.Name) {
			target.ItemList = append(target.ItemList[:i], target.ItemList[i+1:]...)
			deleted++
		}
	}
	if deleted > 0 {
		target.SendBagItemsFull(server)
	}
	return true
}

// cmdClrPassword 清除目标仓库密码（Delphi CmdClearHumanPassword，
// 分发 ObjBase.pas:8209-8211，实现 11088-11110；命令表 CLRPASSWORD
// 'ClrPassword' 权限 10，M2Share.pas:2749，ini 键 StorageClearPassword 5462）。
// Delphi 清空 m_sStoragePwd/m_boPasswordLocked/m_boUnLockStoragePwd；
// Go 对应 StoragePassword/StoragePwdLocked/storagePwdFail
//（playobject.go:212-220，逻辑见 storagesystem.go）。
func (p *PlayObject) cmdClrPassword(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		// '?' 参数由注册表 usage 帮助统一拦截；无参分支保留 Delphi
		// 双行文案（ObjBase.pas:11095-11097）。
		p.sysMsg(server, "清除玩家的仓库密码！！！")
		p.sysMsg(server, "命令格式: @ClrPassword 人物名称")
		return true
	}
	target := p.Engine.GetPlayerByName(params[0])
	if target == nil {
		return true // Delphi 目标不在线静默退出（ObjBase.pas:11100-11103）
	}
	target.StoragePassword = ""
	target.StoragePwdLocked = false
	target.storagePwdFail = 0
	target.sysMsgGreen(server, "你的保护密码已被清除！！！")
	p.sysMsgGreen(server, fmt.Sprintf("%s的保护密码已被清除！！！", params[0]))
	return true
}

// cmdRefineWeapon 调整手上武器属性（Delphi CmdRefineWeapon，
// 分发 ObjBase.pas:8244-8246，实现 14106-14124；命令表 REFINEWEAPON
// 'RefineWeapon' 权限 10，M2Share.pas:2858，ini 键 RefineWeapon 6349）。
//
// 编码与 Delphi TUserItem.btValue / Go upgrade.go 一致：
// BtValue[0..2]=DC/MC/SC 升级加点（upgrade.go CheckWeaponUpgradeStatus
// 成功后同样写入 0/1/2；RecalcAbilitys 消费，itemsystem.go:784-793），
// BtValue[5]=HIT（Delphi 武器 AC2=准确：GetItemAddValue 加 btValue[5]
// 到 AC2，ApplyItemParameters AC2→wHitPoint，ItmUnit.pas:125,561）。
// 注意：Go RecalcAbilitys 尚未消费武器 BtValue[5] 为 HIT
//（与 drops.go:57、npcscript.go UPGRADEITEMS 写入现状一致，属既有缺口；
// Delphi 原版该值也只作用于发送客户端的物品面板）。
//
// 升级 DB 记录：Go 武器升级仅在 NPC 处"寄存待取"阶段存在 DB 记录
//（upgrade.go SaveNpcUpgrade），彼时武器已不在身上；手上武器
// 的 BtValue 随角色存档持久化，无需额外同步。
func (p *PlayObject) cmdRefineWeapon(server *netserver.TCPServer, params []string) bool {
	// Delphi Str_ToInt(sParamN, 0)（ObjBase.pas:8246）。
	nDc, nMc, nSc, nHit := 0, 0, 0, 0
	if len(params) > 0 {
		nDc, _ = strconv.Atoi(params[0])
	}
	if len(params) > 1 {
		nMc, _ = strconv.Atoi(params[1])
	}
	if len(params) > 2 {
		nSc, _ = strconv.Atoi(params[2])
	}
	if len(params) > 3 {
		nHit, _ = strconv.Atoi(params[3])
	}
	if nDc+nMc+nSc > 10 {
		return true // Delphi 静默退出（ObjBase.pas:14113）
	}
	weapon := p.UseItems[protocol.UWeapon]
	if weapon == nil || weapon.WIndex <= 0 {
		return true // Delphi wIndex<=0 静默退出（ObjBase.pas:14114）
	}
	// Integer→Byte 隐式截断，与 Delphi 一致（负值/超界取低 8 位）。
	weapon.BtValue[0] = byte(nDc)
	weapon.BtValue[1] = byte(nMc)
	weapon.BtValue[2] = byte(nSc)
	weapon.BtValue[5] = byte(nHit)
	p.SendUseItemsFull(server) // Delphi SendUpdateItem（ObjBase.pas:14119）
	p.RecalcAbilitys()
	p.SendAbility(server)    // Delphi RM_ABILITY（ObjBase.pas:14121）
	p.SendSubAbility(server) // Delphi RM_SUBABILITY（ObjBase.pas:14122）
	// "MC" 后缺冒号为 Delphi 原文（ObjBase.pas:14123），照抄。
	log.Logf(log.LevelInfo, "GM", "[武器调整]%s DC:%d MC%d SC%d HIT:%d", p.Name, nDc, nMc, nSc, nHit)
	return true
}

// 自定义物品名内存登记表（Delphi ItemUnit.m_ItemNameList +
// AddCustomItemName/SaveCustomItemName，ItmUnit.pas:767-793）。
// Delphi 持久化到 CustomItemName 文件并在发送物品名时套用
//（GetItemName，ItmUnit.pas:65）；Go 未落盘、客户端显示未接入，
// 仅内存登记（重启即失）。
var (
	gmCustomItemNameMu sync.Mutex
	gmCustomItemNames  = make(map[int32]string) // MakeIndex -> 自定义名
)

// cmdChangeItemName 自定义物品名（Delphi CmdChangeItemName，
// 分发 ObjBase.pas:8382-8384，实现 10791-10817；命令表 CHANGEITEMNAME
// 'ChangeItemName' 表权限 10，M2Share.pas:2837，ini 键 ChangeItemName 6163；
// 实现内 m_btPermission<6 硬编码静默退出，ObjBase.pas:10795）。
// 参数：物品编号(MakeIndex) 物品ID号(wIndex) 物品名称。
func (p *PlayObject) cmdChangeItemName(server *netserver.TCPServer, params []string) bool {
	if len(params) < 3 {
		p.sysMsg(server, "命令格式: @ChangeItemName 物品编号 物品ID号 物品名称")
		return true
	}
	// Delphi Str_ToInt(..., -1)；makeIndex<=0 或 itemIndex<0 → 格式提示
	//（ObjBase.pas:10800-10807）。
	makeIndex, err1 := strconv.Atoi(params[0])
	itemIndex, err2 := strconv.Atoi(params[1])
	if err1 != nil || err2 != nil || makeIndex <= 0 || itemIndex < 0 {
		p.sysMsg(server, "命令格式: @ChangeItemName 物品编号 物品ID号 物品名称")
		return true
	}
	gmCustomItemNameMu.Lock()
	// Delphi AddCustomItemName 以 (MakeIndex, ItemIndex) 双键查重；
	// MakeIndex 为全局唯一实例 ID，Go 以 MakeIndex 单键等价实现。
	_, exists := gmCustomItemNames[int32(makeIndex)]
	if !exists {
		gmCustomItemNames[int32(makeIndex)] = params[2]
	}
	gmCustomItemNameMu.Unlock()
	if exists {
		p.sysMsg(server, "此物品，已经设置了其它的名称！！！")
		return true
	}
	// Delphi 成功即 SaveCustomItemName 落盘（ObjBase.pas:10811）；Go 仅内存。
	p.sysMsgGreen(server, "物品名称设置成功。")
	return true
}
