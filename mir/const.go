package mir

import "github.com/yenkeia/mirgo/setting"

const DataRange = 20

const (
	StopGameServerClosed              = 0
	StopGameDoubleLogin               = 1
	StopGameChatMessageTooLong        = 2
	StopGameServerCrashed             = 3
	StopGameKickedByAdmin             = 4
	StopGameMaximumConnectionsReached = 5
	StopGameWrongClientVersion        = 10
	StopGameDisconnected              = 20
	StopGameConnectionTimedOut        = 21
	StopGameUserClosedGame            = 22
	StopGameUserReturnedToSelectChar  = 23
	StopGameUnknown                   = 24
)

type DefaultNPCType byte

const (
	DefaultNPCTypeLogin DefaultNPCType = iota
	DefaultNPCTypeLevelUp
	DefaultNPCTypeUseItem
	DefaultNPCTypeMapCoord
	DefaultNPCTypeMapEnter
	DefaultNPCTypeDie
	DefaultNPCTypeTrigger
	DefaultNPCTypeCustomCommand
	DefaultNPCTypeOnAcceptQuest
	DefaultNPCTypeOnFinishQuest
	DefaultNPCTypeDaily
	DefaultNPCTypeTalkMonster
)

var settings = setting.DefaultSettings()
