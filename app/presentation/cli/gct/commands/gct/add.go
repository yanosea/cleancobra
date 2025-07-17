package gct

import (
	"github.com/spf13/cobra"
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

// NewAddCommand creates the add command for the gct CLI application
func NewAddCommand(
	cobraProxy proxy.Cobra,
	addUseCase *application.AddTodoUseCase,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("add <description>")
	cmd.SetShort("Add a new todo")
	cmd.SetLong(`Add a new todo with the specified description.

The description is required and cannot be empty.

Examples:
  gct add "Buy groceries"
  gct add "Complete project documentation"`)
	cmd.SetArgs(cobraProxy.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	
	// Set the run function
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		return RunAdd(addUseCase, presenter, args[0])
	})
	
	return cmd
}

// RunAdd executes the add todo functionality
func RunAdd(addUseCase *application.AddTodoUseCase, presenter *presenter.TodoPresenter, description string) error {
	todo, err := addUseCase.Run(description)
	if err != nil {
		presenter.ShowError(err)
		return err
	}
	
	presenter.ShowAddSuccess(todo)
	return nil
}