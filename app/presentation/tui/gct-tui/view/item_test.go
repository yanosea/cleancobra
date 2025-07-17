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

func TestNewItemView(t *testing.T) {
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

			// Setup mock expectations for style creation
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			if itemView == nil {
				t.Error("NewItemView returned nil")
			}
			if itemView.lipgloss != mockLipgloss {
				t.Error("lipgloss proxy not set correctly")
			}
		})
	}
}

func TestItemView_RenderItem(t *testing.T) {
	tests := []struct {
		name     string
		item     *model.ItemModel
		width    int
		want     string
		wantErr  bool
	}{
		{
			name: "positive testing with incomplete todo",
			item: model.NewItemModel(&domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}),
			width:   50,
			want:    "rendered content",
			wantErr: false,
		},
		{
			name: "positive testing with completed todo",
			item: model.NewItemModel(&domain.Todo{
				ID:          2,
				Description: "Completed todo",
				Done:        true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}),
			width:   50,
			want:    "rendered content",
			wantErr: false,
		},
		{
			name:    "negative testing (nil item failed)",
			item:    nil,
			width:   50,
			want:    "",
			wantErr: false,
		},
		{
			name:    "negative testing (item with nil todo failed)",
			item:    model.NewItemModel(nil),
			width:   50,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup mock expectations for NewItemView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			if tt.item != nil && tt.item.Todo() != nil {
				// Setup expectations for rendering
				mockStyle.EXPECT().Width(tt.width - 2).Return(mockStyle).Times(1)
				mockStyle.EXPECT().Render(gomock.Any()).Return("rendered content").Times(1)
			}

			result := itemView.RenderItem(tt.item, tt.width)

			if tt.want == "" && result != "" {
				t.Errorf("RenderItem() = %v, want empty string", result)
			}
			if tt.want != "" && result == "" {
				t.Errorf("RenderItem() = empty string, want non-empty")
			}
		})
	}
}

func TestItemView_RenderItemWithSelection(t *testing.T) {
	tests := []struct {
		name       string
		item       *model.ItemModel
		isSelected bool
		width      int
		want       string
		wantErr    bool
	}{
		{
			name: "positive testing with selected item",
			item: model.NewItemModel(&domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}),
			isSelected: true,
			width:      50,
			want:       "rendered content",
			wantErr:    false,
		},
		{
			name: "positive testing with unselected item",
			item: model.NewItemModel(&domain.Todo{
				ID:          2,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}),
			isSelected: false,
			width:      50,
			want:       "rendered content",
			wantErr:    false,
		},
		{
			name:       "negative testing (nil item failed)",
			item:       nil,
			isSelected: false,
			width:      50,
			want:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup mock expectations for NewItemView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			if tt.item != nil && tt.item.Todo() != nil {
				// Setup expectations for rendering with selection
				mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(1)
				mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(1)
				if tt.isSelected {
					mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
					mockStyle.EXPECT().Foreground(lipgloss.Color("15")).Return(mockStyle).Times(1)
				}
				mockStyle.EXPECT().Width(tt.width - 2).Return(mockStyle).Times(1)
				mockStyle.EXPECT().Render(gomock.Any()).Return("rendered content").Times(1)
			}

			result := itemView.RenderItemWithSelection(tt.item, tt.isSelected, tt.width)

			if tt.want == "" && result != "" {
				t.Errorf("RenderItemWithSelection() = %v, want empty string", result)
			}
			if tt.want != "" && result == "" {
				t.Errorf("RenderItemWithSelection() = empty string, want non-empty")
			}
		})
	}
}

func TestItemView_RenderCompactItem(t *testing.T) {
	tests := []struct {
		name       string
		item       *model.ItemModel
		isSelected bool
		width      int
		want       string
		wantErr    bool
	}{
		{
			name: "positive testing with normal width",
			item: model.NewItemModel(&domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}),
			isSelected: false,
			width:      50,
			want:       "rendered content",
			wantErr:    false,
		},
		{
			name: "positive testing with long description truncation",
			item: model.NewItemModel(&domain.Todo{
				ID:          2,
				Description: "This is a very long todo description that should be truncated",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}),
			isSelected: false,
			width:      20,
			want:       "rendered content",
			wantErr:    false,
		},
		{
			name:       "negative testing (nil item failed)",
			item:       nil,
			isSelected: false,
			width:      50,
			want:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup mock expectations for NewItemView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			if tt.item != nil && tt.item.Todo() != nil {
				// Setup expectations for compact rendering
				mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(1)
				if tt.item.Todo().Done {
					mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
				}
				if tt.isSelected {
					mockStyle.EXPECT().Bold(true).Return(mockStyle).Times(1)
				}
				mockStyle.EXPECT().Render(gomock.Any()).Return("rendered content").Times(1)
			}

			result := itemView.RenderCompactItem(tt.item, tt.isSelected, tt.width)

			if tt.want == "" && result != "" {
				t.Errorf("RenderCompactItem() = %v, want empty string", result)
			}
			if tt.want != "" && result == "" {
				t.Errorf("RenderCompactItem() = empty string, want non-empty")
			}
		})
	}
}

