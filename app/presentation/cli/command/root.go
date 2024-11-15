package command

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/yanosea/cleancobra/config"
	"github.com/yanosea/cleancobra/presentation/cli/command/todo"
)

type CleanCobra struct {
	Out            io.Writer
	ErrOut         io.Writer
	Config         *config.Config
	NewRootCommand func(conf *config.Config) *cobra.Command
}

func NewCleanCobra(conf *config.Config) *CleanCobra {
	return &CleanCobra{
		Out:            os.Stdout,
		ErrOut:         os.Stderr,
		Config:         conf,
		NewRootCommand: NewRootCommand,
	}
}

func NewRootCommand(conf *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleancobra",
		Short: "A clean architecture application",
	}

	cmd.AddCommand(
		todo.NewTodoCommand(&conf.Todo),
	)

	return cmd
}

func (g *CleanCobra) Run() int {
	rootCmd := g.NewRootCommand(g.Config)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(g.ErrOut, color.RedString(err.Error()))
		return 1
	}
	return 0
}
