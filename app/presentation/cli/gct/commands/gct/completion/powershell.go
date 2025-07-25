package completion

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewPowershellCompletionCommand creates the powershell completion subcommand
func NewPowershellCompletionCommand(cobraProxy proxy.Cobra, rootCmd proxy.Command) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse(powershellCompletionUse)
	cmd.SetShort(powershellCompletionShort)
	cmd.SetLong(powershellCompletionLong)
	cmd.SetArgs(cobraProxy.NoArgs())
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(func(_ *cobra.Command, _ []string) error {
		return runPowershellCompletion(rootCmd)
	})

	return cmd
}

func runPowershellCompletion(rootCmd proxy.Command) error {
	return rootCmd.GenPowerShellCompletion(proxy.Stdout)
}

const (
	powershellCompletionUse   string = "powershell"
	powershellCompletionShort string = "🔧🪟 Generate powershell completion script"
	powershellCompletionLong  string = `🔧🪟 Generate powershell completion script for gct.

To load completions:

PowerShell:
  PS> gct completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> gct completion powershell > gct.ps1
  # and source this file from your PowerShell profile.`
)
