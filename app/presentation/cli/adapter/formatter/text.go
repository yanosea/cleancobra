package formatter

import (
	"cleancobra/app/presentation/cli/adapter/dto"
	"fmt"
	"strings"
)

type TextFormatter struct{}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

func (f *TextFormatter) Format(todos []dto.Todo) (string, error) {
	if len(todos) == 0 {
		return "No todos found", nil
	}

	var result strings.Builder
	for _, todo := range todos {
		status := "[ ]"
		if todo.Done {
			status = Green("[✓]")
		}
		result.WriteString(fmt.Sprintf("%s %s (ID: %s)\n", status, todo.Title, todo.ID))
	}
	return result.String(), nil
}
