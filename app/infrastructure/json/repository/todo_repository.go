package repository

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/yanosea/gct/app/config"
	todoDomain "github.com/yanosea/gct/app/domain/todo"

	"github.com/yanosea/gct/pkg/proxy"
	"github.com/yanosea/gct/pkg/utility"
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
		return nil, err
	}

	var todos []*todoDomain.Todo
	if err := r.json.Unmarshal(file, &todos); err != nil {
		return nil, err
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

	return nil, todoDomain.TodoNotFoundError{ID: id}
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

	return todoDomain.TodoNotFoundError{ID: todo.ID}
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

	return todoDomain.TodoNotFoundError{ID: id}
}

func (r *TodoRepository) FindByQuery(query todoDomain.TodoQuery) ([]*todoDomain.Todo, error) {
	todos, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	// Filter by completion status
	if query.Done != nil {
		filtered := make([]*todoDomain.Todo, 0)
		for _, todo := range todos {
			if todo.Done == *query.Done {
				filtered = append(filtered, todo)
			}
		}
		todos = filtered
	}

	// Sort todos
	r.sortTodos(todos, query.SortBy, query.Order)

	// Apply pagination
	if query.Offset > 0 && query.Offset < len(todos) {
		todos = todos[query.Offset:]
	}
	if query.Limit > 0 && query.Limit < len(todos) {
		todos = todos[:query.Limit]
	}

	return todos, nil
}

func (r *TodoRepository) Count() (int, error) {
	todos, err := r.FindAll()
	if err != nil {
		return 0, err
	}
	return len(todos), nil
}

func (r *TodoRepository) CountByQuery(query todoDomain.TodoQuery) (int, error) {
	todos, err := r.FindAll()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, todo := range todos {
		if query.Done == nil || todo.Done == *query.Done {
			count++
		}
	}

	return count, nil
}

func (r *TodoRepository) sortTodos(todos []*todoDomain.Todo, sortBy, order string) {
	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	sort.Slice(todos, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "title":
			less = strings.ToLower(todos[i].Title) < strings.ToLower(todos[j].Title)
		case "created_at":
			less = todos[i].CreatedAt.Before(todos[j].CreatedAt)
		default:
			less = todos[i].CreatedAt.Before(todos[j].CreatedAt)
		}

		if order == "desc" {
			return !less
		}
		return less
	})
}

func (r *TodoRepository) writeTodos(todos []*todoDomain.Todo) error {
	if err := r.os.MkdirAll(filepath.Dir(r.dbFilePath), 0755); err != nil {
		return err
	}

	file, err := r.json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	if err := r.os.WriteFile(r.dbFilePath, file, 0644); err != nil {
		return err
	}

	return nil
}
