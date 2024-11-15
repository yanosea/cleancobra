package todo

import (
	todoApp "cleancobra/app/application/usecase/todo"
	"github.com/spf13/cobra"
)

func NewToggleCommand(tu todoApp.ToggleTodoUseCase) *cobra.Command {
	return &cobra.Command{
		Use:   "toggle [id]",
		Short: "Toggle todo status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tu.Run(args[0])
		},
	}
}
