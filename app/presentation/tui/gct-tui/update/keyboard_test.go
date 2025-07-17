package update

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/mock/gomock"
)

// createKeyMsg creates a proper tea.KeyMsg for testing
func createKeyMsg(keyString string) tea.KeyMsg {
	switch keyString {
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	default:
		// For single character keys
		if len(keyString) == 1 {
			return tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune(keyString),
			}
		}
		// For other special keys, just use runes
		return tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(keyString),
		}
	}
}

func createTestStateModel(t *testing.T) (*model.StateModel, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	
	// Create mock use cases
	mockAddUseCase := &application.AddTodoUseCase{}
	mockListUseCase := &application.ListTodoUseCase{}
	mockToggleUseCase := &application.ToggleTodoUseCase{}
	mockDeleteUseCase := &application.DeleteTodoUseCase{}
	
	// Create mock bubbles
	mockBubbles := proxy.NewMockBubbles(ctrl)
	mockTextInput := proxy.NewMockTextInput(ctrl)
	
	mockBubbles.EXPECT().NewTextInput().Return(mockTextInput)
	mockTextInput.EXPECT().SetPlaceholder("Enter todo description...")
	mockTextInput.EXPECT().SetCharLimit(500)
	mockTextInput.EXPECT().SetWidth(50)
	
	// Allow common input operations that might be called during mode changes
	mockTextInput.EXPECT().Blur().AnyTimes()
	mockTextInput.EXPECT().SetValue(gomock.Any()).AnyTimes()
	mockTextInput.EXPECT().Value().Return("").AnyTimes()
	mockTextInput.EXPECT().Focus().Return(nil).AnyTimes()
	
	stateModel := model.NewStateModel(
		mockAddUseCase,
		mockListUseCase,
		mockToggleUseCase,
		mockDeleteUseCase,
		mockBubbles,
	)
	
	return stateModel, ctrl
}

func TestKeyboardHandler(t *testing.T) {
	tests := []struct {
		name         string
		mode         model.Mode
		keyString    string
		expectQuit   bool
		expectMode   model.Mode
		expectCmd    bool
	}{
		{
			name:       "positive testing (normal mode quit)",
			mode:       model.ModeNormal,
			keyString:  "q",
			expectQuit: true,
			expectMode: model.ModeNormal,
			expectCmd:  true,
		},
		{
			name:       "positive testing (normal mode ctrl+c)",
			mode:       model.ModeNormal,
			keyString:  "ctrl+c",
			expectQuit: true,
			expectMode: model.ModeNormal,
			expectCmd:  true,
		},
		{
			name:       "positive testing (normal mode up)",
			mode:       model.ModeNormal,
			keyString:  "up",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (normal mode k)",
			mode:       model.ModeNormal,
			keyString:  "k",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (normal mode down)",
			mode:       model.ModeNormal,
			keyString:  "down",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (normal mode j)",
			mode:       model.ModeNormal,
			keyString:  "j",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (normal mode add)",
			mode:       model.ModeNormal,
			keyString:  "a",
			expectQuit: false,
			expectMode: model.ModeInput,
			expectCmd:  false,
		},
		{
			name:       "positive testing (input mode enter)",
			mode:       model.ModeInput,
			keyString:  "enter",
			expectQuit: false,
			expectMode: model.ModeInput, // Stays in input mode when no value
			expectCmd:  false,
		},
		{
			name:       "positive testing (input mode escape)",
			mode:       model.ModeInput,
			keyString:  "esc",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (edit mode enter)",
			mode:       model.ModeEdit,
			keyString:  "enter",
			expectQuit: false,
			expectMode: model.ModeEdit,
			expectCmd:  false,
		},
		{
			name:       "positive testing (edit mode escape)",
			mode:       model.ModeEdit,
			keyString:  "esc",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (confirmation mode y)",
			mode:       model.ModeConfirmation,
			keyString:  "y",
			expectQuit: false,
			expectMode: model.ModeConfirmation,
			expectCmd:  false,
		},
		{
			name:       "positive testing (confirmation mode n)",
			mode:       model.ModeConfirmation,
			keyString:  "n",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModel(t)
			defer ctrl.Finish()
			
			stateModel.SetMode(tt.mode)
			
			keyMsg := createKeyMsg(tt.keyString)
			
			result := KeyboardHandler(stateModel, keyMsg)
			
			if result.Model == nil {
				t.Error("KeyboardHandler() returned nil model")
			}
			
			if tt.expectQuit && result.Cmd == nil {
				t.Error("KeyboardHandler() should return quit command")
			}
			
			if !tt.expectQuit && tt.expectCmd && result.Cmd == nil {
				t.Error("KeyboardHandler() should return a command")
			}
			
			if result.Model.Mode() != tt.expectMode {
				t.Errorf("KeyboardHandler() mode = %v, want %v", result.Model.Mode(), tt.expectMode)
			}
		})
	}
}

