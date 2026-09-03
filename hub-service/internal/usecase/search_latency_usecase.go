package usecase

import (
	"context"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	_err "common/domain/err"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

const (
	searchLatencyFlushDelay = 10 * time.Second
	// Cap samples per path/method/service trong 1 cửa sổ flush để tránh OOM khi traffic cao.
	searchLatencySampleCap = 20000
)

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

type SearchLatencyUsecase struct {
	repo repo.ISearchLatencyRepo

	mu         sync.Mutex
	aggs       map[searchLatencyKey]*searchLatencyAgg
	jobPending bool

	flushDelay time.Duration
	now        func() time.Time
	afterFunc  func(time.Duration, func()) *time.Timer
}

type ISearchLatencyUsecase interface {
	Track(ctx context.Context, req *dto.TrackSearchLatencyRequest) *_err.ErrorDTO
	GetDashboard(ctx context.Context, req *dto.GetSearchLatencyDashboardRequest) (*dto.GetSearchLatencyDashboardResponse, *_err.ErrorDTO)
	GetDailyStats(ctx context.Context, req *dto.GetSearchLatencyDailyStatsRequest) (*dto.GetSearchLatencyDailyStatsResponse, *_err.ErrorDTO)
}

// @bind: internal/usecase/ISearchLatencyUsecase
func NewSearchLatencyUsecase(repository repo.ISearchLatencyRepo) ISearchLatencyUsecase {
	return &SearchLatencyUsecase{
		repo:       repository,
		aggs:       make(map[searchLatencyKey]*searchLatencyAgg),
		flushDelay: searchLatencyFlushDelay,
		now:        time.Now,
		afterFunc:  time.AfterFunc,
	}
}

func (u *SearchLatencyUsecase) Track(_ context.Context, req *dto.TrackSearchLatencyRequest) *_err.ErrorDTO {
	if u == nil || req == nil {
		return nil
	}
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	path := strings.TrimSpace(req.Path)
	if method == "" || path == "" || req.DurationMs < 0 {
		return nil
	}
	service := strings.TrimSpace(req.Service)
	if service == "" {
		service = "unknown"
	}

	key := searchLatencyKey{Method: method, Path: path, Service: service}

	u.mu.Lock()
	defer u.mu.Unlock()

	agg, ok := u.aggs[key]
	if !ok {
		agg = &searchLatencyAgg{}
		u.aggs[key] = agg
	}
	agg.addSample(req.DurationMs)

	if !u.jobPending {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	return nil
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
	// Reservoir sampling: giữ phân phối mẫu khi vượt cap.
	j := rand.Int63n(a.Count)
	if j < int64(searchLatencySampleCap) {
		a.Samples[j] = durationMs
	}
}

func (a *searchLatencyAgg) merge(other *searchLatencyAgg) {
	if other == nil || other.Count <= 0 {
		return
	}
	if a.Count == 0 {
		*a = *other
		a.Samples = append([]int64(nil), other.Samples...)
		return
	}
	a.Count += other.Count
	a.SumDurationMs += other.SumDurationMs
	if other.MinDurationMs < a.MinDurationMs {
		a.MinDurationMs = other.MinDurationMs
	}
	if other.MaxDurationMs > a.MaxDurationMs {
		a.MaxDurationMs = other.MaxDurationMs
	}
	for _, sample := range other.Samples {
		if len(a.Samples) < searchLatencySampleCap {
			a.Samples = append(a.Samples, sample)
			continue
		}
		j := rand.Int63n(a.Count)
		if j < int64(searchLatencySampleCap) {
			a.Samples[j] = sample
		}
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
	// Nearest-rank: ceil(p/100 * n) - 1
	rank := (p*n + 99) / 100
	if rank < 1 {
		rank = 1
	}
	if rank > n {
		rank = n
	}
	return sorted[rank-1]
}

func (u *SearchLatencyUsecase) flush() {
	u.mu.Lock()
	snapshot := u.aggs
	u.aggs = make(map[searchLatencyKey]*searchLatencyAgg)
	u.mu.Unlock()

	now := u.now().UTC()
	flushFailed := false

	if len(snapshot) > 0 && u.repo != nil {
		logs := make([]*domain.SearchLatencyLog, 0, len(snapshot))
		for key, agg := range snapshot {
			if agg == nil || agg.Count <= 0 {
				continue
			}
			p50, p95, p99 := agg.percentiles()
			logs = append(logs, &domain.SearchLatencyLog{
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
			if err := u.repo.Flush(context.Background(), logs); err != nil {
				log.Printf("search_latency flush error: %v", err)
				flushFailed = true
				u.mu.Lock()
				for key, agg := range snapshot {
					existing, ok := u.aggs[key]
					if !ok {
						u.aggs[key] = agg
						continue
					}
					existing.merge(agg)
				}
				u.mu.Unlock()
			}
		}
	}

	u.mu.Lock()
	u.jobPending = false
	if flushFailed || len(u.aggs) > 0 {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	u.mu.Unlock()
}

func (u *SearchLatencyUsecase) GetDashboard(ctx context.Context, req *dto.GetSearchLatencyDashboardRequest) (*dto.GetSearchLatencyDashboardResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetSearchLatencyDashboardRequest{}
	}
	req.Normalize()

	fromTime := requestStatsFromTime(req.LastDates)
	resp, err := u.repo.GetSummary(ctx, fromTime, req.Path)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi thống kê latency search: " + err.Error()}
	}
	if resp == nil {
		resp = &dto.GetSearchLatencyDashboardResponse{Path: req.Path}
	}
	resp.LastDates = req.LastDates
	resp.Path = req.Path
	return resp, nil
}

func (u *SearchLatencyUsecase) GetDailyStats(ctx context.Context, req *dto.GetSearchLatencyDailyStatsRequest) (*dto.GetSearchLatencyDailyStatsResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetSearchLatencyDailyStatsRequest{}
	}
	req.Normalize()

	fromTime := requestStatsFromTime(req.LastDates)

	var (
		rows []dto.SearchLatencyDailyItem
		err  error
		data []dto.SearchLatencyDailyItem
	)

	if req.Gap == dto.SearchLatencyGapMinutes {
		rows, err = u.repo.GetMinuteStats(ctx, fromTime, req.Path)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy latency theo phút: " + err.Error()}
		}
		data = buildSearchLatencyMinuteSeries(fromTime, rows)
	} else if req.Gap == dto.RequestStatsGapHours {
		rows, err = u.repo.GetHourlyStats(ctx, fromTime, req.Path)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy latency theo giờ: " + err.Error()}
		}
		data = buildSearchLatencyHourlySeries(fromTime, rows)
	} else {
		rows, err = u.repo.GetDailyStats(ctx, fromTime, req.Path)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy latency theo ngày: " + err.Error()}
		}
		data = buildSearchLatencyDailySeries(fromTime, req.LastDates, rows)
	}

	return &dto.GetSearchLatencyDailyStatsResponse{
		Data:      data,
		LastDates: req.LastDates,
		Gap:       req.Gap,
		Path:      req.Path,
	}, nil
}

