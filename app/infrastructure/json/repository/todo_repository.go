package repository

import (
	"path/filepath"
	"strings"

	"github.com/yanosea/cleancobra/app/config"
	todoDomain "github.com/yanosea/cleancobra/app/domain/todo"

	"github.com/yanosea/cleancobra/pkg/errors"
	"github.com/yanosea/cleancobra/pkg/proxy"
	"github.com/yanosea/cleancobra/pkg/utility"
)

const (
	dbFileName = "todos.json"
)

type TodoRepository struct {
	dbFilePath string
	json       proxy.Json
	os         proxy.Os
}

func NewTodoRepository(
	conf *config.TodoConfig,
	fileutil utility.FileUtil,
	json proxy.Json,
	os proxy.Os,
) (todoDomain.TodoRepository, error) {
	xdgDataHome, err := fileutil.GetXDGDataHome()
	if err != nil {
		return nil, err
	}
	dbFileDirPath := strings.Replace(conf.DBDirPath, "XDG_DATA_HOME", xdgDataHome, 1)
	if err := fileutil.MkdirIfNotExist(dbFileDirPath); err != nil {
		return nil, err
	}
	dbFilePath := filepath.Join(dbFileDirPath, dbFileName)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		if err := fileutil.InitializeJSONFile(dbFilePath, []*todoDomain.Todo{}); err != nil {
			return nil, err
		}
	}
	return &TodoRepository{
		dbFilePath: dbFilePath,
		json:       json,
		os:         os,
	}, nil
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
	file, err := r.os.ReadFile(r.dbFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var todos []*todoDomain.Todo
	if err := r.json.Unmarshal(file, &todos); err != nil {
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

func (r *TodoRepository) writeTodos(todos []*todoDomain.Todo) error {
	if err := r.os.MkdirAll(filepath.Dir(r.dbFilePath), 0755); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	file, err := r.json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal todos")
	}

	if err := r.os.WriteFile(r.dbFilePath, file, 0644); err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}
