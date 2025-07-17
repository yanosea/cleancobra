package container

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/infrastructure"
)

func TestNewContainer(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		cleanup func()
		wantErr bool
	}{
		{
			name: "positive testing",
			setup: func() {
				// Set up a temporary directory for testing
				tempDir := t.TempDir()
				os.Setenv("GCT_DATA_FILE", filepath.Join(tempDir, "test_todos.json"))
			},
			cleanup: func() {
				os.Unsetenv("GCT_DATA_FILE")
			},
			wantErr: false,
		},
		{
			name: "positive testing with XDG_DATA_HOME",
			setup: func() {
				// Set up XDG_DATA_HOME for testing
				tempDir := t.TempDir()
				os.Setenv("XDG_DATA_HOME", tempDir)
				os.Unsetenv("GCT_DATA_FILE") // Ensure GCT_DATA_FILE is not set
			},
			cleanup: func() {
				os.Unsetenv("XDG_DATA_HOME")
			},
			wantErr: false,
		},
		{
			name: "positive testing with fallback path",
			setup: func() {
				// Ensure both environment variables are unset to test fallback
				os.Unsetenv("GCT_DATA_FILE")
				os.Unsetenv("XDG_DATA_HOME")
			},
			cleanup: func() {
				// No cleanup needed
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			container, err := NewContainer()

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewContainer() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewContainer() unexpected error: %v", err)
				return
			}

			if container == nil {
				t.Error("NewContainer() returned nil container")
				return
			}

			// Verify container has all required components
			if container.config == nil {
				t.Error("Container config is nil")
			}

			if container.proxies == nil {
				t.Error("Container proxies is nil")
			}

			if container.repository == nil {
				t.Error("Container repository is nil")
			}

			if container.useCases == nil {
				t.Error("Container useCases is nil")
			}
		})
	}
}

func TestContainer_GetUseCases(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			tempDir := t.TempDir()
			os.Setenv("GCT_DATA_FILE", filepath.Join(tempDir, "test_todos.json"))
			defer os.Unsetenv("GCT_DATA_FILE")

			container, err := NewContainer()
			if err != nil {
				t.Fatalf("NewContainer() error: %v", err)
			}

			useCases := container.GetUseCases()

			if useCases == nil {
				t.Error("GetUseCases() returned nil")
				return
			}

			// Verify all use cases are present
			if useCases.AddTodo == nil {
				t.Error("AddTodo use case is nil")
			}

			if useCases.DeleteTodo == nil {
				t.Error("DeleteTodo use case is nil")
			}

			if useCases.ListTodo == nil {
				t.Error("ListTodo use case is nil")
			}

			if useCases.ToggleTodo == nil {
				t.Error("ToggleTodo use case is nil")
			}

			// Verify use cases are not nil (type checking is done at compile time)
			// Since these are concrete types, we just need to verify they're not nil
		})
	}
}

func TestContainer_GetRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			tempDir := t.TempDir()
			os.Setenv("GCT_DATA_FILE", filepath.Join(tempDir, "test_todos.json"))
			defer os.Unsetenv("GCT_DATA_FILE")

			container, err := NewContainer()
			if err != nil {
				t.Fatalf("NewContainer() error: %v", err)
			}

			repository := container.GetRepository()

			if repository == nil {
				t.Error("GetRepository() returned nil")
				return
			}

			// Verify repository implements TodoRepository interface
			if _, ok := repository.(domain.TodoRepository); !ok {
				t.Error("Repository does not implement TodoRepository interface")
			}

			// Verify repository is of correct concrete type
			if _, ok := repository.(*infrastructure.JSONRepository); !ok {
				t.Error("Repository is not of correct concrete type")
			}
		})
	}
}

func TestContainer_GetConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			tempDir := t.TempDir()
			testDataFile := filepath.Join(tempDir, "test_todos.json")
			os.Setenv("GCT_DATA_FILE", testDataFile)
			defer os.Unsetenv("GCT_DATA_FILE")

			container, err := NewContainer()
			if err != nil {
				t.Fatalf("NewContainer() error: %v", err)
			}

			cfg := container.GetConfig()

			if cfg == nil {
				t.Error("GetConfig() returned nil")
				return
			}

			// Config type is verified at compile time since it's a concrete type

			// Verify config has expected data file path
			if cfg.DataFile != testDataFile {
				t.Errorf("Config DataFile = %v, want %v", cfg.DataFile, testDataFile)
			}
		})
	}
}

