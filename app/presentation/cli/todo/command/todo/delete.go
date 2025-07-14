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

func NewDeleteCommand(
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
	cmd.SetUse("delete [id]")
	cmd.SetShort("Delete a todo")
	cmd.SetArgs(cobra.ExactArgs(1))
	cmd.PersistentFlags().StringVarP(
		&format,
		"format",
		"f",
		conf.OutputFormat,
		"Output format (text|json)",
	)
	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			todoRepo, err := todoRepo.NewTodoRepository(
				conf,
				fileutil,
				json,
				os,
			)
			if err != nil {
				return err
			}

			uc := todoApp.NewDeleteTodoUseCase(todoRepo)
			dto, err := uc.Run(args[0])
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
