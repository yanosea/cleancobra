package command

import (
	o "os"

	"github.com/yanosea/gct/app/config"
	"github.com/yanosea/gct/app/presentation/cli/todo/formatter"
	"github.com/yanosea/gct/app/presentation/cli/todo/presenter"

	"github.com/yanosea/gct/pkg/proxy"
	"github.com/yanosea/gct/pkg/utility"
)

type Cli struct {
	Cobra          proxy.Cobra
	Envconfig      proxy.Envconfig
	Json           proxy.Json
	Os             proxy.Os
	FileUtil       utility.FileUtil
	Config         *config.TodoConfig
	NewRootCommand func(
		cobra proxy.Cobra,
		envconfig proxy.Envconfig,
		json proxy.Json,
		os proxy.Os,
		fileutil utility.FileUtil,
		conf *config.TodoConfig,
		output *string,
	) proxy.Command
}

var (
	output string
)

func NewCli(
	cobra proxy.Cobra,
	envconfig proxy.Envconfig,
	json proxy.Json,
	os proxy.Os,
	fileutil utility.FileUtil,
) *Cli {
	return &Cli{
		Cobra:          cobra,
		Envconfig:      envconfig,
		Json:           json,
		Os:             os,
		FileUtil:       fileutil,
		Config:         nil,
		NewRootCommand: NewRootCommand,
	}
}

func (c *Cli) Run() int {
	configurator := config.NewConfigurator(c.Envconfig)
	conf, err := configurator.GetConfig()
	if err != nil {
		return 1
	}
	rootCmd := c.NewRootCommand(
		c.Cobra,
		c.Envconfig,
		c.Json,
		c.Os,
		c.FileUtil,
		conf,
		&output,
	)

	out := o.Stdout
	exitCode := 0

	if err := rootCmd.Execute(); err != nil {
		output = formatter.AppendErrorToOutput(err, output)
		out = o.Stderr
		exitCode = 1
	}

	presenter.Present(out, output)

	return exitCode
}
