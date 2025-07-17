package view

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/mock/gomock"
)

func TestNewHeaderView(t *testing.T) {
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

			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			headerView := NewHeaderView(mockLipgloss)

			if headerView == nil {
				t.Error("NewHeaderView() returned nil")
			}
		})
	}
}

func TestHeaderView_Render(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected string
	}{
		{
			name:     "positive testing",
			width:    80,
			expected: "rendered header",
		},
		{
			name:     "positive testing (narrow width)",
			width:    40,
			expected: "rendered header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewHeaderView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			// Setup for Render
			mockLipgloss.EXPECT().Left().Return(lipgloss.Left)
			mockLipgloss.EXPECT().JoinHorizontal(lipgloss.Left, gomock.Any(), gomock.Any(), gomock.Any()).Return("joined header")
			mockStyle.EXPECT().Width(tt.width).Return(mockStyle)
			mockStyle.EXPECT().Render("joined header").Return(tt.expected)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()

			headerView := NewHeaderView(mockLipgloss)
			result := headerView.Render(stateModel, tt.width)

			if result != tt.expected {
				t.Errorf("Render() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHeaderView_RenderCompact(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected string
	}{
		{
			name:     "positive testing",
			width:    40,
			expected: "compact header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewHeaderView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			// Setup for RenderCompact
			mockLipgloss.EXPECT().Left().Return(lipgloss.Left)
			mockLipgloss.EXPECT().JoinHorizontal(lipgloss.Left, gomock.Any(), gomock.Any(), gomock.Any()).Return("compact joined header")
			mockStyle.EXPECT().Width(tt.width).Return(mockStyle)
			mockStyle.EXPECT().Render("compact joined header").Return(tt.expected)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()

			headerView := NewHeaderView(mockLipgloss)
			result := headerView.RenderCompact(stateModel, tt.width)

			if result != tt.expected {
				t.Errorf("RenderCompact() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHeaderView_RenderWithMode(t *testing.T) {
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
			expected: "header with mode",
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			width:    80,
			expected: "header with mode",
		},
		{
			name:     "positive testing (edit mode)",
			mode:     model.ModeEdit,
			width:    80,
			expected: "header with mode",
		},
		{
			name:     "positive testing (confirmation mode)",
			mode:     model.ModeConfirmation,
			width:    80,
			expected: "header with mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewHeaderView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			// Setup for RenderWithMode
			mockLipgloss.EXPECT().Left().Return(lipgloss.Left)
			mockLipgloss.EXPECT().JoinHorizontal(lipgloss.Left, gomock.Any(), gomock.Any(), gomock.Any()).Return("mode header")
			mockStyle.EXPECT().Width(tt.width).Return(mockStyle)
			mockStyle.EXPECT().Render("mode header").Return(tt.expected)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetMode(tt.mode)

			headerView := NewHeaderView(mockLipgloss)
			result := headerView.RenderWithMode(stateModel, tt.width)

			if result != tt.expected {
				t.Errorf("RenderWithMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHeaderView_GetCompletedCount(t *testing.T) {
	tests := []struct {
		name     string
		todos    []*model.ItemModel
		expected int
	}{
		{
			name:     "positive testing (no todos)",
			todos:    []*model.ItemModel{},
			expected: 0,
		},
		{
			name: "positive testing (no completed todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Todo 1",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
				model.NewItemModel(&domain.Todo{
					ID:          2,
					Description: "Todo 2",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			expected: 0,
		},
		{
			name: "positive testing (some completed todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Todo 1",
					Done:        true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
				model.NewItemModel(&domain.Todo{
					ID:          2,
					Description: "Todo 2",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
				model.NewItemModel(&domain.Todo{
					ID:          3,
					Description: "Todo 3",
					Done:        true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			expected: 2,
		},
		{
			name: "positive testing (all completed todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Todo 1",
					Done:        true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
				model.NewItemModel(&domain.Todo{
					ID:          2,
					Description: "Todo 2",
					Done:        true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewHeaderView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			headerView := NewHeaderView(mockLipgloss)
			result := headerView.GetCompletedCount(tt.todos)

			if result != tt.expected {
				t.Errorf("GetCompletedCount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHeaderView_getModeIndicator(t *testing.T) {
	tests := []struct {
		name     string
		mode     model.Mode
		expected string
	}{
		{
			name:     "positive testing (normal mode)",
			mode:     model.ModeNormal,
			expected: "NORMAL",
		},
		{
			name:     "positive testing (input mode)",
			mode:     model.ModeInput,
			expected: "ADD",
		},
		{
			name:     "positive testing (edit mode)",
			mode:     model.ModeEdit,
			expected: "EDIT",
		},
		{
			name:     "positive testing (confirmation mode)",
			mode:     model.ModeConfirmation,
			expected: "CONFIRM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewHeaderView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			headerView := NewHeaderView(mockLipgloss)
			result := headerView.getModeIndicator(tt.mode)

			if result != tt.expected {
				t.Errorf("getModeIndicator() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHeaderView_GetHeight(t *testing.T) {
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

			// Setup for NewHeaderView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			headerView := NewHeaderView(mockLipgloss)
			result := headerView.GetHeight()

			if result != tt.expected {
				t.Errorf("GetHeight() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHeaderView_SetStyle(t *testing.T) {
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

			// Setup for NewHeaderView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
			mockStyle.EXPECT().Bold(true).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle)

			headerView := NewHeaderView(mockLipgloss)
			headerView.SetStyle(newStyle)

			// Test that the style was set (we can't directly verify this without exposing the field)
			// This test mainly ensures the method doesn't panic
		})
	}
}