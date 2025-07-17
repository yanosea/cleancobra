package model

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewItemModel(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
		want *ItemModel
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			want: &ItemModel{
				todo:     &domain.Todo{ID: 1, Description: "Test todo", Done: false},
				selected: false,
				editing:  false,
			},
		},
		{
			name: "positive testing with nil todo",
			todo: nil,
			want: &ItemModel{
				todo:     nil,
				selected: false,
				editing:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewItemModel(tt.todo)
			
			if got.selected != tt.want.selected {
				t.Errorf("NewItemModel().selected = %v, want %v", got.selected, tt.want.selected)
			}
			if got.editing != tt.want.editing {
				t.Errorf("NewItemModel().editing = %v, want %v", got.editing, tt.want.editing)
			}
			if tt.todo != nil && got.todo != nil {
				if got.todo.ID != tt.want.todo.ID {
					t.Errorf("NewItemModel().todo.ID = %v, want %v", got.todo.ID, tt.want.todo.ID)
				}
				if got.todo.Description != tt.want.todo.Description {
					t.Errorf("NewItemModel().todo.Description = %v, want %v", got.todo.Description, tt.want.todo.Description)
				}
			}
		})
	}
}

func TestItemModel_Todo(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
		want *domain.Todo
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			want: &domain.Todo{ID: 1, Description: "Test todo", Done: false},
		},
		{
			name: "positive testing with nil todo",
			todo: nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(tt.todo)
			got := m.Todo()
			
			if tt.want == nil && got != nil {
				t.Errorf("ItemModel.Todo() = %v, want nil", got)
			} else if tt.want != nil && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("ItemModel.Todo().ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Description != tt.want.Description {
					t.Errorf("ItemModel.Todo().Description = %v, want %v", got.Description, tt.want.Description)
				}
			}
		})
	}
}

func TestItemModel_SetTodo(t *testing.T) {
	tests := []struct {
		name    string
		initial *domain.Todo
		newTodo *domain.Todo
	}{
		{
			name: "positive testing",
			initial: &domain.Todo{
				ID:          1,
				Description: "Initial todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			newTodo: &domain.Todo{
				ID:          2,
				Description: "Updated todo",
				Done:        true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name:    "positive testing with nil initial",
			initial: nil,
			newTodo: &domain.Todo{
				ID:          1,
				Description: "New todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(tt.initial)
			m.SetTodo(tt.newTodo)
			
			got := m.Todo()
			if got != tt.newTodo {
				t.Errorf("ItemModel.SetTodo() did not update todo correctly")
			}
		})
	}
}

func TestItemModel_IsSelected(t *testing.T) {
	tests := []struct {
		name     string
		selected bool
		want     bool
	}{
		{
			name:     "positive testing selected true",
			selected: true,
			want:     true,
		},
		{
			name:     "positive testing selected false",
			selected: false,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(nil)
			m.SetSelected(tt.selected)
			
			got := m.IsSelected()
			if got != tt.want {
				t.Errorf("ItemModel.IsSelected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemModel_SetSelected(t *testing.T) {
	tests := []struct {
		name     string
		selected bool
	}{
		{
			name:     "positive testing set selected true",
			selected: true,
		},
		{
			name:     "positive testing set selected false",
			selected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(nil)
			m.SetSelected(tt.selected)
			
			if m.IsSelected() != tt.selected {
				t.Errorf("ItemModel.SetSelected(%v) did not set selection correctly", tt.selected)
			}
		})
	}
}

func TestItemModel_IsEditing(t *testing.T) {
	tests := []struct {
		name    string
		editing bool
		want    bool
	}{
		{
			name:    "positive testing editing true",
			editing: true,
			want:    true,
		},
		{
			name:    "positive testing editing false",
			editing: false,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(nil)
			m.SetEditing(tt.editing)
			
			got := m.IsEditing()
			if got != tt.want {
				t.Errorf("ItemModel.IsEditing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemModel_SetEditing(t *testing.T) {
	tests := []struct {
		name    string
		editing bool
	}{
		{
			name:    "positive testing set editing true",
			editing: true,
		},
		{
			name:    "positive testing set editing false",
			editing: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(nil)
			m.SetEditing(tt.editing)
			
			if m.IsEditing() != tt.editing {
				t.Errorf("ItemModel.SetEditing(%v) did not set editing state correctly", tt.editing)
			}
		})
	}
}

func TestItemModel_Toggle(t *testing.T) {
	tests := []struct {
		name        string
		todo        *domain.Todo
		wantDone    bool
		shouldPanic bool
	}{
		{
			name: "positive testing toggle false to true",
			todo: &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantDone:    true,
			shouldPanic: false,
		},
		{
			name: "positive testing toggle true to false",
			todo: &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantDone:    false,
			shouldPanic: false,
		},
		{
			name:        "positive testing with nil todo",
			todo:        nil,
			wantDone:    false,
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(tt.todo)
			
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("ItemModel.Toggle() should have panicked")
					}
				}()
			}
			
			m.Toggle()
			
			if tt.todo != nil && m.Todo().Done != tt.wantDone {
				t.Errorf("ItemModel.Toggle() done = %v, want %v", m.Todo().Done, tt.wantDone)
			}
		})
	}
}

func TestItemModel_UpdateDescription(t *testing.T) {
	tests := []struct {
		name        string
		todo        *domain.Todo
		description string
		wantErr     bool
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Original description",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			description: "Updated description",
			wantErr:     false,
		},
		{
			name: "negative testing (empty description failed)",
			todo: &domain.Todo{
				ID:          1,
				Description: "Original description",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			description: "",
			wantErr:     true,
		},
		{
			name:        "positive testing with nil todo",
			todo:        nil,
			description: "New description",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(tt.todo)
			
			err := m.UpdateDescription(tt.description)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemModel.UpdateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if !tt.wantErr && tt.todo != nil && m.Todo().Description != tt.description {
				t.Errorf("ItemModel.UpdateDescription() description = %v, want %v", m.Todo().Description, tt.description)
			}
		})
	}
}

func TestItemModel_Init(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(tt.todo)
			
			cmd := m.Init()
			if cmd != nil {
				t.Errorf("ItemModel.Init() should return nil command")
			}
		})
	}
}

