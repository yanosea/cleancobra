package command

import (
	c "github.com/spf13/cobra"

	"github.com/yanosea/cleancobra/app/config"
	"github.com/yanosea/cleancobra/app/presentation/cli/todo/command/todo"

	"github.com/yanosea/cleancobra/pkg/proxy"
	"github.com/yanosea/cleancobra/pkg/utility"
)

func NewRootCommand(
	cobra proxy.Cobra,
	envconfig proxy.Envconfig,
	json proxy.Json,
	os proxy.Os,
	fileutil utility.FileUtil,
	conf *config.TodoConfig,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetSilenceErrors(true)
	cmd.SetUse("todo")
	cmd.SetShort("A clean architecture TODO application")

	listCmd := todo.NewListCommand(
		cobra,
		json,
		os,
		fileutil,
		conf,
		output,
	)

	cmd.AddCommand(
		todo.NewAddCommand(
			cobra,
			json,
			os,
			fileutil,
			conf,
			output,
		),
		todo.NewDeleteCommand(
			cobra,
			json,
			os,
			fileutil,
			conf,
			output,
		),
		todo.NewToggleCommand(
			cobra,
			json,
			os,
			fileutil,
			conf,
			output,
		),
		listCmd,
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runRoot(cmd, args, listCmd)
		},
	)

	return cmd
}

// runRoot runs the root command with default list behavior
func runRoot(
	cmd *c.Command,
	args []string,
	listCmd proxy.Command,
) error {
	return listCmd.RunE(cmd, args)
}
