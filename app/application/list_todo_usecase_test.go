package application

import (
	"testing"
	"time"

	"github.com/yanosea/gct/app/domain"
	"go.uber.org/mock/gomock"
)

func TestNewListTodoUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockTodoRepository(ctrl)
	useCase := NewListTodoUseCase(mockRepo)

	if useCase == nil {
		t.Error("NewListTodoUseCase() returned nil")
	}
	if useCase.repository != mockRepo {
		t.Error("NewListTodoUseCase() repository not set correctly")
	}
}

func TestListTodoUseCase_Run(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func(*domain.MockTodoRepository)
		want      []domain.Todo
		wantErr   bool
	}{
		{
			name: "positive testing (empty todo list)",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().FindAll().Return([]domain.Todo{}, nil)
			},
			want:    []domain.Todo{},
			wantErr: false,
		},
		{
			name: "positive testing (single todo)",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				todos := []domain.Todo{
					{
						ID:          1,
						Description: "Buy groceries",
						Done:        false,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}
				mockRepo.EXPECT().FindAll().Return(todos, nil)
			},
			want: []domain.Todo{
				{
					ID:          1,
					Description: "Buy groceries",
					Done:        false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		{
			name: "positive testing (multiple todos)",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
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
						Description: "Clean house",
						Done:        true,
						CreatedAt:   now,
						UpdatedAt:   now.Add(time.Hour),
					},
				}
				mockRepo.EXPECT().FindAll().Return(todos, nil)
			},
			want: []domain.Todo{
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
					UpdatedAt:   now.Add(time.Hour),
				},
			},
			wantErr: false,
		},
		{
			name: "negative testing (repository error)",
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().FindAll().Return(nil, domain.NewDomainError(
					domain.ErrorTypeFileSystem,
					"failed to read file",
					nil,
				))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			useCase := NewListTodoUseCase(mockRepo)
			got, err := useCase.Run()

			if (err != nil) != tt.wantErr {
				t.Errorf("ListTodoUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ListTodoUseCase.Run() count = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, todo := range got {
				if i < len(tt.want) {
					if todo.ID != tt.want[i].ID {
						t.Errorf("ListTodoUseCase.Run() todo[%d].ID = %v, want %v", i, todo.ID, tt.want[i].ID)
					}
					if todo.Description != tt.want[i].Description {
						t.Errorf("ListTodoUseCase.Run() todo[%d].Description = %v, want %v", i, todo.Description, tt.want[i].Description)
					}
					if todo.Done != tt.want[i].Done {
						t.Errorf("ListTodoUseCase.Run() todo[%d].Done = %v, want %v", i, todo.Done, tt.want[i].Done)
					}
					if !todo.CreatedAt.Equal(tt.want[i].CreatedAt) {
						t.Errorf("ListTodoUseCase.Run() todo[%d].CreatedAt = %v, want %v", i, todo.CreatedAt, tt.want[i].CreatedAt)
					}
					if !todo.UpdatedAt.Equal(tt.want[i].UpdatedAt) {
						t.Errorf("ListTodoUseCase.Run() todo[%d].UpdatedAt = %v, want %v", i, todo.UpdatedAt, tt.want[i].UpdatedAt)
					}
				}
			}
		})
	}
}