func buildSearchLatencyDailySeries(fromTime time.Time, lastDates uint32, rows []dto.SearchLatencyDailyItem) []dto.SearchLatencyDailyItem {
	byDate := make(map[string]dto.SearchLatencyDailyItem, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row
	}
	data := make([]dto.SearchLatencyDailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		item := byDate[date]
		data = append(data, copySearchLatencyItem(date, item))
	}
	return data
}

func buildSearchLatencyHourlySeries(fromTime time.Time, rows []dto.SearchLatencyDailyItem) []dto.SearchLatencyDailyItem {
	byHour := make(map[string]dto.SearchLatencyDailyItem, len(rows))
	for _, row := range rows {
		byHour[row.Date] = row
	}

	loc := fromTime.Location()
	now := time.Now().In(loc)
	endHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	startHour := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)

	data := make([]dto.SearchLatencyDailyItem, 0)
	for cursor := startHour; !cursor.After(endHour); cursor = cursor.Add(time.Hour) {
		key := cursor.Format("2006-01-02 15:00")
		item := byHour[key]
		data = append(data, copySearchLatencyItem(key, item))
	}
	return data
}

func buildSearchLatencyMinuteSeries(fromTime time.Time, rows []dto.SearchLatencyDailyItem) []dto.SearchLatencyDailyItem {
	byMinute := make(map[string]dto.SearchLatencyDailyItem, len(rows))
	for _, row := range rows {
		byMinute[row.Date] = row
	}

	loc := fromTime.Location()
	now := time.Now().In(loc)
	endMinute := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, loc)
	startMinute := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)

	data := make([]dto.SearchLatencyDailyItem, 0)
	for cursor := startMinute; !cursor.After(endMinute); cursor = cursor.Add(time.Minute) {
		key := cursor.Format("2006-01-02 15:04")
		item := byMinute[key]
		data = append(data, copySearchLatencyItem(key, item))
	}
	return data
}

func copySearchLatencyItem(date string, item dto.SearchLatencyDailyItem) dto.SearchLatencyDailyItem {
	return dto.SearchLatencyDailyItem{
		Date:          date,
		Count:         item.Count,
		AvgDurationMs: item.AvgDurationMs,
		MinDurationMs: item.MinDurationMs,
		MaxDurationMs: item.MaxDurationMs,
		P50DurationMs: item.P50DurationMs,
		P95DurationMs: item.P95DurationMs,
		P99DurationMs: item.P99DurationMs,
	}
}
