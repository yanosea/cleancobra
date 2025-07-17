package completion

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewBashCompletionCommand(t *testing.T) {
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
			rootCmd := cobraProxy.NewCommand()

			result := NewBashCompletionCommand(cobraProxy, rootCmd)

			if result == nil {
				t.Error("Expected command to be created, got nil")
			}
		})
	}
}

func TestRunBashCompletion(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*proxy.MockCommand)
		wantErr   bool
	}{
		{
			name: "positive testing",
			setupMock: func(mockCmd *proxy.MockCommand) {
				mockCmd.EXPECT().GenBashCompletion(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "negative testing (generation failed)",
			setupMock: func(mockCmd *proxy.MockCommand) {
				mockCmd.EXPECT().GenBashCompletion(gomock.Any()).Return(errors.New("generation failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRootCmd := proxy.NewMockCommand(ctrl)
			tt.setupMock(mockRootCmd)

			err := runBashCompletion(mockRootCmd)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}