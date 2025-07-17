package model

import (
	"github.com/yanosea/gct/pkg/proxy"
)

// InputState manages text input state for the TUI
type InputState struct {
	textInput proxy.TextInput
	focused   bool
}

// NewInputState creates a new input state with the given bubbles proxy
func NewInputState(bubbles proxy.Bubbles) *InputState {
	input := bubbles.NewTextInput()
	input.SetPlaceholder("Enter todo description...")
	input.SetCharLimit(500)
	input.SetWidth(50)
	
	return &InputState{
		textInput: input,
		focused:   false,
	}
}

// TextInput returns the underlying text input model
func (i *InputState) TextInput() proxy.TextInput {
	return i.textInput
}

// Value returns the current input value
func (i *InputState) Value() string {
	return i.textInput.Value()
}

// SetValue sets the input value
func (i *InputState) SetValue(value string) {
	i.textInput.SetValue(value)
}

// Focus focuses the input
func (i *InputState) Focus() {
	i.textInput.Focus()
	i.focused = true
}

// Blur blurs the input
func (i *InputState) Blur() {
	i.textInput.Blur()
	i.focused = false
}

// IsFocused returns whether the input is focused
func (i *InputState) IsFocused() bool {
	return i.focused
}

// Clear clears the input value
func (i *InputState) Clear() {
	i.textInput.SetValue("")
}

// SetWidth sets the input width
func (i *InputState) SetWidth(width int) {
	i.textInput.SetWidth(width)
}

// SetPlaceholder sets the input placeholder
func (i *InputState) SetPlaceholder(placeholder string) {
	i.textInput.SetPlaceholder(placeholder)
}

// View returns the input view
func (i *InputState) View() string {
	return i.textInput.View()
}