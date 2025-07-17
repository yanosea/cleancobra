package application

import (
	"github.com/yanosea/gct/app/domain"
)

// DeleteTodoUseCase handles the business logic for deleting todos
type DeleteTodoUseCase struct {
	repository domain.TodoRepository
}

// NewDeleteTodoUseCase creates a new DeleteTodoUseCase instance
func NewDeleteTodoUseCase(repository domain.TodoRepository) *DeleteTodoUseCase {
	return &DeleteTodoUseCase{
		repository: repository,
	}
}

// Run executes the delete todo use case
func (uc *DeleteTodoUseCase) Run(id int) error {
	// Validate input
	if id <= 0 {
		return domain.NewDomainError(
			domain.ErrorTypeInvalidInput,
			"todo ID must be positive",
			nil,
		)
	}

	// Delete the todo by ID
	err := uc.repository.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}