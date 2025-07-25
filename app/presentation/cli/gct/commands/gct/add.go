package gct

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewAddCommand creates the add command for the gct CLI application
func NewAddCommand(
	container *container.Container,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := container.GetProxies().Cobra.NewCommand()
	cmd.SetUse(addUse)
	cmd.SetShort(addShort)
	cmd.SetLong(addLong)
	cmd.SetArgs(container.GetProxies().Cobra.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	cmd.SetRunE(func(_ *cobra.Command, args []string) error {
		return runAdd(container.GetUseCases().AddTodo, presenter, args[0])
	})

	return cmd
}

// runAdd executes the add todo functionality
func runAdd(
	addUseCase *application.AddTodoUseCase,
	presenter *presenter.TodoPresenter,
	description string,
) error {
	todo, err := addUseCase.Run(description)
	if err != nil {
		presenter.ShowError(err)
		return err
	}

	presenter.ShowAddSuccess(todo)
	return nil
}

const (
	addUse   string = "add"
	addShort string = "✅➕ Add a new todo"
	addLong  string = `✅➕ Add a new todo with the specified description.

The description is required and cannot be empty.

Examples:
  gct add "Buy groceries"
  gct add "Complete project documentation"`
)
