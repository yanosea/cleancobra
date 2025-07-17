package view

import (
	"testing"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/mock/gomock"
)

func createTestStateModelForView(t *testing.T) (*model.StateModel, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	
	// Create mock use cases
	mockAddUseCase := &application.AddTodoUseCase{}
	mockListUseCase := &application.ListTodoUseCase{}
	mockToggleUseCase := &application.ToggleTodoUseCase{}
	mockDeleteUseCase := &application.DeleteTodoUseCase{}
	
	// Create mock bubbles
	mockBubbles := proxy.NewMockBubbles(ctrl)
	mockTextInput := proxy.NewMockTextInput(ctrl)
	
	mockBubbles.EXPECT().NewTextInput().Return(mockTextInput)
	mockTextInput.EXPECT().SetPlaceholder("Enter todo description...")
	mockTextInput.EXPECT().SetCharLimit(500)
	mockTextInput.EXPECT().SetWidth(50)
	
	// Allow common text input operations that might be called during mode changes
	mockTextInput.EXPECT().Blur().AnyTimes()
	mockTextInput.EXPECT().SetValue(gomock.Any()).AnyTimes()
	mockTextInput.EXPECT().Focus().AnyTimes()
	
	stateModel := model.NewStateModel(
		mockAddUseCase,
		mockListUseCase,
		mockToggleUseCase,
		mockDeleteUseCase,
		mockBubbles,
	)
	
	return stateModel, ctrl
}

func TestNewFooterView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			footerView := NewFooterView(mockLipgloss)

			if footerView == nil {
				t.Error("NewFooterView() returned nil")
			}
		})
	}
}

