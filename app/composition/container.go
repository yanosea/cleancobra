package composition

import (
	"github.com/yanosea/gct/app/application/gct"
	"github.com/yanosea/gct/app/config"
	domainStorage "github.com/yanosea/gct/app/domain/storage"
	todoDomain "github.com/yanosea/gct/app/domain/todo"
	"github.com/yanosea/gct/app/infrastructure/json/repository"
	"github.com/yanosea/gct/app/infrastructure/storage"
	"github.com/yanosea/gct/pkg/proxy"
	"github.com/yanosea/gct/pkg/utility"
)

// Container holds all dependencies for dependency injection
// This is the composition root where all dependencies are wired together
type Container struct {
	// Config
	config *config.TodoConfig

	// Proxies (external dependencies)
	cobra     proxy.Cobra
	envconfig proxy.Envconfig
	json      proxy.Json
	os        proxy.Os

	// Utilities
	fileUtil utility.FileUtil

	// Storage (infrastructure layer)
	fileStorage domainStorage.FileStorage

	// Repositories (infrastructure layer)
	todoRepo todoDomain.TodoRepository

	// Use Cases (application layer)
	addTodoUseCase    *gct.AddTodoUseCase
	listTodoUseCase   *gct.ListTodoUseCase
	toggleTodoUseCase *gct.ToggleTodoUseCase
	deleteTodoUseCase *gct.DeleteTodoUseCase
}

// NewContainer creates a new dependency injection container
// This is the composition root of the application
func NewContainer() *Container {
	return &Container{}
}

// InitializeProxies initializes all external proxy dependencies
func (c *Container) InitializeProxies() {
	c.cobra = proxy.NewCobra()
	c.envconfig = proxy.NewEnvconfig()
	c.json = proxy.NewJson()
	c.os = proxy.NewOs()
	c.fileUtil = utility.NewFileUtil(c.os, c.json)
}

// InitializeConfig initializes application configuration
func (c *Container) InitializeConfig() error {
	configurator := config.NewConfigurator(c.envconfig)
	config, err := configurator.GetConfig()
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

// InitializeStorage initializes storage implementations
func (c *Container) InitializeStorage() {
	c.fileStorage = storage.NewOSFileStorage(c.os)
}

// InitializeRepositories initializes all repository implementations
// This is where we wire domain interfaces to infrastructure implementations
func (c *Container) InitializeRepositories() error {
	todoRepo, err := repository.NewTodoRepository(
		c.config,
		c.fileUtil,
		c.json,
		c.os,
	)
	if err != nil {
		return err
	}
	c.todoRepo = todoRepo
	return nil
}

// InitializeUseCases initializes all application use cases
// This is where we inject repositories into use cases
func (c *Container) InitializeUseCases() {
	c.addTodoUseCase = gct.NewAddTodoUseCase(c.todoRepo)
	c.listTodoUseCase = gct.NewListTodoUseCase(c.todoRepo)
	c.toggleTodoUseCase = gct.NewToggleTodoUseCase(c.todoRepo)
	c.deleteTodoUseCase = gct.NewDeleteTodoUseCase(c.todoRepo)
}

// Initialize initializes all dependencies in the correct order
// This method orchestrates the entire dependency graph construction
func (c *Container) Initialize() error {
	c.InitializeProxies()

	if err := c.InitializeConfig(); err != nil {
		return err
	}

	c.InitializeStorage()

	if err := c.InitializeRepositories(); err != nil {
		return err
	}

	c.InitializeUseCases()
	return nil
}

func (c *Container) GetCobra() proxy.Cobra         { return c.cobra }
func (c *Container) GetEnvconfig() proxy.Envconfig { return c.envconfig }
func (c *Container) GetJson() proxy.Json           { return c.json }
func (c *Container) GetOs() proxy.Os               { return c.os }
func (c *Container) GetFileUtil() utility.FileUtil { return c.fileUtil }
