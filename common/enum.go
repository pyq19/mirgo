package common

type MirGender uint8

const (
	MirGenderMale   MirGender = 0
	MirGenderFemale MirGender = 1
)

type MirClass uint8

const (
	MirClassWarrior MirClass = iota
	MirClassWizard
	MirClassTaoist
	MirClassAssassin
	MirClassArcher
)

type LightSetting uint8

const (
	LightSettingNormal  LightSetting = 0
	LightSettingDawn                 = 1
	LightSettingDay                  = 2
	LightSettingEvening              = 3
	LightSettingNight                = 4
)

type MirDirection uint8

const (
	MirDirectionUp        MirDirection = 0
	MirDirectionUpRight                = 1
	MirDirectionRight                  = 2
	MirDirectionDownRight              = 3
	MirDirectionDown                   = 4
	MirDirectionDownLeft               = 5
	MirDirectionLeft                   = 6
	MirDirectionUpLeft                 = 7
)

type ItemType uint8

const (
	ItemTypeNothing          ItemType = 0
	ItemTypeWeapon                    = 1
	ItemTypeArmour                    = 2
	ItemTypeHelmet                    = 4
	ItemTypeNecklace                  = 5
	ItemTypeBracelet                  = 6
	ItemTypeRing                      = 7
	ItemTypeAmulet                    = 8
	ItemTypeBelt                      = 9
	ItemTypeBoots                     = 10
	ItemTypeStone                     = 11
	ItemTypeTorch                     = 12
	ItemTypePotion                    = 13
	ItemTypeOre                       = 14
	ItemTypeMeat                      = 15
	ItemTypeCraftingMaterial          = 16
	ItemTypeScroll                    = 17
	ItemTypeGem                       = 18
	ItemTypeMount                     = 19
	ItemTypeBook                      = 20
	ItemTypeScript                    = 21
	ItemTypeReins                     = 22
	ItemTypeBells                     = 23
	ItemTypeSaddle                    = 24
	ItemTypeRibbon                    = 25
	ItemTypeMask                      = 26
	ItemTypeFood                      = 27
	ItemTypeHook                      = 28
	ItemTypeFloat                     = 29
	ItemTypeBait                      = 30
	ItemTypeFinder                    = 31
	ItemTypeReel                      = 32
	ItemTypeFish                      = 33
	ItemTypeQuest                     = 34
	ItemTypeAwakening                 = 35
	ItemTypePets                      = 36
	ItemTypeTransform                 = 37
)

type ItemGrade uint8

const (
	ItemGradeNone      ItemGrade = 0
	ItemGradeCommon              = 1
	ItemGradeRare                = 2
	ItemGradeLegendary           = 3
	ItemGradeMythical            = 4
)

type RequiredType uint8

const (
	RequiredTypeLevel    RequiredType = 0
	RequiredTypeMaxAC                 = 1
	RequiredTypeMaxMAC                = 2
	RequiredTypeMaxDC                 = 3
	RequiredTypeMaxMC                 = 4
	RequiredTypeMaxSC                 = 5
	RequiredTypeMaxLevel              = 6
	RequiredTypeMinAC                 = 7
	RequiredTypeMinMAC                = 8
	RequiredTypeMinDC                 = 9
	RequiredTypeMinMC                 = 10
	RequiredTypeMinSC                 = 11
)

type RequiredClass uint8

const (
	RequiredClassWarrior  RequiredClass = 1
	RequiredClassWizard                 = 2
	RequiredClassTaoist                 = 4
	RequiredClassAssassin               = 8
	RequiredClassArcher                 = 16
	//WarWizTao = Warrior | Wizard | Taoist,
	//None = WarWizTao | Assassin | Archer
)

type RequiredGender uint8

const (
	RequiredGenderMale   RequiredGender = 1
	RequiredGenderFemale                = 2
	//None = Male | Female
)

type ItemSet uint8

