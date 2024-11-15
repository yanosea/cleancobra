package mapper

import (
	"cleancobra/app/domain/model"
	"cleancobra/app/presentation/cli/adapter/dto"
)

func ToDTO(todos []*model.Todo) []dto.Todo {
	dtos := make([]dto.Todo, len(todos))
	for i, todo := range todos {
		dtos[i] = dto.Todo{
			ID:        todo.ID,
			Title:     todo.Title,
			Done:      todo.Done,
			CreatedAt: todo.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return dtos
}
