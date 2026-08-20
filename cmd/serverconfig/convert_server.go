package main

import (
	"fmt"
	"path/filepath"
)

// ServerConfig 表示服务器主配置，节与字段顺序与 cmd/server/config.go 的
// ServerConfig 保持一致，确保 server.jsonc 的阅读顺序即服务端定义顺序。
type ServerConfig struct {
	Server    ServerSection    `json:"server"`
	Database  DatabaseSection  `json:"database"`
	Game      GameSection      `json:"game"`
	Login     LoginSection     `json:"login"`
	Player    PlayerSection    `json:"player"`
	Combat    CombatSection    `json:"combat"`
	PK        PKSection        `json:"pk"`
	Monster   MonsterSection   `json:"monster"`
	MonsterAI MonsterAISection `json:"monsterAI"`
	Magic     MagicSection     `json:"magic"`
	Economy   EconomySection   `json:"economy"`
	Drop      DropSection      `json:"drop"`
	Mining    MiningSection    `json:"mining"`
	Commands  CommandsSection  `json:"commands"`
	Plugins   PluginsSection   `json:"plugins"`
}

// ServerSection 包含服务器标识和网络设置。
type ServerSection struct {
	Name   string       `json:"name"`
	Index  int          `json:"index"`
	Listen ListenConfig `json:"listen"`
	Limits LimitsConfig `json:"limits"`
}

// ListenConfig 包含网络监听设置。
type ListenConfig struct {
	Addr string `json:"addr"`
}

// LimitsConfig 包含服务器容量限制。
type LimitsConfig struct {
	MaxPlayers int `json:"maxPlayers"`
	HumLimit   int `json:"humLimit"`
	MonLimit   int `json:"monLimit"`
}

// DatabaseSection 包含数据库文件路径。
type DatabaseSection struct {
	Path string `json:"path"`
}

// GameSection 包含游戏世界设置。字段集与 cmd/server/config.go 的 Game 节对齐。
type GameSection struct {
	HomeMap            string         `json:"homeMap"`
	HomeX              int            `json:"homeX"`
	HomeY              int            `json:"homeY"`
	RedHomeMap         string         `json:"redHomeMap"`
	RedHomeX           int            `json:"redHomeX"`
	RedHomeY           int            `json:"redHomeY"`
	GroupMembersMax    int            `json:"groupMembersMax"`
	BuildGuild         int            `json:"buildGuild"`
	GuildWarFee        int            `json:"guildWarFee"`
	DisableRun         bool           `json:"disableRun"`
	GMRunAll           bool           `json:"gmRunAll"`
	GMAccounts         map[string]int `json:"gmAccounts"`
	WalkInterval       int64          `json:"walkInterval"`
	RunInterval        int64          `json:"runInterval"`
	SpeedHackKick      bool           `json:"speedHackKick"`
	SpeedHackMax       int            `json:"speedHackMax"`
	ZenLimit           int            `json:"zenLimit"`
	MonGenRate         int            `json:"monGenRate"`
	UserFull           int            `json:"userFull"`
	ZenFastStep        int            `json:"zenFastStep"`
	WalkOnly           bool           `json:"walkOnly"`
	HitIntervalTime    int64          `json:"hitIntervalTime"`
	ActionInterval     int64          `json:"actionInterval"`
	RunLongHitInterval int64          `json:"runLongHitInterval"`
	RunHitInterval     int64          `json:"runHitInterval"`
	WalkHitInterval    int64          `json:"walkHitInterval"`
	RunMagicInterval   int64          `json:"runMagicInterval"`
	StruckTime         int64          `json:"struckTime"`
	// Go 服务端专属配置（Delphi 原版无此项）：受击开启头顶血条，
	// 默认开启；设为 false 回到仅心灵启示开启血条的原版行为。
	ShowHPOnStruck      bool    `json:"showHPOnStruck"`
	TurnInterval        int64   `json:"turnInterval"`
	SpellInterval       int64   `json:"spellInterval"`
	DigUpInterval       int64   `json:"digUpInterval"`
	MagicAttackRange    int     `json:"magicAttackRange"`
	TickInterval        int64   `json:"tickInterval"`
	SaveInterval        int     `json:"saveInterval"`
	GuildSaveInterval   int     `json:"guildSaveInterval"`
	DisableOnlineCount  bool    `json:"disableOnlineCount"`
	SendOnlineTime      int64   `json:"sendOnlineTime"`
	SendOnlineCountRate int     `json:"sendOnlineCountRate"`
	SendOnlineCountMsg  string  `json:"sendOnlineCountMsg"`
	MsgRateLimit        float64 `json:"msgRateLimit"`
	MsgBurst            int     `json:"msgBurst"`
}

// LoginSection 登录注册限制（Go 服务端专属配置节）。
type LoginSection struct {
	MaxAttempts     int `json:"maxAttempts"`
	LockoutSecs     int `json:"lockoutSecs"`
	MaxCharsPerAcct int `json:"maxCharsPerAcct"`
	MinNameLen      int `json:"minNameLen"`
	MaxNameLen      int `json:"maxNameLen"`
	MinAcctLen      int `json:"minAcctLen"`
	MaxAcctLen      int `json:"maxAcctLen"`
	MinPassLen      int `json:"minPassLen"`
	MaxPassLen      int `json:"maxPassLen"`
}

// PlayerSection 玩家属性与掉落设置。
type PlayerSection struct {
	MaxBagSlots        int   `json:"maxBagSlots"`
	MaxStorageSlots    int   `json:"maxStorageSlots"`
	MaxGold            int   `json:"maxGold"`
	MaxTradeItems      int   `json:"maxTradeItems"`
	MaxFriends         int   `json:"maxFriends"`
	MaxApprentices     int   `json:"maxApprentices"`
	MaxSlaves          int   `json:"maxSlaves"`
	MaxSlaveLevel      int   `json:"maxSlaveLevel"`
	BonusPerLevel      int   `json:"bonusPerLevel"`
	ViewRange          int   `json:"viewRange"`
	VisionInterval     int64 `json:"visionInterval"`
	VisionRand         int64 `json:"visionRand"`
	DeathSkeletonDelay int64 `json:"deathSkeletonDelay"`
	AutoReviveDelay    int64 `json:"autoReviveDelay"`
	ReviveHPRatio      int   `json:"reviveHPRatio"`
	HPRegenInterval    int64 `json:"hpRegenInterval"`
	MPRegenInterval    int64 `json:"mpRegenInterval"`
	HPRegenRatio       int   `json:"hpRegenRatio"`
	MPRegenRatio       int   `json:"mpRegenRatio"`
	DropEquipRate      int   `json:"dropEquipRate"`
	DropEquipRateRed   int   `json:"dropEquipRateRed"`
	DropBagRate        int   `json:"dropBagRate"`
	ExpPenaltyRatio    int   `json:"expPenaltyRatio"`
	NPCInteractDist    int   `json:"npcInteractDist"`
	SpeedDecayInterval int64 `json:"speedDecayInterval"`
	BonusHPPerPoint    int   `json:"bonusHPPerPoint"`
	BonusMPPerPoint    int   `json:"bonusMPPerPoint"`
}

