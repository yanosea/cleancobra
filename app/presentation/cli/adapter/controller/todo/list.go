package todo

import (
	todoApp "cleancobra/app/application/usecase/todo"
	"cleancobra/app/presentation/cli/adapter/mapper"
	"cleancobra/app/presentation/cli/adapter/presenter"
	"github.com/spf13/cobra"
)

func NewListCommand(tu todoApp.ListTodoUseCase) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		RunE: func(cmd *cobra.Command, args []string) error {
			todos, err := tu.Run()
			if err != nil {
				return err
			}

			todoDTOs := mapper.ToDTO(todos)

			p := presenter.NewPresenter(format)
			return p.Present(todoDTOs)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text|json)")
	return cmd
}
