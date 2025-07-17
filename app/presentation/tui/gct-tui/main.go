package main

import (
	"os"

	"github.com/yanosea/gct/app/presentation/tui/gct-tui/program"
)

func main() {
	// Initialize program and dependencies
	prog, err := program.InitializeProgram()
	if err != nil {
		os.Exit(1)
	}

	// Execute the program
	if _, err := prog.Run(); err != nil {
		os.Exit(1)
	}
}