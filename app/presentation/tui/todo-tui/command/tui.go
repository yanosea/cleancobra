package command

import (
	"fmt"
	"os"

	todoApp "github.com/yanosea/cleancobra/app/application/todo"
	"github.com/yanosea/cleancobra/app/config"
	"github.com/yanosea/cleancobra/app/infrastructure/json/repository"
	"github.com/yanosea/cleancobra/app/presentation/tui/todo-tui/model"
	"github.com/yanosea/cleancobra/pkg/proxy"
	"github.com/yanosea/cleancobra/pkg/utility"
)

type Tui struct {
	Bubbletea     proxy.Bubbletea
	Envconfig     proxy.Envconfig
	Json          proxy.Json
	Os            proxy.Os
	FileUtil      utility.FileUtil
	Config        *config.TodoConfig
	NewRootRunner func(proxy.Bubbletea, *model.Usecases) *Runner
}

func NewTui(
	bubbletea proxy.Bubbletea,
	envconfig proxy.Envconfig,
	json proxy.Json,
	os proxy.Os,
	fileutil utility.FileUtil,
) *Tui {
	return &Tui{
		Bubbletea:     bubbletea,
		Envconfig:     envconfig,
		Json:          json,
		Os:            os,
		FileUtil:      fileutil,
		Config:        nil,
		NewRootRunner: NewRootRunner,
	}
}

func (t *Tui) Run() int {
	configurator := config.NewConfigurator(t.Envconfig)
	conf, err := configurator.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		return 1
	}
	t.Config = conf

	todoRepo, err := repository.NewTodoRepository(conf, t.FileUtil, t.Json, t.Os)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize repository: %v\n", err)
		return 1
	}

	usecases := &model.Usecases{
		List:   todoApp.NewListTodoUseCase(todoRepo),
		Add:    todoApp.NewAddTodoUseCase(todoRepo),
		Delete: todoApp.NewDeleteTodoUseCase(todoRepo),
		Toggle: todoApp.NewToggleTodoUseCase(todoRepo),
	}

	runner := t.NewRootRunner(t.Bubbletea, usecases)
	return runner.Run()
}
