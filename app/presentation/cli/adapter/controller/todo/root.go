package todo

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"cleancobra/app/presentation/cli/adapter"
)

type GlobalOption struct {
	Out            io.Writer
	ErrOut         io.Writer
	TodoUsecases   *adapter.TodoUsecases
	NewRootCommand func(*adapter.TodoUsecases) *cobra.Command
}

func NewGlobalOption(uc *adapter.TodoUsecases) *GlobalOption {
	return &GlobalOption{
		Out:            os.Stdout,
		ErrOut:         os.Stderr,
		TodoUsecases:   uc,
		NewRootCommand: NewRootCommand,
	}
}

func NewRootCommand(uc *adapter.TodoUsecases) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleancobra",
		Short: "A clean architecture TODO application",
	}

	cmd.AddCommand(
		NewAddCommand(uc.AddTodoUseCase),
		NewListCommand(uc.ListTodoUseCase),
		NewToggleCommand(uc.ToggleTodoUseCase),
		NewDeleteCommand(uc.DeleteTodoUseCase),
	)

	return cmd
}

func (g *GlobalOption) Execute() int {
	rootCmd := g.NewRootCommand(g.TodoUsecases)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(g.ErrOut, color.RedString(err.Error()))
		return 1
	}
	return 0
}
