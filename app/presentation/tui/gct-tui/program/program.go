package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/yanosea/gct/app/container"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
)

// InitializeProgram initializes the TUI application with all dependencies
func InitializeProgram() (*tea.Program, error) {
	// Initialize dependency injection container
	container, err := container.NewContainer()
	if err != nil {
		return nil, err
	}

	// Get dependencies from container
	useCases := container.GetUseCases()
	proxies := container.GetProxies()

	// Create the TUI model with dependencies
	stateModel := model.NewStateModel(
		useCases.AddTodo,
		useCases.ListTodo,
		useCases.ToggleTodo,
		useCases.DeleteTodo,
		proxies.Bubbles,
	)

	// Create the bubbletea program
	program := tea.NewProgram(
		stateModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return program, nil
}
