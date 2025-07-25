package commands

import (
	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/cli/gct/formatter"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"

	"github.com/yanosea/gct/pkg/proxy"
)

// InitializeCommand initializes the CLI application with all commands and dependencies
func InitializeCommand() (proxy.Command, error) {
	// initialize dependency injection container
	c, err := container.NewContainer()
	if err != nil {
		return nil, err
	}

	// initialize presenter
	todoPresenter := presenter.NewTodoPresenter(
		c.GetProxies().Errors,
		c.GetProxies().Fmt,
		c.GetProxies().OS,
		formatter.NewJSONFormatter(c.GetProxies().JSON),
		formatter.NewPlainFormatter(c.GetProxies().Fmt, c.GetProxies().Strconv, c.GetProxies().Strings),
		formatter.NewTableFormatter(c.GetProxies().Color, c.GetProxies().Strings, c.GetProxies().Fmt),
	)

	// initialize root command
	rootCmd := NewRootCommand(
		c,
		todoPresenter,
	)

	return rootCmd, nil
}
