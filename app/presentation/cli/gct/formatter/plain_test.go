package formatter

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
)

func TestNewPlainFormatter(t *testing.T) {
	formatter := NewPlainFormatter()

	if formatter == nil {
		t.Error("NewPlainFormatter() returned nil")
	}
}

func TestPlainFormatter_Format(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		todos   []domain.Todo
		wantErr bool
	}{
		{
			name:    "positive testing (empty todos)",
			todos:   []domain.Todo{},
			wantErr: false,
		},
		{
			name: "positive testing (single incomplete todo)",
			todos: []domain.Todo{
				{
					ID:          1,
					Description: "Buy groceries",
					Done:        false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		{
			name: "positive testing (single completed todo)",
			todos: []domain.Todo{
				{
					ID:          2,
					Description: "Clean house",
					Done:        true,
					CreatedAt:   now,
					UpdatedAt:   now.Add(time.Hour),
				},
			},
			wantErr: false,
		},
		{
			name: "positive testing (multiple todos with mixed status)",
			todos: []domain.Todo{
				{
					ID:          1,
					Description: "Buy groceries",
					Done:        false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          2,
					Description: "Clean house",
					Done:        true,
					CreatedAt:   now,
					UpdatedAt:   now.Add(time.Hour),
				},
				{
					ID:          3,
					Description: "Walk the dog",
					Done:        false,
					CreatedAt:   now.Add(time.Hour * 2),
					UpdatedAt:   now.Add(time.Hour * 2),
				},
			},
			wantErr: false,
		},
		{
			name: "positive testing (todo with long description)",
			todos: []domain.Todo{
				{
					ID:          4,
					Description: "This is a very long description that should be displayed in full without truncation in plain text format",
					Done:        false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		{
			name: "positive testing (todo with special characters)",
			todos: []domain.Todo{
				{
					ID:          5,
					Description: "Review \"important\" document & send email",
					Done:        true,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewPlainFormatter()
			got, err := formatter.Format(tt.todos)

			if (err != nil) != tt.wantErr {
				t.Errorf("PlainFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(tt.todos) == 0 {
					// Empty case should show "No todos found"
					if got != "No todos found." {
						t.Errorf("PlainFormatter.Format() = %v, want 'No todos found.'", got)
					}
				} else {
					// Non-empty case should have each todo on a line
					lines := strings.Split(got, "\n")
					if len(lines) != len(tt.todos) {
						t.Errorf("PlainFormatter.Format() line count = %v, want %v", len(lines), len(tt.todos))
					}

					// Verify each todo appears in output with correct format
					for i, todo := range tt.todos {
						expectedStatus := "[ ]"
						if todo.Done {
							expectedStatus = "[x]"
						}
						
						if i < len(lines) {
							line := lines[i]
							// Check for status indicator
							if !strings.Contains(line, expectedStatus) {
								t.Errorf("PlainFormatter.Format() line %d missing status %v: %v", i, expectedStatus, line)
							}
							// Check for description
							if !strings.Contains(line, todo.Description) {
								t.Errorf("PlainFormatter.Format() line %d missing description %v: %v", i, todo.Description, line)
							}
						}
					}
				}
			}
		})
	}
}

func TestPlainFormatter_FormatSingle(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		todo    domain.Todo
		wantErr bool
	}{
		{
			name: "positive testing (incomplete todo)",
			todo: domain.Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "positive testing (completed todo)",
			todo: domain.Todo{
				ID:          2,
				Description: "Clean house",
				Done:        true,
				CreatedAt:   now,
				UpdatedAt:   now.Add(time.Hour),
			},
			wantErr: false,
		},
		{
			name: "positive testing (todo with special characters)",
			todo: domain.Todo{
				ID:          3,
				Description: "Review \"important\" document & send email",
				Done:        false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "positive testing (todo with long description)",
			todo: domain.Todo{
				ID:          10,
				Description: "This is a very long description that should be displayed in full without truncation in plain text format",
				Done:        true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewPlainFormatter()
			got, err := formatter.FormatSingle(tt.todo)

			if (err != nil) != tt.wantErr {
				t.Errorf("PlainFormatter.FormatSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the todo appears in output with correct format
				expectedStatus := "[ ]"
				if tt.todo.Done {
					expectedStatus = "[x]"
				}
				
				// Check for status indicator
				if !strings.Contains(got, expectedStatus) {
					t.Errorf("PlainFormatter.FormatSingle() missing status %v: %v", expectedStatus, got)
				}
				// Check for description
				if !strings.Contains(got, tt.todo.Description) {
					t.Errorf("PlainFormatter.FormatSingle() missing description %v: %v", tt.todo.Description, got)
				}
				// Check for ID (convert to string properly)
				idStr := fmt.Sprintf("%d", tt.todo.ID)
				if !strings.Contains(got, idStr) {
					t.Errorf("PlainFormatter.FormatSingle() missing ID %v: %v", tt.todo.ID, got)
				}
			}
		})
	}
}