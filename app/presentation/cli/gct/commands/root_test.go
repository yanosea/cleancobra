package commands

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/cli/gct/formatter"
	"github.com/yanosea/gct/app/presentation/cli/gct/presenter"
	"github.com/yanosea/gct/pkg/proxy"
)

func TestNewRootCommand(t *testing.T) {
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
			
			listUseCase := application.NewListTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, fmtProxy, osProxy)

			// Execute
			result := NewRootCommand(cobraProxy, listUseCase, presenter)

			// Verify
			if result == nil {
				t.Error("Expected command to be created, got nil")
			}
		})
	}
}

func TestRunRoot(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		setupMock   func(*domain.MockTodoRepository, *proxy.MockFmt)
		wantErr     bool
		expectError string
	}{
		{
			name:   "positive testing",
			format: "table",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				todos := []domain.Todo{
					{ID: 1, Description: "Test todo", Done: false},
				}
				mockRepo.EXPECT().FindAll().Return(todos, nil)
				mockFmt.EXPECT().Println(gomock.Any())
			},
			wantErr: false,
		},
		{
			name:   "positive testing (json format)",
			format: "json",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				todos := []domain.Todo{
					{ID: 1, Description: "Test todo", Done: false},
				}
				mockRepo.EXPECT().FindAll().Return(todos, nil)
				mockFmt.EXPECT().Println(gomock.Any())
			},
			wantErr: false,
		},
		{
			name:   "positive testing (plain format)",
			format: "plain",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				todos := []domain.Todo{
					{ID: 1, Description: "Test todo", Done: false},
				}
				mockRepo.EXPECT().FindAll().Return(todos, nil)
				mockFmt.EXPECT().Println(gomock.Any())
			},
			wantErr: false,
		},
		{
			name:   "negative testing (repository error failed)",
			format: "table",
			setupMock: func(mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockRepo.EXPECT().FindAll().Return(nil, errors.New("repository error"))
				mockFmt.EXPECT().Printf("Error: %s\n", "repository error")
			},
			wantErr:     true,
			expectError: "repository error",
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
			fmtProxy := proxy.NewFmt()
			
			jsonFormatter := formatter.NewJSONFormatter(jsonProxy)
			tableFormatter := formatter.NewTableFormatter(colorProxy, stringsProxy, fmtProxy)
			plainFormatter := formatter.NewPlainFormatter()
			
			listUseCase := application.NewListTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, mockFmt, mockOS)

			// Set up mock expectations
			tt.setupMock(mockRepository, mockFmt)

			// Execute
			err := runRoot(listUseCase, presenter, tt.format)

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