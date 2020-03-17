package clock

import "time"

// Real time. Not fake time.
type Real struct{}

// Now returns the current time.
func (Real) Now() time.Time {
	return time.Now()
}

// New clock for determining the time.
func New() Clock {
	return Real{}
}

// Clock for requesting the time.
type Clock interface {
	Now() time.Time
}
