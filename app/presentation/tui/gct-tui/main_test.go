package main

import (
	"os"
	"testing"

	"github.com/yanosea/gct/app/presentation/tui/gct-tui/program"
)

func TestMain_Integration(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() error
		wantErr bool
	}{
		{
			name: "positive testing - successful TUI initialization",
			setup: func() error {
				// Set up a temporary data file for testing
				tempFile := "/tmp/gct_test_tui.json"
				os.Setenv("GCT_DATA_FILE", tempFile)
				
				// Clean up any existing test file
				os.Remove(tempFile)
				
				return nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if err := tt.setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Test program initialization (similar to main function logic)
			prog, err := program.InitializeProgram()
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Program initialization failed: %v", err)
			}

			if prog == nil {
				t.Error("Program should not be nil")
			}

			// Cleanup
			if dataFile := os.Getenv("GCT_DATA_FILE"); dataFile != "" {
				os.Remove(dataFile)
			}
			os.Unsetenv("GCT_DATA_FILE")
		})
	}
}

