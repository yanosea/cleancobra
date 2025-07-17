package commands

import (
	"github.com/spf13/cobra"
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/presentation/cli/gct/commands/gct"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

// NewRootCommand creates the root command for the gct CLI application
func NewRootCommand(
	cobraProxy proxy.Cobra,
	listUseCase *application.ListTodoUseCase,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("gct")
	cmd.SetShort("A clean architecture todo application")
	cmd.SetLong(`gct is a todo application built with clean architecture principles.
It provides both CLI and TUI interfaces for managing your todos.

When run without any subcommands, it will list all todos.`)
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	
	// Add global flags
	formatFlag := cmd.PersistentFlags().StringP("format", "f", "table", "Output format (json, table, plain)")
	
	// Set the default run function to execute list command
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		// Get format from flag
		format := "table"
		if formatFlag != nil {
			format = *formatFlag
		}
		return runRoot(listUseCase, presenter, format)
	})
	
	return cmd
}

// runRoot executes the root command functionality (delegates to list command)
func runRoot(listUseCase *application.ListTodoUseCase, presenter *presenter.TodoPresenter, format string) error {
	// Call the actual list command implementation from gct/list.go
	return gct.RunList(listUseCase, presenter, format)
}