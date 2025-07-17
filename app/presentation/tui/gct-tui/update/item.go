package update

import (
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	tea "github.com/charmbracelet/bubbletea"
)

// ItemUpdateResult represents the result of an item update operation
type ItemUpdateResult struct {
	Model *model.ItemModel
	Cmd   tea.Cmd
}

// UpdateItem handles item-specific update operations
func UpdateItem(itemModel *model.ItemModel, msg tea.Msg) ItemUpdateResult {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleItemKeyMsg(itemModel, msg)
	case model.ItemToggleMsg:
		return handleItemToggle(itemModel, msg)
	case model.ItemEditMsg:
		return handleItemEdit(itemModel, msg)
	case model.ItemSelectMsg:
		return handleItemSelect(itemModel, msg)
	case model.ItemUpdateMsg:
		return handleItemUpdate(itemModel, msg)
	}
	
	return ItemUpdateResult{Model: itemModel, Cmd: nil}
}

// handleItemKeyMsg handles keyboard input for individual items
func handleItemKeyMsg(itemModel *model.ItemModel, msg tea.KeyMsg) ItemUpdateResult {
	if !itemModel.IsSelected() {
		return ItemUpdateResult{Model: itemModel, Cmd: nil}
	}
	
	switch msg.String() {
	case " ":
		// Generate toggle command
		return ItemUpdateResult{
			Model: itemModel,
			Cmd: func() tea.Msg {
				return model.ItemToggleMsg{ID: itemModel.Todo().ID}
			},
		}
	case "e":
		// Generate edit command
		return ItemUpdateResult{
			Model: itemModel,
			Cmd: func() tea.Msg {
				return model.ItemEditMsg{ID: itemModel.Todo().ID, Editing: true}
			},
		}
	case "d":
		// Generate delete command
		return ItemUpdateResult{
			Model: itemModel,
			Cmd: func() tea.Msg {
				return ItemDeleteMsg{ID: itemModel.Todo().ID}
			},
		}
	}
	
	return ItemUpdateResult{Model: itemModel, Cmd: nil}
}

// handleItemToggle handles toggle operations for items
func handleItemToggle(itemModel *model.ItemModel, msg model.ItemToggleMsg) ItemUpdateResult {
	if msg.ID == itemModel.Todo().ID {
		itemModel.Toggle()
		
		// Generate async command to persist the toggle
		return ItemUpdateResult{
			Model: itemModel,
			Cmd: func() tea.Msg {
				return ItemToggleAsyncMsg{
					ID:     msg.ID,
					NewDone: itemModel.Todo().Done,
				}
			},
		}
	}
	
	return ItemUpdateResult{Model: itemModel, Cmd: nil}
}

// handleItemEdit handles edit operations for items
func handleItemEdit(itemModel *model.ItemModel, msg model.ItemEditMsg) ItemUpdateResult {
	if msg.ID == itemModel.Todo().ID {
		itemModel.SetEditing(msg.Editing)
		
		// Generate command to notify state model about edit mode change
		return ItemUpdateResult{
			Model: itemModel,
			Cmd: func() tea.Msg {
				return ItemEditStateMsg{
					ID:      msg.ID,
					Editing: msg.Editing,
				}
			},
		}
	}
	
	return ItemUpdateResult{Model: itemModel, Cmd: nil}
}

// handleItemSelect handles selection operations for items
func handleItemSelect(itemModel *model.ItemModel, msg model.ItemSelectMsg) ItemUpdateResult {
	if msg.ID == itemModel.Todo().ID {
		itemModel.SetSelected(msg.Selected)
	}
	
	return ItemUpdateResult{Model: itemModel, Cmd: nil}
}

// handleItemUpdate handles update operations for items
func handleItemUpdate(itemModel *model.ItemModel, msg model.ItemUpdateMsg) ItemUpdateResult {
	if msg.Todo != nil && itemModel.Todo() != nil && msg.Todo.ID == itemModel.Todo().ID {
		itemModel.SetTodo(msg.Todo)
		
		// Generate command to notify about successful update
		return ItemUpdateResult{
			Model: itemModel,
			Cmd: func() tea.Msg {
				return ItemUpdatedAsyncMsg{Todo: msg.Todo}
			},
		}
	}
	
	return ItemUpdateResult{Model: itemModel, Cmd: nil}
}

// BatchUpdateItems updates multiple items with the same message
func BatchUpdateItems(items []*model.ItemModel, msg tea.Msg) ([]tea.Cmd, bool) {
	var commands []tea.Cmd
	updated := false
	
	for _, item := range items {
		result := UpdateItem(item, msg)
		if result.Cmd != nil {
			commands = append(commands, result.Cmd)
			updated = true
		}
	}
	
	return commands, updated
}

// CreateItemAddCommand creates a command to add a new item
func CreateItemAddCommand(description string) tea.Cmd {
	return func() tea.Msg {
		return ItemAddMsg{Description: description}
	}
}

// CreateItemDeleteCommand creates a command to delete an item
func CreateItemDeleteCommand(id int) tea.Cmd {
	return func() tea.Msg {
		return ItemDeleteMsg{ID: id}
	}
}

// CreateItemToggleCommand creates a command to toggle an item
func CreateItemToggleCommand(id int) tea.Cmd {
	return func() tea.Msg {
		return model.ItemToggleMsg{ID: id}
	}
}

// CreateItemSelectCommand creates a command to select an item
func CreateItemSelectCommand(id int, selected bool) tea.Cmd {
	return func() tea.Msg {
		return model.ItemSelectMsg{ID: id, Selected: selected}
	}
}

// Message types for item operations

// ItemAddMsg is sent when a new item should be added
type ItemAddMsg struct {
	Description string
}

// ItemDeleteMsg is sent when an item should be deleted
type ItemDeleteMsg struct {
	ID int
}

// ItemToggleAsyncMsg is sent when an item toggle operation should be persisted
type ItemToggleAsyncMsg struct {
	ID      int
	NewDone bool
}

// ItemEditStateMsg is sent when an item's edit state changes
type ItemEditStateMsg struct {
	ID      int
	Editing bool
}

// ItemUpdatedAsyncMsg is sent when an item has been successfully updated
type ItemUpdatedAsyncMsg struct {
	Todo *domain.Todo
}