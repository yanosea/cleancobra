package domain

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		description string
		want        *Todo
		wantErr     bool
	}{
		{
			name:        "positive testing",
			id:          1,
			description: "Buy groceries",
			want: &Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
			},
			wantErr: false,
		},
		{
			name:        "positive testing with whitespace trimming",
			id:          2,
			description: "  Clean house  ",
			want: &Todo{
				ID:          2,
				Description: "Clean house",
				Done:        false,
			},
			wantErr: false,
		},
		{
			name:        "negative testing (empty description failed)",
			id:          1,
			description: "",
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "negative testing (whitespace only description failed)",
			id:          1,
			description: "   ",
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "negative testing (description too long failed)",
			id:          1,
			description: strings.Repeat("a", 501),
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTodo(tt.id, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ID != tt.want.ID {
					t.Errorf("NewTodo() ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Description != tt.want.Description {
					t.Errorf("NewTodo() Description = %v, want %v", got.Description, tt.want.Description)
				}
				if got.Done != tt.want.Done {
					t.Errorf("NewTodo() Done = %v, want %v", got.Done, tt.want.Done)
				}
				if got.CreatedAt.IsZero() {
					t.Error("NewTodo() CreatedAt should not be zero")
				}
				if got.UpdatedAt.IsZero() {
					t.Error("NewTodo() UpdatedAt should not be zero")
				}
			}
		})
	}
}

