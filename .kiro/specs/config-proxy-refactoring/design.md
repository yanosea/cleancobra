# Design Document

## Overview

This design outlines the refactoring of the `app/config` package to use proxy interfaces instead of direct standard library imports. The refactoring will improve testability while maintaining backward compatibility and following the established proxy pattern used throughout the GCT application.

## Architecture

### Current Architecture
```
config.go
├── Direct imports: os, path/filepath, envconfig
├── Load() function with hardcoded dependencies
└── Helper functions with direct standard library calls
```

### Target Architecture
```
config.go
├── Proxy imports: pkg/proxy.OS, pkg/proxy.Filepath, pkg/proxy.Envconfig
├── Load() function with dependency injection
├── LoadWithDependencies() function for explicit dependency injection
└── Helper functions using injected proxy interfaces
```

## Components and Interfaces

### Config Struct
The `Config` struct will remain unchanged to maintain backward compatibility:
```go
type Config struct {
    DataFile string `envconfig:"GCT_DATA_FILE"`
}
```

### Dependency Structure
A new internal struct will hold the proxy dependencies:
```go
type configDependencies struct {
    os       proxy.OS
    filepath proxy.Filepath
    envconfig proxy.Envconfig
}
```

### Public API Functions

#### Load() Function
- **Purpose**: Maintain backward compatibility
- **Behavior**: Create real proxy implementations and call LoadWithDependencies
- **Signature**: `func Load() (*Config, error)`

#### LoadWithDependencies() Function
- **Purpose**: Enable dependency injection for testing
- **Behavior**: Accept proxy interfaces as parameters
- **Signature**: `func LoadWithDependencies(os proxy.OS, filepath proxy.Filepath, envconfig proxy.Envconfig) (*Config, error)`

### Helper Functions Refactoring

#### getDefaultDataFilePath()
- **Current**: Direct calls to `os.Getenv()` and `os.UserHomeDir()`
- **Refactored**: Accept `os` and `filepath` proxy interfaces as parameters
- **Signature**: `func getDefaultDataFilePath(os proxy.OS, filepath proxy.Filepath) (string, error)`

#### ensureDirectoryExists()
- **Current**: Direct calls to `os.Stat()`, `os.IsNotExist()`, `os.MkdirAll()`
- **Refactored**: Accept `os` proxy interface as parameter
- **Signature**: `func ensureDirectoryExists(dir string, os proxy.OS) error`

## Data Models

### Proxy Interface Usage

#### OS Proxy
Used for:
- Environment variable access (`Getenv`)
- Home directory detection (`UserHomeDir`)
- Directory creation (`MkdirAll`)
- File system operations (`Stat`, `IsNotExist`)

#### Filepath Proxy
Used for:
- Path joining (`Join`)
- Directory extraction (`Dir`)

#### Envconfig Proxy
Used for:
- Configuration processing (`Process`)

## Error Handling

Error handling will remain identical to the current implementation:
- Domain errors for configuration issues
- Wrapped errors with appropriate error types
- Same error messages and error types

## Testing Strategy

### Unit Testing Approach
1. **Positive Tests**: Use real proxy implementations to verify actual behavior
2. **Negative Tests**: Use mocked proxy interfaces to simulate error conditions
3. **Backward Compatibility Tests**: Verify that `Load()` function works identically to before

### Test Structure
```go
// Positive test with real implementations
func TestLoad(t *testing.T) {
    config, err := Load()
    // Verify actual behavior
}

// Negative test with mocked dependencies
func TestLoadWithDependencies_EnvconfigError(t *testing.T) {
    mockEnvconfig := proxy.NewMockEnvconfig(ctrl)
    mockEnvconfig.EXPECT().Process(gomock.Any(), gomock.Any()).Return(errors.New("envconfig error"))
    // Test error handling
}
```

### Mock Usage Guidelines
- Use existing mocks from `pkg/proxy/*_mock.go`
- Mock only for error simulation in negative test cases
- Use real implementations for positive test cases to verify actual behavior

## Implementation Approach

### Phase 1: Add Dependency Injection Support
1. Create `configDependencies` struct
2. Add `LoadWithDependencies()` function
3. Refactor helper functions to accept proxy parameters

### Phase 2: Refactor Existing Functions
1. Update `Load()` to use `LoadWithDependencies()` with real implementations
2. Update `Validate()` method to use proxy interfaces
3. Update all helper functions to use injected dependencies

### Phase 3: Update Imports
1. Remove direct standard library imports
2. Add proxy package imports
3. Verify all functionality works identically

## Backward Compatibility

### Public API Preservation
- `Load()` function signature remains unchanged
- `Config` struct remains unchanged
- `Validate()` method signature remains unchanged
- All public behavior remains identical

### Migration Path
- Existing code continues to work without changes
- New test code can use `LoadWithDependencies()` for dependency injection
- No breaking changes to the public API