package program

import (
	"os"
	"testing"
)

func TestInitializeProgram(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		cleanup func()
		wantErr bool
	}{
		{
			name: "positive testing",
			setup: func() {
				// Set up a temporary data file for testing
				tempFile := "/tmp/gct_test_program_init.json"
				os.Setenv("GCT_DATA_FILE", tempFile)
				os.Remove(tempFile)
			},
			cleanup: func() {
				if dataFile := os.Getenv("GCT_DATA_FILE"); dataFile != "" {
					os.Remove(dataFile)
				}
				os.Unsetenv("GCT_DATA_FILE")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			got, err := InitializeProgram()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("InitializeProgram() returned nil program")
					return
				}
			}
		})
	}
}