const (
	ItemSetNone       ItemSet = 0
	ItemSetSpirit             = 1
	ItemSetRecall             = 2
	ItemSetRedOrchid          = 3
	ItemSetRedFlower          = 4
	ItemSetSmash              = 5
	ItemSetHwanDevil          = 6
	ItemSetPurity             = 7
	ItemSetFiveString         = 8
	ItemSetMundane            = 9
	ItemSetNokChi             = 10
	ItemSetTaoProtect         = 11
	ItemSetMir                = 12
	ItemSetBone               = 13
	ItemSetBug                = 14
	ItemSetWhiteGold          = 15
	ItemSetWhiteGoldH         = 16
	ItemSetRedJade            = 17
	ItemSetRedJadeH           = 18
	ItemSetNephrite           = 19
	ItemSetNephriteH          = 20
	ItemSetWhisker1           = 21
	ItemSetWhisker2           = 22
	ItemSetWhisker3           = 23
	ItemSetWhisker4           = 24
	ItemSetWhisker5           = 25
	ItemSetHyeolryong         = 26
	ItemSetMonitor            = 27
	ItemSetOppressive         = 28
	ItemSetPaeok              = 29
	ItemSetSulgwan            = 30
)

type LevelEffects uint8

const (
	LevelEffectsNone       = 0
	LevelEffectsMist       = 0x0001
	LevelEffectsRedDragon  = 0x0002
	LevelEffectsBlueDragon = 0x0004
)

type ChatType uint8

const (
	ChatTypeNormal       ChatType = 0
	ChatTypeShout                 = 1
	ChatTypeSystem                = 2
	ChatTypeHint                  = 3
	ChatTypeAnnouncement          = 4
	ChatTypeGroup                 = 5
	ChatTypeWhisperIn             = 6
	ChatTypeWhisperOut            = 7
	ChatTypeGuild                 = 8
	ChatTypeTrainer               = 9
	ChatTypeLevelUp               = 10
	ChatTypeSystem2               = 11
	ChatTypeRelationship          = 12
	ChatTypeMentor                = 13
	ChatTypeShout2                = 14
	ChatTypeShout3                = 15
)

type MirGridType uint8

const (
	MirGridTypeNone           MirGridType = 0
	MirGridTypeInventory                  = 1
	MirGridTypeEquipment                  = 2
	MirGridTypeTrade                      = 3
	MirGridTypeStorage                    = 4
	MirGridTypeBuyBack                    = 5
	MirGridTypeDropPanel                  = 6
	MirGridTypeInspect                    = 7
	MirGridTypeTrustMerchant              = 8
	MirGridTypeGuildStorage               = 9
	MirGridTypeGuestTrade                 = 10
	MirGridTypeMount                      = 11
	MirGridTypeFishing                    = 12
	MirGridTypeQuestInventory             = 13
	MirGridTypeAwakenItem                 = 14
	MirGridTypeMail                       = 15
	MirGridTypeRefine                     = 16
	MirGridTypeRenting                    = 17
	MirGridTypeGuestRenting               = 18
	MirGridTypeCraft                      = 19
)

type AttackMode uint8

const (
	AttackModePeace      AttackMode = 0
	AttackModeGroup                 = 1
	AttackModeGuild                 = 2
	AttackModeEnemyGuild            = 3
	AttackModeRedBrown              = 4
	AttackModeAll                   = 5
)

type Spell uint8

