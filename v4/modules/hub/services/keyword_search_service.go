package services

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

const keywordSearchFlushDelay = 20 * time.Second
const keywordSearchMaxLen = 255

type keywordSearchKey struct{ Keyword, Source string }

// KeywordSearchService tracks keyword searches with in-memory batching.
type KeywordSearchService interface {
	Track(keyword, source string, count int64)
	GetStats(ctx context.Context, req dto.GetKeywordSearchStatsQuery) (*dto.GetKeywordSearchStatsResponse, error)
}

type LocalKeywordSearchService struct {
	repo       repository.KeywordSearchRepository
	mu         sync.Mutex
	counts     map[keywordSearchKey]int64
	jobPending bool
	flushDelay time.Duration
	now        func() time.Time
	afterFunc  func(time.Duration, func()) *time.Timer
}

func NewLocalKeywordSearchService(repo repository.KeywordSearchRepository) *LocalKeywordSearchService {
	return &LocalKeywordSearchService{
		repo:       repo,
		counts:     make(map[keywordSearchKey]int64),
		flushDelay: keywordSearchFlushDelay,
		now:        time.Now,
		afterFunc:  time.AfterFunc,
	}
}

func (s *LocalKeywordSearchService) Track(keyword, source string, count int64) {
	keyword = normalizeKw(keyword)
	if keyword == "" {
		return
	}
	source = strings.ToLower(strings.TrimSpace(source))
	if source == "" {
		source = "unknown"
	}
	if count <= 0 {
		count = 1
	}
	key := keywordSearchKey{Keyword: keyword, Source: source}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[key] += count
	if !s.jobPending {
		s.jobPending = true
		s.afterFunc(s.flushDelay, s.flush)
	}
}

func (s *LocalKeywordSearchService) flush() {
	s.mu.Lock()
	snapshot := s.counts
	s.counts = make(map[keywordSearchKey]int64)
	s.mu.Unlock()

	now := s.now().UTC()
	failed := false
	if len(snapshot) > 0 && s.repo != nil {
		logs := make([]*entity.KeywordSearchLog, 0, len(snapshot))
		for k, c := range snapshot {
			if c > 0 && k.Keyword != "" {
				logs = append(logs, &entity.KeywordSearchLog{Keyword: k.Keyword, Source: k.Source, Count: c, Timestamp: now})
			}
		}
		if len(logs) > 0 {
			if err := s.repo.FlushCounts(context.Background(), logs); err != nil {
				log.Printf("keyword_search flush error: %v", err)
				failed = true
				s.mu.Lock()
				for k, c := range snapshot {
					s.counts[k] += c
				}
				s.mu.Unlock()
			}
		}
	}
	s.mu.Lock()
	s.jobPending = false
	if failed || len(s.counts) > 0 {
		s.jobPending = true
		s.afterFunc(s.flushDelay, s.flush)
	}
	s.mu.Unlock()
}

func (s *LocalKeywordSearchService) GetStats(ctx context.Context, req dto.GetKeywordSearchStatsQuery) (*dto.GetKeywordSearchStatsResponse, error) {
	lastDates := req.LastDates
	if lastDates <= 0 {
		lastDates = 3
	}
	if lastDates > 30 {
		lastDates = 30
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	fromTime := time.Now().UTC().Truncate(24 * time.Hour).AddDate(0, 0, -int(lastDates))

	total, err := s.repo.CountFromTime(ctx, fromTime, req.Source)
	if err != nil {
		return nil, err
	}
	top, err := s.repo.GetTopKeywords(ctx, fromTime, req.Source, int(limit))
	if err != nil {
		return nil, err
	}
	if top == nil {
		top = []dto.KeywordSearchItem{}
	}
	return &dto.GetKeywordSearchStatsResponse{
		Data: top, TotalSearches: total, LastDates: lastDates, Limit: limit,
	}, nil
}

func normalizeKw(raw string) string {
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
