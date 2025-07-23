package domain

import (
	"errors"
	"testing"
)

// mockTodoRepository is a simple mock implementation for testing the interface
type mockTodoRepository struct {
	todos   []Todo
	nextID  int
	findErr error
	saveErr error
	delErr  error
}

func newMockTodoRepository() *mockTodoRepository {
	return &mockTodoRepository{
		todos:  make([]Todo, 0),
		nextID: 1,
	}
}

func (m *mockTodoRepository) FindAll() ([]Todo, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.todos, nil
}

func (m *mockTodoRepository) Save(todos ...Todo) ([]Todo, error) {
	if m.saveErr != nil {
		return nil, m.saveErr
	}

	saved := make([]Todo, 0, len(todos))
	for _, todo := range todos {
		if todo.ID == 0 {
			// New todo - assign ID
			todo.ID = m.nextID
			m.nextID++
			m.todos = append(m.todos, todo)
		} else {
			// Update existing todo
			found := false
			for i, existing := range m.todos {
				if existing.ID == todo.ID {
					m.todos[i] = todo
					found = true
					break
				}
			}
			if !found {
				return nil, ErrTodoNotFound
			}
		}
		saved = append(saved, todo)
	}
	return saved, nil
}

func (m *mockTodoRepository) DeleteById(id int) error {
	if m.delErr != nil {
		return m.delErr
	}

	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return ErrTodoNotFound
}

func TestTodoRepository_Interface(t *testing.T) {
	// Test that our mock implements the interface
	var repo TodoRepository = newMockTodoRepository()

	// Test FindAll
	todos, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v, want nil", err)
	}
	if len(todos) != 0 {
		t.Errorf("FindAll() returned %d todos, want 0", len(todos))
	}

	// Test Save with new todo
	newTodo := Todo{
		ID:          0, // New todo
		Description: "Test todo",
		Done:        false,
	}

	saved, err := repo.Save(newTodo)
	if err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}
	if len(saved) != 1 {
		t.Errorf("Save() returned %d todos, want 1", len(saved))
	}
	if saved[0].ID == 0 {
		t.Errorf("Save() did not assign ID to new todo")
	}

	// Test FindAll after save
	todos, err = repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v, want nil", err)
	}
	if len(todos) != 1 {
		t.Errorf("FindAll() returned %d todos, want 1", len(todos))
	}

	// Test Save with existing todo (update)
	existingTodo := saved[0]
	existingTodo.Description = "Updated description"
	existingTodo.Done = true

	updated, err := repo.Save(existingTodo)
	if err != nil {
		t.Errorf("Save() update error = %v, want nil", err)
	}
	if len(updated) != 1 {
		t.Errorf("Save() update returned %d todos, want 1", len(updated))
	}
	if updated[0].Description != "Updated description" {
		t.Errorf("Save() update description = %v, want 'Updated description'", updated[0].Description)
	}
	if updated[0].Done != true {
		t.Errorf("Save() update done = %v, want true", updated[0].Done)
	}

	// Test DeleteById with existing todo
	err = repo.DeleteById(saved[0].ID)
	if err != nil {
		t.Errorf("DeleteById() error = %v, want nil", err)
	}

	// Test FindAll after delete
	todos, err = repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v, want nil", err)
	}
	if len(todos) != 0 {
		t.Errorf("FindAll() returned %d todos after delete, want 0", len(todos))
	}

	// Test DeleteById with non-existent todo
	err = repo.DeleteById(999)
	if err != ErrTodoNotFound {
		t.Errorf("DeleteById() with non-existent ID error = %v, want ErrTodoNotFound", err)
	}
}

func TestTodoRepository_ErrorHandling(t *testing.T) {
	mock := newMockTodoRepository()

	// Test FindAll error
	mock.findErr = errors.New("find error")
	_, err := mock.FindAll()
	if err == nil {
		t.Errorf("FindAll() with error should return error")
	}

	// Test Save error
	mock.findErr = nil
	mock.saveErr = errors.New("save error")
	_, err = mock.Save(Todo{ID: 0, Description: "test"})
	if err == nil {
		t.Errorf("Save() with error should return error")
	}

	// Test DeleteById error
	mock.saveErr = nil
	mock.delErr = errors.New("delete error")
	err = mock.DeleteById(1)
	if err == nil {
		t.Errorf("DeleteById() with error should return error")
	}
}

func TestTodoRepository_SaveMultiple(t *testing.T) {
	repo := newMockTodoRepository()

	// Test saving multiple todos at once
	todos := []Todo{
		{ID: 0, Description: "Todo 1", Done: false},
		{ID: 0, Description: "Todo 2", Done: true},
		{ID: 0, Description: "Todo 3", Done: false},
	}

	saved, err := repo.Save(todos...)
	if err != nil {
		t.Errorf("Save() multiple error = %v, want nil", err)
	}

	if len(saved) != 3 {
		t.Errorf("Save() multiple returned %d todos, want 3", len(saved))
	}

	// Verify all todos got unique IDs
	ids := make(map[int]bool)
	for _, todo := range saved {
		if todo.ID == 0 {
			t.Errorf("Save() multiple did not assign ID to todo")
		}
		if ids[todo.ID] {
			t.Errorf("Save() multiple assigned duplicate ID %d", todo.ID)
		}
		ids[todo.ID] = true
	}

	// Verify all todos are in repository
	allTodos, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v, want nil", err)
	}
	if len(allTodos) != 3 {
		t.Errorf("FindAll() returned %d todos, want 3", len(allTodos))
	}
}

func TestTodoRepository_UpdateNonExistent(t *testing.T) {
	repo := newMockTodoRepository()

	// Try to update a todo that doesn't exist
	nonExistentTodo := Todo{
		ID:          999,
		Description: "Non-existent todo",
		Done:        false,
	}

	_, err := repo.Save(nonExistentTodo)
	if err != ErrTodoNotFound {
		t.Errorf("Save() with non-existent todo error = %v, want ErrTodoNotFound", err)
	}
}
