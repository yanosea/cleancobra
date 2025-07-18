package domain

import (
	"errors"
	"fmt"
)

// ErrorType represents the type of domain error
type ErrorType int

const (
	// ErrorTypeNotFound indicates a resource was not found
	ErrorTypeNotFound ErrorType = iota
	// ErrorTypeInvalidInput indicates invalid input was provided
	ErrorTypeInvalidInput
	// ErrorTypeFileSystem indicates a file system operation failed
	ErrorTypeFileSystem
	// ErrorTypeJSON indicates a JSON operation failed
	ErrorTypeJSON
	// ErrorTypeConfiguration indicates a configuration error
	ErrorTypeConfiguration
)

// String returns the string representation of the error type
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeNotFound:
		return "NotFound"
	case ErrorTypeInvalidInput:
		return "InvalidInput"
	case ErrorTypeFileSystem:
		return "FileSystem"
	case ErrorTypeJSON:
		return "JSON"
	case ErrorTypeConfiguration:
		return "Configuration"
	default:
		return "Unknown"
	}
}

// DomainError represents a domain-specific error
type DomainError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// NewDomainError creates a new domain error
func NewDomainError(errorType ErrorType, message string, cause error) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type.String(), e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

// Unwrap returns the underlying cause error
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target error
func (e *DomainError) Is(target error) bool {
	if t, ok := target.(*DomainError); ok {
		return e.Type == t.Type && e.Message == t.Message
	}
	return false
}

// Predefined domain errors
var (
	// ErrTodoNotFound indicates a todo was not found
	ErrTodoNotFound = NewDomainError(ErrorTypeNotFound, "todo not found", nil)

	// ErrEmptyDescription indicates an empty description was provided
	ErrEmptyDescription = NewDomainError(ErrorTypeInvalidInput, "description cannot be empty", nil)

	// ErrInvalidID indicates an invalid ID was provided
	ErrInvalidID = NewDomainError(ErrorTypeInvalidInput, "invalid todo ID", nil)
)

// IsDomainError checks if an error is a domain error
func IsDomainError(err error) bool {
	var domainErr *DomainError
	return errors.As(err, &domainErr)
}

// GetErrorType returns the error type of a domain error, or ErrorTypeConfiguration for non-domain errors
func GetErrorType(err error) ErrorType {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type
	}
	return ErrorTypeConfiguration
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	return GetErrorType(err) == ErrorTypeNotFound
}

// IsInvalidInputError checks if an error is an invalid input error
func IsInvalidInputError(err error) bool {
	return GetErrorType(err) == ErrorTypeInvalidInput
}

// IsFileSystemError checks if an error is a file system error
func IsFileSystemError(err error) bool {
	return GetErrorType(err) == ErrorTypeFileSystem
}

// IsJSONError checks if an error is a JSON error
func IsJSONError(err error) bool {
	return GetErrorType(err) == ErrorTypeJSON
}

// IsConfigurationError checks if an error is a configuration error
func IsConfigurationError(err error) bool {
	return GetErrorType(err) == ErrorTypeConfiguration
}
