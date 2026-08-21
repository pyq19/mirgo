package main

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/storage"
)

type UserEngine struct {
	PlayObjectList map[int32]*PlayObject
	mu             sync.RWMutex

	Config *ServerConfig

	db     *storage.Database
	mapMgr *MapManager
	ItemDB *ItemDB

	MonsterDB    *MonsterDB
	DropTables   *DropTable
	monGenPath   string
	npcConfigDir string

	Monsters      []*MonsterObject
	Npcs          []*NpcObject
	MonGenList    []MonGenEntry
	nextMonsterID int32
	nextItemID    int32
	currMonGen    int // round-robin 刷怪器索引

	MagicDB *MagicDB
	Parties map[int32]*Party
	Guilds  []*Guild
	Castle  *CastleObject

	NoMonGen bool

	// GM 封禁/发言过滤/怪物攻城（Delphi g_DenyIPLogonList/g_DenyAccountLogonList/
	// g_DenyCharNameLogonList、g_DisableSendMsgList、Mission 集中点，
	// ObjBase.pas:12466-12760、12762-12807、13593-13614）。受 mu 保护。
	DenyIPList      []DenyEntry
	DenyAccountList []DenyEntry
	DenyCharList    []DenyEntry
	DisableSendMsg  []string
	MissionMap      string
	MissionX        int
	MissionY        int

	// denyListPath 永久封禁落盘路径（LoadDenyLists 设置）。
	// Delphi 写 Envir/DenyIPAddrList.txt 等三个文本文件
	//（M2Share.pas:4755-4972）；Go 合并为单个 JSON，见 LoadDenyLists。
	denyListPath string

	// AdminList 管理员列表（Delphi Envir/AdminList.txt，@ReloadAdmin 可重读）。
	// 受 mu 保护。
	AdminList     []AdminInfo
	adminListPath string

	// 重载族（@ReloadItemDB/@ReloadMonsterDB/@ReloadMonItems/@ReloadLineNotice）
	// 所需的配置路径，main.go 启动时赋值。
	configDir     string
	itemDBPath    string
	monsterDBPath string
	dropTablesDir string

	// NoticeText 登录公告（Delphi LineNoticeList，M2Share.pas:3426-3450
	// LoadLineNotice；Go 无滚动公告，以登录公告为最近等价物）。
	// 受 mu 保护，@ReloadLineNotice 热重载。
	NoticeText string

	// SpiritMutinyTick 祈祷生效宝宝叛变到期时间戳（Delphi
	// g_dwSpiritMutinyTick，M2Share.pas:1509；@祈祷生效/@停止叛变 写入，
	// monsterobject.go slaveFollow 消费）。受 mu 保护。
	SpiritMutinyTick int64
}

// GetNoticeText 登录公告读取（mu 保护）。
func (e *UserEngine) GetNoticeText() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.NoticeText
}

// SetNoticeText 登录公告写入（mu 保护，@ReloadLineNotice）。
func (e *UserEngine) SetNoticeText(text string) {
	e.mu.Lock()
	e.NoticeText = text
	e.mu.Unlock()
}

// DenyEntry 封禁条目（Delphi pTDenyIPLogon 等，永久标志决定是否落盘）。
type DenyEntry struct {
	Value     string
	Permanent bool
}

func NewUserEngine(db *storage.Database, mapMgr *MapManager) *UserEngine {
	return &UserEngine{
		PlayObjectList: make(map[int32]*PlayObject),
		db:             db,
		mapMgr:         mapMgr,
		nextMonsterID:  100000,
		nextItemID:     200000,
		Parties:        make(map[int32]*Party),
	}
}

// DB 返回底层存储（城堡等子系统持久化用）。
func (e *UserEngine) DB() *storage.Database {
	return e.db
}

// allocItemID 分配全局唯一物品实例 ID。
func (e *UserEngine) allocItemID() int32 {
	e.mu.Lock()
	defer e.mu.Unlock()
	id := e.nextItemID
	e.nextItemID++
	return id
}

func (e *UserEngine) AddPlayer(player *PlayObject) {
	e.mu.Lock()
	e.PlayObjectList[player.ID] = player
	e.mu.Unlock()
	log.Logf(log.LevelInfo, "UserEngine", "player %s joined (total: %d)", player.Name, len(e.PlayObjectList))
}

func (e *UserEngine) RemovePlayer(id int32) {
	e.mu.Lock()
	delete(e.PlayObjectList, id)
	e.mu.Unlock()
	log.Logf(log.LevelInfo, "UserEngine", "player %d removed (total: %d)", id, len(e.PlayObjectList))
}

func (e *UserEngine) GetPlayer(id int32) *PlayObject {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.PlayObjectList[id]
}

func (e *UserEngine) GetMonster(id int32) *MonsterObject {
	for _, mon := range e.Monsters {
		if mon.ID == id && !mon.Ghost {
			return mon
		}
	}
	return nil
}

