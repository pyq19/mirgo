package cm

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
	MirDirectionCount                  = 8 // 方向个数
)

// NegativeDirection 反方向
func (m MirDirection) NegativeDirection() MirDirection {
	if m > 3 {
		return m - 4
	}
	return m + 4
}

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
	LevelEffectsNone       LevelEffects = 0
	LevelEffectsMist                    = 0x0001
	LevelEffectsRedDragon               = 0x0002
	LevelEffectsBlueDragon              = 0x0004
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

// Spell 技能翻译参考自：https://github.com/cjlaaa/mir2/blob/4132c5b1f198e1353152839d799eefc667003d35/Common.cs#L3753
type Spell uint8

const (
	SpellNone             Spell = 0
	SpellFencing                = 1  // 基本剑术 Warrior 战士
	SpellSlaying                = 2  // 攻杀剑术
	SpellThrusting              = 3  // 刺杀剑术
	SpellHalfMoon               = 4  // 半月弯刀
	SpellShoulderDash           = 5  // 野蛮冲撞
	SpellTwinDrakeBlade         = 6  // 双龙斩
	SpellEntrapment             = 7  // 捕绳剑
	SpellFlamingSword           = 8  // 烈火剑法
	SpellLionRoar               = 9  // 狮子吼
	SpellCrossHalfMoon          = 10 // 圆月弯刀
	SpellBladeAvalanche         = 11 // 攻破斩
	SpellProtectionField        = 12 // 护身气幕
	SpellRage                   = 13 // 剑气爆
	SpellCounterAttack          = 14 // 反击
	SpellSlashingBurst          = 15 // 日闪
	SpellFury                   = 16 // 血龙剑法
	SpellImmortalSkin           = 17
	SpellFireBall               = 31 // 火球术 Wizard 法师
	SpellRepulsion              = 32 // 抗拒火环
	SpellElectricShock          = 33 // 诱惑之光
	SpellGreatFireBall          = 34 // 大火球
	SpellHellFire               = 35 // 地狱火
	SpellThunderBolt            = 36 // 雷电术
	SpellTeleport               = 37 // 瞬息移动
	SpellFireBang               = 38 // 爆裂火焰
	SpellFireWall               = 39 // 火墙
	SpellLightning              = 40 // 疾光电影
	SpellFrostCrunch            = 41 // 寒冰掌
	SpellThunderStorm           = 42 // 地狱雷光
	SpellMagicShield            = 43 // 魔法盾
	SpellTurnUndead             = 44 // 圣言术
	SpellVampirism              = 45 // 嗜血术
	SpellIceStorm               = 46 // 冰咆哮
	SpellFlameDisruptor         = 47 // 火龙术
	SpellMirroring              = 48 // 分身术
	SpellFlameField             = 49 // 火龙气焰
	SpellBlizzard               = 50 // 天霜冰环
	SpellMagicBooster           = 51 // 深延术
	SpellMeteorStrike           = 52 // 流星火雨
	SpellIceThrust              = 53 // 冰焰术
	SpellFastMove               = 54
	SpellStormEscape            = 55
	SpellHealing                = 61 // 治愈术 Taoist 道士
	SpellSpiritSword            = 62 // 精神力战法
	SpellPoisoning              = 63 // 施毒术
	SpellSoulFireBall           = 64 // 灵魂火符
	SpellSummonSkeleton         = 65 // 召唤骷髅
	SpellHiding                 = 67 // 隐身术
	SpellMassHiding             = 68 // 集体隐身术
	SpellSoulShield             = 69 // 幽灵盾
	SpellRevelation             = 70 // 心灵启示
	SpellBlessedArmour          = 71 // 神圣战甲术
	SpellEnergyRepulsor         = 72 // 气功波
	SpellTrapHexagon            = 73 // 困魔咒
	SpellPurification           = 74 // 净化术
	SpellMassHealing            = 75 // 群体治疗术
	SpellHallucination          = 76 // 迷魂术
	SpellUltimateEnhancer       = 77 // 无极真气
	SpellSummonShinsu           = 78 // 召唤神兽
	SpellReincarnation          = 79 // 复活术
	SpellSummonHolyDeva         = 80 // 召唤月灵
	SpellCurse                  = 81 // 诅咒术
	SpellPlague                 = 82 // 瘟疫
	SpellPoisonCloud            = 83 // 毒云
	SpellEnergyShield           = 84 // 阴阳盾
	SpellPetEnhancer            = 85 // 血龙水
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
	PanelTypeBuy PanelType = iota
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

type PetMode uint8

const (
	PetModeBoth       PetMode = 0
	PetModeMoveOnly           = 1
	PetModeAttackOnly         = 2
	PetModeNone               = 3
)

type UserItemType uint8

const (
	UserItemTypeInventory      UserItemType = 0
	UserItemTypeEquipment                   = 1
	UserItemTypeQuestInventory              = 2
	UserItemTypeStorage                     = 3
	UserItemTypeTrade                       = 4
)

type PoisonType uint16

const (
	PoisonTypeNone             PoisonType = 0
	PoisonTypeGreen                       = 1
	PoisonTypeRed                         = 2
	PoisonTypeSlow                        = 4
	PoisonTypeFrozen                      = 8
	PoisonTypeStun                        = 16
	PoisonTypeParalysis                   = 32
	PoisonTypeDelayedExplosion            = 64
	PoisonTypeBleeding                    = 128
	PoisonTypeLRParalysis                 = 256
)

type SpellEffect uint8

const (
	SpellEffectNone SpellEffect = iota
	SpellEffectFatalSword
	SpellEffectTeleport
	SpellEffectHealing
	SpellEffectRedMoonEvil
	SpellEffectTwinDrakeBlade
	SpellEffectMagicShieldUp
	SpellEffectMagicShieldDown
	SpellEffectGreatFoxSpirit
	SpellEffectEntrapment
	SpellEffectReflect
	SpellEffectCritical
	SpellEffectMine
	SpellEffectElementalBarrierUp
	SpellEffectElementalBarrierDown
	SpellEffectDelayedExplosion
	SpellEffectMPEater
	SpellEffectHemorrhage
	SpellEffectBleeding
	SpellEffectAwakeningSuccess
	SpellEffectAwakeningFail
	SpellEffectAwakeningMiss
	SpellEffectAwakeningHit
	SpellEffectStormEscape
	SpellEffectTurtleKing
	SpellEffectBehemoth
	SpellEffectStunned
	SpellEffectIcePillar
)

type BuffType uint8

const (
	BuffTypeNone         BuffType = iota
	BuffTypeTemporalFlux          //magics
	BuffTypeHiding
	BuffTypeHaste
	BuffTypeSwiftFeet
	BuffTypeFury
	BuffTypeSoulShield
	BuffTypeBlessedArmour
	BuffTypeLightBody
	BuffTypeUltimateEnhancer
	BuffTypeProtectionField
	BuffTypeRage
	BuffTypeCurse
	BuffTypeMoonLight
	BuffTypeDarkBody
	BuffTypeConcentration
	BuffTypeVampireShot
	BuffTypePoisonShot
	BuffTypeCounterAttack
	BuffTypeMentalState
	BuffTypeEnergyShield
	BuffTypeMagicBooster
	BuffTypePetEnhancer
	BuffTypeImmortalSkin
	BuffTypeMagicShield
	BuffTypeGameMaster      = 100 //special
	BuffTypeGeneral         = 101
	BuffTypeExp             = 102
	BuffTypeDrop            = 103
	BuffTypeGold            = 104
	BuffTypeBagWeight       = 105
	BuffTypeTransform       = 106
	BuffTypeRelationshipEXP = 107
	BuffTypeMentee          = 108
	BuffTypeMentor          = 109
	BuffTypeGuild           = 110
	BuffTypePrison          = 111
	BuffTypeRested          = 112
	BuffTypeImpact          = 200 //stats
	BuffTypeMagic           = 201
	BuffTypeTaoist          = 202
	BuffTypeStorm           = 203
	BuffTypeHealthAid       = 204
	BuffTypeManaAid         = 205
	BuffTypeDefence         = 206
	BuffTypeMagicDefence    = 207
	BuffTypeWonderDrug      = 208
	BuffTypeKnapsack        = 209
)

type Monster uint16

const (
	MonsterGuard              Monster = 0
	MonsterTaoistGuard                = 1
	MonsterGuard2                     = 2
	MonsterHen                        = 3
	MonsterDeer                       = 4
	MonsterScarecrow                  = 5
	MonsterHookingCat                 = 6
	MonsterRakingCat                  = 7
	MonsterYob                        = 8
	MonsterOma                        = 9
	MonsterCannibalPlant              = 10
	MonsterForestYeti                 = 11
	MonsterSpittingSpider             = 12
	MonsterChestnutTree               = 13
	MonsterEbonyTree                  = 14
	MonsterLargeMushroom              = 15
	MonsterCherryTree                 = 16
	MonsterOmaFighter                 = 17
	MonsterOmaWarrior                 = 18
	MonsterCaveBat                    = 19
	MonsterCaveMaggot                 = 20
	MonsterScorpion                   = 21
	MonsterSkeleton                   = 22
	MonsterBoneFighter                = 23
	MonsterAxeSkeleton                = 24
	MonsterBoneWarrior                = 25
	MonsterBoneElite                  = 26
	MonsterDung                       = 27
	MonsterDark                       = 28
	MonsterWoomaSoldier               = 29
	MonsterWoomaFighter               = 30
	MonsterWoomaWarrior               = 31
	MonsterFlamingWooma               = 32
	MonsterWoomaGuardian              = 33
	MonsterWoomaTaurus                = 34
	MonsterWhimperingBee              = 35
	MonsterGiantWorm                  = 36
	MonsterCentipede                  = 37
	MonsterBlackMaggot                = 38
	MonsterTongs                      = 39
	MonsterEvilTongs                  = 40
	MonsterEvilCentipede              = 41
	MonsterBugBat                     = 42
	MonsterBugBatMaggot               = 43
	MonsterWedgeMoth                  = 44
	MonsterRedBoar                    = 45
	MonsterBlackBoar                  = 46
	MonsterSnakeScorpion              = 47
	MonsterWhiteBoar                  = 48
	MonsterEvilSnake                  = 49
	MonsterBombSpider                 = 50
	MonsterRootSpider                 = 51
	MonsterSpiderBat                  = 52
	MonsterVenomSpider                = 53
	MonsterGangSpider                 = 54
	MonsterGreatSpider                = 55
	MonsterLureSpider                 = 56
	MonsterBigApe                     = 57
	MonsterEvilApe                    = 58
	MonsterGrayEvilApe                = 59
	MonsterRedEvilApe                 = 60
	MonsterCrystalSpider              = 61
	MonsterRedMoonEvil                = 62
	MonsterBigRat                     = 63
	MonsterZumaArcher                 = 64
	MonsterZumaStatue                 = 65
	MonsterZumaGuardian               = 66
	MonsterRedThunderZuma             = 67
	MonsterZumaTaurus                 = 68
	MonsterDigOutZombie               = 69
	MonsterClZombie                   = 70
	MonsterNdZombie                   = 71
	MonsterCrawlerZombie              = 72
	MonsterShamanZombie               = 73
	MonsterGhoul                      = 74
	MonsterKingScorpion               = 75
	MonsterKingHog                    = 76
	MonsterDarkDevil                  = 77
	MonsterBoneFamiliar               = 78
	MonsterShinsu                     = 79
	MonsterShinsu1                    = 80
	MonsterSpiderFrog                 = 81
	MonsterHoroBlaster                = 82
	MonsterBlueHoroBlaster            = 83
	MonsterKekTal                     = 84
	MonsterVioletKekTal               = 85
	MonsterKhazard                    = 86
	MonsterRoninGhoul                 = 87
	MonsterToxicGhoul                 = 88
	MonsterBoneCaptain                = 89
	MonsterBoneSpearman               = 90
	MonsterBoneBlademan               = 91
	MonsterBoneArcher                 = 92
	MonsterBoneLord                   = 93
	MonsterMinotaur                   = 94
	MonsterIceMinotaur                = 95
	MonsterElectricMinotaur           = 96
	MonsterWindMinotaur               = 97
	MonsterFireMinotaur               = 98
	MonsterRightGuard                 = 99
	MonsterLeftGuard                  = 100
	MonsterMinotaurKing               = 101
	MonsterFrostTiger                 = 102
	MonsterSheep                      = 103
	MonsterWolf                       = 104
	MonsterShellNipper                = 105
	MonsterKeratoid                   = 106
	MonsterGiantKeratoid              = 107
	MonsterSkyStinger                 = 108
	MonsterSandWorm                   = 109
	MonsterVisceralWorm               = 110
	MonsterRedSnake                   = 111
	MonsterTigerSnake                 = 112
	MonsterYimoogi                    = 113
	MonsterGiantWhiteSnake            = 114
	MonsterBlueSnake                  = 115
	MonsterYellowSnake                = 116
	MonsterHolyDeva                   = 117
	MonsterAxeOma                     = 118
	MonsterSwordOma                   = 119
	MonsterCrossbowOma                = 120
	MonsterWingedOma                  = 121
	MonsterFlailOma                   = 122
	MonsterOmaGuard                   = 123
	MonsterYinDevilNode               = 124
	MonsterYangDevilNode              = 125
	MonsterOmaKing                    = 126
	MonsterBlackFoxman                = 127
	MonsterRedFoxman                  = 128
	MonsterWhiteFoxman                = 129
	MonsterTrapRock                   = 130
	MonsterGuardianRock               = 131
	MonsterThunderElement             = 132
	MonsterCloudElement               = 133
	MonsterGreatFoxSpirit             = 134
	MonsterHedgeKekTal                = 135
	MonsterBigHedgeKekTal             = 136
	MonsterRedFrogSpider              = 137
	MonsterBrownFrogSpider            = 138
	MonsterArcherGuard                = 139
	MonsterKatanaGuard                = 140
	MonsterArcherGuard2               = 141
	MonsterPig                        = 142
	MonsterBull                       = 143
	MonsterBush                       = 144
	MonsterChristmasTree              = 145
	MonsterHighAssassin               = 146
	MonsterDarkDustPile               = 147
	MonsterDarkBrownWolf              = 148
	MonsterFootball                   = 149
	MonsterGingerBreadman             = 150
	MonsterHalloweenScythe            = 151
	MonsterGhastlyLeecher             = 152
	MonsterCyanoGhast                 = 153
	MonsterMutatedManworm             = 154
	MonsterCrazyManworm               = 155
	MonsterMudPile                    = 156
	MonsterTailedLion                 = 157
	MonsterBehemoth                   = 158 //done BOSS
	MonsterDarkDevourer               = 159 //done
	MonsterPoisonHugger               = 160 //done
	MonsterHugger                     = 161 //done
	MonsterMutatedHugger              = 162 //done
	MonsterDreamDevourer              = 163 //done
	MonsterTreasurebox                = 164 //done
	MonsterSnowPile                   = 165 //done
	MonsterSnowman                    = 166 //done
	MonsterSnowTree                   = 167 //done
	MonsterGiantEgg                   = 168 //done
	MonsterRedTurtle                  = 169 //done
	MonsterGreenTurtle                = 170 //done
	MonsterBlueTurtle                 = 171 //done
	MonsterCatapult                   = 172 //not added frames //special 3 states in 1
	MonsterSabukWallSection           = 173 //not added frames
	MonsterNammandWallSection         = 174 //not added frames
	MonsterSiegeRepairman             = 175 //not added frames
	MonsterBlueSanta                  = 176 //done
	MonsterBattleStandard             = 177 //done
	//MonsterArcherGuard2 = 178,//done
	MonsterRedYimoogi         = 179 //done
	MonsterLionRiderMale      = 180 //frames not added
	MonsterLionRiderFemale    = 181 //frames not added
	MonsterTornado            = 182 //done
	MonsterFlameTiger         = 183 //done
	MonsterWingedTigerLord    = 184 //done BOSS
	MonsterTowerTurtle        = 185 //done
	MonsterFinialTurtle       = 186 //done
	MonsterTurtleKing         = 187 //done BOSS
	MonsterDarkTurtle         = 188 //done
	MonsterLightTurtle        = 189 //done
	MonsterDarkSwordOma       = 190 //done
	MonsterDarkAxeOma         = 191 //done
	MonsterDarkCrossbowOma    = 192 //done
	MonsterDarkWingedOma      = 193 //done
	MonsterBoneWhoo           = 194 //done
	MonsterDarkSpider         = 195 //done
	MonsterViscusWorm         = 196 //done
	MonsterViscusCrawler      = 197 //done
	MonsterCrawlerLave        = 198 //done
	MonsterDarkYob            = 199 //done
	MonsterFlamingMutant      = 200 //FINISH
	MonsterStoningStatue      = 201 //FINISH BOSS
	MonsterFlyingStatue       = 202 //FINISH
	MonsterValeBat            = 203 //done
	MonsterWeaver             = 204 //done
	MonsterVenomWeaver        = 205 //done
	MonsterCrackingWeaver     = 206 //done
	MonsterArmingWeaver       = 207 //done
	MonsterCrystalWeaver      = 208 //done
	MonsterFrozenZumaStatue   = 209 //done
	MonsterFrozenZumaGuardian = 210 //done
	MonsterFrozenRedZuma      = 211 //done
	MonsterGreaterWeaver      = 212 //done
	MonsterSpiderWarrior      = 213 //done
	MonsterSpiderBarbarian    = 214 //done
	MonsterHellSlasher        = 215 //done
	MonsterHellPirate         = 216 //done
	MonsterHellCannibal       = 217 //done
	MonsterHellKeeper         = 218 //done BOSS
	MonsterHellBolt           = 219 //done
	MonsterWitchDoctor        = 220 //done
	MonsterManectricHammer    = 221 //done
	MonsterManectricClub      = 222 //done
	MonsterManectricClaw      = 223 //done
	MonsterManectricStaff     = 224 //done
	MonsterNamelessGhost      = 225 //done
	MonsterDarkGhost          = 226 //done
	MonsterChaosGhost         = 227 //done
	MonsterManectricBlest     = 228 //done
	MonsterManectricKing      = 229 //done
	MonsterFrozenDoor         = 230 //done
	MonsterIcePillar          = 231 //done
	MonsterFrostYeti          = 232 //done
	MonsterManectricSlave     = 233 //done
	MonsterTrollHammer        = 234 //done
	MonsterTrollBomber        = 235 //done
	MonsterTrollStoner        = 236 //done
	MonsterTrollKing          = 237 //done BOSS
	MonsterFlameSpear         = 238 //done
	MonsterFlameMage          = 239 //done
	MonsterFlameScythe        = 240 //done
	MonsterFlameAssassin      = 241 //done
	MonsterFlameQueen         = 242 //finish BOSS
	MonsterHellKnight1        = 243 //done
	MonsterHellKnight2        = 244 //done
	MonsterHellKnight3        = 245 //done
	MonsterHellKnight4        = 246 //done
	MonsterHellLord           = 247 //done BOSS
	MonsterWaterGuard         = 248 //done
	MonsterIceGuard           = 249
	MonsterElementGuard       = 250
	MonsterDemonGuard         = 251
	MonsterKingGuard          = 252
	MonsterSnake10            = 253 //done
	MonsterSnake11            = 254 //done
	MonsterSnake12            = 255 //done
	MonsterSnake13            = 256 //done
	MonsterSnake14            = 257 //done
	MonsterSnake15            = 258 //done
	MonsterSnake16            = 259 //done
	MonsterSnake17            = 260 //done
	MonsterDeathCrawler       = 261
	MonsterBurningZombie      = 262
	MonsterMudZombie          = 263
	MonsterFrozenZombie       = 264
	MonsterUndeadWolf         = 265
	MonsterDemonwolf          = 266
	MonsterWhiteMammoth       = 267
	MonsterDarkBeast          = 268
	MonsterLightBeast         = 269
	MonsterBloodBaboon        = 270
	MonsterHardenRhino        = 271
	MonsterAncientBringer     = 272
	MonsterFightingCat        = 273
	MonsterFireCat            = 274
	MonsterCatWidow           = 275
	MonsterStainHammerCat     = 276
	MonsterBlackHammerCat     = 277
	MonsterStrayCat           = 278
	MonsterCatShaman          = 279
	MonsterJar1               = 280
	MonsterJar2               = 281
	MonsterSeedingsGeneral    = 282
	MonsterRestlessJar        = 283
	MonsterGeneralJinmYo      = 284
	MonsterBunny              = 285
	MonsterTucson             = 286
	MonsterTucsonFighter      = 287
	MonsterTucsonMage         = 288
	MonsterTucsonWarrior      = 289
	MonsterArmadillo          = 290
	MonsterArmadilloElder     = 291
	MonsterTucsonEgg          = 292
	MonsterPlaguedTucson      = 293
	MonsterSandSnail          = 294
	MonsterCannibalTentacles  = 295
	MonsterTucsonGeneral      = 296
	MonsterGasToad            = 297
	MonsterMantis             = 298
	MonsterSwampWarrior       = 299
	MonsterAssassinBird       = 300
	MonsterRhinoWarrior       = 301
	MonsterRhinoPriest        = 302
	MonsterSwampSlime         = 303
	MonsterRockGuard          = 304
	MonsterMudWarrior         = 305
	MonsterSmallPot           = 306
	MonsterTreeQueen          = 307
	MonsterShellFighter       = 308
	MonsterDarkBaboon         = 309
	MonsterTwinHeadBeast      = 310
	MonsterOmaCannibal        = 311
	MonsterOmaBlest           = 312
	MonsterOmaSlasher         = 313
	MonsterOmaAssassin        = 314
	MonsterOmaMage            = 315
	MonsterOmaWitchDoctor     = 316
	MonsterLightningBead      = 317
	MonsterHealingBead        = 318
	MonsterPowerUpBead        = 319
	MonsterDarkOmaKing        = 320
	MonsterCaveMage           = 321
	MonsterMandrill           = 322
	MonsterPlagueCrab         = 323
	MonsterCreeperPlant       = 324
	MonsterFloatingWraith     = 325
	MonsterArmedPlant         = 326
	MonsterAvengerPlant       = 327
	MonsterNadz               = 328
	MonsterAvengingSpirit     = 329
	MonsterAvengingWarrior    = 330
	MonsterAxePlant           = 331
	MonsterWoodBox            = 332
	MonsterClawBeast          = 333
	MonsterKillerPlant        = 334
	MonsterSackWarrior        = 335
	MonsterWereTiger          = 336
	MonsterKingHydrax         = 337
	MonsterHydrax             = 338
	MonsterHornedMage         = 339
	MonsterBasiloid           = 340
	MonsterHornedArcher       = 341
	MonsterColdArcher         = 342
	MonsterHornedWarrior      = 343
	MonsterFloatingRock       = 344
	MonsterScalyBeast         = 345
	MonsterHornedSorceror     = 346
	MonsterBoulderSpirit      = 347
	MonsterHornedCommander    = 348
	MonsterMoonStone          = 349
	MonsterSunStone           = 350
	MonsterLightningStone     = 351
	MonsterTurtlegrass        = 352
	MonsterMantree            = 353
	MonsterBear               = 354
	MonsterLeopard            = 355
	MonsterChieftainArcher    = 356
	MonsterChieftainSword     = 357
	MonsterStoningSpider      = 358 //Archer Spell mob (not yet coded)
	MonsterVampireSpider      = 359 //Archer Spell mob
	MonsterSpittingToad       = 360 //Archer Spell mob
	MonsterSnakeTotem         = 361 //Archer Spell mob
	MonsterCharmedSnake       = 362 //Archer Spell mob
	MonsterFrozenSoldier      = 363
	MonsterFrozenFighter      = 364
	MonsterFrozenArcher       = 365
	MonsterFrozenKnight       = 366
	MonsterFrozenGolem        = 367
	MonsterIcePhantom         = 368
	MonsterSnowWolf           = 369
	MonsterSnowWolfKing       = 370
	MonsterWaterDragon        = 371
	MonsterBlackTortoise      = 372
	MonsterManticore          = 373
	MonsterDragonWarrior      = 374
	MonsterDragonArcher       = 375
	MonsterKirin              = 376
	MonsterGuard3             = 377
	MonsterArcherGuard3       = 378
	MonsterBunny2             = 379
	MonsterFrozenMiner        = 380
	MonsterFrozenAxeman       = 381
	MonsterFrozenMagician     = 382
	MonsterSnowYeti           = 383
	MonsterIceCrystalSoldier  = 384
	MonsterDarkWraith         = 385
	MonsterDarkSpirit         = 386
	MonsterCrystalBeast       = 387
	MonsterRedOrb             = 388
	MonsterBlueOrb            = 389
	MonsterYellowOrb          = 390
	MonsterGreenOrb           = 391
	MonsterWhiteOrb           = 392
	MonsterFatalLotus         = 393
	MonsterAntCommander       = 394
	MonsterCargoBoxwithlogo   = 395
	MonsterDoe                = 396
	MonsterReindeer           = 397 //frames not added
	MonsterAngryReindeer      = 398
	MonsterCargoBox           = 399
	MonsterRam1               = 400
	MonsterRam2               = 401
	MonsterKite               = 402
	MonsterEvilMir            = 900
	MonsterEvilMirBody        = 901
	MonsterDragonStatue       = 902
	MonsterHellBomb1          = 903
	MonsterHellBomb2          = 904
	MonsterHellBomb3          = 905
	MonsterSabukGate          = 950
	MonsterPalaceWallLeft     = 951
	MonsterPalaceWall1        = 952
	MonsterPalaceWall2        = 953
	MonsterGiGateSouth        = 954
	MonsterGiGateEast         = 955
	MonsterGiGateWest         = 956
	MonsterSSabukWall1        = 957
	MonsterSSabukWall2        = 958
	MonsterSSabukWall3        = 959
	MonsterBabyPig            = 10000 //Permanent
	MonsterChick              = 10001 //Special
	MonsterKitten             = 10002 //Permanent
	MonsterBabySkeleton       = 10003 //Special
	MonsterBaekdon            = 10004 //Special
	MonsterWimaen             = 10005 //Event
	MonsterBlackKitten        = 10006 //unknown
	MonsterBabyDragon         = 10007 //unknown
	MonsterOlympicFlame       = 10008 //unknown
	MonsterBabySnowMan        = 10009 //unknown
	MonsterFrog               = 10010 //unknown
	MonsterBabyMonkey         = 10011 //unknown
	MonsterAngryBird          = 10012
	MonsterFoxey              = 10013
)

type DamageType uint8

const (
	DamageTypeHit      DamageType = 0
	DamageTypeMiss                = 1
	DamageTypeCritical            = 2
)

type ObjectType uint8

const (
	ObjectTypeNone     ObjectType = 0
	ObjectTypePlayer              = 1
	ObjectTypeItem                = 2
	ObjectTypeMerchant            = 3
	ObjectTypeSpell               = 4
	ObjectTypeMonster             = 5
	ObjectTypeDeco                = 6
	ObjectTypeCreature            = 7
)

type DefenceType uint8

const (
	DefenceTypeACAgility DefenceType = iota
	DefenceTypeAC
	DefenceTypeMACAgility
	DefenceTypeMAC
	DefenceTypeAgility
	DefenceTypeRepulsion
	DefenceTypeNone
)

const (
	BindModeNone                = 0    //
	BindModeDontDeathdrop       = 1    //0x0001
	BindModeDontDrop            = 2    //0x0002
	BindModeDontSell            = 4    //0x0004
	BindModeDontStore           = 8    //0x0008
	BindModeDontTrade           = 16   //0x0010
	BindModeDontRepair          = 32   //0x0020
	BindModeDontUpgrade         = 64   //0x0040
	BindModeDestroyOnDrop       = 128  //0x0080
	BindModeBreakOnDeath        = 256  //0x0100
	BindModeBindOnEquip         = 512  //0x0200
	BindModeNoSRepair           = 1024 //0x0400
	BindModeNoWeddingRing       = 2048 //0x0800
	BindModeUnableToRent        = 4096
	BindModeUnableToDisassemble = 8192
	BindModeNoMail              = 16384
)

// RankOptions 所在行会里的权限
type RankOptions uint8 // byte

const (
	RankOptionsCanChangeRank    RankOptions = 1
	RankOptionsCanRecruit                   = 2
	RankOptionsCanKick                      = 4
	RankOptionsCanStoreItem                 = 8
	RankOptionsCanRetrieveItem              = 16
	RankOptionsCanAlterAlliance             = 32
	RankOptionsCanChangeNotice              = 64
	RankOptionsCanActivateBuff              = 128
)
