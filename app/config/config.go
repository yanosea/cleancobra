package config

import (
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/yanosea/gct/app/domain"
)

// Config represents the application configuration
type Config struct {
	// DataFile is the path to the JSON file where todos are stored
	DataFile string `envconfig:"GCT_DATA_FILE"`
}

// Load loads the configuration from environment variables and applies defaults
func Load() (*Config, error) {
	var cfg Config

	// Load configuration from environment variables
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeConfiguration,
			"failed to load configuration from environment",
			err,
		)
	}

	// Apply default data file path if not set
	if cfg.DataFile == "" {
		defaultPath, err := getDefaultDataFilePath()
		if err != nil {
			return nil, domain.NewDomainError(
				domain.ErrorTypeConfiguration,
				"failed to determine default data file path",
				err,
			)
		}
		cfg.DataFile = defaultPath
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	if c.DataFile == "" {
		return domain.NewDomainError(
			domain.ErrorTypeInvalidInput,
			"data file path cannot be empty",
			nil,
		)
	}

	// Check if the directory exists or can be created
	dir := filepath.Dir(c.DataFile)
	if err := ensureDirectoryExists(dir); err != nil {
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
func getDefaultDataFilePath() (string, error) {
	// Try XDG_DATA_HOME first
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "gct", "todos.json"), nil
	}

	// Fallback to ~/.local/share/gct/todos.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".local", "share", "gct", "todos.json"), nil
}

// ensureDirectoryExists creates the directory if it doesn't exist
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
