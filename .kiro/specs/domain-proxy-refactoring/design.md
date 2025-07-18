# Design Document

## Overview

This design outlines the refactoring of the `app/domain` package to use proxy interfaces instead of direct standard library imports. The refactoring will improve testability while maintaining backward compatibility and following the established proxy pattern used throughout the GCT application.

## Architecture

### Current Architecture
```
domain/
├── todo.go (direct imports: time, strings, fmt, json)
├── errors.go (direct imports: fmt, errors)
├── repository.go
└── Direct standard library calls throughout
```

### Target Architecture
```
domain/
├── todo.go (proxy imports: pkg/proxy.Time, pkg/proxy.Strings, pkg/proxy.Fmt, pkg/proxy.JSON)
├── errors.go (proxy imports: pkg/proxy.Fmt)
├── repository.go (unchanged)
└── Dependency injection for standard library operations
```

## Components and Interfaces

### Todo Entity Refactoring

#### Current Structure
```go
type Todo struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Refactored Structure
The `Todo` struct will remain unchanged to maintain backward compatibility. New methods with dependency injection will be added:

```go
// Existing methods remain unchanged for backward compatibility
func NewTodo(id int, description string) (*Todo, error)
func (t *Todo) Toggle()
func (t *Todo) UpdateDescription(description string) error

// New methods with dependency injection
func NewTodoWithDeps(id int, description string, timeProxy proxy.Time, stringsProxy proxy.Strings) (*Todo, error)
func (t *Todo) ToggleWithDeps(timeProxy proxy.Time)
func (t *Todo) UpdateDescriptionWithDeps(description string, timeProxy proxy.Time, stringsProxy proxy.Strings) error
func (t *Todo) MarshalJSONWithDeps(jsonProxy proxy.JSON, timeProxy proxy.Time) ([]byte, error)
func (t *Todo) UnmarshalJSONWithDeps(data []byte, jsonProxy proxy.JSON, timeProxy proxy.Time, stringsProxy proxy.Strings) error
```

### Error Handling Refactoring

#### Current Structure
```go
func (e *DomainError) Error() string {
    return fmt.Sprintf("...")
}
```

#### Refactored Structure
```go
// Existing method remains unchanged for backward compatibility
func (e *DomainError) Error() string

// New method with dependency injection
func (e *DomainError) ErrorWithDeps(fmtProxy proxy.Fmt) string
```

## Data Models

### Proxy Interface Usage

#### Time Proxy
Used for:
- Current time generation (`Now`)
- Time parsing (`Parse`)
- Time formatting (`Format` through time constants)

#### Strings Proxy
Used for:
- String trimming (`TrimSpace`)
- String length validation

#### Fmt Proxy
Used for:
- String formatting (`Sprintf`)
- Error message formatting

#### JSON Proxy
Used for:
- JSON marshaling (`Marshal`)
- JSON unmarshaling (`Unmarshal`)

## Error Handling

Error handling will remain identical to the current implementation:
- Same domain error types and messages
- Wrapped errors with appropriate error types
- Backward compatible error behavior

## Testing Strategy

### Unit Testing Approach
1. **Positive Tests**: Use real proxy implementations to verify actual behavior
2. **Negative Tests**: Use mocked proxy interfaces to simulate error conditions
3. **Time-Controlled Tests**: Use mocked time proxy for consistent timestamp testing
4. **Backward Compatibility Tests**: Verify that existing methods work identically to before

### Test Structure
```go
// Positive test with real implementations
func TestNewTodo(t *testing.T) {
    todo, err := NewTodo(1, "Test todo")
    // Verify actual behavior with real time
}

// Time-controlled test with mocked time
func TestNewTodoWithDeps_ControlledTime(t *testing.T) {
    mockTime := proxy.NewMockTime(ctrl)
    fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    mockTime.EXPECT().Now().Return(fixedTime)
    
    todo, err := NewTodoWithDeps(1, "Test", mockTime, proxy.NewStrings())
    // Verify controlled time behavior
}

// Error simulation test
func TestTodo_MarshalJSONWithDeps_Error(t *testing.T) {
    mockJSON := proxy.NewMockJSON(ctrl)
    mockJSON.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("marshal error"))
    // Test error handling
}
```

### Mock Usage Guidelines
- Use existing mocks from `pkg/proxy/*_mock.go`
- Mock only for error simulation and time control in test cases
- Use real implementations for positive test cases to verify actual behavior

## Implementation Approach

### Phase 1: Add Dependency Injection Methods
1. Add `*WithDeps` methods to Todo entity
2. Add `*WithDeps` methods to DomainError
3. Keep existing methods unchanged for backward compatibility

### Phase 2: Update Internal Helper Functions
1. Update `validateDescription` to accept strings proxy
2. Update internal time operations to use time proxy
3. Update internal formatting to use fmt proxy

### Phase 3: Update Imports and Dependencies
1. Add proxy package imports
2. Update existing methods to use new dependency injection methods internally
3. Verify all functionality works identically

## Backward Compatibility

### Public API Preservation
- All existing public methods remain unchanged
- `Todo` struct remains unchanged
- `DomainError` struct and methods remain unchanged
- All public behavior remains identical

### Migration Path
- Existing code continues to work without changes
- New test code can use `*WithDeps` methods for dependency injection
- No breaking changes to the public API

### Internal Implementation Strategy
```go
// Existing public method (unchanged)
func NewTodo(id int, description string) (*Todo, error) {
    // Internally use dependency injection with real implementations
    return NewTodoWithDeps(id, description, proxy.NewTime(), proxy.NewStrings())
}

// New method with dependency injection (for testing)
func NewTodoWithDeps(id int, description string, timeProxy proxy.Time, stringsProxy proxy.Strings) (*Todo, error) {
    // Implementation using injected dependencies
}
```

## File-by-File Refactoring Plan

### todo.go Refactoring
1. Add proxy imports
2. Add `*WithDeps` methods for all public methods
3. Update existing methods to use `*WithDeps` internally
4. Update helper functions to accept proxy parameters

### errors.go Refactoring
1. Add proxy imports
2. Add `ErrorWithDeps` method
3. Update existing `Error` method to use `ErrorWithDeps` internally

## Testing Enhancement Benefits

### Time-Dependent Testing
- Consistent timestamps in tests
- Ability to test time-based logic
- Predictable test results

### Error Path Testing
- JSON marshaling/unmarshaling error simulation
- String operation error simulation
- Comprehensive error handling coverage

### Integration Testing
- Better control over external dependencies
- More reliable test execution
- Enhanced test isolation