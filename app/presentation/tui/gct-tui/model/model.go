package model

import (
	"fmt"
	"strings"

	todoApp "github.com/yanosea/gct/app/application/gct"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/formatter"
	"github.com/yanosea/gct/pkg/proxy"
)

type Model struct {
	state    *State
	usecases *Usecases
}

type Usecases struct {
	List   *todoApp.ListTodoUseCase
	Add    *todoApp.AddTodoUseCase
	Delete *todoApp.DeleteTodoUseCase
	Toggle *todoApp.ToggleTodoUseCase
}

func NewModel(usecases *Usecases) *Model {
	return &Model{
		state:    NewState(),
		usecases: usecases,
	}
}

func (m *Model) Init() proxy.Cmd {
	return m.loadTodos()
}

func (m *Model) Update(msg proxy.Msg) (proxy.Model, proxy.Cmd) {
	return m.updateModel(msg)
}

func (m *Model) View() string {
	return m.renderView()
}

func (m *Model) State() *State {
	return m.state
}

func (m *Model) Usecases() *Usecases {
	return m.usecases
}

func (m *Model) updateModel(msg proxy.Msg) (*Model, proxy.Cmd) {
	state := m.state

	if windowSizeMsg, ok := proxy.IsWindowSizeMsg(msg); ok {
		state.SetDimensions(windowSizeMsg.GetWidth(), windowSizeMsg.GetHeight())
		return m, nil
	}

	if keyMsg, ok := proxy.IsKeyMsg(msg); ok {
		return m.handleKeyPress(keyMsg)
	}

	switch msg := msg.(type) {
	case TodosLoadedMsg:
		state.SetTodos(msg.Todos)
		state.SetError("")
		return m, nil

	case ErrorMsg:
		state.SetError(msg.Error)
		return m, nil

	case SuccessMsg:
		state.SetMessage(msg.Message)
		state.SetError("")
		return m, m.loadTodos()

	default:
		return m, nil
	}
}

func (m *Model) handleKeyPress(keyMsg proxy.KeyMsg) (*Model, proxy.Cmd) {
	switch m.state.Mode() {
	case ModeList:
		return m.handleListMode(keyMsg)
	case ModeAdd:
		return m.handleAddMode(keyMsg)
	case ModeDelete:
		return m.handleDeleteMode(keyMsg)
	default:
		return m, nil
	}
}

func (m *Model) handleListMode(keyMsg proxy.KeyMsg) (*Model, proxy.Cmd) {
	state := m.state

	switch keyMsg.String() {
	case "ctrl+c", "q":
		state.SetQuitting(true)
		return m, proxy.Quit()

	case "up", "k":
		state.MoveCursorUp()

	case "down", "j":
		state.MoveCursorDown()

	case "enter", " ":
		if todo := state.CurrentTodo(); todo != nil {
			return m, m.toggleTodo(todo.ID)
		}

	case "a":
		state.SetMode(ModeAdd)
		state.ResetInput()
		state.ClearMessages()

	case "d":
		if state.CurrentTodo() != nil {
			state.SetMode(ModeDelete)
			state.ResetDeleteButton()
			state.ClearMessages()
		}

	case "r":
		return m, m.loadTodos()
	}

	return m, nil
}

func (m *Model) handleAddMode(keyMsg proxy.KeyMsg) (*Model, proxy.Cmd) {
	state := m.state

	switch keyMsg.String() {
	case "ctrl+c":
		state.SetQuitting(true)
		return m, proxy.Quit()

	case "esc":
		state.SetMode(ModeList)
		state.ResetInput()
		state.ClearMessages()

	case "enter":
		if len(state.Input()) > 0 {
			cmd := m.addTodo(state.Input())
			state.SetMode(ModeList)
			state.ResetInput()
			return m, cmd
		}

	case "backspace":
		state.Backspace()

	default:
		state.AppendToInput(keyMsg.String())
	}

	return m, nil
}

func (m *Model) handleDeleteMode(keyMsg proxy.KeyMsg) (*Model, proxy.Cmd) {
	state := m.state

	switch keyMsg.String() {
	case "ctrl+c":
		state.SetQuitting(true)
		return m, proxy.Quit()

	case "esc":
		state.SetMode(ModeList)
		state.ClearMessages()

	case "left", "right", "tab", "h", "l":
		state.ToggleDeleteButton()

	case "enter":
		if state.ConfirmButtonSelected() {
			if todo := state.CurrentTodo(); todo != nil {
				cmd := m.deleteTodo(todo.ID)
				state.SetMode(ModeList)
				return m, cmd
			}
			return m, nil
		} else {
			state.SetMode(ModeList)
			state.ClearMessages()
			return m, nil
		}

	case "y":
		if todo := state.CurrentTodo(); todo != nil {
			cmd := m.deleteTodo(todo.ID)
			state.SetMode(ModeList)
			return m, cmd
		}

	case "n":
		state.SetMode(ModeList)
		state.ClearMessages()
	}

	return m, nil
}

