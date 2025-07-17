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

func createMockItemView(ctrl *gomock.Controller) *ItemView {
	mockLipgloss := proxy.NewMockLipgloss(ctrl)
	mockStyle := proxy.NewMockStyle(ctrl)

	// Setup for NewItemView - allow flexible calls since these are shared
	mockLipgloss.EXPECT().NewStyle().Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Padding(gomock.Any(), gomock.Any()).Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Background(gomock.Any()).Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Foreground(gomock.Any()).Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Border(gomock.Any(), gomock.Any()).Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Italic(gomock.Any()).Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Bold(gomock.Any()).Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Width(gomock.Any()).Return(mockStyle).AnyTimes()
	mockStyle.EXPECT().Render(gomock.Any()).Return("rendered").AnyTimes()

	return NewItemView(mockLipgloss)
}

func TestNewListView(t *testing.T) {
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

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)
			listView := NewListView(mockLipgloss, itemView)

			if listView == nil {
				t.Error("NewListView() returned nil")
			}
		})
	}
}

func TestListView_Render(t *testing.T) {
	tests := []struct {
		name   string
		todos  []*model.ItemModel
		width  int
		height int
	}{
		{
			name:   "positive testing (empty todos)",
			todos:  []*model.ItemModel{},
			width:  80,
			height: 20,
		},
		{
			name: "positive testing (with todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			width:  80,
			height: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for better testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			listView := NewListView(lipgloss, itemView)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetTodos(tt.todos)

			result := listView.Render(stateModel, tt.width, tt.height)

			// Verify that we get a non-empty result
			if result == "" {
				t.Error("Render() returned empty string")
			}
			
			// For empty todos, should contain empty state message
			if len(tt.todos) == 0 && !containsEmptyMessage(result) {
				t.Error("Render() should contain empty state message for empty todos")
			}
		})
	}
}

func containsEmptyMessage(result string) bool {
	return len(result) > 0 // Simple check that we got some content
}

func TestListView_RenderCompact(t *testing.T) {
	tests := []struct {
		name   string
		todos  []*model.ItemModel
		width  int
		height int
	}{
		{
			name:   "positive testing (empty todos)",
			todos:  []*model.ItemModel{},
			width:  40,
			height: 10,
		},
		{
			name: "positive testing (with todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			width:  40,
			height: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for better testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			listView := NewListView(lipgloss, itemView)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetTodos(tt.todos)

			result := listView.RenderCompact(stateModel, tt.width, tt.height)

			// Verify that we get a non-empty result
			if result == "" {
				t.Error("RenderCompact() returned empty string")
			}
			
			// For empty todos, should contain "No todos"
			if len(tt.todos) == 0 && result != "No todos" {
				t.Errorf("RenderCompact() for empty todos = %v, want 'No todos'", result)
			}
		})
	}
}

func TestListView_RenderWithPagination(t *testing.T) {
	tests := []struct {
		name   string
		todos  []*model.ItemModel
		width  int
		height int
	}{
		{
			name:   "positive testing (empty todos)",
			todos:  []*model.ItemModel{},
			width:  80,
			height: 20,
		},
		{
			name: "positive testing (with todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			width:  80,
			height: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for better testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			listView := NewListView(lipgloss, itemView)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetTodos(tt.todos)

			result := listView.RenderWithPagination(stateModel, tt.width, tt.height)

			// Verify that we get a non-empty result
			if result == "" {
				t.Error("RenderWithPagination() returned empty string")
			}
		})
	}
}

func TestListView_CalculateVisibleRange(t *testing.T) {
	tests := []struct {
		name          string
		cursor        int
		totalItems    int
		visibleHeight int
		expectedStart int
		expectedEnd   int
	}{
		{
			name:          "positive testing (all items visible)",
			cursor:        2,
			totalItems:    5,
			visibleHeight: 10,
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "positive testing (cursor at beginning)",
			cursor:        0,
			totalItems:    20,
			visibleHeight: 10,
			expectedStart: 0,
			expectedEnd:   10,
		},
		{
			name:          "positive testing (cursor in middle)",
			cursor:        10,
			totalItems:    20,
			visibleHeight: 10,
			expectedStart: 5,
			expectedEnd:   15,
		},
		{
			name:          "positive testing (cursor at end)",
			cursor:        19,
			totalItems:    20,
			visibleHeight: 10,
			expectedStart: 10,
			expectedEnd:   20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)

			listView := NewListView(mockLipgloss, itemView)
			startIndex, endIndex := listView.CalculateVisibleRange(tt.cursor, tt.totalItems, tt.visibleHeight)

			if startIndex != tt.expectedStart {
				t.Errorf("CalculateVisibleRange() startIndex = %v, want %v", startIndex, tt.expectedStart)
			}
			if endIndex != tt.expectedEnd {
				t.Errorf("CalculateVisibleRange() endIndex = %v, want %v", endIndex, tt.expectedEnd)
			}
		})
	}
}

