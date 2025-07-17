package completion

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/yanosea/gct/pkg/proxy"
)

// NewPowershellCompletionCommand creates the powershell completion subcommand
func NewPowershellCompletionCommand(cobraProxy proxy.Cobra, rootCmd proxy.Command) proxy.Command {
	cmd := cobraProxy.NewCommand()
	cmd.SetUse("powershell")
	cmd.SetShort("Generate powershell completion script")
	cmd.SetLong(`Generate powershell completion script for gct.

To load completions:

PowerShell:
  PS> gct completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> gct completion powershell > gct.ps1
  # and source this file from your PowerShell profile.`)
	cmd.SetArgs(cobraProxy.NoArgs())
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(func(_ *cobra.Command, _ []string) error {
		return runPowershellCompletion(rootCmd)
	})

	return cmd
}

func runPowershellCompletion(rootCmd proxy.Command) error {
	return rootCmd.GenPowerShellCompletion(os.Stdout)
}