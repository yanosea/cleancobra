//go:generate mockgen -source=io.go -destination=io_mock.go -package=proxy

package proxy

import (
	"io"
)

// IO provides a proxy interface for io package functions
type IO interface {
	ReadAll(r io.Reader) ([]byte, error)
	WriteString(w io.Writer, s string) (n int, err error)
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

// IOImpl implements the IO interface using the standard io package
type IOImpl struct{}

// NewIO creates a new IO implementation
func NewIO() IO {
	return &IOImpl{}
}

func (i *IOImpl) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func (i *IOImpl) WriteString(w io.Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}

func (i *IOImpl) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}