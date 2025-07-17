# Design-First Development Approach

## Core Principle
**ALWAYS consult design documents before writing any code.**

## Required Pre-Implementation Steps

### 1. Document Review Checklist
Before implementing any feature or component, the following documents MUST be reviewed:

- [ ] **requirements.md** - Understand the functional and non-functional requirements
- [ ] **design.md** - Follow the architectural decisions and component structure
- [ ] **tasks.md** - Understand the specific task scope and dependencies

### 2. Design Alignment Verification
After reviewing documents but before writing code:

- [ ] Verify package structure matches design.md
- [ ] Verify naming conventions match design.md
- [ ] Verify component responsibilities match design.md
- [ ] Verify interfaces match design.md

## Implementation Guidelines

### Directory Structure
- **ALWAYS** follow the exact directory structure specified in design.md
- **NEVER** create new directories without checking design.md first
- **NEVER** place files in incorrect directories

### Naming Conventions
- **ALWAYS** follow the exact naming conventions in design.md
- **NEVER** add redundant suffixes or prefixes not specified in design.md
- **NEVER** deviate from the casing conventions in design.md

### Component Responsibilities
- **ALWAYS** ensure components have single responsibilities as defined in design.md
- **NEVER** add functionality that crosses layer boundaries
- **NEVER** implement features not specified in requirements.md

## Code Review Process

### Self-Review Checklist
Before submitting code, verify:

1. **Structure Compliance**
   - [ ] Directory structure matches design.md
   - [ ] File names match design.md conventions
   - [ ] Package structure matches design.md

2. **Interface Compliance**
   - [ ] Public interfaces match design.md specifications
   - [ ] Method signatures match design.md
   - [ ] Data structures match design.md

3. **Functionality Compliance**
   - [ ] Implemented features satisfy requirements.md
   - [ ] No extra features beyond requirements.md
   - [ ] Error handling follows design.md guidelines

## Enforcement

### Pre-Implementation Verification
Before starting any implementation:

```
// REQUIRED: Document review confirmation
I have reviewed:
- requirements.md: [YES/NO]
- design.md: [YES/NO]
- tasks.md: [YES/NO]

// REQUIRED: Design alignment confirmation
The implementation will follow:
- Directory structure: [SPECIFIC STRUCTURE FROM DESIGN.MD]
- Naming conventions: [SPECIFIC CONVENTIONS FROM DESIGN.MD]
- Component responsibilities: [SPECIFIC RESPONSIBILITIES FROM DESIGN.MD]
```

### Post-Implementation Verification
After completing implementation:

```
// REQUIRED: Compliance verification
I have verified that:
- Directory structure matches design.md: [YES/NO]
- File names match design.md conventions: [YES/NO]
- Interfaces match design.md specifications: [YES/NO]
- Functionality satisfies requirements.md: [YES/NO]
```

## Example Application

### Correct Approach - CLI Main Structure
```
// REQUIRED: Document review confirmation
I have reviewed:
- requirements.md: YES
- design.md: YES
- tasks.md: YES

// REQUIRED: Design alignment confirmation
The implementation will follow:
- Directory structure: app/presentation/cli/gct/
- Naming conventions: main.go (simple execution), commands/command.go (initialization)
- Component responsibilities: main.go executes, command.go initializes dependencies
```

### Correct Approach - Formatter Implementation
```
// REQUIRED: Document review confirmation
I have reviewed:
- requirements.md: YES
- design.md: YES
- tasks.md: YES

// REQUIRED: Design alignment confirmation
The implementation will follow:
- Directory structure: app/presentation/cli/gct/formatter/
- Naming conventions: json.go, table.go, plain.go (no _formatter suffix)
- Component responsibilities: Format todos as JSON, table, or plain text
```

### Incorrect Approach (AVOID)
```
// Implementing without checking design.md
// Creating app/presentation/formatter/ directory
// Using _formatter suffix in filenames
// Misplacing components in wrong directories
// Putting initialization logic directly in main.go
```

This steering document is MANDATORY for all development work.