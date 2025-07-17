package gct

import (
	"github.com/spf13/cobra"
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

// NewToggleCommand creates the toggle command for the gct CLI application
func NewToggleCommand(
	cobraProxy proxy.Cobra,
	strconvProxy proxy.Strconv,
	toggleUseCase *application.ToggleTodoUseCase,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("toggle <id>")
	cmd.SetShort("Toggle the completion status of a todo")
	cmd.SetLong(`Toggle the completion status of a todo by its ID.

The ID must be a positive integer that corresponds to an existing todo.

Examples:
  gct toggle 1
  gct toggle 5`)
	cmd.SetArgs(cobraProxy.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	
	// Set the run function
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		return RunToggle(strconvProxy, toggleUseCase, presenter, args[0])
	})
	
	return cmd
}

// RunToggle executes the toggle todo functionality
func RunToggle(strconvProxy proxy.Strconv, toggleUseCase *application.ToggleTodoUseCase, presenter *presenter.TodoPresenter, idStr string) error {
	// Parse the ID from string to int
	id, err := strconvProxy.Atoi(idStr)
	if err != nil {
		presenter.ShowValidationError("invalid todo ID: must be a number")
		return err
	}
	
	// Validate that ID is positive
	if id <= 0 {
		presenter.ShowValidationError("invalid todo ID: must be positive")
		return nil
	}
	
	// Execute the toggle use case
	todo, err := toggleUseCase.Run(id)
	if err != nil {
		presenter.ShowError(err)
		return err
	}
	
	presenter.ShowToggleSuccess(todo)
	return nil
}