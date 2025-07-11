package main

import (
	"flag"
	"fmt"

	"github.com/yanosea/cleancobra/app/presentation/tui/todo-tui/command"
	"github.com/yanosea/cleancobra/pkg/proxy"
	"github.com/yanosea/cleancobra/pkg/utility"
)

const helpText = `A clean architecture TODO application

Usage:
  todo-tui [options]

Available Operations (Interactive):
  ↑/k         Move cursor up
  ↓/j         Move cursor down
  enter/space Toggle todo status
  a           Add a new todo
  d           Delete selected todo
  r           Refresh todo list
  q           Quit application

Options:
  -h, --help  Show this help message`

var (
	bubbletea = proxy.NewBubbletea()
	envconfig = proxy.NewEnvconfig()
	json      = proxy.NewJson()
	os        = proxy.NewOs()
	fileutil  = utility.NewFileUtil(os, proxy.NewJson())
)

func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && args[0] == "help" {
		showHelp = true
	}

	if showHelp {
		fmt.Println(helpText)
		return
	}

	tui := command.NewTui(
		bubbletea,
		envconfig,
		json,
		os,
		fileutil,
	)
	os.Exit(tui.Run())
}
