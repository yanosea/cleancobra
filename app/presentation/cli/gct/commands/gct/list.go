package gct

import (
	"github.com/spf13/cobra"
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

// NewListCommand creates the list command for the gct CLI application
func NewListCommand(
	cobraProxy proxy.Cobra,
	listUseCase *application.ListTodoUseCase,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("list")
	cmd.SetShort("List all todos")
	cmd.SetLong(`List all todos with their current status.

Supports multiple output formats:
- table: Human-readable table format (default)
- json: JSON format for programmatic use
- plain: Simple plain text format`)
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)

	// Add format flag
	formatFlag := cmd.Flags().StringP("format", "f", "table", "Output format (json, table, plain)")

	// Set the run function
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		// Get format from flag
		format := "table"
		if formatFlag != nil {
			format = *formatFlag
		}
		return RunList(listUseCase, presenter, format)
	})

	return cmd
}

// RunList executes the list todos functionality
func RunList(listUseCase *application.ListTodoUseCase, presenter *presenter.TodoPresenter, format string) error {
	todos, err := listUseCase.Run()
	if err != nil {
		presenter.ShowError(err)
		return err
	}

	return presenter.ShowListResults(todos, format)
}
