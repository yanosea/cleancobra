package formatter

import (
	"errors"
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
	"go.uber.org/mock/gomock"
)

func TestNewJSONFormatter(t *testing.T) {
	jsonProxy := proxy.NewJSON()
	formatter := NewJSONFormatter(jsonProxy)

	if formatter == nil {
		t.Error("NewJSONFormatter() returned nil")
	}
	if formatter.jsonProxy != jsonProxy {
		t.Error("NewJSONFormatter() jsonProxy not set correctly")
	}
}

func TestJSONFormatter_Format(t *testing.T) {
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
			name: "positive testing (single todo)",
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
			name: "positive testing (multiple todos)",
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
			},
			wantErr: false,
		},
		{
			name: "positive testing (todo with special characters)",
			todos: []domain.Todo{
				{
					ID:          3,
					Description: "Review \"important\" document & send email",
					Done:        false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonProxy := proxy.NewJSON()
			formatter := NewJSONFormatter(jsonProxy)
			got, err := formatter.Format(tt.todos)

			if (err != nil) != tt.wantErr {
				t.Errorf("JSONFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify it's valid JSON by unmarshaling back
				var result []domain.Todo
				if err := jsonProxy.Unmarshal([]byte(got), &result); err != nil {
					t.Errorf("JSONFormatter.Format() produced invalid JSON: %v", err)
					return
				}

				// Verify the content matches
				if len(result) != len(tt.todos) {
					t.Errorf("JSONFormatter.Format() result count = %v, want %v", len(result), len(tt.todos))
					return
				}

				for i, todo := range result {
					if i < len(tt.todos) {
						if todo.ID != tt.todos[i].ID {
							t.Errorf("JSONFormatter.Format() todo[%d].ID = %v, want %v", i, todo.ID, tt.todos[i].ID)
						}
						if todo.Description != tt.todos[i].Description {
							t.Errorf("JSONFormatter.Format() todo[%d].Description = %v, want %v", i, todo.Description, tt.todos[i].Description)
						}
						if todo.Done != tt.todos[i].Done {
							t.Errorf("JSONFormatter.Format() todo[%d].Done = %v, want %v", i, todo.Done, tt.todos[i].Done)
						}
					}
				}
			}
		})
	}

	// Test marshal error case using mock (異常系でエラーを意図的に起こす)
	t.Run("negative testing (marshal error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJSON := proxy.NewMockJSON(ctrl)
		mockJSON.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("marshal error"))

		formatter := NewJSONFormatter(mockJSON)
		todos := []domain.Todo{{ID: 1, Description: "Test", Done: false}}
		
		got, err := formatter.Format(todos)
		
		if err == nil {
			t.Error("JSONFormatter.Format() expected error but got none")
		}
		if got != "" {
			t.Errorf("JSONFormatter.Format() = %v, want empty string on error", got)
		}
		if !domain.IsDomainError(err) {
			t.Errorf("JSONFormatter.Format() error should be a domain error, got %T", err)
		}
	})
}

func TestJSONFormatter_FormatSingle(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonProxy := proxy.NewJSON()
			formatter := NewJSONFormatter(jsonProxy)
			got, err := formatter.FormatSingle(tt.todo)

			if (err != nil) != tt.wantErr {
				t.Errorf("JSONFormatter.FormatSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify it's valid JSON by unmarshaling back
				var result domain.Todo
				if err := jsonProxy.Unmarshal([]byte(got), &result); err != nil {
					t.Errorf("JSONFormatter.FormatSingle() produced invalid JSON: %v", err)
					return
				}

				// Verify the content matches
				if result.ID != tt.todo.ID {
					t.Errorf("JSONFormatter.FormatSingle() ID = %v, want %v", result.ID, tt.todo.ID)
				}
				if result.Description != tt.todo.Description {
					t.Errorf("JSONFormatter.FormatSingle() Description = %v, want %v", result.Description, tt.todo.Description)
				}
				if result.Done != tt.todo.Done {
					t.Errorf("JSONFormatter.FormatSingle() Done = %v, want %v", result.Done, tt.todo.Done)
				}
			}
		})
	}

	// Test marshal error case using mock (異常系でエラーを意図的に起こす)
	t.Run("negative testing (marshal error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJSON := proxy.NewMockJSON(ctrl)
		mockJSON.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("marshal error"))

		formatter := NewJSONFormatter(mockJSON)
		todo := domain.Todo{ID: 1, Description: "Test", Done: false}
		
		got, err := formatter.FormatSingle(todo)
		
		if err == nil {
			t.Error("JSONFormatter.FormatSingle() expected error but got none")
		}
		if got != "" {
			t.Errorf("JSONFormatter.FormatSingle() = %v, want empty string on error", got)
		}
		if !domain.IsDomainError(err) {
			t.Errorf("JSONFormatter.FormatSingle() error should be a domain error, got %T", err)
		}
	})
}