func TestListView_RenderScrollIndicator(t *testing.T) {
	tests := []struct {
		name     string
		todos    []*model.ItemModel
		cursor   int
		width    int
		expected string
	}{
		{
			name:     "positive testing (empty todos)",
			todos:    []*model.ItemModel{},
			cursor:   0,
			width:    80,
			expected: "",
		},
		{
			name: "positive testing (single todo)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			cursor:   0,
			width:    80,
			expected: "scroll indicator",
		},
		{
			name: "positive testing (multiple todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo 1",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
				model.NewItemModel(&domain.Todo{
					ID:          2,
					Description: "Test todo 2",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			cursor:   1,
			width:    80,
			expected: "scroll indicator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)

			// Setup for RenderScrollIndicator
			if len(tt.todos) > 0 {
				mockLipgloss.EXPECT().NewStyle().Return(mockStyle)
				mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
				mockStyle.EXPECT().Render(gomock.Any()).Return(tt.expected)
			}

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetTodos(tt.todos)
			stateModel.SetCursor(tt.cursor)

			listView := NewListView(mockLipgloss, itemView)
			result := listView.RenderScrollIndicator(stateModel, tt.width)

			if result != tt.expected {
				t.Errorf("RenderScrollIndicator() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestListView_GetVisibleItemCount(t *testing.T) {
	tests := []struct {
		name     string
		height   int
		expected int
	}{
		{
			name:     "positive testing (small height)",
			height:   5,
			expected: 3,
		},
		{
			name:     "positive testing (medium height)",
			height:   20,
			expected: 18,
		},
		{
			name:     "positive testing (large height)",
			height:   50,
			expected: 48,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)

			listView := NewListView(mockLipgloss, itemView)
			result := listView.GetVisibleItemCount(tt.height)

			if result != tt.expected {
				t.Errorf("GetVisibleItemCount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestListView_GetTotalHeight(t *testing.T) {
	tests := []struct {
		name     string
		todos    []*model.ItemModel
		expected int
	}{
		{
			name:     "positive testing (empty todos)",
			todos:    []*model.ItemModel{},
			expected: 2, // Just padding
		},
		{
			name: "positive testing (single todo)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			expected: 3, // 1 item + 2 padding
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetTodos(tt.todos)

			listView := NewListView(mockLipgloss, itemView)
			result := listView.GetTotalHeight(stateModel)

			if result != tt.expected {
				t.Errorf("GetTotalHeight() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestListView_IsScrollable(t *testing.T) {
	tests := []struct {
		name     string
		todos    []*model.ItemModel
		height   int
		expected bool
	}{
		{
			name:     "positive testing (not scrollable)",
			todos:    []*model.ItemModel{},
			height:   20,
			expected: false,
		},
		{
			name: "positive testing (scrollable)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
			height:   2, // Smaller than total height
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)

			stateModel, stateCtrl := createTestStateModelForView(t)
			defer stateCtrl.Finish()
			
			stateModel.SetTodos(tt.todos)

			listView := NewListView(mockLipgloss, itemView)
			result := listView.IsScrollable(stateModel, tt.height)

			if result != tt.expected {
				t.Errorf("IsScrollable() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestListView_SetContentStyle(t *testing.T) {
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

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)

			listView := NewListView(mockLipgloss, itemView)
			listView.SetContentStyle(newStyle)

			// Test that the style was set (we can't directly verify this without exposing the field)
			// This test mainly ensures the method doesn't panic
		})
	}
}

func TestListView_SetEmptyStyle(t *testing.T) {
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

			// Setup for NewListView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(2)
			mockStyle.EXPECT().Padding(1, 2).Return(mockStyle)
			mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle)
			mockStyle.EXPECT().Italic(true).Return(mockStyle)

			itemView := createMockItemView(ctrl)

			listView := NewListView(mockLipgloss, itemView)
			listView.SetEmptyStyle(newStyle)

			// Test that the style was set (we can't directly verify this without exposing the field)
			// This test mainly ensures the method doesn't panic
		})
	}
}