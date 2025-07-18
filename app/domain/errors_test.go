package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		name     string
		errorType ErrorType
		want     string
	}{
		{
			name:     "positive testing - ErrorTypeNotFound",
			errorType: ErrorTypeNotFound,
			want:     "NotFound",
		},
		{
			name:     "positive testing - ErrorTypeInvalidInput",
			errorType: ErrorTypeInvalidInput,
			want:     "InvalidInput",
		},
		{
			name:     "positive testing - ErrorTypeFileSystem",
			errorType: ErrorTypeFileSystem,
			want:     "FileSystem",
		},
		{
			name:     "positive testing - ErrorTypeJSON",
			errorType: ErrorTypeJSON,
			want:     "JSON",
		},
		{
			name:     "positive testing - ErrorTypeConfiguration",
			errorType: ErrorTypeConfiguration,
			want:     "Configuration",
		},
		{
			name:     "positive testing - unknown error type",
			errorType: ErrorType(999),
			want:     "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.errorType.String()
			if got != tt.want {
				t.Errorf("ErrorType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDomainError(t *testing.T) {
	causeErr := errors.New("underlying error")

	tests := []struct {
		name      string
		errorType ErrorType
		message   string
		cause     error
		want      *DomainError
	}{
		{
			name:      "positive testing - with cause error",
			errorType: ErrorTypeNotFound,
			message:   "test message",
			cause:     causeErr,
			want: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "test message",
				Cause:   causeErr,
			},
		},
		{
			name:      "positive testing - without cause error",
			errorType: ErrorTypeInvalidInput,
			message:   "validation failed",
			cause:     nil,
			want: &DomainError{
				Type:    ErrorTypeInvalidInput,
				Message: "validation failed",
				Cause:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDomainError(tt.errorType, tt.message, tt.cause)
			if got.Type != tt.want.Type {
				t.Errorf("NewDomainError().Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.Message != tt.want.Message {
				t.Errorf("NewDomainError().Message = %v, want %v", got.Message, tt.want.Message)
			}
			if got.Cause != tt.want.Cause {
				t.Errorf("NewDomainError().Cause = %v, want %v", got.Cause, tt.want.Cause)
			}
		})
	}
}

func TestDomainError_Error(t *testing.T) {
	causeErr := errors.New("underlying error")

	tests := []struct {
		name        string
		domainError *DomainError
		want        string
	}{
		{
			name: "positive testing - with cause error",
			domainError: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "todo not found",
				Cause:   causeErr,
			},
			want: "NotFound: todo not found (caused by: underlying error)",
		},
		{
			name: "positive testing - without cause error",
			domainError: &DomainError{
				Type:    ErrorTypeInvalidInput,
				Message: "invalid input",
				Cause:   nil,
			},
			want: "InvalidInput: invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.domainError.Error()
			if got != tt.want {
				t.Errorf("DomainError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	causeErr := errors.New("underlying error")

	tests := []struct {
		name        string
		domainError *DomainError
		want        error
	}{
		{
			name: "positive testing - with cause error",
			domainError: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "test message",
				Cause:   causeErr,
			},
			want: causeErr,
		},
		{
			name: "positive testing - without cause error",
			domainError: &DomainError{
				Type:    ErrorTypeInvalidInput,
				Message: "test message",
				Cause:   nil,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.domainError.Unwrap()
			if got != tt.want {
				t.Errorf("DomainError.Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainError_Is(t *testing.T) {
	causeErr := errors.New("underlying error")
	domainErr1 := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "todo not found",
		Cause:   causeErr,
	}
	domainErr2 := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "todo not found",
		Cause:   nil,
	}
	domainErr3 := &DomainError{
		Type:    ErrorTypeInvalidInput,
		Message: "invalid input",
		Cause:   nil,
	}
	regularErr := errors.New("regular error")

	tests := []struct {
		name        string
		domainError *DomainError
		target      error
		want        bool
	}{
		{
			name:        "positive testing - same type and message",
			domainError: domainErr1,
			target:      domainErr2,
			want:        true,
		},
		{
			name:        "positive testing - different type",
			domainError: domainErr1,
			target:      domainErr3,
			want:        false,
		},
		{
			name:        "positive testing - different message",
			domainError: domainErr1,
			target: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "different message",
				Cause:   nil,
			},
			want: false,
		},
		{
			name:        "positive testing - non-domain error",
			domainError: domainErr1,
			target:      regularErr,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.domainError.Is(tt.target)
			if got != tt.want {
				t.Errorf("DomainError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name      string
		err       *DomainError
		wantType  ErrorType
		wantMsg   string
		wantCause error
	}{
		{
			name:      "positive testing - ErrTodoNotFound",
			err:       ErrTodoNotFound,
			wantType:  ErrorTypeNotFound,
			wantMsg:   "todo not found",
			wantCause: nil,
		},
		{
			name:      "positive testing - ErrEmptyDescription",
			err:       ErrEmptyDescription,
			wantType:  ErrorTypeInvalidInput,
			wantMsg:   "description cannot be empty",
			wantCause: nil,
		},
		{
			name:      "positive testing - ErrInvalidID",
			err:       ErrInvalidID,
			wantType:  ErrorTypeInvalidInput,
			wantMsg:   "invalid todo ID",
			wantCause: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Type != tt.wantType {
				t.Errorf("Error type = %v, want %v", tt.err.Type, tt.wantType)
			}
			if tt.err.Message != tt.wantMsg {
				t.Errorf("Error message = %v, want %v", tt.err.Message, tt.wantMsg)
			}
			if tt.err.Cause != tt.wantCause {
				t.Errorf("Error cause = %v, want %v", tt.err.Cause, tt.wantCause)
			}
		})
	}
}

func TestIsDomainError(t *testing.T) {
	domainErr := NewDomainError(ErrorTypeNotFound, "test", nil)
	regularErr := errors.New("regular error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing - domain error",
			err:  domainErr,
			want: true,
		},
		{
			name: "positive testing - regular error",
			err:  regularErr,
			want: false,
		},
		{
			name: "positive testing - nil error",
			err:  nil,
			want: false,
		},
		{
			name: "positive testing - wrapped domain error",
			err:  fmt.Errorf("wrapped: %w", domainErr),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDomainError(tt.err)
			if got != tt.want {
				t.Errorf("IsDomainError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetErrorType(t *testing.T) {
	domainErr := NewDomainError(ErrorTypeNotFound, "test", nil)
	regularErr := errors.New("regular error")

	tests := []struct {
		name string
		err  error
		want ErrorType
	}{
		{
			name: "positive testing - domain error",
			err:  domainErr,
			want: ErrorTypeNotFound,
		},
		{
			name: "positive testing - regular error",
			err:  regularErr,
			want: ErrorTypeConfiguration,
		},
		{
			name: "positive testing - wrapped domain error",
			err:  fmt.Errorf("wrapped: %w", domainErr),
			want: ErrorTypeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetErrorType(tt.err)
			if got != tt.want {
				t.Errorf("GetErrorType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	notFoundErr := NewDomainError(ErrorTypeNotFound, "test", nil)
	invalidInputErr := NewDomainError(ErrorTypeInvalidInput, "test", nil)
	regularErr := errors.New("regular error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing - not found error",
			err:  notFoundErr,
			want: true,
		},
		{
			name: "positive testing - invalid input error",
			err:  invalidInputErr,
			want: false,
		},
		{
			name: "positive testing - regular error",
			err:  regularErr,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotFoundError(tt.err)
			if got != tt.want {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalidInputError(t *testing.T) {
	invalidInputErr := NewDomainError(ErrorTypeInvalidInput, "test", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "test", nil)
	regularErr := errors.New("regular error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing - invalid input error",
			err:  invalidInputErr,
			want: true,
		},
		{
			name: "positive testing - not found error",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "positive testing - regular error",
			err:  regularErr,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInvalidInputError(tt.err)
			if got != tt.want {
				t.Errorf("IsInvalidInputError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFileSystemError(t *testing.T) {
	fileSystemErr := NewDomainError(ErrorTypeFileSystem, "test", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "test", nil)
	regularErr := errors.New("regular error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing - file system error",
			err:  fileSystemErr,
			want: true,
		},
		{
			name: "positive testing - not found error",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "positive testing - regular error",
			err:  regularErr,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFileSystemError(tt.err)
			if got != tt.want {
				t.Errorf("IsFileSystemError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJSONError(t *testing.T) {
	jsonErr := NewDomainError(ErrorTypeJSON, "test", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "test", nil)
	regularErr := errors.New("regular error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing - JSON error",
			err:  jsonErr,
			want: true,
		},
		{
			name: "positive testing - not found error",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "positive testing - regular error",
			err:  regularErr,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsJSONError(tt.err)
			if got != tt.want {
				t.Errorf("IsJSONError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsConfigurationError(t *testing.T) {
	configErr := NewDomainError(ErrorTypeConfiguration, "test", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "test", nil)
	regularErr := errors.New("regular error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing - configuration error",
			err:  configErr,
			want: true,
		},
		{
			name: "positive testing - not found error",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "positive testing - regular error (defaults to configuration)",
			err:  regularErr,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsConfigurationError(tt.err)
			if got != tt.want {
				t.Errorf("IsConfigurationError() = %v, want %v", got, tt.want)
			}
		})
	}
}