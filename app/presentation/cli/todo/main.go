package main

import (
	"github.com/yanosea/gct/app/presentation/cli/todo/command"

	"github.com/yanosea/gct/pkg/proxy"
	"github.com/yanosea/gct/pkg/utility"
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
