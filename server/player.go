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
	AccountId int
	GameStage int
	Session   *cellnet.Session
	Character *common.Character
	Magics    *[]common.MagicInfo
	UserItems *[]common.UserItem
	Cell      *Cell
}
