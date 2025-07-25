package gct

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewDeleteCommand creates the delete command for the gct CLI application
func NewDeleteCommand(
	container *container.Container,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := container.GetProxies().Cobra.NewCommand()
	cmd.SetUse(deleteUse)
	cmd.SetShort(deleteShort)
	cmd.SetLong(deleteLong)
	cmd.SetArgs(container.GetProxies().Cobra.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	cmd.SetRunE(func(_ *cobra.Command, args []string) error {
		return runDelete(container.GetProxies().Strconv, container.GetUseCases().DeleteTodo, presenter, args[0])
	})

	return cmd
}

// runDelete executes the delete todo functionality
func runDelete(
	strconvProxy proxy.Strconv,
	deleteUseCase *application.DeleteTodoUseCase,
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

	// execute the delete use case
	err = deleteUseCase.Run(id)
	if err != nil {
		presenter.ShowError(err)
		return err
	}

	presenter.ShowDeleteSuccess(id)
	return nil
}

const (
	deleteUse   string = "delete"
	deleteShort string = "✅🗑️ Delete a todo"
	deleteLong  string = `✅🗑️ Delete a todo by its ID.

The ID must be a positive integer that corresponds to an existing todo.
Once deleted, the todo cannot be recovered.

Examples:
  gct delete 1
  gct delete 5`
)
