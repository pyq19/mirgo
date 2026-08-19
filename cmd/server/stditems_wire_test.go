package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/netserver"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// 本文件守护的失败模式：SMStdItems 物品定义库超过 protocol.BufferSize
// 时 EncodeBuffer 曾静默返回空串，客户端收不到任何定义，背包图标回退
// 到数据库索引（金创药 idx=1 被画成 Items.wil[1] 的肉图标）。
// 测试在真实 TCP 链路上运行 SendStdItems，按真实客户端 ReadLoop 的方式
// （FrameScanner + '*' 回显）读取，确保整库定义完整送达。

// sendStdItemsOverWire 启动最小 TCPServer，用 SendStdItems 下发给定物品库，
// 返回客户端侧裸连接供读取。
func sendStdItemsOverWire(t *testing.T, items []ItemDef) net.Conn {
	t.Helper()
	log.SetLevel(log.LevelError)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := listener.Addr().String()
	listener.Close()

	server := netserver.NewTCPServer(addr)
	sessCh := make(chan *netserver.Session, 1)
	server.SetConnectHandler(func(s *netserver.Session) {
		select {
		case sessCh <- s:
		default:
		}
	})
	if err := server.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(server.Stop)

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	var session *netserver.Session
	select {
	case session = <-sessCh:
	case <-time.After(3 * time.Second):
		t.Fatal("session not established")
	}

	player := NewPlayObject(session, "StdItemsTester", 1)
	player.ItemDB = &ItemDB{Items: items}
	player.SendStdItems(server)
	return conn
}

// readStdItemsFrames 按真实客户端 ReadLoop 的方式（FrameScanner 切帧 +
// '*' 探针回显，见 cmd/client/main.go ReadLoop）读取 SMStdItems 消息，
// 返回 idx->looks 映射与消息头序列。
func readStdItemsFrames(t *testing.T, conn net.Conn, want int) (map[int]uint16, []protocol.DefaultMessage) {
	t.Helper()
	var scanner protocol.FrameScanner
	defs := make(map[int]uint16)
	var headers []protocol.DefaultMessage
	received := 0
	buf := make([]byte, 4096)
	for received < want {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := conn.Read(buf)
		if n == 0 && err != nil {
			t.Fatalf("read (received %d/%d defs): %v", received, want, err)
		}
		payloads, overflow := scanner.Feed(buf[:n], false, func() {
			conn.Write([]byte{'*'})
		})
		if overflow {
			t.Fatal("receive buffer overflow")
		}
		for _, payload := range payloads {
			if len(payload) < protocol.DefBlockSize {
				continue
			}
			msg := protocol.DecodeMessage(payload[:protocol.DefBlockSize])
			if msg.Ident != protocol.SMStdItems {
				continue
			}
			headers = append(headers, msg)
			body := payload[protocol.DefBlockSize:]
			// 空 body 正是原图标 bug 的失败形态，必须显式失败。
			if body == "" {
				t.Fatal("SMStdItems with empty body — EncodeBuffer silently dropped the payload")
			}
			raw := []byte(protocol.DecodeString(body))
			if len(raw) < 2 {
				t.Fatalf("SMStdItems body too short: %d bytes", len(raw))
			}
			count := int(binary.LittleEndian.Uint16(raw[0:2]))
			off := 2
			for i := 0; i < count; i++ {
				if off+38 > len(raw) {
					t.Fatalf("chunk truncated at record %d (off=%d, len=%d)", i, off, len(raw))
				}
				idx := int(binary.LittleEndian.Uint16(raw[off : off+2]))
				looks := binary.LittleEndian.Uint16(raw[off+2 : off+4])
				nameLen := int(raw[off+37])
				if off+38+nameLen > len(raw) {
					t.Fatalf("chunk truncated in name of record %d", i)
				}
				defs[idx] = looks
				off += 38 + nameLen
			}
			received += count
		}
	}
	if len(defs) != want {
		t.Fatalf("received %d records but only %d unique idx — duplicate idx in item DB", received, len(defs))
	}
	return defs, headers
}

func TestSendStdItemsWireRoundTrip(t *testing.T) {
	// 合成大库：规模超过 BufferSize 多个身位，足以触发分块与背压探针。
	items := []ItemDef{
		{Idx: 1, Name: "金创药(小量)", Looks: 398}, // 用户报告案例
		{Idx: 2, Name: "肉", Looks: 1},
	}
	for i := 3; i < 700; i++ {
		items = append(items, ItemDef{
			Idx:   i,
			Name:  fmt.Sprintf("测试物品%04d", i),
			Looks: uint16(500 + i),
		})
	}

	conn := sendStdItemsOverWire(t, items)
	defs, headers := readStdItemsFrames(t, conn, len(items))

	// 分块元数据：Recog=总数，Param=块序号（0 基递增），Tag=总块数。
	if len(headers) < 2 {
		t.Fatalf("chunks = %d, want >= 2 for %d items", len(headers), len(items))
	}
	for i, h := range headers {
		if h.Recog != int32(len(items)) {
			t.Errorf("chunk %d Recog = %d, want %d", i, h.Recog, len(items))
		}
		if h.Param != uint16(i) {
			t.Errorf("chunk %d Param = %d, want %d", i, h.Param, i)
		}
		if h.Tag != uint16(len(headers)) {
			t.Errorf("chunk %d Tag = %d, want %d", i, h.Tag, len(headers))
		}
	}

	// 全部定义完整到达且 Looks 正确。
	if defs[1] != 398 {
		t.Errorf("金创药(小量) Looks = %d, want 398（bug 症状：回退 idx=1 画成肉图标）", defs[1])
	}
	if defs[2] != 1 {
		t.Errorf("肉 Looks = %d, want 1", defs[2])
	}
	for _, it := range items {
		if defs[it.Idx] != it.Looks {
			t.Fatalf("item idx=%d Looks = %d, want %d", it.Idx, defs[it.Idx], it.Looks)
		}
	}
}

// TestSendStdItemsWireRoundTrip_RealDB 用真实转换产出的物品库跑同一条链路。
// serverconfig/ 为 gitignore 的本地资源，缺失时跳过。
func TestSendStdItemsWireRoundTrip_RealDB(t *testing.T) {
	path := filepath.Join("..", "..", "serverconfig", "items", "std_items.jsonc")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("real item DB not available: %v", err)
	}
	itemDB, err := LoadItemDB(path)
	if err != nil {
		t.Fatalf("LoadItemDB: %v", err)
	}
	if len(itemDB.Items) == 0 {
		t.Skip("empty item DB")
	}

	conn := sendStdItemsOverWire(t, itemDB.Items)
	defs, _ := readStdItemsFrames(t, conn, len(itemDB.Items))

	// 用户报告案例在真实库中的回归锚点。
	if potion := itemDB.GetByName("金创药(小量)"); potion != nil {
		if defs[potion.Idx] != potion.Looks {
			t.Errorf("金创药(小量) Looks = %d, want %d", defs[potion.Idx], potion.Looks)
		}
	}
	for i := range itemDB.Items {
		it := &itemDB.Items[i]
		if defs[it.Idx] != it.Looks {
			t.Fatalf("item idx=%d (%s) Looks = %d, want %d",
				it.Idx, it.Name, defs[it.Idx], it.Looks)
		}
	}
}
