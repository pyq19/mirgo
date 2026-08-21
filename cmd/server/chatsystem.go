package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// 聊天消息 Param 颜色：低字节=前景、高字节=背景的调色板索引
// （客户端 GetRGB ≡ Prguse MainPalette，ClMain.pas:392-398,996）。
// 默认值取自 Delphi M2Share.pas:1838-1857；发送端填色见
// ObjBase.pas:5662-5674。组队色为预留（Go 服务端暂未实现组队聊天）。
// 有意偏差：Delphi 把 cry 并入 SM_HEAR 通道、组队并入 SM_SYSMESSAGE
// 通道（ObjBase.pas:5665,5671，纯靠颜色区分）；Go 保留独立消息 ID，
// 仅对齐颜色下发。
const (
	chatColorHear    = 0xFF<<8 | 0x00 // 附近：黑字亮底
	chatColorCry     = 0x97<<8 | 0x00 // 喊话：黑字青底
	chatColorWhisper = 0xFF<<8 | 0xFC // 私聊：浅字亮底
	chatColorSys     = 0x38<<8 | 0xFF // 系统（红系默认）：白字深底
	chatColorGuild   = 0xFF<<8 | 0xDB // 行会：绿字亮底
	chatColorGroup   = 0xFF<<8 | 0xC4 // 组队：预留
	// SysMsg 颜色（Delphi btBlueMsg/btGreenMsg 默认值，M2Share.pas:1846-1850）
	chatColorSysBlue  = 0xFC<<8 | 0xFF // c_Blue
	chatColorSysGreen = 0xFF<<8 | 0xDB // c_Green
)

func (p *PlayObject) HandleSay(msg SendMessage, server *netserver.TCPServer) {
	text := msg.Msg
	if text == "" {
		return
	}

	if p.ShutupTick > 0 && time.Now().UnixMilli() < p.ShutupTick {
		p.sysMsg(server, "你已被禁言")
		return
	}

	if strings.HasPrefix(text, "@") {
		p.HandleGMCommand(text[1:], server)
		return
	}

	if strings.HasPrefix(text, "/") {
		p.handleSlashCmd(text[1:], server)
		return
	}

	// 敏感词过滤（Delphi RunGate FilterSayMsg 网关层语义，Main.pas:554-579）
	text = filterChatText(text)

	// 发言过滤（Delphi m_boFilterSendMsg，ObjBase.pas:704 注释"发的文字只能
	// 自己看到"，@DisableSendMsg 设置）：位于 @ 命令与 / 命令分发之后，
	// GM 命令不受过滤影响；其后附近发言（8730-8735）、喊话（8713-8718）、
	// 行会消息均不再广播，仅以 SM_HEAR 回显给自己。
	// 偏差：Delphi 编组('!!')/行会('!~')分支未查 m_boFilterSendMsg，
	// Go 按字段语义一并收紧。
	if p.FilterSendMsg {
		hearMsg := protocol.MakeDefaultMsg(protocol.SMHear, p.ID, chatColorHear, 0, 0)
		server.Send(p.Session.ID, hearMsg, protocol.EncodeString(p.Name+"/"+text))
		return
	}

	if strings.HasPrefix(text, "!~") {
		if guildText := text[2:]; guildText != "" {
			gmsg := msg
			gmsg.Msg = guildText
			p.HandleGuildMessage(gmsg, server)
		}
		return
	}

	if strings.HasPrefix(text, "!") {
		if cryText := text[1:]; cryText != "" {
			p.HandleCry(cryText, server)
		}
		return
	}

	if p.envir == nil {
		return
	}

	objs := p.envir.GetRangeObjects(p.CurrX, p.CurrY, viewRange)
	for _, obj := range objs {
		other, ok := obj.(*PlayObject)
		if !ok || other.Ghost {
			continue
		}
		hearMsg := protocol.MakeDefaultMsg(protocol.SMHear, p.ID, chatColorHear, 0, 0)
		body := protocol.EncodeString(p.Name + "/" + text)
		server.Send(other.Session.ID, hearMsg, body)
	}
	log.Logf(log.LevelInfo, "Chat", "%s: %s", p.Name, text)
}

// handleSlashCmd 私聊前缀 "/" 的特殊命令（Delphi ProcessSayMsg，
// ObjBase.pas:8601-8676）：/who 在线人数（权限3）、/total 总在线（权限5），
// 其余 "/名字 内容" 为私聊——Go 未实装私聊通道，静默吞掉。
func (p *PlayObject) handleSlashCmd(text string, server *netserver.TCPServer) {
	cmd := strings.ToLower(strings.TrimSpace(text))
	if cmd != "who" && cmd != "total" {
		return
	}
	whoPerm, totalPerm := 3, 5
	if p.Engine != nil && p.Engine.Config != nil {
		if v, ok := p.Engine.Config.Commands.Permissions["Who"]; ok {
			whoPerm = v
		}
		if v, ok := p.Engine.Config.Commands.Permissions["Total"]; ok {
			totalPerm = v
		}
	}
	count := 0
	if p.Engine != nil {
		count = p.Engine.GetPlayerCount()
	}
	if cmd == "who" {
		if int(p.Permission) < whoPerm {
			p.sysMsg(server, gmMsgPermTooLow)
			return
		}
		p.sysMsg(server, fmt.Sprintf("在线人数: %d", count)) // g_sOnlineCountMsg（M2Share.pas:3230）
		return
	}
	if int(p.Permission) < totalPerm {
		p.sysMsg(server, gmMsgPermTooLow)
		return
	}
	// Go 单进程单服，总在线=本服在线（Delphi 跨服统计无对应）。
	p.sysMsg(server, fmt.Sprintf("总在线人数: %d", count)) // g_sTotalOnlineCountMsg（M2Share.pas:3231）
}