// CombatSection 战斗与技能数值设置。
type CombatSection struct {
	FireHitWindow     int64 `json:"fireHitWindow"`
	FireHitCooldown   int64 `json:"fireHitCooldown"`
	TwinHitCooldown   int64 `json:"twinHitCooldown"`
	MagicShieldRatio  int   `json:"magicShieldRatio"`
	RedPoisonBonus    int   `json:"redPoisonBonus"`
	ParalysisDenom    int   `json:"paralysisDenom"`
	ParalysisDuration int64 `json:"paralysisDuration"`
	MapEnterProtect   int64 `json:"mapEnterProtect"`
	PKProtectLevel    int   `json:"pkProtectLevel"`
	PKProtectDiff     int   `json:"pkProtectDiff"`
	PKRedProtectLevel int   `json:"pkRedProtectLevel"`
	ProjectileDelay   int64 `json:"projectileDelay"`
	MaxLOSCheck       int   `json:"maxLOSCheck"`
}

// PKSection PK 值设置。
type PKSection struct {
	PointsPerKill       int   `json:"pointsPerKill"`
	DecayInterval       int64 `json:"decayInterval"`
	DecayAmount         int   `json:"decayAmount"`
	SelfDefenseDuration int64 `json:"selfDefenseDuration"`
	WeaponDuraLoss      int   `json:"weaponDuraLoss"`
}

// MonsterSection 怪物通用设置。
type MonsterSection struct {
	AITickInterval       int64 `json:"aiTickInterval"`
	CorpseDelay          int64 `json:"corpseDelay"`
	GhostDelay           int64 `json:"ghostDelay"`
	ButchInterval        int64 `json:"butchInterval"`
	HPRegenInterval      int64 `json:"hpRegenInterval"`
	HPRegenDivisor       int   `json:"hpRegenDivisor"`
	OverlapThinkInterval int64 `json:"overlapThinkInterval"`
	SearchInterval       int64 `json:"searchInterval"`
	SearchRand           int64 `json:"searchRand"`
	SearchNoTarget       int64 `json:"searchNoTarget"`
	SearchHasTarget      int64 `json:"searchHasTarget"`
	TargetFocusTimeout   int64 `json:"targetFocusTimeout"`
	TargetLossDist       int   `json:"targetLossDist"`
	DefaultViewRange     int   `json:"defaultViewRange"`
	DefaultWalkStep      int   `json:"defaultWalkStep"`
	DefaultWalkWait      int64 `json:"defaultWalkWait"`
	DefaultWalkSpeed     int64 `json:"defaultWalkSpeed"`
	DefaultAttackSpeed   int64 `json:"defaultAttackSpeed"`
	MinWalkSpeed         int64 `json:"minWalkSpeed"`
	AttackSpeedOffset    int64 `json:"attackSpeedOffset"`
	SpawnMaxTries        int   `json:"spawnMaxTries"`
	SlaveTeleportDist    int   `json:"slaveTeleportDist"`
	SlaveFollowDist      int   `json:"slaveFollowDist"`
	StruckPenaltyBase    int   `json:"struckPenaltyBase"`
	StruckPenaltyLevel   int   `json:"struckPenaltyLevel"`
	StruckPenaltyMin     int   `json:"struckPenaltyMin"`
	TargetSwitchChance   int   `json:"targetSwitchChance"`
	RangedDistance       int   `json:"rangedDistance"`
	AreaDistance         int   `json:"areaDistance"`
	GuardMaxHP           int   `json:"guardMaxHP"`
	GuardViewRange       int   `json:"guardViewRange"`
}

// MonsterAISection 特化怪物 AI 参数（对应 Delphi 各 Race 的行为常量）。
type MonsterAISection struct {
	BurrowTriggerDist    int   `json:"burrowTriggerDist"`
	ReburrowDist         int   `json:"reburrowDist"`
	ExplodeTimer         int64 `json:"explodeTimer"`
	ExplodePowerDiv      int   `json:"explodePowerDiv"`
	ExplodePowerMin      int   `json:"explodePowerMin"`
	DualAxeRange         int   `json:"dualAxeRange"`
	LeechRange           int   `json:"leechRange"`
	LeechBoostRatio      int   `json:"leechBoostRatio"`
	CritRange            int   `json:"critRange"`
	CritChance           int   `json:"critChance"`
	FireballRange        int   `json:"fireballRange"`
	SpitRange            int   `json:"spitRange"`
	PulseRange           int   `json:"pulseRange"`
	LightningRange       int   `json:"lightningRange"`
	CloneThreshold       int   `json:"cloneThreshold"`
	CloneCooldown        int64 `json:"cloneCooldown"`
	SummonMaxMinions     int   `json:"summonMaxMinions"`
	HiveMaxChildren      int   `json:"hiveMaxChildren"`
	CentipedeCooldown    int64 `json:"centipedeCooldown"`
	CentipedeAoEInterval int64 `json:"centipedeAoEInterval"`
	ZumaMaxSlaves        int   `json:"zumaMaxSlaves"`
	CowKingStunSpeed     int64 `json:"cowKingStunSpeed"`
	CowKingRageDuration  int64 `json:"cowKingRageDuration"`
	CowKingBerserkAtk    int64 `json:"cowKingBerserkAtk"`
	CowKingBerserkWalk   int64 `json:"cowKingBerserkWalk"`
	FireAuraDuration     int64 `json:"fireAuraDuration"`
	TransformCooldown    int64 `json:"transformCooldown"`
	BoneKingCooldown     int64 `json:"boneKingCooldown"`
	BoneKingMaxChildren  int   `json:"boneKingMaxChildren"`
	FleeChance           int   `json:"fleeChance"`
}

