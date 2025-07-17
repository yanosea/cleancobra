package application

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"go.uber.org/mock/gomock"
)

func TestNewToggleTodoUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockTodoRepository(ctrl)
	useCase := NewToggleTodoUseCase(mockRepo)

	if useCase == nil {
		t.Error("NewToggleTodoUseCase() returned nil")
	}
	if useCase.repository != mockRepo {
		t.Error("NewToggleTodoUseCase() repository not set correctly")
	}
}

func TestToggleTodoUseCase_Run(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		id          int
		setupMock   func(*domain.MockTodoRepository)
		wantTodo    *domain.Todo
		wantErr     bool
		expectedErr error
	}{
		{
			name: "positive testing (toggle incomplete todo to complete)",
			id:   1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				existingTodos := []domain.Todo{
					{
						ID:          1,
						Description: "Buy groceries",
						Done:        false,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          2,
						Description: "Clean house",
						Done:        true,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}
				mockRepo.EXPECT().FindAll().Return(existingTodos, nil)

				// Expect Save to be called with toggled todo
				toggledTodo := domain.Todo{
					ID:          1,
					Description: "Buy groceries",
					Done:        true, // Toggled to true
					CreatedAt:   now,
					UpdatedAt:   now.Add(time.Second), // UpdatedAt will be updated
				}
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{toggledTodo}, nil)
			},
			wantTodo: &domain.Todo{
				ID:          1,
				Description: "Buy groceries",
				Done:        true,
			},
			wantErr: false,
		},
		{
			name: "positive testing (toggle complete todo to incomplete)",
			id:   2,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				existingTodos := []domain.Todo{
					{
						ID:          1,
						Description: "Buy groceries",
						Done:        false,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          2,
						Description: "Clean house",
						Done:        true,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}
				mockRepo.EXPECT().FindAll().Return(existingTodos, nil)

				// Expect Save to be called with toggled todo
				toggledTodo := domain.Todo{
					ID:          2,
					Description: "Clean house",
					Done:        false, // Toggled to false
					CreatedAt:   now,
					UpdatedAt:   now.Add(time.Second), // UpdatedAt will be updated
				}
				mockRepo.EXPECT().Save(gomock.Any()).Return([]domain.Todo{toggledTodo}, nil)
			},
			wantTodo: &domain.Todo{
				ID:          2,
				Description: "Clean house",
				Done:        false,
			},
			wantErr: false,
		},
		{
			name: "negative testing (invalid ID - zero)",
			id:   0,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				// No repository calls expected
			},
			wantTodo: nil,
			wantErr:  true,
		},
		{
			name: "negative testing (invalid ID - negative)",
			id:   -1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				// No repository calls expected
			},
			wantTodo: nil,
			wantErr:  true,
		},
		{
			name: "negative testing (todo not found)",
			id:   999,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				existingTodos := []domain.Todo{
					{
						ID:          1,
						Description: "Buy groceries",
						Done:        false,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}
				mockRepo.EXPECT().FindAll().Return(existingTodos, nil)
				// No Save call expected since todo not found
			},
			wantTodo:    nil,
			wantErr:     true,
			expectedErr: domain.ErrTodoNotFound,
		},
		{
			name: "negative testing (repository FindAll error)",
			id:   1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().FindAll().Return(nil, domain.NewDomainError(
					domain.ErrorTypeFileSystem,
					"failed to read file",
					nil,
				))
			},
			wantTodo: nil,
			wantErr:  true,
		},
		{
			name: "negative testing (repository Save error)",
			id:   1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				existingTodos := []domain.Todo{
					{
						ID:          1,
						Description: "Buy groceries",
						Done:        false,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}
				mockRepo.EXPECT().FindAll().Return(existingTodos, nil)
				mockRepo.EXPECT().Save(gomock.Any()).Return(nil, domain.NewDomainError(
					domain.ErrorTypeFileSystem,
					"failed to save file",
					nil,
				))
			},
			wantTodo: nil,
			wantErr:  true,
		},
		{
			name: "negative testing (repository returns empty result)",
			id:   1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				existingTodos := []domain.Todo{
					{
						ID:          1,
						Description: "Buy groceries",
						Done:        false,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}
				mockRepo.EXPECT().FindAll().Return(existingTodos, nil)
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

			useCase := NewToggleTodoUseCase(mockRepo)
			got, err := useCase.Run(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("ToggleTodoUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("ToggleTodoUseCase.Run() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if got == nil {
				t.Error("ToggleTodoUseCase.Run() returned nil todo")
				return
			}

			if got.ID != tt.wantTodo.ID {
				t.Errorf("ToggleTodoUseCase.Run() ID = %v, want %v", got.ID, tt.wantTodo.ID)
			}
			if got.Description != tt.wantTodo.Description {
				t.Errorf("ToggleTodoUseCase.Run() Description = %v, want %v", got.Description, tt.wantTodo.Description)
			}
			if got.Done != tt.wantTodo.Done {
				t.Errorf("ToggleTodoUseCase.Run() Done = %v, want %v", got.Done, tt.wantTodo.Done)
			}
			if got.CreatedAt.IsZero() {
				t.Error("ToggleTodoUseCase.Run() CreatedAt should not be zero")
			}
			if got.UpdatedAt.IsZero() {
				t.Error("ToggleTodoUseCase.Run() UpdatedAt should not be zero")
			}
		})
	}
}