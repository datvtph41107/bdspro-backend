package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	_errors "common/errors"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/enums"
	"crm/internal/repo"
)

const seoCountFlushDelay = 10 * time.Second

type seoCountKey struct {
	// SeoID > 0 khi track theo seo_domain.id; khi đó bỏ qua RefType/RefID.
	SeoID   uint64
	RefType uint32
	RefID   uint64
}

// SeoCountUsecase buffers public click increments in-process and flushes to DB
// every 10s. Only one flush job is scheduled at a time; a new job is created
// only after the previous job finishes (and only if there are pending counts).
type SeoCountUsecase struct {
	repo          repo.SeoCountRepo
	seoDomainRepo repo.SeoDomainRepo

	mu         sync.Mutex
	counts     map[seoCountKey]int64
	jobPending bool

	flushDelay time.Duration
	now        func() time.Time
	afterFunc  func(time.Duration, func()) *time.Timer
}

func NewSeoCountUsecase(repository repo.SeoCountRepo, seoDomainRepo repo.SeoDomainRepo) *SeoCountUsecase {
	return &SeoCountUsecase{
		repo:          repository,
		seoDomainRepo: seoDomainRepo,
		counts:        make(map[seoCountKey]int64),
		flushDelay:    seoCountFlushDelay,
		now:           time.Now,
		afterFunc:     time.AfterFunc,
	}
}

// IncrementClick bumps the in-memory click counter for (refType, refId) and
// schedules a flush job if none is pending.
func (u *SeoCountUsecase) IncrementClick(_ context.Context, refType uint32, refID uint64) error {
	if u == nil {
		return nil
	}
	if refID == 0 {
		return _errors.ReturnError(400, "refId is required")
	}
	if !enums.ESEORefType(refType).Validate() {
		return _errors.ReturnError(400, "refType không hợp lệ")
	}

	return u.bump(seoCountKey{RefType: refType, RefID: refID})
}

// IncrementClickBySeoID bumps the counter theo seo_domain.id (không cần refType/refId).
func (u *SeoCountUsecase) IncrementClickBySeoID(ctx context.Context, seoID uint64) error {
	if u == nil {
		return nil
	}
	if seoID == 0 {
		return _errors.ReturnError(400, "id is required")
	}
	if u.seoDomainRepo == nil {
		return _errors.ReturnError(500, "seo domain repo is not configured")
	}
	page, err := u.seoDomainRepo.GetByID(ctx, seoID)
	if err != nil {
		return err
	}
	if page == nil || page.ID == 0 {
		return _errors.ReturnError(404, "seo page not found")
	}
	return u.bump(seoCountKey{SeoID: page.ID})
}

func (u *SeoCountUsecase) bump(key seoCountKey) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.counts[key]++
	if !u.jobPending {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	return nil
}

func (u *SeoCountUsecase) flush() {
	u.mu.Lock()
	snapshot := u.counts
	u.counts = make(map[seoCountKey]int64)
	u.mu.Unlock()

	if len(snapshot) > 0 && u.repo != nil {
		now := u.now().UTC()
		logs := make([]*seo_domain.SeoCountLog, 0, len(snapshot))
		for key, clicked := range snapshot {
			if clicked <= 0 {
				continue
			}
			seoID, err := u.resolveSeoID(context.Background(), key)
			if err != nil {
				log.Printf("seo_count resolve seoId=%d refType=%d refId=%d: %v", key.SeoID, key.RefType, key.RefID, err)
				continue
			}
			logs = append(logs, &seo_domain.SeoCountLog{
				SeoID:     seoID,
				Clicked:   clicked,
				Timestamp: now,
			})
		}
		if len(logs) > 0 {
			if err := u.repo.FlushClicks(context.Background(), logs); err != nil {
				log.Printf("seo_count flush error: %v", err)
				u.mu.Lock()
				for key, clicked := range snapshot {
					u.counts[key] += clicked
				}
				u.jobPending = false
				if len(u.counts) > 0 {
					u.jobPending = true
					u.afterFunc(u.flushDelay, u.flush)
				}
				u.mu.Unlock()
				return
			}
		}
	}

	u.mu.Lock()
	u.jobPending = false
	if len(u.counts) > 0 {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	u.mu.Unlock()
}

func (u *SeoCountUsecase) resolveSeoID(ctx context.Context, key seoCountKey) (uint64, error) {
	if key.SeoID > 0 {
		return key.SeoID, nil
	}
	if u.seoDomainRepo == nil {
		return 0, fmt.Errorf("seo domain repo is not configured")
	}
	page, err := u.seoDomainRepo.GetByRef(ctx, key.RefType, key.RefID)
	if err != nil {
		return 0, err
	}
	if page == nil || page.ID == 0 {
		return 0, fmt.Errorf("seo page not found")
	}
	return page.ID, nil
}
