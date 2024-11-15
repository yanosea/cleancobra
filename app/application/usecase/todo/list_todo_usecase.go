package todo

import (
	"cleancobra/app/domain/model"
	"cleancobra/app/domain/repository"
)

type ListTodoUseCase interface {
	Run() ([]*model.Todo, error)
}

type listTodoUseCase struct {
	repo repository.TodoRepository
}

func NewListTodoUseCase(repo repository.TodoRepository) ListTodoUseCase {
	return &listTodoUseCase{repo: repo}
}

func (uc *listTodoUseCase) Run() ([]*model.Todo, error) {
	return uc.repo.FindAll()
}
