package domain

import (
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestTodoRepository_FindAll(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockTodoRepository)
		wantTodos   []Todo
		wantErr     bool
		expectedErr error
	}{
		{
			name: "positive testing (empty repository)",
			setupMock: func(mock *MockTodoRepository) {
				mock.EXPECT().FindAll().Return([]Todo{}, nil)
			},
			wantTodos: []Todo{},
			wantErr:   false,
		},
		{
			name: "positive testing (repository with todos)",
			setupMock: func(mock *MockTodoRepository) {
				todos := []Todo{
					{ID: 1, Description: "Todo 1", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: 2, Description: "Todo 2", Done: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				}
				mock.EXPECT().FindAll().Return(todos, nil)
			},
			wantTodos: []Todo{
				{ID: 1, Description: "Todo 1", Done: false},
				{ID: 2, Description: "Todo 2", Done: true},
			},
			wantErr: false,
		},
		{
			name: "negative testing (repository error)",
			setupMock: func(mock *MockTodoRepository) {
				mock.EXPECT().FindAll().Return(nil, ErrTodoNotFound)
			},
			wantTodos:   nil,
			wantErr:     true,
			expectedErr: ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			got, err := mockRepo.FindAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("FindAll() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if len(got) != len(tt.wantTodos) {
				t.Errorf("FindAll() count = %v, want %v", len(got), len(tt.wantTodos))
				return
			}
			for i, todo := range got {
				if i < len(tt.wantTodos) {
					if todo.ID != tt.wantTodos[i].ID || todo.Description != tt.wantTodos[i].Description || todo.Done != tt.wantTodos[i].Done {
						t.Errorf("FindAll() todo[%d] = %v, want %v", i, todo, tt.wantTodos[i])
					}
				}
			}
		})
	}
}

func TestTodoRepository_Save(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockTodoRepository)
		inputTodos  []Todo
		wantTodos   []Todo
		wantErr     bool
		expectedErr error
	}{
		{
			name: "positive testing (save new todo)",
			setupMock: func(mock *MockTodoRepository) {
				inputTodo := Todo{ID: 0, Description: "New todo", Done: false}
				savedTodo := Todo{ID: 1, Description: "New todo", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
				mock.EXPECT().Save(inputTodo).Return([]Todo{savedTodo}, nil)
			},
			inputTodos: []Todo{
				{ID: 0, Description: "New todo", Done: false},
			},
			wantTodos: []Todo{
				{ID: 1, Description: "New todo", Done: false},
			},
			wantErr: false,
		},
		{
			name: "positive testing (save multiple new todos)",
			setupMock: func(mock *MockTodoRepository) {
				inputTodos := []Todo{
					{ID: 0, Description: "Todo 1", Done: false},
					{ID: 0, Description: "Todo 2", Done: false},
				}
				savedTodos := []Todo{
					{ID: 1, Description: "Todo 1", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: 2, Description: "Todo 2", Done: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				}
				mock.EXPECT().Save(inputTodos[0], inputTodos[1]).Return(savedTodos, nil)
			},
			inputTodos: []Todo{
				{ID: 0, Description: "Todo 1", Done: false},
				{ID: 0, Description: "Todo 2", Done: false},
			},
			wantTodos: []Todo{
				{ID: 1, Description: "Todo 1", Done: false},
				{ID: 2, Description: "Todo 2", Done: false},
			},
			wantErr: false,
		},
		{
			name: "positive testing (update existing todo)",
			setupMock: func(mock *MockTodoRepository) {
				inputTodo := Todo{ID: 1, Description: "Updated todo", Done: true}
				savedTodo := Todo{ID: 1, Description: "Updated todo", Done: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
				mock.EXPECT().Save(inputTodo).Return([]Todo{savedTodo}, nil)
			},
			inputTodos: []Todo{
				{ID: 1, Description: "Updated todo", Done: true},
			},
			wantTodos: []Todo{
				{ID: 1, Description: "Updated todo", Done: true},
			},
			wantErr: false,
		},
		{
			name: "negative testing (update non-existent todo failed)",
			setupMock: func(mock *MockTodoRepository) {
				inputTodo := Todo{ID: 999, Description: "Non-existent todo", Done: false}
				mock.EXPECT().Save(inputTodo).Return(nil, ErrTodoNotFound)
			},
			inputTodos: []Todo{
				{ID: 999, Description: "Non-existent todo", Done: false},
			},
			wantTodos:   nil,
			wantErr:     true,
			expectedErr: ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			got, err := mockRepo.Save(tt.inputTodos...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("Save() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.wantTodos) {
					t.Errorf("Save() returned count = %v, want %v", len(got), len(tt.wantTodos))
					return
				}
				for i, todo := range got {
					if i < len(tt.wantTodos) {
						if todo.ID != tt.wantTodos[i].ID || todo.Description != tt.wantTodos[i].Description || todo.Done != tt.wantTodos[i].Done {
							t.Errorf("Save() todo[%d] = %v, want %v", i, todo, tt.wantTodos[i])
						}
					}
				}
			}
		})
	}
}

func TestTodoRepository_DeleteById(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockTodoRepository)
		deleteID    int
		wantErr     bool
		expectedErr error
	}{
		{
			name: "positive testing (delete existing todo)",
			setupMock: func(mock *MockTodoRepository) {
				mock.EXPECT().DeleteById(1).Return(nil)
			},
			deleteID: 1,
			wantErr:  false,
		},
		{
			name: "negative testing (delete non-existent todo failed)",
			setupMock: func(mock *MockTodoRepository) {
				mock.EXPECT().DeleteById(999).Return(ErrTodoNotFound)
			},
			deleteID:    999,
			wantErr:     true,
			expectedErr: ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockTodoRepository(ctrl)
			tt.setupMock(mockRepo)

			err := mockRepo.DeleteById(tt.deleteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("DeleteById() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
		})
	}
}