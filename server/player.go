package main

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
)

const (
	LOGIN = iota
	SELECT
	GAME
	DISCONNECTED
)

// Player ...
type Player struct {
	Cell      *Cell
	ID        string
	AccountId int
	GameStage int
	Session   *cellnet.Session
	Character *common.Character
	Magics    *[]common.MagicInfo
	UserItems *[]common.UserItem
}

func (p *Player) Point() *common.Point {
	x := int(p.Character.CurrentLocationX)
	y := int(p.Character.CurrentLocationY)
	return common.NewPoint(x, y)
}
