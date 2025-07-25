//go:generate mockgen -source=sort.go -destination=sort_mock.go -package=proxy

package proxy

import (
	"sort"
)

// Sort provides a proxy interface for sort package functions
type Sort interface {
	Slice(x any, less func(i, j int) bool)
}

// SortImpl implements the Sort interface using the standard sort package
type SortImpl struct{}

// NewSort creates a new Sort implementation
func NewSort() Sort {
	return &SortImpl{}
}

func (s *SortImpl) Slice(x any, less func(i, j int) bool) {
	sort.Slice(x, less)
}
