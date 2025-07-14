package formatter

import (
	"github.com/yanosea/gct/pkg/proxy"
)

type JSONFormatter struct {
	json proxy.Json
}

func NewJSONFormatter(json proxy.Json) *JSONFormatter {
	return &JSONFormatter{
		json: json,
	}
}

func (f *JSONFormatter) Format(result interface{}) (string, error) {
	output, err := f.json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}
