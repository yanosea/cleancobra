//go:generate mockgen -source=bubbles.go -destination=bubbles_mock.go -package=proxy

package proxy

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Bubbles provides a proxy interface for bubbles package functionality
type Bubbles interface {
	NewTextInput() TextInput
}

// TextInput provides a proxy interface for textinput.Model
type TextInput interface {
	SetValue(s string)
	Value() string
	SetPlaceholder(str string)
	Placeholder() string
	Focus() tea.Cmd
	Blur()
	Focused() bool
	SetPrompt(str string)
	Prompt() string
	SetCharLimit(limit int)
	CharLimit() int
	SetWidth(w int)
	Width() int
	Update(msg tea.Msg) (TextInput, tea.Cmd)
	View() string
}

// BubblesImpl implements the Bubbles interface using the bubbles package
type BubblesImpl struct{}

// TextInputImpl implements the TextInput interface wrapping textinput.Model
type TextInputImpl struct {
	model textinput.Model
}

// NewBubbles creates a new Bubbles implementation
func NewBubbles() Bubbles {
	return &BubblesImpl{}
}

func (b *BubblesImpl) NewTextInput() TextInput {
	return &TextInputImpl{model: textinput.New()}
}

func (t *TextInputImpl) SetValue(s string) {
	t.model.SetValue(s)
}

func (t *TextInputImpl) Value() string {
	return t.model.Value()
}

func (t *TextInputImpl) SetPlaceholder(str string) {
	t.model.Placeholder = str
}

func (t *TextInputImpl) Placeholder() string {
	return t.model.Placeholder
}

func (t *TextInputImpl) Focus() tea.Cmd {
	return t.model.Focus()
}

func (t *TextInputImpl) Blur() {
	t.model.Blur()
}

func (t *TextInputImpl) Focused() bool {
	return t.model.Focused()
}

func (t *TextInputImpl) SetPrompt(str string) {
	t.model.Prompt = str
}

func (t *TextInputImpl) Prompt() string {
	return t.model.Prompt
}

func (t *TextInputImpl) SetCharLimit(limit int) {
	t.model.CharLimit = limit
}

func (t *TextInputImpl) CharLimit() int {
	return t.model.CharLimit
}

func (t *TextInputImpl) SetWidth(w int) {
	t.model.Width = w
}

func (t *TextInputImpl) Width() int {
	return t.model.Width
}

func (t *TextInputImpl) Update(msg tea.Msg) (TextInput, tea.Cmd) {
	model, cmd := t.model.Update(msg)
	t.model = model
	return t, cmd
}

func (t *TextInputImpl) View() string {
	return t.model.View()
}