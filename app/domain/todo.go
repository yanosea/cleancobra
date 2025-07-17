package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Todo represents a todo item in the domain
type Todo struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTodo creates a new Todo with the given description
func NewTodo(id int, description string) (*Todo, error) {
	if err := validateDescription(description); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Todo{
		ID:          id,
		Description: strings.TrimSpace(description),
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Toggle toggles the completion status of the todo
func (t *Todo) Toggle() {
	t.Done = !t.Done
	t.UpdatedAt = time.Now()
}

// UpdateDescription updates the description of the todo
func (t *Todo) UpdateDescription(description string) error {
	if err := validateDescription(description); err != nil {
		return err
	}
	
	t.Description = strings.TrimSpace(description)
	t.UpdatedAt = time.Now()
	return nil
}

// Validate validates the todo entity
func (t *Todo) Validate() error {
	if t.ID <= 0 {
		return NewDomainError(ErrorTypeInvalidInput, "todo ID must be positive", nil)
	}
	
	if err := validateDescription(t.Description); err != nil {
		return err
	}
	
	if t.CreatedAt.IsZero() {
		return NewDomainError(ErrorTypeInvalidInput, "created_at cannot be zero", nil)
	}
	
	if t.UpdatedAt.IsZero() {
		return NewDomainError(ErrorTypeInvalidInput, "updated_at cannot be zero", nil)
	}
	
	if t.UpdatedAt.Before(t.CreatedAt) {
		return NewDomainError(ErrorTypeInvalidInput, "updated_at cannot be before created_at", nil)
	}
	
	return nil
}

// validateDescription validates the todo description
func validateDescription(description string) error {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return NewDomainError(ErrorTypeInvalidInput, "description cannot be empty", nil)
	}
	
	if len(trimmed) > 500 {
		return NewDomainError(ErrorTypeInvalidInput, "description cannot exceed 500 characters", nil)
	}
	
	return nil
}

// String returns a string representation of the todo
func (t *Todo) String() string {
	status := "[ ]"
	if t.Done {
		status = "[x]"
	}
	return fmt.Sprintf("%s %d: %s", status, t.ID, t.Description)
}

// MarshalJSON implements custom JSON marshaling
func (t *Todo) MarshalJSON() ([]byte, error) {
	type Alias Todo
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(t),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling
func (t *Todo) UnmarshalJSON(data []byte) error {
	type Alias Todo
	aux := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias: (*Alias)(t),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return NewDomainError(ErrorTypeJSON, "failed to unmarshal todo", err)
	}
	
	var err error
	t.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return NewDomainError(ErrorTypeJSON, "invalid created_at format", err)
	}
	
	t.UpdatedAt, err = time.Parse(time.RFC3339, aux.UpdatedAt)
	if err != nil {
		return NewDomainError(ErrorTypeJSON, "invalid updated_at format", err)
	}
	
	return t.Validate()
}