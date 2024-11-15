package todo

import (
	todoDomain "github.com/yanosea/cleancobra/app/domain/todo"
)

type DeleteTodoUseCase struct {
	todoRepo todoDomain.TodoRepository
}

func NewDeleteTodoUseCase(
	todoRepo todoDomain.TodoRepository,
) *DeleteTodoUseCase {
	return &DeleteTodoUseCase{
		todoRepo: todoRepo,
	}
}

type DeleteTodoUsecaseOutputDto struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt string
}

func (uc *DeleteTodoUseCase) Run(id string) (*DeleteTodoUsecaseOutputDto, error) {
	todo, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := uc.todoRepo.Delete(id); err != nil {
		return nil, err
	}
	return &DeleteTodoUsecaseOutputDto{
		ID:        todo.ID,
		Title:     todo.Title,
		Done:      todo.Done,
		CreatedAt: todo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
