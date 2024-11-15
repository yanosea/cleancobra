package json

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cleancobra/app/domain/model"
	"cleancobra/app/domain/repository"
	"cleancobra/pkg/errors"
)

const (
	appName  = "cleancobra"
	todoFile = "todos.json"
)

type TodoRepository struct {
	filePath string
}

func NewTodoRepository() repository.TodoRepository {
	dataDir := getDataDir()
	return &TodoRepository{
		filePath: filepath.Join(dataDir, todoFile),
	}
}

func (r *TodoRepository) Save(todo *model.Todo) error {
	todos, err := r.FindAll()
	if err != nil {
		return err
	}
	todos = append(todos, todo)
	return r.writeTodos(todos)
}

func (r *TodoRepository) FindAll() ([]*model.Todo, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []*model.Todo{}, nil
	}

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var todos []*model.Todo
	if err := json.Unmarshal(file, &todos); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal todos")
	}

	return todos, nil
}

func (r *TodoRepository) FindByID(id string) (*model.Todo, error) {
	todos, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	for _, todo := range todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return nil, errors.New("todo not found")
}

func (r *TodoRepository) Update(todo *model.Todo) error {
	todos, err := r.FindAll()
	if err != nil {
		return err
	}

	for i, t := range todos {
		if t.ID == todo.ID {
			todos[i] = todo
			return r.writeTodos(todos)
		}
	}

	return errors.New("todo not found")
}

func (r *TodoRepository) Delete(id string) error {
	todos, err := r.FindAll()
	if err != nil {
		return err
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return r.writeTodos(todos)
		}
	}

	return errors.New("todo not found")
}

func getDataDir() string {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		xdgDataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdgDataHome, appName)
}

func (r *TodoRepository) writeTodos(todos []*model.Todo) error {
	if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	file, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal todos")
	}

	if err := os.WriteFile(r.filePath, file, 0644); err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}
