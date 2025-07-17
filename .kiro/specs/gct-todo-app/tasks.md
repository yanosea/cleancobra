# Implementation Plan

- [x] 1. Set up project foundation and proxy layer





  - Create proxy interfaces for all external dependencies (os, filepath, json, time, io, fmt, strings, strconv, cobra, bubbletea, bubbles, lipgloss, color, envconfig)
  - Implement proxy implementations for standard and third-party packages
  - Set up gomock generation with proper //go:generate directives
  - _Requirements: 7.1, 7.2, 8.4_

- [x] 2. Implement domain layer with core business logic






  - Create Todo entity with proper JSON tags and validation
  - Define TodoRepository interface with FindAll, Save, and DeleteById methods
  - Implement domain-specific error types and error handling
  - Write comprehensive unit tests for domain entities
  - _Requirements: 7.1, 9.1, 9.3_

- [x] 3. Create configuration management system (app/config)



  - Implement Config struct with envconfig tags for GCT_DATA_FILE
  - Create configuration loading logic with XDG_DATA_HOME and fallback path support
  - Add configuration validation and error handling
  - Write unit tests for configuration loading scenarios
  - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [x] 4. Implement JSON-based repository infrastructure




  - Create JSONRepository struct implementing TodoRepository interface
  - Implement FindAll method with JSON file reading and parsing
  - Implement Save method supporting both new and existing todos with ID management
  - Implement DeleteById method with proper ID resequencing
  - Add file system error handling and recovery mechanisms
  - Write comprehensive unit tests with file system mocking
  - _Requirements: 6.5, 9.1, 9.2, 9.3, 9.4, 9.5_

- [x] 5. Build application layer use cases



- [x] 5.1 Implement AddTodoUseCase



  - Create AddTodoUseCase struct with repository dependency
  - Implement Run method for adding new todos with validation
  - Add proper error handling for invalid inputs
  - Write table-driven unit tests with positive and negative scenarios
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [x] 5.2 Implement ListTodoUseCase



  - Create ListTodoUseCase struct with repository dependency
  - Implement Run method for retrieving all todos
  - Add handling for empty todo lists
  - Write table-driven unit tests covering all scenarios
  - _Requirements: 2.1, 2.3_

- [x] 5.3 Implement ToggleTodoUseCase



  - Create ToggleTodoUseCase struct with repository dependency
  - Implement Run method for toggling todo completion status
  - Add validation for todo ID existence
  - Write table-driven unit tests with ID validation scenarios
  - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [x] 5.4 Implement DeleteTodoUseCase



  - Create DeleteTodoUseCase struct with repository dependency
  - Implement Run method for deleting todos by ID
  - Add validation for todo ID existence
  - Write table-driven unit tests with ID validation scenarios




  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [x] 6. Create CLI presentation layer formatters
- [x] 6.1 Implement JSON formatter


  - Create JSONFormatter struct for JSON output formatting
  - Implement Format method for todos and single todo items
  - Add proper JSON marshaling with error handling
  - Write unit tests for JSON formatting scenarios
  - _Requirements: 2.2, 10.4_

- [x] 6.2 Implement table formatter








  - Create TableFormatter struct for human-readable table output
  - Implement Format method with colored output using color proxy
  - Add proper column alignment and status indicators
  - Write unit tests for table formatting scenarios
  - _Requirements: 2.4, 2.5, 10.4_

- [x] 6.3 Implement plain text formatter



  - Create PlainFormatter struct for simple text output
  - Implement Format method without colors or special formatting
  - Add basic text representation of todos
  - Write unit tests for plain text formatting scenarios
  - _Requirements: 10.4_

- [x] 7. Build CLI presentation layer components



- [x] 7.1 Implement TodoPresenter




  - Create TodoPresenter struct with formatter dependencies
  - Implement methods for showing add success, list results, toggle success, delete success
  - Add error presentation methods with user-friendly messages
  - Write unit tests for all presentation scenarios
  - _Requirements: 10.1, 10.2_

- [x] 7.2 Implement root command






  - Create NewRootCommand function using cobra proxy
  - Set up root command to execute list command by default
  - Add global flags and configuration
  - Implement command execution logic with proper error handling
  - Write unit tests for root command behavior
  - _Requirements: 2.1, 8.2, 10.1, 10.5_

- [x] 7.3 Implement add command





  - Create NewAddCommand function using cobra proxy pattern
  - Implement runAdd function with use case integration
  - Add argument validation and error handling for missing description
  - Write unit tests for add command scenarios including error cases
  - _Requirements: 1.1, 1.2, 1.3, 7.5, 10.1, 10.2_

