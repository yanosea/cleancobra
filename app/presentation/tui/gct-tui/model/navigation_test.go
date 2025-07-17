package model

import (
	"testing"
)

func TestNewNavigationState(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNavigationState()
			if got == nil {
				t.Error("NewNavigationState() returned nil")
			}
			if got.Cursor() != 0 {
				t.Errorf("NewNavigationState() cursor = %v, want 0", got.Cursor())
			}
			if got.HasSelection() {
				t.Error("NewNavigationState() should not have selection initially")
			}
		})
	}
}

func TestNavigationState_SetCursor(t *testing.T) {
	tests := []struct {
		name     string
		cursor   int
		maxItems int
		want     int
	}{
		{
			name:     "positive testing - valid cursor",
			cursor:   2,
			maxItems: 5,
			want:     2,
		},
		{
			name:     "positive testing - cursor too high",
			cursor:   10,
			maxItems: 5,
			want:     4,
		},
		{
			name:     "positive testing - negative cursor",
			cursor:   -1,
			maxItems: 5,
			want:     0,
		},
		{
			name:     "positive testing - zero max items",
			cursor:   1,
			maxItems: 0,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := NewNavigationState()
			nav.SetCursor(tt.cursor, tt.maxItems)
			got := nav.Cursor()
			if got != tt.want {
				t.Errorf("NavigationState.SetCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationState_MoveCursor(t *testing.T) {
	tests := []struct {
		name         string
		initialPos   int
		maxItems     int
		operation    func(*NavigationState, int)
		expectedPos  int
	}{
		{
			name:        "positive testing - move up",
			initialPos:  2,
			maxItems:    5,
			operation:   (*NavigationState).MoveCursorUp,
			expectedPos: 1,
		},
		{
			name:        "positive testing - move down",
			initialPos:  2,
			maxItems:    5,
			operation:   (*NavigationState).MoveCursorDown,
			expectedPos: 3,
		},
		{
			name:        "positive testing - move up from top",
			initialPos:  0,
			maxItems:    5,
			operation:   (*NavigationState).MoveCursorUp,
			expectedPos: 0,
		},
		{
			name:        "positive testing - move down from bottom",
			initialPos:  4,
			maxItems:    5,
			operation:   (*NavigationState).MoveCursorDown,
			expectedPos: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := NewNavigationState()
			nav.SetCursor(tt.initialPos, tt.maxItems)
			tt.operation(nav, tt.maxItems)
			got := nav.Cursor()
			if got != tt.expectedPos {
				t.Errorf("NavigationState cursor = %v, want %v", got, tt.expectedPos)
			}
		})
	}
}

func TestNavigationState_MoveCursorToEdges(t *testing.T) {
	tests := []struct {
		name        string
		maxItems    int
		operation   func(*NavigationState, int)
		expectedPos int
	}{
		{
			name:        "positive testing - move to top",
			maxItems:    5,
			operation:   func(nav *NavigationState, _ int) { nav.MoveCursorToTop() },
			expectedPos: 0,
		},
		{
			name:        "positive testing - move to bottom",
			maxItems:    5,
			operation:   func(nav *NavigationState, maxItems int) { nav.MoveCursorToBottom(maxItems) },
			expectedPos: 4,
		},
		{
			name:        "positive testing - move to bottom with zero items",
			maxItems:    0,
			operation:   func(nav *NavigationState, maxItems int) { nav.MoveCursorToBottom(maxItems) },
			expectedPos: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := NewNavigationState()
			nav.SetCursor(2, tt.maxItems) // Start in middle
			tt.operation(nav, tt.maxItems)
			got := nav.Cursor()
			if got != tt.expectedPos {
				t.Errorf("NavigationState cursor = %v, want %v", got, tt.expectedPos)
			}
		})
	}
}

func TestNavigationState_Selection(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*NavigationState)
		index    int
		selected bool
		want     bool
	}{
		{
			name:     "positive testing - set selected",
			setup:    func(nav *NavigationState) {},
			index:    1,
			selected: true,
			want:     true,
		},
		{
			name:     "positive testing - set unselected",
			setup:    func(nav *NavigationState) { nav.SetSelected(1, true) },
			index:    1,
			selected: false,
			want:     false,
		},
		{
			name:     "positive testing - check unselected index",
			setup:    func(nav *NavigationState) { nav.SetSelected(1, true) },
			index:    2,
			selected: false,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := NewNavigationState()
			tt.setup(nav)
			nav.SetSelected(tt.index, tt.selected)
			got := nav.IsSelected(tt.index)
			if got != tt.want {
				t.Errorf("NavigationState.IsSelected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationState_ClearSelection(t *testing.T) {
	tests := []struct {
		name string
		setup func(*NavigationState)
	}{
		{
			name: "positive testing - clear multiple selections",
			setup: func(nav *NavigationState) {
				nav.SetSelected(1, true)
				nav.SetSelected(2, true)
				nav.SetSelected(3, true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := NewNavigationState()
			tt.setup(nav)
			
			// Verify selections exist
			if !nav.HasSelection() {
				t.Error("Expected selections to exist before clearing")
			}
			
			nav.ClearSelection()
			
			// Verify all selections are cleared
			if nav.HasSelection() {
				t.Error("Expected no selections after clearing")
			}
			
			indices := nav.GetSelectedIndices()
			if len(indices) != 0 {
				t.Errorf("Expected 0 selected indices, got %d", len(indices))
			}
		})
	}
}

func TestNavigationState_GetSelectedIndices(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*NavigationState)
		wantLen  int
	}{
		{
			name:    "positive testing - no selections",
			setup:   func(nav *NavigationState) {},
			wantLen: 0,
		},
		{
			name: "positive testing - multiple selections",
			setup: func(nav *NavigationState) {
				nav.SetSelected(1, true)
				nav.SetSelected(3, true)
				nav.SetSelected(5, true)
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := NewNavigationState()
			tt.setup(nav)
			got := nav.GetSelectedIndices()
			if len(got) != tt.wantLen {
				t.Errorf("NavigationState.GetSelectedIndices() length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}