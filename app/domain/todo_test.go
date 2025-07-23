package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yanosea/gct/pkg/proxy"
	"go.uber.org/mock/gomock"
)

func TestNewTodo(t *testing.T) {
	// Create mock for time to ensure CreatedAt and UpdatedAt are the same
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTime := proxy.NewMockTime(ctrl)
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockTime.EXPECT().Now().Return(fixedTime).AnyTimes()

	// Initialize with real proxies for testing except for time
	InitializeDomain(mockTime, proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	tests := []struct {
		name        string
		id          int
		description string
		wantErr     bool
		errType     ErrorType
	}{
		{
			name:        "Valid todo",
			id:          1,
			description: "Buy milk",
			wantErr:     false,
		},
		{
			name:        "Valid todo with spaces",
			id:          2,
			description: "  Buy bread  ",
			wantErr:     false,
		},
		{
			name:        "Empty description",
			id:          1,
			description: "",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "Whitespace only description",
			id:          1,
			description: "   ",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "Description too long",
			id:          1,
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
					t.Errorf("NewTodo() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !IsDomainError(err) {
					t.Errorf("NewTodo() error is not a domain error: %v", err)
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("NewTodo() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if todo.ID != tt.id {
				t.Errorf("NewTodo().ID = %v, want %v", todo.ID, tt.id)
			}

			expectedDesc := strings.TrimSpace(tt.description)
			if todo.Description != expectedDesc {
				t.Errorf("NewTodo().Description = %v, want %v", todo.Description, expectedDesc)
			}

			if todo.Done != false {
				t.Errorf("NewTodo().Done = %v, want false", todo.Done)
			}

			if todo.CreatedAt.IsZero() {
				t.Errorf("NewTodo().CreatedAt is zero")
			}

			if todo.UpdatedAt.IsZero() {
				t.Errorf("NewTodo().UpdatedAt is zero")
			}

			if !todo.CreatedAt.Equal(todo.UpdatedAt) {
				t.Errorf("NewTodo() CreatedAt and UpdatedAt should be equal for new todo")
			}
		})
	}
}

func TestTodo_Toggle(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	todo, err := NewTodo(1, "Test todo")
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	originalUpdatedAt := todo.UpdatedAt

	// Sleep a small amount to ensure time difference
	time.Sleep(1 * time.Millisecond)

	// Test toggling from false to true
	todo.Toggle()

	if todo.Done != true {
		t.Errorf("After first toggle, Done = %v, want true", todo.Done)
	}

	if !todo.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt should be updated after toggle")
	}

	secondUpdatedAt := todo.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	// Test toggling from true to false
	todo.Toggle()

	if todo.Done != false {
		t.Errorf("After second toggle, Done = %v, want false", todo.Done)
	}

	if !todo.UpdatedAt.After(secondUpdatedAt) {
		t.Errorf("UpdatedAt should be updated after second toggle")
	}
}

func TestTodo_UpdateDescription(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	todo, err := NewTodo(1, "Original description")
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	originalUpdatedAt := todo.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name        string
		description string
		wantErr     bool
		errType     ErrorType
		expected    string
	}{
		{
			name:        "Valid update",
			description: "New description",
			wantErr:     false,
			expected:    "New description",
		},
		{
			name:        "Update with spaces",
			description: "  Trimmed description  ",
			wantErr:     false,
			expected:    "Trimmed description",
		},
		{
			name:        "Empty description",
			description: "",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "Whitespace only",
			description: "   ",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "Description too long",
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
					t.Errorf("UpdateDescription() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !IsDomainError(err) {
					t.Errorf("UpdateDescription() error is not a domain error: %v", err)
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("UpdateDescription() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if todo.Description != tt.expected {
				t.Errorf("After UpdateDescription(), Description = %v, want %v", todo.Description, tt.expected)
			}

			if !todo.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdatedAt should be updated after description change")
			}

			originalUpdatedAt = todo.UpdatedAt
			time.Sleep(1 * time.Millisecond)
		})
	}
}

