package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
)

// ListView handles the todo list rendering
type ListView struct {
	lipgloss     proxy.Lipgloss
	itemView     *ItemView
	contentStyle proxy.Style
	emptyStyle   proxy.Style
}

// NewListView creates a new ListView with the given dependencies
func NewListView(lg proxy.Lipgloss, itemView *ItemView) *ListView {
	return &ListView{
		lipgloss:     lg,
		itemView:     itemView,
		contentStyle: lg.NewStyle().Padding(1, 2),
		emptyStyle:   lg.NewStyle().Foreground(lipgloss.Color("8")).Italic(true),
	}
}

// Render renders the todo list
func (v *ListView) Render(stateModel *model.StateModel, width, height int) string {
	todos := stateModel.Todos()

	if len(todos) == 0 {
		return v.renderEmptyState(width, height)
	}

	return v.renderTodoList(stateModel, width, height)
}

// RenderCompact renders a compact version of the todo list
func (v *ListView) RenderCompact(stateModel *model.StateModel, width, height int) string {
	todos := stateModel.Todos()

	if len(todos) == 0 {
		return v.emptyStyle.Render("No todos")
	}

	var todoLines []string
	cursor := stateModel.Cursor()
	visibleHeight := height - 1

	startIndex, endIndex := v.CalculateVisibleRange(cursor, len(todos), visibleHeight)

	for i := startIndex; i < endIndex && i < len(todos); i++ {
		isSelected := i == cursor
		todoLine := v.itemView.RenderCompactItem(todos[i], isSelected, width)
		todoLines = append(todoLines, todoLine)
	}

	return strings.Join(todoLines, "\n")
}

// RenderWithPagination renders the list with pagination indicators
func (v *ListView) RenderWithPagination(stateModel *model.StateModel, width, height int) string {
	todos := stateModel.Todos()

	if len(todos) == 0 {
		return v.renderEmptyState(width, height)
	}

	cursor := stateModel.Cursor()
	visibleHeight := height - 4 // Account for pagination indicators

	// Calculate pagination
	startIndex, endIndex := v.CalculateVisibleRange(cursor, len(todos), visibleHeight)

	// Render todos
	var todoLines []string
	for i := startIndex; i < endIndex && i < len(todos); i++ {
		isSelected := i == cursor
		todoLine := v.itemView.RenderItemWithSelection(todos[i], isSelected, width-4)
		todoLines = append(todoLines, todoLine)
	}

	content := strings.Join(todoLines, "\n")

	// Add pagination indicators
	if startIndex > 0 {
		content = v.lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("↑ More items above") + "\n" + content
	}

	if endIndex < len(todos) {
		content = content + "\n" + v.lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("↓ More items below")
	}

	return v.contentStyle.Width(width).Height(height).Render(content)
}

// renderEmptyState renders the empty state when no todos exist
func (v *ListView) renderEmptyState(width, height int) string {
	emptyMessage := "No todos yet. Press 'a' to add your first todo!"

	return v.contentStyle.
		Width(width).
		Height(height).
		AlignHorizontal(v.lipgloss.Center()).
		AlignVertical(v.lipgloss.Center()).
		Render(v.emptyStyle.Render(emptyMessage))
}

// renderTodoList renders the main todo list
func (v *ListView) renderTodoList(stateModel *model.StateModel, width, height int) string {
	todos := stateModel.Todos()
	cursor := stateModel.Cursor()

	// Calculate visible range for scrolling
	visibleHeight := height - 2 // Account for padding
	startIndex, endIndex := v.CalculateVisibleRange(cursor, len(todos), visibleHeight)

	var todoLines []string
	for i := startIndex; i < endIndex && i < len(todos); i++ {
		isSelected := i == cursor
		todoLine := v.itemView.RenderItemWithSelection(todos[i], isSelected, width-4)
		todoLines = append(todoLines, todoLine)
	}

	content := strings.Join(todoLines, "\n")

	return v.contentStyle.Width(width).Height(height).Render(content)
}

// CalculateVisibleRange calculates which todos should be visible based on cursor position
func (v *ListView) CalculateVisibleRange(cursor, totalItems, visibleHeight int) (int, int) {
	if totalItems <= visibleHeight {
		return 0, totalItems
	}

	// Keep cursor in the middle of the visible area when possible
	halfHeight := visibleHeight / 2

	startIndex := cursor - halfHeight
	if startIndex < 0 {
		startIndex = 0
	}

	endIndex := startIndex + visibleHeight
	if endIndex > totalItems {
		endIndex = totalItems
		startIndex = endIndex - visibleHeight
		if startIndex < 0 {
			startIndex = 0
		}
	}

	return startIndex, endIndex
}

// RenderScrollIndicator renders scroll indicators when content is scrolled
func (v *ListView) RenderScrollIndicator(stateModel *model.StateModel, width int) string {
	todos := stateModel.Todos()
	cursor := stateModel.Cursor()

	if len(todos) == 0 {
		return ""
	}

	// Calculate scroll position
	scrollPercent := float64(cursor) / float64(len(todos)-1)
	if len(todos) == 1 {
		scrollPercent = 0
	}

	// Create scroll bar
	barWidth := width - 4
	if barWidth < 10 {
		barWidth = 10
	}

	position := int(scrollPercent * float64(barWidth-1))

	scrollBar := strings.Repeat("─", barWidth)
	if position >= 0 && position < len(scrollBar) {
		scrollBar = scrollBar[:position] + "●" + scrollBar[position+1:]
	}

	return v.lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(scrollBar)
}

// GetVisibleItemCount returns the number of items that can be displayed in the given height
func (v *ListView) GetVisibleItemCount(height int) int {
	return height - 2 // Account for padding
}

// GetTotalHeight returns the total height needed to display all items
func (v *ListView) GetTotalHeight(stateModel *model.StateModel) int {
	todos := stateModel.Todos()
	totalHeight := 0

	for _, todo := range todos {
		totalHeight += v.itemView.GetItemHeight(todo)
	}

	return totalHeight + 2 // Account for padding
}

// IsScrollable returns whether the list needs scrolling
func (v *ListView) IsScrollable(stateModel *model.StateModel, height int) bool {
	return v.GetTotalHeight(stateModel) > height
}

// SetContentStyle allows customization of the content style
func (v *ListView) SetContentStyle(style proxy.Style) {
	v.contentStyle = style
}

// SetEmptyStyle allows customization of the empty state style
func (v *ListView) SetEmptyStyle(style proxy.Style) {
	v.emptyStyle = style
}
