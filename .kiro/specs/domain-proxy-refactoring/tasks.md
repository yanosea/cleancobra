# Implementation Plan

- [x] 1. Add proxy imports and dependency injection structure to todo.go
  - Add imports for pkg/proxy.Time, pkg/proxy.Strings, pkg/proxy.Fmt, pkg/proxy.JSON
  - Remove direct imports of time, strings, fmt, encoding/json packages
  - Verify proxy interfaces are available and compatible
  - _Requirements: 1.1, 3.1, 3.4_

- [ ] 2. Refactor Todo entity constructor with dependency injection
  - [x] 2.1 Create NewTodoWithDeps function with proxy parameters
    - Add function signature accepting id, description, timeProxy, stringsProxy parameters
    - Implement todo creation using proxy.Time.Now() instead of time.Now()
    - Use proxy.Strings.TrimSpace() instead of strings.TrimSpace()
    - Call validateDescriptionWithDeps with proxy parameters
    - _Requirements: 1.1, 1.3, 3.2_

  - [ ] 2.2 Update existing NewTodo function for backward compatibility
    - Modify NewTodo to create real proxy implementations
    - Call NewTodoWithDeps with real proxy instances
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.2, 3.4_

- [ ] 3. Refactor Todo entity methods with dependency injection
  - [ ] 3.1 Create ToggleWithDeps method
    - Add method signature accepting timeProxy parameter
    - Replace time.Now() with proxy.Time.Now()
    - Implement identical toggle logic with injected time dependency
    - _Requirements: 1.1, 1.3, 3.2_

  - [ ] 3.2 Update existing Toggle method for backward compatibility
    - Modify Toggle to create real time proxy implementation
    - Call ToggleWithDeps with real proxy instance
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.3, 3.4_

  - [ ] 3.3 Create UpdateDescriptionWithDeps method
    - Add method signature accepting description, timeProxy, stringsProxy parameters
    - Use proxy.Strings.TrimSpace() for description processing
    - Use proxy.Time.Now() for timestamp updates
    - Call validateDescriptionWithDeps with proxy parameters
    - _Requirements: 1.1, 1.3, 3.2_

  - [ ] 3.4 Update existing UpdateDescription method for backward compatibility
    - Modify UpdateDescription to create real proxy implementations
    - Call UpdateDescriptionWithDeps with real proxy instances
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.3, 3.4_

- [ ] 4. Refactor JSON marshaling methods with dependency injection
  - [ ] 4.1 Create MarshalJSONWithDeps method
    - Add method signature accepting jsonProxy and timeProxy parameters
    - Replace json.Marshal with proxy.JSON.Marshal
    - Use proxy.Time constants for RFC3339 formatting
    - Implement identical marshaling logic with injected dependencies
    - _Requirements: 1.1, 1.4, 3.2_

  - [ ] 4.2 Update existing MarshalJSON method for backward compatibility
    - Modify MarshalJSON to create real proxy implementations
    - Call MarshalJSONWithDeps with real proxy instances
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.3, 3.4_

  - [ ] 4.3 Create UnmarshalJSONWithDeps method
    - Add method signature accepting data, jsonProxy, timeProxy parameters
    - Replace json.Unmarshal with proxy.JSON.Unmarshal
    - Use proxy.Time.Parse for timestamp parsing
    - Implement identical unmarshaling logic with injected dependencies
    - _Requirements: 1.1, 1.4, 3.2_

  - [ ] 4.4 Update existing UnmarshalJSON method for backward compatibility
    - Modify UnmarshalJSON to create real proxy implementations
    - Call UnmarshalJSONWithDeps with real proxy instances
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.3, 3.4_

- [ ] 5. Refactor helper functions with dependency injection
  - [ ] 5.1 Create validateDescriptionWithDeps function
    - Add function signature accepting description and stringsProxy parameters
    - Replace strings.TrimSpace with proxy.Strings.TrimSpace
    - Implement identical validation logic with injected dependencies
    - _Requirements: 1.1, 1.3, 3.2_

  - [ ] 5.2 Update existing validateDescription function for backward compatibility
    - Modify validateDescription to create real strings proxy implementation
    - Call validateDescriptionWithDeps with real proxy instance
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.3, 3.4_

  - [ ] 5.3 Create StringWithDeps method for Todo
    - Add method signature accepting fmtProxy parameter
    - Replace fmt.Sprintf with proxy.Fmt.Sprintf
    - Implement identical string formatting with injected dependencies
    - _Requirements: 1.1, 1.4, 3.2_

  - [ ] 5.4 Update existing String method for backward compatibility
    - Modify String to create real fmt proxy implementation
    - Call StringWithDeps with real proxy instance
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.3, 3.4_

