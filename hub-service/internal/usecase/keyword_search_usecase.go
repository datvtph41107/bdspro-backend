package usecase

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	_err "common/domain/err"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

const keywordSearchFlushDelay = 20 * time.Second
const keywordSearchMaxLen = 255

type keywordSearchKey struct {
	Keyword string
	Source  string
}

type KeywordSearchUsecase struct {
	repo repo.IKeywordSearchRepo

	mu         sync.Mutex
	counts     map[keywordSearchKey]int64
	jobPending bool

	flushDelay time.Duration
	now        func() time.Time
	afterFunc  func(time.Duration, func()) *time.Timer
}

type IKeywordSearchUsecase interface {
	TrackKeywordSearch(ctx context.Context, req *dto.TrackKeywordSearchRequest) *_err.ErrorDTO
	GetStats(ctx context.Context, req *dto.GetKeywordSearchStatsRequest) (*dto.GetKeywordSearchStatsResponse, *_err.ErrorDTO)
}

// @bind: internal/usecase/IKeywordSearchUsecase
func NewKeywordSearchUsecase(repository repo.IKeywordSearchRepo) IKeywordSearchUsecase {
	return &KeywordSearchUsecase{
		repo:       repository,
		counts:     make(map[keywordSearchKey]int64),
		flushDelay: keywordSearchFlushDelay,
		now:        time.Now,
		afterFunc:  time.AfterFunc,
	}
}

func (u *KeywordSearchUsecase) TrackKeywordSearch(_ context.Context, req *dto.TrackKeywordSearchRequest) *_err.ErrorDTO {
	if u == nil || req == nil {
		return nil
	}

	keyword := normalizeKeyword(req.Keyword)
	if keyword == "" {
		return nil
	}

	source := strings.ToLower(strings.TrimSpace(req.Source))
	if source == "" {
		source = "unknown"
	}

	count := req.Count
	if count <= 0 {
		count = 1
	}

	key := keywordSearchKey{Keyword: keyword, Source: source}

	u.mu.Lock()
	defer u.mu.Unlock()

	u.counts[key] += count
	if !u.jobPending {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	return nil
}

func (u *KeywordSearchUsecase) flush() {
	u.mu.Lock()
	snapshot := u.counts
	u.counts = make(map[keywordSearchKey]int64)
	u.mu.Unlock()

	now := u.now().UTC()
	flushFailed := false

	if len(snapshot) > 0 && u.repo != nil {
		logs := make([]*domain.KeywordSearchLog, 0, len(snapshot))
		for key, count := range snapshot {
			if count <= 0 || key.Keyword == "" {
				continue
			}
			logs = append(logs, &domain.KeywordSearchLog{
				Keyword:   key.Keyword,
				Source:    key.Source,
				Count:     count,
				Timestamp: now,
			})
		}
		if len(logs) > 0 {
			if err := u.repo.FlushCounts(context.Background(), logs); err != nil {
				log.Printf("keyword_search flush error: %v", err)
				flushFailed = true
				u.mu.Lock()
				for key, count := range snapshot {
					u.counts[key] += count
				}
				u.mu.Unlock()
			}
		}
	}

	u.mu.Lock()
	u.jobPending = false
	if flushFailed || len(u.counts) > 0 {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	u.mu.Unlock()
}

func (u *KeywordSearchUsecase) GetStats(ctx context.Context, req *dto.GetKeywordSearchStatsRequest) (*dto.GetKeywordSearchStatsResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetKeywordSearchStatsRequest{}
	}
	req.Normalize()

	fromTime := requestStatsFromTime(req.LastDates)

	total, err := u.repo.CountFromTime(ctx, fromTime, req.Source)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi thống kê từ khóa: " + err.Error()}
	}

	top, err := u.repo.GetTopKeywords(ctx, fromTime, req.Source, int(req.Limit))
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy top từ khóa: " + err.Error()}
	}
	if top == nil {
		top = []dto.KeywordSearchItem{}
	}

	return &dto.GetKeywordSearchStatsResponse{
		Data:          top,
		TotalSearches: total,
		LastDates:     req.LastDates,
		Limit:         req.Limit,
	}, nil
}

func normalizeKeyword(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return ""
	}
	s = strings.Join(strings.Fields(s), " ")
	if utf8.RuneCountInString(s) > keywordSearchMaxLen {
		runes := []rune(s)
		s = string(runes[:keywordSearchMaxLen])
	}
	return s
}
