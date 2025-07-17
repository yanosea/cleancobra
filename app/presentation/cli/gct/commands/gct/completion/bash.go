package completion

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewBashCompletionCommand creates the bash completion subcommand
func NewBashCompletionCommand(cobraProxy proxy.Cobra, rootCmd proxy.Command) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("bash")
	cmd.SetShort("Generate bash completion script")
	cmd.SetLong(`Generate bash completion script for gct.

To load completions:

Bash:
  $ source <(gct completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ gct completion bash > /etc/bash_completion.d/gct
  # macOS:
  $ gct completion bash > /usr/local/etc/bash_completion.d/gct

You will need to start a new shell for this setup to take effect.`)
	cmd.SetArgs(cobraProxy.NoArgs())
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(func(_ *cobra.Command, _ []string) error {
		return runBashCompletion(rootCmd)
	})

	return cmd
}

func runBashCompletion(rootCmd proxy.Command) error {
	return rootCmd.GenBashCompletion(os.Stdout)
}