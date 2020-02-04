package setting

import (
	"github.com/joho/godotenv"
	"os"
	"io/ioutil"
	"fmt"
	"github.com/yenkeia/mirgo/setting/env"
	"strings"
	"github.com/yenkeia/mirgo/setting/configure"
	"errors"
	"gopkg.in/yaml.v2"
)


type AppConfig struct {
	Env     string
	IsDebug bool
}

var appConfig AppConfig

var configFile map[string][]byte

func init() {
	// Init config container
	configFile = make(map[string][]byte)
	envLoad()
	configLoad()
	initAppConfig()
	gameLoad()
}


// Must load before loading config
func envLoad() {
	err := godotenv.Load("./.env")
	if err != nil {
		panic("error loading .env file")
	}
}

func initAppConfig() {
	debug := os.Getenv("APP_DEBUG")

	appConfig = AppConfig{
		Env:     os.Getenv("APP_ENV"),
		IsDebug: debug == "true",
	}
}

func configLoad() {
	configPath := "./config"
	pathSep := string(os.PathSeparator)
	_, err := os.Stat(configPath)
	if err != nil && os.IsNotExist(err) {
		panic("config path not exists")
	}
	files, _ := ioutil.ReadDir(configPath)
	for _, f := range files {
		path := configPath + pathSep + f.Name()
		data, fileErr := ioutil.ReadFile(path)
		if fileErr != nil {
			panic(fmt.Sprintf("can not open config %s", f.Name()))
		}

		// Replace config env variable
		data = env.ReplaceConfigEnv(data)
		fileName := strings.Split(f.Name(), ".")[0]

		// Attach data to global config variable
		configFile[fileName] = data
	}
}

func Config(name string, object interface{}) error {
	config, isExists := configFile[name]
	if !isExists {
		return errors.New(fmt.Sprintf("can not found config %s group", name))
	}

	ymlErr := yaml.Unmarshal(config, object)
	if ymlErr != nil {
		panic(fmt.Sprintf("config %s formation is not correct", name))
	}
	return nil
}

func gameLoad() {
	database := configure.GameConfigure{}
	err := Config("game", &database)

	if err != nil {
		panic(err)
	}
	configure.InitGameConfigure(database, appConfig.IsDebug)
}