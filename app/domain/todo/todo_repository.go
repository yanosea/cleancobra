package todo

// TodoQuery represents query parameters for finding todos
type TodoQuery struct {
	Done   *bool  // Filter by completion status
	Limit  int    // Limit number of results (0 = no limit)
	Offset int    // Offset for pagination
	SortBy string // Sort field ("created_at", "title")
	Order  string // Sort order ("asc", "desc")
}

type TodoRepository interface {
	Save(todo *Todo) error
	FindAll() ([]*Todo, error)
	FindByQuery(query TodoQuery) ([]*Todo, error)
	FindByID(id string) (*Todo, error)
	Update(todo *Todo) error
	Delete(id string) error
	Count() (int, error)
	CountByQuery(query TodoQuery) (int, error)
}
