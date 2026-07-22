package clock

import "time"

type systemClock struct{}

// New returns a production clock using system time (UTC).
func New() Clock {
	return systemClock{}
}

func (systemClock) Now() time.Time {
	return time.Now().UTC()
}

func (c systemClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c systemClock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}