func TestTodo_Validate(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		name    string
		todo    Todo
		wantErr bool
		errType ErrorType
	}{
		{
			name: "Valid todo",
			todo: Todo{
				ID:          1,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "Invalid ID - zero",
			todo: Todo{
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
			name: "Invalid ID - negative",
			todo: Todo{
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
			name: "Empty description",
			todo: Todo{
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
			name: "Description too long",
			todo: Todo{
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
			name: "Zero CreatedAt",
			todo: Todo{
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
			name: "Zero UpdatedAt",
			todo: Todo{
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
			name: "UpdatedAt before CreatedAt",
			todo: Todo{
				ID:          1,
				Description: "Valid description",
				Done:        false,
				CreatedAt:   future,
				UpdatedAt:   past,
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
					t.Errorf("Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !IsDomainError(err) {
					t.Errorf("Validate() error is not a domain error: %v", err)
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("Validate() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTodo_String(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	tests := []struct {
		name     string
		todo     Todo
		expected string
	}{
		{
			name: "Incomplete todo",
			todo: Todo{
				ID:          1,
				Description: "Buy milk",
				Done:        false,
			},
			expected: "[ ] 1: Buy milk",
		},
		{
			name: "Complete todo",
			todo: Todo{
				ID:          2,
				Description: "Buy bread",
				Done:        true,
			},
			expected: "[x] 2: Buy bread",
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
	// Initialize with real proxies for testing
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)

	todo := Todo{
		ID:          1,
		Description: "Test todo",
		Done:        true,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	data, err := todo.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Parse the JSON to verify structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Verify fields
	if int(result["id"].(float64)) != 1 {
		t.Errorf("MarshalJSON() id = %v, want 1", result["id"])
	}

	if result["description"] != "Test todo" {
		t.Errorf("MarshalJSON() description = %v, want 'Test todo'", result["description"])
	}

	if result["done"] != true {
		t.Errorf("MarshalJSON() done = %v, want true", result["done"])
	}

	if result["created_at"] != "2023-01-01T12:00:00Z" {
		t.Errorf("MarshalJSON() created_at = %v, want '2023-01-01T12:00:00Z'", result["created_at"])
	}

	if result["updated_at"] != "2023-01-01T13:00:00Z" {
		t.Errorf("MarshalJSON() updated_at = %v, want '2023-01-01T13:00:00Z'", result["updated_at"])
	}
}

func TestTodo_MarshalJSON_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJSON := proxy.NewMockJSON(ctrl)
	mockJSON.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("marshal error"))

	// Initialize with mock JSON proxy for error simulation
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), mockJSON)

	todo := Todo{
		ID:          1,
		Description: "Test todo",
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := todo.MarshalJSON()
	if err == nil {
		t.Errorf("MarshalJSON() error = nil, want error")
	}
}

func TestTodo_UnmarshalJSON(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	tests := []struct {
		name    string
		data    string
		wantErr bool
		errType ErrorType
		verify  func(*testing.T, *Todo)
	}{
		{
			name: "Valid JSON",
			data: `{
				"id": 1,
				"description": "Test todo",
				"done": true,
				"created_at": "2023-01-01T12:00:00Z",
				"updated_at": "2023-01-01T13:00:00Z"
			}`,
			wantErr: false,
			verify: func(t *testing.T, todo *Todo) {
				if todo.ID != 1 {
					t.Errorf("UnmarshalJSON() ID = %v, want 1", todo.ID)
				}
				if todo.Description != "Test todo" {
					t.Errorf("UnmarshalJSON() Description = %v, want 'Test todo'", todo.Description)
				}
				if todo.Done != true {
					t.Errorf("UnmarshalJSON() Done = %v, want true", todo.Done)
				}
				expectedCreated := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
				if !todo.CreatedAt.Equal(expectedCreated) {
					t.Errorf("UnmarshalJSON() CreatedAt = %v, want %v", todo.CreatedAt, expectedCreated)
				}
				expectedUpdated := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)
				if !todo.UpdatedAt.Equal(expectedUpdated) {
					t.Errorf("UnmarshalJSON() UpdatedAt = %v, want %v", todo.UpdatedAt, expectedUpdated)
				}
			},
		},
		{
			name:    "Invalid JSON",
			data:    `{"invalid": json}`,
			wantErr: true,
			errType: ErrorTypeJSON,
		},
		{
			name: "Invalid created_at format",
			data: `{
				"id": 1,
				"description": "Test todo",
				"done": false,
				"created_at": "invalid-date",
				"updated_at": "2023-01-01T13:00:00Z"
			}`,
			wantErr: true,
			errType: ErrorTypeJSON,
		},
		{
			name: "Invalid updated_at format",
			data: `{
				"id": 1,
				"description": "Test todo",
				"done": false,
				"created_at": "2023-01-01T12:00:00Z",
				"updated_at": "invalid-date"
			}`,
			wantErr: true,
			errType: ErrorTypeJSON,
		},
		{
			name: "Invalid todo data (fails validation)",
			data: `{
				"id": 0,
				"description": "Test todo",
				"done": false,
				"created_at": "2023-01-01T12:00:00Z",
				"updated_at": "2023-01-01T13:00:00Z"
			}`,
			wantErr: true,
			errType: ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var todo Todo
			err := todo.UnmarshalJSON([]byte(tt.data))

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalJSON() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !IsDomainError(err) {
					t.Errorf("UnmarshalJSON() error is not a domain error: %v", err)
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("UnmarshalJSON() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.verify != nil {
				tt.verify(t, &todo)
			}
		})
	}
}

func TestTodo_UnmarshalJSON_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJSON := proxy.NewMockJSON(ctrl)
	mockJSON.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(errors.New("unmarshal error"))

	// Initialize with mock JSON proxy for error simulation
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), mockJSON)

	var todo Todo
	err := todo.UnmarshalJSON([]byte(`{"id": 1}`))

	if err == nil {
		t.Errorf("UnmarshalJSON() error = nil, want error")
		return
	}

	if !IsDomainError(err) {
		t.Errorf("UnmarshalJSON() error is not a domain error: %v", err)
		return
	}

	if GetErrorType(err) != ErrorTypeJSON {
		t.Errorf("UnmarshalJSON() error type = %v, want %v", GetErrorType(err), ErrorTypeJSON)
	}
}

func TestValidateDescription(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomain(proxy.NewTime(), proxy.NewStrings(), proxy.NewFmt(), proxy.NewJSON())

	tests := []struct {
		name        string
		description string
		wantErr     bool
		errType     ErrorType
	}{
		{
			name:        "Valid description",
			description: "Valid description",
			wantErr:     false,
		},
		{
			name:        "Description with spaces",
			description: "  Valid description  ",
			wantErr:     false,
		},
		{
			name:        "Empty description",
			description: "",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "Whitespace only",
			description: "   ",
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
		{
			name:        "Description at limit (500 chars)",
			description: strings.Repeat("a", 500),
			wantErr:     false,
		},
		{
			name:        "Description too long (501 chars)",
			description: strings.Repeat("a", 501),
			wantErr:     true,
			errType:     ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateDescription() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !IsDomainError(err) {
					t.Errorf("validateDescription() error is not a domain error: %v", err)
					return
				}
				if GetErrorType(err) != tt.errType {
					t.Errorf("validateDescription() error type = %v, want %v", GetErrorType(err), tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("validateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitializeDomain(t *testing.T) {
	timeProxy := proxy.NewTime()
	stringsProxy := proxy.NewStrings()
	fmtProxy := proxy.NewFmt()
	jsonProxy := proxy.NewJSON()

	InitializeDomain(timeProxy, stringsProxy, fmtProxy, jsonProxy)

	// Test that the proxies are set by creating a todo
	todo, err := NewTodo(1, "Test todo")
	if err != nil {
		t.Errorf("After InitializeDomain, NewTodo() failed: %v", err)
	}

	if todo == nil {
		t.Errorf("After InitializeDomain, NewTodo() returned nil")
	}

	// Test String method works (uses fmt proxy)
	result := todo.String()
	expected := "[ ] 1: Test todo"
	if result != expected {
		t.Errorf("After InitializeDomain, String() = %v, want %v", result, expected)
	}
}
