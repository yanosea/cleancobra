package todo

import (
	todoApp "cleancobra/app/application/usecase/todo"
	"github.com/spf13/cobra"
)

func NewAddCommand(tu todoApp.AddTodoUseCase) *cobra.Command {
	return &cobra.Command{
		Use:   "add [title]",
		Short: "Add a new todo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tu.Run(args[0])
		},
	}
}
