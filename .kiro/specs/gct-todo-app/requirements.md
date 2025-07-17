# Requirements Document

## Introduction

「gct」（Go Clean-Architecture Todo）は、クリーンアーキテクチャの原則に従って構築されたTODOアプリケーションです。CLI（Command Line Interface）とTUI（Text User Interface）の両方のインターフェースを提供し、JSONベースのストレージを使用してTODOを管理します。アプリケーションは環境変数による設定をサポートし、テスタビリティを重視した設計となっています。

## Requirements

### Requirement 1

**User Story:** As a user, I want to add new todos via CLI, so that I can quickly capture tasks that need to be done

#### Acceptance Criteria

1. WHEN user executes `gct add "task description"` THEN system SHALL create a new todo with the provided description
2. WHEN user executes `gct add` without description THEN system SHALL display error message
3. WHEN todo is successfully added THEN system SHALL display confirmation message
4. WHEN todo is added THEN system SHALL save it to the configured JSON file
5. IF JSON file does not exist THEN system SHALL create it with proper directory structure

### Requirement 2

**User Story:** As a user, I want to list all todos via CLI, so that I can see what tasks I need to complete

#### Acceptance Criteria

1. WHEN user executes `gct list` OR `gct` THEN system SHALL display all todos with their status
2. WHEN user executes `gct list --format json` THEN system SHALL output todos in JSON format
3. WHEN user executes `gct list --format plain` THEN system SHALL output todos in plain text format
4. WHEN no todos exist THEN system SHALL display appropriate message
5. WHEN todos exist THEN system SHALL display each todo with ID, description, and completion status
6. WHEN displaying todos THEN system SHALL use clear formatting with colors for better readability

### Requirement 3

**User Story:** As a user, I want to toggle todo completion status via CLI, so that I can mark tasks as done or undone

#### Acceptance Criteria

1. WHEN user executes `gct toggle <id>` THEN system SHALL toggle the completion status of the specified todo
2. WHEN todo ID does not exist THEN system SHALL display error message
3. WHEN todo is successfully toggled THEN system SHALL display confirmation message with new status
4. WHEN todo status is changed THEN system SHALL save changes to JSON file

### Requirement 4

**User Story:** As a user, I want to delete todos via CLI, so that I can remove completed or unwanted tasks

#### Acceptance Criteria

1. WHEN user executes `gct delete <id>` THEN system SHALL remove the specified todo
2. WHEN todo ID does not exist THEN system SHALL display error message
3. WHEN todo is successfully deleted THEN system SHALL display confirmation message
4. WHEN todo is deleted THEN system SHALL save changes to JSON file and update remaining todo IDs

### Requirement 5

**User Story:** As a user, I want to manage todos via TUI, so that I can have an interactive experience for todo management

#### Acceptance Criteria

1. WHEN user executes `gct-tui` THEN system SHALL launch interactive TUI interface
2. WHEN in TUI mode THEN system SHALL display all todos with keyboard navigation support
3. WHEN user presses designated keys THEN system SHALL allow adding, toggling, and deleting todos
4. WHEN changes are made in TUI THEN system SHALL immediately update the display and save to JSON file
5. WHEN user exits TUI THEN system SHALL gracefully close the application

### Requirement 6

**User Story:** As a user, I want to configure data storage location, so that I can control where my todos are stored

#### Acceptance Criteria

1. WHEN `GCT_DATA_FILE` environment variable is set THEN system SHALL use that path for todo storage
2. WHEN `GCT_DATA_FILE` is not set AND `XDG_DATA_HOME` is set THEN system SHALL use `$XDG_DATA_HOME/gct/todos.json`
3. WHEN neither environment variable is set THEN system SHALL use `~/.local/share/gct/todos.json`
4. WHEN data directory does not exist THEN system SHALL create it automatically
5. WHEN JSON file is corrupted THEN system SHALL handle error gracefully and create new file

### Requirement 7

**User Story:** As a developer, I want the application to follow clean architecture principles, so that the code is maintainable and testable

#### Acceptance Criteria

1. WHEN examining code structure THEN system SHALL have separate domain, application, infrastructure, and presentation layers
2. WHEN examining dependencies THEN system SHALL follow dependency inversion principle
3. WHEN running tests THEN system SHALL achieve 100% test coverage for app/ directory code
4. WHEN examining TUI code THEN system SHALL implement both Clean Architecture and ELM architecture patterns
5. WHEN examining CLI code THEN system SHALL have one file per subcommand implementation

### Requirement 8

**User Story:** As a developer, I want proper package organization, so that the code is well-structured and easy to navigate

#### Acceptance Criteria

1. WHEN examining project structure THEN system SHALL organize code according to specified directory layout
2. WHEN examining CLI implementation THEN system SHALL be executable via `go run ./app/presentation/cli/gct/main.go`
3. WHEN examining TUI implementation THEN system SHALL be executable via `go run ./app/presentation/tui/gct-tui/main.go`
4. WHEN examining proxy packages THEN system SHALL provide testable wrappers for standard and third-party packages
5. WHEN examining composition root THEN system SHALL provide proper dependency injection setup

### Requirement 9

**User Story:** As a user, I want reliable data persistence, so that my todos are safely stored and retrieved

#### Acceptance Criteria

1. WHEN todos are modified THEN system SHALL save changes to JSON file immediately
2. WHEN application starts THEN system SHALL load existing todos from JSON file
3. WHEN JSON file has invalid format THEN system SHALL handle error and create new valid file
4. WHEN file system operations fail THEN system SHALL display appropriate error messages
5. WHEN concurrent access occurs THEN system SHALL handle file locking appropriately

### Requirement 10

**User Story:** As a user, I want consistent command-line experience, so that I can efficiently use the CLI interface

#### Acceptance Criteria

1. WHEN using CLI commands THEN system SHALL provide consistent help messages and usage information
2. WHEN command fails THEN system SHALL display clear error messages with suggested solutions
3. WHEN using CLI THEN system SHALL support shell completion for commands and options
4. WHEN user executes `gct completion bash|zsh|fish|powershell` THEN system SHALL generate appropriate completion script for the specified shell
5. WHEN user executes `gct completion` without shell argument THEN system SHALL display error message requiring shell specification
4. WHEN formatting output THEN system SHALL provide both human-readable and JSON formats
5. WHEN no arguments provided to root command THEN system SHALL execute list command by default