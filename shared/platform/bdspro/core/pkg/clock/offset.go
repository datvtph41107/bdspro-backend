package clock

import "time"

// OffsetClock dùng để simulate lệch giờ
type offsetClock struct {
	base   Clock
	offset time.Duration
}

// NewOffset returns a clock shifted by duration d from base.
func NewOffset(base Clock, d time.Duration) Clock {
	return offsetClock{
		base:   base,
		offset: d,
	}
}

func (o offsetClock) Now() time.Time {
	return o.base.Now().Add(o.offset)
}

func (o offsetClock) Since(t time.Time) time.Duration {
	return o.Now().Sub(t)
}

func (o offsetClock) Until(t time.Time) time.Duration {
	return t.Sub(o.Now())
}
