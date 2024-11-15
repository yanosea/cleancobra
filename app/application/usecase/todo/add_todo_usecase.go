package todo

import (
	"cleancobra/app/domain/model"
	"cleancobra/app/domain/repository"
	"cleancobra/pkg/errors"
)

type AddTodoUseCase interface {
	Run(title string) error
}

type addTodoUseCase struct {
	repo repository.TodoRepository
}

func NewAddTodoUseCase(repo repository.TodoRepository) AddTodoUseCase {
	return &addTodoUseCase{repo: repo}
}

func (uc *addTodoUseCase) Run(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	todo := model.NewTodo(title)
	return uc.repo.Save(todo)
}
