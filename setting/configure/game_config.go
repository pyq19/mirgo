package configure

var GameConfigData map[string]GameConfig

type GameConfig struct {
	Addr          string `yaml:"addr"`
	DBPath        string `yaml:"db_path"`
	MapDirPath    string `yaml:"map_dir_path"`
	ScriptDirPath string `yaml:"script_dir_path"`
	DropDirPath   string `yaml:"drop_dir_path"`
	NPCDirPath    string `yaml:"npc_dir_path"`
}


type GameConfigure struct {
	Game map[string]GameConfig `yaml:"game"`
}

func InitGameConfigure(databases GameConfigure, isDebug bool) {

	GameConfigData = make(map[string]GameConfig)
	for key, config := range databases.Game {
		GameConfigData[key] = config
	}
}

func GetGameConfigure(server string)  GameConfig{
	if _, isExists := GameConfigData[server]; isExists == true {
		return GameConfigData[server]
	}
	panic("game config not found2")
}

