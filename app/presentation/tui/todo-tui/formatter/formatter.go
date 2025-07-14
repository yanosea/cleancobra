package formatter

import (
	"fmt"

	todoApp "github.com/yanosea/gct/app/application/todo"
	"github.com/yanosea/gct/pkg/errors"
)

type Formatter interface {
	Format(result any) (string, error)
}

type TuiFormatter struct{}

func NewTuiFormatter() *TuiFormatter {
	return &TuiFormatter{}
}

func (f *TuiFormatter) Format(result any) (string, error) {
	switch v := result.(type) {
	case []*todoApp.ListTodoUsecaseOutputDto:
		return f.formatTodoList(v), nil
	case *todoApp.AddTodoUsecaseOutputDto:
		return f.formatAddResult(v), nil
	case *todoApp.DeleteTodoUsecaseOutputDto:
		return f.formatDeleteResult(v), nil
	case *todoApp.ToggleTodoUsecaseOutputDto:
		return f.formatToggleResult(v), nil
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", result), nil
	}
}

func (f *TuiFormatter) formatTodoList(todos []*todoApp.ListTodoUsecaseOutputDto) string {
	if len(todos) == 0 {
		return "No todos found."
	}

	var result string
	for _, todo := range todos {
		result += FormatTodoItem(todo.Title, todo.Done, false) + "\n"
	}
	return result
}

func (f *TuiFormatter) formatAddResult(result *todoApp.AddTodoUsecaseOutputDto) string {
	return FormatSuccess(fmt.Sprintf("Added todo: %s", result.Title))
}

func (f *TuiFormatter) formatDeleteResult(result *todoApp.DeleteTodoUsecaseOutputDto) string {
	return FormatSuccess(fmt.Sprintf("Deleted todo: %s", result.Title))
}

func (f *TuiFormatter) formatToggleResult(result *todoApp.ToggleTodoUsecaseOutputDto) string {
	status := "incomplete"
	if result.Done {
		status = "complete"
	}
	return FormatSuccess(fmt.Sprintf("Toggled todo: %s (now %s)", result.Title, status))
}

func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "tui", "styled":
		return NewTuiFormatter(), nil
	default:
		return nil, errors.New("invalid format for TUI")
	}
}

func AppendErrorToOutput(err error, output string) string {
	if err == nil && output == "" {
		return ""
	}

	var result string
	if err != nil {
		if output == "" {
			result = fmt.Sprintf("Error: %s", err)
		} else {
			result = fmt.Sprintf("Error: %s\n%s", err, output)
		}
	} else {
		result = output
	}

	if result != "" {
		result = FormatError(result)
	}

	return result
}
