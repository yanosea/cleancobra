package model

import (
	"fmt"
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTodo(title string) *Todo {
	return &Todo{
		ID:        generateUUID(),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