func (e *UserEngine) clearMapMonsters(env *Environment) {
	for _, mon := range e.Monsters {
		if mon.envir == env && !mon.Death && !mon.Ghost && !mon.IsSafeZoneGuard {
			mon.Death = true
			mon.DeathTick = time.Now().UnixMilli()
			mon.WAbil.HP = 0
		}
	}
}

func (e *UserEngine) ProcessHumans(server *netserver.TCPServer) {
	e.mu.RLock()
	players := make([]*PlayObject, 0, len(e.PlayObjectList))
	for _, p := range e.PlayObjectList {
		players = append(players, p)
	}
	e.mu.RUnlock()

	for _, player := range players {
		player.Operate(server)
	}
}

func (e *UserEngine) ProcessDoors(server *netserver.TCPServer, currentTick int64) {
	e.mapMgr.mu.RLock()
	defer e.mapMgr.mu.RUnlock()
	for _, env := range e.mapMgr.maps {
		ProcessDoors(env, currentTick, e.Config, server)
	}
}

func (e *UserEngine) ProcessEvents(server *netserver.TCPServer, now int64) {
	e.mapMgr.mu.RLock()
	defer e.mapMgr.mu.RUnlock()
	for _, env := range e.mapMgr.maps {
		env.ProcessMapEvents(server, now)
	}
}

func (e *UserEngine) ProcessMineRegen(now int64) {
	e.mapMgr.ProcessMineRegen(now)
}

func (e *UserEngine) GetPlayerCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.PlayObjectList)
}

