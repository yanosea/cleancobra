package domain

import (
	"errors"
	"testing"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		name string
		et   ErrorType
		want string
	}{
		{
			name: "positive testing (NotFound)",
			et:   ErrorTypeNotFound,
			want: "NotFound",
		},
		{
			name: "positive testing (InvalidInput)",
			et:   ErrorTypeInvalidInput,
			want: "InvalidInput",
		},
		{
			name: "positive testing (FileSystem)",
			et:   ErrorTypeFileSystem,
			want: "FileSystem",
		},
		{
			name: "positive testing (JSON)",
			et:   ErrorTypeJSON,
			want: "JSON",
		},
		{
			name: "positive testing (Configuration)",
			et:   ErrorTypeConfiguration,
			want: "Configuration",
		},
		{
			name: "positive testing (Unknown)",
			et:   ErrorType(999),
			want: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.et.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDomainError(t *testing.T) {
	cause := errors.New("underlying error")
	
	tests := []struct {
		name      string
		errorType ErrorType
		message   string
		cause     error
		want      *DomainError
	}{
		{
			name:      "positive testing (with cause)",
			errorType: ErrorTypeNotFound,
			message:   "resource not found",
			cause:     cause,
			want: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "resource not found",
				Cause:   cause,
			},
		},
		{
			name:      "positive testing (without cause)",
			errorType: ErrorTypeInvalidInput,
			message:   "invalid input provided",
			cause:     nil,
			want: &DomainError{
				Type:    ErrorTypeInvalidInput,
				Message: "invalid input provided",
				Cause:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDomainError(tt.errorType, tt.message, tt.cause)
			if got.Type != tt.want.Type {
				t.Errorf("NewDomainError() Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.Message != tt.want.Message {
				t.Errorf("NewDomainError() Message = %v, want %v", got.Message, tt.want.Message)
			}
			if got.Cause != tt.want.Cause {
				t.Errorf("NewDomainError() Cause = %v, want %v", got.Cause, tt.want.Cause)
			}
		})
	}
}

func TestDomainError_Error(t *testing.T) {
	cause := errors.New("underlying error")
	
	tests := []struct {
		name string
		err  *DomainError
		want string
	}{
		{
			name: "positive testing (with cause)",
			err: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "resource not found",
				Cause:   cause,
			},
			want: "NotFound: resource not found (caused by: underlying error)",
		},
		{
			name: "positive testing (without cause)",
			err: &DomainError{
				Type:    ErrorTypeInvalidInput,
				Message: "invalid input provided",
				Cause:   nil,
			},
			want: "InvalidInput: invalid input provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	
	tests := []struct {
		name string
		err  *DomainError
		want error
	}{
		{
			name: "positive testing (with cause)",
			err: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "resource not found",
				Cause:   cause,
			},
			want: cause,
		},
		{
			name: "positive testing (without cause)",
			err: &DomainError{
				Type:    ErrorTypeInvalidInput,
				Message: "invalid input provided",
				Cause:   nil,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Unwrap()
			if got != tt.want {
				t.Errorf("Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainError_Is(t *testing.T) {
	err1 := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	err2 := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	err3 := NewDomainError(ErrorTypeInvalidInput, "invalid input", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name   string
		err    *DomainError
		target error
		want   bool
	}{
		{
			name:   "positive testing (same error)",
			err:    err1,
			target: err2,
			want:   true,
		},
		{
			name:   "negative testing (different type failed)",
			err:    err1,
			target: err3,
			want:   false,
		},
		{
			name:   "negative testing (non-domain error failed)",
			err:    err1,
			target: otherErr,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Is(tt.target)
			if got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDomainError(t *testing.T) {
	domainErr := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing (domain error)",
			err:  domainErr,
			want: true,
		},
		{
			name: "negative testing (non-domain error failed)",
			err:  otherErr,
			want: false,
		},
		{
			name: "negative testing (nil error failed)",
			err:  nil,
			want: false,
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
	domainErr := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name string
		err  error
		want ErrorType
	}{
		{
			name: "positive testing (domain error)",
			err:  domainErr,
			want: ErrorTypeNotFound,
		},
		{
			name: "positive testing (non-domain error returns Configuration)",
			err:  otherErr,
			want: ErrorTypeConfiguration,
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
	notFoundErr := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	invalidInputErr := NewDomainError(ErrorTypeInvalidInput, "invalid input", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing (not found error)",
			err:  notFoundErr,
			want: true,
		},
		{
			name: "negative testing (different domain error failed)",
			err:  invalidInputErr,
			want: false,
		},
		{
			name: "negative testing (non-domain error failed)",
			err:  otherErr,
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
	invalidInputErr := NewDomainError(ErrorTypeInvalidInput, "invalid input", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing (invalid input error)",
			err:  invalidInputErr,
			want: true,
		},
		{
			name: "negative testing (different domain error failed)",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "negative testing (non-domain error failed)",
			err:  otherErr,
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
	fileSystemErr := NewDomainError(ErrorTypeFileSystem, "file system error", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing (file system error)",
			err:  fileSystemErr,
			want: true,
		},
		{
			name: "negative testing (different domain error failed)",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "negative testing (non-domain error failed)",
			err:  otherErr,
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
	jsonErr := NewDomainError(ErrorTypeJSON, "JSON error", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing (JSON error)",
			err:  jsonErr,
			want: true,
		},
		{
			name: "negative testing (different domain error failed)",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "negative testing (non-domain error failed)",
			err:  otherErr,
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
	configErr := NewDomainError(ErrorTypeConfiguration, "configuration error", nil)
	notFoundErr := NewDomainError(ErrorTypeNotFound, "resource not found", nil)
	otherErr := errors.New("other error")
	
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "positive testing (configuration error)",
			err:  configErr,
			want: true,
		},
		{
			name: "negative testing (different domain error failed)",
			err:  notFoundErr,
			want: false,
		},
		{
			name: "positive testing (non-domain error returns true)",
			err:  otherErr,
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