func TestItemView_GetItemHeight(t *testing.T) {
	tests := []struct {
		name string
		item *model.ItemModel
		want int
	}{
		{
			name: "positive testing with normal item",
			item: model.NewItemModel(&domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}),
			want: 1,
		},
		{
			name: "positive testing with editing item",
			item: func() *model.ItemModel {
				item := model.NewItemModel(&domain.Todo{
					ID:          2,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				item.SetEditing(true)
				return item
			}(),
			want: 3,
		},
		{
			name: "negative testing (nil item failed)",
			item: nil,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup mock expectations for NewItemView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			result := itemView.GetItemHeight(tt.item)

			if result != tt.want {
				t.Errorf("GetItemHeight() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestItemView_GetMinimumWidth(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "positive testing",
			want: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup mock expectations for NewItemView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			result := itemView.GetMinimumWidth()

			if result != tt.want {
				t.Errorf("GetMinimumWidth() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestItemView_CreateStatusIndicator(t *testing.T) {
	tests := []struct {
		name     string
		done     bool
		selected bool
		want     string
	}{
		{
			name:     "positive testing with completed and selected",
			done:     true,
			selected: true,
			want:     "[✓]",
		},
		{
			name:     "positive testing with incomplete and unselected",
			done:     false,
			selected: false,
			want:     "[ ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup mock expectations for NewItemView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			// Setup expectations for CreateStatusIndicator
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(1)
			if tt.done {
				mockStyle.EXPECT().Foreground(lipgloss.Color("2")).Return(mockStyle).Times(1)
			} else {
				mockStyle.EXPECT().Foreground(lipgloss.Color("8")).Return(mockStyle).Times(1)
			}
			if tt.selected {
				mockStyle.EXPECT().Bold(true).Return(mockStyle).Times(1)
			}
			if tt.done {
				mockStyle.EXPECT().Render("✓").Return("✓").Times(1)
			} else {
				mockStyle.EXPECT().Render(" ").Return(" ").Times(1)
			}

			result := itemView.CreateStatusIndicator(tt.done, tt.selected)

			if result != tt.want {
				t.Errorf("CreateStatusIndicator() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestItemView_CreateSelectionIndicator(t *testing.T) {
	tests := []struct {
		name     string
		selected bool
		want     string
	}{
		{
			name:     "positive testing with selected",
			selected: true,
			want:     "▶",
		},
		{
			name:     "positive testing with unselected",
			selected: false,
			want:     " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLipgloss := proxy.NewMockLipgloss(ctrl)
			mockStyle := proxy.NewMockStyle(ctrl)

			// Setup mock expectations for NewItemView
			mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(4)
			mockStyle.EXPECT().Padding(0, 1).Return(mockStyle).Times(4)
			mockStyle.EXPECT().Background(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Foreground(lipgloss.Color("240")).Return(mockStyle).Times(1)
			mockStyle.EXPECT().Border(lipgloss.RoundedBorder(), true).Return(mockStyle).Times(1)

			itemView := NewItemView(mockLipgloss)

			if tt.selected {
				// Setup expectations for selected indicator
				mockLipgloss.EXPECT().NewStyle().Return(mockStyle).Times(1)
				mockStyle.EXPECT().Foreground(lipgloss.Color("12")).Return(mockStyle).Times(1)
				mockStyle.EXPECT().Bold(true).Return(mockStyle).Times(1)
				mockStyle.EXPECT().Render("▶").Return("▶").Times(1)
			}

			result := itemView.CreateSelectionIndicator(tt.selected)

			if result != tt.want {
				t.Errorf("CreateSelectionIndicator() = %v, want %v", result, tt.want)
			}
		})
	}
}