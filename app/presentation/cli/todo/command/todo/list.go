package todo

import (
	c "github.com/spf13/cobra"

	todoApp "github.com/yanosea/cleancobra/app/application/todo"
	"github.com/yanosea/cleancobra/app/config"
	todoRepo "github.com/yanosea/cleancobra/app/infrastructure/json/repository"
	"github.com/yanosea/cleancobra/app/presentation/cli/todo/formatter"

	"github.com/yanosea/cleancobra/pkg/proxy"
	"github.com/yanosea/cleancobra/pkg/utility"
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
		func(cmd *c.Command, _ []string) error {
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
		},
	)

	return cmd
}
