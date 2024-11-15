package todo

import (
	"cleancobra/app/domain/repository"
)

type ToggleTodoUseCase interface {
	Run(id string) error
}

type toggleTodoUseCase struct {
	repo repository.TodoRepository
}

func NewToggleTodoUseCase(repo repository.TodoRepository) ToggleTodoUseCase {
	return &toggleTodoUseCase{repo: repo}
}

func (uc *toggleTodoUseCase) Run(id string) error {
	todo, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}
	todo.Done = !todo.Done
	return uc.repo.Update(todo)
}
