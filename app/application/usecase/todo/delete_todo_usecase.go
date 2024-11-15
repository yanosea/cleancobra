package todo

import (
	"cleancobra/app/domain/repository"
)

type DeleteTodoUseCase interface {
	Run(id string) error
}

type deleteTodoUseCase struct {
	repo repository.TodoRepository
}

func NewDeleteTodoUseCase(repo repository.TodoRepository) DeleteTodoUseCase {
	return &deleteTodoUseCase{repo: repo}
}

func (uc *deleteTodoUseCase) Run(id string) error {
	return uc.repo.Delete(id)
}
