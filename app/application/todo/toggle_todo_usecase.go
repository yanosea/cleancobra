package todo

import (
	todoDomain "github.com/yanosea/gct/app/domain/todo"
)

type ToggleTodoUseCase struct {
	todoRepo todoDomain.TodoRepository
}

func NewToggleTodoUseCase(
	todoRepo todoDomain.TodoRepository,
) *ToggleTodoUseCase {
	return &ToggleTodoUseCase{
		todoRepo: todoRepo,
	}
}

type ToggleTodoUsecaseOutputDto struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt string
}

func (uc *ToggleTodoUseCase) Run(id string) (*ToggleTodoUsecaseOutputDto, error) {
	todo, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	todo.Done = !todo.Done
	if err := uc.todoRepo.Update(todo); err != nil {
		return nil, err
	}
	return &ToggleTodoUsecaseOutputDto{
		ID:        todo.ID,
		Title:     todo.Title,
		Done:      todo.Done,
		CreatedAt: todo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
