package application

import (
	"testing"

	"github.com/yanosea/gct/app/domain"
	"go.uber.org/mock/gomock"
)

func TestNewDeleteTodoUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockTodoRepository(ctrl)
	useCase := NewDeleteTodoUseCase(mockRepo)

	if useCase == nil {
		t.Error("NewDeleteTodoUseCase() returned nil")
	}
	if useCase.repository != mockRepo {
		t.Error("NewDeleteTodoUseCase() repository not set correctly")
	}
}

func TestDeleteTodoUseCase_Run(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		setupMock   func(*domain.MockTodoRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "positive testing (delete existing todo)",
			id:   1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().DeleteById(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "positive testing (delete another existing todo)",
			id:   5,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().DeleteById(5).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "negative testing (invalid ID - zero)",
			id:   0,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				// No repository call expected
			},
			wantErr: true,
		},
		{
			name: "negative testing (invalid ID - negative)",
			id:   -1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				// No repository call expected
			},
			wantErr: true,
		},
		{
			name: "negative testing (todo not found)",
			id:   999,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().DeleteById(999).Return(domain.ErrTodoNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrTodoNotFound,
		},
		{
			name: "negative testing (repository error)",
			id:   1,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().DeleteById(1).Return(domain.NewDomainError(
					domain.ErrorTypeFileSystem,
					"failed to delete from file",
					nil,
				))
			},
			wantErr: true,
		},
		{
			name: "negative testing (repository JSON error)",
			id:   2,
			setupMock: func(mockRepo *domain.MockTodoRepository) {
				mockRepo.EXPECT().DeleteById(2).Return(domain.NewDomainError(
					domain.ErrorTypeJSON,
					"failed to encode JSON",
					nil,
				))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := domain.NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			useCase := NewDeleteTodoUseCase(mockRepo)
			err := useCase.Run(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTodoUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("DeleteTodoUseCase.Run() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}