package view

import (
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"

	"github.com/yanosea/gct/pkg/proxy"
)

// ItemView handles rendering of individual todo items
type ItemView struct {
	completedStyleProxy proxy.Style
	editingStyleProxy   proxy.Style
	lipglossProxy       proxy.Lipgloss
	normalStyleProxy    proxy.Style
	selectedStyleProxy  proxy.Style
}

// NewItemView creates a new ItemView with the given lipgloss proxy
func NewItemView(lipgrossProxy proxy.Lipgloss) *ItemView {
	return &ItemView{
		lipglossProxy:       lipgrossProxy,
		normalStyleProxy:    lipgrossProxy.NewStyle().Padding(0, 1),
		selectedStyleProxy:  lipgrossProxy.NewStyle().Padding(0, 1).Background(lipgrossProxy.Color("240")),
		completedStyleProxy: lipgrossProxy.NewStyle().Padding(0, 1).Foreground(lipgrossProxy.Color("240")),
		editingStyleProxy:   lipgrossProxy.NewStyle().Padding(0, 1).Border(lipgrossProxy.RoundedBorder(), true),
	}
}

// RenderItem renders a single todo item with appropriate styling
func (v *ItemView) RenderItem(item *model.ItemModel, width int) string {
	if item == nil || item.Todo() == nil {
		return ""
	}

	todo := item.Todo()

	// Status indicator
	status := "[ ]"
	if todo.Done {
		status = "[✓]"
	}

	// Build the content
	content := status + " " + todo.Description

	// Apply styling based on item state
	style := v.getItemStyle(item)

	// Adjust width for responsive layout
	if width > 0 {
		style = style.Width(width - 2) // Account for padding
	}

	return style.Render(content)
}

// RenderItemWithSelection renders an item with selection highlighting
func (v *ItemView) RenderItemWithSelection(item *model.ItemModel, isSelected bool, width int) string {
	if item == nil || item.Todo() == nil {
		return ""
	}

	todo := item.Todo()

	// Status indicator with better Unicode symbols
	status := "○"
	if todo.Done {
		status = "●"
	}

	// Selection indicator
	cursor := "  "
	if isSelected {
		cursor = "▶ "
	}

	// Build the content
	content := cursor + status + " " + todo.Description

	// Apply styling
	style := v.getItemStyleWithSelection(item, isSelected)

	// Adjust width for responsive layout
	if width > 0 {
		style = style.Width(width - 2)
	}

	return style.Render(content)
}

// RenderCompactItem renders a compact version of the item for smaller terminals
func (v *ItemView) RenderCompactItem(item *model.ItemModel, isSelected bool, width int) string {
	if item == nil || item.Todo() == nil {
		return ""
	}

	todo := item.Todo()

	// Compact status indicator
	status := "□"
	if todo.Done {
		status = "■"
	}

	// Selection indicator
	cursor := " "
	if isSelected {
		cursor = ">"
	}

	// Truncate description if too long
	description := todo.Description
	maxDescLength := width - 10 // Account for cursor, status, and padding
	if maxDescLength > 0 && len(description) > maxDescLength {
		description = description[:maxDescLength-3] + "..."
	}

	content := cursor + status + " " + description

	// Apply minimal styling for compact view
	style := v.lipglossProxy.NewStyle()
	if isSelected {
		style = style.Bold(true)
	}
	if todo.Done {
		style = style.Foreground(v.lipglossProxy.Color("240")) // Gray color
	}

	return style.Render(content)
}

// getItemStyle returns the appropriate style for an item based on its state
func (v *ItemView) getItemStyle(item *model.ItemModel) proxy.Style {
	if item.IsEditing() {
		return v.editingStyleProxy
	}

	if item.Todo().Done {
		return v.completedStyleProxy
	}

	if item.IsSelected() {
		return v.selectedStyleProxy
	}

	return v.normalStyleProxy
}

// getItemStyleWithSelection returns the appropriate style with selection highlighting
func (v *ItemView) getItemStyleWithSelection(item *model.ItemModel, isSelected bool) proxy.Style {
	baseStyle := v.lipglossProxy.NewStyle().Padding(0, 1)

	if item.IsEditing() {
		baseStyle = baseStyle.Border(v.lipglossProxy.RoundedBorder(), true)
	}

	if isSelected {
		baseStyle = baseStyle.Background(v.lipglossProxy.Color("240")).
			Foreground(v.lipglossProxy.Color("15")) // Highlight selected item
	}

	if item.Todo().Done {
		baseStyle = baseStyle.Foreground(v.lipglossProxy.Color("240")).
			Strikethrough(true) // Gray and strikethrough for completed
	}

	return baseStyle
}

// GetItemHeight returns the height needed to render an item
func (v *ItemView) GetItemHeight(item *model.ItemModel) int {
	if item == nil {
		return 0
	}

	// Most items are single line, but editing mode might need more space
	if item.IsEditing() {
		return 3 // Account for border
	}

	return 1
}

// GetMinimumWidth returns the minimum width needed to render items properly
func (v *ItemView) GetMinimumWidth() int {
	return 20 // Minimum width for reasonable display
}

// CreateStatusIndicator creates a styled status indicator
func (v *ItemView) CreateStatusIndicator(done bool, selected bool) string {
	var indicator string
	var style proxy.Style

	if done {
		indicator = "✓"
		style = v.lipglossProxy.NewStyle().Foreground(v.lipglossProxy.Color("2")) // Green
	} else {
		indicator = " "
		style = v.lipglossProxy.NewStyle().Foreground(v.lipglossProxy.Color("8")) // Gray
	}

	if selected {
		style = style.Bold(true)
	}

	return "[" + style.Render(indicator) + "]"
}

// CreateSelectionIndicator creates a selection cursor indicator
func (v *ItemView) CreateSelectionIndicator(selected bool) string {
	if selected {
		return v.lipglossProxy.NewStyle().
			Foreground(v.lipglossProxy.Color("12")).
			Bold(true).
			Render("▶")
	}
	return " "
}