- [x] 7.4 Implement list command




  - Create NewListCommand function using cobra proxy pattern
  - Implement runList function with format flag support (json, table, plain)
  - Add output formatting logic with presenter integration
  - Write unit tests for list command with all format options
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 7.5, 10.4_

- [x] 7.5 Implement toggle command





  - Create NewToggleCommand function using cobra proxy pattern
  - Implement runToggle function with ID validation
  - Add proper error handling for invalid IDs
  - Write unit tests for toggle command scenarios
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 7.5, 10.1, 10.2_

- [x] 7.6 Implement delete command





  - Create NewDeleteCommand function using cobra proxy pattern
  - Implement runDelete function with ID validation
  - Add confirmation and error handling
  - Write unit tests for delete command scenarios
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 7.5, 10.1, 10.2_

- [x] 7.7 Implement shell completion commands





- [x] 7.7.1 Create completion parent command


  - Create NewCompletionCommand function using cobra proxy pattern
  - Set up parent command that requires subcommand (bash, zsh, fish, powershell)
  - Add help text and usage information for completion setup
  - Write unit tests for completion parent command
  - _Requirements: 10.3_

- [x] 7.7.2 Implement bash completion subcommand


  - Create NewBashCompletionCommand function using cobra proxy pattern
  - Implement runBashCompletion function to generate bash completion script
  - Add proper bash completion script generation with cobra functionality
  - Write unit tests for bash completion generation
  - _Requirements: 10.3_

- [x] 7.7.3 Implement zsh completion subcommand


  - Create NewZshCompletionCommand function using cobra proxy pattern
  - Implement runZshCompletion function to generate zsh completion script
  - Add proper zsh completion script generation with cobra functionality
  - Write unit tests for zsh completion generation
  - _Requirements: 10.3_

- [x] 7.7.4 Implement fish completion subcommand


  - Create NewFishCompletionCommand function using cobra proxy pattern
  - Implement runFishCompletion function to generate fish completion script
  - Add proper fish completion script generation with cobra functionality
  - Write unit tests for fish completion generation
  - _Requirements: 10.3_

- [x] 7.7.5 Implement powershell completion subcommand


  - Create NewPowershellCompletionCommand function using cobra proxy pattern
  - Implement runPowershellCompletion function to generate powershell completion script
  - Add proper powershell completion script generation with cobra functionality
  - Write unit tests for powershell completion generation
  - _Requirements: 10.3_

- [x] 8. Create dependency injection container (app/container)





  - Implement Container struct with all dependency management
  - Create Proxies struct to organize all proxy dependencies
  - Create UseCases struct to organize all use case dependencies
  - Implement NewContainer function with proper initialization order
  - Add getter methods for accessing configured dependencies
  - Write comprehensive unit tests for container initialization and dependency wiring
  - _Requirements: 7.2, 8.4_

- [x] 9. Create composition root and CLI main






  - Create CLI main.go with simple execution logic (InitializeCommand + Execute)
  - Create commands/command.go with InitializeCommand function for dependency injection
  - Implement proper initialization and command wiring with container
  - Add graceful error handling and exit codes
  - Write integration tests for CLI application startup and command initialization
  - _Requirements: 7.2, 8.2, 10.3_

- [x] 10. Implement TUI model layer (ELM Architecture)





- [x] 10.1 Create TUI item model


  - Implement item model for TUI with bubbletea integration
  - Add model state management for individual todo items
  - Create model initialization and update methods
  - Write unit tests for item model behavior
  - _Requirements: 5.2, 5.3, 7.1, 7.4_

- [x] 10.2 Create TUI state model


  - Implement main application state model with todo list state
  - Add cursor management, selection state, and input modes
  - Create model initialization with use case dependencies
  - Write unit tests for application state model management
  - _Requirements: 5.2, 5.3, 7.1, 7.4_

- [x] 11. Implement TUI update layer (ELM Architecture)





- [x] 11.1 Create TUI item update logic


  - Implement item-specific update functions for state changes
  - Add message handling for item operations (add, toggle, delete)
  - Create command generation for async operations
  - Write unit tests for item update logic
  - _Requirements: 5.4, 7.1, 7.4_

- [x] 11.2 Create TUI handler update logic


  - Implement main update function with message routing
  - Add keyboard event handling for navigation and actions
  - Create mode switching logic (normal, input, confirmation)
  - Write unit tests for handler update scenarios
  - _Requirements: 5.2, 5.3, 5.4, 7.1, 7.4_

