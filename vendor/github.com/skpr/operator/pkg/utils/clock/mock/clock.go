package mock

import (
	"time"
)

// Clock which can be manipulated.
type Clock struct {
	Time time.Time
}

// Now returns the manipulated time.
func (c *Clock) Now() time.Time {
	return c.Time
}

// Update manipulates the time.
func (c *Clock) Update(stamp string) error {
	t, err := time.Parse(time.RFC3339, stamp)
	if err != nil {
		return err
	}

	c.Time = t

	return nil
}

// New mock clock.
func New(stamp string) (*Clock, error) {
	t, err := time.Parse(time.RFC3339, stamp)
	if err != nil {
		return nil, err
	}

	return &Clock{Time: t}, nil
}
