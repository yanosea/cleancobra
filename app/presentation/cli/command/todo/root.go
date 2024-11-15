package todo

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/cleancobra/config"
)

func NewTodoCommand(conf *config.TodoConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "todo",
		Short: "A clean architecture TODO application",
	}

	cmd.AddCommand(
		NewAddCommand(conf),
		NewDeleteCommand(conf),
		NewToggleCommand(conf),
		NewListCommand(conf),
	)

	return cmd
}
