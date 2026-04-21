package clock

import "time"

// FixedClock dùng cho test
type FixedClock struct {
	now time.Time
}

// NewFixed creates a FixedClock at a specific time.
func NewFixed(t time.Time) *FixedClock {
	return &FixedClock{now: t}
}

func (f *FixedClock) Now() time.Time {
	return f.now
}

func (f *FixedClock) Since(t time.Time) time.Duration {
	return f.now.Sub(t)
}

func (f *FixedClock) Until(t time.Time) time.Duration {
	return t.Sub(f.now)
}

// Advance moves the clock forward.
func (f *FixedClock) Advance(d time.Duration) {
	f.now = f.now.Add(d)
}
