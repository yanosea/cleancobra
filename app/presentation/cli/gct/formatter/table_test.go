package formatter

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewTableFormatter(t *testing.T) {
	colorProxy := proxy.NewColor()
	stringsProxy := proxy.NewStrings()
	fmtProxy := proxy.NewFmt()
	formatter := NewTableFormatter(colorProxy, stringsProxy, fmtProxy)

	if formatter == nil {
		t.Error("NewTableFormatter() returned nil")
	}
	if formatter.colorProxy != colorProxy {
		t.Error("NewTableFormatter() colorProxy not set correctly")
	}
	if formatter.stringsProxy != stringsProxy {
		t.Error("NewTableFormatter() stringsProxy not set correctly")
	}
	if formatter.fmtProxy != fmtProxy {
		t.Error("NewTableFormatter() fmtProxy not set correctly")
	}
}

func TestTableFormatter_Format(t *testing.T) {
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
					Description: "Write documentation",
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
					ID:          4,
					Description: "Review \"important\" document & send email",
					Done:        false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		{
			name: "positive testing (todo with long description)",
			todos: []domain.Todo{
				{
					ID:          5,
					Description: "This is a very long todo description that should be handled properly by the table formatter without breaking the layout",
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
			colorProxy := proxy.NewColor()
			stringsProxy := proxy.NewStrings()
			fmtProxy := proxy.NewFmt()
			formatter := NewTableFormatter(colorProxy, stringsProxy, fmtProxy)
			got, err := formatter.Format(tt.todos)

			if (err != nil) != tt.wantErr {
				t.Errorf("TableFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify output is not empty
				if got == "" {
					t.Error("TableFormatter.Format() returned empty string")
					return
				}

				// For empty todos, should show "No todos found."
				if len(tt.todos) == 0 {
					if !stringsProxy.Contains(got, "No todos found.") {
						t.Errorf("TableFormatter.Format() for empty todos should contain 'No todos found.', got: %s", got)
					}
					return
				}

				// For non-empty todos, should contain header
				if !stringsProxy.Contains(got, "ID") || !stringsProxy.Contains(got, "Status") || !stringsProxy.Contains(got, "Description") {
					t.Errorf("TableFormatter.Format() should contain table headers, got: %s", got)
				}

				// Should contain separator line
				if !stringsProxy.Contains(got, "----") {
					t.Errorf("TableFormatter.Format() should contain separator line, got: %s", got)
				}

				// Verify each todo appears in output
				for _, todo := range tt.todos {
					if !stringsProxy.Contains(got, todo.Description) {
						t.Errorf("TableFormatter.Format() should contain todo description '%s', got: %s", todo.Description, got)
					}
				}
			}
		})
	}
}

func TestTableFormatter_FormatSingle(t *testing.T) {
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
				ID:          4,
				Description: "This is a very long todo description that should be handled properly by the table formatter without breaking the layout",
				Done:        true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorProxy := proxy.NewColor()
			stringsProxy := proxy.NewStrings()
			fmtProxy := proxy.NewFmt()
			formatter := NewTableFormatter(colorProxy, stringsProxy, fmtProxy)
			got, err := formatter.FormatSingle(tt.todo)

			if (err != nil) != tt.wantErr {
				t.Errorf("TableFormatter.FormatSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify output is not empty
				if got == "" {
					t.Error("TableFormatter.FormatSingle() returned empty string")
					return
				}

				// Should contain header
				if !stringsProxy.Contains(got, "ID") || !stringsProxy.Contains(got, "Status") || !stringsProxy.Contains(got, "Description") {
					t.Errorf("TableFormatter.FormatSingle() should contain table headers, got: %s", got)
				}

				// Should contain separator line
				if !stringsProxy.Contains(got, "----") {
					t.Errorf("TableFormatter.FormatSingle() should contain separator line, got: %s", got)
				}

				// Should contain the todo description
				if !stringsProxy.Contains(got, tt.todo.Description) {
					t.Errorf("TableFormatter.FormatSingle() should contain todo description '%s', got: %s", tt.todo.Description, got)
				}
			}
		})
	}
}

func TestTableFormatter_formatTodoLine(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		todo domain.Todo
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorProxy := proxy.NewColor()
			stringsProxy := proxy.NewStrings()
			fmtProxy := proxy.NewFmt()
			formatter := NewTableFormatter(colorProxy, stringsProxy, fmtProxy)
			got := formatter.formatTodoLine(tt.todo)

			// Verify output is not empty
			if got == "" {
				t.Error("TableFormatter.formatTodoLine() returned empty string")
				return
			}

			// Should contain the todo description
			if !stringsProxy.Contains(got, tt.todo.Description) {
				t.Errorf("TableFormatter.formatTodoLine() should contain todo description '%s', got: %s", tt.todo.Description, got)
			}

			// Should contain appropriate status indicator
			if tt.todo.Done {
				if !stringsProxy.Contains(got, "Done") {
					t.Errorf("TableFormatter.formatTodoLine() for completed todo should contain 'Done', got: %s", got)
				}
			} else {
				if !stringsProxy.Contains(got, "Todo") {
					t.Errorf("TableFormatter.formatTodoLine() for incomplete todo should contain 'Todo', got: %s", got)
				}
			}
		})
	}
}