package application

import (
	"github.com/yanosea/gct/app/domain"
)

// ListTodoUseCase handles the business logic for retrieving all todos
type ListTodoUseCase struct {
	repository domain.TodoRepository
}

// NewListTodoUseCase creates a new ListTodoUseCase instance
func NewListTodoUseCase(repository domain.TodoRepository) *ListTodoUseCase {
	return &ListTodoUseCase{
		repository: repository,
	}
}

// Run executes the list todos use case
func (uc *ListTodoUseCase) Run() ([]domain.Todo, error) {
	// Retrieve all todos from repository
	todos, err := uc.repository.FindAll()
	if err != nil {
		return nil, err
	}

	// Return todos (empty slice if no todos exist)
	return todos, nil
}