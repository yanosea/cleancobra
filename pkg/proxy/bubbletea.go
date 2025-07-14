package proxy

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Bubbletea interface {
	NewProgram(model Model, opts ...ProgramOption) Program
	Quit() Msg
	WithAltScreen() ProgramOption
}

type bubbleteaProxy struct{}

func NewBubbletea() Bubbletea {
	return &bubbleteaProxy{}
}

func (*bubbleteaProxy) NewProgram(model Model, opts ...ProgramOption) Program {
	teaOpts := make([]tea.ProgramOption, len(opts))
	for i, opt := range opts {
		if proxyOpt, ok := opt.(*programOptionProxy); ok {
			teaOpts[i] = proxyOpt.option
		}
	}

	adapter := &modelAdapter{model: model}
	return &programProxy{Program: tea.NewProgram(adapter, teaOpts...)}
}

func (*bubbleteaProxy) Quit() Msg {
	return tea.Quit()
}

func (*bubbleteaProxy) WithAltScreen() ProgramOption {
	return &programOptionProxy{option: tea.WithAltScreen()}
}

type Program interface {
	Kill()
	Run() (tea.Model, error)
	Send(msg tea.Msg)
	Wait() tea.Model
}

type programProxy struct {
	Program *tea.Program
}

func (p *programProxy) Kill() {
	p.Program.Kill()
}

func (p *programProxy) Run() (tea.Model, error) {
	return p.Program.Run()
}

func (p *programProxy) Send(msg tea.Msg) {
	p.Program.Send(msg)
}

func (p *programProxy) Wait() tea.Model {
	p.Program.Wait()
	return nil
}

type KeyMsg interface {
	String() string
}

type keyMsgProxy struct {
	tea.KeyMsg
}

func (k *keyMsgProxy) String() string {
	return k.KeyMsg.String()
}

func NewKeyMsg(keyMsg tea.KeyMsg) KeyMsg {
	return &keyMsgProxy{KeyMsg: keyMsg}
}

type WindowSizeMsg interface {
	GetHeight() int
	GetWidth() int
}

type windowSizeMsgProxy struct {
	tea.WindowSizeMsg
}

func (w *windowSizeMsgProxy) GetHeight() int {
	return w.Height
}

func (w *windowSizeMsgProxy) GetWidth() int {
	return w.Width
}

func NewWindowSizeMsg(windowSizeMsg tea.WindowSizeMsg) WindowSizeMsg {
	return &windowSizeMsgProxy{WindowSizeMsg: windowSizeMsg}
}

type modelAdapter struct {
	model Model
}

func (m *modelAdapter) Init() tea.Cmd {
	cmd := m.model.Init()
	if cmd == nil {
		return nil
	}
	return func() tea.Msg {
		return cmd()
	}
}

func (m *modelAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := m.model.Update(msg)
	adapter := &modelAdapter{model: updatedModel}
	if cmd == nil {
		return adapter, nil
	}
	return adapter, func() tea.Msg {
		return cmd()
	}
}

func (m *modelAdapter) View() string {
	return m.model.View()
}

func Quit() Cmd {
	return func() Msg {
		return tea.Quit()
	}
}

type Model interface {
	Init() Cmd
	Update(Msg) (Model, Cmd)
	View() string
}

type Cmd func() Msg

type Msg any

type ProgramOption any

type programOptionProxy struct {
	option tea.ProgramOption
}

func IsKeyMsg(msg tea.Msg) (KeyMsg, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return NewKeyMsg(keyMsg), true
	}
	return nil, false
}

func IsWindowSizeMsg(msg tea.Msg) (WindowSizeMsg, bool) {
	if windowSizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		return NewWindowSizeMsg(windowSizeMsg), true
	}
	return nil, false
}