- [x] 12. Implement TUI view layer (ELM Architecture)





- [x] 12.1 Create TUI item view components


  - Implement todo item rendering with lipgloss styling
  - Add status indicators and selection highlighting
  - Create responsive layout for different terminal sizes
  - Write unit tests for item view rendering
  - _Requirements: 5.2, 5.3, 7.1, 7.4_

- [x] 12.2 Create TUI layout view


  - Implement main application layout with header, todo list, and footer
  - Add help text and status information display
  - Create input field rendering for add mode
  - Write unit tests for layout view composition
  - _Requirements: 5.2, 5.3, 5.4, 7.1, 7.4_

- [x] 13. Create TUI composition root and main





  - Implement dependency injection container for TUI application
  - Create TUI main.go with bubbletea program initialization
  - Add proper cleanup and error handling
  - Write integration tests for TUI application startup
  - _Requirements: 5.1, 7.2, 8.3_

- [x] 13.1 Refactor TUI to CLI-like structure







  - Refactor TUI main.go to simple execution pattern like CLI
  - Create program/program.go for initialization logic (like CLI commands/command.go)
  - Split large files into focused, single-responsibility modules
  - Maintain model/update/view structure while improving organization
  - _Requirements: 7.1, 7.2, 8.1_

- [x] 14. Add comprehensive test coverage and validation



- [x] 14.1 Create tests for TUI input model


  - Implement unit tests for InputState struct with all methods (NewInputState, Value, SetValue, Focus, Blur, etc.)
  - Test input state management including focus/blur behavior and value manipulation
  - Test input configuration methods (SetWidth, SetPlaceholder, Clear)
  - Write table-driven tests covering positive and negative scenarios
  - _Requirements: 5.2, 7.3_

- [x] 14.2 Create tests for TUI messages model


  - Implement unit tests for all message types (TodosLoadedMsg, TodoAddedMsg, TodoToggledMsg, etc.)
  - Test message struct initialization and field access
  - Test message type definitions and their usage patterns
  - Write comprehensive tests for all message variants
  - _Requirements: 5.2, 7.3_

- [x] 14.3 Create tests for TUI keyboard handler


  - Implement unit tests for KeyboardHandler function with all mode scenarios
  - Test keyboard input routing for normal, input, edit, and confirmation modes
  - Test key mapping and command generation for all supported keys
  - Test mode transitions triggered by keyboard input
  - Write table-driven tests for each mode's key handling
  - _Requirements: 5.2, 5.3, 5.4, 7.3_

- [x] 14.4 Create tests for TUI operations handler


  - Implement unit tests for OperationsHandler struct and all handler methods
  - Test todo operation handling (HandleTodosLoaded, HandleTodoAdded, HandleTodoToggled, etc.)
  - Test async operation message creation and handling
  - Test error handling and state management during operations
  - Write comprehensive tests for all operation scenarios
  - _Requirements: 5.2, 5.3, 5.4, 7.3_

- [x] 14.5 Create tests for TUI footer view


  - Implement unit tests for FooterView struct with all rendering methods
  - Test help text generation for different modes (normal, input, edit, confirmation)
  - Test compact and expanded footer rendering
  - Test scroll indicator integration and style customization
  - Write table-driven tests for all footer rendering scenarios
  - _Requirements: 5.2, 7.3_

- [x] 14.6 Create tests for TUI header view


  - Implement unit tests for HeaderView struct with all rendering methods
  - Test header rendering with todo count and completion status
  - Test compact header rendering and mode indicators
  - Test completed count calculation and style customization
  - Write comprehensive tests for all header rendering scenarios
  - _Requirements: 5.2, 7.3_

- [x] 14.7 Create tests for TUI list view


  - Implement unit tests for ListView struct with all rendering methods
  - Test todo list rendering including empty state and pagination
  - Test visible range calculation and scroll indicator functionality
  - Test compact list rendering and item height calculations
  - Write table-driven tests for all list rendering scenarios
  - _Requirements: 5.2, 5.3, 7.3_

- [ ] 14.8 Run comprehensive test coverage analysis
  - Run test coverage analysis to ensure 100% coverage for app/ directory
  - Add missing test cases identified by coverage analysis
  - Implement integration tests for end-to-end scenarios
  - Add benchmark tests for performance-critical operations
  - _Requirements: 7.3_

- [ ] 15. Create build and development tooling
  - Add Makefile or build scripts for common development tasks
  - Set up gomock generation automation
  - Create development documentation and usage examples
  - Add linting and formatting configuration
  - _Requirements: 8.1, 8.2, 8.3_