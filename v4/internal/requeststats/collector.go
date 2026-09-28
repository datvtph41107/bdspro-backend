package requeststats

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	flushDelay     = 10 * time.Second
	bucketDuration = 5 * time.Minute
)

type statsKey struct {
	Method     string
	Path       string
	Service    string
	StatusCode int32
}

type userKey struct {
	ProfileID      uint64
	OrganizationID uint64
}

// TrackEvent is one HTTP request observation from middleware.
type TrackEvent struct {
	Method         string
	Path           string
	Service        string
	StatusCode     int32
	DurationMs     int64
	ProfileID      uint64
	OrganizationID uint64
	IsSpecAudit    bool
}

// Collector buffers request stats/audit in-process and flushes to the local DB
// every ~10s (same behavior as hub-service RequestStatsUsecase).
type Collector struct {
	store *Store

	mu         sync.Mutex
	counts     map[statsKey]int64
	userCounts map[userKey]int64
	audits     []*RequestAuditLog
	jobPending bool

	flushDelay time.Duration
	now        func() time.Time
	afterFunc  func(time.Duration, func()) *time.Timer
}

func NewCollector(store *Store) *Collector {
	return &Collector{
		store:      store,
		counts:     make(map[statsKey]int64),
		userCounts: make(map[userKey]int64),
		audits:     make([]*RequestAuditLog, 0),
		flushDelay: flushDelay,
		now:        time.Now,
		afterFunc:  time.AfterFunc,
	}
}

// Track buffers one request. Spec-audit URLs write individual audit rows;
// other traffic aggregates into stats + user-interaction counts.
func (c *Collector) Track(event TrackEvent) {
	if c == nil {
		return
	}
	method := strings.ToUpper(strings.TrimSpace(event.Method))
	path := strings.TrimSpace(event.Path)
	if method == "" || path == "" {
		return
	}
	service := strings.TrimSpace(event.Service)
	if service == "" {
		service = "unknown"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if event.IsSpecAudit {
		c.audits = append(c.audits, &RequestAuditLog{
			Method:     method,
			Path:       path,
			Service:    service,
			StatusCode: event.StatusCode,
			DurationMs: event.DurationMs,
			Timestamp:  c.now().UTC(),
		})
	} else {
		c.counts[statsKey{
			Method:     method,
			Path:       path,
			Service:    service,
			StatusCode: event.StatusCode,
		}]++
		if event.ProfileID > 0 {
			c.userCounts[userKey{
				ProfileID:      event.ProfileID,
				OrganizationID: event.OrganizationID,
			}]++
		}
	}

	if !c.jobPending {
		c.jobPending = true
		c.afterFunc(c.flushDelay, c.flush)
	}
}

func (c *Collector) flush() {
	c.mu.Lock()
	snapshot := c.counts
	c.counts = make(map[statsKey]int64)
	userSnapshot := c.userCounts
	c.userCounts = make(map[userKey]int64)
	auditSnapshot := c.audits
	c.audits = make([]*RequestAuditLog, 0, len(auditSnapshot))
	c.mu.Unlock()

	now := c.now().UTC()
	bucket := now.Truncate(bucketDuration)
	flushFailed := false

	if len(snapshot) > 0 && c.store != nil {
		logs := make([]*RequestStatsLog, 0, len(snapshot))
		for key, count := range snapshot {
			if count <= 0 {
				continue
			}
			logs = append(logs, &RequestStatsLog{
				Method:     key.Method,
				Path:       key.Path,
				Service:    key.Service,
				StatusCode: key.StatusCode,
				Count:      count,
				Timestamp:  now,
				BucketHour: bucket,
			})
		}
		if len(logs) > 0 {
			if err := c.store.FlushCounts(context.Background(), logs); err != nil {
				log.Printf("request_stats flush error: %v", err)
				flushFailed = true
				c.mu.Lock()
				for key, count := range snapshot {
					c.counts[key] += count
				}
				c.mu.Unlock()
			}
		}
	}

	if len(userSnapshot) > 0 && c.store != nil {
		logs := make([]*UserInteractionLog, 0, len(userSnapshot))
		for key, count := range userSnapshot {
			if count <= 0 || key.ProfileID == 0 {
				continue
			}
			logs = append(logs, &UserInteractionLog{
				ProfileID:      key.ProfileID,
				OrganizationID: key.OrganizationID,
				Count:          count,
				Timestamp:      now,
				BucketHour:     bucket,
			})
		}
		if len(logs) > 0 {
			if err := c.store.FlushUserInteractions(context.Background(), logs); err != nil {
				log.Printf("user_interaction flush error: %v", err)
				flushFailed = true
				c.mu.Lock()
				for key, count := range userSnapshot {
					c.userCounts[key] += count
				}
				c.mu.Unlock()
			}
		}
	}

	if len(auditSnapshot) > 0 && c.store != nil {
		if err := c.store.FlushAuditLogs(context.Background(), auditSnapshot); err != nil {
			log.Printf("request_audit flush error: %v", err)
			flushFailed = true
			c.mu.Lock()
			c.audits = append(auditSnapshot, c.audits...)
			c.mu.Unlock()
		}
	}

	c.mu.Lock()
	c.jobPending = false
	if flushFailed || len(c.counts) > 0 || len(c.userCounts) > 0 || len(c.audits) > 0 {
		c.jobPending = true
		c.afterFunc(c.flushDelay, c.flush)
	}
	c.mu.Unlock()
}
