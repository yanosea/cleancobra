package main

import (
	"os"

	"github.com/fatih/color"

	"github.com/yanosea/cleancobra/config"
	"github.com/yanosea/cleancobra/presentation/cli/command"

	"github.com/yanosea/cleancobra-pkg/util"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		util.PrintlnWithWriter(os.Stderr, color.RedString(err.Error()))
		os.Exit(1)
	}
	cc := command.NewCleanCobra(conf)
	os.Exit(cc.Run())
}
