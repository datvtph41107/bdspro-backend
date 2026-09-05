package scheduler

import (
	"context"
	"log"
	"time"
	"user/internal/interface/providers"
)

// ZnsScheduler is a process-owned actor. Run blocks until ctx cancellation;
// constructors never start scheduler goroutines.
type ZnsScheduler struct {
	znsProvider providers.IZnsProvider
}

const (
	timeLayout = "2006-01-02 15:04:05"
)

func NewZnsScheduler(znsProvider providers.IZnsProvider) *ZnsScheduler {
	return &ZnsScheduler{znsProvider: znsProvider}
}

func (s *ZnsScheduler) Run(ctx context.Context) {
	log.Println("[ZNS Scheduler] started")
	for {
		next := s.getNextMidnightTime()
		timer := time.NewTimer(time.Until(next))
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			log.Println("[ZNS Scheduler] stopped")
			return
		case <-timer.C:
			s.refreshTokenJob(ctx)
		}
	}
}

func (s *ZnsScheduler) getNextMidnightTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}

func (s *ZnsScheduler) refreshTokenJob(ctx context.Context) {
	jobCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if _, _, err := s.znsProvider.RefreshAccessToken(jobCtx); err != nil {
		log.Printf("[ZNS Scheduler] refresh failed: %v", err)
		return
	}
	// Tokens are credentials. Never log token values or prefixes.
	log.Printf("[ZNS Scheduler] refresh succeeded at %s", time.Now().Format(timeLayout))
}
