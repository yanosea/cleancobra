package todo

// MockTodoRepository is a mock implementation of TodoRepository for testing
type MockTodoRepository struct {
	todos   []*Todo
	saveErr error
	findErr error
}

// NewMockTodoRepository creates a new mock repository
func NewMockTodoRepository() *MockTodoRepository {
	return &MockTodoRepository{
		todos: make([]*Todo, 0),
	}
}

// SetSaveError sets an error to be returned by Save method
func (m *MockTodoRepository) SetSaveError(err error) {
	m.saveErr = err
}

// SetFindError sets an error to be returned by Find methods
func (m *MockTodoRepository) SetFindError(err error) {
	m.findErr = err
}

// AddTodo adds a todo to the mock repository
func (m *MockTodoRepository) AddTodo(todo *Todo) {
	m.todos = append(m.todos, todo)
}

// ClearTodos clears all todos from the mock repository
func (m *MockTodoRepository) ClearTodos() {
	m.todos = make([]*Todo, 0)
}

func (m *MockTodoRepository) Save(todo *Todo) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.todos = append(m.todos, todo)
	return nil
}

func (m *MockTodoRepository) FindAll() ([]*Todo, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.todos, nil
}

func (m *MockTodoRepository) FindByID(id string) (*Todo, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for _, todo := range m.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return nil, TodoNotFoundError{ID: id}
}

func (m *MockTodoRepository) Update(todo *Todo) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	for i, t := range m.todos {
		if t.ID == todo.ID {
			m.todos[i] = todo
			return nil
		}
	}
	return TodoNotFoundError{ID: todo.ID}
}

func (m *MockTodoRepository) Delete(id string) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return TodoNotFoundError{ID: id}
}

func (m *MockTodoRepository) FindByQuery(query TodoQuery) ([]*Todo, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}

	todos := make([]*Todo, 0)
	for _, todo := range m.todos {
		if query.Done == nil || todo.Done == *query.Done {
			todos = append(todos, todo)
		}
	}

	// Apply pagination
	if query.Offset > 0 && query.Offset < len(todos) {
		todos = todos[query.Offset:]
	}
	if query.Limit > 0 && query.Limit < len(todos) {
		todos = todos[:query.Limit]
	}

	return todos, nil
}

func (m *MockTodoRepository) Count() (int, error) {
	if m.findErr != nil {
		return 0, m.findErr
	}
	return len(m.todos), nil
}

func (m *MockTodoRepository) CountByQuery(query TodoQuery) (int, error) {
	if m.findErr != nil {
		return 0, m.findErr
	}

	count := 0
	for _, todo := range m.todos {
		if query.Done == nil || todo.Done == *query.Done {
			count++
		}
	}
	return count, nil
}