const (
	SpellNone             Spell = 0
	SpellFencing                = 1 //Warrior
	SpellSlaying                = 2
	SpellThrusting              = 3
	SpellHalfMoon               = 4
	SpellShoulderDash           = 5
	SpellTwinDrakeBlade         = 6
	SpellEntrapment             = 7
	SpellFlamingSword           = 8
	SpellLionRoar               = 9
	SpellCrossHalfMoon          = 10
	SpellBladeAvalanche         = 11
	SpellProtectionField        = 12
	SpellRage                   = 13
	SpellCounterAttack          = 14
	SpellSlashingBurst          = 15
	SpellFury                   = 16
	SpellImmortalSkin           = 17
	SpellFireBall               = 31 //Wizard
	SpellRepulsion              = 32
	SpellElectricShock          = 33
	SpellGreatFireBall          = 34
	SpellHellFire               = 35
	SpellThunderBolt            = 36
	SpellTeleport               = 37
	SpellFireBang               = 38
	SpellFireWall               = 39
	SpellLightning              = 40
	SpellFrostCrunch            = 41
	SpellThunderStorm           = 42
	SpellMagicShield            = 43
	SpellTurnUndead             = 44
	SpellVampirism              = 45
	SpellIceStorm               = 46
	SpellFlameDisruptor         = 47
	SpellMirroring              = 48
	SpellFlameField             = 49
	SpellBlizzard               = 50
	SpellMagicBooster           = 51
	SpellMeteorStrike           = 52
	SpellIceThrust              = 53
	SpellFastMove               = 54
	SpellStormEscape            = 55
	SpellHealing                = 61 //Taoist
	SpellSpiritSword            = 62
	SpellPoisoning              = 63
	SpellSoulFireBall           = 64
	SpellSummonSkeleton         = 65
	SpellHiding                 = 67
	SpellMassHiding             = 68
	SpellSoulShield             = 69
	SpellRevelation             = 70
	SpellBlessedArmour          = 71
	SpellEnergyRepulsor         = 72
	SpellTrapHexagon            = 73
	SpellPurification           = 74
	SpellMassHealing            = 75
	SpellHallucination          = 76
	SpellUltimateEnhancer       = 77
	SpellSummonShinsu           = 78
	SpellReincarnation          = 79
	SpellSummonHolyDeva         = 80
	SpellCurse                  = 81
	SpellPlague                 = 82
	SpellPoisonCloud            = 83
	SpellEnergyShield           = 84
	SpellPetEnhancer            = 85
	SpellHealingCircle          = 86
	SpellFatalSword             = 91 //Assassin
	SpellDoubleSlash            = 92
	SpellHaste                  = 93
	SpellFlashDash              = 94
	SpellLightBody              = 95
	SpellHeavenlySword          = 96
	SpellFireBurst              = 97
	SpellTrap                   = 98
	SpellPoisonSword            = 99
	SpellMoonLight              = 100
	SpellMPEater                = 101
	SpellSwiftFeet              = 102
	SpellDarkBody               = 103
	SpellHemorrhage             = 104
	SpellCrescentSlash          = 105
	SpellMoonMist               = 106
	SpellFocus                  = 121 //Archer
	SpellStraightShot           = 122
	SpellDoubleShot             = 123
	SpellExplosiveTrap          = 124
	SpellDelayedExplosion       = 125
	SpellMeditation             = 126
	SpellBackStep               = 127
	SpellElementalShot          = 128
	SpellConcentration          = 129
	SpellStonetrap              = 130
	SpellElementalBarrier       = 131
	SpellSummonVampire          = 132
	SpellVampireShot            = 133
	SpellSummonToad             = 134
	SpellPoisonShot             = 135
	SpellCrippleShot            = 136
	SpellSummonSnakes           = 137
	SpellNapalmShot             = 138
	SpellOneWithNature          = 139
	SpellBindingShot            = 140
	SpellMentalState            = 141
	SpellBlink                  = 151 //Custom
	SpellPortal                 = 152
	SpellBattleCry              = 153
	SpellDigOutZombie           = 200 //Map Events
	SpellRubble                 = 201
	SpellMapLightning           = 202
	SpellMapLava                = 203
	SpellMapQuake1              = 204
	SpellMapQuake2              = 205
)

type PanelType uint8

const (
	PanelTypeBuy PanelType = 0
	PanelTypeSell
	PanelTypeRepair
	PanelTypeSpecialRepair
	PanelTypeConsign
	PanelTypeCraft
	PanelTypeRefine
	PanelTypeCheckRefine
	PanelTypeDisassemble
	PanelTypeDowngrade
	PanelTypeReset
	PanelTypeCollectRefine
	PanelTypeReplaceWedRing
)

type CellAttribute uint8

const (
	CellAttributeWalk     CellAttribute = 0
	CellAttributeHighWall               = 1
	CellAttributeLowWall                = 2
)
