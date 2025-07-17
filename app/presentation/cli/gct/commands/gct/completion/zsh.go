package completion

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewZshCompletionCommand creates the zsh completion subcommand
func NewZshCompletionCommand(cobraProxy proxy.Cobra, rootCmd proxy.Command) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("zsh")
	cmd.SetShort("Generate zsh completion script")
	cmd.SetLong(`Generate zsh completion script for gct.

To load completions:

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ gct completion zsh > "${fpath[1]}/_gct"

  # You will need to start a new shell for this setup to take effect.`)
	cmd.SetArgs(cobraProxy.NoArgs())
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(func(_ *cobra.Command, _ []string) error {
		return runZshCompletion(rootCmd)
	})

	return cmd
}

func runZshCompletion(rootCmd proxy.Command) error {
	return rootCmd.GenZshCompletion(os.Stdout)
}