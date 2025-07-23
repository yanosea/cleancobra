//go:generate mockgen -source=errors.go -destination=errors_mock.go -package=proxy

package proxy

import (
	"errors"
)

// Errors provides a proxy interface for errors package functions
type Errors interface {
	New(text string) error
	Is(err, target error) bool
	As(err error, target any) bool
	Unwrap(err error) error
}

// ErrorsImpl implements the Errors interface using the standard errors package
type ErrorsImpl struct{}

// NewErrors creates a new Errors implementation
func NewErrors() Errors {
	return &ErrorsImpl{}
}

func (e *ErrorsImpl) New(text string) error {
	return errors.New(text)
}

func (e *ErrorsImpl) Is(err, target error) bool {
	return errors.Is(err, target)
}

func (e *ErrorsImpl) As(err error, target any) bool {
	return errors.As(err, target)
}

func (e *ErrorsImpl) Unwrap(err error) error {
	return errors.Unwrap(err)
}
