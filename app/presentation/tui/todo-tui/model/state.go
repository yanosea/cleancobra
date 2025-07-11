package model

import todoApp "github.com/yanosea/cleancobra/app/application/todo"

type Mode int

const (
	ModeList Mode = iota
	ModeAdd
	ModeDelete
)

type State struct {
	mode     Mode
	cursor   int
	input    string
	width    int
	height   int
	quitting bool

	todos   []*todoApp.ListTodoUsecaseOutputDto
	error   string
	message string

	confirmButtonSelected bool
}

func NewState() *State {
	return &State{
		mode:                  ModeList,
		cursor:                0,
		input:                 "",
		width:                 80,
		height:                24,
		quitting:              false,
		todos:                 make([]*todoApp.ListTodoUsecaseOutputDto, 0),
		error:                 "",
		message:               "",
		confirmButtonSelected: false,
	}
}

func (s *State) Mode() Mode        { return s.mode }
func (s *State) SetMode(mode Mode) { s.mode = mode }

func (s *State) Cursor() int          { return s.cursor }
func (s *State) SetCursor(cursor int) { s.cursor = cursor }
func (s *State) MoveCursorUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}
func (s *State) MoveCursorDown() {
	if s.cursor < len(s.todos)-1 {
		s.cursor++
	}
}

func (s *State) Input() string             { return s.input }
func (s *State) SetInput(input string)     { s.input = input }
func (s *State) ResetInput()               { s.input = "" }
func (s *State) AppendToInput(text string) { s.input += text }
func (s *State) Backspace() {
	if len(s.input) > 0 {
		s.input = s.input[:len(s.input)-1]
	}
}

func (s *State) Width() int  { return s.width }
func (s *State) Height() int { return s.height }
func (s *State) SetDimensions(width, height int) {
	s.width = width
	s.height = height
}

func (s *State) Quitting() bool            { return s.quitting }
func (s *State) SetQuitting(quitting bool) { s.quitting = quitting }

func (s *State) Todos() []*todoApp.ListTodoUsecaseOutputDto         { return s.todos }
func (s *State) SetTodos(todos []*todoApp.ListTodoUsecaseOutputDto) { s.todos = todos }
func (s *State) CurrentTodo() *todoApp.ListTodoUsecaseOutputDto {
	if s.cursor >= 0 && s.cursor < len(s.todos) {
		return s.todos[s.cursor]
	}
	return nil
}

func (s *State) Error() string             { return s.error }
func (s *State) Message() string           { return s.message }
func (s *State) SetError(error string)     { s.error = error }
func (s *State) SetMessage(message string) { s.message = message }
func (s *State) ClearMessages() {
	s.error = ""
	s.message = ""
}

func (s *State) ConfirmButtonSelected() bool            { return s.confirmButtonSelected }
func (s *State) SetConfirmButtonSelected(selected bool) { s.confirmButtonSelected = selected }
func (s *State) ToggleDeleteButton()                    { s.confirmButtonSelected = !s.confirmButtonSelected }
func (s *State) ResetDeleteButton()                     { s.confirmButtonSelected = false }

