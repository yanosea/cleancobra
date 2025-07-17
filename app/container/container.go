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
	OS       proxy.OS
	Filepath proxy.Filepath
	JSON     proxy.JSON
	Time     proxy.Time
	IO       proxy.IO
	Fmt      proxy.Fmt
	Strings  proxy.Strings
	Strconv  proxy.Strconv
	Cobra    proxy.Cobra
	Bubbletea proxy.Bubbletea
	Bubbles  proxy.Bubbles
	Lipgloss proxy.Lipgloss
	Color    proxy.Color
	Envconfig proxy.Envconfig
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
		OS:        proxy.NewOS(),
		Filepath:  proxy.NewFilepath(),
		JSON:      proxy.NewJSON(),
		Time:      proxy.NewTime(),
		IO:        proxy.NewIO(),
		Fmt:       proxy.NewFmt(),
		Strings:   proxy.NewStrings(),
		Strconv:   proxy.NewStrconv(),
		Cobra:     proxy.NewCobra(),
		Bubbletea: proxy.NewBubbletea(),
		Bubbles:   proxy.NewBubbles(),
		Lipgloss:  proxy.NewLipgloss(),
		Color:     proxy.NewColor(),
		Envconfig: proxy.NewEnvconfig(),
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize repository with required proxies
	repository := infrastructure.NewJSONRepository(
		cfg.DataFile,
		proxies.OS,
		proxies.JSON,
	)

	// Initialize use cases with repository dependency
	useCases := &UseCases{
		AddTodo:    application.NewAddTodoUseCase(repository),
		DeleteTodo: application.NewDeleteTodoUseCase(repository),
		ListTodo:   application.NewListTodoUseCase(repository),
		ToggleTodo: application.NewToggleTodoUseCase(repository),
	}

	return &Container{
		config:     cfg,
		proxies:    proxies,
		repository: repository,
		useCases:   useCases,
	}, nil
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