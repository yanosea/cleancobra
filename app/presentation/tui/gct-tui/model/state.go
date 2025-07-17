package model

import (
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/pkg/proxy"
	tea "github.com/charmbracelet/bubbletea"
)

// StateModel represents the main application state model for the TUI
type StateModel struct {
	// Todo management
	todos []*ItemModel
	
	// UI state components
	mode       Mode
	navigation *NavigationState
	input      *InputState
	
	// Use cases for business logic
	addUseCase    *application.AddTodoUseCase
	listUseCase   *application.ListTodoUseCase
	toggleUseCase *application.ToggleTodoUseCase
	deleteUseCase *application.DeleteTodoUseCase
	
	// Confirmation state
	confirmationMessage string
	confirmationAction  func() tea.Cmd
	
	// Error state
	errorMessage string
	
	// Dimensions
	width  int
	height int
}

// NewStateModel creates a new StateModel with the given use cases
func NewStateModel(
	addUseCase *application.AddTodoUseCase,
	listUseCase *application.ListTodoUseCase,
	toggleUseCase *application.ToggleTodoUseCase,
	deleteUseCase *application.DeleteTodoUseCase,
	bubbles proxy.Bubbles,
) *StateModel {
	return &StateModel{
		todos:         make([]*ItemModel, 0),
		mode:          ModeNormal,
		navigation:    NewNavigationState(),
		input:         NewInputState(bubbles),
		addUseCase:    addUseCase,
		listUseCase:   listUseCase,
		toggleUseCase: toggleUseCase,
		deleteUseCase: deleteUseCase,
		width:         80,
		height:        24,
	}
}

// Todos returns the list of todo item models
func (m *StateModel) Todos() []*ItemModel {
	return m.todos
}

// SetTodos updates the list of todo item models
func (m *StateModel) SetTodos(todos []*ItemModel) {
	m.todos = todos
	// Adjust cursor if it's out of bounds using navigation state
	m.navigation.SetCursor(m.navigation.Cursor(), len(m.todos))
}

// LoadTodos loads todos from the use case and converts them to item models
func (m *StateModel) LoadTodos() tea.Cmd {
	return func() tea.Msg {
		todos, err := m.listUseCase.Run()
		if err != nil {
			return ErrorMsg{Error: err}
		}
		
		itemModels := make([]*ItemModel, len(todos))
		for i, todo := range todos {
			itemModels[i] = NewItemModel(&todo)
		}
		
		return TodosLoadedMsg{Todos: itemModels}
	}
}

// Cursor returns the current cursor position
func (m *StateModel) Cursor() int {
	return m.navigation.Cursor()
}

// SetCursor sets the cursor position
func (m *StateModel) SetCursor(cursor int) {
	m.navigation.SetCursor(cursor, len(m.todos))
	
	// Update selection state
	m.clearSelection()
	if cursor >= 0 && cursor < len(m.todos) {
		m.todos[cursor].SetSelected(true)
		m.navigation.SetSelected(cursor, true)
	}
}

// MoveCursorUp moves the cursor up by one position
func (m *StateModel) MoveCursorUp() {
	m.navigation.MoveCursorUp(len(m.todos))
	m.updateSelectionFromNavigation()
}

// MoveCursorDown moves the cursor down by one position
func (m *StateModel) MoveCursorDown() {
	m.navigation.MoveCursorDown(len(m.todos))
	m.updateSelectionFromNavigation()
}

// MoveCursorToTop moves the cursor to the first position
func (m *StateModel) MoveCursorToTop() {
	m.navigation.MoveCursorToTop()
	m.updateSelectionFromNavigation()
}

// MoveCursorToBottom moves the cursor to the last position
func (m *StateModel) MoveCursorToBottom() {
	m.navigation.MoveCursorToBottom(len(m.todos))
	m.updateSelectionFromNavigation()
}

// clearSelection clears all selection states
func (m *StateModel) clearSelection() {
	for _, todo := range m.todos {
		todo.SetSelected(false)
	}
	m.navigation.ClearSelection()
}

// updateSelectionFromNavigation updates todo selection based on navigation state
func (m *StateModel) updateSelectionFromNavigation() {
	m.clearSelection()
	cursor := m.navigation.Cursor()
	if cursor >= 0 && cursor < len(m.todos) {
		m.todos[cursor].SetSelected(true)
		m.navigation.SetSelected(cursor, true)
	}
}

// Mode returns the current mode
func (m *StateModel) Mode() Mode {
	return m.mode
}

