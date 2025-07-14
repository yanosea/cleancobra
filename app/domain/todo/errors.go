package todo

import "fmt"

// TodoNotFoundError represents an error when a todo is not found
type TodoNotFoundError struct {
	ID string
}

func (e TodoNotFoundError) Error() string {
	return fmt.Sprintf("todo with ID %s not found", e.ID)
}

// InvalidTodoError represents an error when todo data is invalid
type InvalidTodoError struct {
	Field   string
	Message string
}

func (e InvalidTodoError) Error() string {
	return fmt.Sprintf("invalid todo %s: %s", e.Field, e.Message)
}

// TodoAlreadyExistsError represents an error when trying to create a duplicate todo
type TodoAlreadyExistsError struct {
	ID string
}

func (e TodoAlreadyExistsError) Error() string {
	return fmt.Sprintf("todo with ID %s already exists", e.ID)
}
