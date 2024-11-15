package adapter

import (
	todoApp "cleancobra/app/application/usecase/todo"
	"cleancobra/app/infrastructure/repository/json"
)

type TodoUsecases struct {
	AddTodoUseCase    todoApp.AddTodoUseCase
	DeleteTodoUseCase todoApp.DeleteTodoUseCase
	ListTodoUseCase   todoApp.ListTodoUseCase
	ToggleTodoUseCase todoApp.ToggleTodoUseCase
}

func NewTodoUseCases() *TodoUsecases {
	repo := json.NewTodoRepository()
	return &TodoUsecases{
		AddTodoUseCase:    todoApp.NewAddTodoUseCase(repo),
		DeleteTodoUseCase: todoApp.NewDeleteTodoUseCase(repo),
		ListTodoUseCase:   todoApp.NewListTodoUseCase(repo),
		ToggleTodoUseCase: todoApp.NewToggleTodoUseCase(repo),
	}
}
