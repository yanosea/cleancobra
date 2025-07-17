//go:generate mockgen -source=time.go -destination=time_mock.go -package=proxy

package proxy

import (
	"time"
)

// Time provides a proxy interface for time package functions
type Time interface {
	Now() time.Time
	Parse(layout, value string) (time.Time, error)
	Since(t time.Time) time.Duration
}

// TimeImpl implements the Time interface using the standard time package
type TimeImpl struct{}

// NewTime creates a new Time implementation
func NewTime() Time {
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