// MagicSection 技能数值设置。
type MagicSection struct {
	FireWallDuration int64 `json:"fireWallDuration"`
	HealingPoolCap   int   `json:"healingPoolCap"`
	TrainThreshold0  int   `json:"trainThreshold0"`
	TrainThreshold1  int   `json:"trainThreshold1"`
	TrainThreshold2  int   `json:"trainThreshold2"`
	SummonPetHPBase  int   `json:"summonPetHPBase"`
	SummonPetHPPerLv int   `json:"summonPetHPPerLv"`
	SummonPetDCBase  int   `json:"summonPetDCBase"`
	SummonPetDCPerLv int   `json:"summonPetDCPerLv"`
	AngelHPBase      int   `json:"angelHPBase"`
	AngelHPPerLv     int   `json:"angelHPPerLv"`
	AngelDCBase      int   `json:"angelDCBase"`
	AngelDCPerLv     int   `json:"angelDCPerLv"`
}

// EconomySection 经济系统设置。
type EconomySection struct {
	GuildCreateCost      int    `json:"guildCreateCost"`
	GuildWarDuration     int64  `json:"guildWarDuration"`
	UpgradeFee           int    `json:"upgradeFee"`
	UpgradeWaitTime      int64  `json:"upgradeWaitTime"`
	UpgradeMaxPoints     int    `json:"upgradeMaxPoints"`
	UpgradeMaterial      string `json:"upgradeMaterial"`
	UpgradeCurseChance   int    `json:"upgradeCurseChance"`
	RepairDuraDivisor    int    `json:"repairDuraDivisor"`
	SpecialRepairMult    int    `json:"specialRepairMult"`
	CastleDiscount       int    `json:"castleDiscount"`
	CastleMinRate        int    `json:"castleMinRate"`
	CastleTaxRate        int    `json:"castleTaxRate"`
	DrugBasePrice        int    `json:"drugBasePrice"`
	SpouseRecallCooldown int64  `json:"spouseRecallCooldown"`
	DoorCloseDelay       int64  `json:"doorCloseDelay"`
	PileStonesDuration   int64  `json:"pileStonesDuration"`
}

// DropSection 掉落规则设置。
type DropSection struct {
	EquipUpgChances     []int `json:"equipUpgChances"`
	EquipDuraMin        int   `json:"equipDuraMin"`
	EquipDuraRand       int   `json:"equipDuraRand"`
	AddValueChance      int   `json:"addValueChance"`
	MaxGoldPiles        int   `json:"maxGoldPiles"`
	MaxGoldPerPile      int   `json:"maxGoldPerPile"`
	FallbackGoldRate    int   `json:"fallbackGoldRate"`
	FallbackItemRate    int   `json:"fallbackItemRate"`
	FallbackExp         int   `json:"fallbackExp"`
	ControlDropItem     bool  `json:"controlDropItem"`
	GroundItemDespawnMs int64 `json:"groundItemDespawnMs"`
}

// MiningSection 挖矿设置。
type MiningSection struct {
	StoneRate int `json:"stoneRate"`
	OreRate   int `json:"oreRate"`
	DuraLoss  int `json:"duraLoss"`
}

// CommandsSection 包含 GM 命令定义。
type CommandsSection struct {
	Names       map[string]string `json:"names"`
	Permissions map[string]int    `json:"permissions"`
}

// PluginsSection 包含插件开关。
type PluginsSection struct {
	Enabled map[string]bool `json:"enabled"`
}

