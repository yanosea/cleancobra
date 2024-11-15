package todo

import (
	"fmt"
	"time"

	"github.com/yanosea/cleancobra-pkg/errors"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTodo(title string) (*Todo, error) {
	if title == "" {
		return nil, errors.New("title is empty")
	}
	return &Todo{
		ID:        generateUUID(),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}, nil
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
