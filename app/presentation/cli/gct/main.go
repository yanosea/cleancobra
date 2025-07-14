package main

import (
	"github.com/yanosea/gct/app/presentation/cli/gct/command"

	"github.com/yanosea/gct/pkg/proxy"
	"github.com/yanosea/gct/pkg/utility"
)

type TodoCliParams struct {
	Cobra     proxy.Cobra
	Envconfig proxy.Envconfig
	Json      proxy.Json
	Os        proxy.Os
	FileUtil  utility.FileUtil
}

var (
	exit          = os.Exit
	os            = proxy.NewOs()
	todoCliParams = TodoCliParams{
		Cobra:     proxy.NewCobra(),
		Envconfig: proxy.NewEnvconfig(),
		Json:      proxy.NewJson(),
		Os:        os,
		FileUtil:  utility.NewFileUtil(os, proxy.NewJson()),
	}
)

func main() {
	cli := command.NewCli(
		todoCliParams.Cobra,
	)
	if exitCode := cli.Init(
		todoCliParams.Envconfig,
		todoCliParams.Json,
		todoCliParams.Os,
		todoCliParams.FileUtil,
	); exitCode != 0 {
		exit(exitCode)
	}

	os.Exit(cli.Run())
}
