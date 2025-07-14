package todo

import (
	"fmt"
	"time"

	"github.com/yanosea/gct/pkg/errors"
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
	now := time.Now()
	return &Todo{
		ID:        generateUUID(now),
		Title:     title,
		Done:      false,
		CreatedAt: now,
	}, nil
}

func generateUUID(now time.Time) string {
	return fmt.Sprintf("%d", now.UnixNano())
}
