//go:generate mockgen -source=fmt.go -destination=fmt_mock.go -package=proxy

package proxy

import (
	"fmt"
	"io"
)

// Fmt provides a proxy interface for fmt package functions
type Fmt interface {
	Sprintf(format string, a ...any) string
	Printf(format string, a ...any) (n int, err error)
	Println(a ...any) (n int, err error)
	Print(a ...any) (n int, err error)
	Fprintf(w io.Writer, format string, a ...any) (n int, err error)
	Fprintln(w io.Writer, a ...any) (n int, err error)
	Fprint(w io.Writer, a ...any) (n int, err error)
	Errorf(format string, a ...any) error
}

// FmtImpl implements the Fmt interface using the standard fmt package
type FmtImpl struct{}

// NewFmt creates a new Fmt implementation
func NewFmt() Fmt {
	return &FmtImpl{}
}

func (f *FmtImpl) Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func (f *FmtImpl) Printf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}

func (f *FmtImpl) Println(a ...any) (n int, err error) {
	return fmt.Println(a...)
}

func (f *FmtImpl) Print(a ...any) (n int, err error) {
	return fmt.Print(a...)
}

func (f *FmtImpl) Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return fmt.Fprintf(w, format, a...)
}

func (f *FmtImpl) Fprintln(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprintln(w, a...)
}

func (f *FmtImpl) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, a...)
}

func (f *FmtImpl) Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}
