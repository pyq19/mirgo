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
