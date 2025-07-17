package completion

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewCompletionCommand creates the completion parent command
func NewCompletionCommand(cobraProxy proxy.Cobra) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("completion")
	cmd.SetShort("Generate shell completion scripts")
	cmd.SetLong(`Generate shell completion scripts for gct.

The completion command allows you to generate shell completion scripts for bash, zsh, fish, and powershell.

To load completions:

Bash:
  $ source <(gct completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ gct completion bash > /etc/bash_completion.d/gct
  # macOS:
  $ gct completion bash > /usr/local/etc/bash_completion.d/gct

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ gct completion zsh > "${fpath[1]}/_gct"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ gct completion fish | source

  # To load completions for each session, execute once:
  $ gct completion fish > ~/.config/fish/completions/gct.fish

PowerShell:
  PS> gct completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> gct completion powershell > gct.ps1
  # and source this file from your PowerShell profile.
`)
	cmd.SetArgs(cobraProxy.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(func(_ *cobra.Command, args []string) error {
		return runCompletion(args[0])
	})

	return cmd
}

func runCompletion(shell string) error {
	return fmt.Errorf("completion command requires a subcommand: bash, zsh, fish, or powershell")
}