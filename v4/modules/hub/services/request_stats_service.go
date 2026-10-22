package services

import (
	"context"

	"mdn4/internal/requeststats"
	"mdn4/modules/hub/dto"
)

// RequestStatsService owns request_stats / audit read+track flows for modules/hub.
type RequestStatsService interface {
	Track(ctx context.Context, req dto.TrackRequestBody)
	GetDashboard(ctx context.Context, req dto.GetRequestStatsDashboardRequest) (*dto.RequestStatsDashboardResponse, error)
	GetDailyStats(ctx context.Context, req dto.GetRequestDailyStatsRequest) (*dto.RequestDailyStatsResponse, error)
	GetUserInteractionDashboard(ctx context.Context, req dto.GetUserInteractionDashboardRequest) (*dto.UserInteractionDashboardResponse, error)
	GetUserInteractionDailyStats(ctx context.Context, req dto.GetUserInteractionDailyStatsRequest) (*dto.UserInteractionDailyStatsResponse, error)
	GetAuditLogs(ctx context.Context, req dto.GetRequestAuditLogsRequest) (*dto.RequestAuditLogsResponse, error)
	Collector() *requeststats.Collector
}

// LocalRequestStatsService reads/writes request stats in the local v4 database.
// @bind: mdn4/modules/hub/services.RequestStatsService
type LocalRequestStatsService struct {
	collector *requeststats.Collector
	store     *requeststats.Store
}

func NewLocalRequestStatsService(collector *requeststats.Collector, store *requeststats.Store) *LocalRequestStatsService {
	return &LocalRequestStatsService{collector: collector, store: store}
}

func (s *LocalRequestStatsService) Collector() *requeststats.Collector {
	if s == nil {
		return nil
	}
	return s.collector
}

func (s *LocalRequestStatsService) Track(_ context.Context, req dto.TrackRequestBody) {
	if s == nil || s.collector == nil {
		return
	}
	s.collector.Track(requeststats.TrackEvent{
		Method:         req.Method,
		Path:           req.Path,
		Service:        req.Service,
		StatusCode:     req.StatusCode,
		DurationMs:     req.DurationMs,
		ProfileID:      req.ProfileID,
		OrganizationID: req.OrganizationID,
		IsSpecAudit:    req.IsSpecAudit,
	})
}

func (s *LocalRequestStatsService) GetDashboard(ctx context.Context, req dto.GetRequestStatsDashboardRequest) (*dto.RequestStatsDashboardResponse, error) {
	req.Normalize()
	from := requeststats.FromTime(req.LastDates)
	total, err := s.store.CountFromTime(ctx, from)
	if err != nil {
		return nil, err
	}
	top, err := s.store.GetTopPaths(ctx, from, requeststats.TopPathLimit)
	if err != nil {
		return nil, err
	}
	byService, err := s.store.GetByService(ctx, from)
	if err != nil {
		return nil, err
	}
	data := make([]dto.RequestTopPathItem, 0, len(top))
	for _, item := range top {
		data = append(data, dto.RequestTopPathItem{
			Method: item.Method, Path: item.Path, Service: item.Service, Count: item.Count,
		})
	}
	services := make([]dto.RequestServiceCountItem, 0, len(byService))
	for _, item := range byService {
		services = append(services, dto.RequestServiceCountItem{Service: item.Service, Count: item.Count})
	}
	return &dto.RequestStatsDashboardResponse{
		TotalRequests: total,
		Data:          data,
		ByService:     services,
		LastDates:     req.LastDates,
	}, nil
}

