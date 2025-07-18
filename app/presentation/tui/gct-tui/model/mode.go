package model

// Mode represents the current mode of the TUI application
type Mode int

const (
	// ModeNormal is the default navigation mode
	ModeNormal Mode = iota
	// ModeInput is the mode for adding new todos
	ModeInput
	// ModeConfirmation is the mode for confirming actions like deletion
	ModeConfirmation
	// ModeEdit is the mode for editing existing todos
	ModeEdit
)

// String returns the string representation of the mode
func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "Normal"
	case ModeInput:
		return "Input"
	case ModeConfirmation:
		return "Confirmation"
	case ModeEdit:
		return "Edit"
	default:
		return "Unknown"
	}
}
