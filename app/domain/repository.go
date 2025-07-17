package domain

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=domain

// TodoRepository defines the interface for todo persistence operations
type TodoRepository interface {
	// FindAll retrieves all todos from the repository
	FindAll() ([]Todo, error)
	
	// Save saves one or more todos to the repository
	// For new todos (ID = 0), it assigns new IDs and saves them
	// For existing todos (ID > 0), it updates them
	// Returns the saved todos with their assigned/updated IDs
	Save(todos ...Todo) ([]Todo, error)
	
	// DeleteById removes a todo with the specified ID from the repository
	// Returns ErrTodoNotFound if the todo doesn't exist
	DeleteById(id int) error
}