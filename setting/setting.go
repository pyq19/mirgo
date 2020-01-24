package setting

import "os"

import "github.com/yenkeia/mirgo/common"

var (
	Conf      config
	BaseStats map[common.MirClass]baseStats
)

func init() {
	gopath := os.Getenv("GOPATH")
	Conf = config{
		Addr:          "0.0.0.0:7000",
		DBPath:        gopath + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite",
		MapDirPath:    gopath + "/src/github.com/yenkeia/mirgo/dotnettools/database/Maps/",
		ScriptDirPath: gopath + "/src/github.com/yenkeia/mirgo/script/",
	}
	// TODO
	BaseStats[common.MirClassWarrior] = baseStats{}
	BaseStats[common.MirClassWizard] = baseStats{}
	BaseStats[common.MirClassTaoist] = baseStats{}
}

type config struct {
	Addr          string
	DBPath        string
	MapDirPath    string
	ScriptDirPath string
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
