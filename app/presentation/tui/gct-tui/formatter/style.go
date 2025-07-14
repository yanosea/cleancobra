package formatter

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Align(lipgloss.Center).
			Padding(0, 1)

	TodoItemStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	CompletedTodoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Strikethrough(true).
				Padding(0, 1).
				Margin(0, 0, 1, 0)

	CheckboxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	UncheckboxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("205")).
				Padding(0, 1)

	ButtonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 2).
			Margin(0, 1)

	ActiveButtonStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("205")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 2).
				Margin(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Margin(1, 0)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	WarningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2).
			Margin(1, 0).
			Align(lipgloss.Center)

	DangerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	HighlightedTodoStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("52")).
				Foreground(lipgloss.Color("255")).
				Padding(0, 1).
				Margin(0, 0, 1, 0).
				Bold(true)

	ConfirmButtonStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("196")).
				Foreground(lipgloss.Color("255")).
				Padding(0, 3).
				Margin(0, 1).
				Bold(true)

	CancelButtonStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("240")).
				Foreground(lipgloss.Color("255")).
				Padding(0, 3).
				Margin(0, 1)

	SelectedConfirmButtonStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("9")).
					Foreground(lipgloss.Color("255")).
					Padding(0, 3).
					Margin(0, 1).
					Bold(true).
					Blink(true)

	SelectedCancelButtonStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("62")).
					Foreground(lipgloss.Color("255")).
					Padding(0, 3).
					Margin(0, 1).
					Bold(true)
)

func FormatTodoItem(title string, done bool, selected bool) string {
	var checkbox string
	style := TodoItemStyle

	if done {
		checkbox = CheckboxStyle.Render("[✓]")
		style = CompletedTodoStyle
	} else {
		checkbox = UncheckboxStyle.Render("[ ]")
	}

	text := checkbox + " " + title

	if selected {
		style = style.Background(lipgloss.Color("62"))
	}

	return style.Render(text)
}

func FormatHeader(text string) string {
	return HeaderStyle.Render(text)
}

func FormatError(text string) string {
	return ErrorStyle.Render("Error: " + text)
}

func FormatSuccess(text string) string {
	return SuccessStyle.Render(text)
}

func FormatHelp(text string) string {
	return HelpStyle.Render(text)
}

func FormatInput(text string) string {
	return FocusedInputStyle.Render(text + "█")
}

func FormatWarningBox(content string) string {
	return WarningBoxStyle.Render(content)
}

func FormatDanger(text string) string {
	return DangerStyle.Render(text)
}

func FormatHighlightedTodo(title string, done bool) string {
	checkbox := "[ ]"
	if done {
		checkbox = "[✓]"
	}
	text := checkbox + " " + title
	return HighlightedTodoStyle.Render(text)
}

func FormatConfirmButton(text string, selected bool) string {
	if selected {
		return SelectedConfirmButtonStyle.Render(text)
	}
	return ConfirmButtonStyle.Render(text)
}

func FormatCancelButton(text string, selected bool) string {
	if selected {
		return SelectedCancelButtonStyle.Render(text)
	}
	return CancelButtonStyle.Render(text)
}
