# Testing Philosophy and Mock Usage Guidelines

## Core Testing Philosophy

### Real Implementation First

- **Default approach**: Always prefer real implementations over mocks
- **Value real behavior**: Tests should verify actual functionality, not just method calls
- **Integration over isolation**: Test components working together when possible

### Mock Usage Principles

#### When Mocks Are Justified

1. **Error simulation**: When testing error handling paths that are difficult to trigger naturally
2. **External dependencies**: File system, network, database operations that are unreliable in tests
3. **Unavoidable dependencies**: Application layer testing use cases with repository dependencies

#### When Mocks Are NOT Justified

1. **Domain layer**: Pure functions, entities, value objects - use real implementations
2. **Simple operations**: JSON marshaling, string manipulation, calculations
3. **Configuration**: Environment variables, file paths - use real values
4. **Positive test cases**: Verify actual behavior with real implementations

### Test Structure Guidelines

#### Positive Test Cases

```go
// GOOD: Use real implementation
func TestJSONFormatter_Format(t *testing.T) {
    jsonProxy := proxy.NewJSON()  // Real implementation
    formatter := NewJSONFormatter(jsonProxy)
    result, err := formatter.Format(todos)
    // Verify actual JSON output
}
```

#### Negative Test Cases

```go
// ACCEPTABLE: Mock only for error simulation
func TestJSONFormatter_Format_Error(t *testing.T) {
    mockJSON := proxy.NewMockJSON(ctrl)
    mockJSON.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("marshal error"))
    // Test error handling
}
```

### Anti-Patterns to Avoid

#### Over-Mocking

```go
// BAD: Mocking everything unnecessarily
func TestAddTodo(t *testing.T) {
    mockTime := NewMockTime(ctrl)
    mockValidator := NewMockValidator(ctrl)
    mockLogger := NewMockLogger(ctrl)
    // Too many mocks for simple functionality
}
```

#### Mock-Heavy Tests

```go
// BAD: Testing mock interactions instead of behavior
func TestUseCase(t *testing.T) {
    mock.EXPECT().Method1().Times(1)
    mock.EXPECT().Method2().Times(1)
    // Only verifying method calls, not actual behavior
}
```

### Preferred Test Patterns

#### Behavior Verification

```go
// GOOD: Verify actual output/behavior
func TestCalculateTotal(t *testing.T) {
    calculator := NewCalculator()
    result := calculator.Calculate(items)
    if result != expectedTotal {
        t.Errorf("got %v, want %v", result, expectedTotal)
    }
}
```

#### State Verification

```go
// GOOD: Verify state changes
func TestAddTodo(t *testing.T) {
    useCase := NewAddTodoUseCase(realRepository)
    todo, err := useCase.Run("Buy milk")
    // Verify the todo was created correctly
    if todo.Description != "Buy milk" {
        t.Errorf("got %v, want %v", todo.Description, "Buy milk")
    }
}
```

### Implementation Guidelines

#### Layer-Specific Approaches

**Domain Layer**

- Use real implementations exclusively
- Test pure business logic
- No external dependencies to mock

**Application Layer**

- Mock repository interfaces (unavoidable dependency)
- Use real domain entities
- Focus on business flow testing

**Infrastructure Layer**

- Mock external systems (file system, network)
- Use real domain entities
- Test integration points

**Presentation Layer**

- Mock external dependencies only for error cases
- Use real formatters and converters
- Test actual output format

### Quality Metrics

#### Good Test Indicators

- Tests pass with real implementations
- Tests verify actual behavior/output
- Tests are readable and maintainable
- Minimal mock setup code

#### Warning Signs

- More mock setup than actual test logic
- Tests only verify method calls
- Brittle tests that break on refactoring
- Complex mock interaction chains

### Enforcement

This philosophy should guide all testing decisions. When in doubt:

1. Try real implementation first
2. Add mocks only when absolutely necessary
3. Prefer behavior verification over interaction verification
4. Keep tests simple and focused on actual functionality
5. Write a `*_test.go` for each `*_.go` file

- e.g.) `main_test.go` for `main.go`

6. Write a test method for each method in the `*_.go` file

- e.g.) There are three methods below in `use_case.go` ...
    - `Run()`, `Get()`, `Delete()`
    - then, write three test methods below in `use_case_test.go` ...
        - `TestRun()`, `TestGet()`, `TestDelete()`