// HandleCry 喊话：广播到同地图所有玩家。
func (p *PlayObject) HandleCry(text string, server *netserver.TCPServer) {
	if p.Engine == nil {
		return
	}
	cryMsg := protocol.MakeDefaultMsg(protocol.SMCry, p.ID, chatColorCry, 0, 0)
	body := protocol.EncodeString(p.Name + "/" + text)
	p.Engine.mu.RLock()
	for _, other := range p.Engine.PlayObjectList {
		if other.Ghost || other.envir != p.envir {
			continue
		}
		// Delphi CryCry 接收方检查（UsrEngn.pas:1516-1522）：
		// m_boBanShout=FALSE（@LetShout 关闭群聊，反义命名 TRUE=允许）
		// 的玩家收不到喊话。
		if !other.BanShout {
			continue
		}
		server.Send(other.Session.ID, cryMsg, body)
	}
	p.Engine.mu.RUnlock()
	log.Logf(log.LevelInfo, "Chat", "%s cries: %s", p.Name, text)
}

// HandleWhisper 私聊：body 格式为 "目标名/消息内容"。
func (p *PlayObject) HandleWhisper(msg SendMessage, server *netserver.TCPServer) {
	body := msg.Msg
	slashIdx := strings.IndexByte(body, '/')
	if slashIdx <= 0 || slashIdx >= len(body)-1 {
		return
	}
	targetName := body[:slashIdx]
	text := filterChatText(body[slashIdx+1:])

	// 禁言玩家不能发起私聊（Delphi ProcessSayMsg '/' 分支：
	// if not m_boFilterSendMsg then Whisper(...)，ObjBase.pas:8674-8676，
	// 静默丢弃）。
	if p.FilterSendMsg {
		return
	}

	if p.Engine == nil {
		return
	}
	target := p.Engine.GetPlayerByName(targetName)
	if target == nil || target.Ghost {
		p.sysMsg(server, "对方不在线")
		return
	}
	// Delphi Whisper 转发检查（ObjBase.pas:2181-2185）：目标关闭私聊
	//（m_boHearWhisper=false，@AllowMsg）或已拉黑发送者（IsBlockWhisper，
	// ObjBase.pas:2217-2230，@PrvMsg）→ 不转发，向发送者反馈
	// g_sUserDenyWhisperMsg='  拒绝私聊...'（M2Share.pas:3192）。
	if !target.HearWhisper || target.isBlockWhisper(p.Name) {
		p.sysMsg(server, target.Name+"  拒绝私聊...")
		return
	}

	whisperMsg := protocol.MakeDefaultMsg(protocol.SMWhisper, p.ID, chatColorWhisper, 0, 0)
	// 发给目标：显示发送者名字
	server.Send(target.Session.ID, whisperMsg, protocol.EncodeString(p.Name+"/"+text))
	// 发回给自己：显示目标名字
	server.Send(p.Session.ID, whisperMsg, protocol.EncodeString(targetName+"/"+text))
	// GM 侦听私聊转发（Delphi m_GetWhisperHuman，ObjBase.pas:2190-2202）。
	p.forwardWhisperWatchers(server, target, text)
	log.Logf(log.LevelInfo, "Chat", "%s whispers to %s: %s", p.Name, targetName, text)
}

type Party struct {
	Leader  int32
	Members []int32
}

// partyOf 返回玩家所在的队伍（不在任何队伍中则返回 nil）。
func (p *PlayObject) partyOf() *Party {
	if p.Engine == nil {
		return nil
	}
	for _, party := range p.Engine.Parties {
		for _, id := range party.Members {
			if id == p.ID {
				return party
			}
		}
	}
	return nil
}

// broadcastPartyMembers 向所有成员发送完整成员名列表
// （客户端 GroupMembers 由该消息重建；队长在前，与 Delphi 顺序一致）。
func (p *PlayObject) broadcastPartyMembers(party *Party, server *netserver.TCPServer) {
	var sb strings.Builder
	for i, id := range party.Members {
		if i > 0 {
			sb.WriteByte('\n')
		}
		name := "?"
		if m := p.Engine.GetPlayer(id); m != nil {
			name = m.Name
		}
		sb.WriteString(name)
	}
	body := protocol.EncodeString(sb.String())
	for _, id := range party.Members {
		if m := p.Engine.GetPlayer(id); m != nil {
			resp := protocol.MakeDefaultMsg(protocol.SMGroupMembers, party.Leader, 0, 0, 0)
			server.Send(m.Session.ID, resp, body)
		}
	}
}

