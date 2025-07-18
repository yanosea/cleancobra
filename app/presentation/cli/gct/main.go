package main

import (
	"os"

	"github.com/yanosea/gct/app/presentation/cli/gct/commands"
)

func main() {
	// Initialize commands and dependencies
	cmd, err := commands.InitializeCommand()
	if err != nil {
		os.Exit(1)
	}

	// Execute the command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
