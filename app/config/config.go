package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Todo TodoConfig
}

type TodoConfig struct {
	DBDirPath    string `envconfig:"CLEANCOBRA_TODO_DB_DIR_PATH" default:"XDG_DATA_HOME/cleancobra/todos"`
	OutputFormat string `envconfig:"CLEANCOBRA_TODO_OUTPUT_FORMAT" default:"text"`
}

var (
	once   sync.Once
	config Config
)

func GetConfig() *Config {
	once.Do(func() {
		if err := envconfig.Process("", &config); err != nil {
			panic(err)
		}
	})
	return &config
}
