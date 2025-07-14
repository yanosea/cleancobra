package gct

import (
	todoDomain "github.com/yanosea/gct/app/domain/todo"
)

type ListTodoUseCase struct {
	todoRepo todoDomain.TodoRepository
}

func NewListTodoUseCase(
	todoRepo todoDomain.TodoRepository,
) *ListTodoUseCase {
	return &ListTodoUseCase{
		todoRepo: todoRepo,
	}
}

type ListTodoUsecaseOutputDto struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt string
}

func (uc *ListTodoUseCase) Run() ([]*ListTodoUsecaseOutputDto, error) {
	todo, err := uc.todoRepo.FindAll()
	if err != nil {
		return nil, err
	}
	todoDto := make([]*ListTodoUsecaseOutputDto, len(todo))
	for i, t := range todo {
		todoDto[i] = &ListTodoUsecaseOutputDto{
			ID:        t.ID,
			Title:     t.Title,
			Done:      t.Done,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return todoDto, nil
}
