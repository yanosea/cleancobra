package domain

import (
	"time"

	"github.com/yanosea/gct/pkg/proxy"
)

// Todo represents a todo item in the domain
type Todo struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTodoWithDeps creates a new Todo with the given description using injected dependencies
func NewTodoWithDeps(id int, description string, timeProxy proxy.TimeProxy, stringsProxy proxy.Strings) (*Todo, error) {
	if err := validateDescriptionWithDeps(description, stringsProxy); err != nil {
		return nil, err
	}

	now := timeProxy.Now()
	return &Todo{
		ID:          id,
		Description: stringsProxy.TrimSpace(description),
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// NewTodo creates a new Todo with the given description
func NewTodo(id int, description string) (*Todo, error) {
	// Create real proxy implementations for backward compatibility
	timeProxy := proxy.NewTime()
	stringsProxy := proxy.NewStrings()

	// Call NewTodoWithDeps with real proxy instances
	return NewTodoWithDeps(id, description, timeProxy, stringsProxy)
}

// Toggle toggles the completion status of the todo
func (t *Todo) Toggle() {
	t.Done = !t.Done
	timeProxy := proxy.NewTime()
	t.UpdatedAt = timeProxy.Now()
}

// UpdateDescription updates the description of the todo
func (t *Todo) UpdateDescription(description string) error {
	if err := validateDescription(description); err != nil {
		return err
	}

	stringsProxy := proxy.NewStrings()
	t.Description = stringsProxy.TrimSpace(description)
	timeProxy := proxy.NewTime()
	t.UpdatedAt = timeProxy.Now()
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

// validateDescriptionWithDeps validates the todo description using injected dependencies
func validateDescriptionWithDeps(description string, stringsProxy proxy.Strings) error {
	trimmed := stringsProxy.TrimSpace(description)
	if trimmed == "" {
		return NewDomainError(ErrorTypeInvalidInput, "description cannot be empty", nil)
	}

	if len(trimmed) > 500 {
		return NewDomainError(ErrorTypeInvalidInput, "description cannot exceed 500 characters", nil)
	}

	return nil
}

// validateDescription validates the todo description
func validateDescription(description string) error {
	stringsProxy := proxy.NewStrings()
	return validateDescriptionWithDeps(description, stringsProxy)
}

// String returns a string representation of the todo
func (t *Todo) String() string {
	status := "[ ]"
	if t.Done {
		status = "[x]"
	}
	fmtProxy := proxy.NewFmt()
	return fmtProxy.Sprintf("%s %d: %s", status, t.ID, t.Description)
}

// MarshalJSON implements custom JSON marshaling
func (t *Todo) MarshalJSON() ([]byte, error) {
	type Alias Todo
	jsonProxy := proxy.NewJSON()
	return jsonProxy.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(t),
		CreatedAt: t.CreatedAt.Format(proxy.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(proxy.RFC3339),
	})
}

// UnmarshalJSON implements custom JSON unmarshalling
func (t *Todo) UnmarshalJSON(data []byte) error {
	type Alias Todo
	aux := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias: (*Alias)(t),
	}

	jsonProxy := proxy.NewJSON()
	if err := jsonProxy.Unmarshal(data, &aux); err != nil {
		return NewDomainError(ErrorTypeJSON, "failed to unmarshal todo", err)
	}

	var err error
	timeProxy := proxy.NewTime()
	t.CreatedAt, err = timeProxy.Parse(proxy.RFC3339, aux.CreatedAt)
	if err != nil {
		return NewDomainError(ErrorTypeJSON, "invalid created_at format", err)
	}

	t.UpdatedAt, err = timeProxy.Parse(proxy.RFC3339, aux.UpdatedAt)
	if err != nil {
		return NewDomainError(ErrorTypeJSON, "invalid updated_at format", err)
	}

	return t.Validate()
}