// SetMode sets the current mode
func (m *StateModel) SetMode(mode Mode) {
	m.mode = mode
	
	switch mode {
	case ModeInput:
		m.input.Clear()
		m.input.Focus()
	case ModeEdit:
		cursor := m.navigation.Cursor()
		if cursor >= 0 && cursor < len(m.todos) && m.todos[cursor].Todo() != nil {
			m.input.SetValue(m.todos[cursor].Todo().Description)
			m.input.Focus()
		}
	default:
		m.input.Blur()
	}
}

// Input returns the input state
func (m *StateModel) Input() *InputState {
	return m.input
}

// AddTodo adds a new todo with the given description
func (m *StateModel) AddTodo(description string) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.addUseCase.Run(description)
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return TodoAddedMsg{Todo: todo}
	}
}

// ToggleTodo toggles the completion status of the todo at the current cursor position
func (m *StateModel) ToggleTodo() tea.Cmd {
	cursor := m.navigation.Cursor()
	if cursor < 0 || cursor >= len(m.todos) || m.todos[cursor].Todo() == nil {
		return nil
	}
	
	todoID := m.todos[cursor].Todo().ID
	return func() tea.Msg {
		todo, err := m.toggleUseCase.Run(todoID)
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return TodoToggledMsg{Todo: todo}
	}
}

// DeleteTodo deletes the todo at the current cursor position
func (m *StateModel) DeleteTodo() tea.Cmd {
	cursor := m.navigation.Cursor()
	if cursor < 0 || cursor >= len(m.todos) || m.todos[cursor].Todo() == nil {
		return nil
	}
	
	todoID := m.todos[cursor].Todo().ID
	return func() tea.Msg {
		err := m.deleteUseCase.Run(todoID)
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return TodoDeletedMsg{ID: todoID}
	}
}

// UpdateTodo updates the description of the todo at the current cursor position
func (m *StateModel) UpdateTodo(description string) tea.Cmd {
	cursor := m.navigation.Cursor()
	if cursor < 0 || cursor >= len(m.todos) || m.todos[cursor].Todo() == nil {
		return nil
	}
	
	todo := m.todos[cursor].Todo()
	if err := todo.UpdateDescription(description); err != nil {
		return func() tea.Msg {
			return ErrorMsg{Error: err}
		}
	}
	
	// Save the updated todo
	return func() tea.Msg {
		// Note: We need to save through repository, but for now we'll simulate
		// In a real implementation, we'd need an UpdateTodoUseCase
		return TodoUpdatedMsg{Todo: todo}
	}
}

// SetConfirmation sets up a confirmation dialog
func (m *StateModel) SetConfirmation(message string, action func() tea.Cmd) {
	m.confirmationMessage = message
	m.confirmationAction = action
	m.SetMode(ModeConfirmation)
}

// ConfirmationMessage returns the current confirmation message
func (m *StateModel) ConfirmationMessage() string {
	return m.confirmationMessage
}

// ExecuteConfirmation executes the confirmed action
func (m *StateModel) ExecuteConfirmation() tea.Cmd {
	if m.confirmationAction != nil {
		cmd := m.confirmationAction()
		m.confirmationAction = nil
		m.confirmationMessage = ""
		m.SetMode(ModeNormal)
		return cmd
	}
	return nil
}

// CancelConfirmation cancels the confirmation dialog
func (m *StateModel) CancelConfirmation() {
	m.confirmationAction = nil
	m.confirmationMessage = ""
	m.SetMode(ModeNormal)
}

// SetError sets an error message
func (m *StateModel) SetError(err error) {
	if err != nil {
		m.errorMessage = err.Error()
	} else {
		m.errorMessage = ""
	}
}

// ErrorMessage returns the current error message
func (m *StateModel) ErrorMessage() string {
	return m.errorMessage
}

// ClearError clears the current error message
func (m *StateModel) ClearError() {
	m.errorMessage = ""
}

// SetSize sets the dimensions of the TUI
func (m *StateModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Adjust input width based on available space
	if width > 10 {
		m.input.SetWidth(width - 10)
	}
}

// Width returns the current width
func (m *StateModel) Width() int {
	return m.width
}

// Height returns the current height
func (m *StateModel) Height() int {
	return m.height
}

// Init implements tea.Model interface
func (m *StateModel) Init() tea.Cmd {
	return m.LoadTodos()
}

