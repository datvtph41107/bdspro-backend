package clock

import "time"

// Clock là abstraction cho time
type Clock interface {
	Now() time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
}
