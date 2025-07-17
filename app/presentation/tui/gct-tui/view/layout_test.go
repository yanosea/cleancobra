package view

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewLayoutView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for better testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)

			layoutView := NewLayoutView(lipgloss, itemView)

			if layoutView == nil {
				t.Error("NewLayoutView returned nil")
			}
			if layoutView.lipgloss != lipgloss {
				t.Error("lipgloss proxy not set correctly")
			}
			if layoutView.itemView != itemView {
				t.Error("itemView not set correctly")
			}
			if layoutView.headerView == nil {
				t.Error("headerView not initialized")
			}
			if layoutView.footerView == nil {
				t.Error("footerView not initialized")
			}
			if layoutView.listView == nil {
				t.Error("listView not initialized")
			}
		})
	}
}

func TestLayoutView_Render(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			layoutView := NewLayoutView(lipgloss, itemView)

			// Create a simple state model for testing
			bubbles := proxy.NewBubbles()
			stateModel := model.NewStateModel(nil, nil, nil, nil, bubbles)
			stateModel.SetSize(80, 24)

			result := layoutView.Render(stateModel)

			if result == "" {
				t.Error("Render() returned empty string")
			}
		})
	}
}

func TestLayoutView_RenderCompact(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			layoutView := NewLayoutView(lipgloss, itemView)

			// Create a simple state model for testing
			bubbles := proxy.NewBubbles()
			stateModel := model.NewStateModel(nil, nil, nil, nil, bubbles)
			stateModel.SetSize(30, 8) // Small size for compact mode

			result := layoutView.RenderCompact(stateModel)

			if result == "" {
				t.Error("RenderCompact() returned empty string")
			}
		})
	}
}

func TestLayoutView_calculateVisibleRange(t *testing.T) {
	tests := []struct {
		name          string
		cursor        int
		totalItems    int
		visibleHeight int
		wantStart     int
		wantEnd       int
	}{
		{
			name:          "positive testing with all items visible",
			cursor:        2,
			totalItems:    5,
			visibleHeight: 10,
			wantStart:     0,
			wantEnd:       5,
		},
		{
			name:          "positive testing with scrolling needed",
			cursor:        10,
			totalItems:    20,
			visibleHeight: 5,
			wantStart:     8,
			wantEnd:       13,
		},
		{
			name:          "positive testing with cursor at beginning",
			cursor:        0,
			totalItems:    20,
			visibleHeight: 5,
			wantStart:     0,
			wantEnd:       5,
		},
		{
			name:          "positive testing with cursor at end",
			cursor:        19,
			totalItems:    20,
			visibleHeight: 5,
			wantStart:     15,
			wantEnd:       20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)

			// Use ListView for CalculateVisibleRange since it was moved there
			listView := NewListView(lipgloss, itemView)
			startIndex, endIndex := listView.CalculateVisibleRange(tt.cursor, tt.totalItems, tt.visibleHeight)

			if startIndex != tt.wantStart {
				t.Errorf("calculateVisibleRange() startIndex = %v, want %v", startIndex, tt.wantStart)
			}
			if endIndex != tt.wantEnd {
				t.Errorf("calculateVisibleRange() endIndex = %v, want %v", endIndex, tt.wantEnd)
			}
		})
	}
}

func TestLayoutView_getCompletedCount(t *testing.T) {
	tests := []struct {
		name  string
		todos []*model.ItemModel
		want  int
	}{
		{
			name: "positive testing with mixed todos",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{ID: 1, Description: "Todo 1", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
				model.NewItemModel(&domain.Todo{ID: 2, Description: "Todo 2", Done: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
				model.NewItemModel(&domain.Todo{ID: 3, Description: "Todo 3", Done: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
			},
			want: 2,
		},
		{
			name: "positive testing with no completed todos",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{ID: 1, Description: "Todo 1", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
				model.NewItemModel(&domain.Todo{ID: 2, Description: "Todo 2", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
			},
			want: 0,
		},
		{
			name:  "positive testing with empty todos",
			todos: []*model.ItemModel{},
			want:  0,
		},
		{
			name: "negative testing (todos with nil todo failed)",
			todos: []*model.ItemModel{
				model.NewItemModel(nil),
				model.NewItemModel(&domain.Todo{ID: 1, Description: "Todo 1", Done: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for testing
			lipgloss := proxy.NewLipgloss()

			// Use HeaderView for GetCompletedCount since it was moved there
			headerView := NewHeaderView(lipgloss)
			result := headerView.GetCompletedCount(tt.todos)

			if result != tt.want {
				t.Errorf("getCompletedCount() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestLayoutView_GetMinimumSize(t *testing.T) {
	tests := []struct {
		name       string
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "positive testing",
			wantWidth:  40,
			wantHeight: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			layoutView := NewLayoutView(lipgloss, itemView)

			width, height := layoutView.GetMinimumSize()

			if width != tt.wantWidth {
				t.Errorf("GetMinimumSize() width = %v, want %v", width, tt.wantWidth)
			}
			if height != tt.wantHeight {
				t.Errorf("GetMinimumSize() height = %v, want %v", height, tt.wantHeight)
			}
		})
	}
}

func TestLayoutView_IsCompactMode(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   bool
	}{
		{
			name:   "positive testing with large terminal",
			width:  80,
			height: 24,
			want:   false,
		},
		{
			name:   "positive testing with small width",
			width:  30,
			height: 24,
			want:   true,
		},
		{
			name:   "positive testing with small height",
			width:  80,
			height: 8,
			want:   true,
		},
		{
			name:   "positive testing with minimum size",
			width:  40,
			height: 10,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			layoutView := NewLayoutView(lipgloss, itemView)

			result := layoutView.IsCompactMode(tt.width, tt.height)

			if result != tt.want {
				t.Errorf("IsCompactMode() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestLayoutView_RenderScrollIndicator(t *testing.T) {
	tests := []struct {
		name   string
		todos  []*model.ItemModel
		cursor int
		width  int
		want   string
	}{
		{
			name: "positive testing with todos",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{ID: 1, Description: "Todo 1", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
				model.NewItemModel(&domain.Todo{ID: 2, Description: "Todo 2", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
				model.NewItemModel(&domain.Todo{ID: 3, Description: "Todo 3", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}),
			},
			cursor: 1,
			width:  50,
			want:   "non-empty",
		},
		{
			name:   "positive testing with empty todos",
			todos:  []*model.ItemModel{},
			cursor: 0,
			width:  50,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use real implementations for testing
			lipgloss := proxy.NewLipgloss()
			itemView := NewItemView(lipgloss)
			layoutView := NewLayoutView(lipgloss, itemView)

			// Create a simple state model for testing
			bubbles := proxy.NewBubbles()
			stateModel := model.NewStateModel(nil, nil, nil, nil, bubbles)
			stateModel.SetTodos(tt.todos)
			stateModel.SetCursor(tt.cursor)

			result := layoutView.RenderScrollIndicator(stateModel, tt.width)

			if tt.want == "" && result != "" {
				t.Errorf("RenderScrollIndicator() = %v, want empty string", result)
			}
			if tt.want != "" && result == "" {
				t.Errorf("RenderScrollIndicator() = empty string, want non-empty")
			}
		})
	}
}