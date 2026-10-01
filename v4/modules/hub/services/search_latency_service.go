package services

import (
	"context"
	"strings"
	"time"

	"userv4/modules/hub/dto"
	"userv4/modules/hub/repository"
)

const defaultSearchPath = "/v2/tqd/discovery/search"

// SearchLatencyService provides search latency dashboard queries.
type SearchLatencyService interface {
	GetDashboard(ctx context.Context, req dto.GetSearchLatencyDashboardQuery) (*dto.SearchLatencyDashboardResponse, error)
	GetDailyStats(ctx context.Context, req dto.GetSearchLatencyDailyQuery) (*dto.SearchLatencyDailyResponse, error)
}

type LocalSearchLatencyService struct {
	repo repository.SearchLatencyRepository
}

func NewLocalSearchLatencyService(repo repository.SearchLatencyRepository) *LocalSearchLatencyService {
	return &LocalSearchLatencyService{repo: repo}
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
