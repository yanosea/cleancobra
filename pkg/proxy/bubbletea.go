//go:generate mockgen -source=bubbletea.go -destination=bubbletea_mock.go -package=proxy

package proxy

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Bubbletea provides a proxy interface for bubbletea package functionality
type Bubbletea interface {
	NewProgram(model tea.Model, opts ...tea.ProgramOption) Program
	WithAltScreen() tea.ProgramOption
	WithMouseCellMotion() tea.ProgramOption
}

// Program provides a proxy interface for tea.Program
type Program interface {
	Start() (tea.Model, error)
	Run() (tea.Model, error)
	Send(msg tea.Msg)
	Quit()
}

// BubbleteaImpl implements the Bubbletea interface using the bubbletea package
type BubbleteaImpl struct{}

// ProgramImpl implements the Program interface wrapping tea.Program
type ProgramImpl struct {
	program *tea.Program
}

// NewBubbletea creates a new Bubbletea implementation
func NewBubbletea() Bubbletea {
	return &BubbleteaImpl{}
}

func (b *BubbleteaImpl) NewProgram(model tea.Model, opts ...tea.ProgramOption) Program {
	return &ProgramImpl{program: tea.NewProgram(model, opts...)}
}

func (b *BubbleteaImpl) WithAltScreen() tea.ProgramOption {
	return tea.WithAltScreen()
}

func (b *BubbleteaImpl) WithMouseCellMotion() tea.ProgramOption {
	return tea.WithMouseCellMotion()
}

func (p *ProgramImpl) Start() (tea.Model, error) {
	model, err := p.program.Run()
	return model, err
}

func (p *ProgramImpl) Run() (tea.Model, error) {
	return p.program.Run()
}

func (p *ProgramImpl) Send(msg tea.Msg) {
	p.program.Send(msg)
}

func (p *ProgramImpl) Quit() {
	p.program.Quit()
}