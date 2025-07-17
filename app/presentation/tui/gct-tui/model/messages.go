package model

import (
	"github.com/yanosea/gct/app/domain"
)

// Message types for the TUI application

// TodosLoadedMsg is sent when todos are loaded from the repository
type TodosLoadedMsg struct {
	Todos []*ItemModel
}

// TodoAddedMsg is sent when a new todo is added
type TodoAddedMsg struct {
	Todo *domain.Todo
}

// TodoToggledMsg is sent when a todo's completion status is toggled
type TodoToggledMsg struct {
	Todo *domain.Todo
}

// TodoDeletedMsg is sent when a todo is deleted
type TodoDeletedMsg struct {
	ID int
}

// TodoUpdatedMsg is sent when a todo is updated
type TodoUpdatedMsg struct {
	Todo *domain.Todo
}

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	Error error
}

// ItemToggleMsg is a message to toggle an item's completion status
type ItemToggleMsg struct {
	ID int
}

// ItemEditMsg is a message to set an item's editing state
type ItemEditMsg struct {
	ID      int
	Editing bool
}

// ItemSelectMsg is a message to set an item's selection state
type ItemSelectMsg struct {
	ID       int
	Selected bool
}

// ItemUpdateMsg is a message to update an item's todo data
type ItemUpdateMsg struct {
	Todo *domain.Todo
}

// ConfirmationMsg is sent when a confirmation dialog is requested
type ConfirmationMsg struct {
	Message string
	Action  func() interface{}
}

// ConfirmationAcceptedMsg is sent when a confirmation is accepted
type ConfirmationAcceptedMsg struct{}

// ConfirmationCancelledMsg is sent when a confirmation is cancelled
type ConfirmationCancelledMsg struct{}

// ModeChangeMsg is sent when the application mode changes
type ModeChangeMsg struct {
	Mode Mode
}

// NavigationMsg is sent when navigation occurs
type NavigationMsg struct {
	Direction string // "up", "down", "top", "bottom"
}

// ClearErrorMsg is sent to clear error messages
type ClearErrorMsg struct{}

// RefreshMsg is sent to refresh the todo list
type RefreshMsg struct{}