func (m *Model) loadTodos() proxy.Cmd {
	return func() proxy.Msg {
		output, err := m.usecases.List.Run()
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}
		return TodosLoadedMsg{Todos: output}
	}
}

func (m *Model) addTodo(title string) proxy.Cmd {
	return func() proxy.Msg {
		_, err := m.usecases.Add.Run(title)
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}
		return SuccessMsg{Message: fmt.Sprintf("Added todo: %s", title)}
	}
}

func (m *Model) deleteTodo(id string) proxy.Cmd {
	return func() proxy.Msg {
		output, err := m.usecases.Delete.Run(id)
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}
		return SuccessMsg{Message: fmt.Sprintf("Deleted todo: %s", output.Title)}
	}
}

func (m *Model) toggleTodo(id string) proxy.Cmd {
	return func() proxy.Msg {
		output, err := m.usecases.Toggle.Run(id)
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}

		status := "incomplete"
		if output.Done {
			status = "complete"
		}

		return SuccessMsg{Message: fmt.Sprintf("Marked todo as %s: %s", status, output.Title)}
	}
}

func (m *Model) renderView() string {
	state := m.state

	if state.Quitting() {
		return "Goodbye!\n"
	}

	var content strings.Builder

	header := formatter.FormatHeader("📝 Todo TUI")
	content.WriteString(header + "\n\n")

	if state.Error() != "" {
		content.WriteString(formatter.FormatError(state.Error()) + "\n\n")
	}

	if state.Message() != "" {
		content.WriteString(formatter.FormatSuccess(state.Message()) + "\n\n")
	}

	switch state.Mode() {
	case ModeList:
		content.WriteString(m.renderListView())
	case ModeAdd:
		content.WriteString(m.renderAddView())
	case ModeDelete:
		content.WriteString(m.renderDeleteView())
	}

	content.WriteString("\n" + m.renderHelpView())

	return formatter.BaseStyle.
		Width(state.Width() - 4).
		Height(state.Height() - 4).
		Render(content.String())
}

func (m *Model) renderListView() string {
	var content strings.Builder
	state := m.state

	if len(state.Todos()) == 0 {
		content.WriteString("No todos found. Press 'a' to add a new todo.\n")
	} else {
		for i, todo := range state.Todos() {
			selected := i == state.Cursor()
			todoItem := formatter.FormatTodoItem(todo.Title, todo.Done, selected)
			content.WriteString(todoItem + "\n")
		}
	}

	return content.String()
}

func (m *Model) renderAddView() string {
	state := m.state

	return fmt.Sprintf(`Add a new todo:

%s
`, formatter.FormatInput(state.Input()))
}

func (m *Model) renderDeleteView() string {
	state := m.state

	todo := state.CurrentTodo()
	if todo == nil {
		return "No todo selected.\n"
	}

	warningHeader := formatter.FormatDanger("⚠️ DELETE CONFIRMATION ⚠️")

	questionText := "Are you sure you want to delete this todo?"

	highlightedTodo := formatter.FormatHighlightedTodo(todo.Title, todo.Done)

	confirmButton := formatter.FormatConfirmButton("YES, DELETE", state.ConfirmButtonSelected())
	cancelButton := formatter.FormatCancelButton("NO, CANCEL", !state.ConfirmButtonSelected())
	buttonsRow := confirmButton + " " + cancelButton

	instructions := `This action cannot be undone!

← → or TAB or h/l: Switch between buttons
ENTER: Execute selected action
ESC: Cancel deletion`
	warningBox := formatter.FormatWarningBox(instructions)

	return fmt.Sprintf(`%s

%s

%s

%s

%s`, warningHeader, questionText, highlightedTodo, buttonsRow, warningBox)
}

func (m *Model) renderHelpView() string {
	switch m.state.Mode() {
	case ModeList:
		return formatter.FormatHelp("↑/k: up • ↓/j: down • enter/space: toggle • a: add • d: delete • r: refresh • q: quit")
	case ModeAdd:
		return formatter.FormatHelp("enter: add todo • esc: cancel • ctrl+c: quit")
	case ModeDelete:
		return formatter.FormatHelp("←→/tab/h/l: switch buttons • enter: execute • y: quick confirm • n/esc: cancel • ctrl+c: quit")
	default:
		return ""
	}
}
