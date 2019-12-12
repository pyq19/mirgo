package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
)

type Environ struct {
	Game    *Game
	GameDB  *GameDB
	Players *[]Player // 总玩家
}

type GameDB struct {
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

// TODO
func (d GameDB) Init() {

}

type Player struct {
	Session   *cellnet.Session
	Character *common.Character
	Magics    *[]common.MagicInfo
	UserItems *[]common.UserItem
}

func (g *Game) NewEnviron() (env *Environ) {
	env = new(Environ)
	env.Game = g

	gameDB := new(GameDB)
	gameDB.Init()

	players := make([]Player, 0)
	env.Players = &players
	return
}
