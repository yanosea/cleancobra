package gct

import (
	"testing"

	todoDomain "github.com/yanosea/gct/app/domain/todo"
)

func TestAddTodoUseCase_Run_Success(t *testing.T) {
	// Arrange
	mockRepo := todoDomain.NewMockTodoRepository()
	useCase := NewAddTodoUseCase(mockRepo)
	title := "Test Todo"

	// Act
	result, err := useCase.Run(title)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result, got nil")
	}
	if result.Title != title {
		t.Errorf("Expected title %s, got %s", title, result.Title)
	}

	// Verify todo was saved
	todos, _ := mockRepo.FindAll()
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}
}

func TestAddTodoUseCase_Run_EmptyTitle(t *testing.T) {
	// Arrange
	mockRepo := todoDomain.NewMockTodoRepository()
	useCase := NewAddTodoUseCase(mockRepo)

	// Act
	result, err := useCase.Run("")

	// Assert
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
	if result != nil {
		t.Error("Expected nil result for error case")
	}

	// Verify no todo was saved
	todos, _ := mockRepo.FindAll()
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestAddTodoUseCase_Run_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := todoDomain.NewMockTodoRepository()
	mockRepo.SetSaveError(todoDomain.TodoAlreadyExistsError{ID: "test"})
	useCase := NewAddTodoUseCase(mockRepo)

	// Act
	result, err := useCase.Run("Test Todo")

	// Assert
	if err == nil {
		t.Error("Expected repository error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result for error case")
	}
}
