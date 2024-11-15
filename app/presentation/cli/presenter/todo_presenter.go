package presenter

import (
	"fmt"
	"github.com/yanosea/cleancobra/presentation/cli/formatter"

	"github.com/yanosea/cleancobra-pkg/errors"
)

type Presenter interface {
	Present(result interface{}) error
}

type todoPresenter struct {
	formatter formatter.Formatter
}

func (p *todoPresenter) Present(result interface{}) error {
	output, err := p.formatter.Format(result)
	if err != nil {
		return err
	}
	if _, err := fmt.Println(output); err != nil {
		return err
	}
	return nil
}

func NewPresenter(format string) (Presenter, error) {
	var f formatter.Formatter
	switch format {
	case "json":
		f = formatter.NewJSONFormatter()
	case "text":
		f = formatter.NewTextFormatter()
	default:
		return nil, errors.New("invalid format")
	}
	return &todoPresenter{formatter: f}, nil
}
