package gct

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

func TestNewDeleteCommand(t *testing.T) {
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
			
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, fmtProxy, osProxy)

			// Execute
			result := NewDeleteCommand(cobraProxy, strconvProxy, deleteUseCase, presenter)

			// Verify
			if result == nil {
				t.Error("Expected command to be created, got nil")
			}
		})
	}
}

func TestRunDelete(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		setupMock   func(*proxy.MockStrconv, *domain.MockTodoRepository, *proxy.MockFmt)
		wantErr     bool
		expectError string
	}{
		{
			name:  "positive testing",
			idStr: "1",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("1").Return(1, nil)
				mockRepo.EXPECT().DeleteById(1).Return(nil)
				mockFmt.EXPECT().Printf("Todo %d deleted successfully\n", 1)
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
				mockRepo.EXPECT().DeleteById(999).Return(domain.NewDomainError(domain.ErrorTypeNotFound, "todo not found", nil))
				mockFmt.EXPECT().Printf("Error: %s\n", "todo not found")
			},
			wantErr:     true,
			expectError: "NotFound: todo not found",
		},
		{
			name:  "negative testing (repository error failed)",
			idStr: "1",
			setupMock: func(mockStrconv *proxy.MockStrconv, mockRepo *domain.MockTodoRepository, mockFmt *proxy.MockFmt) {
				mockStrconv.EXPECT().Atoi("1").Return(1, nil)
				mockRepo.EXPECT().DeleteById(1).Return(errors.New("repository error"))
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
			
			deleteUseCase := application.NewDeleteTodoUseCase(mockRepository)
			presenter := presenter.NewTodoPresenter(jsonFormatter, tableFormatter, plainFormatter, mockFmt, mockOS)

			// Set up mock expectations
			tt.setupMock(mockStrconv, mockRepository, mockFmt)

			// Execute
			err := RunDelete(mockStrconv, deleteUseCase, presenter, tt.idStr)

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