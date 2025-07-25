package container

import (
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/config"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/infrastructure"

	"github.com/yanosea/gct/pkg/proxy"
)

// Container manages all application dependencies
type Container struct {
	config     *config.Config
	proxies    *Proxies
	repository domain.TodoRepository
	useCases   *UseCases
}

// Proxies organizes all proxy dependencies
type Proxies struct {
	Bubbles   proxy.Bubbles
	Bubbletea proxy.Bubbletea
	Cobra     proxy.Cobra
	Color     proxy.Color
	Envconfig proxy.Envconfig
	Errors    proxy.Errors
	Filepath  proxy.Filepath
	Fmt       proxy.Fmt
	IO        proxy.IO
	JSON      proxy.JSON
	Lipgloss  proxy.Lipgloss
	OS        proxy.OS
	Sort      proxy.Sort
	Strconv   proxy.Strconv
	Strings   proxy.Strings
	Time      proxy.Time
}

// UseCases organizes all use case dependencies
type UseCases struct {
	AddTodo    *application.AddTodoUseCase
	DeleteTodo *application.DeleteTodoUseCase
	ListTodo   *application.ListTodoUseCase
	ToggleTodo *application.ToggleTodoUseCase
}

// NewContainer creates a new Container with all dependencies properly initialized
func NewContainer() (*Container, error) {
	// Initialize proxies first
	proxies := &Proxies{
		Bubbles:   proxy.NewBubbles(),
		Bubbletea: proxy.NewBubbletea(),
		Cobra:     proxy.NewCobra(),
		Color:     proxy.NewColor(),
		Envconfig: proxy.NewEnvconfig(),
		Errors:    proxy.NewErrors(),
		Filepath:  proxy.NewFilepath(),
		Fmt:       proxy.NewFmt(),
		IO:        proxy.NewIO(),
		JSON:      proxy.NewJSON(),
		Lipgloss:  proxy.NewLipgloss(),
		OS:        proxy.NewOS(),
		Sort:      proxy.NewSort(),
		Strconv:   proxy.NewStrconv(),
		Strings:   proxy.NewStrings(),
		Time:      proxy.NewTime(),
	}

	// initialize domain
	domain.InitializeDomain(
		proxies.Fmt,
		proxies.JSON,
		proxies.Strings,
		proxies.Time,
	)

	// initialize domain error
	domain.InitializeDomainErrors(
		proxies.Errors,
		proxies.Fmt,
	)

	// load configuration
	configurator := config.NewConfigurator(
		proxies.Envconfig,
		proxies.Filepath,
		proxies.OS,
	)
	cfg, err := configurator.Load()
	if err != nil {
		return nil, err
	}

	// initialize repository with required proxies
	repository := infrastructure.NewJSONRepository(
		cfg.DataFile,
		proxies.Filepath,
		proxies.JSON,
		proxies.OS,
		proxies.Sort,
	)

	// initialize use cases with repository dependency
	useCases := &UseCases{
		AddTodo:    application.NewAddTodoUseCase(repository),
		DeleteTodo: application.NewDeleteTodoUseCase(repository),
		ListTodo:   application.NewListTodoUseCase(repository),
		ToggleTodo: application.NewToggleTodoUseCase(repository),
	}

	container := &Container{
		config:     cfg,
		proxies:    proxies,
		repository: repository,
		useCases:   useCases,
	}

	return container, nil
}

// GetUseCases returns the use cases struct
func (c *Container) GetUseCases() *UseCases {
	return c.useCases
}

// GetRepository returns the todo repository
func (c *Container) GetRepository() domain.TodoRepository {
	return c.repository
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// GetProxies returns the proxies struct
func (c *Container) GetProxies() *Proxies {
	return c.proxies
}
