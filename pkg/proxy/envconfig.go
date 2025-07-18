//go:generate mockgen -source=envconfig.go -destination=envconfig_mock.go -package=proxy

package proxy

import (
	"github.com/kelseyhightower/envconfig"
)

// Envconfig provides a proxy interface for envconfig package functionality
type Envconfig interface {
	Process(prefix string, spec interface{}) error
}

// EnvconfigImpl implements the Envconfig interface using the envconfig package
type EnvconfigImpl struct{}

// NewEnvconfig creates a new Envconfig implementation
func NewEnvconfig() Envconfig {
	return &EnvconfigImpl{}
}

func (e *EnvconfigImpl) Process(prefix string, spec interface{}) error {
	return envconfig.Process(prefix, spec)
}
