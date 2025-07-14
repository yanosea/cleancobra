package todo

import (
	todoDomain "github.com/yanosea/gct/app/domain/todo"
)

type AddTodoUseCase struct {
	todoRepo todoDomain.TodoRepository
}

func NewAddTodoUseCase(
	todoRepo todoDomain.TodoRepository,
) *AddTodoUseCase {
	return &AddTodoUseCase{
		todoRepo: todoRepo,
	}
}

type AddTodoUsecaseOutputDto struct {
	ID        string
	Title     string
	CreatedAt string
}

func (uc *AddTodoUseCase) Run(title string) (*AddTodoUsecaseOutputDto, error) {
	todo, err := todoDomain.NewTodo(title)
	if err != nil {
		return nil, err
	}
	if err := uc.todoRepo.Save(todo); err != nil {
		return nil, err
	}
	return &AddTodoUsecaseOutputDto{
		ID:        todo.ID,
		Title:     todo.Title,
		CreatedAt: todo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
