package application

import (
	"github.com/yanosea/gct/app/domain"
)

// AddTodoUseCase handles the business logic for adding new todos
type AddTodoUseCase struct {
	repository domain.TodoRepository
}

// NewAddTodoUseCase creates a new AddTodoUseCase instance
func NewAddTodoUseCase(repository domain.TodoRepository) *AddTodoUseCase {
	return &AddTodoUseCase{
		repository: repository,
	}
}

// Run executes the add todo use case
func (uc *AddTodoUseCase) Run(description string) (*domain.Todo, error) {
	// validate input
	if err := validateDescription(description); err != nil {
		return nil, err
	}

	// create new todo with current timestamp
	todo, err := domain.NewTodo(0, description)
	if err != nil {
		return nil, err
	}

	// save to repository
	savedTodos, err := uc.repository.Save(*todo)
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

// validateDescription validates the todo description
func validateDescription(description string) error {
	if description == "" {
		return domain.ErrEmptyDescription
	}

	if len(description) > 500 {
		return domain.NewDomainError(
			domain.ErrorTypeInvalidInput,
			"description cannot exceed 500 characters",
			nil,
		)
	}

	return nil
}
