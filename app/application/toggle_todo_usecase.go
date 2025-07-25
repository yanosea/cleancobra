package application

import (
	"github.com/yanosea/gct/app/domain"
)

// ToggleTodoUseCase handles the business logic for toggling todo completion status
type ToggleTodoUseCase struct {
	repository domain.TodoRepository
}

// NewToggleTodoUseCase creates a new ToggleTodoUseCase instance
func NewToggleTodoUseCase(repository domain.TodoRepository) *ToggleTodoUseCase {
	return &ToggleTodoUseCase{
		repository: repository,
	}
}

// Run executes the toggle todo use case
func (uc *ToggleTodoUseCase) Run(id int) (*domain.Todo, error) {
	// validate input
	if id <= 0 {
		return nil, domain.NewDomainError(
			domain.ErrorTypeInvalidInput,
			"todo ID must be positive",
			nil,
		)
	}

	// get all todos to find the one to toggle
	todos, err := uc.repository.FindAll()
	if err != nil {
		return nil, err
	}

	// find the todo with the specified ID
	var todoToToggle *domain.Todo
	for i, todo := range todos {
		if todo.ID == id {
			todoToToggle = &todos[i]
			break
		}
	}

	if todoToToggle == nil {
		return nil, domain.ErrTodoNotFound
	}

	// toggle the completion status
	todoToToggle.Toggle()

	// save the updated todo
	savedTodos, err := uc.repository.Save(*todoToToggle)
	if err != nil {
		return nil, err
	}

	if len(savedTodos) == 0 {
		return nil, domain.NewDomainError(
			domain.ErrorTypeConfiguration,
			"repository returned empty result",
			nil,
		)
	}

	return &savedTodos[0], nil
}
