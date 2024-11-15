package todo

import (
	"github.com/spf13/cobra"

	todoApp "github.com/yanosea/cleancobra/application/todo"
	"github.com/yanosea/cleancobra/config"
	todoRepo "github.com/yanosea/cleancobra/infrastructure/json/repository"
	"github.com/yanosea/cleancobra/presentation/cli/presenter"
)

func NewListCommand(conf *config.TodoConfig) *cobra.Command {
	var format = conf.OutputFormat
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		RunE: func(cmd *cobra.Command, args []string) error {
			todoRepo, err := todoRepo.NewTodoRepository(conf)
			if err != nil {
				return err
			}

			uc := todoApp.NewListTodoUseCase(todoRepo)
			dto, err := uc.Run()
			if err != nil {
				return err
			}

			p, err := presenter.NewPresenter(format)
			if err != nil {
				return err
			}

			return p.Present(dto)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", conf.OutputFormat, "Output format (text|json)")
	return cmd
}
