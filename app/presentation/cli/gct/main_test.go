package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestMainFunction(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "positive testing - help command",
			args:           []string{"--help"},
			expectError:    false,
			expectedOutput: "gct is a todo application built with clean architecture principles",
		},
		{
			name:           "positive testing - add command help",
			args:           []string{"add", "--help"},
			expectError:    false,
			expectedOutput: "Add a new todo",
		},
		{
			name:           "negative testing - invalid command failed",
			args:           []string{"invalid"},
			expectError:    true,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the CLI application
			execName := "gct-main-test.exe"
			buildCmd := exec.Command("go", "build", "-o", execName, ".")
			if err := buildCmd.Run(); err != nil {
				t.Fatalf("Failed to build CLI application: %v", err)
			}
			defer os.Remove(execName)

			// Execute the CLI application with test arguments
			absPath, _ := filepath.Abs(execName)
			cmd := exec.Command(absPath, tt.args...)
			output, err := cmd.CombinedOutput()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, string(output))
			}

			// Check output expectation
			if tt.expectedOutput != "" && !strings.Contains(string(output), tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got: %s", tt.expectedOutput, string(output))
			}
		})
	}
}

func TestMainInitialization(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		expectError bool
	}{
		{
			name: "positive testing - successful initialization",
			setupEnv: func() {
				// No special setup needed for default configuration
			},
			cleanupEnv: func() {
				// No cleanup needed
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			tt.setupEnv()
			defer tt.cleanupEnv()

			// Build the CLI application
			execName := "gct-init-main-test.exe"
			buildCmd := exec.Command("go", "build", "-o", execName, ".")
			if err := buildCmd.Run(); err != nil {
				t.Fatalf("Failed to build CLI application: %v", err)
			}
			defer os.Remove(execName)

			// Execute the CLI application with --help to test initialization
			absPath, _ := filepath.Abs(execName)
			cmd := exec.Command(absPath, "--help")
			_, err := cmd.CombinedOutput()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error during initialization: %v", err)
			}
		})
	}
}