package gct

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewToggleCommand creates the toggle command for the gct CLI application
func NewToggleCommand(
	container *container.Container,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := container.GetProxies().Cobra.NewCommand()
	cmd.SetUse(toggleUse)
	cmd.SetShort(toggleShort)
	cmd.SetLong(toggleLong)
	cmd.SetArgs(container.GetProxies().Cobra.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		return runToggle(container.GetProxies().Strconv, container.GetUseCases().ToggleTodo, presenter, args[0])
	})

	return cmd
}

// runToggle executes the toggle todo functionality
func runToggle(
	strconvProxy proxy.Strconv,
	toggleUseCase *application.ToggleTodoUseCase,
	presenter *presenter.TodoPresenter,
	idStr string,
) error {
	// parse the ID from string to int
	id, err := strconvProxy.Atoi(idStr)
	if err != nil {
		presenter.ShowValidationError("invalid todo ID: must be a number")
		return err
	}

	// validate that ID is positive
	if id <= 0 {
		presenter.ShowValidationError("invalid todo ID: must be positive")
		return nil
	}

	// execute the toggle use case
	todo, err := toggleUseCase.Run(id)
	if err != nil {
		presenter.ShowError(err)
		return err
	}

	presenter.ShowToggleSuccess(todo)
	return nil
}

const (
	toggleUse   string = "toggle"
	toggleShort string = "✅♻️ Toggle the completion status of a todo"
	toggleLong  string = `✅♻️ Toggle the completion status of a todo by its ID.

The ID must be a positive integer that corresponds to an existing todo.

Examples:
  gct toggle 1
  gct toggle 5`
)
