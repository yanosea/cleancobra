package domain

import (
	"errors"
	"testing"

	"github.com/yanosea/gct/pkg/proxy"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		name      string
		errorType ErrorType
		expected  string
	}{
		{
			name:      "ErrorTypeNotFound",
			errorType: ErrorTypeNotFound,
			expected:  "NotFound",
		},
		{
			name:      "ErrorTypeInvalidInput",
			errorType: ErrorTypeInvalidInput,
			expected:  "InvalidInput",
		},
		{
			name:      "ErrorTypeFileSystem",
			errorType: ErrorTypeFileSystem,
			expected:  "FileSystem",
		},
		{
			name:      "ErrorTypeJSON",
			errorType: ErrorTypeJSON,
			expected:  "JSON",
		},
		{
			name:      "ErrorTypeConfiguration",
			errorType: ErrorTypeConfiguration,
			expected:  "Configuration",
		},
		{
			name:      "Unknown error type",
			errorType: ErrorType(999),
			expected:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errorType.String()
			if result != tt.expected {
				t.Errorf("ErrorType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewDomainError(t *testing.T) {
	tests := []struct {
		name      string
		errorType ErrorType
		message   string
		cause     error
	}{
		{
			name:      "Error without cause",
			errorType: ErrorTypeNotFound,
			message:   "test message",
			cause:     nil,
		},
		{
			name:      "Error with cause",
			errorType: ErrorTypeInvalidInput,
			message:   "test message with cause",
			cause:     errors.New("underlying error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewDomainError(tt.errorType, tt.message, tt.cause)

			if err.Type != tt.errorType {
				t.Errorf("NewDomainError().Type = %v, want %v", err.Type, tt.errorType)
			}
			if err.Message != tt.message {
				t.Errorf("NewDomainError().Message = %v, want %v", err.Message, tt.message)
			}
			if err.Cause != tt.cause {
				t.Errorf("NewDomainError().Cause = %v, want %v", err.Cause, tt.cause)
			}
		})
	}
}

func TestDomainError_Error(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	tests := []struct {
		name     string
		err      *DomainError
		expected string
	}{
		{
			name: "Error without cause",
			err: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "test message",
				Cause:   nil,
			},
			expected: "NotFound: test message",
		},
		{
			name: "Error with cause",
			err: &DomainError{
				Type:    ErrorTypeInvalidInput,
				Message: "test message",
				Cause:   errors.New("underlying error"),
			},
			expected: "InvalidInput: test message (caused by: underlying error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("DomainError.Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "test message",
		Cause:   cause,
	}

	result := err.Unwrap()
	if result != cause {
		t.Errorf("DomainError.Unwrap() = %v, want %v", result, cause)
	}

	// Test with nil cause
	errNoCause := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "test message",
		Cause:   nil,
	}

	result = errNoCause.Unwrap()
	if result != nil {
		t.Errorf("DomainError.Unwrap() = %v, want nil", result)
	}
}

func TestDomainError_Is(t *testing.T) {
	err1 := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "test message",
		Cause:   nil,
	}

	err2 := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "test message",
		Cause:   nil,
	}

	err3 := &DomainError{
		Type:    ErrorTypeInvalidInput,
		Message: "test message",
		Cause:   nil,
	}

	err4 := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "different message",
		Cause:   nil,
	}

	nonDomainErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      *DomainError
		target   error
		expected bool
	}{
		{
			name:     "Same type and message",
			err:      err1,
			target:   err2,
			expected: true,
		},
		{
			name:     "Different type",
			err:      err1,
			target:   err3,
			expected: false,
		},
		{
			name:     "Different message",
			err:      err1,
			target:   err4,
			expected: false,
		},
		{
			name:     "Non-domain error",
			err:      err1,
			target:   nonDomainErr,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Is(tt.target)
			if result != tt.expected {
				t.Errorf("DomainError.Is() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsDomainError(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	domainErr := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "test message",
		Cause:   nil,
	}

	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Domain error",
			err:      domainErr,
			expected: true,
		},
		{
			name:     "Regular error",
			err:      regularErr,
			expected: false,
		},
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDomainError(tt.err)
			if result != tt.expected {
				t.Errorf("IsDomainError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetErrorType(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	domainErr := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "test message",
		Cause:   nil,
	}

	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		expected ErrorType
	}{
		{
			name:     "Domain error",
			err:      domainErr,
			expected: ErrorTypeNotFound,
		},
		{
			name:     "Regular error",
			err:      regularErr,
			expected: ErrorTypeConfiguration,
		},
		{
			name:     "Nil error",
			err:      nil,
			expected: ErrorTypeConfiguration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorType(tt.err)
			if result != tt.expected {
				t.Errorf("GetErrorType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	notFoundErr := &DomainError{Type: ErrorTypeNotFound, Message: "not found", Cause: nil}
	invalidInputErr := &DomainError{Type: ErrorTypeInvalidInput, Message: "invalid", Cause: nil}
	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Not found error",
			err:      notFoundErr,
			expected: true,
		},
		{
			name:     "Invalid input error",
			err:      invalidInputErr,
			expected: false,
		},
		{
			name:     "Regular error",
			err:      regularErr,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFoundError(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotFoundError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsInvalidInputError(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	invalidInputErr := &DomainError{Type: ErrorTypeInvalidInput, Message: "invalid", Cause: nil}
	notFoundErr := &DomainError{Type: ErrorTypeNotFound, Message: "not found", Cause: nil}
	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Invalid input error",
			err:      invalidInputErr,
			expected: true,
		},
		{
			name:     "Not found error",
			err:      notFoundErr,
			expected: false,
		},
		{
			name:     "Regular error",
			err:      regularErr,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInvalidInputError(tt.err)
			if result != tt.expected {
				t.Errorf("IsInvalidInputError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsFileSystemError(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	fileSystemErr := &DomainError{Type: ErrorTypeFileSystem, Message: "file error", Cause: nil}
	notFoundErr := &DomainError{Type: ErrorTypeNotFound, Message: "not found", Cause: nil}
	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "File system error",
			err:      fileSystemErr,
			expected: true,
		},
		{
			name:     "Not found error",
			err:      notFoundErr,
			expected: false,
		},
		{
			name:     "Regular error",
			err:      regularErr,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFileSystemError(tt.err)
			if result != tt.expected {
				t.Errorf("IsFileSystemError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsJSONError(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	jsonErr := &DomainError{Type: ErrorTypeJSON, Message: "json error", Cause: nil}
	notFoundErr := &DomainError{Type: ErrorTypeNotFound, Message: "not found", Cause: nil}
	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "JSON error",
			err:      jsonErr,
			expected: true,
		},
		{
			name:     "Not found error",
			err:      notFoundErr,
			expected: false,
		},
		{
			name:     "Regular error",
			err:      regularErr,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsJSONError(tt.err)
			if result != tt.expected {
				t.Errorf("IsJSONError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsConfigurationError(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	configErr := &DomainError{Type: ErrorTypeConfiguration, Message: "config error", Cause: nil}
	notFoundErr := &DomainError{Type: ErrorTypeNotFound, Message: "not found", Cause: nil}
	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Configuration error",
			err:      configErr,
			expected: true,
		},
		{
			name:     "Not found error",
			err:      notFoundErr,
			expected: false,
		},
		{
			name:     "Regular error (treated as configuration)",
			err:      regularErr,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsConfigurationError(tt.err)
			if result != tt.expected {
				t.Errorf("IsConfigurationError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInitializeDomainErrors(t *testing.T) {
	errorsProxy := proxy.NewErrors()
	fmtProxy := proxy.NewFmt()

	InitializeDomainErrors(errorsProxy, fmtProxy)

	// Test that the proxies are set by using them
	err := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "test",
		Cause:   nil,
	}

	// This should not panic if proxies are properly initialized
	result := err.Error()
	expected := "NotFound: test"
	if result != expected {
		t.Errorf("After InitializeDomainErrors, Error() = %v, want %v", result, expected)
	}
}

func TestPredefinedErrors(t *testing.T) {
	// Initialize with real proxies for testing
	InitializeDomainErrors(proxy.NewErrors(), proxy.NewFmt())

	tests := []struct {
		name    string
		err     *DomainError
		errType ErrorType
		message string
	}{
		{
			name:    "ErrTodoNotFound",
			err:     ErrTodoNotFound,
			errType: ErrorTypeNotFound,
			message: "todo not found",
		},
		{
			name:    "ErrEmptyDescription",
			err:     ErrEmptyDescription,
			errType: ErrorTypeInvalidInput,
			message: "description cannot be empty",
		},
		{
			name:    "ErrInvalidID",
			err:     ErrInvalidID,
			errType: ErrorTypeInvalidInput,
			message: "invalid todo ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Type != tt.errType {
				t.Errorf("%s.Type = %v, want %v", tt.name, tt.err.Type, tt.errType)
			}
			if tt.err.Message != tt.message {
				t.Errorf("%s.Message = %v, want %v", tt.name, tt.err.Message, tt.message)
			}
			if tt.err.Cause != nil {
				t.Errorf("%s.Cause = %v, want nil", tt.name, tt.err.Cause)
			}
		})
	}
}
