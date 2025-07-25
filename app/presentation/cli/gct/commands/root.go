package commands

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/cli/gct/commands/gct"
	"github.com/yanosea/gct/app/presentation/cli/gct/commands/gct/completion"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewRootCommand creates the root command for the gct CLI application
func NewRootCommand(
	container *container.Container,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := container.GetProxies().Cobra.NewCommand()
	cmd.SetUse(rootUse)
	cmd.SetShort(rootShort)
	cmd.SetLong(rootLong)
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	cmd.AddCommand(gct.NewAddCommand(
		container,
		presenter,
	))
	listCmd := gct.NewListCommand(
		container,
		presenter,
	)
	cmd.AddCommand(listCmd)
	cmd.AddCommand(gct.NewToggleCommand(
		container,
		presenter,
	))
	cmd.AddCommand(gct.NewDeleteCommand(
		container,
		presenter,
	))
	cmd.AddCommand(completion.NewCompletionCommand(
		cmd,
		container,
		presenter,
	))
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		return runRoot(listCmd, args)
	})

	return cmd
}

// runRoot executes the root command functionality (delegates to list command)
func runRoot(listCmd proxy.Command, args []string) error {
	// root command delegates to list command
	listCmd.RunE(nil, args)

	return nil
}

const (
	rootUse   string = "gct"
	rootShort string = "✅ A clean architecture todo application"
	rootLong  string = `✅ gct is a todo application built with clean architecture principles.
It provides both CLI and TUI interfaces for managing your todos.

When run without any subcommands, it will list all todos.`
)
