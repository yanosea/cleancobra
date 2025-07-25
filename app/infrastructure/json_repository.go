package infrastructure

import (
	"github.com/yanosea/gct/app/domain"

	"github.com/yanosea/gct/pkg/proxy"
)

// JSONRepository implements TodoRepository interface using JSON file storage
type JSONRepository struct {
	filePath      string
	filepathProxy proxy.Filepath
	jsonProxy     proxy.JSON
	osProxy       proxy.OS
	sortProxy     proxy.Sort
}

// NewJSONRepository creates a new JSONRepository instance
func NewJSONRepository(
	filePath string,
	filepathProxy proxy.Filepath,
	jsonProxy proxy.JSON,
	osProxy proxy.OS,
	sortProxy proxy.Sort,
) *JSONRepository {
	return &JSONRepository{
		filePath:      filePath,
		filepathProxy: filepathProxy,
		jsonProxy:     jsonProxy,
		osProxy:       osProxy,
		sortProxy:     sortProxy,
	}
}

// FindAll retrieves all todos from the JSON file
func (r *JSONRepository) FindAll() ([]domain.Todo, error) {
	// check if file exists
	if _, err := r.osProxy.Stat(r.filePath); r.osProxy.IsNotExist(err) {
		// file doesn't exist, return empty slice
		return []domain.Todo{}, nil
	} else if err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to check file status",
			err,
		)
	}

	// read file content
	file, err := r.osProxy.OpenFile(r.filePath, proxy.ORdOnly, 0644)
	if err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to open file for reading",
			err,
		)
	}
	defer file.Close()

	// decode JSON
	var todos []domain.Todo
	if err := r.jsonProxy.NewDecoder(file).Decode(&todos); err != nil {
		return nil, domain.NewDomainError(
			domain.ErrorTypeJSON,
			"failed to decode JSON data",
			err,
		)
	}

	// sort todos by ID for consistent ordering
	r.sortProxy.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos, nil
}

// Save saves one or more todos to the JSON file
func (r *JSONRepository) Save(todos ...domain.Todo) ([]domain.Todo, error) {
	if len(todos) == 0 {
		return []domain.Todo{}, nil
	}

	// load existing todos
	existingTodos, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	// create a map for quick lookup of existing todos
	todoMap := make(map[int]domain.Todo)
	for _, todo := range existingTodos {
		todoMap[todo.ID] = todo
	}

	// find the next available ID
	nextID := r.getNextID(existingTodos)

	var savedTodos []domain.Todo
	for _, todo := range todos {
		if todo.ID == 0 {
			// new todo - assign ID
			todo.ID = nextID
			nextID++
		} else {
			// update existing todo - check if it exists
			if _, exists := todoMap[todo.ID]; !exists {
				return nil, domain.ErrTodoNotFound
			}
		}

		// add or update in the map
		todoMap[todo.ID] = todo
		savedTodos = append(savedTodos, todo)
	}

	// convert map back to slice
	allTodos := make([]domain.Todo, 0, len(todoMap))
	for _, todo := range todoMap {
		allTodos = append(allTodos, todo)
	}

	// sort by ID for consistent ordering
	r.sortProxy.Slice(allTodos, func(i, j int) bool {
		return allTodos[i].ID < allTodos[j].ID
	})

	// write to file
	if err := r.writeToFile(allTodos); err != nil {
		return nil, err
	}

	return savedTodos, nil
}

// DeleteById removes a todo with the specified ID from the JSON file
func (r *JSONRepository) DeleteById(id int) error {
	// load existing todos
	existingTodos, err := r.FindAll()
	if err != nil {
		return err
	}

	// find and remove the todo
	var updatedTodos []domain.Todo
	found := false
	for _, todo := range existingTodos {
		if todo.ID == id {
			found = true
			continue // skip this todo (delete it)
		}
		updatedTodos = append(updatedTodos, todo)
	}

	if !found {
		return domain.ErrTodoNotFound
	}

	// resequence IDs to maintain consecutive numbering
	for i := range updatedTodos {
		updatedTodos[i].ID = i + 1
	}

	// write updated todos back to file
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
	// ensure directory exists
	fileDir := r.filepathProxy.Dir(r.filePath)
	if err := r.osProxy.MkdirAll(fileDir, 0755); err != nil {
		return domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to create directory",
			err,
		)
	}

	// create or truncate file
	file, err := r.osProxy.OpenFile(r.filePath, proxy.OCreate|proxy.OWrOnly|proxy.OTrunc, 0644)
	if err != nil {
		return domain.NewDomainError(
			domain.ErrorTypeFileSystem,
			"failed to open file for writing",
			err,
		)
	}
	defer file.Close()

	// encode JSON with indentation for readability
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
