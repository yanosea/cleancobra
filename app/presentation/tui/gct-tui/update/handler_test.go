package update

import (
	"errors"
	"testing"
	"time"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name        string
		msg         tea.Msg
		wantCmd     bool
		wantQuit    bool
		setupState  func(*model.StateModel)
		verifyState func(*testing.T, *model.StateModel)
	}{
		{
			name:    "positive testing - window size message updates dimensions",
			msg:     tea.WindowSizeMsg{Width: 100, Height: 50},
			wantCmd: false,
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Width() != 100 || sm.Height() != 50 {
					t.Errorf("Expected dimensions 100x50, got %dx%d", sm.Width(), sm.Height())
				}
			},
		},
		{
			name:     "positive testing - quit key in normal mode",
			msg:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantCmd:  true,
			wantQuit: true,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
			},
		},
		{
			name:    "positive testing - up key moves cursor up",
			msg:     tea.KeyMsg{Type: tea.KeyUp},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos(createTestTodos(3))
				sm.SetCursor(1)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Cursor() != 0 {
					t.Errorf("Expected cursor at 0, got %d", sm.Cursor())
				}
			},
		},
		{
			name:    "positive testing - down key moves cursor down",
			msg:     tea.KeyMsg{Type: tea.KeyDown},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos(createTestTodos(3))
				sm.SetCursor(0)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Cursor() != 1 {
					t.Errorf("Expected cursor at 1, got %d", sm.Cursor())
				}
			},
		},
		{
			name:    "positive testing - a key switches to input mode",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeInput {
					t.Errorf("Expected ModeInput, got %v", sm.Mode())
				}
			},
		},
		{
			name:    "positive testing - e key switches to edit mode when todos exist",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos(createTestTodos(1))
				sm.SetCursor(0)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeEdit {
					t.Errorf("Expected ModeEdit, got %v", sm.Mode())
				}
			},
		},
		{
			name:    "positive testing - d key switches to confirmation mode",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos(createTestTodos(1))
				sm.SetCursor(0)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeConfirmation {
					t.Errorf("Expected ModeConfirmation, got %v", sm.Mode())
				}
			},
		},
		{
			name:    "positive testing - g key moves to top",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos(createTestTodos(3))
				sm.SetCursor(2)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Cursor() != 0 {
					t.Errorf("Expected cursor at 0, got %d", sm.Cursor())
				}
			},
		},
		{
			name:    "positive testing - G key moves to bottom",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos(createTestTodos(3))
				sm.SetCursor(0)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Cursor() != 2 {
					t.Errorf("Expected cursor at 2, got %d", sm.Cursor())
				}
			},
		},
		{
			name:    "positive testing - esc in input mode cancels",
			msg:     tea.KeyMsg{Type: tea.KeyEsc},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeInput)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeNormal {
					t.Errorf("Expected ModeNormal, got %v", sm.Mode())
				}
			},
		},
		{
			name:    "positive testing - y in confirmation mode executes action",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantCmd: true,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeConfirmation)
				sm.SetConfirmation("Test confirmation", func() tea.Cmd {
					return func() tea.Msg { return model.TodoDeletedMsg{ID: 1} }
				})
			},
		},
		{
			name:    "positive testing - n in confirmation mode cancels",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeConfirmation)
				sm.SetConfirmation("Test confirmation", func() tea.Cmd {
					return func() tea.Msg { return model.TodoDeletedMsg{ID: 1} }
				})
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeNormal {
					t.Errorf("Expected ModeNormal, got %v", sm.Mode())
				}
			},
		},
		{
			name:    "positive testing - todos loaded message updates state",
			msg:     model.TodosLoadedMsg{Todos: createTestTodos(2)},
			wantCmd: false,
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if len(sm.Todos()) != 2 {
					t.Errorf("Expected 2 todos, got %d", len(sm.Todos()))
				}
				if sm.Cursor() != 0 {
					t.Errorf("Expected cursor at 0, got %d", sm.Cursor())
				}
			},
		},
		{
			name:    "positive testing - todo added message triggers reload",
			msg:     model.TodoAddedMsg{Todo: &domain.Todo{ID: 1, Description: "New todo"}},
			wantCmd: true,
		},
		{
			name:    "positive testing - todo toggled message updates specific todo",
			msg:     model.TodoToggledMsg{Todo: &domain.Todo{ID: 1, Description: "Test", Done: true}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetTodos(createTestTodos(1))
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if !sm.Todos()[0].Todo().Done {
					t.Error("Expected todo to be marked as done")
				}
			},
		},
		{
			name:    "positive testing - error message sets error and returns to normal mode",
			msg:     model.ErrorMsg{Error: errors.New("test error")},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeInput)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeNormal {
					t.Errorf("Expected ModeNormal, got %v", sm.Mode())
				}
				if sm.ErrorMessage() == "" {
					t.Error("Expected error message to be set")
				}
			},
		},
		{
			name:    "negative testing - e key ignored when no todos exist",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos([]*model.ItemModel{})
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeNormal {
					t.Errorf("Expected to stay in ModeNormal, got %v", sm.Mode())
				}
			},
		},
		{
			name:    "negative testing - d key ignored when no todos exist",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
				sm.SetTodos([]*model.ItemModel{})
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeNormal {
					t.Errorf("Expected to stay in ModeNormal, got %v", sm.Mode())
				}
			},
		},
		{
			name:    "negative testing - unknown key in normal mode ignored",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
			wantCmd: false,
			setupState: func(sm *model.StateModel) {
				sm.SetMode(model.ModeNormal)
			},
			verifyState: func(t *testing.T, sm *model.StateModel) {
				if sm.Mode() != model.ModeNormal {
					t.Errorf("Expected to stay in ModeNormal, got %v", sm.Mode())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh state model for each test
			stateModel := createTestStateModelWithRealDeps()
			
			if tt.setupState != nil {
				tt.setupState(stateModel)
			}

			result := UpdateHandler(stateModel, tt.msg)

			if result.Model == nil {
				t.Error("Expected model to be returned")
			}

			if tt.wantCmd && result.Cmd == nil {
				t.Error("Expected command to be generated")
			}

			if !tt.wantCmd && result.Cmd != nil {
				t.Error("Expected no command to be generated")
			}

			if tt.wantQuit && result.Cmd != nil {
				// Check if it's a quit command by executing it
				if msg := result.Cmd(); msg != tea.Quit() {
					t.Error("Expected quit command")
				}
			}

			if tt.verifyState != nil {
				tt.verifyState(t, result.Model)
			}
		})
	}
}