// ConvertServer 转换服务器主配置文件。
func ConvertServer(inputDir, outputDir string) error {
	// 解析 !setup.txt
	setupFile := filepath.Join(inputDir, "!setup.txt")
	setup, err := ParseINI(setupFile)
	if err != nil {
		return fmt.Errorf("parsing !setup.txt: %w", err)
	}

	// 解析 Command.ini
	cmdFile := filepath.Join(inputDir, "Command.ini")
	commands, err := ParseINI(cmdFile)
	if err != nil {
		return fmt.Errorf("parsing Command.ini: %w", err)
	}

	// Parse 系统插件.ini
	pluginFile := filepath.Join(inputDir, "系统插件.ini")
	plugins, err := ParseINI(pluginFile)
	if err != nil {
		// 插件文件为可选项
		plugins = make(map[string]map[string]string)
	}

	homeMap := getINIValue(setup, "Setup", "HomeMap", "0")
	homeX := getINIInt(setup, "Setup", "HomeX", 289)
	homeY := getINIInt(setup, "Setup", "HomeY", 618)

	// 构建服务器配置。新增节/字段的默认值取自 cmd/server/config.go 的
	// Get* 兜底值，写出来后服务端行为不变（getter 仅在零值时兜底）。
	config := ServerConfig{
		Server: ServerSection{
			Name:  getINIValue(setup, "Server", "ServerName", "Mir2 Server"),
			Index: getINIInt(setup, "Server", "ServerIndex", 0),
			Listen: ListenConfig{
				Addr: "0.0.0.0:7000",
			},
			Limits: LimitsConfig{
				MaxPlayers: getINIInt(setup, "Server", "UserFull", 10000),
				HumLimit:   getINIInt(setup, "Server", "HumLimit", 30),
				MonLimit:   getINIInt(setup, "Server", "MonLimit", 30),
			},
		},
		Database: DatabaseSection{
			Path: "serverdata/mir2.db",
		},
		Game: GameSection{
			HomeMap:            homeMap,
			HomeX:              homeX,
			HomeY:              homeY,
			RedHomeMap:         getINIValue(setup, "Setup", "sRedHomeMap", homeMap),
			RedHomeX:           getINIInt(setup, "Setup", "nRedHomeX", homeX),
			RedHomeY:           getINIInt(setup, "Setup", "nRedHomeY", homeY),
			GroupMembersMax:    getINIInt(setup, "Setup", "GroupMembersMax", 10),
			BuildGuild:         getINIInt(setup, "Setup", "BuildGuild", 1000000),
			GuildWarFee:        getINIInt(setup, "Setup", "GuildWarFee", 30000),
			DisableRun:         getINIBool(setup, "Setup", "DiableHumanRun", false),
			GMRunAll:           getINIBool(setup, "Setup", "GMRunAll", true),
			GMAccounts:         make(map[string]int),
			WalkInterval:       int64(getINIInt(setup, "Setup", "WalkIntervalTime", 600)),
			RunInterval:        int64(getINIInt(setup, "Setup", "RunIntervalTime", 600)),
			SpeedHackKick:      getINIBool(setup, "Setup", "KickOverSpeed", false),
			SpeedHackMax:       getINIInt(setup, "Setup", "OverSpeedKickCount", 4),
			ZenLimit:           getINIInt(setup, "Server", "ZenLimit", 50),
			MonGenRate:         getINIInt(setup, "Setup", "MonGenRate", 10),
			UserFull:           getINIInt(setup, "Server", "UserFull", 1000),
			ZenFastStep:        getINIInt(setup, "Server", "ZenFastStep", 300),
			WalkOnly:           getINIBool(setup, "Setup", "WalkOnly", false),
			HitIntervalTime:    int64(getINIInt(setup, "Setup", "HitIntervalTime", 1400)),
			ActionInterval:     int64(getINIInt(setup, "Setup", "ActionIntervalTime", 400)),
			RunLongHitInterval: int64(getINIInt(setup, "Setup", "RunLongHitIntervalTime", 800)),
			RunHitInterval:     int64(getINIInt(setup, "Setup", "RunHitIntervalTime", 800)),
			WalkHitInterval:    int64(getINIInt(setup, "Setup", "WalkHitIntervalTime", 800)),
			RunMagicInterval:   int64(getINIInt(setup, "Setup", "RunMagicIntervalTime", 900)),
			StruckTime:         int64(getINIInt(setup, "Setup", "StruckTime", 500)),
			ShowHPOnStruck:     true,
			TurnInterval:       int64(getINIInt(setup, "Setup", "TurnIntervalTime", 300)),
			SpellInterval:      int64(getINIInt(setup, "Setup", "SpellIntervalTime", 600)),
			DigUpInterval:      int64(getINIInt(setup, "Setup", "DigUpIntervalTime", 1000)),
			MagicAttackRange:   getINIInt(setup, "Setup", "MagicAttackRange", 12),
			// 以下为 Go 服务端专属参数（Delphi 无对应 INI 键）
			TickInterval:        100,
			SaveInterval:        300,
			GuildSaveInterval:   600,
			DisableOnlineCount:  !getINIBool(setup, "Setup", "SendOnlineCount", true),
			SendOnlineTime:      int64(getINIInt(setup, "Setup", "SendOnlineTime", 300000)),
			SendOnlineCountRate: getINIInt(setup, "Setup", "SendOnlineCountRate", 10),
			SendOnlineCountMsg:  getINIValue(setup, "Setup", "SendOnlineCountMsg", "当前在线人数: %c"),
			MsgRateLimit:        60,
			MsgBurst:            40,
		},
		Login: LoginSection{
			MaxAttempts:     5,
			LockoutSecs:     60,
			MaxCharsPerAcct: 2,
			MinNameLen:      3,
			MaxNameLen:      14,
			MinAcctLen:      3,
			MaxAcctLen:      10,
			MinPassLen:      3,
			MaxPassLen:      10,
		},
		Player: PlayerSection{
			MaxBagSlots:        46,
			MaxStorageSlots:    39,
			MaxGold:            10000000,
			MaxTradeItems:      12,
			MaxFriends:         50,
			MaxApprentices:     3,
			MaxSlaves:          2,
			MaxSlaveLevel:      7,
			BonusPerLevel:      3,
			ViewRange:          12,
			VisionInterval:     2000,
			VisionRand:         2000,
			DeathSkeletonDelay: 10000,
			AutoReviveDelay:    180000,
			ReviveHPRatio:      2,
			HPRegenInterval:    10000,
			MPRegenInterval:    10000,
			HPRegenRatio:       20,
			MPRegenRatio:       15,
			DropEquipRate:      30,
			DropEquipRateRed:   15,
			DropBagRate:        10,
			ExpPenaltyRatio:    20,
			NPCInteractDist:    15,
			SpeedDecayInterval: 10000,
			BonusHPPerPoint:    5,
			BonusMPPerPoint:    5,
		},
		Combat: CombatSection{
			FireHitWindow:     20000,
			FireHitCooldown:   10000,
			TwinHitCooldown:   60000,
			MagicShieldRatio:  150,
			RedPoisonBonus:    5,
			ParalysisDenom:    5,
			ParalysisDuration: 50,
			MapEnterProtect:   3000,
			PKProtectLevel:    7,
			PKProtectDiff:     10,
			PKRedProtectLevel: 15,
			ProjectileDelay:   600,
			MaxLOSCheck:       13,
		},
		PK: PKSection{
			PointsPerKill:       100,
			DecayInterval:       120000,
			DecayAmount:         1,
			SelfDefenseDuration: 60000,
			WeaponDuraLoss:      100,
		},
		Monster: MonsterSection{
			AITickInterval:       250,
			CorpseDelay:          180000,
			GhostDelay:           300000,
			ButchInterval:        600,
			HPRegenInterval:      6000,
			HPRegenDivisor:       75,
			OverlapThinkInterval: 3000,
			SearchInterval:       3000,
			SearchRand:           2000,
			SearchNoTarget:       1000,
			SearchHasTarget:      8000,
			TargetFocusTimeout:   30000,
			TargetLossDist:       15,
			DefaultViewRange:     5,
			DefaultWalkStep:      3,
			DefaultWalkWait:      1000,
			DefaultWalkSpeed:     1400,
			DefaultAttackSpeed:   2000,
			MinWalkSpeed:         200,
			AttackSpeedOffset:    900,
			SpawnMaxTries:        31,
			SlaveTeleportDist:    20,
			SlaveFollowDist:      3,
			StruckPenaltyBase:    150,
			StruckPenaltyLevel:   130,
			StruckPenaltyMin:     20,
			TargetSwitchChance:   6,
			RangedDistance:       5,
			AreaDistance:         6,
			GuardMaxHP:           65535,
			GuardViewRange:       7,
		},
		MonsterAI: MonsterAISection{
			BurrowTriggerDist:    3,
			ReburrowDist:         4,
			ExplodeTimer:         60000,
			ExplodePowerDiv:      2,
			ExplodePowerMin:      50,
			DualAxeRange:         7,
			LeechRange:           2,
			LeechBoostRatio:      150,
			CritRange:            7,
			CritChance:           4,
			FireballRange:        8,
			SpitRange:            2,
			PulseRange:           16,
			LightningRange:       8,
			CloneThreshold:       3,
			CloneCooldown:        20000,
			SummonMaxMinions:     3,
			HiveMaxChildren:      15,
			CentipedeCooldown:    10000,
			CentipedeAoEInterval: 3000,
			ZumaMaxSlaves:        30,
			CowKingStunSpeed:     10000,
			CowKingRageDuration:  8000,
			CowKingBerserkAtk:    500,
			CowKingBerserkWalk:   400,
			FireAuraDuration:     20000,
			TransformCooldown:    5000,
			BoneKingCooldown:     15000,
			BoneKingMaxChildren:  8,
			FleeChance:           30,
		},
		Magic: MagicSection{
			FireWallDuration: 8000,
			HealingPoolCap:   300,
			TrainThreshold0:  20,
			TrainThreshold1:  50,
			TrainThreshold2:  100,
			SummonPetHPBase:  50,
			SummonPetHPPerLv: 5,
			SummonPetDCBase:  5,
			SummonPetDCPerLv: 2,
			AngelHPBase:      80,
			AngelHPPerLv:     8,
			AngelDCBase:      8,
			AngelDCPerLv:     1,
		},
		Economy: EconomySection{
			GuildCreateCost:      1000000,
			GuildWarDuration:     10800000,
			UpgradeFee:           10000,
			UpgradeWaitTime:      3600000,
			UpgradeMaxPoints:     7,
			UpgradeMaterial:      "黑铁矿",
			UpgradeCurseChance:   30,
			RepairDuraDivisor:    30,
			SpecialRepairMult:    3,
			CastleDiscount:       80,
			CastleMinRate:        60,
			CastleTaxRate:        5,
			DrugBasePrice:        500,
			SpouseRecallCooldown: 60000,
			DoorCloseDelay:       5000,
			PileStonesDuration:   300000,
		},
		Drop: DropSection{
			EquipUpgChances:     []int{15, 24, 20, 30, 40},
			EquipDuraMin:        20,
			EquipDuraRand:       80,
			AddValueChance:      10,
			MaxGoldPiles:        17,
			MaxGoldPerPile:      2000,
			FallbackGoldRate:    30,
			FallbackItemRate:    10,
			FallbackExp:         10,
			ControlDropItem:     false,
			GroundItemDespawnMs: 3600000,
		},
		Mining: MiningSection{
			StoneRate: 4,
			OreRate:   12,
			DuraLoss:  100,
		},
		Commands: CommandsSection{
			Names:       getINISection(commands, "Command"),
			Permissions: getINISectionInt(commands, "Permission"),
		},
		Plugins: PluginsSection{
			// 系统插件.ini 的节名为 [setup]（非 [Plugins]）。
			Enabled: getINISectionBool(plugins, "setup"),
		},
	}

	// 序列化并写入输出文件
	outputFile := filepath.Join(outputDir, "server.jsonc")
	comment := "服务器主配置\n来源: asset/server/!setup.txt + Command.ini + 系统插件.ini\n说明: 合并后的服务端配置，移除了多进程地址配置；未标注来源的节为 Go 服务端专属参数（默认值与 cmd/server/config.go 一致）"

	data, err := MarshalJSONC(config, serverComments, false)
	if err != nil {
		return fmt.Errorf("marshaling server config: %w", err)
	}
	return WriteJSONC(outputFile, data, comment)
}

