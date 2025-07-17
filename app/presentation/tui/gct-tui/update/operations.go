package update

import (
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	tea "github.com/charmbracelet/bubbletea"
)

// OperationsHandler handles todo operations (add, delete, toggle, update)
type OperationsHandler struct{}

// NewOperationsHandler creates a new operations handler
func NewOperationsHandler() *OperationsHandler {
	return &OperationsHandler{}
}

// HandleTodosLoaded handles the todos loaded message
func (h *OperationsHandler) HandleTodosLoaded(stateModel *model.StateModel, msg model.TodosLoadedMsg) HandlerUpdateResult {
	stateModel.SetTodos(msg.Todos)
	if len(stateModel.Todos()) > 0 {
		stateModel.SetCursor(0)
	}
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// HandleTodoAdded handles the todo added message
func (h *OperationsHandler) HandleTodoAdded(stateModel *model.StateModel, msg model.TodoAddedMsg) HandlerUpdateResult {
	// Reload todos to get the updated list
	return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.LoadTodos()}
}

// HandleTodoToggled handles the todo toggled message
func (h *OperationsHandler) HandleTodoToggled(stateModel *model.StateModel, msg model.TodoToggledMsg) HandlerUpdateResult {
	// Update the specific todo in our list
	for _, itemModel := range stateModel.Todos() {
		if itemModel.Todo() != nil && itemModel.Todo().ID == msg.Todo.ID {
			itemModel.SetTodo(msg.Todo)
			break
		}
	}
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// HandleTodoDeleted handles the todo deleted message
func (h *OperationsHandler) HandleTodoDeleted(stateModel *model.StateModel, msg model.TodoDeletedMsg) HandlerUpdateResult {
	// Reload todos to get the updated list
	return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.LoadTodos()}
}

// HandleTodoUpdated handles the todo updated message
func (h *OperationsHandler) HandleTodoUpdated(stateModel *model.StateModel, msg model.TodoUpdatedMsg) HandlerUpdateResult {
	// Update the specific todo in our list
	for _, itemModel := range stateModel.Todos() {
		if itemModel.Todo() != nil && itemModel.Todo().ID == msg.Todo.ID {
			itemModel.SetTodo(msg.Todo)
			break
		}
	}
	stateModel.SetMode(model.ModeNormal)
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// HandleError handles error messages
func (h *OperationsHandler) HandleError(stateModel *model.StateModel, msg model.ErrorMsg) HandlerUpdateResult {
	stateModel.SetError(msg.Error)
	stateModel.SetMode(model.ModeNormal)
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// HandleItemAdd handles item add messages from item update layer
func (h *OperationsHandler) HandleItemAdd(stateModel *model.StateModel, msg ItemAddMsg) HandlerUpdateResult {
	return HandlerUpdateResult{Model: stateModel, Cmd: stateModel.AddTodo(msg.Description)}
}

// HandleItemDelete handles item delete messages from item update layer
func (h *OperationsHandler) HandleItemDelete(stateModel *model.StateModel, msg ItemDeleteMsg) HandlerUpdateResult {
	// Find the todo by ID and set up confirmation
	for _, itemModel := range stateModel.Todos() {
		if itemModel.Todo() != nil && itemModel.Todo().ID == msg.ID {
			stateModel.SetConfirmation(
				"Delete '"+itemModel.Todo().Description+"'?",
				stateModel.DeleteTodo,
			)
			break
		}
	}
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// HandleItemToggleAsync handles async toggle operations from item update layer
func (h *OperationsHandler) HandleItemToggleAsync(stateModel *model.StateModel, msg ItemToggleAsyncMsg) HandlerUpdateResult {
	// Find and toggle the specific todo
	for _, itemModel := range stateModel.Todos() {
		if itemModel.Todo() != nil && itemModel.Todo().ID == msg.ID {
			// The toggle has already been applied to the model, now persist it
			return HandlerUpdateResult{
				Model: stateModel,
				Cmd:   stateModel.ToggleTodo(),
			}
		}
	}
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// HandleItemEditState handles edit state changes from item update layer
func (h *OperationsHandler) HandleItemEditState(stateModel *model.StateModel, msg ItemEditStateMsg) HandlerUpdateResult {
	if msg.Editing {
		// Find the todo and switch to edit mode
		for i, itemModel := range stateModel.Todos() {
			if itemModel.Todo() != nil && itemModel.Todo().ID == msg.ID {
				stateModel.SetCursor(i)
				stateModel.SetMode(model.ModeEdit)
				break
			}
		}
	} else {
		stateModel.SetMode(model.ModeNormal)
	}
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// HandleItemUpdatedAsync handles async update operations from item update layer
func (h *OperationsHandler) HandleItemUpdatedAsync(stateModel *model.StateModel, msg ItemUpdatedAsyncMsg) HandlerUpdateResult {
	// The update has been applied, just acknowledge it
	return HandlerUpdateResult{Model: stateModel, Cmd: nil}
}

// CreateModeTransitionCommand creates a command to transition between modes
func CreateModeTransitionCommand(mode model.Mode) tea.Cmd {
	return func() tea.Msg {
		return ModeTransitionMsg{Mode: mode}
	}
}

// CreateNavigationCommand creates a command for navigation actions
func CreateNavigationCommand(action NavigationAction, value int) tea.Cmd {
	return func() tea.Msg {
		return NavigationMsg{Action: action, Value: value}
	}
}

// CreateAsyncOperationCommand creates a command for async operations
func CreateAsyncOperationCommand(operation AsyncOperation, data interface{}) tea.Cmd {
	return func() tea.Msg {
		return AsyncOperationMsg{Operation: operation, Data: data}
	}
}

// Message types for operations

// ModeTransitionMsg is sent when the application should transition between modes
type ModeTransitionMsg struct {
	Mode model.Mode
}

// NavigationMsg is sent for navigation operations
type NavigationMsg struct {
	Action NavigationAction
	Value  int
}

// NavigationAction represents different navigation actions
type NavigationAction int

const (
	NavigationUp NavigationAction = iota
	NavigationDown
	NavigationTop
	NavigationBottom
	NavigationTo
)

// AsyncOperationMsg is sent for async operations
type AsyncOperationMsg struct {
	Operation AsyncOperation
	Data      interface{}
}

// AsyncOperation represents different async operations
type AsyncOperation int

const (
	AsyncOperationAdd AsyncOperation = iota
	AsyncOperationToggle
	AsyncOperationDelete
	AsyncOperationUpdate
	AsyncOperationLoad
)