func TestContainer_GetProxies(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			tempDir := t.TempDir()
			os.Setenv("GCT_DATA_FILE", filepath.Join(tempDir, "test_todos.json"))
			defer os.Unsetenv("GCT_DATA_FILE")

			container, err := NewContainer()
			if err != nil {
				t.Fatalf("NewContainer() error: %v", err)
			}

			proxies := container.GetProxies()

			if proxies == nil {
				t.Error("GetProxies() returned nil")
				return
			}

			// Verify all proxies are present and not nil
			if proxies.OS == nil {
				t.Error("OS proxy is nil")
			}

			if proxies.Filepath == nil {
				t.Error("Filepath proxy is nil")
			}

			if proxies.JSON == nil {
				t.Error("JSON proxy is nil")
			}

			if proxies.Time == nil {
				t.Error("Time proxy is nil")
			}

			if proxies.IO == nil {
				t.Error("IO proxy is nil")
			}

			if proxies.Fmt == nil {
				t.Error("Fmt proxy is nil")
			}

			if proxies.Strings == nil {
				t.Error("Strings proxy is nil")
			}

			if proxies.Strconv == nil {
				t.Error("Strconv proxy is nil")
			}

			if proxies.Cobra == nil {
				t.Error("Cobra proxy is nil")
			}

			if proxies.Bubbletea == nil {
				t.Error("Bubbletea proxy is nil")
			}

			if proxies.Bubbles == nil {
				t.Error("Bubbles proxy is nil")
			}

			if proxies.Lipgloss == nil {
				t.Error("Lipgloss proxy is nil")
			}

			if proxies.Color == nil {
				t.Error("Color proxy is nil")
			}

			if proxies.Envconfig == nil {
				t.Error("Envconfig proxy is nil")
			}
		})
	}
}

func TestContainer_DependencyWiring(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing dependency wiring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			tempDir := t.TempDir()
			os.Setenv("GCT_DATA_FILE", filepath.Join(tempDir, "test_todos.json"))
			defer os.Unsetenv("GCT_DATA_FILE")

			container, err := NewContainer()
			if err != nil {
				t.Fatalf("NewContainer() error: %v", err)
			}

			// Test that use cases can actually work with the wired dependencies
			// This is an integration test to verify the dependency wiring is correct

			// Test AddTodoUseCase
			addUseCase := container.GetUseCases().AddTodo
			todo, err := addUseCase.Run("Test todo")
			if err != nil {
				t.Errorf("AddTodoUseCase.Run() error: %v", err)
			}
			if todo == nil {
				t.Error("AddTodoUseCase.Run() returned nil todo")
			}
			if todo != nil && todo.Description != "Test todo" {
				t.Errorf("AddTodoUseCase.Run() todo description = %v, want %v", todo.Description, "Test todo")
			}

			// Test ListTodoUseCase
			listUseCase := container.GetUseCases().ListTodo
			todos, err := listUseCase.Run()
			if err != nil {
				t.Errorf("ListTodoUseCase.Run() error: %v", err)
			}
			if len(todos) != 1 {
				t.Errorf("ListTodoUseCase.Run() returned %d todos, want 1", len(todos))
			}

			// Test ToggleTodoUseCase
			if len(todos) > 0 {
				toggleUseCase := container.GetUseCases().ToggleTodo
				toggledTodo, err := toggleUseCase.Run(todos[0].ID)
				if err != nil {
					t.Errorf("ToggleTodoUseCase.Run() error: %v", err)
				}
				if toggledTodo == nil {
					t.Error("ToggleTodoUseCase.Run() returned nil todo")
				}
				if toggledTodo != nil && toggledTodo.Done != true {
					t.Errorf("ToggleTodoUseCase.Run() todo.Done = %v, want true", toggledTodo.Done)
				}
			}

			// Test DeleteTodoUseCase
			if len(todos) > 0 {
				deleteUseCase := container.GetUseCases().DeleteTodo
				err := deleteUseCase.Run(todos[0].ID)
				if err != nil {
					t.Errorf("DeleteTodoUseCase.Run() error: %v", err)
				}

				// Verify todo was deleted
				remainingTodos, err := listUseCase.Run()
				if err != nil {
					t.Errorf("ListTodoUseCase.Run() after delete error: %v", err)
				}
				if len(remainingTodos) != 0 {
					t.Errorf("After delete, found %d todos, want 0", len(remainingTodos))
				}
			}
		})
	}
}