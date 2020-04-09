package game

import (
	"github.com/yenkeia/mirgo/setting"
)

var settings *setting.Settings

const DataRange = 20

const MaxGroup = 5 // 小队的最大人数

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

const (
	AIDeer1                     = 1
	AIDeer2                     = 2
	AITree                      = 3
	AISpittingSpider            = 4
	AICannibalPlant             = 5
	AIGuard1                    = 6
	AICaveMaggot                = 7
	AIAxeSkeleton               = 8
	AIHarvestMonster            = 9
	AIFlamingWooma              = 10
	AIWoomaTaurus               = 11
	AIBugBagMaggot              = 12
	AIRedMoonEvil               = 13
	AIEvilCentipede             = 14
	AIZumaMonster               = 15
	AIRedThunderZuma            = 16
	AIZumaTaurus                = 17
	AIShinsu                    = 18
	AIKingScorpion              = 19
	AIDarkDevil                 = 20
	AIIncarnatedGhoul           = 21
	AIIncarnatedZT              = 22
	AIBoneFamiliar              = 23
	AIDigOutZombie              = 24
	AIRevivingZombie            = 25
	AIShamanZombie              = 26
	AIKhazard                   = 27
	AIToxicGhoul                = 28
	AIBoneSpearman              = 29
	AIBoneLord                  = 30
	AIRightGuard                = 31
	AILeftGuard                 = 32
	AIMinotaurKing              = 33
	AIFrostTiger                = 34
	AISandWorm                  = 35
	AIYimoogi                   = 36
	AICrystalSpider             = 37
	AIHolyDeva                  = 38
	AIRootSpider                = 39
	AIBombSpider                = 40
	AIYinDevilNode              = 41
	AIOmaKing                   = 43
	AIBlackFoxman               = 44
	AIRedFoxman                 = 45
	AIWhiteFoxman               = 46
	AITrapRock                  = 47
	AIGuardianRock              = 48
	AIThunderElement            = 49
	AIGreatFoxSpirit            = 50
	AIHedgeKekTal               = 51
	AIEvilMir                   = 52
	AIEvilMirBody               = 53
	AIDragonStatue              = 54
	AIHumanWizard               = 55
	AITrainer                   = 56
	AITownArcher                = 57
	AIGuard2                    = 58
	AIHumanAssassin             = 59
	AIVampireSpider             = 60
	AISpittingToad              = 61
	AISnakeTotem                = 62
	AICharmedSnake              = 63
	AIIntelligentCreatureObject = 64
	AIMutatedManworm            = 65
	AICrazyManworm              = 66
	AIDarkDevourer              = 67
	AIFootball                  = 68
	AIPoisonHugger              = 69
	AIHugger                    = 70
	AIBehemoth                  = 71
	AIFinialTurtle              = 72
	AITurtleKing                = 73
	AILightTurtle               = 74
	AIWitchDoctor               = 75
	AIHellSlasher               = 76
	AIHellPirate                = 77
	AIHellCannibal              = 78
	AIHellKeeper                = 79
	AIConquestArcher            = 80
	AIGate                      = 81
	AIWall                      = 82
	AITornado                   = 83
	AIWingedTigerLord           = 84
	AIManectricClaw             = 86
	AIManectricBlest            = 87
	AIManectricKing             = 88
	AIIcePillar                 = 89
	AITrollBomber               = 90
	AITrollKing                 = 91
	AIFlameSpear                = 92
	AIFlameMage                 = 93
	AIFlameScythe               = 94
	AIFlameAssassin             = 95
	AIFlameQueen                = 96
	AIHellKnight                = 97
	AIHellLord                  = 98
	AIHellBomb                  = 99
	AIVenomSpider               = 100
)

const (
	EquipmentSlotWeapon    = 0
	EquipmentSlotArmour    = 1
	EquipmentSlotHelmet    = 2
	EquipmentSlotTorch     = 3
	EquipmentSlotNecklace  = 4
	EquipmentSlotBraceletL = 5
	EquipmentSlotBraceletR = 6
	EquipmentSlotRingL     = 7
	EquipmentSlotRingR     = 8
	EquipmentSlotAmulet    = 9
	EquipmentSlotBelt      = 10
	EquipmentSlotBoots     = 11
	EquipmentSlotStone     = 12
	EquipmentSlotMount     = 13
)
