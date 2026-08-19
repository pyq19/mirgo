package main

import (
	"testing"

	"github.com/pyq0109/mirgo/internal/protocol"
)

func TestTakeOnSlotMatches(t *testing.T) {
	cases := []struct {
		stdMode byte
		slot    int
		want    bool
	}{
		{5, protocol.UWeapon, true},
		{6, protocol.UWeapon, true},
		{5, protocol.UDress, false},
		{10, protocol.UDress, true},
		{15, protocol.UHelmet, true},
		{19, protocol.UNecklace, true},
		{30, protocol.URightHand, true},
		{22, protocol.URingL, true},
		{22, protocol.URingR, true},
		{23, protocol.URingL, true},
		{22, protocol.UNecklace, false},
		{24, protocol.UArmRingL, true},
		{24, protocol.UArmRingR, true},
		{26, protocol.UArmRingR, true},
		// 护身符额外接受左手镯 (FState:3326-3335)。
		{25, protocol.UBujuk, true},
		{25, protocol.UArmRingL, true},
		{25, protocol.UArmRingR, false},
		{51, protocol.UBujuk, true},
		{51, protocol.UArmRingL, true},
		{52, protocol.UBoots, true},
		{53, protocol.UCharm, true},
		{54, protocol.UBelt, true},
	}
	for _, tt := range cases {
		if got := takeOnSlotMatches(tt.stdMode, tt.slot); got != tt.want {
			t.Errorf("takeOnSlotMatches(%d, %d) = %v, want %v", tt.stdMode, tt.slot, got, tt.want)
		}
	}
}

func TestDressSexOK(t *testing.T) {
	cases := []struct {
		stdMode byte
		sex     int
		want    bool
	}{
		{10, 0, true},
		{10, 1, false},
		{11, 1, true},
		{11, 0, false},
		{5, 1, true}, // 非衣服不校验
	}
	for _, tt := range cases {
		if got := dressSexOK(tt.stdMode, tt.sex); got != tt.want {
			t.Errorf("dressSexOK(%d, %d) = %v, want %v", tt.stdMode, tt.sex, got, tt.want)
		}
	}
}
