package completion

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewFishCompletionCommand creates the fish completion subcommand
func NewFishCompletionCommand(cobraProxy proxy.Cobra, rootCmd proxy.Command) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("fish")
	cmd.SetShort("Generate fish completion script")
	cmd.SetLong(`Generate fish completion script for gct.

To load completions:

Fish:
  $ gct completion fish | source

  # To load completions for each session, execute once:
  $ gct completion fish > ~/.config/fish/completions/gct.fish`)
	cmd.SetArgs(cobraProxy.NoArgs())
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(func(_ *cobra.Command, _ []string) error {
		return runFishCompletion(rootCmd)
	})

	return cmd
}

func runFishCompletion(rootCmd proxy.Command) error {
	return rootCmd.GenFishCompletion(os.Stdout, true)
}