package storage

import (
	"errors"
	"os"
	"testing"
	"time"

	domainStorage "github.com/yanosea/gct/app/domain/storage"
)

// MockOs implements proxy.Os for testing
type MockOs struct {
	readFileFunc    func(string) ([]byte, error)
	writeFileFunc   func(string, []byte, os.FileMode) error
	mkdirAllFunc    func(string, os.FileMode) error
	statFunc        func(string) (os.FileInfo, error)
	isNotExistFunc  func(error) bool
	exitFunc        func(int)
	getenvFunc      func(string) string
	userHomeDirFunc func() (string, error)
}

func (m *MockOs) ReadFile(filename string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(filename)
	}
	return []byte("test content"), nil
}

func (m *MockOs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.writeFileFunc != nil {
		return m.writeFileFunc(filename, data, perm)
	}
	return nil
}

func (m *MockOs) MkdirAll(path string, perm os.FileMode) error {
	if m.mkdirAllFunc != nil {
		return m.mkdirAllFunc(path, perm)
	}
	return nil
}

func (m *MockOs) Stat(name string) (os.FileInfo, error) {
	if m.statFunc != nil {
		return m.statFunc(name)
	}
	return &mockFileInfo{name: "test.txt", size: 100, isDir: false}, nil
}

func (m *MockOs) IsNotExist(err error) bool {
	if m.isNotExistFunc != nil {
		return m.isNotExistFunc(err)
	}
	return false
}

func (m *MockOs) Exit(code int) {
	if m.exitFunc != nil {
		m.exitFunc(code)
	}
	// Do nothing in tests - we don't want to actually exit
}

func (m *MockOs) Getenv(key string) string {
	if m.getenvFunc != nil {
		return m.getenvFunc(key)
	}
	return ""
}

func (m *MockOs) UserHomeDir() (string, error) {
	if m.userHomeDirFunc != nil {
		return m.userHomeDirFunc()
	}
	return "/home/test", nil
}

// mockFileInfo implements os.FileInfo for testing
type mockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func TestOSFileStorage_ReadFile_Success(t *testing.T) {
	// Arrange
	mockOs := &MockOs{
		readFileFunc: func(filename string) ([]byte, error) {
			if filename == "test.txt" {
				return []byte("test content"), nil
			}
			return nil, errors.New("file not found")
		},
	}
	storage := NewOSFileStorage(mockOs)

	// Act
	content, err := storage.ReadFile("test.txt")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("Expected 'test content', got %s", string(content))
	}
}

func TestOSFileStorage_WriteFile_Success(t *testing.T) {
	// Arrange
	var writtenData []byte
	var writtenPerm os.FileMode
	mockOs := &MockOs{
		writeFileFunc: func(filename string, data []byte, perm os.FileMode) error {
			writtenData = data
			writtenPerm = perm
			return nil
		},
	}
	storage := NewOSFileStorage(mockOs)

	// Act
	err := storage.WriteFile("test.txt", []byte("test data"), 0644)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(writtenData) != "test data" {
		t.Errorf("Expected 'test data', got %s", string(writtenData))
	}
	if writtenPerm != 0644 {
		t.Errorf("Expected permission 0644, got %v", writtenPerm)
	}
}

func TestOSFileStorage_Stat_Success(t *testing.T) {
	// Arrange
	mockOs := &MockOs{
		statFunc: func(name string) (os.FileInfo, error) {
			return &mockFileInfo{name: "test.txt", size: 100, isDir: false}, nil
		},
	}
	storage := NewOSFileStorage(mockOs)

	// Act
	info, err := storage.Stat("test.txt")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if info.Name() != "test.txt" {
		t.Errorf("Expected name 'test.txt', got %s", info.Name())
	}
	if info.Size() != 100 {
		t.Errorf("Expected size 100, got %d", info.Size())
	}
	if info.IsDir() {
		t.Error("Expected file, got directory")
	}
}

func TestOSFileStorage_IsNotExist(t *testing.T) {
	// Arrange
	mockOs := &MockOs{
		isNotExistFunc: func(err error) bool {
			return err.Error() == "file not found"
		},
	}
	storage := NewOSFileStorage(mockOs)

	// Act
	result := storage.IsNotExist(errors.New("file not found"))

	// Assert
	if !result {
		t.Error("Expected true for 'file not found' error")
	}
}

func TestNewOSFileStorage(t *testing.T) {
	// Arrange
	mockOs := &MockOs{}

	// Act
	storage := NewOSFileStorage(mockOs)

	// Assert
	if storage == nil {
		t.Error("Expected storage instance, got nil")
	}

	// Verify it implements the interface
	var _ domainStorage.FileStorage = storage
}
