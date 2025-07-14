package formatter

import (
	"fmt"

	"github.com/yanosea/gct/pkg/errors"
	"github.com/yanosea/gct/pkg/proxy"
)

type Formatter interface {
	Format(result interface{}) (string, error)
}

func NewFormatter(
	format string,
	json proxy.Json,
) (Formatter, error) {
	var f Formatter
	switch format {
	case "json":
		f = NewJSONFormatter(json)
	case "text":
		f = NewTextFormatter()
	default:
		return nil, errors.New("invalid format")
	}
	return f, nil
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
		result = Red(result)
	}

	return result
}
