package services

import (
	"context"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

const (
	defaultSearchPath       = "/v2/tqd/discovery/search"
	searchLatencyFlushDelay = 10 * time.Second
	searchLatencySampleCap  = 20000
)

// SearchLatencyService provides search latency track + dashboard queries.
type SearchLatencyService interface {
	Track(method, path, service string, durationMs int64)
	GetDashboard(ctx context.Context, req dto.GetSearchLatencyDashboardQuery) (*dto.SearchLatencyDashboardResponse, error)
	GetDailyStats(ctx context.Context, req dto.GetSearchLatencyDailyQuery) (*dto.SearchLatencyDailyResponse, error)
}

type searchLatencyKey struct {
	Method  string
	Path    string
	Service string
}

type searchLatencyAgg struct {
	Count         int64
	SumDurationMs int64
	MinDurationMs int64
	MaxDurationMs int64
	Samples       []int64
}

type LocalSearchLatencyService struct {
	repo repository.SearchLatencyRepository

	mu         sync.Mutex
	aggs       map[searchLatencyKey]*searchLatencyAgg
	jobPending bool
	flushDelay time.Duration
	now        func() time.Time
	afterFunc  func(time.Duration, func()) *time.Timer
}

func NewLocalSearchLatencyService(repo repository.SearchLatencyRepository) *LocalSearchLatencyService {
	return &LocalSearchLatencyService{
		repo:       repo,
		aggs:       make(map[searchLatencyKey]*searchLatencyAgg),
		flushDelay: searchLatencyFlushDelay,
		now:        time.Now,
		afterFunc:  time.AfterFunc,
	}
}

func (s *LocalSearchLatencyService) Track(method, path, service string, durationMs int64) {
	if s == nil {
		return
	}
	method = strings.ToUpper(strings.TrimSpace(method))
	path = strings.TrimSpace(path)
	if method == "" || path == "" || durationMs < 0 {
		return
	}
	service = strings.TrimSpace(service)
	if service == "" {
		service = "unknown"
	}
	key := searchLatencyKey{Method: method, Path: path, Service: service}

	s.mu.Lock()
	defer s.mu.Unlock()
	agg, ok := s.aggs[key]
	if !ok {
		agg = &searchLatencyAgg{}
		s.aggs[key] = agg
	}
	agg.addSample(durationMs)
	if !s.jobPending {
		s.jobPending = true
		s.afterFunc(s.flushDelay, s.flush)
	}
}

func (a *searchLatencyAgg) addSample(durationMs int64) {
	if a.Count == 0 {
		a.Count = 1
		a.SumDurationMs = durationMs
		a.MinDurationMs = durationMs
		a.MaxDurationMs = durationMs
		a.Samples = []int64{durationMs}
		return
	}
	a.Count++
	a.SumDurationMs += durationMs
	if durationMs < a.MinDurationMs {
		a.MinDurationMs = durationMs
	}
	if durationMs > a.MaxDurationMs {
		a.MaxDurationMs = durationMs
	}
	if len(a.Samples) < searchLatencySampleCap {
		a.Samples = append(a.Samples, durationMs)
		return
	}
	j := rand.Int63n(a.Count)
	if j < int64(searchLatencySampleCap) {
		a.Samples[j] = durationMs
	}
}

func (a *searchLatencyAgg) percentiles() (p50, p95, p99 int64) {
	if a == nil || len(a.Samples) == 0 {
		return 0, 0, 0
	}
	sorted := append([]int64(nil), a.Samples...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return percentileMs(sorted, 50), percentileMs(sorted, 95), percentileMs(sorted, 99)
}

func percentileMs(sorted []int64, p int) int64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return sorted[0]
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[n-1]
	}
	rank := (p*n + 99) / 100
	if rank < 1 {
		rank = 1
	}
	if rank > n {
		rank = n
	}
	return sorted[rank-1]
}

