package model

// NavigationState manages cursor position and selection state
type NavigationState struct {
	cursor   int
	selected map[int]struct{}
}

// NewNavigationState creates a new navigation state
func NewNavigationState() *NavigationState {
	return &NavigationState{
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

// Cursor returns the current cursor position
func (n *NavigationState) Cursor() int {
	return n.cursor
}

// SetCursor sets the cursor position with bounds checking
func (n *NavigationState) SetCursor(cursor, maxItems int) {
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= maxItems {
		cursor = maxItems - 1
	}
	if cursor < 0 {
		cursor = 0
	}
	n.cursor = cursor
}

// MoveCursorUp moves the cursor up by one position
func (n *NavigationState) MoveCursorUp(maxItems int) {
	n.SetCursor(n.cursor-1, maxItems)
}

// MoveCursorDown moves the cursor down by one position
func (n *NavigationState) MoveCursorDown(maxItems int) {
	n.SetCursor(n.cursor+1, maxItems)
}

// MoveCursorToTop moves the cursor to the first position
func (n *NavigationState) MoveCursorToTop() {
	n.cursor = 0
}

// MoveCursorToBottom moves the cursor to the last position
func (n *NavigationState) MoveCursorToBottom(maxItems int) {
	if maxItems > 0 {
		n.cursor = maxItems - 1
	} else {
		n.cursor = 0
	}
}

// IsSelected returns whether the given index is selected
func (n *NavigationState) IsSelected(index int) bool {
	_, exists := n.selected[index]
	return exists
}

// SetSelected sets the selection state for the given index
func (n *NavigationState) SetSelected(index int, selected bool) {
	if selected {
		n.selected[index] = struct{}{}
	} else {
		delete(n.selected, index)
	}
}

// ClearSelection clears all selections
func (n *NavigationState) ClearSelection() {
	n.selected = make(map[int]struct{})
}

// GetSelectedIndices returns a slice of all selected indices
func (n *NavigationState) GetSelectedIndices() []int {
	indices := make([]int, 0, len(n.selected))
	for index := range n.selected {
		indices = append(indices, index)
	}
	return indices
}

// HasSelection returns whether any items are selected
func (n *NavigationState) HasSelection() bool {
	return len(n.selected) > 0
}