- [ ] 6. Refactor errors.go with dependency injection
  - [ ] 6.1 Add proxy imports to errors.go
    - Add import for pkg/proxy.Fmt
    - Remove direct import of fmt package
    - Verify proxy interface compatibility
    - _Requirements: 1.1, 3.1, 3.4_

  - [ ] 6.2 Create ErrorWithDeps method for DomainError
    - Add method signature accepting fmtProxy parameter
    - Replace fmt.Sprintf with proxy.Fmt.Sprintf
    - Implement identical error formatting with injected dependencies
    - _Requirements: 1.1, 1.4, 3.2_

  - [ ] 6.3 Update existing Error method for backward compatibility
    - Modify Error to create real fmt proxy implementation
    - Call ErrorWithDeps with real proxy instance
    - Ensure identical behavior to original implementation
    - _Requirements: 2.1, 2.4, 3.4_

- [ ] 7. Write comprehensive unit tests for todo.go
  - [ ] 7.1 Write positive tests using real proxy implementations
    - Test NewTodo and NewTodoWithDeps with real implementations
    - Test Toggle and ToggleWithDeps with real implementations
    - Test UpdateDescription and UpdateDescriptionWithDeps with real implementations
    - Test MarshalJSON and MarshalJSONWithDeps with real implementations
    - Test UnmarshalJSON and UnmarshalJSONWithDeps with real implementations
    - Verify actual behavior matches original implementation
    - _Requirements: 2.2, 2.3, 4.1, 4.2, 4.3, 4.4_

  - [ ] 7.2 Write time-controlled tests using mocked time proxy
    - Test NewTodoWithDeps with fixed time for consistent timestamps
    - Test ToggleWithDeps with controlled time progression
    - Test UpdateDescriptionWithDeps with predictable time updates
    - Test JSON marshaling with controlled time formatting
    - _Requirements: 4.1, 4.2, 4.3_

  - [ ] 7.3 Write negative tests using mocked proxy interfaces
    - Test JSON marshaling error handling with mock JSON proxy
    - Test JSON unmarshaling error handling with mock JSON proxy
    - Test time parsing error handling with mock time proxy
    - Test string operation edge cases with mock strings proxy
    - _Requirements: 1.1, 1.4, 4.4_

  - [ ] 7.4 Write backward compatibility tests
    - Verify NewTodo behavior matches original implementation
    - Verify Toggle behavior matches original implementation
    - Verify UpdateDescription behavior matches original implementation
    - Verify MarshalJSON behavior matches original implementation
    - Verify UnmarshalJSON behavior matches original implementation
    - Test that existing test cases continue to pass
    - _Requirements: 2.1, 2.2, 2.3_

- [ ] 8. Write comprehensive unit tests for errors.go
  - [ ] 8.1 Write positive tests using real proxy implementations
    - Test Error and ErrorWithDeps with real fmt proxy
    - Verify actual error message formatting behavior
    - Test error wrapping and unwrapping functionality
    - _Requirements: 2.4, 4.4_

  - [ ] 8.2 Write negative tests using mocked proxy interfaces
    - Test error formatting with mock fmt proxy
    - Test complex error message scenarios
    - Verify error type checking functionality
    - _Requirements: 1.4, 4.4_

  - [ ] 8.3 Write backward compatibility tests
    - Verify Error method behavior matches original implementation
    - Test that existing error handling continues to work
    - Verify domain error helper functions remain unchanged
    - _Requirements: 2.1, 2.4_

- [ ] 9. Integration testing and validation
  - [ ] 9.1 Run existing test suite to ensure no regressions
    - Execute all existing domain layer tests
    - Verify no test failures or behavior changes
    - Confirm backward compatibility is maintained
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

  - [ ] 9.2 Validate proxy interface usage consistency
    - Verify all proxy interfaces are used correctly
    - Check that proxy constructors follow established patterns
    - Ensure consistent error handling across all methods
    - _Requirements: 3.1, 3.2, 3.4_

  - [ ] 9.3 Performance validation
    - Compare performance of original vs refactored methods
    - Ensure no significant performance degradation
    - Verify memory usage remains consistent
    - _Requirements: 2.1, 2.2, 2.3, 2.4_