func (s *LocalSearchLatencyService) flush() {
	s.mu.Lock()
	snapshot := s.aggs
	s.aggs = make(map[searchLatencyKey]*searchLatencyAgg)
	s.mu.Unlock()

	now := s.now().UTC()
	failed := false
	if len(snapshot) > 0 && s.repo != nil {
		logs := make([]*entity.SearchLatencyLog, 0, len(snapshot))
		for key, agg := range snapshot {
			if agg == nil || agg.Count <= 0 {
				continue
			}
			p50, p95, p99 := agg.percentiles()
			logs = append(logs, &entity.SearchLatencyLog{
				Method:        key.Method,
				Path:          key.Path,
				Service:       key.Service,
				Count:         agg.Count,
				SumDurationMs: agg.SumDurationMs,
				MinDurationMs: agg.MinDurationMs,
				MaxDurationMs: agg.MaxDurationMs,
				P50DurationMs: p50,
				P95DurationMs: p95,
				P99DurationMs: p99,
				Timestamp:     now,
			})
		}
		if len(logs) > 0 {
			if err := s.repo.Flush(context.Background(), logs); err != nil {
				log.Printf("search_latency flush error: %v", err)
				failed = true
				s.mu.Lock()
				for key, agg := range snapshot {
					if existing, ok := s.aggs[key]; ok && existing != nil {
						existing.Count += agg.Count
						existing.SumDurationMs += agg.SumDurationMs
						if agg.MinDurationMs < existing.MinDurationMs {
							existing.MinDurationMs = agg.MinDurationMs
						}
						if agg.MaxDurationMs > existing.MaxDurationMs {
							existing.MaxDurationMs = agg.MaxDurationMs
						}
						existing.Samples = append(existing.Samples, agg.Samples...)
					} else {
						s.aggs[key] = agg
					}
				}
				s.mu.Unlock()
			}
		}
	}

	s.mu.Lock()
	s.jobPending = false
	if failed || len(s.aggs) > 0 {
		s.jobPending = true
		s.afterFunc(s.flushDelay, s.flush)
	}
	s.mu.Unlock()
}

func (s *LocalSearchLatencyService) GetDashboard(ctx context.Context, req dto.GetSearchLatencyDashboardQuery) (*dto.SearchLatencyDashboardResponse, error) {
	lastDates, fromTime, path := normalizeLatencyQuery(req.LastDates, req.Path, 30)
	resp, err := s.repo.GetSummary(ctx, fromTime, path)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		resp = &dto.SearchLatencyDashboardResponse{}
	}
	resp.Path = path
	resp.LastDates = lastDates
	return resp, nil
}

func (s *LocalSearchLatencyService) GetDailyStats(ctx context.Context, req dto.GetSearchLatencyDailyQuery) (*dto.SearchLatencyDailyResponse, error) {
	maxDays := 30
	if strings.EqualFold(req.Gap, "minutes") && maxDays > 2 {
		maxDays = 2
	}
	lastDates, fromTime, path := normalizeLatencyQuery(req.LastDates, req.Path, int32(maxDays))

	var items []dto.SearchLatencyDailyItem
	var err error
	switch strings.ToLower(req.Gap) {
	case "hours":
		items, err = s.repo.GetHourlyStats(ctx, fromTime, path)
	case "minutes":
		items, err = s.repo.GetMinuteStats(ctx, fromTime, path)
	default:
		items, err = s.repo.GetDailyStats(ctx, fromTime, path)
		req.Gap = "dates"
	}
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.SearchLatencyDailyItem{}
	}
	return &dto.SearchLatencyDailyResponse{
		Data: items, LastDates: lastDates, Gap: req.Gap, Path: path,
	}, nil
}

func normalizeLatencyQuery(lastDates int32, path string, maxDays int32) (int32, time.Time, string) {
	if lastDates <= 0 {
		lastDates = 3
	}
	if lastDates > maxDays {
		lastDates = maxDays
	}
	if strings.TrimSpace(path) == "" {
		path = defaultSearchPath
	}
	fromTime := time.Now().UTC().Truncate(24 * time.Hour).AddDate(0, 0, -int(lastDates))
	return lastDates, fromTime, path
}
