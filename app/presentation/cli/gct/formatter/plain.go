package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yanosea/gct/app/domain"
)

// PlainFormatter formats todos as simple plain text output
type PlainFormatter struct{}

// NewPlainFormatter creates a new PlainFormatter instance
func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

// Format formats a slice of todos as plain text string
func (f *PlainFormatter) Format(todos []domain.Todo) (string, error) {
	if len(todos) == 0 {
		return "No todos found.", nil
	}

	var result strings.Builder
	
	for _, todo := range todos {
		line := f.formatTodoLine(todo)
		result.WriteString(line)
		result.WriteString("\n")
	}
	
	return strings.TrimSpace(result.String()), nil
}

// FormatSingle formats a single todo as plain text string
func (f *PlainFormatter) FormatSingle(todo domain.Todo) (string, error) {
	return f.formatTodoLine(todo), nil
}

// formatTodoLine formats a single todo as a plain text line
func (f *PlainFormatter) formatTodoLine(todo domain.Todo) string {
	// Status indicator
	status := "[ ]"
	if todo.Done {
		status = "[x]"
	}
	
	// Format: [x] 1: Buy groceries
	return fmt.Sprintf("%s %s: %s", 
		status, 
		strconv.Itoa(todo.ID), 
		todo.Description,
	)
}