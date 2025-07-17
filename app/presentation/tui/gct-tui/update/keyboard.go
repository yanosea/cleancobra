package update

import (
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyboardHandler handles keyboard input routing based on current mode
func KeyboardHandler(stateModel *model.StateModel, msg tea.KeyMsg) HandlerUpdateResult {
	switch stateModel.Mode() {
	case model.ModeNormal:
		return handleNormalModeKeys(stateModel, msg)
	case model.ModeInput:
		return handleInputModeKeys(stateModel, msg)
	case model.ModeEdit:
		return handleEditModeKeys(stateModel, msg)
	case model.ModeConfirmation:
		return handleConfirmationModeKeys(stateModel, msg)
	}
	
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// handleNormalModeKeys handles keyboard input in normal navigation mode
func handleNormalModeKeys(stateModel *model.StateModel, msg tea.KeyMsg) HandlerUpdateResult {
	switch msg.String() {
	case "ctrl+c", "q":
		return HandlerUpdateResult{Model: stateModel, Cmd: tea.Quit}
		
	case "up", "k":
		stateModel.MoveCursorUp()
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "down", "j":
		stateModel.MoveCursorDown()
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "g":
		stateModel.MoveCursorToTop()
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "G":
		stateModel.MoveCursorToBottom()
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case " ":
		return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.ToggleTodo()}
		
	case "a":
		stateModel.SetMode(model.ModeInput)
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "e":
		if len(stateModel.Todos()) > 0 && stateModel.Cursor() >= 0 && stateModel.Cursor() < len(stateModel.Todos()) {
			stateModel.SetMode(model.ModeEdit)
		}
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "d":
		if len(stateModel.Todos()) > 0 && stateModel.Cursor() >= 0 && stateModel.Cursor() < len(stateModel.Todos()) {
			todo := stateModel.Todos()[stateModel.Cursor()].Todo()
			if todo != nil {
				stateModel.SetConfirmation(
					"Delete '"+todo.Description+"'?",
					stateModel.DeleteTodo,
				)
			}
		}
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "r":
		return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.LoadTodos()}
		
	case "esc":
		stateModel.ClearError()
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
	}
	
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// handleInputModeKeys handles keyboard input in add todo mode
func handleInputModeKeys(stateModel *model.StateModel, msg tea.KeyMsg) HandlerUpdateResult {
	switch msg.String() {
	case "enter":
		description := stateModel.Input().Value()
		if description != "" {
			stateModel.SetMode(model.ModeNormal)
			return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.AddTodo(description)}
		}
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "esc":
		stateModel.SetMode(model.ModeNormal)
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
	}
	
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// handleEditModeKeys handles keyboard input in edit todo mode
func handleEditModeKeys(stateModel *model.StateModel, msg tea.KeyMsg) HandlerUpdateResult {
	switch msg.String() {
	case "enter":
		description := stateModel.Input().Value()
		if description != "" {
			return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.UpdateTodo(description)}
		}
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
		
	case "esc":
		stateModel.SetMode(model.ModeNormal)
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
	}
	
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// handleConfirmationModeKeys handles keyboard input in confirmation mode
func handleConfirmationModeKeys(stateModel *model.StateModel, msg tea.KeyMsg) HandlerUpdateResult {
	switch msg.String() {
	case "y", "Y":
		return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.ExecuteConfirmation()}
		
	case "n", "N", "esc":
		stateModel.CancelConfirmation()
		return HandlerUpdateResult{Model: stateModel, Cmd: nil}
	}
	
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

