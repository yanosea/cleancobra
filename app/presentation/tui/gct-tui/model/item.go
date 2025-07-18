package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanosea/gct/app/domain"
)

// ItemModel represents the model for individual todo items in the TUI
type ItemModel struct {
	todo     *domain.Todo
	selected bool
	editing  bool
}

// NewItemModel creates a new ItemModel with the given todo
func NewItemModel(todo *domain.Todo) *ItemModel {
	return &ItemModel{
		todo:     todo,
		selected: false,
		editing:  false,
	}
}

// Todo returns the underlying todo entity
func (m *ItemModel) Todo() *domain.Todo {
	return m.todo
}

// SetTodo updates the underlying todo entity
func (m *ItemModel) SetTodo(todo *domain.Todo) {
	m.todo = todo
}

// IsSelected returns whether this item is currently selected
func (m *ItemModel) IsSelected() bool {
	return m.selected
}

// SetSelected sets the selection state of this item
func (m *ItemModel) SetSelected(selected bool) {
	m.selected = selected
}

// IsEditing returns whether this item is currently being edited
func (m *ItemModel) IsEditing() bool {
	return m.editing
}

// SetEditing sets the editing state of this item
func (m *ItemModel) SetEditing(editing bool) {
	m.editing = editing
}

// Toggle toggles the completion status of the todo item
func (m *ItemModel) Toggle() {
	if m.todo != nil {
		m.todo.Toggle()
	}
}

// UpdateDescription updates the description of the todo item
func (m *ItemModel) UpdateDescription(description string) error {
	if m.todo != nil {
		return m.todo.UpdateDescription(description)
	}
	return nil
}

// Init implements tea.Model interface for individual item initialization
func (m *ItemModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model interface for individual item updates
func (m *ItemModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case ItemToggleMsg:
		if msg.ID == m.todo.ID {
			m.Toggle()
		}
		return m, nil
	case ItemEditMsg:
		if msg.ID == m.todo.ID {
			m.SetEditing(msg.Editing)
		}
		return m, nil
	case ItemSelectMsg:
		if msg.ID == m.todo.ID {
			m.SetSelected(msg.Selected)
		}
		return m, nil
	}
	return m, nil
}

// View implements tea.Model interface for individual item rendering
func (m *ItemModel) View() string {
	if m.todo == nil {
		return ""
	}

	// Basic rendering - will be enhanced in view layer
	status := "[ ]"
	if m.todo.Done {
		status = "[x]"
	}

	prefix := "  "
	if m.selected {
		prefix = "> "
	}

	return prefix + status + " " + m.todo.Description
}

// handleKeyMsg handles keyboard input for individual items
func (m *ItemModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.selected {
		return m, nil
	}

	switch msg.String() {
	case " ":
		// Toggle completion status
		return m, func() tea.Msg {
			return ItemToggleMsg{ID: m.todo.ID}
		}
	case "e":
		// Enter edit mode
		return m, func() tea.Msg {
			return ItemEditMsg{ID: m.todo.ID, Editing: true}
		}
	}

	return m, nil
}
