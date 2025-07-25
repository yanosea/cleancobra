package completion

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewFishCompletionCommand creates the fish completion subcommand
func NewFishCompletionCommand(cobraProxy proxy.Cobra, rootCmd proxy.Command) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse(fishCompletionUse)
	cmd.SetShort(fishCompletionShort)
	cmd.SetLong(fishCompletionLong)
	cmd.SetArgs(cobraProxy.NoArgs())
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(func(_ *cobra.Command, _ []string) error {
		return runFishCompletion(rootCmd)
	})

	return cmd
}

func runFishCompletion(rootCmd proxy.Command) error {
	return rootCmd.GenFishCompletion(proxy.Stdout, true)
}

const (
	fishCompletionUse   string = "fish"
	fishCompletionShort string = "🔧🐟 Generate fish completion script"
	fishCompletionLong  string = `🔧🐟 Generate fish completion script for gct.

To load completions:

Fish:
  $ gct completion fish | source

  # To load completions for each session, execute once:
  $ gct completion fish > ~/.config/fish/completions/gct.fish`
)
