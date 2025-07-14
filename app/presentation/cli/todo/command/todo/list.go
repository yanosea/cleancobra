package todo

import (
	c "github.com/spf13/cobra"

	todoApp "github.com/yanosea/gct/app/application/todo"
	"github.com/yanosea/gct/app/config"
	todoRepo "github.com/yanosea/gct/app/infrastructure/json/repository"
	"github.com/yanosea/gct/app/presentation/cli/todo/formatter"

	"github.com/yanosea/gct/pkg/proxy"
	"github.com/yanosea/gct/pkg/utility"
)

func NewListCommand(
	cobra proxy.Cobra,
	json proxy.Json,
	os proxy.Os,
	fileutil utility.FileUtil,
	conf *config.TodoConfig,
	output *string,
) proxy.Command {
	var format = conf.OutputFormat
	cmd := cobra.NewCommand()
	cmd.SetSilenceErrors(true)
	cmd.SetUse("list")
	cmd.SetShort("List all todos")
	cmd.PersistentFlags().StringVarP(
		&format,
		"format",
		"f",
		conf.OutputFormat,
		"Output format (text|json)",
	)
	cmd.SetRunE(
		func(_ *c.Command, _ []string) error {
			return runList(format, json, os, fileutil, conf, output)
		},
	)

	return cmd
}

func runList(
	format string,
	json proxy.Json,
	os proxy.Os,
	fileutil utility.FileUtil,
	conf *config.TodoConfig,
	output *string,
) error {
	todoRepo, err := todoRepo.NewTodoRepository(
		conf,
		fileutil,
		json,
		os,
	)
	if err != nil {
		return err
	}

	uc := todoApp.NewListTodoUseCase(todoRepo)
	dto, err := uc.Run()
	if err != nil {
		return err
	}

	f, err := formatter.NewFormatter(format, json)
	if err != nil {
		return err
	}

	o, err := f.Format(dto)
	if err != nil {
		return err
	}

	*output = o

	return nil
}
