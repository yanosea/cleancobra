package gct

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewListCommand creates the list command for the gct CLI application
func NewListCommand(
	container *container.Container,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := container.GetProxies().Cobra.NewCommand()
	cmd.SetUse(listUse)
	cmd.SetShort(listShort)
	cmd.SetLong(listLong)
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	format := cmd.Flags().StringP(
		"format",
		"f",
		"table",
		"Output format (json, table, plain)",
	)
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		return runList(container.GetUseCases().ListTodo, presenter, *format)
	})

	return cmd
}

// RunList executes the list todos functionality
func runList(
	listUseCase *application.ListTodoUseCase,
	presenter *presenter.TodoPresenter,
	format string,
) error {
	todos, err := listUseCase.Run()
	if err != nil {
		presenter.ShowError(err)
		return err
	}

	return presenter.ShowListResults(todos, format)
}

const (
	listUse   string = "list"
	listShort string = "✅📜 List all todos"
	listLong  string = `✅📜 List all todos with their current status.

Supports multiple output formats:
  - table: Human-readable table format (default)
  - json:  JSON format for programmatic use
  - plain: Simple plain text format`
)
