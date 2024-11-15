package todo

import (
	todoApp "cleancobra/app/application/usecase/todo"
	"github.com/spf13/cobra"
)

func NewDeleteCommand(tu todoApp.DeleteTodoUseCase) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a todo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tu.Run(args[0])
		},
	}
}
