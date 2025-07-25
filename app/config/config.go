package config

import (
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
)

// Configurator provides access to the configuration
type Configurator struct {
	Envconfig proxy.Envconfig
	Filepath  proxy.Filepath
	OS        proxy.OS
}

// NewConfigurator creates a new Configurator
func NewConfigurator(
	envconfigProxy proxy.Envconfig,
	filepathProxy proxy.Filepath,
	osProxy proxy.OS,
) *Configurator {
	return &Configurator{
		Envconfig: envconfigProxy,
		Filepath:  filepathProxy,
		OS:        osProxy,
	}
}

// Config represents the application configuration
type Config struct {
	// DataFile is the path to the JSON file where todos are stored
	DataFile string `envconfig:"GCT_DATA_FILE"`
}

// Load loads the configuration
func (c *Configurator) Load() (*Config, error) {
	var cfg Config

	// load configuration from environment variables
	if err := c.Envconfig.Process("", &cfg); err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeConfiguration,
			"failed to load configuration from environment",
			err,
		)
	}

	// apply default data file path if not set
	if cfg.DataFile == "" {
		defaultPath, err := getDefaultDataFilePath(c.OS, c.Filepath)
		if err != nil {
			return nil, domain.NewDomainError(
				domain.ErrorTypeConfiguration,
				"failed to determine default data file path",
				err,
			)
		}
		cfg.DataFile = defaultPath
	}

	// validate configuration
	if err := cfg.Validate(c.OS, c.Filepath); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate validates the configuration values
func (c *Config) Validate(osProxy proxy.OS, filepathProxy proxy.Filepath) error {
	// check if the data file path is set
	if c.DataFile == "" {
		return domain.NewDomainError(
			domain.ErrorTypeInvalidInput,
			"data file path cannot be empty",
			nil,
		)
	}

	// check if the directory exists or can be created
	dir := filepathProxy.Dir(c.DataFile)
	if err := ensureDirectoryExists(dir, osProxy); err != nil {
		return domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to ensure data directory exists",
			err,
		)
	}

	return nil
}

// getDefaultDataFilePath returns the default path for the data file
// Following XDG Base Directory Specification with fallback to home directory
func getDefaultDataFilePath(osProxy proxy.OS, filepathProxy proxy.Filepath) (string, error) {
	// try XDG_DATA_HOME first
	if xdgDataHome := osProxy.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepathProxy.Join(xdgDataHome, "gct", "todos.json"), nil
	}

	// fallback to ~/.local/share/gct/todos.json
	homeDir, err := osProxy.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepathProxy.Join(homeDir, ".local", "share", "gct", "todos.json"), nil
}

// ensureDirectoryExists creates the directory if it doesn't exist
func ensureDirectoryExists(dir string, osProxy proxy.OS) error {
	if _, statErr := osProxy.Stat(dir); osProxy.IsNotExist(statErr) {
		if mkdirErr := osProxy.MkdirAll(dir, 0755); mkdirErr != nil {
			return mkdirErr
		}
	} else if statErr != nil {
		return statErr
	}

	return nil
}
