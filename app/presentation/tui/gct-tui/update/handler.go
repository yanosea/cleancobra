package update

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
)

// HandlerUpdateResult represents the result of a handler update operation
type HandlerUpdateResult struct {
	Model *model.StateModel
	Cmd   tea.Cmd
}

// UpdateHandler handles the main update function with message routing
func UpdateHandler(stateModel *model.StateModel, msg tea.Msg) HandlerUpdateResult {
	// Create operations handler
	opsHandler := NewOperationsHandler()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return handleWindowSize(stateModel, msg)
	case tea.KeyMsg:
		return KeyboardHandler(stateModel, msg)
	case model.TodosLoadedMsg:
		return opsHandler.HandleTodosLoaded(stateModel, msg)
	case model.TodoAddedMsg:
		return opsHandler.HandleTodoAdded(stateModel, msg)
	case model.TodoToggledMsg:
		return opsHandler.HandleTodoToggled(stateModel, msg)
	case model.TodoDeletedMsg:
		return opsHandler.HandleTodoDeleted(stateModel, msg)
	case model.TodoUpdatedMsg:
		return opsHandler.HandleTodoUpdated(stateModel, msg)
	case model.ErrorMsg:
		return opsHandler.HandleError(stateModel, msg)
	case ItemAddMsg:
		return opsHandler.HandleItemAdd(stateModel, msg)
	case ItemDeleteMsg:
		return opsHandler.HandleItemDelete(stateModel, msg)
	case ItemToggleAsyncMsg:
		return opsHandler.HandleItemToggleAsync(stateModel, msg)
	case ItemEditStateMsg:
		return opsHandler.HandleItemEditState(stateModel, msg)
	case ItemUpdatedAsyncMsg:
		return opsHandler.HandleItemUpdatedAsync(stateModel, msg)
	}

	// Handle input updates when in input modes
	if stateModel.Mode() == model.ModeInput || stateModel.Mode() == model.ModeEdit {
		return handleInputUpdate(stateModel, msg)
	}

	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// handleWindowSize handles window size changes
func handleWindowSize(stateModel *model.StateModel, msg tea.WindowSizeMsg) HandlerUpdateResult {
	stateModel.SetSize(msg.Width, msg.Height)
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// handleInputUpdate handles text input updates in input/edit modes
func handleInputUpdate(stateModel *model.StateModel, msg tea.Msg) HandlerUpdateResult {
	var cmd tea.Cmd
	textInput := stateModel.Input().TextInput()
	textInput, cmd = textInput.Update(msg)

	return HandlerUpdateResult{Model: stateModel, Cmd: cmd}
}
