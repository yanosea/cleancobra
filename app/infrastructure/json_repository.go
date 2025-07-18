package infrastructure

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
)

// JSONRepository implements TodoRepository interface using JSON file storage
type JSONRepository struct {
	filePath  string
	osProxy   proxy.OS
	jsonProxy proxy.JSON
}

// NewJSONRepository creates a new JSONRepository instance
func NewJSONRepository(filePath string, osProxy proxy.OS, jsonProxy proxy.JSON) *JSONRepository {
	return &JSONRepository{
		filePath:  filePath,
		osProxy:   osProxy,
		jsonProxy: jsonProxy,
	}
}

// FindAll retrieves all todos from the JSON file
func (r *JSONRepository) FindAll() ([]domain.Todo, error) {
	// Check if file exists
	if _, err := r.osProxy.Stat(r.filePath); r.osProxy.IsNotExist(err) {
		// File doesn't exist, return empty slice
		return []domain.Todo{}, nil
	} else if err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to check file status",
			err,
		)
	}

	// Read file content
	file, err := r.osProxy.OpenFile(r.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to open file for reading",
			err,
		)
	}
	defer file.Close()

	// Decode JSON
	var todos []domain.Todo
	if err := r.jsonProxy.NewDecoder(file).Decode(&todos); err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeJSON,
			"failed to decode JSON data",
			err,
		)
	}

	// Sort todos by ID for consistent ordering
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos, nil
}

// Save saves one or more todos to the JSON file
func (r *JSONRepository) Save(todos ...domain.Todo) ([]domain.Todo, error) {
	if len(todos) == 0 {
		return []domain.Todo{}, nil
	}

	// Load existing todos
	existingTodos, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup of existing todos
	todoMap := make(map[int]domain.Todo)
	for _, todo := range existingTodos {
		todoMap[todo.ID] = todo
	}

	// Find the next available ID
	nextID := r.getNextID(existingTodos)

	var savedTodos []domain.Todo
	for _, todo := range todos {
		if todo.ID == 0 {
			// New todo - assign ID
			todo.ID = nextID
			nextID++
		} else {
			// Update existing todo - check if it exists
			if _, exists := todoMap[todo.ID]; !exists {
				return nil, domain.ErrTodoNotFound
			}
		}

		// Add or update in the map
		todoMap[todo.ID] = todo
		savedTodos = append(savedTodos, todo)
	}

	// Convert map back to slice
	allTodos := make([]domain.Todo, 0, len(todoMap))
	for _, todo := range todoMap {
		allTodos = append(allTodos, todo)
	}

	// Sort by ID for consistent ordering
	sort.Slice(allTodos, func(i, j int) bool {
		return allTodos[i].ID < allTodos[j].ID
	})

	// Write to file
	if err := r.writeToFile(allTodos); err != nil {
		return nil, err
	}

	return savedTodos, nil
}

// DeleteById removes a todo with the specified ID from the JSON file
func (r *JSONRepository) DeleteById(id int) error {
	// Load existing todos
	existingTodos, err := r.FindAll()
	if err != nil {
		return err
	}

	// Find and remove the todo
	var updatedTodos []domain.Todo
	found := false
	for _, todo := range existingTodos {
		if todo.ID == id {
			found = true
			continue // Skip this todo (delete it)
		}
		updatedTodos = append(updatedTodos, todo)
	}

	if !found {
		return domain.ErrTodoNotFound
	}

	// Resequence IDs to maintain consecutive numbering
	for i := range updatedTodos {
		updatedTodos[i].ID = i + 1
	}

	// Write updated todos back to file
	return r.writeToFile(updatedTodos)
}

// getNextID finds the next available ID
func (r *JSONRepository) getNextID(todos []domain.Todo) int {
	if len(todos) == 0 {
		return 1
	}

	maxID := 0
	for _, todo := range todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}

	return maxID + 1
}

// writeToFile writes todos to the JSON file
func (r *JSONRepository) writeToFile(todos []domain.Todo) error {
	// Ensure directory exists
	fileDir := filepath.Dir(r.filePath)
	if err := r.osProxy.MkdirAll(fileDir, 0755); err != nil {
		return domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to create directory",
			err,
		)
	}

	// Create or truncate file
	file, err := r.osProxy.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to open file for writing",
			err,
		)
	}
	defer file.Close()

	// Encode JSON with indentation for readability
	encoder := r.jsonProxy.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(todos); err != nil {
		return domain.NewDomainError(
			domain.ErrorTypeJSON,
			"failed to encode JSON data",
			err,
		)
	}

	return nil
}