// serverComments 是 server.jsonc 的中文注释表。
// 键为点号路径；文案来源：cmd/server/config.go 字段/getter 注释、
// 服务端消费逻辑、Delphi M2Share.pas 对应键。
var serverComments = map[string]string{
	// --- server ---
	"server":                   "服务器标识与网络设置",
	"server.name":              "服务器名，显示在客户端服务器列表（Delphi ServerName）",
	"server.index":             "服务器序号（Delphi ServerIndex）",
	"server.listen":            "网络监听（Delphi 原 8 进程合并为单端口监听）",
	"server.listen.addr":       "监听地址 host:port",
	"server.limits":            "容量限制",
	"server.limits.maxPlayers": "服务器容量上限（Delphi [Server] UserFull）",
	"server.limits.humLimit":   "单地图玩家数上限（Delphi HumLimit；服务端暂未消费，保留）",
	"server.limits.monLimit":   "单地图怪物数上限（Delphi MonLimit；服务端暂未消费，保留）",

	// --- database ---
	"database":      "数据库设置",
	"database.path": "SQLite 数据库文件路径（Go 服务端以 SQLite 取代 Delphi 平面文件数据库）",

	// --- game ---
	"game":                     "游戏世界设置（多数字段对应 Delphi g_Config，来自 !setup.txt）",
	"game.homeMap":             "回城点地图编号（回城卷目的地 / 新角色出生点）",
	"game.homeX":               "回城点 X 坐标",
	"game.homeY":               "回城点 Y 坐标",
	"game.redHomeMap":          "红名玩家回城地图（Delphi sRedHomeMap；PKLevel>=2 时使用回城卷的目的地）",
	"game.redHomeX":            "红名回城点 X 坐标（Delphi nRedHomeX）",
	"game.redHomeY":            "红名回城点 Y 坐标（Delphi nRedHomeY）",
	"game.groupMembersMax":     "队伍最大成员数（Delphi GroupMembersMax）",
	"game.buildGuild":          "创建行会金币费用（Delphi BuildGuild；Go 服务端实际消费 economy.guildCreateCost）",
	"game.guildWarFee":         "宣布行会战争金币费用（Delphi GuildWarFee）",
	"game.disableRun":          "全局禁止跑步（Delphi DiableHumanRun，键名为 Delphi 原版拼写）",
	"game.gmRunAll":            "GM 是否允许全部跑位命令（Delphi GMRunAll）",
	"game.gmAccounts":          "GM 账号 → 权限等级映射（0-10，越高权限越大）；空表表示不设 GM",
	"game.walkInterval":        "玩家走路最小间隔 ms（Delphi WalkIntervalTime，限速判定用）",
	"game.runInterval":         "玩家跑步最小间隔 ms（Delphi RunIntervalTime，限速判定用）",
	"game.speedHackKick":       "超速（加速外挂）是否踢下线（Delphi KickOverSpeed）",
	"game.speedHackMax":        "连续超速达到该次数后踢下线（Delphi OverSpeedKickCount）",
	"game.zenLimit":            "单次刷怪 tick 的时间预算 ms，超时中止本轮刷怪（Delphi g_dwZenLimit）",
	"game.monGenRate":          "刷怪数量比率，10 为基准，实际数量按 10/rate 缩放（Delphi nMonGenRate）",
	"game.userFull":            "在线人数达到该值后刷怪开始加速（Delphi nUserFull；与 limits.maxPlayers 共用 INI 键）",
	"game.zenFastStep":         "在线人数每多出该数量，刷怪间隔加速一档（Delphi nZenFastStep）",
	"game.walkOnly":            "只允许走路禁止跑步（permission>=10 的 GM 不受限）",
	"game.hitIntervalTime":     "普通攻击最小间隔 ms（Delphi HitIntervalTime，限速判定用）",
	"game.actionInterval":      "通用动作最小间隔 ms（Delphi dwActionIntervalTime）",
	"game.runLongHitInterval":  "跑位刺杀最小间隔 ms（Delphi dwRunLongHitIntervalTime）",
	"game.runHitInterval":      "跑位普攻最小间隔 ms（Delphi dwRunHitIntervalTime）",
	"game.walkHitInterval":     "走位普攻最小间隔 ms（Delphi dwWalkHitIntervalTime）",
	"game.runMagicInterval":    "跑位施法最小间隔 ms（Delphi dwRunMagicIntervalTime）",
	"game.struckTime":          "被击硬直时间 ms（Delphi StruckTime）",
	"game.showHPOnStruck":      "受击开启头顶血条（Go 专属：Delphi 原版仅心灵启示开启血条；false 回到原版行为）",
	"game.turnInterval":        "转身最小间隔 ms（Delphi TurnIntervalTime）",
	"game.spellInterval":       "施法最小间隔 ms（Delphi SpellIntervalTime）",
	"game.digUpInterval":       "挖肉最小间隔 ms（Delphi DigUpIntervalTime）",
	"game.magicAttackRange":    "技能锁定目标的最大距离（格）",
	"game.tickInterval":        "服务器主循环心跳间隔 ms（Go 专属）",
	"game.saveInterval":        "玩家数据定期保存间隔，tickInterval 的倍数（Go 专属；300 ≈ 30 秒）",
	"game.guildSaveInterval":   "行会数据定期保存间隔，tickInterval 的倍数（Go 专属；600 ≈ 60 秒）",
	"game.disableOnlineCount":  "禁止在线人数广播（Delphi SendOnlineCount 的反义）",
	"game.sendOnlineTime":      "在线人数广播间隔 ms（Delphi dwSendOnlineTime）",
	"game.sendOnlineCountRate": "在线人数广播频率系数，实际间隔 = sendOnlineTime × rate/10（Delphi nSendOnlineCountRate）",
	"game.sendOnlineCountMsg":  "在线人数广播文案，%c 替换为人数（Delphi g_sSendOnlineCountMsg）",
	"game.msgRateLimit":        "每连接入站消息限速，条/秒（Go 专属网关补偿层）",
	"game.msgBurst":            "每连接入站消息令牌桶突发容量（Go 专属网关补偿层）",

	// --- login ---
	"login":                 "登录注册限制（Go 专属配置节）",
	"login.maxAttempts":     "密码重试次数上限，超过后临时锁定",
	"login.lockoutSecs":     "超过重试上限后的锁定时长（秒）",
	"login.maxCharsPerAcct": "每个账号最多可创建角色数",
	"login.minNameLen":      "角色名最小长度",
	"login.maxNameLen":      "角色名最大长度",
	"login.minAcctLen":      "账号名最小长度",
	"login.maxAcctLen":      "账号名最大长度",
	"login.minPassLen":      "密码最小长度",
	"login.maxPassLen":      "密码最大长度",

	// --- player ---
	"player":                    "玩家设置",
	"player.maxBagSlots":        "背包格数（Delphi 背包 46 格）",
	"player.maxStorageSlots":    "仓库格数（Delphi 运行时上限 39）",
	"player.maxGold":            "金币持有上限（Delphi nHumanMaxGold）",
	"player.maxTradeItems":      "单笔交易最多可放置物品数",
	"player.maxFriends":         "好友列表上限",
	"player.maxApprentices":     "师父最多可收徒弟数",
	"player.maxSlaves":          "玩家最多可拥有召唤兽/宠物数",
	"player.maxSlaveLevel":      "召唤兽/宠物等级上限",
	"player.bonusPerLevel":      "每升一级获得的属性点数",
	"player.viewRange":          "玩家视野半径（格）",
	"player.visionInterval":     "视野刷新基础间隔 ms",
	"player.visionRand":         "视野刷新随机抖动上限 ms（实际间隔 = visionInterval + rand(visionRand)）",
	"player.deathSkeletonDelay": "死亡后骨架化时间 ms（骨架化后不再播放死亡动画）",
	"player.autoReviveDelay":    "死亡后自动复活等待时间 ms",
	"player.reviveHPRatio":      "复活时 HP 恢复为上限的 1/N",
	"player.hpRegenInterval":    "HP 自然恢复间隔 ms",
	"player.mpRegenInterval":    "MP 自然恢复间隔 ms",
	"player.hpRegenRatio":       "每次 HP 自然恢复量为上限的 1/N",
	"player.mpRegenRatio":       "每次 MP 自然恢复量为上限的 1/N",
	"player.dropEquipRate":      "死亡时身上装备掉落概率 1/N",
	"player.dropEquipRateRed":   "红名玩家死亡时装备掉落概率 1/N（比白名更高）",
	"player.dropBagRate":        "死亡时背包物品掉落概率 1/N",
	"player.expPenaltyRatio":    "死亡经验损失比例为当前经验的 1/N",
	"player.npcInteractDist":    "NPC 对话交互最大距离（格）",
	"player.speedDecayInterval": "超速计数衰减间隔 ms（间隔内无新超速则计数 -1）",
	"player.bonusHPPerPoint":    "每点体质属性增加的 HP",
	"player.bonusMPPerPoint":    "每点魔力属性增加的 MP",

	// --- combat ---
	"combat":                   "战斗设置",
	"combat.fireHitWindow":     "烈火剑法蓄力状态有效窗口 ms（窗口内下次普攻附带火系加成）",
	"combat.fireHitCooldown":   "烈火剑法触发冷却 ms",
	"combat.twinHitCooldown":   "双刀斩（TwinHit）冷却 ms",
	"combat.magicShieldRatio":  "魔法盾额外耗蓝百分比：受击耗蓝 = 伤害 × (1 + ratio/100)",
	"combat.redPoisonBonus":    "红毒增伤：伤害额外附加 伤害/ratio",
	"combat.paralysisDenom":    "麻痹抗性判定分母：rand(抗性 + denom) == 0 时麻痹生效",
	"combat.paralysisDuration": "麻痹持续时间（内部时间单位）",
	"combat.mapEnterProtect":   "进入地图后的 PK 保护时间 ms（保护期内不可被攻击）",
	"combat.pkProtectLevel":    "低等级 PK 保护：攻击者等级超过该值才计入 PK 值判定",
	"combat.pkProtectDiff":     "PK 保护的等级差阈值",
	"combat.pkRedProtectLevel": "红名玩家不再受 PK 保护的等级",
	"combat.projectileDelay":   "投射物（火球等）飞行延迟 ms",
	"combat.maxLOSCheck":       "视线（LOS）检测的最大格数",

	// --- pk ---
	"pk":                     "PK 值设置",
	"pk.pointsPerKill":       "击杀白名玩家增加的 PK 值",
	"pk.decayInterval":       "PK 值自然衰减间隔 ms",
	"pk.decayAmount":         "每次衰减的 PK 值",
	"pk.selfDefenseDuration": "正当防卫判定窗口 ms（被攻击后该时长内反击不计 PK 值）",
	"pk.weaponDuraLoss":      "PK 杀人后武器耐久损失",

	// --- monster ---
	"monster":                      "怪物通用设置",
	"monster.aiTickInterval":       "怪物 AI 思考间隔 ms",
	"monster.corpseDelay":          "尸体保留时间 ms（到期后骨架化）",
	"monster.ghostDelay":           "幽灵态（不可见残留）清除时间 ms",
	"monster.butchInterval":        "屠宰（Alt 挖肉）最小间隔 ms（Delphi dwTurnIntervalTime）",
	"monster.hpRegenInterval":      "怪物 HP 自然恢复间隔 ms",
	"monster.hpRegenDivisor":       "怪物每次 HP 恢复量为上限的 1/N",
	"monster.overlapThinkInterval": "怪物重叠时的避让思考间隔 ms",
	"monster.searchInterval":       "索敌扫描基础间隔 ms",
	"monster.searchRand":           "索敌扫描随机抖动上限 ms",
	"monster.searchNoTarget":       "无目标时的索敌间隔 ms",
	"monster.searchHasTarget":      "有目标时的索敌间隔 ms",
	"monster.targetFocusTimeout":   "锁定单一目标的专注超时 ms（超时后重新选择目标）",
	"monster.targetLossDist":       "目标超出该距离（格）后丢失仇恨",
	"monster.defaultViewRange":     "怪物默认视野半径（格，数据库未指定时使用）",
	"monster.defaultWalkStep":      "怪物默认连续行走步数",
	"monster.defaultWalkWait":      "怪物默认行走停顿 ms",
	"monster.defaultWalkSpeed":     "怪物默认行走速度 ms/步",
	"monster.defaultAttackSpeed":   "怪物默认攻击速度 ms/次",
	"monster.minWalkSpeed":         "怪物行走速度下限 ms/步",
	"monster.attackSpeedOffset":    "怪物攻击速度随机偏移基准 ms",
	"monster.spawnMaxTries":        "刷怪选点的最大尝试次数",
	"monster.slaveTeleportDist":    "召唤兽跟随距离超过该值时瞬移到主人身边（格）",
	"monster.slaveFollowDist":      "召唤兽与主人的最大跟随距离（格）",
	"monster.struckPenaltyBase":    "怪物被击攻击延迟惩罚基准 ms（Delphi: += base - min(level*4, penaltyLevel)，下限 penaltyMin）",
	"monster.struckPenaltyLevel":   "被击惩罚的等级减免上限",
	"monster.struckPenaltyMin":     "被击惩罚下限 ms",
	"monster.targetSwitchChance":   "被击时 1/N 概率切换仇恨目标",
	"monster.rangedDistance":       "远程怪物默认攻击距离（格）",
	"monster.areaDistance":         "范围攻击怪物的默认作用半径（格）",
	"monster.guardMaxHP":           "城门守卫 HP 上限",
	"monster.guardViewRange":       "城门守卫视野半径（格）",

	// --- monsterAI ---
	"monsterAI":                      "特化怪物 AI 参数（对应 Delphi 各 Race 的行为常量）",
	"monsterAI.burrowTriggerDist":    "潜地怪（如沙漠蠕虫）目标接近到该距离时钻出地面（格）",
	"monsterAI.reburrowDist":         "潜地怪目标远离超过该距离时重新钻地（格）",
	"monsterAI.explodeTimer":         "自爆怪引爆倒计时 ms",
	"monsterAI.explodePowerDiv":      "自爆伤害 = 自身 HP / 该值",
	"monsterAI.explodePowerMin":      "自爆伤害下限",
	"monsterAI.dualAxeRange":         "双斧骷髅类远程攻击距离（格）",
	"monsterAI.leechRange":           "吸血怪攻击距离（格）",
	"monsterAI.leechBoostRatio":      "吸血怪伤害加成百分比",
	"monsterAI.critRange":            "暴击类怪物触发距离（格）",
	"monsterAI.critChance":           "暴击概率 1/N",
	"monsterAI.fireballRange":        "火球类怪物远程攻击距离（格）",
	"monsterAI.spitRange":            "喷吐类怪物攻击距离（格）",
	"monsterAI.pulseRange":           "脉冲类怪物范围攻击半径（格）",
	"monsterAI.lightningRange":       "闪电类怪物远程攻击距离（格）",
	"monsterAI.cloneThreshold":       "分裂怪 HP 低于上限的 1/N 时触发分裂，分裂体 HP 同为上限的 1/N",
	"monsterAI.cloneCooldown":        "分裂冷却 ms",
	"monsterAI.summonMaxMinions":     "召唤类怪物最多仆从数",
	"monsterAI.hiveMaxChildren":      "巢穴类怪物（母虫）最多子虫数",
	"monsterAI.centipedeCooldown":    "蜈蚣精特化技能冷却 ms",
	"monsterAI.centipedeAoEInterval": "蜈蚣精范围攻击间隔 ms",
	"monsterAI.zumaMaxSlaves":        "祖玛系怪物最多召唤护卫数",
	"monsterAI.cowKingStunSpeed":     "牛魔王眩晕持续 ms",
	"monsterAI.cowKingRageDuration":  "牛魔王狂暴持续 ms",
	"monsterAI.cowKingBerserkAtk":    "牛魔王狂暴攻击速度 ms/次",
	"monsterAI.cowKingBerserkWalk":   "牛魔王狂暴行走速度 ms/步",
	"monsterAI.fireAuraDuration":     "火焰光环持续 ms",
	"monsterAI.transformCooldown":    "变身类怪物变身冷却 ms",
	"monsterAI.boneKingCooldown":     "骨王召唤冷却 ms",
	"monsterAI.boneKingMaxChildren":  "骨王最多召唤随从数",
	"monsterAI.fleeChance":           "鹿类怪物（Race 52）1/N 概率生成逃跑型",

	// --- magic ---
	"magic":                  "技能数值设置",
	"magic.fireWallDuration": "火墙持续时间 ms",
	"magic.healingPoolCap":   "治疗池（群体治愈类）单次累计治疗上限",
	"magic.trainThreshold0":  "技能 1 级所需熟练度",
	"magic.trainThreshold1":  "技能 2 级所需熟练度",
	"magic.trainThreshold2":  "技能 3 级所需熟练度",
	"magic.summonPetHPBase":  "召唤宠物基础 HP",
	"magic.summonPetHPPerLv": "召唤宠物每技能级 +HP",
	"magic.summonPetDCBase":  "召唤宠物基础攻击（DC）",
	"magic.summonPetDCPerLv": "召唤宠物每技能级 +DC",
	"magic.angelHPBase":      "召唤天使基础 HP",
	"magic.angelHPPerLv":     "召唤天使每技能级 +HP",
	"magic.angelDCBase":      "召唤天使基础攻击（DC）",
	"magic.angelDCPerLv":     "召唤天使每技能级 +DC",

	// --- economy ---
	"economy":                      "经济系统设置",
	"economy.guildCreateCost":      "创建行会金币费用",
	"economy.guildWarDuration":     "行会战争持续时间 ms（默认 3 小时）",
	"economy.upgradeFee":           "武器升级手续费（金币）",
	"economy.upgradeWaitTime":      "武器升级等待时间 ms（默认 1 小时）",
	"economy.upgradeMaxPoints":     "武器升级最高点数",
	"economy.upgradeMaterial":      "武器升级所需矿石名",
	"economy.upgradeCurseChance":   "武器升级失败时附加诅咒的概率 %",
	"economy.repairDuraDivisor":    "普通修理的耐久损耗：损耗 = 已损失耐久 / 该值",
	"economy.specialRepairMult":    "特殊修理费用倍率（相对普通修理）",
	"economy.castleDiscount":       "沙巴克城内商店价格百分比（80 = 八折）",
	"economy.castleMinRate":        "沙巴克折扣下限百分比",
	"economy.castleTaxRate":        "沙巴克交易税率 %",
	"economy.drugBasePrice":        "制药基础价格（金币）",
	"economy.spouseRecallCooldown": "夫妻传送（天地合一）冷却 ms",
	"economy.doorCloseDelay":       "城门自动关闭延迟 ms",
	"economy.pileStonesDuration":   "挖矿落石堆存在时长 ms",

	// --- drop ---
	"drop":                     "掉落设置",
	"drop.equipUpgChances":     "掉落装备随机强化的 5 档权重（randomUpgradeItem 用）",
	"drop.equipDuraMin":        "掉落装备耐久基数",
	"drop.equipDuraRand":       "掉落装备耐久随机上限（实际耐久 = 基数 + rand(该值)）",
	"drop.addValueChance":      "掉落装备 1/N 概率出现附加属性",
	"drop.maxGoldPiles":        "金币掉落时最多堆数",
	"drop.maxGoldPerPile":      "每堆金币数量上限",
	"drop.fallbackGoldRate":    "无掉落表怪物掉金币概率 %",
	"drop.fallbackItemRate":    "无掉落表怪物掉物品概率 %",
	"drop.fallbackExp":         "无掉落表怪物的额外经验补偿",
	"drop.controlDropItem":     "廉价物品管控（Delphi boControlDropItem）：开启后低价物品直接删除不落地",
	"drop.groundItemDespawnMs": "地面物品消失时间 ms（Delphi dwClearDropOnFloorItemTime）",

	// --- mining ---
	"mining":           "挖矿设置",
	"mining.stoneRate": "挖出石头的概率 1/N",
	"mining.oreRate":   "挖出矿石的概率 1/N",
	"mining.duraLoss":  "每次挥镐的镐子耐久损耗",

	// --- commands ---
	"commands":             "GM 命令设置（来源: Command.ini）",
	"commands.names":       "GM 命令名 → 实际执行命令名映射",
	"commands.permissions": "GM 命令 → 所需权限等级（0-10）",

	// --- plugins ---
	"plugins":         "插件设置（来源: 系统插件.ini）",
	"plugins.enabled": "插件名 → 是否启用",
}

