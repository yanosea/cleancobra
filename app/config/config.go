package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"

	"github.com/yanosea/cleancobra-pkg/errors"
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
	err    error
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		var processErr error
		if processErr = envconfig.Process("", &config); processErr != nil {
			err = errors.Wrap(processErr, "failed to process environment")
		}
	})
	return &config, err
}
