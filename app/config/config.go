package config

import (
	"github.com/yanosea/gct/pkg/errors"
	"github.com/yanosea/gct/pkg/proxy"
)

type Configurator interface {
	GetConfig() (*TodoConfig, error)
}

type configurator struct {
	envconfig proxy.Envconfig
}

func NewConfigurator(
	ep proxy.Envconfig,
) Configurator {
	return &configurator{
		envconfig: ep,
	}
}

type TodoConfig struct {
	DBDirPath    string `envconfig:"GCT_DB_DIR_PATH" default:"XDG_DATA_HOME/gct"`
	OutputFormat string `envconfig:"GCT_OUTPUT_FORMAT" default:"text"`
}

func (c *configurator) GetConfig() (*TodoConfig, error) {
	var config TodoConfig
	var err error
	if err = c.envconfig.Process("", &config); err != nil {
		return nil, errors.Wrap(err, "failed to process environment")
	}
	return &config, err
}
