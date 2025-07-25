package completion

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewCompletionCommand creates the completion parent command
func NewCompletionCommand(
	rootCmd proxy.Command,
	container *container.Container,
	presenter *presenter.TodoPresenter,
) proxy.Command {
	cmd := container.GetProxies().Cobra.NewCommand()
	cmd.SetUse(completionUse)
	cmd.SetShort(completionShort)
	cmd.SetLong(completionLong)
	cmd.SetArgs(container.GetProxies().Cobra.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.AddCommand(NewBashCompletionCommand(container.GetProxies().Cobra, rootCmd))
	cmd.AddCommand(NewZshCompletionCommand(container.GetProxies().Cobra, rootCmd))
	cmd.AddCommand(NewFishCompletionCommand(container.GetProxies().Cobra, rootCmd))
	cmd.AddCommand(NewPowershellCompletionCommand(container.GetProxies().Cobra, rootCmd))
	rootCmd.SetRunE(func(_ *cobra.Command, args []string) error {
		return runCompletion(presenter)
	})

	return cmd
}

func runCompletion(presenter *presenter.TodoPresenter) error {
	presenter.ShowError(domain.ErrNoSubCommand)

	return nil
}

const (
	completionUse   string = "completion"
	completionShort string = "🔧 Generate shell completion scripts"
	completionLong  string = `🔧 Generate shell completion scripts for gct.

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
`
)
