package view

import (
	"fmt"
	"strings"

	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	"github.com/charmbracelet/lipgloss"
)

// LayoutView handles the main application layout rendering
type LayoutView struct {
	lipgloss   proxy.Lipgloss
	headerView *HeaderView
	footerView *FooterView
	listView   *ListView
	itemView   *ItemView
	
	// Layout styles
	contentStyle  proxy.Style
	inputStyle    proxy.Style
	errorStyle    proxy.Style
	confirmStyle  proxy.Style
	helpStyle     proxy.Style
}

// NewLayoutView creates a new LayoutView with the given dependencies
func NewLayoutView(lg proxy.Lipgloss, itemView *ItemView) *LayoutView {
	headerView := NewHeaderView(lg)
	footerView := NewFooterView(lg)
	listView := NewListView(lg, itemView)
	
	return &LayoutView{
		lipgloss:     lg,
		headerView:   headerView,
		footerView:   footerView,
		listView:     listView,
		itemView:     itemView,
		contentStyle: lg.NewStyle().Padding(1, 2),
		inputStyle:   lg.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(0, 1),
		errorStyle:   lg.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Padding(0, 1),
		confirmStyle: lg.NewStyle().Foreground(lipgloss.Color("3")).Bold(true).Padding(0, 1),
		helpStyle:    lg.NewStyle().Foreground(lipgloss.Color("8")).Italic(true),
	}
}

// Render renders the complete application layout
func (v *LayoutView) Render(stateModel *model.StateModel) string {
	width := stateModel.Width()
	height := stateModel.Height()
	
	// Build layout components using new component structure
	header := v.headerView.Render(stateModel, width)
	content := v.renderContent(stateModel, width, height-4) // Reserve space for header and footer
	footer := v.footerView.Render(stateModel, width)
	
	// Combine components vertically
	return v.lipgloss.JoinVertical(
		v.lipgloss.Left(),
		header,
		content,
		footer,
	)
}



// renderContent renders the main content area
func (v *LayoutView) renderContent(stateModel *model.StateModel, width, height int) string {
	switch stateModel.Mode() {
	case model.ModeInput:
		return v.renderInputMode(stateModel, width, height)
	case model.ModeEdit:
		return v.renderEditMode(stateModel, width, height)
	case model.ModeConfirmation:
		return v.renderConfirmationMode(stateModel, width, height)
	default:
		return v.renderNormalMode(stateModel, width, height)
	}
}

// renderNormalMode renders the normal todo list view
func (v *LayoutView) renderNormalMode(stateModel *model.StateModel, width, height int) string {
	// Use the new ListView component
	content := v.listView.Render(stateModel, width, height)
	
	// Add error message if present
	if errorMsg := stateModel.ErrorMessage(); errorMsg != "" {
		content += "\n\n" + v.errorStyle.Render("Error: "+errorMsg)
	}
	
	return content
}

// renderInputMode renders the add todo input mode
func (v *LayoutView) renderInputMode(stateModel *model.StateModel, width, height int) string {
	prompt := "Add new todo:"
	inputField := stateModel.Input().View()
	
	content := v.lipgloss.JoinVertical(
		v.lipgloss.Left(),
		prompt,
		"",
		v.inputStyle.Width(width-6).Render(inputField),
		"",
		v.helpStyle.Render("Press Enter to add, Esc to cancel"),
	)
	
	return v.contentStyle.
		Width(width).
		Height(height).
		AlignVertical(v.lipgloss.Center()).
		Render(content)
}

// renderEditMode renders the edit todo mode
func (v *LayoutView) renderEditMode(stateModel *model.StateModel, width, height int) string {
	prompt := "Edit todo:"
	inputField := stateModel.Input().View()
	
	// Show current todo being edited
	currentTodo := ""
	if cursor := stateModel.Cursor(); cursor >= 0 && cursor < len(stateModel.Todos()) {
		if todo := stateModel.Todos()[cursor].Todo(); todo != nil {
			currentTodo = fmt.Sprintf("Current: %s", todo.Description)
		}
	}
	
	content := v.lipgloss.JoinVertical(
		v.lipgloss.Left(),
		prompt,
		v.helpStyle.Render(currentTodo),
		"",
		v.inputStyle.Width(width-6).Render(inputField),
		"",
		v.helpStyle.Render("Press Enter to save, Esc to cancel"),
	)
	
	return v.contentStyle.
		Width(width).
		Height(height).
		AlignVertical(v.lipgloss.Center()).
		Render(content)
}

// renderConfirmationMode renders the confirmation dialog
func (v *LayoutView) renderConfirmationMode(stateModel *model.StateModel, width, height int) string {
	message := stateModel.ConfirmationMessage()
	
	content := v.lipgloss.JoinVertical(
		v.lipgloss.Center(),
		v.confirmStyle.Render(message),
		"",
		v.helpStyle.Render("Press 'y' to confirm, 'n' or Esc to cancel"),
	)
	
	return v.contentStyle.
		Width(width).
		Height(height).
		AlignHorizontal(v.lipgloss.Center()).
		AlignVertical(v.lipgloss.Center()).
		Render(content)
}



// RenderCompact renders a compact version of the layout for smaller terminals
func (v *LayoutView) RenderCompact(stateModel *model.StateModel) string {
	width := stateModel.Width()
	height := stateModel.Height()
	
	// Use new component structure for compact rendering
	header := v.headerView.RenderCompact(stateModel, width)
	
	// Compact content
	var content string
	switch stateModel.Mode() {
	case model.ModeInput:
		content = "Add: " + stateModel.Input().View()
	case model.ModeEdit:
		content = "Edit: " + stateModel.Input().View()
	case model.ModeConfirmation:
		content = stateModel.ConfirmationMessage() + " (y/n)"
	default:
		content = v.listView.RenderCompact(stateModel, width, height-2)
	}
	
	// Use new component structure for compact footer
	footer := v.footerView.RenderCompact(stateModel, width)
	
	return v.lipgloss.JoinVertical(
		v.lipgloss.Left(),
		header,
		content,
		footer,
	)
}



// GetMinimumSize returns the minimum terminal size needed for proper display
func (v *LayoutView) GetMinimumSize() (width, height int) {
	return 40, 10 // Minimum 40x10 for reasonable display
}

// IsCompactMode determines if compact mode should be used based on terminal size
func (v *LayoutView) IsCompactMode(width, height int) bool {
	minWidth, minHeight := v.GetMinimumSize()
	return width < minWidth || height < minHeight
}

// RenderScrollIndicator renders scroll indicators when content is scrolled
func (v *LayoutView) RenderScrollIndicator(stateModel *model.StateModel, width int) string {
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
	
	return v.helpStyle.Render(scrollBar)
}