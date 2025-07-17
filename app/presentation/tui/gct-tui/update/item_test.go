package update

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdateItem(t *testing.T) {
	// Create test todo
	testTodo := &domain.Todo{
		ID:          1,
		Description: "Test todo",
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	tests := []struct {
		name        string
		itemModel   *model.ItemModel
		msg         tea.Msg
		wantCmd     bool
		wantCmdType interface{}
	}{
		{
			name:        "positive testing - key message with space toggles item",
			itemModel:   createSelectedItemModel(testTodo),
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}},
			wantCmd:     true,
			wantCmdType: model.ItemToggleMsg{},
		},
		{
			name:        "positive testing - key message with e triggers edit",
			itemModel:   createSelectedItemModel(testTodo),
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},
			wantCmd:     true,
			wantCmdType: model.ItemEditMsg{},
		},
		{
			name:        "positive testing - key message with d triggers delete",
			itemModel:   createSelectedItemModel(testTodo),
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantCmd:     true,
			wantCmdType: ItemDeleteMsg{},
		},
		{
			name:        "positive testing - item toggle message toggles item",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemToggleMsg{ID: 1},
			wantCmd:     true,
			wantCmdType: ItemToggleAsyncMsg{},
		},
		{
			name:        "positive testing - item edit message sets editing state",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemEditMsg{ID: 1, Editing: true},
			wantCmd:     true,
			wantCmdType: ItemEditStateMsg{},
		},
		{
			name:        "positive testing - item select message sets selection",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemSelectMsg{ID: 1, Selected: true},
			wantCmd:     false,
			wantCmdType: nil,
		},
		{
			name:        "positive testing - item update message updates todo",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemUpdateMsg{Todo: testTodo},
			wantCmd:     true,
			wantCmdType: ItemUpdatedAsyncMsg{},
		},
		{
			name:        "negative testing - unselected item ignores key messages",
			itemModel:   createItemModel(testTodo),
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}},
			wantCmd:     false,
			wantCmdType: nil,
		},
		{
			name:        "negative testing - wrong ID in toggle message ignored",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemToggleMsg{ID: 999},
			wantCmd:     false,
			wantCmdType: nil,
		},
		{
			name:        "negative testing - wrong ID in edit message ignored",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemEditMsg{ID: 999, Editing: true},
			wantCmd:     false,
			wantCmdType: nil,
		},
		{
			name:        "negative testing - wrong ID in select message ignored",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemSelectMsg{ID: 999, Selected: true},
			wantCmd:     false,
			wantCmdType: nil,
		},
		{
			name:        "negative testing - wrong ID in update message ignored",
			itemModel:   createItemModel(testTodo),
			msg:         model.ItemUpdateMsg{Todo: &domain.Todo{ID: 999}},
			wantCmd:     false,
			wantCmdType: nil,
		},
		{
			name:        "negative testing - unknown message type ignored",
			itemModel:   createItemModel(testTodo),
			msg:         tea.WindowSizeMsg{},
			wantCmd:     false,
			wantCmdType: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateItem(tt.itemModel, tt.msg)
			
			if result.Model == nil {
				t.Error("Expected model to be returned")
			}
			
			if tt.wantCmd && result.Cmd == nil {
				t.Error("Expected command to be generated")
			}
			
			if !tt.wantCmd && result.Cmd != nil {
				t.Error("Expected no command to be generated")
			}
			
			if tt.wantCmd && result.Cmd != nil {
				// Execute command to get message type
				msg := result.Cmd()
				if msg == nil {
					t.Error("Expected command to return a message")
				} else {
					// Check message type matches expected
					switch tt.wantCmdType.(type) {
					case model.ItemToggleMsg:
						if _, ok := msg.(model.ItemToggleMsg); !ok {
							t.Errorf("Expected ItemToggleMsg, got %T", msg)
						}
					case model.ItemEditMsg:
						if _, ok := msg.(model.ItemEditMsg); !ok {
							t.Errorf("Expected ItemEditMsg, got %T", msg)
						}
					case ItemDeleteMsg:
						if _, ok := msg.(ItemDeleteMsg); !ok {
							t.Errorf("Expected ItemDeleteMsg, got %T", msg)
						}
					case ItemToggleAsyncMsg:
						if _, ok := msg.(ItemToggleAsyncMsg); !ok {
							t.Errorf("Expected ItemToggleAsyncMsg, got %T", msg)
						}
					case ItemEditStateMsg:
						if _, ok := msg.(ItemEditStateMsg); !ok {
							t.Errorf("Expected ItemEditStateMsg, got %T", msg)
						}
					case ItemUpdatedAsyncMsg:
						if _, ok := msg.(ItemUpdatedAsyncMsg); !ok {
							t.Errorf("Expected ItemUpdatedAsyncMsg, got %T", msg)
						}
					}
				}
			}
		})
	}
}

