package command

import (
	o "os"

	"github.com/yanosea/gct/app/config"
	"github.com/yanosea/gct/app/presentation/cli/todo/formatter"
	"github.com/yanosea/gct/app/presentation/cli/todo/presenter"

	"github.com/yanosea/gct/pkg/proxy"
	"github.com/yanosea/gct/pkg/utility"
)

var (
	output string
	NewCli CreateCliFunc = newCli
)

type Cli interface {
	Init(envconfig proxy.Envconfig, json proxy.Json, os proxy.Os, fileUtil utility.FileUtil) int
	Run() int
}

type cli struct {
	Cobra       proxy.Cobra
	RootCommand proxy.Command
}

type CreateCliFunc func(cobra proxy.Cobra) Cli

func newCli(cobra proxy.Cobra) Cli {
	return &cli{
		Cobra:       cobra,
		RootCommand: nil,
	}
}

func (c *cli) Init(
	envconfig proxy.Envconfig,
	json proxy.Json,
	os proxy.Os,
	fileUtil utility.FileUtil,
) int {
	configurator := config.NewConfigurator(envconfig)
	conf, err := configurator.GetConfig()
	if err != nil {
		output = formatter.AppendErrorToOutput(err, output)
		presenter.Present(o.Stderr, output)
		return 1
	}

	c.RootCommand = NewRootCommand(
		c.Cobra,
		json,
		os,
		fileUtil,
		conf,
		&output,
	)

	return 0
}

func (c *cli) Run() int {
	out := o.Stdout
	exitCode := 0

	if err := c.RootCommand.Execute(); err != nil {
		output = formatter.AppendErrorToOutput(err, output)
		out = o.Stderr
		exitCode = 1
	}

	presenter.Present(out, output)

	return exitCode
}
