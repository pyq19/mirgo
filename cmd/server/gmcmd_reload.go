package main

import (
	"fmt"

	"github.com/pyq0109/mirgo/internal/netserver"
)

// GM 命令——重载族移植。
// Delphi 参照：ObjBase.pas（分发 7976-7990 常规段、8420-8477 调试段）+
// M2Share.pas g_GameCommand（2802-2804,2859-2869）与文案常量（3315-3321）。
// 注册顺序 = Delphi 分发链出现顺序。
// 未移植：ReloadRobot/ReloadRobotManage（ObjBase.pas:8435-8443，
// ObjRobot.pas:382-423）——Go 无机器人子系统（Robot.txt/TRobotObject），
// 无行为可闭环，与原版空壳同策略留档不实装。

func init() {
	// --- 常规 GM 段（ObjBase.pas:7976-7990，nPermissionMin 检查）---
	registerGMCommand(&gmCmdDef{key: "ReloadGuild", name: "ReloadGuild", permMin: 10, usage: "行会名称", handler: (*PlayObject).cmdReloadGuild})
	registerGMCommand(&gmCmdDef{key: "ReloadLineNotice", name: "ReloadLineNotice", permMin: 10, handler: (*PlayObject).cmdReloadLineNotice})
	registerGMCommand(&gmCmdDef{key: "ReloadAbuse", name: "ReloadAbuse", permMin: 10, handler: (*PlayObject).cmdReloadAbuse})
	// --- 调试段（ObjBase.pas:8420 起，内层 >=5/boTestServer 门槛内）---
	// ReloadAdmin 实现内另有 m_btPermission<6 静默退出（ObjBase.pas:14168）。
	registerGMCommand(&gmCmdDef{key: "ReloadAdmin", name: "ReloadAdmin", permMin: 6, permFail: permFailSilent, handler: (*PlayObject).cmdReloadAdmin})
	registerGMCommand(&gmCmdDef{key: "ReloadManage", name: "ReloadManage", permMin: 10, handler: (*PlayObject).cmdReloadManage})
	// ReloadMonItems 实现内另有 m_btPermission<6 静默退出（ObjBase.pas:14290）。
	registerGMCommand(&gmCmdDef{key: "ReloadMonItems", name: "ReloadMonItems", permMin: 6, permFail: permFailSilent, handler: (*PlayObject).cmdReloadMonItems})
	// ReloadItemDB/ReloadMonsterDB/ReLoadMinMap 为分发即实现（无权限检查），
	// 仅受调试段内层 >=5 门槛约束 → permMin 5 + permFailFallback。
	registerGMCommand(&gmCmdDef{key: "ReloadItemDB", name: "ReloadItemDB", permMin: 5, permFail: permFailFallback, handler: (*PlayObject).cmdReloadItemDB})
	registerGMCommand(&gmCmdDef{key: "ReloadMonsterDB", name: "ReloadMonsterDB", permMin: 5, permFail: permFailFallback, handler: (*PlayObject).cmdReloadMonsterDB})
	registerGMCommand(&gmCmdDef{key: "ReLoadMinMap", name: "ReLoadMinMap", permMin: 5, permFail: permFailFallback, handler: (*PlayObject).cmdReloadMinMap})
}

// cmdReloadGuild 从存档重载行会（Delphi CmdReloadGuild，
// ObjBase.pas:14174-14204 → Guild.LoadGuild，Guild.pas:429-436 读
// GuildBase/Guilds/<名称>.txt/.ini）。
// 偏差：Delphi 行会为文本文件、Go 为 SQLite，重载 = 以 DB 记录重建内存
// 行会（含成员/职位/公告/行会战/联盟）；DB 无记录时保留内存现状仍报成功
// （Delphi FindGuild 命中后 LoadGuildFile 缺失同样走成功分支）。
// 原版怪癖照抄：FindGuild 区分大小写（Guild.pas:171 '=' 比较）、
// 成功提示用 c_Red（g_sGameCommandReloadGuildSuccessMsg，M2Share.pas:3318）。
func (p *PlayObject) cmdReloadGuild(server *netserver.TCPServer, params []string) bool {
	if len(params) < 1 {
		p.gmParamHelp(server, "ReloadGuild", "行会名称")
		return true
	}
	e := p.Engine
	if e.FindGuild(params[0]) == nil {
		p.sysMsg(server, fmt.Sprintf("未找到行会%s！！！", params[0]))
		return true
	}
	if e.db != nil {
		if records, err := e.db.LoadGuilds(); err == nil {
			for _, rec := range records {
				if rec.Name == params[0] {
					fresh := guildFromRecord(rec)
					e.mu.Lock()
					for i, g := range e.Guilds {
						if g.Name == params[0] {
							e.Guilds[i] = fresh
							break
						}
					}
					e.mu.Unlock()
					break
				}
			}
		}
	}
	p.sysMsg(server, fmt.Sprintf("行会%s重加载成功...", params[0]))
	return true
}

// cmdReloadLineNotice 重载登录公告（Delphi CmdReloadLineNotice，
// ObjBase.pas:14212-14232 → LoadLineNotice，M2Share.pas:3426-3450 读
// Notice/LineNotice.txt 进滚动公告表）。
// 偏差：Go 无滚动公告系统，以登录公告（notice.jsonc）为最近等价物；
// Delphi 文件缺失报失败，Go loadNoticeText 回退默认文案恒成功，
// 失败分支不可达（文案常量保留备查，M2Share.pas:3321）。
func (p *PlayObject) cmdReloadLineNotice(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.gmParamHelp(server, "ReloadLineNotice", "")
		return true
	}
	p.Engine.SetNoticeText(loadNoticeText(p.Engine.configDir))
	p.sysMsgGreen(server, "重新加载公告设置信息完成。")
	return true
}