func TestBatchUpdateItems(t *testing.T) {
	// Create test todos
	todo1 := &domain.Todo{ID: 1, Description: "Todo 1", Done: false}
	todo2 := &domain.Todo{ID: 2, Description: "Todo 2", Done: false}
	
	items := []*model.ItemModel{
		createSelectedItemModel(todo1),
		createSelectedItemModel(todo2),
	}
	
	tests := []struct {
		name         string
		items        []*model.ItemModel
		msg          tea.Msg
		wantCommands int
		wantUpdated  bool
	}{
		{
			name:         "positive testing - key message generates commands for selected items",
			items:        items,
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}},
			wantCommands: 2,
			wantUpdated:  true,
		},
		{
			name:         "positive testing - toggle message affects matching item",
			items:        items,
			msg:          model.ItemToggleMsg{ID: 1},
			wantCommands: 1,
			wantUpdated:  true,
		},
		{
			name:         "negative testing - unknown message generates no commands",
			items:        items,
			msg:          tea.WindowSizeMsg{},
			wantCommands: 0,
			wantUpdated:  false,
		},
		{
			name:         "negative testing - empty items list generates no commands",
			items:        []*model.ItemModel{},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}},
			wantCommands: 0,
			wantUpdated:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands, updated := BatchUpdateItems(tt.items, tt.msg)
			
			if len(commands) != tt.wantCommands {
				t.Errorf("Expected %d commands, got %d", tt.wantCommands, len(commands))
			}
			
			if updated != tt.wantUpdated {
				t.Errorf("Expected updated=%v, got %v", tt.wantUpdated, updated)
			}
		})
	}
}

func TestCreateItemCommands(t *testing.T) {
	tests := []struct {
		name        string
		createCmd   func() tea.Cmd
		expectedMsg interface{}
	}{
		{
			name:        "positive testing - CreateItemAddCommand creates ItemAddMsg",
			createCmd:   func() tea.Cmd { return CreateItemAddCommand("Test todo") },
			expectedMsg: ItemAddMsg{Description: "Test todo"},
		},
		{
			name:        "positive testing - CreateItemDeleteCommand creates ItemDeleteMsg",
			createCmd:   func() tea.Cmd { return CreateItemDeleteCommand(1) },
			expectedMsg: ItemDeleteMsg{ID: 1},
		},
		{
			name:        "positive testing - CreateItemToggleCommand creates ItemToggleMsg",
			createCmd:   func() tea.Cmd { return CreateItemToggleCommand(1) },
			expectedMsg: model.ItemToggleMsg{ID: 1},
		},
		{
			name:        "positive testing - CreateItemSelectCommand creates ItemSelectMsg",
			createCmd:   func() tea.Cmd { return CreateItemSelectCommand(1, true) },
			expectedMsg: model.ItemSelectMsg{ID: 1, Selected: true},
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
			case ItemAddMsg:
				if actual, ok := msg.(ItemAddMsg); !ok {
					t.Errorf("Expected ItemAddMsg, got %T", msg)
				} else if actual.Description != expected.Description {
					t.Errorf("Expected description %q, got %q", expected.Description, actual.Description)
				}
			case ItemDeleteMsg:
				if actual, ok := msg.(ItemDeleteMsg); !ok {
					t.Errorf("Expected ItemDeleteMsg, got %T", msg)
				} else if actual.ID != expected.ID {
					t.Errorf("Expected ID %d, got %d", expected.ID, actual.ID)
				}
			case model.ItemToggleMsg:
				if actual, ok := msg.(model.ItemToggleMsg); !ok {
					t.Errorf("Expected ItemToggleMsg, got %T", msg)
				} else if actual.ID != expected.ID {
					t.Errorf("Expected ID %d, got %d", expected.ID, actual.ID)
				}
			case model.ItemSelectMsg:
				if actual, ok := msg.(model.ItemSelectMsg); !ok {
					t.Errorf("Expected ItemSelectMsg, got %T", msg)
				} else if actual.ID != expected.ID || actual.Selected != expected.Selected {
					t.Errorf("Expected ID %d, Selected %v, got ID %d, Selected %v", 
						expected.ID, expected.Selected, actual.ID, actual.Selected)
				}
			}
		})
	}
}

// Helper functions for creating test models

func createItemModel(todo *domain.Todo) *model.ItemModel {
	return model.NewItemModel(todo)
}

func createSelectedItemModel(todo *domain.Todo) *model.ItemModel {
	itemModel := model.NewItemModel(todo)
	itemModel.SetSelected(true)
	return itemModel
}