package setting

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/util"
)

type Conf struct {
	DataPath string
	Acceptor string // websocket | tcp(defulat)
}

type Settings struct {
	Addr          string
	DBPath        string
	AccountDBPath string
	MapDirPath    string
	DropDirPath   string
	EnvirPath     string
	ConfigsPath   string
	RoutePath     string
	Acceptor      string

	BaseStats           map[cm.MirClass]baseStats
	MagicResistWeight   int
	Guild_RequiredLevel int // 创建行会需要的等级
}

func Must() *Settings {
	s, err := New()
	if err != nil {
		panic("配置初始化失败:" + err.Error())
	}
	return s
}

func New() (*Settings, error) {
	file := "./config.toml"

	conf := Conf{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		checkdir, err := filepath.Abs("../dotnettools")
		if err != nil {
			return nil, errors.New("没有配置")
		}

		if util.IsDir(checkdir) {
			conf.DataPath = checkdir
		}
	}

	err = toml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	BaseStats := make(map[cm.MirClass]baseStats)
	BaseStats[cm.MirClassWarrior] = baseStats{
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
	BaseStats[cm.MirClassWizard] = baseStats{
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
	BaseStats[cm.MirClassTaoist] = baseStats{
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

	return &Settings{
		Addr:              "0.0.0.0:7000",
		DBPath:            filepath.Join(conf.DataPath, "/mir.sqlite"),
		AccountDBPath:     filepath.Join(conf.DataPath, "/account.sqlite"),
		MapDirPath:        filepath.Join(conf.DataPath, "/Maps/"),
		DropDirPath:       filepath.Join(conf.DataPath, "/Envir/Drops/"),
		EnvirPath:         filepath.Join(conf.DataPath, "/Envir/"),
		ConfigsPath:       filepath.Join(conf.DataPath, "/Configs/"),
		RoutePath:         filepath.Join(conf.DataPath, "/Envir/Routes/"),
		BaseStats:         BaseStats,
		Acceptor:          conf.Acceptor,
		MagicResistWeight: 10,
	}, nil
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
