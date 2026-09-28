package requeststats

import (
	"testing"
	"time"
)

func TestCollectorTrackAggregatesAndAudits(t *testing.T) {
	var flushCalls int
	c := &Collector{
		counts:     make(map[statsKey]int64),
		userCounts: make(map[userKey]int64),
		audits:     make([]*RequestAuditLog, 0),
		flushDelay: time.Hour,
		now:        func() time.Time { return time.Date(2026, 9, 25, 10, 0, 0, 0, time.UTC) },
		afterFunc: func(d time.Duration, f func()) *time.Timer {
			flushCalls++
			return time.NewTimer(d)
		},
	}

	c.Track(TrackEvent{
		Method: "get", Path: "/v2/hub/provinces", Service: "hub", StatusCode: 200,
		ProfileID: 7, OrganizationID: 3,
	})
	c.Track(TrackEvent{
		Method: "GET", Path: "/v2/hub/provinces", Service: "hub", StatusCode: 200,
		ProfileID: 7, OrganizationID: 3,
	})
	c.Track(TrackEvent{
		Method: "GET", Path: "/v2/crm/sitemap.xml", Service: "crm", StatusCode: 200,
		DurationMs: 12, IsSpecAudit: true,
	})

	c.mu.Lock()
	defer c.mu.Unlock()
	if got := c.counts[statsKey{Method: "GET", Path: "/v2/hub/provinces", Service: "hub", StatusCode: 200}]; got != 2 {
		t.Fatalf("stats count = %d, want 2", got)
	}
	if got := c.userCounts[userKey{ProfileID: 7, OrganizationID: 3}]; got != 2 {
		t.Fatalf("user count = %d, want 2", got)
	}
	if len(c.audits) != 1 {
		t.Fatalf("audits = %d, want 1", len(c.audits))
	}
	if flushCalls != 1 {
		t.Fatalf("flush scheduled %d times, want 1", flushCalls)
	}
}
