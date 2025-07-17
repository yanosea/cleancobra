package view

import (
	"fmt"
	"strings"

	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	"github.com/charmbracelet/lipgloss"
)

// HeaderView handles the application header rendering
type HeaderView struct {
	lipgloss    proxy.Lipgloss
	headerStyle proxy.Style
}

// NewHeaderView creates a new HeaderView with the given lipgloss proxy
func NewHeaderView(lg proxy.Lipgloss) *HeaderView {
	return &HeaderView{
		lipgloss:    lg,
		headerStyle: lg.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Padding(0, 1),
	}
}

// Render renders the application header
func (v *HeaderView) Render(stateModel *model.StateModel, width int) string {
	title := "GCT - Go Clean-Architecture Todo"
	todoCount := len(stateModel.Todos())
	completedCount := v.GetCompletedCount(stateModel.Todos())
	
	status := fmt.Sprintf("(%d/%d todos)", completedCount, todoCount)
	
	// Create header with title and status
	paddingLength := width - len(title) - len(status) - 4
	if paddingLength < 0 {
		paddingLength = 0
	}
	headerContent := v.lipgloss.JoinHorizontal(
		v.lipgloss.Left(),
		title,
		strings.Repeat(" ", paddingLength), // Padding
		status,
	)
	
	return v.headerStyle.Width(width).Render(headerContent)
}

// RenderCompact renders a compact header for smaller terminals
func (v *HeaderView) RenderCompact(stateModel *model.StateModel, width int) string {
	todoCount := len(stateModel.Todos())
	completedCount := v.GetCompletedCount(stateModel.Todos())
	
	title := "GCT"
	status := fmt.Sprintf("(%d/%d)", completedCount, todoCount)
	
	// Create compact header
	paddingLength := width - len(title) - len(status) - 4
	if paddingLength < 0 {
		paddingLength = 0
	}
	headerContent := v.lipgloss.JoinHorizontal(
		v.lipgloss.Left(),
		title,
		strings.Repeat(" ", paddingLength),
		status,
	)
	
	return v.headerStyle.Width(width).Render(headerContent)
}

// RenderWithMode renders header with current mode indication
func (v *HeaderView) RenderWithMode(stateModel *model.StateModel, width int) string {
	title := "GCT"
	todoCount := len(stateModel.Todos())
	completedCount := v.GetCompletedCount(stateModel.Todos())
	
	// Add mode indicator
	modeIndicator := v.getModeIndicator(stateModel.Mode())
	status := fmt.Sprintf("%s (%d/%d)", modeIndicator, completedCount, todoCount)
	
	// Create header with mode and status
	paddingLength := width - len(title) - len(status) - 4
	if paddingLength < 0 {
		paddingLength = 0
	}
	headerContent := v.lipgloss.JoinHorizontal(
		v.lipgloss.Left(),
		title,
		strings.Repeat(" ", paddingLength),
		status,
	)
	
	return v.headerStyle.Width(width).Render(headerContent)
}

// GetCompletedCount returns the number of completed todos
func (v *HeaderView) GetCompletedCount(todos []*model.ItemModel) int {
	count := 0
	for _, todo := range todos {
		if todo.Todo() != nil && todo.Todo().Done {
			count++
		}
	}
	return count
}

// getModeIndicator returns a visual indicator for the current mode
func (v *HeaderView) getModeIndicator(mode model.Mode) string {
	switch mode {
	case model.ModeInput:
		return "ADD"
	case model.ModeEdit:
		return "EDIT"
	case model.ModeConfirmation:
		return "CONFIRM"
	default:
		return "NORMAL"
	}
}

// GetHeight returns the height needed for the header
func (v *HeaderView) GetHeight() int {
	return 1
}

// SetStyle allows customization of the header style
func (v *HeaderView) SetStyle(style proxy.Style) {
	v.headerStyle = style
}