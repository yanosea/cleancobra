package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yanosea/gct/app/domain"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalDataFile := os.Getenv("GCT_DATA_FILE")
	originalXDGDataHome := os.Getenv("XDG_DATA_HOME")
	defer func() {
		os.Setenv("GCT_DATA_FILE", originalDataFile)
		os.Setenv("XDG_DATA_HOME", originalXDGDataHome)
	}()

	tests := []struct {
		name           string
		setupEnv       func()
		wantDataFile   string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "positive testing (custom data file from environment)",
			setupEnv: func() {
				// Create a temporary directory for testing
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, "custom_todos.json")
				os.Setenv("GCT_DATA_FILE", testFile)
				os.Unsetenv("XDG_DATA_HOME")
			},
			wantDataFile: "", // Will be set dynamically in test
			wantErr:      false,
		},
		{
			name: "positive testing (XDG_DATA_HOME path)",
			setupEnv: func() {
				tempDir := t.TempDir()
				os.Unsetenv("GCT_DATA_FILE")
				os.Setenv("XDG_DATA_HOME", tempDir)
			},
			wantDataFile: "", // Will be set dynamically in test
			wantErr:      false,
		},
		{
			name: "positive testing (fallback to home directory)",
			setupEnv: func() {
				os.Unsetenv("GCT_DATA_FILE")
				os.Unsetenv("XDG_DATA_HOME")
			},
			wantDataFile: "", // Will be set dynamically in test
			wantErr:      false,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.expectedErrMsg != "" && !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Load() error = %v, expected to contain %v", err, tt.expectedErrMsg)
				}
				return
			}

			if got == nil {
				t.Error("Load() returned nil config")
				return
			}

			// Verify data file path is set
			if got.DataFile == "" {
				t.Error("Load() DataFile is empty")
			}

			// For custom data file test, verify it matches what we set
			if tt.name == "positive testing (custom data file from environment)" {
				expectedPath := os.Getenv("GCT_DATA_FILE")
				if got.DataFile != expectedPath {
					t.Errorf("Load() DataFile = %v, want %v", got.DataFile, expectedPath)
				}
			}

			// For XDG_DATA_HOME test, verify it uses XDG path
			if tt.name == "positive testing (XDG_DATA_HOME path)" {
				xdgDataHome := os.Getenv("XDG_DATA_HOME")
				expectedPath := filepath.Join(xdgDataHome, "gct", "todos.json")
				if got.DataFile != expectedPath {
					t.Errorf("Load() DataFile = %v, want %v", got.DataFile, expectedPath)
				}
			}

			// For fallback test, verify it uses home directory
			if tt.name == "positive testing (fallback to home directory)" {
				homeDir, _ := os.UserHomeDir()
				expectedPath := filepath.Join(homeDir, ".local", "share", "gct", "todos.json")
				if got.DataFile != expectedPath {
					t.Errorf("Load() DataFile = %v, want %v", got.DataFile, expectedPath)
				}
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "positive testing (valid config)",
			config: &Config{
				DataFile: filepath.Join(t.TempDir(), "todos.json"),
			},
			wantErr: false,
		},
		{
			name: "negative testing (empty data file)",
			config: &Config{
				DataFile: "",
			},
			wantErr: true,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				// Verify it's a domain error
				if !domain.IsDomainError(err) {
					t.Errorf("Config.Validate() error should be a domain error, got %T", err)
				}
			}
		})
	}
}

func TestGetDefaultDataFilePath(t *testing.T) {
	// Save original environment
	originalXDGDataHome := os.Getenv("XDG_DATA_HOME")
	defer func() {
		os.Setenv("XDG_DATA_HOME", originalXDGDataHome)
	}()

	tests := []struct {
		name     string
		setupEnv func()
		wantPath func() string
		wantErr  bool
	}{
		{
			name: "positive testing (XDG_DATA_HOME set)",
			setupEnv: func() {
				os.Setenv("XDG_DATA_HOME", "/tmp/xdg_data")
			},
			wantPath: func() string {
				return filepath.Join("/tmp/xdg_data", "gct", "todos.json")
			},
			wantErr: false,
		},
		{
			name: "positive testing (XDG_DATA_HOME not set, fallback to home)",
			setupEnv: func() {
				os.Unsetenv("XDG_DATA_HOME")
			},
			wantPath: func() string {
				homeDir, _ := os.UserHomeDir()
				return filepath.Join(homeDir, ".local", "share", "gct", "todos.json")
			},
			wantErr: false,
		},
		{
			name: "positive testing (XDG_DATA_HOME empty, fallback to home)",
			setupEnv: func() {
				os.Setenv("XDG_DATA_HOME", "")
			},
			wantPath: func() string {
				homeDir, _ := os.UserHomeDir()
				return filepath.Join(homeDir, ".local", "share", "gct", "todos.json")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			got, err := getDefaultDataFilePath()
			if (err != nil) != tt.wantErr {
				t.Errorf("getDefaultDataFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				want := tt.wantPath()
				if got != want {
					t.Errorf("getDefaultDataFilePath() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	tests := []struct {
		name    string
		setupDir func() string
		wantErr bool
	}{
		{
			name: "positive testing (directory does not exist)",
			setupDir: func() string {
				tempDir := t.TempDir()
				return filepath.Join(tempDir, "new_directory")
			},
			wantErr: false,
		},
		{
			name: "positive testing (directory already exists)",
			setupDir: func() string {
				return t.TempDir() // Already exists
			},
			wantErr: false,
		},
		{
			name: "positive testing (nested directory creation)",
			setupDir: func() string {
				tempDir := t.TempDir()
				return filepath.Join(tempDir, "level1", "level2", "level3")
			},
			wantErr: false,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir()

			err := ensureDirectoryExists(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureDirectoryExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify directory was created
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					t.Errorf("ensureDirectoryExists() directory %v was not created", dir)
				}
			}
		})
	}
}