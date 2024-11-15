package presenter

import (
	"cleancobra/app/presentation/cli/adapter/dto"
	"cleancobra/app/presentation/cli/adapter/formatter"
	"fmt"
)

type Presenter interface {
	Present(todos []dto.Todo) error
}

type todoPresenter struct {
	formatter interface {
		Format(todos []dto.Todo) (string, error)
	}
}

func (p *todoPresenter) Present(todos []dto.Todo) error {
	output, err := p.formatter.Format(todos)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}

func NewPresenter(format string) Presenter {
	switch format {
	case "json":
		return &todoPresenter{formatter: formatter.NewJSONFormatter()}
	case "text":
		return &todoPresenter{formatter: formatter.NewTextFormatter()}
	default:
		return &todoPresenter{formatter: formatter.NewTextFormatter()}
	}
}
