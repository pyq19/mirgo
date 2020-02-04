package setting

import (
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/setting/configure"
)

var (
	Conf      config
	BaseStats map[common.MirClass]baseStats
)

func init() {
	gameConfig := configure.GetGameConfigure("server")
	Conf = config{
		Addr:         gameConfig.Addr,
		DBPath:       gameConfig.DBPath,
		MapDirPath:   gameConfig.MapDirPath,
		ScriptDirPath: gameConfig.ScriptDirPath,
		DropDirPath:   gameConfig.DropDirPath,
		NPCDirPath:    gameConfig.NPCDirPath,
	}
	BaseStats = make(map[common.MirClass]baseStats)
	BaseStats[common.MirClassWarrior] = baseStats{
		HpGain:              4,
		HpGainRate:          4.5,
		MpGainRate:          0,
		BagWeightGain:       3,
		WearWeightGain:      20,
		HandWeightGain:      13,
		MinAc:               0,
		MaxAc:               7,
		MinMac:              0,
		MaxMac:              0,
		MinDc:               5,
		MaxDc:               5,
		MinMc:               0,
		MaxMc:               0,
		MinSc:               0,
		MaxSc:               0,
		StartAgility:        15,
		StartAccuracy:       5,
		StartCriticalRate:   0,
		StartCriticalDamage: 0,
		CritialRateGain:     0,
		CriticalDamageGain:  0,
	}
	BaseStats[common.MirClassWizard] = baseStats{
		HpGain:              15,
		HpGainRate:          1.8,
		MpGainRate:          0,
		BagWeightGain:       5,
		WearWeightGain:      100,
		HandWeightGain:      90,
		MinAc:               0,
		MaxAc:               0,
		MinMac:              0,
		MaxMac:              0,
		MinDc:               7,
		MaxDc:               7,
		MinMc:               7,
		MaxMc:               7,
		MinSc:               0,
		MaxSc:               0,
		StartAgility:        15,
		StartAccuracy:       5,
		StartCriticalRate:   0,
		StartCriticalDamage: 0,
		CritialRateGain:     0,
		CriticalDamageGain:  0,
	}
	BaseStats[common.MirClassTaoist] = baseStats{
		HpGain:              6,
		HpGainRate:          2.5,
		MpGainRate:          0,
		BagWeightGain:       4,
		WearWeightGain:      50,
		HandWeightGain:      42,
		MinAc:               0,
		MaxAc:               0,
		MinMac:              12,
		MaxMac:              6,
		MinDc:               7,
		MaxDc:               7,
		MinMc:               0,
		MaxMc:               0,
		MinSc:               7,
		MaxSc:               7,
		StartAgility:        18,
		StartAccuracy:       5,
		StartCriticalRate:   0,
		StartCriticalDamage: 0,
		CritialRateGain:     0,
		CriticalDamageGain:  0,
	}
}

type config struct {
	Addr          string
	DBPath        string
	MapDirPath    string
	ScriptDirPath string
	DropDirPath   string
	NPCDirPath    string
}

type baseStats struct {
	HpGain              float32
	HpGainRate          float32
	MpGainRate          float32
	BagWeightGain       float32
	WearWeightGain      float32
	HandWeightGain      float32
	MinAc               int
	MaxAc               int
	MinMac              int
	MaxMac              int
	MinDc               int
	MaxDc               int
	MinMc               int
	MaxMc               int
	MinSc               int
	MaxSc               int
	StartAgility        int
	StartAccuracy       int
	StartCriticalRate   int
	StartCriticalDamage int
	CritialRateGain     float32
	CriticalDamageGain  float32
}
