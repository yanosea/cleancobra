package proxy

import (
	"github.com/kelseyhightower/envconfig"
)

type Envconfig interface {
	Process(prefix string, spec interface{}) error
}

type envconfigProxy struct{}

func NewEnvconfig() Envconfig {
	return &envconfigProxy{}
}

func (envconfigProxy) Process(prefix string, spec interface{}) error {
	return envconfig.Process(prefix, spec)
}
