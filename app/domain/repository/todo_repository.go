package repository

import "cleancobra/app/domain/model"

type TodoRepository interface {
	Save(todo *model.Todo) error
	FindAll() ([]*model.Todo, error)
	FindByID(id string) (*model.Todo, error)
	Update(todo *model.Todo) error
	Delete(id string) error
}
