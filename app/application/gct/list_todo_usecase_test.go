package gct

import (
	"testing"

	todoDomain "github.com/yanosea/gct/app/domain/todo"
)

func TestListTodoUseCase_Run_Success(t *testing.T) {
	// Arrange
	mockRepo := todoDomain.NewMockTodoRepository()
	useCase := NewListTodoUseCase(mockRepo)

	// Add test todos
	todo1, _ := todoDomain.NewTodo("Todo 1")
	todo2, _ := todoDomain.NewTodo("Todo 2")
	todo2.Done = true
	mockRepo.AddTodo(todo1)
	mockRepo.AddTodo(todo2)

	// Act
	result, err := useCase.Run()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result, got nil")
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(result))
	}
	if result[0].Title != "Todo 1" {
		t.Errorf("Expected first todo title 'Todo 1', got %s", result[0].Title)
	}
	if result[1].Done != true {
		t.Errorf("Expected second todo to be done, got %v", result[1].Done)
	}
}

func TestListTodoUseCase_Run_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := todoDomain.NewMockTodoRepository()
	useCase := NewListTodoUseCase(mockRepo)

	// Act
	result, err := useCase.Run()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result, got nil")
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(result))
	}
}

func TestListTodoUseCase_Run_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := todoDomain.NewMockTodoRepository()
	mockRepo.SetFindError(todoDomain.TodoNotFoundError{ID: "test"})
	useCase := NewListTodoUseCase(mockRepo)

	// Act
	result, err := useCase.Run()

	// Assert
	if err == nil {
		t.Error("Expected repository error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result for error case")
	}
}
