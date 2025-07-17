# Mock Implementation Guidelines

## Overview
This document defines the strict guidelines for mock implementation in the GCT project to ensure consistency and maintainability.

## Mock Usage Rules

### 1. Minimize Mock Usage
- **PREFER** real implementations over mocks whenever possible
- Use mocks **ONLY** for error cases where it's difficult to trigger errors naturally
- Use mocks **ONLY** when testing external dependencies that are hard to control

### 2. When Mocks Are Appropriate
- **Error simulation**: When you need to test error handling but can't easily cause real errors
- **External dependencies**: File system, network, database operations that are hard to control in tests
- **Application layer**: When testing use cases that depend on repositories (unavoidable dependency)

### 3. When Mocks Are NOT Appropriate
- **Domain layer**: Pure functions and entities should use real implementations
- **Configuration layer**: Use real environment variables and file operations
- **Simple operations**: JSON marshaling, string operations, calculations
- **Positive test cases**: Use real implementations to verify actual behavior

### 4. Mock Implementation Rules
- Check for existing mock files in `pkg/proxy/` directory first
- Generate new mocks with mockgen only: `mockgen -source=interface.go -destination=interface_mock.go -package=packagename`
- **NEVER** implement custom mock structs in test files
- Follow pkg/proxy mock file style for consistency

### 5. Test Structure
- **Positive cases**: Use real implementations to verify actual behavior
- **Negative cases**: Use mocks only when necessary to simulate errors
- Always verify the actual output/behavior, not just that methods were called

## Example Mock Generation Commands

```bash
# For domain interfaces
mockgen -source=app/domain/repository.go -destination=app/domain/repository_mock.go -package=domain

# For service interfaces  
mockgen -source=app/service/interface.go -destination=app/service/interface_mock.go -package=service
```

## Verification Checklist

Before implementing any test:
- [ ] Checked pkg/proxy for existing mocks
- [ ] Used existing mocks if available
- [ ] Generated new mocks with mockgen if needed
- [ ] No custom mock implementations in test files
- [ ] Mock file style matches pkg/proxy standards

## Enforcement
This guideline is **MANDATORY** for all development. Any custom mock implementation will be rejected and must be replaced with mockgen-generated mocks.