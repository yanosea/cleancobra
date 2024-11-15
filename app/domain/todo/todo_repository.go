package todo

type TodoRepository interface {
	Save(todo *Todo) error
	FindAll() ([]*Todo, error)
	FindByID(id string) (*Todo, error)
	Update(todo *Todo) error
	Delete(id string) error
}
