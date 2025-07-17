package presenter

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/cli/gct/formatter"
	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewTodoPresenter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFmt := proxy.NewMockFmt(ctrl)
	mockOS := proxy.NewMockOS(ctrl)
	mockColor := proxy.NewMockColor(ctrl)
	mockStrings := proxy.NewMockStrings(ctrl)
	mockJSON := proxy.NewMockJSON(ctrl)

	jsonFormatter := formatter.NewJSONFormatter(mockJSON)
	tableFormatter := formatter.NewTableFormatter(mockColor, mockStrings, mockFmt)
	plainFormatter := formatter.NewPlainFormatter()

	presenter := NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, mockFmt, mockOS)

	if presenter == nil {
		t.Error("NewTodoPresenter should return a non-nil presenter")
	}
	if presenter.jsonFormatter != jsonFormatter {
		t.Error("jsonFormatter should be set correctly")
	}
	if presenter.tableFormatter != tableFormatter {
		t.Error("tableFormatter should be set correctly")
	}
	if presenter.plainFormatter != plainFormatter {
		t.Error("plainFormatter should be set correctly")
	}
	if presenter.fmtProxy != mockFmt {
		t.Error("fmtProxy should be set correctly")
	}
	if presenter.osProxy != mockOS {
		t.Error("osProxy should be set correctly")
	}
}

func TestTodoPresenter_ShowAddSuccess(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)

			mockFmt.EXPECT().Printf("Todo added successfully: %s (ID: %d)\n", tt.todo.Description, tt.todo.ID).Times(1)

			presenter := &TodoPresenter{
				fmtProxy: mockFmt,
				osProxy:  mockOS,
			}

			presenter.ShowAddSuccess(tt.todo)
		})
	}
}

