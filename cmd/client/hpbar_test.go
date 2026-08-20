package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/protocol"
)

// TestReadyActionOpenHealthStateBit 验证移动/出现消息按 charStatus 的
// STATE_OPENHEATH 位（$02）重置头顶血条开关（Delphi ReadyAction，
// Actor.pas:1452-1463）。
func TestReadyActionOpenHealthStateBit(t *testing.T) {
	moveIdents := []int{
		protocol.SMWalk, protocol.SMRun, protocol.SMTurn, protocol.SMRush,
		protocol.SMRushKung, protocol.SMBackStep, protocol.SMHorseRun,
		protocol.SMDigUp, protocol.SMAlive,
	}
	for _, ident := range moveIdents {
		// 带 STATE_OPENHEATH 位 → 开启
		a := &Actor{}
		a.ReadyAction(ActorMsg{Ident: ident, State: 0x02})
		if !a.ShowHP {
			t.Errorf("Ident=%d State=0x02: ShowHP = false, want true", ident)
		}
		// 无该位 → 关闭
		a.ReadyAction(ActorMsg{Ident: ident, State: 0})
		if a.ShowHP {
			t.Errorf("Ident=%d State=0: ShowHP = true, want false", ident)
		}
	}

	// 死亡中的角色不重置（Delphi not m_boDeath 守卫，Actor.pas:1452）
	d := &Actor{Death: true, ShowHP: true}
	d.ReadyAction(ActorMsg{Ident: protocol.SMWalk, State: 0})
	if !d.ShowHP {
		t.Error("死亡角色不应被移动消息重置血条开关")
	}
}

// TestParseCharDesc 验证 TCharDesc{feature, charStatus} 解析。
func TestParseCharDesc(t *testing.T) {
	// 8 字节小端 {feature=0x11223344, status=0x00000002}
	body := string([]byte{0x44, 0x33, 0x22, 0x11, 0x02, 0x00, 0x00, 0x00})
	feature, status, ok := parseCharDesc(body)
	if !ok {
		t.Fatal("8 字节 body 应解析成功")
	}
	if feature != 0x11223344 {
		t.Errorf("feature = %x, want 0x11223344", feature)
	}
	if status != 0x02 {
		t.Errorf("status = %x, want 0x02", status)
	}

	// 短 body 失败
	if _, _, ok := parseCharDesc("abc"); ok {
		t.Error("短 body 不应解析成功")
	}
}
