package domain

import (
	"github.com/yanosea/gct/pkg/proxy"
)

// Domain proxies for dependency injection
var (
	// dtp is a proxy for the time package for dependency injection
	dtp proxy.Time
	// dsp is a proxy for the strings package for dependency injection
	dsp proxy.Strings
	// dfp is a proxy for the fmt package for dependency injection
	dfp proxy.Fmt
	// djp is a proxy for the json package for dependency injection
	djp proxy.JSON
)

// Todo represents a todo item in the domain
type Todo struct {
	ID          int             `json:"id"`
	Description string          `json:"description"`
	Done        bool            `json:"done"`
	CreatedAt   proxy.TimeAlias `json:"created_at"`
	UpdatedAt   proxy.TimeAlias `json:"updated_at"`
}

// InitializeDomain sets the proxy for the domain
func InitializeDomain(timeProxy proxy.Time, stringsProxy proxy.Strings, fmtProxy proxy.Fmt, jsonProxy proxy.JSON) {
	dtp = timeProxy
	dsp = stringsProxy
	dfp = fmtProxy
	djp = jsonProxy
}

// NewTodo creates a new Todo with the given description
func NewTodo(id int, description string) (*Todo, error) {
	if err := validateDescription(description); err != nil {
		return nil, err
	}

	trimmedDesc := dsp.TrimSpace(description)
	return &Todo{
		ID:          id,
		Description: trimmedDesc,
		Done:        false,
		CreatedAt:   dtp.Now(),
		UpdatedAt:   dtp.Now(),
	}, nil
}

// Toggle toggles the completion status of the todo
func (t *Todo) Toggle() {
	t.Done = !t.Done
	t.UpdatedAt = dtp.Now()
}

// UpdateDescription updates the description of the todo
func (t *Todo) UpdateDescription(description string) error {
	if err := validateDescription(description); err != nil {
		return err
	}

	t.Description = dsp.TrimSpace(description)
	t.UpdatedAt = dtp.Now()
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
	trimmed := dsp.TrimSpace(description)
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
	return dfp.Sprintf("%s %d: %s", status, t.ID, t.Description)
}

// MarshalJSON implements custom JSON marshaling
func (t *Todo) MarshalJSON() ([]byte, error) {
	type Alias Todo
	return djp.Marshal(&struct {
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

	if err := djp.Unmarshal(data, &aux); err != nil {
		return NewDomainError(ErrorTypeJSON, "failed to unmarshal todo", err)
	}

	var err error
	t.CreatedAt, err = dtp.Parse(proxy.RFC3339, aux.CreatedAt)
	if err != nil {
		return NewDomainError(ErrorTypeJSON, "invalid created_at format", err)
	}

	t.UpdatedAt, err = dtp.Parse(proxy.RFC3339, aux.UpdatedAt)
	if err != nil {
		return NewDomainError(ErrorTypeJSON, "invalid updated_at format", err)
	}

	return t.Validate()
}