// Update implements tea.Model interface
func (m *StateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil
		
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
		
	case TodosLoadedMsg:
		m.SetTodos(msg.Todos)
		if len(m.todos) > 0 {
			m.SetCursor(0)
		}
		return m, nil
		
	case TodoAddedMsg:
		// Reload todos to get the updated list
		return m, m.LoadTodos()
		
	case TodoToggledMsg:
		// Update the specific todo in our list
		for _, itemModel := range m.todos {
			if itemModel.Todo() != nil && itemModel.Todo().ID == msg.Todo.ID {
				itemModel.SetTodo(msg.Todo)
				break
			}
		}
		return m, nil
		
	case TodoDeletedMsg:
		// Reload todos to get the updated list
		return m, m.LoadTodos()
		
	case TodoUpdatedMsg:
		// Update the specific todo in our list
		for _, itemModel := range m.todos {
			if itemModel.Todo() != nil && itemModel.Todo().ID == msg.Todo.ID {
				itemModel.SetTodo(msg.Todo)
				break
			}
		}
		m.SetMode(ModeNormal)
		return m, nil
		
	case ErrorMsg:
		m.SetError(msg.Error)
		m.SetMode(ModeNormal)
		return m, nil
	}
	
	// Handle input updates when in input modes
	if m.mode == ModeInput || m.mode == ModeEdit {
		var cmd tea.Cmd
		textInput := m.input.TextInput()
		textInput, cmd = textInput.Update(msg)
		return m, cmd
	}
	
	return m, nil
}

// View implements tea.Model interface
func (m *StateModel) View() string {
	// Basic view implementation - will be enhanced in view layer
	switch m.mode {
	case ModeInput:
		return "Add Todo: " + m.input.View()
	case ModeEdit:
		return "Edit Todo: " + m.input.View()
	case ModeConfirmation:
		return m.confirmationMessage + " (y/n)"
	default:
		view := "Todos:\n"
		cursor := m.navigation.Cursor()
		for i, todo := range m.todos {
			prefix := "  "
			if i == cursor {
				prefix = "> "
			}
			view += prefix + todo.View() + "\n"
		}
		if m.errorMessage != "" {
			view += "\nError: " + m.errorMessage
		}
		return view
	}
}

// handleKeyMsg handles keyboard input based on current mode
func (m *StateModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeNormal:
		return m.handleNormalModeKeys(msg)
	case ModeInput:
		return m.handleInputModeKeys(msg)
	case ModeEdit:
		return m.handleEditModeKeys(msg)
	case ModeConfirmation:
		return m.handleConfirmationModeKeys(msg)
	}
	return m, nil
}

// handleNormalModeKeys handles keyboard input in normal mode
func (m *StateModel) handleNormalModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		m.MoveCursorUp()
		return m, nil
	case "down", "j":
		m.MoveCursorDown()
		return m, nil
	case "g":
		m.MoveCursorToTop()
		return m, nil
	case "G":
		m.MoveCursorToBottom()
		return m, nil
	case " ":
		return m, m.ToggleTodo()
	case "a":
		m.SetMode(ModeInput)
		return m, nil
	case "e":
		if len(m.todos) > 0 {
			m.SetMode(ModeEdit)
		}
		return m, nil
	case "d":
		cursor := m.navigation.Cursor()
		if len(m.todos) > 0 && cursor >= 0 && cursor < len(m.todos) {
			todo := m.todos[cursor].Todo()
			if todo != nil {
				m.SetConfirmation(
					"Delete '"+todo.Description+"'?",
					m.DeleteTodo,
				)
			}
		}
		return m, nil
	case "r":
		return m, m.LoadTodos()
	case "esc":
		m.ClearError()
		return m, nil
	}
	return m, nil
}

// handleInputModeKeys handles keyboard input in input mode
func (m *StateModel) handleInputModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		description := m.input.Value()
		if description != "" {
			m.SetMode(ModeNormal)
			return m, m.AddTodo(description)
		}
		return m, nil
	case "esc":
		m.SetMode(ModeNormal)
		return m, nil
	}
	return m, nil
}

// handleEditModeKeys handles keyboard input in edit mode
func (m *StateModel) handleEditModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		description := m.input.Value()
		if description != "" {
			return m, m.UpdateTodo(description)
		}
		return m, nil
	case "esc":
		m.SetMode(ModeNormal)
		return m, nil
	}
	return m, nil
}

// handleConfirmationModeKeys handles keyboard input in confirmation mode
func (m *StateModel) handleConfirmationModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m, m.ExecuteConfirmation()
	case "n", "N", "esc":
		m.CancelConfirmation()
		return m, nil
	}
	return m, nil
}

