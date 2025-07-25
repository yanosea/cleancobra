//go:generate mockgen -source=strings.go -destination=strings_mock.go -package=proxy

package proxy

import (
	"strings"
)

// Strings provides a proxy interface for strings package functions
type Strings interface {
	TrimSpace(s string) string
	Split(s, sep string) []string
	Join(elems []string, sep string) string
	Contains(s, substr string) bool
	HasPrefix(s, prefix string) bool
	HasSuffix(s, suffix string) bool
	ToLower(s string) string
	ToUpper(s string) string
	Repeat(s string, count int) string
	Replace(s, old, new string, n int) string
	ReplaceAll(s, old, new string) string
}

// StringsImpl implements the Strings interface using the standard strings package
type StringsImpl struct{}

// NewStrings creates a new Strings implementation
func NewStrings() Strings {
	return &StringsImpl{}
}

func (s *StringsImpl) TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

func (s *StringsImpl) Split(str, sep string) []string {
	return strings.Split(str, sep)
}

func (s *StringsImpl) Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func (s *StringsImpl) Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func (s *StringsImpl) HasPrefix(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

func (s *StringsImpl) HasSuffix(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

func (s *StringsImpl) ToLower(str string) string {
	return strings.ToLower(str)
}

func (s *StringsImpl) ToUpper(str string) string {
	return strings.ToUpper(str)
}

func (s *StringsImpl) Repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

func (s *StringsImpl) Replace(str, old, new string, n int) string {
	return strings.Replace(str, old, new, n)
}

func (s *StringsImpl) ReplaceAll(str, old, new string) string {
	return strings.ReplaceAll(str, old, new)
}
