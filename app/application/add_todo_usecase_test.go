package application

import (
	"strings"
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"go.uber.org/mock/gomock"
)

func TestNewAddTodoUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockTodoRepository(ctrl)
	useCase := NewAddTodoUseCase(mockRepo)

	if useCase == nil {
		t.Error("NewAddTodoUseCase() returned nil")
	}
	if useCase.repository != mockRepo {
		t.Error("NewAddTodoUseCase() repository not set correctly")
	}
}

func TestAddTodoUseCase_Run(t *testing.T) {
	tests := []struct {
		name        string
		description string
		setupMock   func(*domain.MockTodoRepository)
		wantTodo    *domain.Todo
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "positive testing (valid description)",
			description: "Buy groceries",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				savedTodo := domain.Todo{
					ID:          1,
					Description: "Buy groceries",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{savedTodo}, nil)
			},
			wantTodo: &domain.Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        false,
			},
			wantErr: false,
		},
		{
			name:        "positive testing (description with whitespace)",
			description: "  Clean house  ",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				savedTodo := domain.Todo{
					ID:          1,
					Description: "  Clean house  ",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{savedTodo}, nil)
			},
			wantTodo: &domain.Todo{
				ID:          1,
				Description: "  Clean house  ",
				Done:        false,
			},
			wantErr: false,
		},
		{
			name:        "negative testing (empty description)",
			description: "",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				// No repository call expected
			},
			wantTodo:    nil,
			wantErr:     true,
			expectedErr: domain.ErrEmptyDescription,
		},
		{
			name:        "negative testing (description too long)",
			description: strings.Repeat("a", 501),
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				// No repository call expected
			},
			wantTodo: nil,
			wantErr:  true,
		},
		{
			name:        "negative testing (repository error)",
			description: "Valid description",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().Save(gomock.Any()).Return(nil, domain.ErrTodoNotFound)
			},
			wantTodo:    nil,
			wantErr:     true,
			expectedErr: domain.ErrTodoNotFound,
		},
		{
			name:        "negative testing (repository returns empty result)",
			description: "Valid description",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{}, nil)
			},
			wantTodo: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			useCase := NewAddTodoUseCase(mockRepo)
			got, err := useCase.Run(tt.description)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddTodoUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("AddTodoUseCase.Run() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if got == nil {
				t.Error("AddTodoUseCase.Run() returned nil todo")
				return
			}

			// Basic validation for successful case
			if got.ID <= 0 {
				t.Errorf("AddTodoUseCase.Run() ID should be positive, got %v", got.ID)
			}
			if got.Description != tt.description {
				t.Errorf("AddTodoUseCase.Run() Description = %v, want %v", got.Description, tt.description)
			}
			if got.Done != false {
				t.Errorf("AddTodoUseCase.Run() Done should be false for new todo, got %v", got.Done)
			}
			if got.CreatedAt.IsZero() {
				t.Error("AddTodoUseCase.Run() CreatedAt should not be zero")
			}
			if got.UpdatedAt.IsZero() {
				t.Error("AddTodoUseCase.Run() UpdatedAt should not be zero")
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "positive testing (valid description)",
			description: "Buy groceries",
			wantErr:     false,
		},
		{
			name:        "positive testing (description with whitespace)",
			description: "  Clean house  ",
			wantErr:     false,
		},
		{
			name:        "positive testing (maximum length description)",
			description: strings.Repeat("a", 500),
			wantErr:     false,
		},
		{
			name:        "negative testing (empty description)",
			description: "",
			wantErr:     true,
			expectedErr: domain.ErrEmptyDescription,
		},
		{
			name:        "negative testing (description too long)",
			description: strings.Repeat("a", 501),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("validateDescription() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}