func (e *UserEngine) GetPlayerByName(name string) *PlayObject {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, p := range e.PlayObjectList {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (e *UserEngine) SaveAllPlayers(db *storage.Database) {
	e.mu.RLock()
	players := make([]*PlayObject, 0, len(e.PlayObjectList))
	for _, p := range e.PlayObjectList {
		players = append(players, p)
	}
	e.mu.RUnlock()
	for _, p := range players {
		saveCharacterData(db, p)
	}
}

func (e *UserEngine) ProcessNpcs() {
	for _, npc := range e.Npcs {
		if npc.IsMerchant && len(npc.RefillConfig) > 0 {
			npc.RefillGoods(e.ItemDB, e)
		}
		npc.SaveData(e.db)
	}
}

func (e *UserEngine) ProcessNpcIdle() {
	for _, npc := range e.Npcs {
		if npc.IsMerchant {
			// 城堡战 NPC 隐藏（Delphi ObjNpc.pas:1642-1657）
			if npc.Castle && e.Castle != nil && e.Castle.IsAtWar() {
				if !npc.FixedHideMode {
					npc.SendRefMsg(RM_DISAPPEAR, 0, 0, 0, "")
					npc.FixedHideMode = true
				}
			} else {
				if npc.FixedHideMode {
					npc.FixedHideMode = false
					npc.SendRefMsg(RM_HIT, npc.Dir, npc.CurrX, npc.CurrY, "")
				}
			}
			npc.idleAnimate()
		}
	}
}

func (e *UserEngine) InvalidateAllNpcScripts() {
	for _, npc := range e.Npcs {
		npc.InvalidateScript()
	}
	log.Logf(log.LevelInfo, "NPC", "all NPC scripts invalidated (%d NPCs)", len(e.Npcs))
}

// awardExpForMonster — 怪物被毒杀时，将经验奖励给 LastHiterID。
func (e *UserEngine) awardExpForMonster(mon *MonsterObject, server *netserver.TCPServer) {
	if mon.LastHiterID == 0 {
		return
	}
	killer := e.GetPlayer(mon.LastHiterID)
	if killer == nil || killer.Ghost || killer.Death {
		return
	}
	killer.awardExp(server, mon)
}

// CountMapHumans 统计指定地图的在线玩家数（Delphi ObjNpc.pas:11131-11153）。
func (e *UserEngine) CountMapHumans(mapName string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	count := 0
	for _, p := range e.PlayObjectList {
		if p.MapName == mapName && !p.Ghost {
			count++
		}
	}
	return count
}

// CountMapMonsters 统计指定地图的怪物数（Delphi ObjNpc.pas:11157-11182）。
func (e *UserEngine) CountMapMonsters(mapName string) int {
	count := 0
	for _, mon := range e.Monsters {
		if mon.MapName == mapName && !mon.Death && !mon.Ghost {
			count++
		}
	}
	return count
}

// denyListFile 永久封禁落盘格式（只存 Permanent 条目）。
// Go 方案：Delphi 把三张永久封禁表分别写入 Envir/DenyIPAddrList.txt、
// DenyAccountList.txt、DenyChrNameList.txt（纯文本每行一条，
// SaveDenyIPAddrList 等只写 Objects[i]<>0 的永久条目，
// M2Share.pas:4803/4948/4876），启动时 LoadDeny* 读回
// （svMain.pas:614-616）；Go 用单个 JSON 文件等价替代。
type denyListFile struct {
	IP      []string `json:"ip"`
	Account []string `json:"account"`
	Char    []string `json:"char"`
}

// LoadDenyLists 启动时加载永久封禁表（Delphi LoadDenyIPAddrList/
// LoadDenyAccountList/LoadDenyChrNameList，M2Share.pas:4755/4900/4828）。
// 文件不存在按空表处理（Delphi 则创建空文件，Go 推迟到首次保存）。
func (e *UserEngine) LoadDenyLists(path string) {
	e.mu.Lock()
	e.denyListPath = path
	e.mu.Unlock()
	data, err := os.ReadFile(path)
	if err != nil {
		log.Logf(log.LevelInfo, "UserEngine", "deny lists: %s absent, starting empty", path)
		return
	}
	var doc denyListFile
	if err := json.Unmarshal(data, &doc); err != nil {
		log.Logf(log.LevelWarn, "UserEngine", "deny lists: parse %s failed: %v", path, err)
		return
	}
	e.mu.Lock()
	for _, v := range doc.IP {
		if v = strings.TrimSpace(v); v != "" {
			e.DenyIPList = append(e.DenyIPList, DenyEntry{Value: v, Permanent: true})
		}
	}
	for _, v := range doc.Account {
		if v = strings.TrimSpace(v); v != "" {
			e.DenyAccountList = append(e.DenyAccountList, DenyEntry{Value: v, Permanent: true})
		}
	}
	for _, v := range doc.Char {
		if v = strings.TrimSpace(v); v != "" {
			e.DenyCharList = append(e.DenyCharList, DenyEntry{Value: v, Permanent: true})
		}
	}
	nIP, nAcc, nChar := len(e.DenyIPList), len(e.DenyAccountList), len(e.DenyCharList)
	e.mu.Unlock()
	log.Logf(log.LevelInfo, "UserEngine", "deny lists loaded: %d IP, %d account, %d char", nIP, nAcc, nChar)
}

// SaveDenyLists 把三张封禁表的永久条目写入 LoadDenyLists 记录的路径
// （Delphi SaveDenyIPAddrList/SaveDenyAccountList/SaveDenyChrNameList，
// M2Share.pas:4803/4948/4876）。未加载过（无路径）时静默跳过。
func (e *UserEngine) SaveDenyLists() {
	e.mu.RLock()
	path := e.denyListPath
	doc := denyListFile{IP: []string{}, Account: []string{}, Char: []string{}}
	for _, d := range e.DenyIPList {
		if d.Permanent {
			doc.IP = append(doc.IP, d.Value)
		}
	}
	for _, d := range e.DenyAccountList {
		if d.Permanent {
			doc.Account = append(doc.Account, d.Value)
		}
	}
	for _, d := range e.DenyCharList {
		if d.Permanent {
			doc.Char = append(doc.Char, d.Value)
		}
	}
	e.mu.RUnlock()
	if path == "" {
		return
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		log.Logf(log.LevelWarn, "UserEngine", "deny lists: marshal failed: %v", err)
		return
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Logf(log.LevelWarn, "UserEngine", "deny lists: write %s failed: %v", path, err)
		return
	}
	log.Logf(log.LevelInfo, "UserEngine", "deny lists saved: %d IP, %d account, %d char",
		len(doc.IP), len(doc.Account), len(doc.Char))
}

// denyEntryHit 封禁表命中查询（Delphi GetDenyIPaddrList/GetDenyAccountList/
// GetDenyChrNameList，M2Share.pas:4784/4929/4857：CompareText 不敏感）。
func denyEntryHit(list []DenyEntry, value string) bool {
	for _, d := range list {
		if strings.EqualFold(d.Value, value) {
			return true
		}
	}
	return false
}

// CheckDenyLogon 登录封禁拦截判定（Delphi TPlayObject.CheckDenyLogon，
// ObjBase.pas:26763-26783，UserLogon 内调用，16610）：按 IP→帐号→角色名
// 顺序检查三张封禁表，命中返回对应拒绝文案（g_sYourIPaddrDenyLogon/
// g_sYourAccountDenyLogon/g_sYourCharNameDenyLogon，M2Share.pas:3187-3189），
// 未命中返回空串。纯函数以便单测。
func (e *UserEngine) CheckDenyLogon(ip, account, charName string) string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if denyEntryHit(e.DenyIPList, ip) {
		return "你当前登录的IP地址已被禁止登录了！！！"
	}
	if denyEntryHit(e.DenyAccountList, account) {
		return "你当前登录的帐号已被禁止登录了！！！"
	}
	if denyEntryHit(e.DenyCharList, charName) {
		return "你当前登录的人物已被禁止登录了！！！"
	}
	return ""
}

// isMsgDisabled 禁言名单查询（Delphi GetDisableSendMsgList，
// M2Share.pas:4678-4691：CompareText 不敏感）。
func (e *UserEngine) isMsgDisabled(name string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, n := range e.DisableSendMsg {
		if strings.EqualFold(n, name) {
			return true
		}
	}
	return false
}