func TestTodoPresenter_ShowListResults(t *testing.T) {
	now := time.Now()
	todos := []domain.Todo{
		{
			ID:          1,
			Description: "Buy groceries",
			Done:        false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Description: "Walk the dog",
			Done:        true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name       string
		todos      []domain.Todo
		format     string
		setupMocks func(*proxy.MockFmt, *proxy.MockOS, *proxy.MockJSON, *proxy.MockColor, *proxy.MockStrings)
		wantErr    bool
	}{
		{
			name:   "positive testing with json format",
			todos:  todos,
			format: "json",
			setupMocks: func(mockFmt *proxy.MockFmt, mockOS *proxy.MockOS, mockJSON *proxy.MockJSON, mockColor *proxy.MockColor, mockStrings *proxy.MockStrings) {
				mockJSON.EXPECT().Marshal(todos).Return([]byte(`[{"id":1,"description":"Buy groceries","done":false}]`), nil)
				mockFmt.EXPECT().Println(`[{"id":1,"description":"Buy groceries","done":false}]`).Times(1)
			},
			wantErr: false,
		},
		{
			name:   "positive testing with table format",
			todos:  todos,
			format: "table",
			setupMocks: func(mockFmt *proxy.MockFmt, mockOS *proxy.MockOS, mockJSON *proxy.MockJSON, mockColor *proxy.MockColor, mockStrings *proxy.MockStrings) {
				// Mock table formatter calls
				mockColor.EXPECT().Cyan("ID").Return("ID")
				mockColor.EXPECT().Cyan("Status").Return("Status")
				mockColor.EXPECT().Cyan("Description").Return("Description")
				mockFmt.EXPECT().Sprintf("%-4s %-6s %s", "ID", "Status", "Description").Return("ID   Status Description")
				mockFmt.EXPECT().Sprintf("%-4s %-6s %s", "----", "------", "-----------").Return("---- ------ -----------")
				
				// Mock for first todo (not done)
				mockColor.EXPECT().Red("✗ Todo").Return("✗ Todo")
				mockFmt.EXPECT().Sprintf("%-4d %-6s %s", 1, "✗ Todo", "Buy groceries").Return("1    ✗ Todo Buy groceries")
				
				// Mock for second todo (done)
				mockColor.EXPECT().Green("✓ Done").Return("✓ Done")
				mockColor.EXPECT().Green("Walk the dog").Return("Walk the dog")
				mockFmt.EXPECT().Sprintf("%-4d %-6s %s", 2, "✓ Done", "Walk the dog").Return("2    ✓ Done Walk the dog")
				
				mockStrings.EXPECT().Join(gomock.Any(), "\n").Return("formatted table")
				mockFmt.EXPECT().Println("formatted table").Times(1)
			},
			wantErr: false,
		},
		{
			name:   "positive testing with plain format",
			todos:  todos,
			format: "plain",
			setupMocks: func(mockFmt *proxy.MockFmt, mockOS *proxy.MockOS, mockJSON *proxy.MockJSON, mockColor *proxy.MockColor, mockStrings *proxy.MockStrings) {
				mockFmt.EXPECT().Println("[ ] 1: Buy groceries\n[x] 2: Walk the dog").Times(1)
			},
			wantErr: false,
		},
		{
			name:   "positive testing with default format (table)",
			todos:  todos,
			format: "unknown",
			setupMocks: func(mockFmt *proxy.MockFmt, mockOS *proxy.MockOS, mockJSON *proxy.MockJSON, mockColor *proxy.MockColor, mockStrings *proxy.MockStrings) {
				// Mock table formatter calls (same as table format test)
				mockColor.EXPECT().Cyan("ID").Return("ID")
				mockColor.EXPECT().Cyan("Status").Return("Status")
				mockColor.EXPECT().Cyan("Description").Return("Description")
				mockFmt.EXPECT().Sprintf("%-4s %-6s %s", "ID", "Status", "Description").Return("ID   Status Description")
				mockFmt.EXPECT().Sprintf("%-4s %-6s %s", "----", "------", "-----------").Return("---- ------ -----------")
				
				mockColor.EXPECT().Red("✗ Todo").Return("✗ Todo")
				mockFmt.EXPECT().Sprintf("%-4d %-6s %s", 1, "✗ Todo", "Buy groceries").Return("1    ✗ Todo Buy groceries")
				
				mockColor.EXPECT().Green("✓ Done").Return("✓ Done")
				mockColor.EXPECT().Green("Walk the dog").Return("Walk the dog")
				mockFmt.EXPECT().Sprintf("%-4d %-6s %s", 2, "✓ Done", "Walk the dog").Return("2    ✓ Done Walk the dog")
				
				mockStrings.EXPECT().Join(gomock.Any(), "\n").Return("formatted table")
				mockFmt.EXPECT().Println("formatted table").Times(1)
			},
			wantErr: false,
		},
		{
			name:   "negative testing (json format failed)",
			todos:  todos,
			format: "json",
			setupMocks: func(mockFmt *proxy.MockFmt, mockOS *proxy.MockOS, mockJSON *proxy.MockJSON, mockColor *proxy.MockColor, mockStrings *proxy.MockStrings) {
				mockJSON.EXPECT().Marshal(todos).Return(nil, errors.New("marshal error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)
			mockJSON := proxy.NewMockJSON(ctrl)
			mockColor := proxy.NewMockColor(ctrl)
			mockStrings := proxy.NewMockStrings(ctrl)

			tt.setupMocks(mockFmt, mockOS, mockJSON, mockColor, mockStrings)

			jsonFormatter := formatter.NewJSONFormatter(mockJSON)
			tableFormatter := formatter.NewTableFormatter(mockColor, mockStrings, mockFmt)
			plainFormatter := formatter.NewPlainFormatter()

			presenter := NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, mockFmt, mockOS)

			err := presenter.ShowListResults(tt.todos, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("ShowListResults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTodoPresenter_ShowToggleSuccess(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
	}{
		{
			name: "positive testing with completed todo",
			todo: &domain.Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "positive testing with incomplete todo",
			todo: &domain.Todo{
				ID:          2,
				Description: "Walk the dog",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)

			status := "incomplete"
			if tt.todo.Done {
				status = "complete"
			}
			mockFmt.EXPECT().Printf("Todo %d marked as %s: %s\n", tt.todo.ID, status, tt.todo.Description).Times(1)

			presenter := &TodoPresenter{
				fmtProxy: mockFmt,
				osProxy:  mockOS,
			}

			presenter.ShowToggleSuccess(tt.todo)
		})
	}
}

func TestTodoPresenter_ShowDeleteSuccess(t *testing.T) {
	tests := []struct {
		name   string
		todoID int
	}{
		{
			name:   "positive testing",
			todoID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)

			mockFmt.EXPECT().Printf("Todo %d deleted successfully\n", tt.todoID).Times(1)

			presenter := &TodoPresenter{
				fmtProxy: mockFmt,
				osProxy:  mockOS,
			}

			presenter.ShowDeleteSuccess(tt.todoID)
		})
	}
}

func TestTodoPresenter_ShowError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		setupMocks func(*proxy.MockFmt)
	}{
		{
			name: "positive testing with nil error",
			err:  nil,
			setupMocks: func(mockFmt *proxy.MockFmt) {
				// No calls expected for nil error
			},
		},
		{
			name: "positive testing with not found error",
			err:  domain.ErrTodoNotFound,
			setupMocks: func(mockFmt *proxy.MockFmt) {
				mockFmt.EXPECT().Printf("Error: %s\n", "todo not found").Times(1)
			},
		},
		{
			name: "positive testing with invalid input error",
			err:  domain.ErrEmptyDescription,
			setupMocks: func(mockFmt *proxy.MockFmt) {
				mockFmt.EXPECT().Printf("Error: %s\n", "description cannot be empty").Times(1)
			},
		},
		{
			name: "positive testing with file system error",
			err:  domain.NewDomainError(domain.ErrorTypeFileSystem, "file not found", nil),
			setupMocks: func(mockFmt *proxy.MockFmt) {
				mockFmt.EXPECT().Printf("File system error: %s\n", "file not found").Times(1)
			},
		},
		{
			name: "positive testing with JSON error",
			err:  domain.NewDomainError(domain.ErrorTypeJSON, "invalid JSON", nil),
			setupMocks: func(mockFmt *proxy.MockFmt) {
				mockFmt.EXPECT().Printf("JSON error: %s\n", "invalid JSON").Times(1)
			},
		},
		{
			name: "positive testing with configuration error",
			err:  domain.NewDomainError(domain.ErrorTypeConfiguration, "config missing", nil),
			setupMocks: func(mockFmt *proxy.MockFmt) {
				mockFmt.EXPECT().Printf("Configuration error: %s\n", "config missing").Times(1)
			},
		},
		{
			name: "positive testing with non-domain error",
			err:  errors.New("generic error"),
			setupMocks: func(mockFmt *proxy.MockFmt) {
				mockFmt.EXPECT().Printf("Error: %s\n", "generic error").Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)

			tt.setupMocks(mockFmt)

			presenter := &TodoPresenter{
				fmtProxy: mockFmt,
				osProxy:  mockOS,
			}

			presenter.ShowError(tt.err)
		})
	}
}

func TestTodoPresenter_ShowUsageError(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "positive testing",
			message: "missing required argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)

			mockFmt.EXPECT().Printf("Usage error: %s\n", tt.message).Times(1)

			presenter := &TodoPresenter{
				fmtProxy: mockFmt,
				osProxy:  mockOS,
			}

			presenter.ShowUsageError(tt.message)
		})
	}
}

func TestTodoPresenter_ShowValidationError(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "positive testing",
			message: "invalid todo ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)

			mockFmt.EXPECT().Printf("Validation error: %s\n", tt.message).Times(1)

			presenter := &TodoPresenter{
				fmtProxy: mockFmt,
				osProxy:  mockOS,
			}

			presenter.ShowValidationError(tt.message)
		})
	}
}