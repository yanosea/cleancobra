package config

import (
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
)

// Config represents the application configuration
type Config struct {
	// DataFile is the path to the JSON file where todos are stored
	DataFile string `envconfig:"GCT_DATA_FILE"`
}

// Load loads the configuration from environment variables and applies defaults
func Load() (*Config, error) {
	// Create real proxy implementations
	osProxy := proxy.NewOS()
	filepathProxy := proxy.NewFilepath()
	envconfigProxy := proxy.NewEnvconfig()

	// Call LoadWithDependencies with real proxy instances
	return LoadWithDependencies(osProxy, filepathProxy, envconfigProxy)
}

// LoadWithDependencies loads the configuration with injected dependencies for testing
func LoadWithDependencies(osProxy proxy.OS, filepathProxy proxy.Filepath, envconfigProxy proxy.Envconfig) (*Config, error) {
	var cfg Config

	// Load configuration from environment variables
	if err := envconfigProxy.Process("", &cfg); err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeConfiguration,
			"failed to load configuration from environment",
			err,
		)
	}

	// Apply default data file path if not set
	if cfg.DataFile == "" {
		defaultPath, err := getDefaultDataFilePath(osProxy, filepathProxy)
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
	if err := cfg.ValidateWithDeps(osProxy, filepathProxy); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Create real proxy implementations for backward compatibility
	osProxy := proxy.NewOS()
	filepathProxy := proxy.NewFilepath()

	// Use ValidateWithDeps with real proxy instances
	return c.ValidateWithDeps(osProxy, filepathProxy)
}

// ValidateWithDeps validates the configuration values with injected dependencies
func (c *Config) ValidateWithDeps(osProxy proxy.OS, filepathProxy proxy.Filepath) error {
	if c.DataFile == "" {
		return domain.NewDomainError(
			domain.ErrorTypeInvalidInput,
			"data file path cannot be empty",
			nil,
		)
	}

	// Check if the directory exists or can be created
	dir := filepathProxy.Dir(c.DataFile)
	if err := ensureDirectoryExistsWithDeps(dir, osProxy); err != nil {
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
	// Try XDG_DATA_HOME first
	if xdgDataHome := osProxy.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepathProxy.Join(xdgDataHome, "gct", "todos.json"), nil
	}

	// Fallback to ~/.local/share/gct/todos.json
	homeDir, err := osProxy.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepathProxy.Join(homeDir, ".local", "share", "gct", "todos.json"), nil
}

// ensureDirectoryExistsWithDeps creates the directory if it doesn't exist using injected dependencies
func ensureDirectoryExistsWithDeps(dir string, osProxy proxy.OS) error {
	if _, err := osProxy.Stat(dir); osProxy.IsNotExist(err) {
		if err := osProxy.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
