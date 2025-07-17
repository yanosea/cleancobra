package infrastructure

import (
	"errors"
	"os"
	"testing"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/pkg/proxy"
	"go.uber.org/mock/gomock"
)

func TestNewJSONRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	osProxy := proxy.NewMockOS(ctrl)
	jsonProxy := proxy.NewMockJSON(ctrl)
	filePath := "/test/path/todos.json"

	repo := NewJSONRepository(filePath, osProxy, jsonProxy)

	if repo == nil {
		t.Error("NewJSONRepository() returned nil")
	}
	if repo.filePath != filePath {
		t.Errorf("NewJSONRepository() filePath = %v, want %v", repo.filePath, filePath)
	}
}

func TestJSONRepository_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*proxy.MockOS, *proxy.MockJSON)
		want      []domain.Todo
		wantErr   bool
	}{
		{
			name: "positive testing (file does not exist)",
			setupMock: func(osProxy *proxy.MockOS, jsonProxy *proxy.MockJSON) {
				osProxy.EXPECT().Stat("/test/todos.json").Return(nil, os.ErrNotExist)
				osProxy.EXPECT().IsNotExist(os.ErrNotExist).Return(true)
			},
			want:    []domain.Todo{},
			wantErr: false,
		},
		{
			name: "negative testing (file stat error)",
			setupMock: func(osProxy *proxy.MockOS, jsonProxy *proxy.MockJSON) {
				statErr := errors.New("stat error")
				osProxy.EXPECT().Stat("/test/todos.json").Return(nil, statErr)
				osProxy.EXPECT().IsNotExist(statErr).Return(false)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative testing (file open error)",
			setupMock: func(osProxy *proxy.MockOS, jsonProxy *proxy.MockJSON) {
				openErr := errors.New("open error")
				osProxy.EXPECT().Stat("/test/todos.json").Return(nil, nil)
				osProxy.EXPECT().IsNotExist(gomock.Any()).Return(false)
				osProxy.EXPECT().OpenFile("/test/todos.json", os.O_RDONLY, os.FileMode(0644)).Return(nil, openErr)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			osProxy := proxy.NewMockOS(ctrl)
			jsonProxy := proxy.NewMockJSON(ctrl)
			tt.setupMock(osProxy, jsonProxy)

			repo := NewJSONRepository("/test/todos.json", osProxy, jsonProxy)
			got, err := repo.FindAll()

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("FindAll() count = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, todo := range got {
				if i < len(tt.want) {
					if todo.ID != tt.want[i].ID || todo.Description != tt.want[i].Description || todo.Done != tt.want[i].Done {
						t.Errorf("FindAll() todo[%d] = %v, want %v", i, todo, tt.want[i])
					}
				}
			}
		})
	}
}

func TestJSONRepository_Save(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*proxy.MockOS, *proxy.MockJSON)
		inputTodos []domain.Todo
		want       []domain.Todo
		wantErr    bool
	}{
		{
			name: "positive testing (empty input)",
			setupMock: func(osProxy *proxy.MockOS, jsonProxy *proxy.MockJSON) {
				// No mocks needed for empty input
			},
			inputTodos: []domain.Todo{},
			want:       []domain.Todo{},
			wantErr:    false,
		},
		{
			name: "negative testing (update non-existent todo)",
			setupMock: func(osProxy *proxy.MockOS, jsonProxy *proxy.MockJSON) {
				// FindAll call - file doesn't exist
				osProxy.EXPECT().Stat("/test/todos.json").Return(nil, os.ErrNotExist)
				osProxy.EXPECT().IsNotExist(os.ErrNotExist).Return(true)
			},
			inputTodos: []domain.Todo{
				{ID: 999, Description: "Non-existent todo", Done: false},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			osProxy := proxy.NewMockOS(ctrl)
			jsonProxy := proxy.NewMockJSON(ctrl)
			tt.setupMock(osProxy, jsonProxy)

			repo := NewJSONRepository("/test/todos.json", osProxy, jsonProxy)
			got, err := repo.Save(tt.inputTodos...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("Save() count = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, todo := range got {
				if i < len(tt.want) {
					if todo.ID != tt.want[i].ID || todo.Description != tt.want[i].Description || todo.Done != tt.want[i].Done {
						t.Errorf("Save() todo[%d] = %v, want %v", i, todo, tt.want[i])
					}
				}
			}
		})
	}
}

func TestJSONRepository_DeleteById(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*proxy.MockOS, *proxy.MockJSON)
		deleteID  int
		wantErr   bool
	}{
		{
			name: "negative testing (delete non-existent todo)",
			setupMock: func(osProxy *proxy.MockOS, jsonProxy *proxy.MockJSON) {
				// FindAll call - file doesn't exist
				osProxy.EXPECT().Stat("/test/todos.json").Return(nil, os.ErrNotExist)
				osProxy.EXPECT().IsNotExist(os.ErrNotExist).Return(true)
			},
			deleteID: 999,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			osProxy := proxy.NewMockOS(ctrl)
			jsonProxy := proxy.NewMockJSON(ctrl)
			tt.setupMock(osProxy, jsonProxy)

			repo := NewJSONRepository("/test/todos.json", osProxy, jsonProxy)
			err := repo.DeleteById(tt.deleteID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONRepository_getNextID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	osProxy := proxy.NewMockOS(ctrl)
	jsonProxy := proxy.NewMockJSON(ctrl)
	repo := NewJSONRepository("/test/todos.json", osProxy, jsonProxy)

	tests := []struct {
		name  string
		todos []domain.Todo
		want  int
	}{
		{
			name:  "positive testing (empty todos)",
			todos: []domain.Todo{},
			want:  1,
		},
		{
			name: "positive testing (todos with sequential IDs)",
			todos: []domain.Todo{
				{ID: 1, Description: "Todo 1"},
				{ID: 2, Description: "Todo 2"},
			},
			want: 3,
		},
		{
			name: "positive testing (todos with non-sequential IDs)",
			todos: []domain.Todo{
				{ID: 1, Description: "Todo 1"},
				{ID: 5, Description: "Todo 5"},
				{ID: 3, Description: "Todo 3"},
			},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.getNextID(tt.todos)
			if got != tt.want {
				t.Errorf("getNextID() = %v, want %v", got, tt.want)
			}
		})
	}
}