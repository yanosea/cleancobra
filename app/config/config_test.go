package config

import (
	"os"
	"testing"

	"github.com/yanosea/gct/pkg/proxy"
)

func TestLoad(t *testing.T) {
	// Test LoadWithDependencies with real implementations
	osProxy := proxy.NewOS()
	filepathProxy := proxy.NewFilepath()
	envconfigProxy := proxy.NewEnvconfig()

	// Test Load function with real implementations
	config, err := Load(osProxy, filepathProxy, envconfigProxy)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if config == nil {
		t.Fatal("Load() returned nil config")
	}

	// Verify that DataFile is set (either from env or default)
	if config.DataFile == "" {
		t.Error("Load() returned config with empty DataFile")
	}
}

func TestLoadWithDependencies(t *testing.T) {
	// Test LoadWithDependencies with real implementations
	osProxy := proxy.NewOS()
	filepathProxy := proxy.NewFilepath()
	envconfigProxy := proxy.NewEnvconfig()

	config, err := Load(osProxy, filepathProxy, envconfigProxy)
	if err != nil {
		t.Fatalf("LoadWithDependencies() failed: %v", err)
	}

	if config == nil {
		t.Fatal("LoadWithDependencies() returned nil config")
	}

	// Verify that DataFile is set (either from env or default)
	if config.DataFile == "" {
		t.Error("LoadWithDependencies() returned config with empty DataFile")
	}
}

func TestLoadWithDependencies_WithEnvVar(t *testing.T) {
	// Set environment variable for test
	testDataFile := "/tmp/test_todos.json"
	os.Setenv("GCT_DATA_FILE", testDataFile)
	defer os.Unsetenv("GCT_DATA_FILE")

	// Test LoadWithDependencies with real implementations
	osProxy := proxy.NewOS()
	filepathProxy := proxy.NewFilepath()
	envconfigProxy := proxy.NewEnvconfig()

	config, err := Load(osProxy, filepathProxy, envconfigProxy)
	if err != nil {
		t.Fatalf("LoadWithDependencies() failed: %v", err)
	}

	if config.DataFile != testDataFile {
		t.Errorf("LoadWithDependencies() DataFile = %v, want %v", config.DataFile, testDataFile)
	}
}
