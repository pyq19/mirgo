package main

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// TestInitMiniMaps 验证 mini_map.jsonc 的加载与 Environment.MinMap 写入
// （含 JSONC 注释剥离、未知地图跳过）。
func TestInitMiniMaps(t *testing.T) {
	tmpDir := t.TempDir()
	mapsDir := filepath.Join(tmpDir, "maps")
	if err := os.MkdirAll(mapsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `// 小地图映射（测试）
{
  "miniMaps": [
    {"mapName": "0", "miniMapId": 5},
    {"mapName": "D001", "miniMapId": 1},
    {"mapName": "missing", "miniMapId": 9}
  ]
}`
	if err := os.WriteFile(filepath.Join(mapsDir, "mini_map.jsonc"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewMapManager("")
	m.maps["0"] = &Environment{Name: "0"}
	m.maps["D001"] = &Environment{Name: "D001"}

	m.InitMiniMaps(tmpDir)

	if got := m.maps["0"].MinMap; got != 5 {
		t.Errorf("map 0 MinMap = %d, want 5", got)
	}
	if got := m.maps["D001"].MinMap; got != 1 {
		t.Errorf("map D001 MinMap = %d, want 1", got)
	}
	if _, ok := m.maps["missing"]; ok {
		t.Error("unknown map 'missing' should not create an environment")
	}
}

// TestInitMiniMapsMissingFile 缺少配置文件时不 panic、不写入。
func TestInitMiniMapsMissingFile(t *testing.T) {
	m := NewMapManager("")
	m.maps["0"] = &Environment{Name: "0"}
	m.InitMiniMaps(t.TempDir()) // 空目录，无 maps/mini_map.jsonc
	if m.maps["0"].MinMap != 0 {
		t.Errorf("MinMap = %d, want 0 when config missing", m.maps["0"].MinMap)
	}
}

// ============================================================================
// HandleWantMinimap — 真实 TCP 回环验证响应消息字段
// ============================================================================

// startMinimapServer 启动回环 TCP 服务端并返回首个连入的会话。
func startMinimapServer(t *testing.T) (*netserver.TCPServer, *netserver.Session, net.Conn) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := listener.Addr().String()
	listener.Close()

	sessCh := make(chan *netserver.Session, 1)
	srv := netserver.NewTCPServer(addr)
	srv.SetConnectHandler(func(s *netserver.Session) { sessCh <- s })
	if err := srv.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(srv.Stop)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	select {
	case session := <-sessCh:
		return srv, session, conn
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for server session")
		return nil, nil, nil
	}
}

// readFrame 从连接读取一帧 "#<payload>!" 并解码消息头。
func readFrame(t *testing.T, conn net.Conn) protocol.DefaultMessage {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("read frame: %v", err)
	}
	data := buf[:n]
	if len(data) < 3 || data[0] != '#' || data[len(data)-1] != '!' {
		t.Fatalf("invalid frame: %q", string(data))
	}
	payload := string(data[1 : len(data)-1])
	if len(payload) < protocol.DefBlockSize {
		t.Fatalf("payload too short: %d bytes", len(payload))
	}
	return protocol.DecodeMessage(payload[:protocol.DefBlockSize])
}

// TestHandleWantMinimapOK 验证有小地图时回 SM_READMINIMAP_OK，
// 且索引放在 Param 字段（Delphi SendDefMessage(SM_READMINIMAP_OK, 0, nMinMap, 0, 0, '')）。
func TestHandleWantMinimapOK(t *testing.T) {
	srv, session, conn := startMinimapServer(t)

	p := NewPlayObject(session, "tester", 1)
	p.envir = &Environment{Name: "0", MinMap: 7}

	p.HandleWantMinimap(srv)

	msg := readFrame(t, conn)
	if msg.Ident != protocol.SMReadMinimapOK {
		t.Fatalf("Ident = %d, want %d (SMReadMinimapOK)", msg.Ident, protocol.SMReadMinimapOK)
	}
	if msg.Param != 7 {
		t.Errorf("Param = %d, want 7 (minimap index rides in Param)", msg.Param)
	}
	if msg.Recog != 0 {
		t.Errorf("Recog = %d, want 0", msg.Recog)
	}
}

// TestHandleWantMinimapFail 验证 MinMap=0 时回 SM_READMINIMAP_FAIL。
func TestHandleWantMinimapFail(t *testing.T) {
	srv, session, conn := startMinimapServer(t)

	p := NewPlayObject(session, "tester", 1)
	p.envir = &Environment{Name: "wild", MinMap: 0}

	p.HandleWantMinimap(srv)

	msg := readFrame(t, conn)
	if msg.Ident != protocol.SMReadMinimapFail {
		t.Fatalf("Ident = %d, want %d (SMReadMinimapFail)", msg.Ident, protocol.SMReadMinimapFail)
	}
}

// TestHandleWantMinimapNilEnvir 验证 envir 为 nil 时不 panic 且回 Fail。
func TestHandleWantMinimapNilEnvir(t *testing.T) {
	srv, session, conn := startMinimapServer(t)

	p := NewPlayObject(session, "tester", 1)
	// p.envir 保持 nil

	p.HandleWantMinimap(srv)

	msg := readFrame(t, conn)
	if msg.Ident != protocol.SMReadMinimapFail {
		t.Fatalf("Ident = %d, want %d (SMReadMinimapFail)", msg.Ident, protocol.SMReadMinimapFail)
	}
}

// TestWantMinimapAllowedWhileDead 验证死亡守卫白名单：
// 死亡玩家发 CMWantMinimap 仍会被处理（纯展示消息）。
func TestWantMinimapAllowedWhileDead(t *testing.T) {
	srv, session, conn := startMinimapServer(t)

	p := NewPlayObject(session, "tester", 1)
	p.envir = &Environment{Name: "0", MinMap: 3}
	p.Death = true

	p.ProcessMessage(SendMessage{Ident: protocol.CMWantMinimap}, srv)

	msg := readFrame(t, conn)
	if msg.Ident != protocol.SMReadMinimapOK {
		t.Fatalf("Ident = %d, want %d (dead player's CMWantMinimap must not be dropped)",
			msg.Ident, protocol.SMReadMinimapOK)
	}
	if msg.Param != 3 {
		t.Errorf("Param = %d, want 3", msg.Param)
	}
}

// TestDeathGateAllowsRMRelay 回归：玩家死亡后 RM 中继必须正常下发。
// 受害者自己的 RM_DEATH 是客户端 SM_NOWDEATH（死亡动画+灰屏）的唯一
// 来源——此前死亡守卫把 RM 与 CM 一并丢弃，客户端永远不进入死亡流程。
func TestDeathGateAllowsRMRelay(t *testing.T) {
	srv, session, conn := startMinimapServer(t)

	env := newTestEnv(40, 40)
	p := NewPlayObject(session, "victim", 1)
	p.envir = env
	p.CurrX, p.CurrY = 20, 20
	env.AddObject(20, 20, OS_MOVINGOBJECT, p)
	p.Death = true

	// broadcastDeathMsg(justDied=true) 入队的消息形态
	p.ProcessMessage(SendMessage{
		Ident: RM_DEATH, SourceID: p.ID,
		Param1: 20, Param2: 20, Dir: 3, Param3: 1,
	}, srv)

	msg := readFrame(t, conn)
	if msg.Ident != protocol.SMNowDeath {
		t.Fatalf("Ident = %d, want %d (SMNowDeath)——死亡玩家自己的 RM_DEATH 不得被守卫丢弃",
			msg.Ident, protocol.SMNowDeath)
	}
	if msg.Recog != p.ID {
		t.Errorf("Recog = %d, want %d", msg.Recog, p.ID)
	}
}

// TestDeathGateBlocksCMAction 回归：死亡守卫仍须拦截 CM 动作消息。
func TestDeathGateBlocksCMAction(t *testing.T) {
	srv, session, conn := startMinimapServer(t)

	p := NewPlayObject(session, "victim", 1)
	p.envir = newTestEnv(40, 40)
	p.CurrX, p.CurrY = 20, 20
	p.Death = true

	p.ProcessMessage(SendMessage{Ident: protocol.CMWalk, Param1: 5, Param2: 20, Param3: 21}, srv)

	conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	buf := make([]byte, 256)
	if n, err := conn.Read(buf); err == nil && n > 0 {
		t.Fatalf("死亡后 CM 动作不应有任何响应，收到 %q", buf[:n])
	}
}