func TestHandleNormalModeKeys(t *testing.T) {
	tests := []struct {
		name         string
		keyString    string
		expectQuit   bool
		expectMode   model.Mode
		expectCmd    bool
		setupTodos   bool
	}{
		{
			name:       "positive testing (quit with q)",
			keyString:  "q",
			expectQuit: true,
			expectMode: model.ModeNormal,
			expectCmd:  true,
		},
		{
			name:       "positive testing (quit with ctrl+c)",
			keyString:  "ctrl+c",
			expectQuit: true,
			expectMode: model.ModeNormal,
			expectCmd:  true,
		},
		{
			name:       "positive testing (move up)",
			keyString:  "up",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (move up with k)",
			keyString:  "k",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (move down)",
			keyString:  "down",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (move down with j)",
			keyString:  "j",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (go to top)",
			keyString:  "g",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (go to bottom)",
			keyString:  "G",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (toggle todo)",
			keyString:  " ",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  true,
			setupTodos: true,
		},
		{
			name:       "positive testing (add todo)",
			keyString:  "a",
			expectQuit: false,
			expectMode: model.ModeInput,
			expectCmd:  false,
		},
		{
			name:       "positive testing (edit todo)",
			keyString:  "e",
			expectQuit: false,
			expectMode: model.ModeEdit,
			expectCmd:  false,
			setupTodos: true,
		},
		{
			name:       "positive testing (delete todo)",
			keyString:  "d",
			expectQuit: false,
			expectMode: model.ModeConfirmation,
			expectCmd:  false,
			setupTodos: true,
		},
		{
			name:       "positive testing (refresh)",
			keyString:  "r",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  true,
		},
		{
			name:       "positive testing (clear error)",
			keyString:  "esc",
			expectQuit: false,
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModel(t)
			defer ctrl.Finish()
			
			if tt.setupTodos {
				// Add a test todo
				testTodo := &domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				itemModel := model.NewItemModel(testTodo)
				stateModel.SetTodos([]*model.ItemModel{itemModel})
				stateModel.SetCursor(0)
			}
			
			keyMsg := createKeyMsg(tt.keyString)
			
			result := handleNormalModeKeys(stateModel, keyMsg)
			
			if result.Model == nil {
				t.Error("handleNormalModeKeys() returned nil model")
			}
			
			if tt.expectQuit && result.Cmd == nil {
				t.Error("handleNormalModeKeys() should return quit command")
			}
			
			if !tt.expectQuit && tt.expectCmd && result.Cmd == nil {
				t.Error("handleNormalModeKeys() should return a command")
			}
			
			if result.Model.Mode() != tt.expectMode {
				t.Errorf("handleNormalModeKeys() mode = %v, want %v", result.Model.Mode(), tt.expectMode)
			}
		})
	}
}

