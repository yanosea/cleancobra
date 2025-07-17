# Development Guide

## Overview

This project implements a TODO application following Clean Architecture principles with both CLI and TUI interfaces.

## Project Structure

```
app/
├── config/              # Configuration management
├── container/           # Dependency injection container
├── domain/              # Domain layer - core business logic
├── application/         # Application layer - use cases
├── infrastructure/      # Infrastructure layer - external interfaces
└── presentation/        # Presentation layer
    ├── cli/             # CLI interface (Cobra)
    └── tui/             # TUI interface (Bubbletea + ELM)
```

## Development Setup

1. **Install dependencies:**
   ```bash
   make setup
   ```

2. **Generate mocks:**
   ```bash
   make generate
   ```

3. **Run tests:**
   ```bash
   make test
   ```

4. **Build applications:**
   ```bash
   make build
   ```

## Testing

### Running Tests
- All tests: `make test`
- With coverage: `make test-coverage`
- Specific package: `go test ./app/domain -v`

### Test Philosophy
- Prefer real implementations over mocks
- Use mocks only for error simulation and external dependencies
- Follow table-driven test patterns
- Aim for 100% coverage in app/ directory

### Mock Generation
Mocks are generated using `gomock`:
```bash
//go:generate mockgen -source=interface.go -destination=interface_mock.go -package=packagename
```

## Architecture

### Clean Architecture Layers

1. **Domain Layer** (`app/domain/`)
   - Pure business logic
   - No external dependencies
   - Entities and repository interfaces

2. **Application Layer** (`app/application/`)
   - Use cases and business workflows
   - Orchestrates domain objects
   - Depends only on domain layer

3. **Infrastructure Layer** (`app/infrastructure/`)
   - External concerns (file system, databases)
   - Implements domain interfaces
   - JSON repository implementation

4. **Presentation Layer** (`app/presentation/`)
   - User interfaces (CLI and TUI)
   - Depends on application layer
   - Handles user input/output

### TUI Architecture (ELM Pattern)

The TUI follows the ELM (Elm Architecture) pattern:

- **Model** (`model/`): Application state
- **Update** (`update/`): State transitions and message handling  
- **View** (`view/`): UI rendering

## Usage Examples

### CLI Usage
```bash
# List todos
./bin/gct list

# Add todo
./bin/gct add "Buy groceries"

# Toggle todo completion
./bin/gct toggle 1

# Delete todo
./bin/gct delete 1

# JSON output
./bin/gct list --format json
```

### TUI Usage
```bash
# Launch interactive TUI
./bin/gct-tui
```

TUI Controls:
- `↑/k`: Move up
- `↓/j`: Move down  
- `Space`: Toggle todo
- `a`: Add todo
- `e`: Edit todo
- `d`: Delete todo
- `q`: Quit

## Configuration

Set data file location via environment variable:
```bash
export GCT_DATA_FILE=/path/to/todos.json
```

Default locations:
1. `$GCT_DATA_FILE` (if set)
2. `$XDG_DATA_HOME/gct/todos.json` (if XDG_DATA_HOME is set)
3. `~/.local/share/gct/todos.json` (fallback)

## Contributing

1. Follow Clean Architecture principles
2. Write tests for all new code
3. Use table-driven tests
4. Generate mocks with `gomock`
5. Run `make format` before committing
6. Ensure `make test` passes