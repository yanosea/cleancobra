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

func TestNewAddCommand(t *testing.T) {
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
			
			addUseCase := application.NewAddTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, fmtProxy, osProxy)

			// Execute
			result := NewAddCommand(cobraProxy, addUseCase, presenter)

			// Verify
			if result == nil {
				t.Error("Expected command to be created, got nil")
			}
		})
	}
}

func TestRunAdd(t *testing.T) {
	tests := []struct {
		name        string
		description string
		setupMock   func(*domain.MockTodoRepository, *proxy.MockFmt)
		wantErr     bool
		expectError string
	}{
		{
			name:        "positive testing",
			description: "Buy groceries",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				expectedTodo := domain.Todo{
					ID:          1,
					Description: "Buy groceries",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{expectedTodo}, nil)
				mockFmt.EXPECT().Printf("Todo added successfully: %s (ID: %d)\n", "Buy groceries", 1)
			},
			wantErr: false,
		},
		{
			name:        "positive testing (long description)",
			description: "This is a very long todo description that should still be valid as long as it doesn't exceed the maximum character limit",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				expectedTodo := domain.Todo{
					ID:          2,
					Description: "This is a very long todo description that should still be valid as long as it doesn't exceed the maximum character limit",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{expectedTodo}, nil)
				mockFmt.EXPECT().Printf("Todo added successfully: %s (ID: %d)\n", expectedTodo.Description, 2)
			},
			wantErr: false,
		},
		{
			name:        "negative testing (empty description failed)",
			description: "",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockFmt.EXPECT().Printf("Error: %s\n", "description cannot be empty")
			},
			wantErr:     true,
			expectError: "InvalidInput: description cannot be empty",
		},
		{
			name:        "negative testing (repository error failed)",
			description: "Valid description",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockRepo.EXPECT().Save(gomock.Any()).Return(nil, errors.New("repository error"))
				mockFmt.EXPECT().Printf("Error: %s\n", "repository error")
			},
			wantErr:     true,
			expectError: "repository error",
		},
		{
			name:        "negative testing (repository returns empty result failed)",
			description: "Valid description",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{}, nil)
				mockFmt.EXPECT().Printf("Configuration error: %s\n", "repository returned empty result")
			},
			wantErr:     true,
			expectError: "Configuration: repository returned empty result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mocks
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
			
			addUseCase := application.NewAddTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, mockFmt, mockOS)

			// Set up mock expectations
			tt.setupMock(mockRepository, mockFmt)

			// Execute
			err := RunAdd(addUseCase, presenter, tt.description)

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