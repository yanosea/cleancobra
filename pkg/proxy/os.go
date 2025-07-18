//go:generate mockgen -source=os.go -destination=os_mock.go -package=proxy

package proxy

import (
	"io/fs"
	"os"
)

// OS provides a proxy interface for os package functions
type OS interface {
	Getenv(key string) string
	UserHomeDir() (string, error)
	MkdirAll(path string, perm fs.FileMode) error
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	Remove(name string) error
	Stat(name string) (fs.FileInfo, error)
	IsNotExist(err error) bool
}

// OSImpl implements the OS interface using the standard os package
type OSImpl struct{}

// NewOS creates a new OS implementation
func NewOS() OS {
	return &OSImpl{}
}

func (o *OSImpl) Getenv(key string) string {
	return os.Getenv(key)
}

func (o *OSImpl) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (o *OSImpl) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o *OSImpl) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (o *OSImpl) Remove(name string) error {
	return os.Remove(name)
}

func (o *OSImpl) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (o *OSImpl) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
