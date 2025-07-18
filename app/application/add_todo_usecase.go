package application

import (
	"time"

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
	// Validate input
	if err := validateDescription(description); err != nil {
		return nil, err
	}

	// Create new todo with current timestamp
	now := time.Now()
	todo := domain.Todo{
		ID:          0, // Will be assigned by repository
		Description: description,
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Note: Skip validation for new todos (ID=0) as ID will be assigned by repository

	// Save to repository
	savedTodos, err := uc.repository.Save(todo)
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
