package main

import (
	"os"

	"github.com/yanosea/cleancobra/config"
	"github.com/yanosea/cleancobra/presentation/cli/command"
)

func main() {
	conf := config.GetConfig()
	cc := command.NewCleanCobra(conf)
	os.Exit(cc.Run())
}
