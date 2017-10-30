package conf

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/logs"
)

// MainConfig config from toml
type MainConfig struct {
	TulingURL  string
	TulingKeys []string

	DBMaxOpenConns int
	DBMaxIdleConns int
	UserDB         string
}

// Config global config
var Config *MainConfig

func init() {
	Config = &MainConfig{}
	_, thisFilePath, _, _ := runtime.Caller(0)
	file := filepath.Join(filepath.Dir(thisFilePath), "config.toml")
	if _, err := toml.DecodeFile(file, Config); err != nil {
		logs.Error("Parse config.toml Failed")
		os.Exit(1)
	}
	logs.Info("config: %+v", Config)
}
