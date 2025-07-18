//go:generate mockgen -source=time.go -destination=time_mock.go -package=proxy

package proxy

import (
	"time"
)

// Time type alias for time.Time
type Time = time.Time

// Duration type alias for time.Duration  
type Duration = time.Duration

// Common time constants
const (
	RFC3339 = time.RFC3339
)

// TimeProxy provides a proxy interface for time package functions
type TimeProxy interface {
	Now() Time
	Parse(layout, value string) (Time, error)
	Since(t Time) Duration
}

// TimeImpl implements the TimeProxy interface using the standard time package
type TimeImpl struct{}

// NewTime creates a new TimeProxy implementation
func NewTime() TimeProxy {
	return &TimeImpl{}
}

func (t *TimeImpl) Now() time.Time {
	return time.Now()
}

func (t *TimeImpl) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func (t *TimeImpl) Since(tm time.Time) time.Duration {
	return time.Since(tm)
}
