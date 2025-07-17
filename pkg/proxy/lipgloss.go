//go:generate mockgen -source=lipgloss.go -destination=lipgloss_mock.go -package=proxy

package proxy

import (
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss provides a proxy interface for lipgloss package functionality
type Lipgloss interface {
	NewStyle() Style
	JoinVertical(pos lipgloss.Position, strs ...string) string
	JoinHorizontal(pos lipgloss.Position, strs ...string) string
	Width(str string) int
	Height(str string) int
	Left() lipgloss.Position
	Center() lipgloss.Position
	Right() lipgloss.Position
	Top() lipgloss.Position
	Bottom() lipgloss.Position
}

// Style provides a proxy interface for lipgloss.Style
type Style interface {
	Render(strs ...string) string
	Width(i int) Style
	Height(i int) Style
	Padding(i ...int) Style
	Margin(i ...int) Style
	Border(border lipgloss.Border, sides ...bool) Style
	BorderStyle(border lipgloss.Border) Style
	BorderTop(v bool) Style
	BorderRight(v bool) Style
	BorderBottom(v bool) Style
	BorderLeft(v bool) Style
	Foreground(c lipgloss.TerminalColor) Style
	Background(c lipgloss.TerminalColor) Style
	Bold(v bool) Style
	Italic(v bool) Style
	Underline(v bool) Style
	Strikethrough(v bool) Style
	Align(p lipgloss.Position) Style
	AlignHorizontal(p lipgloss.Position) Style
	AlignVertical(p lipgloss.Position) Style
}

// LipglossImpl implements the Lipgloss interface using the lipgloss package
type LipglossImpl struct{}

// StyleImpl implements the Style interface wrapping lipgloss.Style
type StyleImpl struct {
	style lipgloss.Style
}

// NewLipgloss creates a new Lipgloss implementation
func NewLipgloss() Lipgloss {
	return &LipglossImpl{}
}

func (l *LipglossImpl) NewStyle() Style {
	return &StyleImpl{style: lipgloss.NewStyle()}
}

func (l *LipglossImpl) JoinVertical(pos lipgloss.Position, strs ...string) string {
	return lipgloss.JoinVertical(pos, strs...)
}

func (l *LipglossImpl) JoinHorizontal(pos lipgloss.Position, strs ...string) string {
	return lipgloss.JoinHorizontal(pos, strs...)
}

func (l *LipglossImpl) Width(str string) int {
	return lipgloss.Width(str)
}

func (l *LipglossImpl) Height(str string) int {
	return lipgloss.Height(str)
}

func (l *LipglossImpl) Left() lipgloss.Position {
	return lipgloss.Left
}

func (l *LipglossImpl) Center() lipgloss.Position {
	return lipgloss.Center
}

func (l *LipglossImpl) Right() lipgloss.Position {
	return lipgloss.Right
}

func (l *LipglossImpl) Top() lipgloss.Position {
	return lipgloss.Top
}

func (l *LipglossImpl) Bottom() lipgloss.Position {
	return lipgloss.Bottom
}

func (s *StyleImpl) Render(strs ...string) string {
	return s.style.Render(strs...)
}

func (s *StyleImpl) Width(i int) Style {
	return &StyleImpl{style: s.style.Width(i)}
}

func (s *StyleImpl) Height(i int) Style {
	return &StyleImpl{style: s.style.Height(i)}
}

func (s *StyleImpl) Padding(i ...int) Style {
	return &StyleImpl{style: s.style.Padding(i...)}
}

func (s *StyleImpl) Margin(i ...int) Style {
	return &StyleImpl{style: s.style.Margin(i...)}
}

func (s *StyleImpl) Border(border lipgloss.Border, sides ...bool) Style {
	return &StyleImpl{style: s.style.Border(border, sides...)}
}

func (s *StyleImpl) BorderStyle(border lipgloss.Border) Style {
	return &StyleImpl{style: s.style.BorderStyle(border)}
}

func (s *StyleImpl) BorderTop(v bool) Style {
	return &StyleImpl{style: s.style.BorderTop(v)}
}

func (s *StyleImpl) BorderRight(v bool) Style {
	return &StyleImpl{style: s.style.BorderRight(v)}
}

func (s *StyleImpl) BorderBottom(v bool) Style {
	return &StyleImpl{style: s.style.BorderBottom(v)}
}

func (s *StyleImpl) BorderLeft(v bool) Style {
	return &StyleImpl{style: s.style.BorderLeft(v)}
}

func (s *StyleImpl) Foreground(c lipgloss.TerminalColor) Style {
	return &StyleImpl{style: s.style.Foreground(c)}
}

func (s *StyleImpl) Background(c lipgloss.TerminalColor) Style {
	return &StyleImpl{style: s.style.Background(c)}
}

func (s *StyleImpl) Bold(v bool) Style {
	return &StyleImpl{style: s.style.Bold(v)}
}

func (s *StyleImpl) Italic(v bool) Style {
	return &StyleImpl{style: s.style.Italic(v)}
}

func (s *StyleImpl) Underline(v bool) Style {
	return &StyleImpl{style: s.style.Underline(v)}
}

func (s *StyleImpl) Strikethrough(v bool) Style {
	return &StyleImpl{style: s.style.Strikethrough(v)}
}

func (s *StyleImpl) Align(p lipgloss.Position) Style {
	return &StyleImpl{style: s.style.Align(p)}
}

func (s *StyleImpl) AlignHorizontal(p lipgloss.Position) Style {
	return &StyleImpl{style: s.style.AlignHorizontal(p)}
}

func (s *StyleImpl) AlignVertical(p lipgloss.Position) Style {
	return &StyleImpl{style: s.style.AlignVertical(p)}
}