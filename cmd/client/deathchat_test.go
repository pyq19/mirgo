package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/engine"
	"github.com/pyq0109/mirgo/internal/protocol"
)

// Go 专属扩展回归：玩家自己死亡（SM_NOWDEATH）时聊天框出红字提示；
// 骨架化转换的 SM_DEATH 与其他玩家的死亡均不提示。

func newDeathChatHandler() *NetHandler {
	h := &NetHandler{playScene: NewPlayScene(nil, &engine.ResourceManager{}, "")}
	me := NewActor(100, 330, 330, 3)
	me.Type = ActorHuman
	h.playScene.State.Actors.Add(me)
	h.playScene.State.MySelf = me
	return h
}

func TestMyDeathAddsChatHint(t *testing.T) {
	h := newDeathChatHandler()

	msg := protocol.DefaultMessage{Recog: 100, Ident: protocol.SMNowDeath, Param: 330, Tag: 330}
	h.HandleMessage(msg, "", "")

	if !h.playScene.deathGray {
		t.Fatal("自己死亡应开启死亡灰屏")
	}
	n := len(h.playScene.chatMessages)
	if n == 0 {
		t.Fatal("玩家自己死亡应在聊天框追加死亡提示")
	}
	if got := h.playScene.chatMessages[n-1].Text; got != "你死了." {
		t.Fatalf("死亡提示文案 = %q, want %q", got, "你死了.")
	}
}

func TestMySkeletonTransitionNoDuplicateHint(t *testing.T) {
	h := newDeathChatHandler()
	h.playScene.State.MySelf.Death = true // 已死亡，收到骨架化转换

	msg := protocol.DefaultMessage{Recog: 100, Ident: protocol.SMDeath, Param: 330, Tag: 330}
	h.HandleMessage(msg, "", "")

	if n := len(h.playScene.chatMessages); n != 0 {
		t.Fatalf("骨架化转换的 SMDeath 不应重复死亡提示, got %d 条", n)
	}
}

func TestOtherActorDeathNoChatHint(t *testing.T) {
	h := newDeathChatHandler()
	other := NewActor(200, 331, 330, 5)
	other.Type = ActorHuman
	h.playScene.State.Actors.Add(other)

	msg := protocol.DefaultMessage{Recog: 200, Ident: protocol.SMNowDeath, Param: 331, Tag: 330}
	h.HandleMessage(msg, "", "")

	if n := len(h.playScene.chatMessages); n != 0 {
		t.Fatalf("其他玩家死亡不应在聊天框提示, got %d 条", n)
	}
}
