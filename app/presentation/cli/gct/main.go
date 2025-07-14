package main

import (
	"os"

	"github.com/yanosea/gct/app/composition"
	"github.com/yanosea/gct/app/presentation/cli/gct/command"
)

var (
	exit = os.Exit
)

func main() {
	// Initialize dependency injection container (composition root)
	diContainer := composition.NewContainer()
	if err := diContainer.Initialize(); err != nil {
		exit(1)
	}

	// Create CLI with dependencies from container
	cli := command.NewCli(
		diContainer.GetCobra(),
	)

	if exitCode := cli.Init(
		diContainer.GetEnvconfig(),
		diContainer.GetJson(),
		diContainer.GetOs(),
		diContainer.GetFileUtil(),
	); exitCode != 0 {
		exit(exitCode)
	}

	exit(cli.Run())
}
