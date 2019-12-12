package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
)

type Environ struct {
	Game    *Game
	Players *[]Player // 总玩家
}

type ServerDB struct {
	Basic         common.Basic
	GameShopItems []common.GameShopItem
	ItemInfos     []common.ItemInfo
	MagicInfos    []common.MagicInfo
	MapInfos      []common.MapInfo
	MonsterInfos  []common.MonsterInfo
	MovementInfos []common.MovementInfo
	NpcInfos      []common.NpcInfo
	QuestInfos    []common.QuestInfo
	RespawnInfos  []common.RespawnInfo
	SafeZoneInfos []common.SafeZoneInfo
}

type Player struct {
	Session   *cellnet.Session
	Character *common.Character
	Magics    *[]common.MagicInfo
	UserItems *[]common.UserItem
}

// TODO
func (g *Game) NewEnviron() (env *Environ) {
	env = new(Environ)
	env.Game = g
	return
}