func (s *LocalRequestStatsService) GetDailyStats(ctx context.Context, req dto.GetRequestDailyStatsRequest) (*dto.RequestDailyStatsResponse, error) {
	req.Normalize()
	from := requeststats.FromTime(req.LastDates)
	var (
		rows []requeststats.DailyItem
		err  error
		data []requeststats.DailyItem
	)
	switch req.Gap {
	case dto.RequestStatsGapMinutes:
		rows, err = s.store.GetMinuteStats(ctx, from)
		if err != nil {
			return nil, err
		}
		data = requeststats.BuildMinuteSeries(from, rows)
	case dto.RequestStatsGapHours:
		rows, err = s.store.GetHourlyStats(ctx, from)
		if err != nil {
			return nil, err
		}
		data = requeststats.BuildHourlySeries(from, rows)
	default:
		rows, err = s.store.GetDailyStats(ctx, from)
		if err != nil {
			return nil, err
		}
		data = requeststats.BuildDailySeries(from, req.LastDates, rows)
	}
	out := make([]dto.RequestDailyItem, 0, len(data))
	for _, item := range data {
		out = append(out, dto.RequestDailyItem{Date: item.Date, Count: item.Count})
	}
	return &dto.RequestDailyStatsResponse{
		Data: out, LastDates: req.LastDates, Gap: req.Gap,
	}, nil
}

func (s *LocalRequestStatsService) GetUserInteractionDashboard(ctx context.Context, req dto.GetUserInteractionDashboardRequest) (*dto.UserInteractionDashboardResponse, error) {
	req.Normalize()
	from := requeststats.FromTime(req.LastDates)
	total, err := s.store.CountUserInteractionsFromTime(ctx, from)
	if err != nil {
		return nil, err
	}
	active, err := s.store.CountActiveUsersFromTime(ctx, from)
	if err != nil {
		return nil, err
	}
	top, err := s.store.GetTopUsers(ctx, from, requeststats.TopUserLimit)
	if err != nil {
		return nil, err
	}
	data := make([]dto.UserInteractionTopItem, 0, len(top))
	for _, item := range top {
		data = append(data, dto.UserInteractionTopItem{
			ProfileID: item.ProfileID, OrganizationID: item.OrganizationID, Count: item.Count,
		})
	}
	return &dto.UserInteractionDashboardResponse{
		TotalInteractions: total,
		ActiveUsers:       active,
		Data:              data,
		LastDates:         req.LastDates,
	}, nil
}

func (s *LocalRequestStatsService) GetUserInteractionDailyStats(ctx context.Context, req dto.GetUserInteractionDailyStatsRequest) (*dto.UserInteractionDailyStatsResponse, error) {
	req.Normalize()
	from := requeststats.FromTime(req.LastDates)
	var (
		rows []requeststats.UserDailyItem
		err  error
		data []requeststats.UserDailyItem
	)
	switch req.Gap {
	case dto.RequestStatsGapMinutes:
		rows, err = s.store.GetUserMinuteStats(ctx, from)
		if err != nil {
			return nil, err
		}
		data = requeststats.BuildUserMinuteSeries(from, rows)
	case dto.RequestStatsGapHours:
		rows, err = s.store.GetUserHourlyStats(ctx, from)
		if err != nil {
			return nil, err
		}
		data = requeststats.BuildUserHourlySeries(from, rows)
	default:
		rows, err = s.store.GetUserDailyStats(ctx, from)
		if err != nil {
			return nil, err
		}
		data = requeststats.BuildUserDailySeries(from, req.LastDates, rows)
	}
	out := make([]dto.UserInteractionDailyItem, 0, len(data))
	for _, item := range data {
		out = append(out, dto.UserInteractionDailyItem{
			Date: item.Date, Interactions: item.Interactions, ActiveUsers: item.ActiveUsers,
		})
	}
	return &dto.UserInteractionDailyStatsResponse{
		Data: out, LastDates: req.LastDates, Gap: req.Gap,
	}, nil
}

func (s *LocalRequestStatsService) GetAuditLogs(ctx context.Context, req dto.GetRequestAuditLogsRequest) (*dto.RequestAuditLogsResponse, error) {
	req.Normalize()
	rows, total, err := s.store.GetAuditLogs(ctx, req.Path, req.Offset(), req.Limit())
	if err != nil {
		return nil, err
	}
	data := make([]dto.RequestAuditLogItem, 0, len(rows))
	for _, item := range rows {
		data = append(data, dto.RequestAuditLogItem{
			ID: item.ID, Method: item.Method, Path: item.Path, Service: item.Service,
			StatusCode: item.StatusCode, DurationMs: item.DurationMs, Timestamp: item.Timestamp,
		})
	}
	return &dto.RequestAuditLogsResponse{
		Data: data, Total: total, Page: req.Page, Size: req.Size,
	}, nil
}