// cmdReloadAbuse 重载脏话过滤表（Delphi CmdReloadAbuse，
// ObjBase.pas:14151-14164）。
// Delphi 原版为空壳（权限+'?'帮助后无任何逻辑）；Go 按字段语义实装真重载
// （wordfilter.go 从 serverconfig/WordFilter.txt 重读过滤词表），
// 文案为 Go 自拟。
func (p *PlayObject) cmdReloadAbuse(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.gmParamHelp(server, "ReloadAbuse", "")
		return true
	}
	loadWordFilter(p.Engine.configDir)
	p.sysMsgGreen(server, "重新加载脏话过滤列表完成。")
	return true
}

// cmdReloadAdmin 重载管理员列表（Delphi CmdReLoadAdmin，
// ObjBase.pas:14166-14172 → FrmDB.LoadAdminList，LocalDB.pas:82-153 读
// Envir/AdminList.txt）。Go 对应 admin_list.jsonc 重读
// （adminlist.go ReloadAdminList）。
func (p *PlayObject) cmdReloadAdmin(server *netserver.TCPServer, params []string) bool {
	p.Engine.ReloadAdminList()
	p.sysMsgGreen(server, "管理员列表重新加载成功...")
	return true
}

// cmdReloadManage 重载管理/功能 NPC 脚本（Delphi CmdReloadManage，
// ObjBase.pas:14234-14268：无参=g_ManageNPC 'QManage' 登录脚本，
// 有参=g_FunctionNPC 'QFunction' 功能脚本，ClearScript+LoadNpcScript）。
// 偏差：Go 无 QManage/QFunction 全局单例（同 gotoNpcLabel 的说明），
// 退化为按类别作废全部 NPC 脚本缓存：无参=管理型（非商人），
// 有参=商人。Delphi 失败文案亦为 c_Green（怪癖），Go 恒成功。
func (p *PlayObject) cmdReloadManage(server *netserver.TCPServer, params []string) bool {
	if len(params) > 0 && params[0] != "" && params[0][0] == '?' {
		p.gmParamHelp(server, "ReloadManage", "")
		return true
	}
	merchant := len(params) > 0
	e := p.Engine
	e.mu.RLock()
	for _, npc := range e.Npcs {
		if npc.IsMerchant == merchant {
			npc.InvalidateScript()
		}
	}
	e.mu.RUnlock()
	if merchant {
		p.sysMsgGreen(server, "重新加载功能脚本完成...")
	} else {
		p.sysMsgGreen(server, "重新加载登录脚本完成...")
	}
	return true
}

// cmdReloadMonItems 重载怪物爆物品表（Delphi CmdReloadMonItems，
// ObjBase.pas:14288-14304 → FrmDB.LoadMonitems 逐怪读
// Envir/MonItems/<怪名>.txt）。Go 对应 mon_items/*.jsonc 整表重读。
// 原版怪癖照抄：失败文案亦 c_Green（M2Share.pas 无失败常量，原文即绿）。
func (p *PlayObject) cmdReloadMonItems(server *netserver.TCPServer, params []string) bool {
	e := p.Engine
	tables := LoadDropTables(e.dropTablesDir)
	if tables == nil {
		p.sysMsgGreen(server, "怪物爆物品列表重加载失败！！！")
		return true
	}
	e.mu.Lock()
	e.DropTables = tables
	e.mu.Unlock()
	p.sysMsgGreen(server, "怪物爆物品列表重加载完成...")
	return true
}

// cmdReloadItemDB 重载物品数据库（Delphi 分发即实现，ObjBase.pas:8454-8458
// → FrmDB.LoadItemsDB，LocalDB.pas:250 起 SQL 重建 STDITEMS）。
// Go 从 std_items.jsonc 重读（itemdb.go LoadItemDB）。
// 偏差：Delphi SQL 层无失败分支，Go 读文件可失败，补红字失败文案。
func (p *PlayObject) cmdReloadItemDB(server *netserver.TCPServer, params []string) bool {
	e := p.Engine
	db, err := LoadItemDB(e.itemDBPath)
	if err != nil {
		p.sysMsg(server, "物品数据库重新加载失败！！！")
		return true
	}
	e.mu.Lock()
	e.ItemDB = db
	e.mu.Unlock()
	p.sysMsgGreen(server, "物品数据库重新加载完成。")
	return true
}

// cmdReloadMonsterDB 重载怪物数据库（Delphi 分发即实现，
// ObjBase.pas:8466-8470 → FrmDB.LoadMonsterDB，LocalDB.pas:1288 起）。
// Go 从 monster_db.jsonc 重读（monsterdb.go LoadMonsterDB）。
func (p *PlayObject) cmdReloadMonsterDB(server *netserver.TCPServer, params []string) bool {
	e := p.Engine
	db, err := LoadMonsterDB(e.monsterDBPath)
	if err != nil {
		p.sysMsg(server, "怪物数据库重新加载失败！！！")
		return true
	}
	e.mu.Lock()
	e.MonsterDB = db
	e.mu.Unlock()
	p.sysMsgGreen(server, "怪物数据库重新加载完成。")
	return true
}

// cmdReloadMinMap 重载小地图配置（Delphi 分发即实现但为空壳，
// ObjBase.pas:8472-8477：LoadMinMap/ReSetMinMap 两行全被注释，仅发
// '小地图配置重新加载完成。'）。Go 有 mini_map.jsonc 基建，实装真重载。
func (p *PlayObject) cmdReloadMinMap(server *netserver.TCPServer, params []string) bool {
	p.MapMgr.InitMiniMaps(p.Engine.configDir)
	p.sysMsgGreen(server, "小地图配置重新加载完成。")
	return true
}
