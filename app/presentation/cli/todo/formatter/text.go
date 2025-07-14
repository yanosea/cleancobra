package formatter

import (
	"fmt"
	"strings"

	todoApp "github.com/yanosea/gct/app/application/todo"

	"github.com/yanosea/gct/pkg/errors"
)

type TextFormatter struct{}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

func (f *TextFormatter) Format(result interface{}) (string, error) {
	switch v := result.(type) {
	case *todoApp.AddTodoUsecaseOutputDto:
		return fmt.Sprintf("Added todo : %s (ID: %s, CREATED AT: %s)", v.Title, v.ID, v.CreatedAt), nil
	case *todoApp.DeleteTodoUsecaseOutputDto:
		status := "[ ]"
		if v.Done {
			status = Green("[✓]")
		}
		return fmt.Sprintf("Deleted todo : %s %s (ID: %s, CREATED AT: %s)", status, v.Title, v.ID, v.CreatedAt), nil
	case *todoApp.ToggleTodoUsecaseOutputDto:
		status := "[ ]"
		if v.Done {
			status = Green("[✓]")
		}
		return fmt.Sprintf("Toggled todo : %s %s (ID: %s, CREATED AT: %s)", status, v.Title, v.ID, v.CreatedAt), nil
	case []*todoApp.ListTodoUsecaseOutputDto:
		if len(v) == 0 {
			return "No todos found", nil
		}
		var result = strings.Builder{}
		for i, todo := range v {
			status := "[ ]"
			if todo.Done {
				status = Green("[✓]")
			}
			result.WriteString(fmt.Sprintf("%s %s (ID: %s, CREATED AT: %s)", status, todo.Title, todo.ID, todo.CreatedAt))
			if i != len(v)-1 {
				result.WriteString("\n")
			}
		}
		return result.String(), nil
	default:
		return "", errors.New("unsupported result type")
	}
}
