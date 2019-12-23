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

//func (p *Player) SetCell(c *Cell) {
//	p.Cell = c
//	p.Character.CurrentMapId = int32(c.Map.Id)
//	p.Character.CurrentLocationX = int32(c.Point.X)
//	p.Character.CurrentLocationY = int32(c.Point.Y)
//}
