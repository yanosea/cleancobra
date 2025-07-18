package formatter

import (
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
)

// JSONFormatter formats todos as JSON output
type JSONFormatter struct {
	jsonProxy proxy.JSON
}

// NewJSONFormatter creates a new JSONFormatter instance
func NewJSONFormatter(jsonProxy proxy.JSON) *JSONFormatter {
	return &JSONFormatter{
		jsonProxy: jsonProxy,
	}
}

// Format formats a slice of todos as JSON string
func (f *JSONFormatter) Format(todos []domain.Todo) (string, error) {
	data, err := f.jsonProxy.Marshal(todos)
	if err != nil {
		return "", domain.NewDomainError(
			domain.ErrorTypeJSON,
			"failed to marshal todos to JSON",
			err,
		)
	}

	return string(data), nil
}

// FormatSingle formats a single todo as JSON string
func (f *JSONFormatter) FormatSingle(todo domain.Todo) (string, error) {
	data, err := f.jsonProxy.Marshal(todo)
	if err != nil {
		return "", domain.NewDomainError(
			domain.ErrorTypeJSON,
			"failed to marshal todo to JSON",
			err,
		)
	}

	return string(data), nil
}
