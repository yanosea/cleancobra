package domain

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		description string
		wantErr     bool
		errType     ErrorType
	}{
		{
			name:        "positive testing",
			id:          1,
			description: "Buy groceries",
			wantErr:     false,
		},
		{
			name:        "positive testing (whitespace trimming)",
			id:          2,
			description: "  Clean house  ",
			wantErr:     false,
		},
		{
			name:        "negative testing (empty description failed)",
			id:          3,
			description: "",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "negative testing (whitespace only description failed)",
			id:          4,
			description: "   ",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "negative testing (description too long failed)",
			id:          5,
			description: strings.Repeat("a", 501),
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := NewTodo(tt.id, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTodo() expected error but got none")
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("NewTodo() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTodo() unexpected error = %v", err)
				return
			}

			if todo.ID != tt.id {
				t.Errorf("NewTodo() ID = %v, want %v", todo.ID, tt.id)
			}

			expectedDesc := strings.TrimSpace(tt.description)
			if todo.Description != expectedDesc {
				t.Errorf("NewTodo() Description = %v, want %v", todo.Description, expectedDesc)
			}

			if todo.Done != false {
				t.Errorf("NewTodo() Done = %v, want %v", todo.Done, false)
			}

			if todo.CreatedAt.IsZero() {
				t.Errorf("NewTodo() CreatedAt should not be zero")
			}

			if todo.UpdatedAt.IsZero() {
				t.Errorf("NewTodo() UpdatedAt should not be zero")
			}

			if !todo.CreatedAt.Equal(todo.UpdatedAt) {
				t.Errorf("NewTodo() CreatedAt and UpdatedAt should be equal for new todo")
			}
		})
	}
}

func TestNewTodoWithDeps(t *testing.T) {
	// Import proxy package for testing
	timeProxy := proxy.NewTime()
	stringsProxy := proxy.NewStrings()
	
	tests := []struct {
		name        string
		id          int
		description string
		wantErr     bool
		errType     ErrorType
	}{
		{
			name:        "positive testing",
			id:          1,
			description: "Buy groceries",
			wantErr:     false,
		},
		{
			name:        "positive testing (whitespace trimming)",
			id:          2,
			description: "  Clean house  ",
			wantErr:     false,
		},
		{
			name:        "negative testing (empty description failed)",
			id:          3,
			description: "",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "negative testing (whitespace only description failed)",
			id:          4,
			description: "   ",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "negative testing (description too long failed)",
			id:          5,
			description: strings.Repeat("a", 501),
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := NewTodoWithDeps(tt.id, tt.description, timeProxy, stringsProxy)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTodoWithDeps() expected error but got none")
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("NewTodoWithDeps() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTodoWithDeps() unexpected error = %v", err)
				return
			}

			if todo.ID != tt.id {
				t.Errorf("NewTodoWithDeps() ID = %v, want %v", todo.ID, tt.id)
			}

			expectedDesc := strings.TrimSpace(tt.description)
			if todo.Description != expectedDesc {
				t.Errorf("NewTodoWithDeps() Description = %v, want %v", todo.Description, expectedDesc)
			}

			if todo.Done != false {
				t.Errorf("NewTodoWithDeps() Done = %v, want %v", todo.Done, false)
			}

			if todo.CreatedAt.IsZero() {
				t.Errorf("NewTodoWithDeps() CreatedAt should not be zero")
			}

			if todo.UpdatedAt.IsZero() {
				t.Errorf("NewTodoWithDeps() UpdatedAt should not be zero")
			}

			if !todo.CreatedAt.Equal(todo.UpdatedAt) {
				t.Errorf("NewTodoWithDeps() CreatedAt and UpdatedAt should be equal for new todo")
			}
		})
	}
}

