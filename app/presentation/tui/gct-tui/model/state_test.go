package model

import (
	"errors"
	"testing"
	"time"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/mock/gomock"
)

func TestNewStateModel(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			got := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)

			if got == nil {
				t.Errorf("NewStateModel() returned nil")
			}
			if got.Mode() != ModeNormal {
				t.Errorf("NewStateModel() mode = %v, want %v", got.Mode(), ModeNormal)
			}
			if got.Cursor() != 0 {
				t.Errorf("NewStateModel() cursor = %v, want %v", got.Cursor(), 0)
			}
			if len(got.Todos()) != 0 {
				t.Errorf("NewStateModel() todos length = %v, want %v", len(got.Todos()), 0)
			}
		})
	}
}

func TestStateModel_SetTodos(t *testing.T) {
	tests := []struct {
		name  string
		todos []*ItemModel
		want  int
	}{
		{
			name: "positive testing with todos",
			todos: []*ItemModel{
				NewItemModel(&domain.Todo{ID: 1, Description: "Test 1"}),
				NewItemModel(&domain.Todo{ID: 2, Description: "Test 2"}),
			},
			want: 2,
		},
		{
			name:  "positive testing with empty todos",
			todos: []*ItemModel{},
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			m.SetTodos(tt.todos)

			if len(m.Todos()) != tt.want {
				t.Errorf("StateModel.SetTodos() todos length = %v, want %v", len(m.Todos()), tt.want)
			}
		})
	}
}

func TestStateModel_LoadTodos(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*domain.MockTodoRepository)
		wantErr   bool
	}{
		{
			name: "positive testing",
			setupMock: func(m *domain.MockTodoRepository) {
				todos := []domain.Todo{
					{ID: 1, Description: "Test 1", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: 2, Description: "Test 2", Done: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				}
				m.EXPECT().FindAll().Return(todos, nil)
			},
			wantErr: false,
		},
		{
			name: "negative testing (repository error failed)",
			setupMock: func(m *domain.MockTodoRepository) {
				m.EXPECT().FindAll().Return(nil, errors.New("repository error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			cmd := m.LoadTodos()

			if cmd == nil {
				t.Errorf("StateModel.LoadTodos() returned nil command")
				return
			}

			msg := cmd()
			switch msg := msg.(type) {
			case TodosLoadedMsg:
				if tt.wantErr {
					t.Errorf("StateModel.LoadTodos() expected error but got success")
				}
			case ErrorMsg:
				if !tt.wantErr {
					t.Errorf("StateModel.LoadTodos() unexpected error: %v", msg.Error)
				}
			default:
				t.Errorf("StateModel.LoadTodos() unexpected message type: %T", msg)
			}
		})
	}
}

func TestStateModel_SetCursor(t *testing.T) {
	tests := []struct {
		name      string
		todos     []*ItemModel
		cursor    int
		wantCursor int
	}{
		{
			name: "positive testing valid cursor",
			todos: []*ItemModel{
				NewItemModel(&domain.Todo{ID: 1, Description: "Test 1"}),
				NewItemModel(&domain.Todo{ID: 2, Description: "Test 2"}),
			},
			cursor:    1,
			wantCursor: 1,
		},
		{
			name: "positive testing cursor too high",
			todos: []*ItemModel{
				NewItemModel(&domain.Todo{ID: 1, Description: "Test 1"}),
			},
			cursor:    5,
			wantCursor: 0,
		},
		{
			name: "positive testing negative cursor",
			todos: []*ItemModel{
				NewItemModel(&domain.Todo{ID: 1, Description: "Test 1"}),
			},
			cursor:    -1,
			wantCursor: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			m.SetTodos(tt.todos)
			m.SetCursor(tt.cursor)

			if m.Cursor() != tt.wantCursor {
				t.Errorf("StateModel.SetCursor() cursor = %v, want %v", m.Cursor(), tt.wantCursor)
			}
		})
	}
}

func TestStateModel_MoveCursor(t *testing.T) {
	tests := []struct {
		name        string
		todos       []*ItemModel
		initialCursor int
		operation   string
		wantCursor  int
	}{
		{
			name: "positive testing move cursor up",
			todos: []*ItemModel{
				NewItemModel(&domain.Todo{ID: 1, Description: "Test 1"}),
				NewItemModel(&domain.Todo{ID: 2, Description: "Test 2"}),
			},
			initialCursor: 1,
			operation:    "up",
			wantCursor:   0,
		},
		{
			name: "positive testing move cursor down",
			todos: []*ItemModel{
				NewItemModel(&domain.Todo{ID: 1, Description: "Test 1"}),
				NewItemModel(&domain.Todo{ID: 2, Description: "Test 2"}),
			},
			initialCursor: 0,
			operation:    "down",
			wantCursor:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			m.SetTodos(tt.todos)
			m.SetCursor(tt.initialCursor)

			switch tt.operation {
			case "up":
				m.MoveCursorUp()
			case "down":
				m.MoveCursorDown()
			}

			if m.Cursor() != tt.wantCursor {
				t.Errorf("StateModel.MoveCursor%s() cursor = %v, want %v", tt.operation, m.Cursor(), tt.wantCursor)
			}
		})
	}
}