func TestCreateCommands(t *testing.T) {
	tests := []struct {
		name        string
		createCmd   func() tea.Cmd
		expectedMsg interface{}
	}{
		{
			name:        "positive testing - CreateModeTransitionCommand creates ModeTransitionMsg",
			createCmd:   func() tea.Cmd { return CreateModeTransitionCommand(model.ModeInput) },
			expectedMsg: ModeTransitionMsg{Mode: model.ModeInput},
		},
		{
			name:        "positive testing - CreateNavigationCommand creates NavigationMsg",
			createCmd:   func() tea.Cmd { return CreateNavigationCommand(NavigationUp, 1) },
			expectedMsg: NavigationMsg{Action: NavigationUp, Value: 1},
		},
		{
			name:        "positive testing - CreateAsyncOperationCommand creates AsyncOperationMsg",
			createCmd:   func() tea.Cmd { return CreateAsyncOperationCommand(AsyncOperationAdd, "test") },
			expectedMsg: AsyncOperationMsg{Operation: AsyncOperationAdd, Data: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.createCmd()
			if cmd == nil {
				t.Error("Expected command to be created")
				return
			}

			msg := cmd()
			if msg == nil {
				t.Error("Expected command to return a message")
				return
			}

			switch expected := tt.expectedMsg.(type) {
			case ModeTransitionMsg:
				if actual, ok := msg.(ModeTransitionMsg); !ok {
					t.Errorf("Expected ModeTransitionMsg, got %T", msg)
				} else if actual.Mode != expected.Mode {
					t.Errorf("Expected mode %v, got %v", expected.Mode, actual.Mode)
				}
			case NavigationMsg:
				if actual, ok := msg.(NavigationMsg); !ok {
					t.Errorf("Expected NavigationMsg, got %T", msg)
				} else if actual.Action != expected.Action || actual.Value != expected.Value {
					t.Errorf("Expected action %v, value %d, got action %v, value %d",
						expected.Action, expected.Value, actual.Action, actual.Value)
				}
			case AsyncOperationMsg:
				if actual, ok := msg.(AsyncOperationMsg); !ok {
					t.Errorf("Expected AsyncOperationMsg, got %T", msg)
				} else if actual.Operation != expected.Operation || actual.Data != expected.Data {
					t.Errorf("Expected operation %v, data %v, got operation %v, data %v",
						expected.Operation, expected.Data, actual.Operation, actual.Data)
				}
			}
		})
	}
}

// Helper functions for creating test data

func createTestStateModelWithRealDeps() *model.StateModel {
	// Create a simple mock bubbles implementation for testing
	mockBubbles := &simpleMockBubbles{}
	
	return model.NewStateModel(
		&application.AddTodoUseCase{},
		&application.ListTodoUseCase{},
		&application.ToggleTodoUseCase{},
		&application.DeleteTodoUseCase{},
		mockBubbles,
	)
}

// Simple mock implementation for testing
type simpleMockBubbles struct{}

func (m *simpleMockBubbles) NewTextInput() proxy.TextInput {
	return &simpleMockTextInput{value: ""}
}

type simpleMockTextInput struct {
	value       string
	focused     bool
	placeholder string
	prompt      string
	charLimit   int
	width       int
}

func (m *simpleMockTextInput) SetValue(s string)        { m.value = s }
func (m *simpleMockTextInput) Value() string            { return m.value }
func (m *simpleMockTextInput) SetPlaceholder(str string) { m.placeholder = str }
func (m *simpleMockTextInput) Placeholder() string      { return m.placeholder }
func (m *simpleMockTextInput) Focus() tea.Cmd           { m.focused = true; return nil }
func (m *simpleMockTextInput) Blur()                    { m.focused = false }
func (m *simpleMockTextInput) Focused() bool            { return m.focused }
func (m *simpleMockTextInput) SetPrompt(str string)     { m.prompt = str }
func (m *simpleMockTextInput) Prompt() string           { return m.prompt }
func (m *simpleMockTextInput) SetCharLimit(limit int)   { m.charLimit = limit }
func (m *simpleMockTextInput) CharLimit() int           { return m.charLimit }
func (m *simpleMockTextInput) SetWidth(w int)           { m.width = w }
func (m *simpleMockTextInput) Width() int               { return m.width }
func (m *simpleMockTextInput) Update(msg tea.Msg) (proxy.TextInput, tea.Cmd) {
	return m, nil
}
func (m *simpleMockTextInput) View() string { return m.value }

func createTestTodos(count int) []*model.ItemModel {
	todos := make([]*model.ItemModel, count)
	for i := 0; i < count; i++ {
		todo := &domain.Todo{
			ID:          i + 1,
			Description: "Test todo " + string(rune('A'+i)),
			Done:        false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		todos[i] = model.NewItemModel(todo)
	}
	return todos
}