func TestTodo_Toggle(t *testing.T) {
	todo, err := NewTodo(1, "Test todo")
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	originalUpdatedAt := todo.UpdatedAt
	
	// Sleep briefly to ensure time difference
	time.Sleep(time.Millisecond)

	// Test toggling from false to true
	todo.Toggle()
	if !todo.Done {
		t.Errorf("Toggle() Done = %v, want %v", todo.Done, true)
	}
	if !todo.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Toggle() should update UpdatedAt timestamp")
	}

	updatedAt1 := todo.UpdatedAt
	time.Sleep(time.Millisecond)

	// Test toggling from true to false
	todo.Toggle()
	if todo.Done {
		t.Errorf("Toggle() Done = %v, want %v", todo.Done, false)
	}
	if !todo.UpdatedAt.After(updatedAt1) {
		t.Errorf("Toggle() should update UpdatedAt timestamp on second toggle")
	}
}

func TestTodo_UpdateDescription(t *testing.T) {
	todo, err := NewTodo(1, "Original description")
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	originalUpdatedAt := todo.UpdatedAt
	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		description string
		wantErr     bool
		errType     ErrorType
	}{
		{
			name:        "positive testing",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "positive testing (whitespace trimming)",
			description: "  Trimmed description  ",
			wantErr:     false,
		},
		{
			name:        "negative testing (empty description failed)",
			description: "",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "negative testing (whitespace only description failed)",
			description: "   ",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "negative testing (description too long failed)",
			description: strings.Repeat("a", 501),
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := todo.UpdateDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateDescription() expected error but got none")
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("UpdateDescription() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateDescription() unexpected error = %v", err)
				return
			}

			expectedDesc := strings.TrimSpace(tt.description)
			if todo.Description != expectedDesc {
				t.Errorf("UpdateDescription() Description = %v, want %v", todo.Description, expectedDesc)
			}

			if !todo.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateDescription() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestTodo_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		todo    *Todo
		wantErr bool
		errType ErrorType
	}{
		{
			name: "positive testing",
			todo: &Todo{
				ID:          1,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "negative testing (invalid ID zero failed)",
			todo: &Todo{
				ID:          0,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
		{
			name: "negative testing (invalid ID negative failed)",
			todo: &Todo{
				ID:          -1,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
		{
			name: "negative testing (empty description failed)",
			todo: &Todo{
				ID:          1,
				Description: "",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
		{
			name: "negative testing (description too long failed)",
			todo: &Todo{
				ID:          1,
				Description: strings.Repeat("a", 501),
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
		{
			name: "negative testing (zero created_at failed)",
			todo: &Todo{
				ID:          1,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   time.Time{},
				UpdatedAt:   now,
			},
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
		{
			name: "negative testing (zero updated_at failed)",
			todo: &Todo{
				ID:          1,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   time.Time{},
			},
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
		{
			name: "negative testing (updated_at before created_at failed)",
			todo: &Todo{
				ID:          1,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now.Add(-time.Hour),
			},
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.todo.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("Validate() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}

func TestTodo_String(t *testing.T) {
	tests := []struct {
		name     string
		todo     *Todo
		expected string
	}{
		{
			name: "positive testing (incomplete todo)",
			todo: &Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
			},
			expected: "[ ] 1: Buy groceries",
		},
		{
			name: "positive testing (completed todo)",
			todo: &Todo{
				ID:          2,
				Description: "Clean house",
				Done:        true,
			},
			expected: "[x] 2: Clean house",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.todo.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTodo_MarshalJSON(t *testing.T) {
	now := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	todo := &Todo{
		ID:          1,
		Description: "Test todo",
		Done:        true,
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Hour),
	}

	data, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("MarshalJSON() unexpected error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Verify all fields are present and correct
	if result["id"] != float64(1) {
		t.Errorf("MarshalJSON() id = %v, want %v", result["id"], 1)
	}

	if result["description"] != "Test todo" {
		t.Errorf("MarshalJSON() description = %v, want %v", result["description"], "Test todo")
	}

	if result["done"] != true {
		t.Errorf("MarshalJSON() done = %v, want %v", result["done"], true)
	}

	expectedCreatedAt := now.Format(time.RFC3339)
	if result["created_at"] != expectedCreatedAt {
		t.Errorf("MarshalJSON() created_at = %v, want %v", result["created_at"], expectedCreatedAt)
	}

	expectedUpdatedAt := now.Add(time.Hour).Format(time.RFC3339)
	if result["updated_at"] != expectedUpdatedAt {
		t.Errorf("MarshalJSON() updated_at = %v, want %v", result["updated_at"], expectedUpdatedAt)
	}
}

func TestTodo_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		errType ErrorType
	}{
		{
			name: "positive testing",
			json: `{
				"id": 1,
				"description": "Test todo",
				"done": true,
				"created_at": "2023-12-25T10:30:00Z",
				"updated_at": "2023-12-25T11:30:00Z"
			}`,
			wantErr: false,
		},
		{
			name: "negative testing (invalid JSON format missing comma failed)",
			json: `{
				"id": 1,
				"description": "Test todo"
				"done": true
			}`,
			wantErr: true,
			errType: ErrorTypeConfiguration, // json.Unmarshal returns regular error, not domain error
		},
		{
			name: "negative testing (invalid JSON format triggers domain error failed)",
			json: `{
				"id": "not-a-number",
				"description": "Test todo",
				"done": true,
				"created_at": "2023-12-25T10:30:00Z",
				"updated_at": "2023-12-25T11:30:00Z"
			}`,
			wantErr: true,
			errType: ErrorTypeJSON, // This should trigger the domain error in UnmarshalJSON
		},
		{
			name: "negative testing (JSON unmarshal error in custom method failed)",
			json: `{
				"id": 1,
				"description": "Test todo",
				"done": "not-a-boolean",
				"created_at": "2023-12-25T10:30:00Z",
				"updated_at": "2023-12-25T11:30:00Z"
			}`,
			wantErr: true,
			errType: ErrorTypeJSON, // This should trigger the "failed to unmarshal todo" error
		},
		{
			name: "negative testing (invalid created_at format failed)",
			json: `{
				"id": 1,
				"description": "Test todo",
				"done": true,
				"created_at": "invalid-date",
				"updated_at": "2023-12-25T11:30:00Z"
			}`,
			wantErr: true,
			errType: ErrorTypeJSON,
		},
		{
			name: "negative testing (invalid updated_at format failed)",
			json: `{
				"id": 1,
				"description": "Test todo",
				"done": true,
				"created_at": "2023-12-25T10:30:00Z",
				"updated_at": "invalid-date"
			}`,
			wantErr: true,
			errType: ErrorTypeJSON,
		},
		{
			name: "negative testing (validation failure empty description failed)",
			json: `{
				"id": 1,
				"description": "",
				"done": true,
				"created_at": "2023-12-25T10:30:00Z",
				"updated_at": "2023-12-25T11:30:00Z"
			}`,
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
		{
			name: "negative testing (validation failure invalid ID failed)",
			json: `{
				"id": 0,
				"description": "Test todo",
				"done": true,
				"created_at": "2023-12-25T10:30:00Z",
				"updated_at": "2023-12-25T11:30:00Z"
			}`,
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var todo Todo
			err := json.Unmarshal([]byte(tt.json), &todo)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalJSON() expected error but got none")
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("UnmarshalJSON() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("UnmarshalJSON() unexpected error = %v", err)
				return
			}

			// Verify the unmarshaled todo is valid
			if err := todo.Validate(); err != nil {
				t.Errorf("UnmarshalJSON() resulted in invalid todo: %v", err)
			}

			// Verify specific values for the valid case
			if tt.name == "valid JSON" {
				if todo.ID != 1 {
					t.Errorf("UnmarshalJSON() ID = %v, want %v", todo.ID, 1)
				}
				if todo.Description != "Test todo" {
					t.Errorf("UnmarshalJSON() Description = %v, want %v", todo.Description, "Test todo")
				}
				if todo.Done != true {
					t.Errorf("UnmarshalJSON() Done = %v, want %v", todo.Done, true)
				}
			}
		})
	}
}