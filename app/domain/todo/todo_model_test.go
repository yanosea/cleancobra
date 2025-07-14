package todo

import (
	"testing"
)

func TestNewTodo_Success(t *testing.T) {
	// Arrange
	title := "Test Todo"

	// Act
	todo, err := NewTodo(title)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if todo == nil {
		t.Error("Expected todo, got nil")
	}
	if todo.Title != title {
		t.Errorf("Expected title %s, got %s", title, todo.Title)
	}
	if todo.Done != false {
		t.Errorf("Expected Done to be false, got %v", todo.Done)
	}
	if todo.ID == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestNewTodo_EmptyTitle(t *testing.T) {
	// Act
	todo, err := NewTodo("")

	// Assert
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
	if todo != nil {
		t.Error("Expected nil todo for error case")
	}

	// Check if it's the correct error type
	if invalidErr, ok := err.(InvalidTodoError); ok {
		if invalidErr.Field != "title" {
			t.Errorf("Expected field 'title', got %s", invalidErr.Field)
		}
	} else {
		t.Errorf("Expected InvalidTodoError, got %T", err)
	}
}

func TestTodoErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "TodoNotFoundError",
			err:      TodoNotFoundError{ID: "123"},
			expected: "todo with ID 123 not found",
		},
		{
			name:     "InvalidTodoError",
			err:      InvalidTodoError{Field: "title", Message: "cannot be empty"},
			expected: "invalid todo title: cannot be empty",
		},
		{
			name:     "TodoAlreadyExistsError",
			err:      TodoAlreadyExistsError{ID: "456"},
			expected: "todo with ID 456 already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message %s, got %s", tt.expected, tt.err.Error())
			}
		})
	}
}
