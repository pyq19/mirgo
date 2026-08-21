package main

import (
	"net"
	"os"
	"strings"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
)

// AdminList：Delphi Envir/AdminList.txt 的移植。
// 加载 = TFrmDB.LoadAdminList（LocalDB.pas:82-154），
// 授权+IP 校验 = TPlayObject.UserLogon（ObjBase.pas:16384-16392）+
// GetHumPermission（UsrEngn.pas:2174-2197）+ CompareIPaddr（M2Share.pas:11036）。
// 数据来源 serverconfig/guards/admin_list.jsonc（serverconfig 转换工具生成）。

// AdminInfo 管理员条目（Delphi TAdminInfo）。
type AdminInfo struct {
	Name  string `json:"name"`  // 角色名
	Level int    `json:"level"` // 权限等级 1-10（AdminList.txt 首字符标记）
	IP    string `json:"ip"`    // 允许登录的 IP，'*' 开头通配
}

func loadAdminList(path string) []AdminInfo {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var doc struct {
		Admins []AdminInfo `json:"admins"`
	}
	if err := parseJSONC(data, &doc); err != nil {
		log.Logf(log.LevelWarn, "Server", "failed to parse admin list %s: %v", path, err)
		return nil
	}
	// Delphi WL 版行为：名字或 IP 为空的条目跳过（LocalDB.pas:135-137）。
	list := make([]AdminInfo, 0, len(doc.Admins))
	for _, a := range doc.Admins {
		if a.Name == "" || a.IP == "" {
			continue
		}
		if a.Level < 1 || a.Level > 10 {
			a.Level = 10
		}
		list = append(list, a)
	}
	log.Logf(log.LevelInfo, "Server", "admin list loaded: %d entries", len(list))
	return list
}

// compareIPaddr Delphi CompareIPaddr（M2Share.pas:11036-11056）：
// 任一为空不匹配；目标 '*' 开头全匹配；含 '*' 前缀匹配；否则忽略大小写相等。
func compareIPaddr(src, dst string) bool {
	if src == "" || dst == "" {
		return false
	}
	if dst[0] == '*' {
		return true
	}
	if i := strings.IndexByte(dst, '*'); i > 0 {
		return len(src) >= i && strings.EqualFold(src[:i], dst[:i])
	}
	return strings.EqualFold(src, dst)
}

// sessionRemoteIP 取会话客户端 IP（不含端口）。
func sessionRemoteIP(session *netserver.Session) string {
	if session == nil || session.Conn == nil {
		return ""
	}
	host, _, err := net.SplitHostPort(session.Conn.RemoteAddr().String())
	if err != nil {
		return session.Conn.RemoteAddr().String()
	}
	return host
}

// getAdminList 返回管理员列表快照（mu 保护，登录校验用）。
func (e *UserEngine) getAdminList() []AdminInfo {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]AdminInfo(nil), e.AdminList...)
}

// ReloadAdminList 重读 admin_list.jsonc（Delphi @ReloadAdmin，
// ObjBase.pas:14166-14173；跨服消息 213 Go 无跨服架构跳过）。
func (e *UserEngine) ReloadAdminList() int {
	list := loadAdminList(e.adminListPath)
	e.mu.Lock()
	e.AdminList = list
	e.mu.Unlock()
	return len(list)
}
