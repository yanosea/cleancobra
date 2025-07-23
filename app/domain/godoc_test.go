package domain

import (
	"testing"
)

func TestPackageDocumentation(t *testing.T) {
	// This test verifies that the domain package is properly documented
	// and can be imported without issues.

	// Test that we can access package-level constants and types
	errorTypes := []ErrorType{
		ErrorTypeNotFound,
		ErrorTypeInvalidInput,
		ErrorTypeFileSystem,
		ErrorTypeJSON,
		ErrorTypeConfiguration,
	}

	for i, et := range errorTypes {
		if int(et) != i {
			t.Errorf("ErrorType constant %v has unexpected value %d, want %d", et, int(et), i)
		}
	}

	// Test that predefined errors are accessible
	predefinedErrors := []*DomainError{
		ErrTodoNotFound,
		ErrEmptyDescription,
		ErrInvalidID,
	}

	for _, err := range predefinedErrors {
		if err == nil {
			t.Errorf("Predefined error is nil")
		}
		if err.Type < ErrorTypeNotFound || err.Type > ErrorTypeConfiguration {
			t.Errorf("Predefined error has invalid type: %v", err.Type)
		}
		if err.Message == "" {
			t.Errorf("Predefined error has empty message")
		}
	}
}

func TestPackageTypes(t *testing.T) {
	// Test that all main types can be instantiated

	// Test ErrorType
	var et ErrorType = ErrorTypeNotFound
	if et.String() != "NotFound" {
		t.Errorf("ErrorType.String() = %v, want 'NotFound'", et.String())
	}

	// Test DomainError
	domainErr := &DomainError{
		Type:    ErrorTypeInvalidInput,
		Message: "test message",
		Cause:   nil,
	}
	if domainErr.Type != ErrorTypeInvalidInput {
		t.Errorf("DomainError.Type = %v, want ErrorTypeInvalidInput", domainErr.Type)
	}

	// Test Todo (basic instantiation)
	todo := &Todo{
		ID:          1,
		Description: "test",
		Done:        false,
	}
	if todo.ID != 1 {
		t.Errorf("Todo.ID = %v, want 1", todo.ID)
	}
}

func TestPackageInterfaces(t *testing.T) {
	// Test that TodoRepository interface can be referenced
	var repo TodoRepository
	if repo != nil {
		t.Errorf("Uninitialized TodoRepository should be nil")
	}

	// Test that we can assign a mock implementation
	mock := newMockTodoRepository()
	repo = mock
	if repo == nil {
		t.Errorf("TodoRepository assignment failed")
	}

	// Test interface methods are callable
	todos, err := repo.FindAll()
	if err != nil {
		t.Errorf("TodoRepository.FindAll() error = %v, want nil", err)
	}
	if todos == nil {
		t.Errorf("TodoRepository.FindAll() returned nil slice")
	}
}
