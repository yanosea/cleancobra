package formatter

import (
	"cleancobra/app/presentation/cli/adapter/dto"
	"encoding/json"
)

type JSONFormatter struct{}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (f *JSONFormatter) Format(todos []dto.Todo) (string, error) {
	output, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}
