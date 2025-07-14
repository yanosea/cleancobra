package storage

import (
	"os"

	domainStorage "github.com/yanosea/gct/app/domain/storage"
	"github.com/yanosea/gct/pkg/proxy"
)

// OSFileStorage implements FileStorage using the proxy.Os interface
// This provides better testability and consistency with the rest of the application
type OSFileStorage struct {
	osProxy proxy.Os
}

// NewOSFileStorage creates a new OSFileStorage with dependency injection
func NewOSFileStorage(osProxy proxy.Os) domainStorage.FileStorage {
	return &OSFileStorage{
		osProxy: osProxy,
	}
}

func (fs *OSFileStorage) ReadFile(filename string) ([]byte, error) {
	return fs.osProxy.ReadFile(filename)
}

func (fs *OSFileStorage) WriteFile(filename string, data []byte, perm int) error {
	return fs.osProxy.WriteFile(filename, data, os.FileMode(perm))
}

func (fs *OSFileStorage) MkdirAll(path string, perm int) error {
	return fs.osProxy.MkdirAll(path, os.FileMode(perm))
}

func (fs *OSFileStorage) Stat(name string) (domainStorage.FileInfo, error) {
	info, err := fs.osProxy.Stat(name)
	if err != nil {
		return nil, err
	}
	return &osFileInfo{info}, nil
}

func (fs *OSFileStorage) IsNotExist(err error) bool {
	return fs.osProxy.IsNotExist(err)
}

// osFileInfo wraps os.FileInfo to implement domain.FileInfo
type osFileInfo struct {
	info os.FileInfo
}

func (fi *osFileInfo) Name() string {
	return fi.info.Name()
}

func (fi *osFileInfo) Size() int64 {
	return fi.info.Size()
}

func (fi *osFileInfo) IsDir() bool {
	return fi.info.IsDir()
}
