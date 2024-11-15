package main

import (
	"github.com/yanosea/cleancobra/app/presentation/cli/todo/command"

	"github.com/yanosea/cleancobra/pkg/proxy"
	"github.com/yanosea/cleancobra/pkg/utility"
)

var (
	cobra     = proxy.NewCobra()
	envconfig = proxy.NewEnvconfig()
	json      = proxy.NewJson()
	os        = proxy.NewOs()
	fileutil  = utility.NewFileUtil(os, proxy.NewJson())
)

func main() {
	cli := command.NewCli(
		cobra,
		envconfig,
		json,
		os,
		fileutil,
	)
	os.Exit(cli.Run())
}