func TestFooterView_Render(t *testing.T) {
	tests := []struct {
		name     string
		mode     model.Mode
		width    int
		expected string
	}{
		{
			name:     "positive testing (normal mode)",
			mode:     model.ModeNormal,
			width:    80,
			expected: "rendered footer",
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			width:    80,
			expected: "rendered footer",
		},
		{
			name:     "positive testing (edit mode)",
			mode:     model.ModeEdit,
			width:    80,
			expected: "rendered footer",
		},
		{
			name:     "positive testing (confirmation mode)",
			mode:     model.ModeConfirmation,
			width:    80,
			expected: "rendered footer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			// Setup for Render
			mockStyle.EXPECT().Width(tt.width).Return(mockStyle)
			mockStyle.EXPECT().Render(gomock.Any()).Return(tt.expected)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetMode(tt.mode)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.Render(stateModel, tt.width)

			if result != tt.expected {
				t.Errorf("Render() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_RenderCompact(t *testing.T) {
	tests := []struct {
		name     string
		mode     model.Mode
		width    int
		expected string
	}{
		{
			name:     "positive testing (normal mode)",
			mode:     model.ModeNormal,
			width:    40,
			expected: "compact footer",
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			width:    40,
			expected: "compact footer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			// Setup for RenderCompact
			mockStyle.EXPECT().Width(tt.width).Return(mockStyle)
			mockStyle.EXPECT().Render(gomock.Any()).Return(tt.expected)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetMode(tt.mode)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.RenderCompact(stateModel, tt.width)

			if result != tt.expected {
				t.Errorf("RenderCompact() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_RenderWithScrollIndicator(t *testing.T) {
	tests := []struct {
		name            string
		mode            model.Mode
		width           int
		scrollIndicator string
		expected        string
	}{
		{
			name:            "positive testing",
			mode:            model.ModeNormal,
			width:           80,
			scrollIndicator: "scroll indicator",
			expected:        "footer with scroll",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			// Setup for RenderWithScrollIndicator
			mockLipgloss.EXPECT().Left().Return(lipgloss.Left)
			mockLipgloss.EXPECT().JoinVertical(lipgloss.Left, tt.scrollIndicator, gomock.Any()).Return("joined content")
			mockStyle.EXPECT().Width(tt.width).Return(mockStyle)
			mockStyle.EXPECT().Render("joined content").Return(tt.expected)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetMode(tt.mode)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.RenderWithScrollIndicator(stateModel, tt.width, tt.scrollIndicator)

			if result != tt.expected {
				t.Errorf("RenderWithScrollIndicator() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_getHelpText(t *testing.T) {
	tests := []struct {
		name     string
		mode     model.Mode
		expected string
	}{
		{
			name:     "positive testing (normal mode)",
			mode:     model.ModeNormal,
			expected: "↑/k: Up • ↓/j: Down • g: Top • G: Bottom • Space: Toggle • a: Add • e: Edit • d: Delete • r: Refresh • q: Quit",
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			expected: "Enter: Add • Esc: Cancel",
		},
		{
			name:     "positive testing (edit mode)",
			mode:     model.ModeEdit,
			expected: "Enter: Save • Esc: Cancel",
		},
		{
			name:     "positive testing (confirmation mode)",
			mode:     model.ModeConfirmation,
			expected: "y: Yes • n/Esc: No",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.getHelpText(tt.mode)

			if result != tt.expected {
				t.Errorf("getHelpText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_getCompactHelpText(t *testing.T) {
	tests := []struct {
		name     string
		mode     model.Mode
		expected string
	}{
		{
			name:     "positive testing (normal mode)",
			mode:     model.ModeNormal,
			expected: "↑↓:Move Space:Toggle a:Add e:Edit d:Del q:Quit",
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			expected: "Enter:Add Esc:Cancel",
		},
		{
			name:     "positive testing (edit mode)",
			mode:     model.ModeEdit,
			expected: "Enter:Save Esc:Cancel",
		},
		{
			name:     "positive testing (confirmation mode)",
			mode:     model.ModeConfirmation,
			expected: "y:Yes n:No",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.getCompactHelpText(tt.mode)

			if result != tt.expected {
				t.Errorf("getCompactHelpText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_RenderModeSpecificHelp(t *testing.T) {
	tests := []struct {
		name     string
		mode     model.Mode
		width    int
		expected string
	}{
		{
			name:     "positive testing (normal mode)",
			mode:     model.ModeNormal,
			width:    80,
			expected: "mode specific help",
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			width:    80,
			expected: "mode specific help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			// Setup for RenderModeSpecificHelp
			// Allow multiple render calls - some for helpStyle, final one for footerStyle
			gomock.InOrder(
				mockStyle.EXPECT().Render(gomock.Any()).Return("styled line").MinTimes(1),
				mockStyle.EXPECT().Width(tt.width).Return(mockStyle),
				mockStyle.EXPECT().Render(gomock.Any()).Return(tt.expected),
			)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetMode(tt.mode)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.RenderModeSpecificHelp(stateModel, tt.width)

			if result != tt.expected {
				t.Errorf("RenderModeSpecificHelp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_GetHeight(t *testing.T) {
	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "positive testing",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.GetHeight()

			if result != tt.expected {
				t.Errorf("GetHeight() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_GetExpandedHeight(t *testing.T) {
	tests := []struct {
		name     string
		mode     model.Mode
		expected int
	}{
		{
			name:     "positive testing (normal mode)",
			mode:     model.ModeNormal,
			expected: 3,
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			expected: 2,
		},
		{
			name:     "positive testing (edit mode)",
			mode:     model.ModeEdit,
			expected: 2,
		},
		{
			name:     "positive testing (confirmation mode)",
			mode:     model.ModeConfirmation,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			footerView := NewFooterView(mockLipgloss)
			result := footerView.GetExpandedHeight(tt.mode)

			if result != tt.expected {
				t.Errorf("GetExpandedHeight() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFooterView_SetStyle(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)
			newStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			footerView := NewFooterView(mockLipgloss)
			footerView.SetStyle(newStyle)

			// Test that the style was set (we can't directly verify this without exposing the field)
			// This test mainly ensures the method doesn't panic
		})
	}
}

func TestFooterView_SetHelpStyle(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)
			newStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewFooterView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			footerView := NewFooterView(mockLipgloss)
			footerView.SetHelpStyle(newStyle)

			// Test that the style was set (we can't directly verify this without exposing the field)
			// This test mainly ensures the method doesn't panic
		})
	}
}