package formatter

import (
	"github.com/yanosea/gct/app/domain"

	"github.com/yanosea/gct/pkg/proxy"
)

// PlainFormatter formats todos as simple plain text output
type PlainFormatter struct {
	fmtProxy     proxy.Fmt
	strconvProxy proxy.Strconv
	stringsProxy proxy.Strings
}

// NewPlainFormatter creates a new PlainFormatter instance
func NewPlainFormatter(
	fmtProxy proxy.Fmt,
	strconvProxy proxy.Strconv,
	stringsProxy proxy.Strings,
) *PlainFormatter {
	return &PlainFormatter{
		fmtProxy:     fmtProxy,
		strconvProxy: strconvProxy,
		stringsProxy: stringsProxy,
	}
}

// Format formats a slice of todos as plain text string
func (f *PlainFormatter) Format(todos []domain.Todo) (string, error) {
	if len(todos) == 0 {
		return "No todos found.", nil
	}

	var lines []string

	for _, todo := range todos {
		line := f.formatTodoLine(todo)
		lines = append(lines, line)
	}

	result := f.stringsProxy.Join(lines, "\n")
	return f.stringsProxy.TrimSpace(result), nil
}

// formatTodoLine formats a single todo as a plain text line
func (f *PlainFormatter) formatTodoLine(todo domain.Todo) string {
	// Status indicator
	status := "[ ]"
	if todo.Done {
		status = "[x]"
	}

	// Format: [x] 1: Buy groceries
	return f.fmtProxy.Sprintf("%s %s: %s",
		status,
		f.strconvProxy.Itoa(todo.ID),
		todo.Description,
	)
}
