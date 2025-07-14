package storage

// FileStorage defines the interface for file operations
type FileStorage interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm int) error
	MkdirAll(path string, perm int) error
	Stat(name string) (FileInfo, error)
	IsNotExist(err error) bool
}

// FileInfo represents file information
type FileInfo interface {
	Name() string
	Size() int64
	IsDir() bool
}