func TestItemModel_Update(t *testing.T) {
	todo := &domain.Todo{
		ID:          1,
		Description: "Test todo",
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		msg     tea.Msg
		wantCmd bool
	}{
		{
			name:    "positive testing ItemToggleMsg",
			msg:     ItemToggleMsg{ID: 1},
			wantCmd: false,
		},
		{
			name:    "positive testing ItemEditMsg",
			msg:     ItemEditMsg{ID: 1, Editing: true},
			wantCmd: false,
		},
		{
			name:    "positive testing ItemSelectMsg",
			msg:     ItemSelectMsg{ID: 1, Selected: true},
			wantCmd: false,
		},
		{
			name:    "positive testing unknown message",
			msg:     "unknown",
			wantCmd: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(todo)
			
			model, cmd := m.Update(tt.msg)
			
			if model == nil {
				t.Errorf("ItemModel.Update() returned nil model")
			}
			
			if (cmd != nil) != tt.wantCmd {
				t.Errorf("ItemModel.Update() cmd = %v, wantCmd %v", cmd != nil, tt.wantCmd)
			}
		})
	}
}

func TestItemModel_View(t *testing.T) {
	tests := []struct {
		name     string
		todo     *domain.Todo
		selected bool
		want     string
	}{
		{
			name: "positive testing incomplete todo not selected",
			todo: &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			selected: false,
			want:     "  [ ] Test todo",
		},
		{
			name: "positive testing complete todo selected",
			todo: &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			selected: true,
			want:     "> [x] Test todo",
		},
		{
			name:     "positive testing nil todo",
			todo:     nil,
			selected: false,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(tt.todo)
			m.SetSelected(tt.selected)
			
			got := m.View()
			if got != tt.want {
				t.Errorf("ItemModel.View() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestItemModel_handleKeyMsg(t *testing.T) {
	todo := &domain.Todo{
		ID:          1,
		Description: "Test todo",
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name     string
		selected bool
		key      string
		wantCmd  bool
	}{
		{
			name:     "positive testing space key when selected",
			selected: true,
			key:      " ",
			wantCmd:  true,
		},
		{
			name:     "positive testing e key when selected",
			selected: true,
			key:      "e",
			wantCmd:  true,
		},
		{
			name:     "positive testing key when not selected",
			selected: false,
			key:      " ",
			wantCmd:  false,
		},
		{
			name:     "positive testing unknown key when selected",
			selected: true,
			key:      "x",
			wantCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewItemModel(todo)
			m.SetSelected(tt.selected)
			
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			model, cmd := m.handleKeyMsg(keyMsg)
			
			if model == nil {
				t.Errorf("ItemModel.handleKeyMsg() returned nil model")
			}
			
			if (cmd != nil) != tt.wantCmd {
				t.Errorf("ItemModel.handleKeyMsg() cmd = %v, wantCmd %v", cmd != nil, tt.wantCmd)
			}
		})
	}
}