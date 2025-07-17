package view

import (
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	"github.com/charmbracelet/lipgloss"
)

// FooterView handles the application footer rendering
type FooterView struct {
	lipgloss    proxy.Lipgloss
	footerStyle proxy.Style
	helpStyle   proxy.Style
}

// NewFooterView creates a new FooterView with the given lipgloss proxy
func NewFooterView(lg proxy.Lipgloss) *FooterView {
	return &FooterView{
		lipgloss:    lg,
		footerStyle: lg.NewStyle().Foreground(lipgloss.Color("8")).Padding(0, 1),
		helpStyle:   lg.NewStyle().Foreground(lipgloss.Color("8")).Italic(true),
	}
}

// Render renders the application footer with help text
func (v *FooterView) Render(stateModel *model.StateModel, width int) string {
	helpText := v.getHelpText(stateModel.Mode())
	return v.footerStyle.Width(width).Render(helpText)
}

// RenderCompact renders a compact footer for smaller terminals
func (v *FooterView) RenderCompact(stateModel *model.StateModel, width int) string {
	helpText := v.getCompactHelpText(stateModel.Mode())
	return v.footerStyle.Width(width).Render(helpText)
}

// RenderWithScrollIndicator renders footer with scroll indicator
func (v *FooterView) RenderWithScrollIndicator(stateModel *model.StateModel, width int, scrollIndicator string) string {
	helpText := v.getHelpText(stateModel.Mode())
	
	// Combine help text with scroll indicator
	content := v.lipgloss.JoinVertical(
		v.lipgloss.Left(),
		scrollIndicator,
		helpText,
	)
	
	return v.footerStyle.Width(width).Render(content)
}

// getHelpText returns appropriate help text based on current mode
func (v *FooterView) getHelpText(mode model.Mode) string {
	switch mode {
	case model.ModeInput:
		return "Enter: Add • Esc: Cancel"
	case model.ModeEdit:
		return "Enter: Save • Esc: Cancel"
	case model.ModeConfirmation:
		return "y: Yes • n/Esc: No"
	default:
		return "↑/k: Up • ↓/j: Down • g: Top • G: Bottom • Space: Toggle • a: Add • e: Edit • d: Delete • r: Refresh • q: Quit"
	}
}

// getCompactHelpText returns compact help text for smaller terminals
func (v *FooterView) getCompactHelpText(mode model.Mode) string {
	switch mode {
	case model.ModeInput:
		return "Enter:Add Esc:Cancel"
	case model.ModeEdit:
		return "Enter:Save Esc:Cancel"
	case model.ModeConfirmation:
		return "y:Yes n:No"
	default:
		return "↑↓:Move Space:Toggle a:Add e:Edit d:Del q:Quit"
	}
}

// RenderModeSpecificHelp renders help text specific to the current context
func (v *FooterView) RenderModeSpecificHelp(stateModel *model.StateModel, width int) string {
	var helpLines []string
	
	switch stateModel.Mode() {
	case model.ModeNormal:
		helpLines = []string{
			"Navigation: ↑/k (up), ↓/j (down), g (top), G (bottom)",
			"Actions: Space (toggle), a (add), e (edit), d (delete), r (refresh)",
			"Other: q (quit), Esc (clear error)",
		}
	case model.ModeInput:
		helpLines = []string{
			"Type your todo description and press Enter to add",
			"Press Esc to cancel without adding",
		}
	case model.ModeEdit:
		helpLines = []string{
			"Edit the todo description and press Enter to save",
			"Press Esc to cancel without saving changes",
		}
	case model.ModeConfirmation:
		helpLines = []string{
			"Press 'y' or 'Y' to confirm the action",
			"Press 'n', 'N', or Esc to cancel",
		}
	}
	
	content := ""
	for i, line := range helpLines {
		if i > 0 {
			content += "\n"
		}
		content += v.helpStyle.Render(line)
	}
	
	return v.footerStyle.Width(width).Render(content)
}

// GetHeight returns the height needed for the footer
func (v *FooterView) GetHeight() int {
	return 1
}

// GetExpandedHeight returns the height needed for expanded help
func (v *FooterView) GetExpandedHeight(mode model.Mode) int {
	switch mode {
	case model.ModeNormal:
		return 3
	default:
		return 2
	}
}

// SetStyle allows customization of the footer style
func (v *FooterView) SetStyle(style proxy.Style) {
	v.footerStyle = style
}

// SetHelpStyle allows customization of the help text style
func (v *FooterView) SetHelpStyle(style proxy.Style) {
	v.helpStyle = style
}