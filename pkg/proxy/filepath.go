//go:generate mockgen -source=filepath.go -destination=filepath_mock.go -package=proxy

package proxy

import (
	"path/filepath"
)

// Filepath provides a proxy interface for filepath package functions
type Filepath interface {
	Join(elem ...string) string
	Dir(path string) string
	Base(path string) string
	Abs(path string) (string, error)
	Clean(path string) string
}

// FilepathImpl implements the Filepath interface using the standard filepath package
type FilepathImpl struct{}

// NewFilepath creates a new Filepath implementation
func NewFilepath() Filepath {
	return &FilepathImpl{}
}

func (f *FilepathImpl) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *FilepathImpl) Dir(path string) string {
	return filepath.Dir(path)
}

func (f *FilepathImpl) Base(path string) string {
	return filepath.Base(path)
}

func (f *FilepathImpl) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (f *FilepathImpl) Clean(path string) string {
	return filepath.Clean(path)
}