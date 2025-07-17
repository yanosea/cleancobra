package gct

import (
	"github.com/spf13/cobra"
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

// NewDeleteCommand creates the delete command for the gct CLI application
func NewDeleteCommand(
	cobraProxy proxy.Cobra,
	strconvProxy proxy.Strconv,
	deleteUseCase *application.DeleteTodoUseCase,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("delete <id>")
	cmd.SetShort("Delete a todo")
	cmd.SetLong(`Delete a todo by its ID.

The ID must be a positive integer that corresponds to an existing todo.
Once deleted, the todo cannot be recovered.

Examples:
  gct delete 1
  gct delete 5`)
	cmd.SetArgs(cobraProxy.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.SetSilenceUsage(true)
	
	// Set the run function
	cmd.SetRunE(func(cobraCmd *cobra.Command, args []string) error {
		return RunDelete(strconvProxy, deleteUseCase, presenter, args[0])
	})
	
	return cmd
}

// RunDelete executes the delete todo functionality
func RunDelete(strconvProxy proxy.Strconv, deleteUseCase *application.DeleteTodoUseCase, presenter *presenter.TodoPresenter, idStr string) error {
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
	
	// Execute the delete use case
	err = deleteUseCase.Run(id)
	if err != nil {
		presenter.ShowError(err)
		return err
	}
	
	presenter.ShowDeleteSuccess(id)
	return nil
}