func TestStateModel_SetMode(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want Mode
	}{
		{
			name: "positive testing set normal mode",
			mode: ModeNormal,
			want: ModeNormal,
		},
		{
			name: "positive testing set input mode",
			mode: ModeInput,
			want: ModeInput,
		},
		{
			name: "positive testing set edit mode",
			mode: ModeEdit,
			want: ModeEdit,
		},
		{
			name: "positive testing set confirmation mode",
			mode: ModeConfirmation,
			want: ModeConfirmation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			m.SetMode(tt.mode)

			if m.Mode() != tt.want {
				t.Errorf("StateModel.SetMode() mode = %v, want %v", m.Mode(), tt.want)
			}
		})
	}
}

func TestStateModel_AddTodo(t *testing.T) {
	tests := []struct {
		name        string
		description string
		setupMock   func(*domain.MockTodoRepository)
		wantErr     bool
	}{
		{
			name:        "positive testing",
			description: "New todo",
			setupMock: func(m *domain.MockTodoRepository) {
				todo := domain.Todo{ID: 1, Description: "New todo", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
				m.EXPECT().Save(gomock.Any()).Return([]domain.Todo{todo}, nil)
			},
			wantErr: false,
		},
		{
			name:        "negative testing (empty description failed)",
			description: "",
			setupMock:   func(m *domain.MockTodoRepository) {},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			cmd := m.AddTodo(tt.description)

			if cmd == nil {
				t.Errorf("StateModel.AddTodo() returned nil command")
				return
			}

			msg := cmd()
			switch msg := msg.(type) {
			case TodoAddedMsg:
				if tt.wantErr {
					t.Errorf("StateModel.AddTodo() expected error but got success")
				}
			case ErrorMsg:
				if !tt.wantErr {
					t.Errorf("StateModel.AddTodo() unexpected error: %v", msg.Error)
				}
			default:
				t.Errorf("StateModel.AddTodo() unexpected message type: %T", msg)
			}
		})
	}
}

func TestStateModel_SetError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "positive testing with error",
			err:  errors.New("test error"),
			want: "test error",
		},
		{
			name: "positive testing with nil error",
			err:  nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			m.SetError(tt.err)

			if m.ErrorMessage() != tt.want {
				t.Errorf("StateModel.SetError() error message = %v, want %v", m.ErrorMessage(), tt.want)
			}
		})
	}
}

func TestStateModel_SetSize(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "positive testing",
			width:  100,
			height: 50,
		},
		{
			name:   "positive testing small size",
			width:  10,
			height: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			m.SetSize(tt.width, tt.height)

			if m.Width() != tt.width {
				t.Errorf("StateModel.SetSize() width = %v, want %v", m.Width(), tt.width)
			}
			if m.Height() != tt.height {
				t.Errorf("StateModel.SetSize() height = %v, want %v", m.Height(), tt.height)
			}
		})
	}
}

func TestStateModel_Init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			cmd := m.Init()

			if cmd == nil {
				t.Errorf("StateModel.Init() returned nil command")
			}
		})
	}
}

func TestStateModel_Update(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.Msg
	}{
		{
			name: "positive testing window size message",
			msg:  tea.WindowSizeMsg{Width: 100, Height: 50},
		},
		{
			name: "positive testing todos loaded message",
			msg: TodosLoadedMsg{
				Todos: []*ItemModel{
					NewItemModel(&domain.Todo{ID: 1, Description: "Test"}),
				},
			},
		},
		{
			name: "positive testing error message",
			msg:  ErrorMsg{Error: errors.New("test error")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			model, cmd := m.Update(tt.msg)

			if model == nil {
				t.Errorf("StateModel.Update() returned nil model")
			}

			// Commands can be nil for some message types
			_ = cmd
		})
	}
}

func TestStateModel_View(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
	}{
		{
			name: "positive testing normal mode",
			mode: ModeNormal,
		},
		{
			name: "positive testing input mode",
			mode: ModeInput,
		},
		{
			name: "positive testing edit mode",
			mode: ModeEdit,
		},
		{
			name: "positive testing confirmation mode",
			mode: ModeConfirmation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			addUseCase := application.NewAddTodoUseCase(mockRepo)
			listUseCase := application.NewListTodoUseCase(mockRepo)
			toggleUseCase := application.NewToggleTodoUseCase(mockRepo)
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepo)
			bubbles := proxy.NewBubbles()

			m := NewStateModel(addUseCase, listUseCase, toggleUseCase, deleteUseCase, bubbles)
			m.SetMode(tt.mode)

			view := m.View()
			if view == "" {
				t.Errorf("StateModel.View() returned empty string")
			}
		})
	}
}