// HandleCreateGroup 创建队伍，并可选邀请指定玩家
// （Delphi SendCreateGroup 携带名称，FState:5523-5536）。
func (p *PlayObject) HandleCreateGroup(msg SendMessage, server *netserver.TCPServer) {
	if p.partyOf() != nil || p.Engine == nil {
		resp := protocol.MakeDefaultMsg(protocol.SMCreateGroupFail, 0, 0, 0, 0)
		server.Send(p.Session.ID, resp, "")
		return
	}
	party := &Party{
		Leader:  p.ID,
		Members: []int32{p.ID},
	}
	if targetName := strings.TrimSpace(msg.Msg); targetName != "" {
		target := p.Engine.GetPlayerByName(targetName)
		if target == nil || !target.AllowGroup || target.partyOf() != nil {
			resp := protocol.MakeDefaultMsg(protocol.SMCreateGroupFail, 0, 0, 0, 0)
			server.Send(p.Session.ID, resp, "")
			return
		}
		party.Members = append(party.Members, target.ID)
	}
	p.Engine.Parties[p.ID] = party
	resp := protocol.MakeDefaultMsg(protocol.SMCreateGroupOK, 0, 0, 0, 0)
	server.Send(p.Session.ID, resp, "")
	p.broadcastPartyMembers(party, server)
	log.Logf(log.LevelInfo, "Party", "%s created a party", p.Name)
}

// HandleGroupMode 切换允许组队标志（CMGroupMode，Param=1/0）。
func (p *PlayObject) HandleGroupMode(msg SendMessage, server *netserver.TCPServer) {
	p.AllowGroup = msg.Param1 != 0
	resp := protocol.MakeDefaultMsg(protocol.SMGroupModeChanged, int32(msg.Param1), 0, 0, 0)
	server.Send(p.Session.ID, resp, "")
}

// HandleAddGroupMember 邀请指定玩家加入调用者的队伍。
func (p *PlayObject) HandleAddGroupMember(msg SendMessage, server *netserver.TCPServer) {
	party := p.partyOf()
	if party == nil || party.Leader != p.ID || p.Engine == nil {
		resp := protocol.MakeDefaultMsg(protocol.SMGroupAddMemFail, 0, 0, 0, 0)
		server.Send(p.Session.ID, resp, "")
		return
	}
	targetName := strings.TrimSpace(msg.Msg)
	target := p.Engine.GetPlayerByName(targetName)
	if target == nil || !target.AllowGroup || target.partyOf() != nil || len(party.Members) >= p.Engine.Config.GetGroupMembersMax() {
		resp := protocol.MakeDefaultMsg(protocol.SMGroupAddMemFail, 0, 0, 0, 0)
		server.Send(p.Session.ID, resp, "")
		return
	}
	party.Members = append(party.Members, target.ID)
	resp := protocol.MakeDefaultMsg(protocol.SMGroupAddMemOK, 0, 0, 0, 0)
	server.Send(p.Session.ID, resp, protocol.EncodeString(target.Name))
	p.broadcastPartyMembers(party, server)
	log.Logf(log.LevelInfo, "Party", "%s invited %s to party", p.Name, target.Name)
}

// HandleDelGroupMember 将指定玩家从调用者的队伍中移除；
// 只剩一人时队伍解散。
func (p *PlayObject) HandleDelGroupMember(msg SendMessage, server *netserver.TCPServer) {
	party := p.partyOf()
	if party == nil || party.Leader != p.ID {
		resp := protocol.MakeDefaultMsg(protocol.SMGroupDelMemFail, 0, 0, 0, 0)
		server.Send(p.Session.ID, resp, "")
		return
	}
	targetName := strings.TrimSpace(msg.Msg)
	idx := -1
	for i, id := range party.Members {
		if m := p.Engine.GetPlayer(id); m != nil && m.Name == targetName {
			idx = i
			break
		}
	}
	if idx < 0 {
		resp := protocol.MakeDefaultMsg(protocol.SMGroupDelMemFail, 0, 0, 0, 0)
		server.Send(p.Session.ID, resp, "")
		return
	}
	party.Members = append(party.Members[:idx], party.Members[idx+1:]...)
	resp := protocol.MakeDefaultMsg(protocol.SMGroupDelMemOK, 0, 0, 0, 0)
	server.Send(p.Session.ID, resp, protocol.EncodeString(targetName))
	if removed := p.Engine.GetPlayerByName(targetName); removed != nil {
		cancel := protocol.MakeDefaultMsg(protocol.SMGroupCancel, 0, 0, 0, 0)
		server.Send(removed.Session.ID, cancel, "")
	}
	if len(party.Members) <= 1 {
		delete(p.Engine.Parties, party.Leader)
		cancel := protocol.MakeDefaultMsg(protocol.SMGroupCancel, 0, 0, 0, 0)
		server.Send(p.Session.ID, cancel, "")
		return
	}
	p.broadcastPartyMembers(party, server)
	log.Logf(log.LevelInfo, "Party", "%s removed %s from party", p.Name, targetName)
}
