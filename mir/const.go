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

const (
	MainKey           = "[@MAIN]"
	BuyKey            = "[@BUY]"
	SellKey           = "[@SELL]"
	BuySellKey        = "[@BUYSELL]"
	RepairKey         = "[@REPAIR]"
	SRepairKey        = "[@SREPAIR]"
	RefineKey         = "[@REFINE]"
	RefineCheckKey    = "[@REFINECHECK]"
	RefineCollectKey  = "[@REFINECOLLECT]"
	ReplaceWedRingKey = "[@REPLACEWEDDINGRING]"
	BuyBackKey        = "[@BUYBACK]"
	StorageKey        = "[@STORAGE]"
	ConsignKey        = "[@CONSIGN]"
	MarketKey         = "[@MARKET]"
	ConsignmentsKey   = "[@CONSIGNMENT]"
	CraftKey          = "[@CRAFT]"
	TradeKey          = "[TRADE]"
	RecipeKey         = "[RECIPE]"
	TypeKey           = "[TYPES]"
	QuestKey          = "[QUESTS]"
	GuildCreateKey    = "[@CREATEGUILD]"
	RequestWarKey     = "[@REQUESTWAR]"
	SendParcelKey     = "[@SENDPARCEL]"
	CollectParcelKey  = "[@COLLECTPARCEL]"
	AwakeningKey      = "[@AWAKENING]"
	DisassembleKey    = "[@DISASSEMBLE]"
	DowngradeKey      = "[@DOWNGRADE]"
	ResetKey          = "[@RESET]"
	PearlBuyKey       = "[@PEARLBUY]"
	BuyUsedKey        = "[@BUYUSED]"
)

var settings = setting.DefaultSettings()