func TestHandleInputModeKeys(t *testing.T) {
	tests := []struct {
		name       string
		keyString  string
		inputValue string
		expectMode model.Mode
		expectCmd  bool
	}{
		{
			name:       "positive testing (enter with value)",
			keyString:  "enter",
			inputValue: "", // Mock returns empty string, so no mode change or command
			expectMode: model.ModeInput,
			expectCmd:  false,
		},
		{
			name:       "positive testing (enter without value)",
			keyString:  "enter",
			inputValue: "",
			expectMode: model.ModeInput,
			expectCmd:  false,
		},
		{
			name:       "positive testing (escape)",
			keyString:  "esc",
			inputValue: "Some text",
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModel(t)
			defer ctrl.Finish()
			
			stateModel.SetMode(model.ModeInput)
			
			// Mock the input value
			mockTextInput := stateModel.Input().TextInput().(*proxy.MockTextInput)
			mockTextInput.EXPECT().Value().Return(tt.inputValue).AnyTimes()
			
			keyMsg := createKeyMsg(tt.keyString)
			
			result := handleInputModeKeys(stateModel, keyMsg)
			
			if result.Model == nil {
				t.Error("handleInputModeKeys() returned nil model")
			}
			
			if tt.expectCmd && result.Cmd == nil {
				t.Error("handleInputModeKeys() should return a command")
			}
			
			if result.Model.Mode() != tt.expectMode {
				t.Errorf("handleInputModeKeys() mode = %v, want %v", result.Model.Mode(), tt.expectMode)
			}
		})
	}
}

func TestHandleEditModeKeys(t *testing.T) {
	tests := []struct {
		name       string
		keyString  string
		inputValue string
		expectMode model.Mode
		expectCmd  bool
	}{
		{
			name:       "positive testing (enter with value)",
			keyString:  "enter",
			inputValue: "", // Mock returns empty string, so no command
			expectMode: model.ModeEdit,
			expectCmd:  false,
		},
		{
			name:       "positive testing (enter without value)",
			keyString:  "enter",
			inputValue: "",
			expectMode: model.ModeEdit,
			expectCmd:  false,
		},
		{
			name:       "positive testing (escape)",
			keyString:  "esc",
			inputValue: "Some text",
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModel(t)
			defer ctrl.Finish()
			
			stateModel.SetMode(model.ModeEdit)
			
			// Mock the input value
			mockTextInput := stateModel.Input().TextInput().(*proxy.MockTextInput)
			mockTextInput.EXPECT().Value().Return(tt.inputValue).AnyTimes()
			
			keyMsg := createKeyMsg(tt.keyString)
			
			result := handleEditModeKeys(stateModel, keyMsg)
			
			if result.Model == nil {
				t.Error("handleEditModeKeys() returned nil model")
			}
			
			if tt.expectCmd && result.Cmd == nil {
				t.Error("handleEditModeKeys() should return a command")
			}
			
			if result.Model.Mode() != tt.expectMode {
				t.Errorf("handleEditModeKeys() mode = %v, want %v", result.Model.Mode(), tt.expectMode)
			}
		})
	}
}

func TestHandleConfirmationModeKeys(t *testing.T) {
	tests := []struct {
		name       string
		keyString  string
		expectMode model.Mode
		expectCmd  bool
	}{
		{
			name:       "positive testing (yes with y)",
			keyString:  "y",
			expectMode: model.ModeNormal, // ExecuteConfirmation changes mode to Normal
			expectCmd:  true,
		},
		{
			name:       "positive testing (yes with Y)",
			keyString:  "Y",
			expectMode: model.ModeNormal, // ExecuteConfirmation changes mode to Normal
			expectCmd:  true,
		},
		{
			name:       "positive testing (no with n)",
			keyString:  "n",
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (no with N)",
			keyString:  "N",
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
		{
			name:       "positive testing (cancel with esc)",
			keyString:  "esc",
			expectMode: model.ModeNormal,
			expectCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModel(t)
			defer ctrl.Finish()
			
			stateModel.SetMode(model.ModeConfirmation)
			
			// Set up a confirmation action
			stateModel.SetConfirmation("Test confirmation", func() tea.Cmd {
				return func() tea.Msg {
					return model.TodoDeletedMsg{ID: 1}
				}
			})
			
			keyMsg := createKeyMsg(tt.keyString)
			
			result := handleConfirmationModeKeys(stateModel, keyMsg)
			
			if result.Model == nil {
				t.Error("handleConfirmationModeKeys() returned nil model")
			}
			
			if tt.expectCmd && result.Cmd == nil {
				t.Error("handleConfirmationModeKeys() should return a command")
			}
			
			if result.Model.Mode() != tt.expectMode {
				t.Errorf("handleConfirmationModeKeys() mode = %v, want %v", result.Model.Mode(), tt.expectMode)
			}
		})
	}
}