package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
)

type Environ struct {
	Game          *Game
	Basic         common.Basic
	Players       *[]Player
	Maps          *[]Map
	Monsters      *[]Monster
	NPCs          *[]NPC
	GameShopItems *[]common.GameShopItem
	ItemInfos     *[]common.ItemInfo
	MagicInfos    *[]common.MagicInfo
	MovementInfos *[]common.MovementInfo
	QuestInfos    *[]common.QuestInfo
}

type Map struct {
	Info      *common.MapInfo
	SafeZones *[]common.SafeZoneInfo
}

type MapObject interface {
}

type Player struct {
	MapObject
	Session   *cellnet.Session
	Character *common.Character
	UserItems *[]common.UserItem
}

type Monster struct {
	MapObject
	Info *common.MonsterInfo
}

type NPC struct {
	MapObject
	Info *common.NpcInfo
}

// TODO
func (g *Game) NewEnviron() (env *Environ) {
	env = new(Environ)
	env.Game = g
	return
}
