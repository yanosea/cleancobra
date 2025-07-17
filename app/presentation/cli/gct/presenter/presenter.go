package presenter

import (
	"errors"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/cli/gct/formatter"
	"github.com/yanosea/gct/pkg/proxy"
)

// Formatter interface defines the contract for todo formatters
type Formatter interface {
	Format(todos []domain.Todo) (string, error)
	FormatSingle(todo domain.Todo) (string, error)
}

// TodoPresenter handles presentation logic for todo operations
type TodoPresenter struct {
	jsonFormatter  *formatter.JSONFormatter
	tableFormatter *formatter.TableFormatter
	plainFormatter *formatter.PlainFormatter
	fmtProxy       proxy.Fmt
	osProxy        proxy.OS
}

// NewTodoPresenter creates a new TodoPresenter instance
func NewTodoPresenter(
	jsonFormatter *formatter.JSONFormatter,
	tableFormatter *formatter.TableFormatter,
	plainFormatter *formatter.PlainFormatter,
	fmtProxy proxy.Fmt,
	osProxy proxy.OS,
) *TodoPresenter {
	return &TodoPresenter{
		jsonFormatter:  jsonFormatter,
		tableFormatter: tableFormatter,
		plainFormatter: plainFormatter,
		fmtProxy:       fmtProxy,
		osProxy:        osProxy,
	}
}

// ShowAddSuccess displays success message for adding a todo
func (p *TodoPresenter) ShowAddSuccess(todo *domain.Todo) {
	p.fmtProxy.Printf("Todo added successfully: %s (ID: %d)\n", todo.Description, todo.ID)
}

// ShowListResults displays the list of todos with the specified format
func (p *TodoPresenter) ShowListResults(todos []domain.Todo, format string) error {
	var formatter Formatter
	
	switch format {
	case "json":
		formatter = p.jsonFormatter
	case "table":
		formatter = p.tableFormatter
	case "plain":
		formatter = p.plainFormatter
	default:
		formatter = p.tableFormatter // default to table format
	}
	
	output, err := formatter.Format(todos)
	if err != nil {
		return err
	}
	
	p.fmtProxy.Println(output)
	return nil
}

// ShowToggleSuccess displays success message for toggling a todo
func (p *TodoPresenter) ShowToggleSuccess(todo *domain.Todo) {
	status := "incomplete"
	if todo.Done {
		status = "complete"
	}
	p.fmtProxy.Printf("Todo %d marked as %s: %s\n", todo.ID, status, todo.Description)
}

// ShowDeleteSuccess displays success message for deleting a todo
func (p *TodoPresenter) ShowDeleteSuccess(todoID int) {
	p.fmtProxy.Printf("Todo %d deleted successfully\n", todoID)
}

// ShowError displays user-friendly error messages
func (p *TodoPresenter) ShowError(err error) {
	if err == nil {
		return
	}
	
	var domainErr *domain.DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Type {
		case domain.ErrorTypeNotFound:
			p.fmtProxy.Printf("Error: %s\n", domainErr.Message)
		case domain.ErrorTypeInvalidInput:
			p.fmtProxy.Printf("Error: %s\n", domainErr.Message)
		case domain.ErrorTypeFileSystem:
			p.fmtProxy.Printf("File system error: %s\n", domainErr.Message)
		case domain.ErrorTypeJSON:
			p.fmtProxy.Printf("JSON error: %s\n", domainErr.Message)
		case domain.ErrorTypeConfiguration:
			p.fmtProxy.Printf("Configuration error: %s\n", domainErr.Message)
		default:
			p.fmtProxy.Printf("Unknown error: %s\n", domainErr.Message)
		}
	} else {
		p.fmtProxy.Printf("Error: %s\n", err.Error())
	}
}

// ShowUsageError displays usage error with helpful message
func (p *TodoPresenter) ShowUsageError(message string) {
	p.fmtProxy.Printf("Usage error: %s\n", message)
}

// ShowValidationError displays validation error message
func (p *TodoPresenter) ShowValidationError(message string) {
	p.fmtProxy.Printf("Validation error: %s\n", message)
}