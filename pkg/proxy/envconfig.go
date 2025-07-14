package proxy

import (
	"github.com/kelseyhightower/envconfig"
)

type Envconfig interface {
	Process(prefix string, spec any) error
}

type envconfigProxy struct{}

func NewEnvconfig() Envconfig {
	return &envconfigProxy{}
}

func (envconfigProxy) Process(prefix string, spec any) error {
	return envconfig.Process(prefix, spec)
}
