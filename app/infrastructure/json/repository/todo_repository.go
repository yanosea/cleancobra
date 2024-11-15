package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/yanosea/cleancobra/config"
	todoDomain "github.com/yanosea/cleancobra/domain/todo"

	"github.com/yanosea/cleancobra-pkg/errors"
)

const (
	dbFileName = "todos.json"
)

type TodoRepository struct {
	dbFilePath string
}

func NewTodoRepository(conf *config.TodoConfig) todoDomain.TodoRepository {
	dbFileDirPath := getDBDirPath(conf.DBDirPath)
	mkdirIfNotExist(dbFileDirPath)
	dbFilePath := filepath.Join(dbFileDirPath, dbFileName)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		initializeDBFile(dbFilePath)
	}
	return &TodoRepository{
		dbFilePath: dbFilePath,
	}
}

func (r *TodoRepository) Save(todo *todoDomain.Todo) error {
	todos, err := r.FindAll()
	if err != nil {
		return err
	}
	todos = append(todos, todo)
	return r.writeTodos(todos)
}

func (r *TodoRepository) FindAll() ([]*todoDomain.Todo, error) {
	file, err := os.ReadFile(r.dbFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var todos []*todoDomain.Todo
	if err := json.Unmarshal(file, &todos); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal todos")
	}

	return todos, nil
}

func (r *TodoRepository) FindByID(id string) (*todoDomain.Todo, error) {
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

func (r *TodoRepository) Update(todo *todoDomain.Todo) error {
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

func getDBDirPath(dbDirPathConfig string) string {
	if dbDirPathConfig == "" {
		dbDirPathConfig = "XDG_DATA_HOME/todos"
	}

	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "~"
		}
		xdgDataHome = filepath.Join(homeDir, ".local", "share")
	}

	dataDir := strings.Replace(dbDirPathConfig, "XDG_DATA_HOME", xdgDataHome, 1)

	return dataDir
}

func (r *TodoRepository) writeTodos(todos []*todoDomain.Todo) error {
	if err := os.MkdirAll(filepath.Dir(r.dbFilePath), 0755); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	file, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal todos")
	}

	if err := os.WriteFile(r.dbFilePath, file, 0644); err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}

func mkdirIfNotExist(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			panic(err)
		}
	}
}

func initializeDBFile(dbFilePath string) {
	emptyTodos := []*todoDomain.Todo{}
	file, err := json.MarshalIndent(emptyTodos, "", "  ")
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal empty todos"))
	}

	if err := os.WriteFile(dbFilePath, file, 0644); err != nil {
		panic(errors.Wrap(err, "failed to create initial todos file"))
	}
}
