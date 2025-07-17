package completion

import (
	"testing"

	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewCompletionCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create real implementation without mocks for positive testing
			cobraProxy := proxy.NewCobra()

			result := NewCompletionCommand(cobraProxy)

			if result == nil {
				t.Error("Expected command to be created, got nil")
			}
		})
	}
}

func TestRunCompletion(t *testing.T) {
	tests := []struct {
		name    string
		shell   string
		wantErr bool
	}{
		{
			name:    "negative testing (bash shell failed)",
			shell:   "bash",
			wantErr: true,
		},
		{
			name:    "negative testing (zsh shell failed)",
			shell:   "zsh",
			wantErr: true,
		},
		{
			name:    "negative testing (fish shell failed)",
			shell:   "fish",
			wantErr: true,
		},
		{
			name:    "negative testing (powershell shell failed)",
			shell:   "powershell",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCompletion(tt.shell)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if err != nil && err.Error() != "completion command requires a subcommand: bash, zsh, fish, or powershell" {
					t.Errorf("Expected error message to contain 'completion command requires a subcommand', got %q", err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}