// 辅助函数

func getINIValue(ini map[string]map[string]string, section, key, defaultVal string) string {
	if sec, ok := ini[section]; ok {
		if val, ok := sec[key]; ok {
			return val
		}
	}
	return defaultVal
}

func getINIInt(ini map[string]map[string]string, section, key string, defaultVal int) int {
	if sec, ok := ini[section]; ok {
		if val, ok := sec[key]; ok {
			var n int
			if _, err := fmt.Sscanf(val, "%d", &n); err == nil {
				return n
			}
		}
	}
	return defaultVal
}

func getINISection(ini map[string]map[string]string, section string) map[string]string {
	if sec, ok := ini[section]; ok {
		result := make(map[string]string)
		for k, v := range sec {
			result[k] = v
		}
		return result
	}
	return make(map[string]string)
}

func getINISectionInt(ini map[string]map[string]string, section string) map[string]int {
	if sec, ok := ini[section]; ok {
		result := make(map[string]int)
		for k, v := range sec {
			var n int
			if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
				result[k] = n
			}
		}
		return result
	}
	return make(map[string]int)
}

func getINIBool(ini map[string]map[string]string, section, key string, defaultVal bool) bool {
	if sec, ok := ini[section]; ok {
		if val, ok := sec[key]; ok {
			return val == "1" || val == "true" || val == "TRUE"
		}
	}
	return defaultVal
}

func getINISectionBool(ini map[string]map[string]string, section string) map[string]bool {
	if sec, ok := ini[section]; ok {
		result := make(map[string]bool)
		for k, v := range sec {
			result[k] = v == "1" || v == "true" || v == "TRUE"
		}
		return result
	}
	return make(map[string]bool)
}