func TestTodo_Toggle(t *testing.T) {
	tests := []struct {
		name         string
		initialDone  bool
		expectedDone bool
	}{
		{
			name:         "positive testing (toggle from false to true)",
			initialDone:  false,
			expectedDone: true,
		},
		{
			name:         "positive testing (toggle from true to false)",
			initialDone:  true,
			expectedDone: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := &Todo{
				ID:          1,
				Description: "Test todo",
				Done:        tt.initialDone,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			
			oldUpdatedAt := todo.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure UpdatedAt changes
			
			todo.Toggle()
			
			if todo.Done != tt.expectedDone {
				t.Errorf("Toggle() Done = %v, want %v", todo.Done, tt.expectedDone)
			}
			if !todo.UpdatedAt.After(oldUpdatedAt) {
				t.Error("Toggle() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestTodo_UpdateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{
			name:        "positive testing",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "positive testing with whitespace trimming",
			description: "  Updated description  ",
			wantErr:     false,
		},
		{
			name:        "negative testing (empty description failed)",
			description: "",
			wantErr:     true,
		},
		{
			name:        "negative testing (whitespace only description failed)",
			description: "   ",
			wantErr:     true,
		},
		{
			name:        "negative testing (description too long failed)",
			description: strings.Repeat("a", 501),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := &Todo{
				ID:          1,
				Description: "Original description",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			
			oldUpdatedAt := todo.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure UpdatedAt changes
			
			err := todo.UpdateDescription(tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				expectedDesc := strings.TrimSpace(tt.description)
				if todo.Description != expectedDesc {
					t.Errorf("UpdateDescription() Description = %v, want %v", todo.Description, expectedDesc)
				}
				if !todo.UpdatedAt.After(oldUpdatedAt) {
					t.Error("UpdateDescription() should update UpdatedAt timestamp")
				}
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
	}{
		{
			name: "positive testing",
			todo: &Todo{
				ID:          1,
				Description: "Valid todo",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "negative testing (invalid ID failed)",
			todo: &Todo{
				ID:          0,
				Description: "Valid todo",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: true,
		},
		{
			name: "negative testing (negative ID failed)",
			todo: &Todo{
				ID:          -1,
				Description: "Valid todo",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: true,
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
		},
		{
			name: "negative testing (zero CreatedAt failed)",
			todo: &Todo{
				ID:          1,
				Description: "Valid todo",
				Done:        false,
				CreatedAt:   time.Time{},
				UpdatedAt:   now,
			},
			wantErr: true,
		},
		{
			name: "negative testing (zero UpdatedAt failed)",
			todo: &Todo{
				ID:          1,
				Description: "Valid todo",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   time.Time{},
			},
			wantErr: true,
		},
		{
			name: "negative testing (UpdatedAt before CreatedAt failed)",
			todo: &Todo{
				ID:          1,
				Description: "Valid todo",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now.Add(-time.Hour),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.todo.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTodo_String(t *testing.T) {
	tests := []struct {
		name string
		todo *Todo
		want string
	}{
		{
			name: "positive testing (incomplete todo)",
			todo: &Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
			},
			want: "[ ] 1: Buy groceries",
		},
		{
			name: "positive testing (completed todo)",
			todo: &Todo{
				ID:          2,
				Description: "Clean house",
				Done:        true,
			},
			want: "[x] 2: Clean house",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.todo.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodo_MarshalJSON(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name    string
		todo    *Todo
		want    string
		wantErr bool
	}{
		{
			name: "positive testing",
			todo: &Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			want:    `{"id":1,"description":"Buy groceries","done":false,"created_at":"2025-01-15T10:00:00Z","updated_at":"2025-01-15T10:00:00Z"}`,
			wantErr: false,
		},
		{
			name: "positive testing (completed todo)",
			todo: &Todo{
				ID:          2,
				Description: "Clean house",
				Done:        true,
				CreatedAt:   now,
				UpdatedAt:   now.Add(time.Hour),
			},
			want:    `{"id":2,"description":"Clean house","done":true,"created_at":"2025-01-15T10:00:00Z","updated_at":"2025-01-15T11:00:00Z"}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.todo.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestTodo_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    *Todo
		wantErr bool
	}{
		{
			name: "positive testing",
			data: `{"id":1,"description":"Buy groceries","done":false,"created_at":"2025-01-15T10:00:00Z","updated_at":"2025-01-15T10:00:00Z"}`,
			want: &Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
				CreatedAt:   time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "positive testing (completed todo)",
			data: `{"id":2,"description":"Clean house","done":true,"created_at":"2025-01-15T10:00:00Z","updated_at":"2025-01-15T11:00:00Z"}`,
			want: &Todo{
				ID:          2,
				Description: "Clean house",
				Done:        true,
				CreatedAt:   time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:    "negative testing (invalid JSON failed)",
			data:    `{"id":1,"description":"Buy groceries","done":false,"created_at":"invalid","updated_at":"2025-01-15T10:00:00Z"}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative testing (invalid created_at format failed)",
			data:    `{"id":1,"description":"Buy groceries","done":false,"created_at":"invalid-date","updated_at":"2025-01-15T10:00:00Z"}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative testing (invalid updated_at format failed)",
			data:    `{"id":1,"description":"Buy groceries","done":false,"created_at":"2025-01-15T10:00:00Z","updated_at":"invalid-date"}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative testing (validation failed)",
			data:    `{"id":0,"description":"Buy groceries","done":false,"created_at":"2025-01-15T10:00:00Z","updated_at":"2025-01-15T10:00:00Z"}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var todo Todo
			err := json.Unmarshal([]byte(tt.data), &todo)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if todo.ID != tt.want.ID {
					t.Errorf("UnmarshalJSON() ID = %v, want %v", todo.ID, tt.want.ID)
				}
				if todo.Description != tt.want.Description {
					t.Errorf("UnmarshalJSON() Description = %v, want %v", todo.Description, tt.want.Description)
				}
				if todo.Done != tt.want.Done {
					t.Errorf("UnmarshalJSON() Done = %v, want %v", todo.Done, tt.want.Done)
				}
				if !todo.CreatedAt.Equal(tt.want.CreatedAt) {
					t.Errorf("UnmarshalJSON() CreatedAt = %v, want %v", todo.CreatedAt, tt.want.CreatedAt)
				}
				if !todo.UpdatedAt.Equal(tt.want.UpdatedAt) {
					t.Errorf("UnmarshalJSON() UpdatedAt = %v, want %v", todo.UpdatedAt, tt.want.UpdatedAt)
				}
			}
		})
	}
}