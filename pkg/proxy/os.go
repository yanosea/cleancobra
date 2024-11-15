package proxy

import (
	"os"
)

type Os interface {
	Exit(code int)
	Getenv(key string) string
	IsNotExist(err error) bool
	MkdirAll(path string, perm os.FileMode) error
	ReadFile(filename string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
	UserHomeDir() (string, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type osProxy struct{}

func NewOs() Os {
	return &osProxy{}
}

func (osProxy) Exit(code int) {
	os.Exit(code)
}

func (osProxy) Getenv(key string) string {
	return os.Getenv(key)
}

func (osProxy) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (osProxy) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (osProxy) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (osProxy) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (osProxy) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (osProxy) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
