package gct

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/cli/gct/formatter"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewToggleCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create real implementations without mocks for positive testing
			cobraProxy := proxy.NewCobra()
			strconvProxy := proxy.NewStrconv()
			
			// Create real formatters and use case with real repository
			jsonProxy := proxy.NewJSON()
			colorProxy := proxy.NewColor()
			stringsProxy := proxy.NewStrings()
			fmtProxy := proxy.NewFmt()
			osProxy := proxy.NewOS()
			
			jsonFormatter := formatter.NewJSONFormatter(jsonProxy)
			tableFormatter := formatter.NewTableFormatter(colorProxy, stringsProxy, fmtProxy)
			plainFormatter := formatter.NewPlainFormatter()
			
			// For this test, we need a real repository implementation
			// Since we don't have one yet, we'll create a simple mock just for this test
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := domain.NewMockTodoRepository(ctrl)
			
			toggleUseCase := application.NewToggleTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, fmtProxy, osProxy)

			// Execute
			result := NewToggleCommand(cobraProxy, strconvProxy, toggleUseCase, presenter)

			// Verify
			if result == nil {
				t.Error("Expected command to be created, got nil")
			}
		})
	}
}

func TestRunToggle(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		setupMock   func(*proxy.MockStrconv, *domain.MockTodoRepository, *proxy.MockFmt)
		wantErr     bool
		expectError string
	}{
		{
			name:  "positive testing (toggle incomplete to complete)",
			idStr: "1",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("1").Return(1, nil)
				existingTodo := domain.Todo{
					ID:          1,
					Description: "Buy groceries",
					Done:        false, // initially incomplete
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				expectedTodo := domain.Todo{
					ID:          1,
					Description: "Buy groceries",
					Done:        true, // toggled to complete
					CreatedAt:   existingTodo.CreatedAt,
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().FindAll().Return([]domain.Todo{existingTodo}, nil)
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{expectedTodo}, nil)
				mockFmt.EXPECT().Printf("Todo %d marked as %s: %s\n", 1, "complete", "Buy groceries")
			},
			wantErr: false,
		},
		{
			name:  "positive testing (toggle complete to incomplete)",
			idStr: "2",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("2").Return(2, nil)
				existingTodo := domain.Todo{
					ID:          2,
					Description: "Complete project",
					Done:        true, // initially complete
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				expectedTodo := domain.Todo{
					ID:          2,
					Description: "Complete project",
					Done:        false, // toggled to incomplete
					CreatedAt:   existingTodo.CreatedAt,
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().FindAll().Return([]domain.Todo{existingTodo}, nil)
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{expectedTodo}, nil)
				mockFmt.EXPECT().Printf("Todo %d marked as %s: %s\n", 2, "incomplete", "Complete project")
			},
			wantErr: false,
		},
		{
			name:  "negative testing (invalid ID string failed)",
			idStr: "abc",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("abc").Return(0, errors.New("strconv.Atoi: parsing \"abc\": invalid syntax"))
				mockFmt.EXPECT().Printf("Validation error: %s\n", "invalid todo ID: must be a number")
			},
			wantErr:     true,
			expectError: "strconv.Atoi: parsing \"abc\": invalid syntax",
		},
		{
			name:  "negative testing (zero ID failed)",
			idStr: "0",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("0").Return(0, nil)
				mockFmt.EXPECT().Printf("Validation error: %s\n", "invalid todo ID: must be positive")
			},
			wantErr: false, // We handle this gracefully without returning error
		},
		{
			name:  "negative testing (negative ID failed)",
			idStr: "-1",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("-1").Return(-1, nil)
				mockFmt.EXPECT().Printf("Validation error: %s\n", "invalid todo ID: must be positive")
			},
			wantErr: false, // We handle this gracefully without returning error
		},
		{
			name:  "negative testing (todo not found failed)",
			idStr: "999",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("999").Return(999, nil)
				mockRepo.EXPECT().FindAll().Return([]domain.Todo{}, nil) // empty list, todo not found
				mockFmt.EXPECT().Printf("Error: %s\n", "todo not found")
			},
			wantErr:     true,
			expectError: "NotFound: todo not found",
		},
		{
			name:  "negative testing (repository FindAll error failed)",
			idStr: "1",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("1").Return(1, nil)
				mockRepo.EXPECT().FindAll().Return(nil, errors.New("repository error"))
				mockFmt.EXPECT().Printf("Error: %s\n", "repository error")
			},
			wantErr:     true,
			expectError: "repository error",
		},
		{
			name:  "negative testing (repository Save error failed)",
			idStr: "1",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("1").Return(1, nil)
				existingTodo := domain.Todo{
					ID:          1,
					Description: "Buy groceries",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().FindAll().Return([]domain.Todo{existingTodo}, nil)
				mockRepo.EXPECT().Save(gomock.Any()).Return(nil, errors.New("save error"))
				mockFmt.EXPECT().Printf("Error: %s\n", "save error")
			},
			wantErr:     true,
			expectError: "save error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mocks
			mockStrconv := proxy.NewMockStrconv(ctrl)
			mockRepository := domain.NewMockTodoRepository(ctrl)
			mockFmt := proxy.NewMockFmt(ctrl)
			mockOS := proxy.NewMockOS(ctrl)

			// Create real formatters and use case
			jsonProxy := proxy.NewJSON()
			colorProxy := proxy.NewColor()
			stringsProxy := proxy.NewStrings()
			
			jsonFormatter := formatter.NewJSONFormatter(jsonProxy)
			tableFormatter := formatter.NewTableFormatter(colorProxy, stringsProxy, proxy.NewFmt())
			plainFormatter := formatter.NewPlainFormatter()
			
			toggleUseCase := application.NewToggleTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, mockFmt, mockOS)

			// Set up mock expectations
			tt.setupMock(mockStrconv, mockRepository, mockFmt)

			// Execute
			err := RunToggle(mockStrconv, toggleUseCase, presenter, tt.idStr)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if tt.expectError != "" && err.Error() != tt.expectError {
					t.Errorf("Expected error message %q, got %q", tt.expectError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}