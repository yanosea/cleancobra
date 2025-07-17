package formatter

import (
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
)

// TableFormatter formats todos as human-readable table output
type TableFormatter struct {
	colorProxy   proxy.Color
	stringsProxy proxy.Strings
	fmtProxy     proxy.Fmt
}

// NewTableFormatter creates a new TableFormatter instance
func NewTableFormatter(colorProxy proxy.Color, stringsProxy proxy.Strings, fmtProxy proxy.Fmt) *TableFormatter {
	return &TableFormatter{
		colorProxy:   colorProxy,
		stringsProxy: stringsProxy,
		fmtProxy:     fmtProxy,
	}
}

// Format formats a slice of todos as a human-readable table string
func (f *TableFormatter) Format(todos []domain.Todo) (string, error) {
	if len(todos) == 0 {
		return f.colorProxy.Yellow("No todos found."), nil
	}

	var result []string
	
	// Add header
	header := f.fmtProxy.Sprintf("%-4s %-6s %s", 
		f.colorProxy.Cyan("ID"), 
		f.colorProxy.Cyan("Status"), 
		f.colorProxy.Cyan("Description"))
	result = append(result, header)
	
	// Add separator line
	separator := f.fmtProxy.Sprintf("%-4s %-6s %s", "----", "------", "-----------")
	result = append(result, separator)

	// Add todos
	for _, todo := range todos {
		line := f.formatTodoLine(todo)
		result = append(result, line)
	}

	return f.stringsProxy.Join(result, "\n"), nil
}

// FormatSingle formats a single todo as a human-readable table string
func (f *TableFormatter) FormatSingle(todo domain.Todo) (string, error) {
	var result []string
	
	// Add header
	header := f.fmtProxy.Sprintf("%-4s %-6s %s", 
		f.colorProxy.Cyan("ID"), 
		f.colorProxy.Cyan("Status"), 
		f.colorProxy.Cyan("Description"))
	result = append(result, header)
	
	// Add separator line
	separator := f.fmtProxy.Sprintf("%-4s %-6s %s", "----", "------", "-----------")
	result = append(result, separator)

	// Add the todo
	line := f.formatTodoLine(todo)
	result = append(result, line)

	return f.stringsProxy.Join(result, "\n"), nil
}

// formatTodoLine formats a single todo as a table row
func (f *TableFormatter) formatTodoLine(todo domain.Todo) string {
	var status string
	var description string
	
	if todo.Done {
		status = f.colorProxy.Green("✓ Done")
		description = f.colorProxy.Green(todo.Description)
	} else {
		status = f.colorProxy.Red("✗ Todo")
		description = todo.Description
	}

	return f.fmtProxy.Sprintf("%-4d %-6s %s", todo.ID, status, description)
}