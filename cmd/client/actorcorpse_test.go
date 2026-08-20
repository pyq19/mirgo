package main

import (
	"encoding/binary"
	"testing"

	"github.com/pyq0109/mirgo/internal/engine"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// 批1客户端回归：后进入视野的玩家收到带 TCharDesc body 的 SM_DEATH/
// SM_SKELETON 时创建尸体/骨架 actor（Delphi PlayScn.pas:2321-2338）。

func charDescBody(raceImg int, appr int) string {
	feature := uint32(raceImg&0xFF) | uint32(appr&0xFFFF)<<16
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[0:4], feature)
	return string(buf)
}

func TestSMDeathCreatesCorpseActor(t *testing.T) {
	h := &NetHandler{playScene: NewPlayScene(nil, &engine.ResourceManager{}, "")}

	// 稻草人 feature：raceImg=18 低字节、appr=27 高字（types.go:378-381）
	msg := protocol.DefaultMessage{Recog: 555, Ident: protocol.SMDeath, Param: 10, Tag: 12, Series: 3}
	h.HandleMessage(msg, charDescBody(18, 27), charDescBody(18, 27))

	actor := h.playScene.State.Actors.Get(555)
	if actor == nil {
		t.Fatal("带 body 的 SMDeath 应对未知 actor 创建尸体")
	}
	if !actor.Death {
		t.Fatal("尸体 actor Death 应为 true")
	}
	if actor.Race != 18 || actor.Appearance != 27 {
		t.Errorf("actor Race/Appearance = %d/%d, want 18/27", actor.Race, actor.Appearance)
	}

	// 空 body（旧消息兼容）不创建
	msg2 := protocol.DefaultMessage{Recog: 556, Ident: protocol.SMDeath, Param: 10, Tag: 12}
	h.HandleMessage(msg2, "", "")
	if h.playScene.State.Actors.Get(556) != nil {
		t.Fatal("空 body 的 SMDeath 不应创建 actor")
	}
}

func TestSMSkeletonCreatesCorpseActor(t *testing.T) {
	h := &NetHandler{playScene: NewPlayScene(nil, &engine.ResourceManager{}, "")}

	msg := protocol.DefaultMessage{Recog: 777, Ident: protocol.SMSkeleton, Param: 5, Tag: 6, Series: 2}
	h.HandleMessage(msg, charDescBody(11, 160), charDescBody(11, 160))

	actor := h.playScene.State.Actors.Get(777)
	if actor == nil {
		t.Fatal("带 body 的 SMSkeleton 应对未知 actor 创建骨架")
	}
	if !actor.Death || !actor.Skeleton {
		t.Fatal("骨架 actor 应 Death=true 且 Skeleton=true（PlayScn.pas:2333-2336）")
	}
}

// TestSkeletonCorpseIdleFrame 骨架尸体的空闲帧 = ActDeath.start
//（Delphi TSkeletonOma.GetDefaultFrame，AxeMon.pas:420-422），
// 普通尸体 = ActDie 末帧（Actor.pas:2104-2107）。
func TestSkeletonCorpseIdleFrame(t *testing.T) {
	a := NewActor(1, 0, 0, 0)
	a.Type = ActorMonster
	a.MonAction = &MA11
	a.Death = true

	a.Skeleton = true
	if got, want := a.getEndFrame(), MA11.ActDeath.Start; got != want {
		t.Errorf("骨架空闲帧 = %d, want %d（ActDeath.start）", got, want)
	}

	a.Skeleton = false
	_, end := CalcFrame(MA11.ActDie, a.Dir)
	if got := a.getEndFrame(); got != end {
		t.Errorf("普通尸体空闲帧 = %d, want %d（ActDie 末帧）", got, end)
	}
}
