package commands

import (
	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/cli/gct/commands/gct"
	"github.com/yanosea/gct/app/presentation/cli/gct/commands/gct/completion"
	"github.com/yanosea/gct/app/presentation/cli/gct/formatter"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

// InitializeCommand initializes the CLI application with all commands and dependencies
func InitializeCommand() (proxy.Command, error) {
	// Initialize dependency injection container
	c, err := container.NewContainer()
	if err != nil {
		return nil, err
	}

	// Initialize formatters
	jsonFormatter := formatter.NewJSONFormatter(c.GetProxies().JSON)
	tableFormatter := formatter.NewTableFormatter(c.GetProxies().Color, c.GetProxies().Strings, c.GetProxies().Fmt)
	plainFormatter := formatter.NewPlainFormatter()

	// Initialize presenter
	todoPresenter := presenter.NewTodoPresenter(
		jsonFormatter,
		tableFormatter,
		plainFormatter,
		c.GetProxies().Fmt,
		c.GetProxies().OS,
	)

	// Initialize root command
	rootCmd := NewRootCommand(
		c.GetProxies().Cobra,
		c.GetUseCases().ListTodo,
		todoPresenter,
	)

	// Add subcommands
	rootCmd.AddCommand(gct.NewAddCommand(
		c.GetProxies().Cobra,
		c.GetUseCases().AddTodo,
		todoPresenter,
	))

	rootCmd.AddCommand(gct.NewListCommand(
		c.GetProxies().Cobra,
		c.GetUseCases().ListTodo,
		todoPresenter,
	))

	rootCmd.AddCommand(gct.NewToggleCommand(
		c.GetProxies().Cobra,
		c.GetProxies().Strconv,
		c.GetUseCases().ToggleTodo,
		todoPresenter,
	))

	rootCmd.AddCommand(gct.NewDeleteCommand(
		c.GetProxies().Cobra,
		c.GetProxies().Strconv,
		c.GetUseCases().DeleteTodo,
		todoPresenter,
	))

	// Add completion command with subcommands
	completionCmd := completion.NewCompletionCommand(c.GetProxies().Cobra)
	completionCmd.AddCommand(completion.NewBashCompletionCommand(c.GetProxies().Cobra, rootCmd))
	completionCmd.AddCommand(completion.NewZshCompletionCommand(c.GetProxies().Cobra, rootCmd))
	completionCmd.AddCommand(completion.NewFishCompletionCommand(c.GetProxies().Cobra, rootCmd))
	completionCmd.AddCommand(completion.NewPowershellCompletionCommand(c.GetProxies().Cobra, rootCmd))
	rootCmd.AddCommand(completionCmd)

	return rootCmd, nil
}