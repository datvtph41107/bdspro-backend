package usecase

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	_err "common/domain/err"
	"hub/internal/domain"
	"hub/internal/dto"
	iprovider "hub/internal/interface"
	"hub/internal/repo"
)

const requestStatsFlushDelay = 10 * time.Second
const requestStatsBucketDuration = 5 * time.Minute
const requestStatsTopPathLimit = 30
const userInteractionTopLimit = 30

type requestStatsKey struct {
	Method     string
	Path       string
	Service    string
	StatusCode int32
}

type userInteractionKey struct {
	ProfileID      uint64
	OrganizationID uint64
}

// RequestStatsUsecase buffers gateway request increments in-process and flushes
// to DB every 10s. Flush upserts theo bucket 5 phút: cùng method/path/service/status
// trong khung 5 phút hiện tại thì tăng count và cập nhật timestamp = lần flush gần
// nhất; bucket mới thì INSERT. Spec URLs (is_spec_audit) vẫn ghi từng lần (audit).
type RequestStatsUsecase struct {
	repo       repo.IRequestStatsRepo
	userClient iprovider.IUserClient

	mu         sync.Mutex
	counts     map[requestStatsKey]int64
	userCounts map[userInteractionKey]int64
	audits     []*domain.RequestAuditLog
	jobPending bool

	flushDelay time.Duration
	now        func() time.Time
	afterFunc  func(time.Duration, func()) *time.Timer
}

type IRequestStatsUsecase interface {
	TrackRequest(ctx context.Context, req *dto.TrackRequestRequest) *_err.ErrorDTO
	GetDashboard(ctx context.Context, req *dto.GetRequestStatsDashboardRequest) (*dto.GetRequestStatsDashboardResponse, *_err.ErrorDTO)
	GetDailyStats(ctx context.Context, req *dto.GetRequestDailyStatsRequest) (*dto.GetRequestDailyStatsResponse, *_err.ErrorDTO)
	GetUserInteractionDashboard(ctx context.Context, req *dto.GetUserInteractionDashboardRequest) (*dto.GetUserInteractionDashboardResponse, *_err.ErrorDTO)
	GetUserInteractionDailyStats(ctx context.Context, req *dto.GetUserInteractionDailyStatsRequest) (*dto.GetUserInteractionDailyStatsResponse, *_err.ErrorDTO)
	GetAuditLogs(ctx context.Context, req *dto.GetRequestAuditLogsRequest) (*dto.GetRequestAuditLogsResponse, *_err.ErrorDTO)
}

// @bind: internal/usecase/IRequestStatsUsecase
func NewRequestStatsUsecase(repository repo.IRequestStatsRepo, userClient iprovider.IUserClient) IRequestStatsUsecase {
	return &RequestStatsUsecase{
		repo:       repository,
		userClient: userClient,
		counts:     make(map[requestStatsKey]int64),
		userCounts: make(map[userInteractionKey]int64),
		audits:     make([]*domain.RequestAuditLog, 0),
		flushDelay: requestStatsFlushDelay,
		now:        time.Now,
		afterFunc:  time.AfterFunc,
	}
}

func (u *RequestStatsUsecase) TrackRequest(_ context.Context, req *dto.TrackRequestRequest) *_err.ErrorDTO {
	if u == nil || req == nil {
		return nil
	}

	method := strings.ToUpper(strings.TrimSpace(req.Method))
	path := strings.TrimSpace(req.Path)
	if method == "" || path == "" {
		return nil
	}

	service := strings.TrimSpace(req.Service)
	if service == "" {
		service = "unknown"
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	// Spec URLs: ghi từng lần gọi (audit), không gộp count; flush sau ~10s.
	if req.IsSpecAudit {
		u.audits = append(u.audits, &domain.RequestAuditLog{
			Method:     method,
			Path:       path,
			Service:    service,
			StatusCode: req.StatusCode,
			DurationMs: req.DurationMs,
			Timestamp:  u.now().UTC(),
		})
	} else {
		key := requestStatsKey{
			Method:     method,
			Path:       path,
			Service:    service,
			StatusCode: req.StatusCode,
		}
		u.counts[key]++
		if req.ProfileID > 0 {
			u.userCounts[userInteractionKey{
				ProfileID:      req.ProfileID,
				OrganizationID: req.OrganizationID,
			}]++
		}
	}

	if !u.jobPending {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	return nil
}

func (u *RequestStatsUsecase) flush() {
	u.mu.Lock()
	snapshot := u.counts
	u.counts = make(map[requestStatsKey]int64)
	userSnapshot := u.userCounts
	u.userCounts = make(map[userInteractionKey]int64)
	auditSnapshot := u.audits
	u.audits = make([]*domain.RequestAuditLog, 0, len(auditSnapshot))
	u.mu.Unlock()

	// timestamp = lần flush gần nhất; bucket_hour = khóa 5 phút (chỉ INSERT khi bucket mới).
	now := u.now().UTC()
	bucketHour := now.Truncate(requestStatsBucketDuration)
	flushFailed := false

	if len(snapshot) > 0 && u.repo != nil {
		logs := make([]*domain.RequestStatsLog, 0, len(snapshot))
		for key, count := range snapshot {
			if count <= 0 {
				continue
			}
			logs = append(logs, &domain.RequestStatsLog{
				Method:     key.Method,
				Path:       key.Path,
				Service:    key.Service,
				StatusCode: key.StatusCode,
				Count:      count,
				Timestamp:  now,
				BucketHour: bucketHour,
			})
		}
		if len(logs) > 0 {
			if err := u.repo.FlushCounts(context.Background(), logs); err != nil {
				log.Printf("request_stats flush error: %v", err)
				flushFailed = true
				u.mu.Lock()
				for key, count := range snapshot {
					u.counts[key] += count
				}
				u.mu.Unlock()
			}
		}
	}

	if len(userSnapshot) > 0 && u.repo != nil {
		logs := make([]*domain.UserInteractionLog, 0, len(userSnapshot))
		for key, count := range userSnapshot {
			if count <= 0 || key.ProfileID == 0 {
				continue
			}
			logs = append(logs, &domain.UserInteractionLog{
				ProfileID:      key.ProfileID,
				OrganizationID: key.OrganizationID,
				Count:          count,
				Timestamp:      now,
				BucketHour:     bucketHour,
			})
		}
		if len(logs) > 0 {
			if err := u.repo.FlushUserInteractions(context.Background(), logs); err != nil {
				log.Printf("user_interaction flush error: %v", err)
				flushFailed = true
				u.mu.Lock()
				for key, count := range userSnapshot {
					u.userCounts[key] += count
				}
				u.mu.Unlock()
			}
		}
	}

	if len(auditSnapshot) > 0 && u.repo != nil {
		if err := u.repo.FlushAuditLogs(context.Background(), auditSnapshot); err != nil {
			log.Printf("request_audit flush error: %v", err)
			flushFailed = true
			u.mu.Lock()
			u.audits = append(auditSnapshot, u.audits...)
			u.mu.Unlock()
		}
	}

	u.mu.Lock()
	u.jobPending = false
	if flushFailed || len(u.counts) > 0 || len(u.userCounts) > 0 || len(u.audits) > 0 {
		u.jobPending = true
		u.afterFunc(u.flushDelay, u.flush)
	}
	u.mu.Unlock()
}

func (u *RequestStatsUsecase) GetDashboard(ctx context.Context, req *dto.GetRequestStatsDashboardRequest) (*dto.GetRequestStatsDashboardResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetRequestStatsDashboardRequest{}
	}
	req.Normalize()

	fromTime := requestStatsFromTime(req.LastDates)

	total, err := u.repo.CountFromTime(ctx, fromTime)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi thống kê request: " + err.Error()}
	}

	topPaths, err := u.repo.GetTopPaths(ctx, fromTime, requestStatsTopPathLimit)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy top path: " + err.Error()}
	}
	if topPaths == nil {
		topPaths = []dto.RequestTopPathItem{}
	}

	byService, err := u.repo.GetByService(ctx, fromTime)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi thống kê theo service: " + err.Error()}
	}
	if byService == nil {
		byService = []dto.RequestServiceCountItem{}
	}

	return &dto.GetRequestStatsDashboardResponse{
		TotalRequests: total,
		Data:          topPaths,
		ByService:     byService,
		LastDates:     req.LastDates,
	}, nil
}

func (u *RequestStatsUsecase) GetDailyStats(ctx context.Context, req *dto.GetRequestDailyStatsRequest) (*dto.GetRequestDailyStatsResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetRequestDailyStatsRequest{}
	}
	req.Normalize()

	fromTime := requestStatsFromTime(req.LastDates)

	var (
		rows []dto.RequestDailyItem
		err  error
		data []dto.RequestDailyItem
	)

	if req.Gap == dto.RequestStatsGapMinutes {
		rows, err = u.repo.GetMinuteStats(ctx, fromTime)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy thống kê theo phút: " + err.Error()}
		}
		data = buildRequestMinuteSeries(fromTime, rows)
	} else if req.Gap == dto.RequestStatsGapHours {
		rows, err = u.repo.GetHourlyStats(ctx, fromTime)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy thống kê theo giờ: " + err.Error()}
		}
		data = buildRequestHourlySeries(fromTime, rows)
	} else {
		rows, err = u.repo.GetDailyStats(ctx, fromTime)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy thống kê theo ngày: " + err.Error()}
		}
		data = buildRequestDailySeries(fromTime, req.LastDates, rows)
	}

	return &dto.GetRequestDailyStatsResponse{
		Data:      data,
		LastDates: req.LastDates,
		Gap:       req.Gap,
	}, nil
}

func (u *RequestStatsUsecase) GetUserInteractionDashboard(ctx context.Context, req *dto.GetUserInteractionDashboardRequest) (*dto.GetUserInteractionDashboardResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetUserInteractionDashboardRequest{}
	}
	req.Normalize()

	fromTime := requestStatsFromTime(req.LastDates)

	total, err := u.repo.CountUserInteractionsFromTime(ctx, fromTime)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi thống kê tương tác user: " + err.Error()}
	}

	activeUsers, err := u.repo.CountActiveUsersFromTime(ctx, fromTime)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi đếm active users: " + err.Error()}
	}

	topUsers, err := u.repo.GetTopUsers(ctx, fromTime, userInteractionTopLimit)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy top users: " + err.Error()}
	}
	if topUsers == nil {
		topUsers = []dto.UserInteractionTopItem{}
	}

	u.enrichUserInteractionProfiles(ctx, topUsers)

	return &dto.GetUserInteractionDashboardResponse{
		TotalInteractions: total,
		ActiveUsers:       activeUsers,
		Data:              topUsers,
		LastDates:         req.LastDates,
	}, nil
}

func (u *RequestStatsUsecase) GetUserInteractionDailyStats(ctx context.Context, req *dto.GetUserInteractionDailyStatsRequest) (*dto.GetUserInteractionDailyStatsResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetUserInteractionDailyStatsRequest{}
	}
	req.Normalize()

	fromTime := requestStatsFromTime(req.LastDates)

	var (
		rows []dto.UserInteractionDailyItem
		err  error
		data []dto.UserInteractionDailyItem
	)

	if req.Gap == dto.RequestStatsGapMinutes {
		rows, err = u.repo.GetUserMinuteStats(ctx, fromTime)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy thống kê user theo phút: " + err.Error()}
		}
		data = buildUserInteractionMinuteSeries(fromTime, rows)
	} else if req.Gap == dto.RequestStatsGapHours {
		rows, err = u.repo.GetUserHourlyStats(ctx, fromTime)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy thống kê user theo giờ: " + err.Error()}
		}
		data = buildUserInteractionHourlySeries(fromTime, rows)
	} else {
		rows, err = u.repo.GetUserDailyStats(ctx, fromTime)
		if err != nil {
			return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy thống kê user theo ngày: " + err.Error()}
		}
		data = buildUserInteractionDailySeries(fromTime, req.LastDates, rows)
	}

	return &dto.GetUserInteractionDailyStatsResponse{
		Data:      data,
		LastDates: req.LastDates,
		Gap:       req.Gap,
	}, nil
}

func (u *RequestStatsUsecase) GetAuditLogs(ctx context.Context, req *dto.GetRequestAuditLogsRequest) (*dto.GetRequestAuditLogsResponse, *_err.ErrorDTO) {
	if req == nil {
		req = &dto.GetRequestAuditLogsRequest{}
	}
	req.Normalize()

	data, total, err := u.repo.GetAuditLogs(ctx, req.Path, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: "Lỗi khi lấy audit log: " + err.Error()}
	}
	if data == nil {
		data = []dto.RequestAuditLogItem{}
	}

	return &dto.GetRequestAuditLogsResponse{
		Data:  data,
		Total: total,
		Page:  req.GetPage(),
		Size:  uint32(req.GetLimit()),
	}, nil
}

func requestStatsFromTime(lastDates uint32) time.Time {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)
	}
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return today.AddDate(0, 0, -int(lastDates)+1)
}

func (u *RequestStatsUsecase) enrichUserInteractionProfiles(ctx context.Context, items []dto.UserInteractionTopItem) {
	if u == nil || u.userClient == nil || len(items) == 0 {
		return
	}

	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		if item.ProfileID > 0 {
			ids = append(ids, item.ProfileID)
		}
	}
	if len(ids) == 0 {
		return
	}

	profiles, err := u.userClient.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Printf("request_stats enrich profiles error: %v", err)
		return
	}
	if len(profiles) == 0 {
		return
	}

	for i := range items {
		p, ok := profiles[items[i].ProfileID]
		if !ok || p == nil {
			continue
		}
		items[i].FullName = p.GetFullName()
		items[i].Avatar = p.GetAvatar()
		items[i].Phone = p.GetPhone()
		items[i].TickVerified = p.GetTickVerified()
	}
}

func buildRequestDailySeries(fromTime time.Time, lastDates uint32, rows []dto.RequestDailyItem) []dto.RequestDailyItem {
	byDate := make(map[string]uint64, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row.Count
	}

	data := make([]dto.RequestDailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		data = append(data, dto.RequestDailyItem{
			Date:  date,
			Count: byDate[date],
		})
	}
	return data
}

func buildRequestHourlySeries(fromTime time.Time, rows []dto.RequestDailyItem) []dto.RequestDailyItem {
	byHour := make(map[string]uint64, len(rows))
	for _, row := range rows {
		byHour[row.Date] = row.Count
	}

	loc := fromTime.Location()
	now := time.Now().In(loc)
	endHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	startHour := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)

	hours := int(endHour.Sub(startHour).Hours()) + 1
	if hours < 1 {
		hours = 1
	}

	data := make([]dto.RequestDailyItem, 0, hours)
	for cursor := startHour; !cursor.After(endHour); cursor = cursor.Add(time.Hour) {
		key := cursor.Format("2006-01-02 15:00")
		data = append(data, dto.RequestDailyItem{
			Date:  key,
			Count: byHour[key],
		})
	}
	return data
}

func buildRequestMinuteSeries(fromTime time.Time, rows []dto.RequestDailyItem) []dto.RequestDailyItem {
	byMinute := make(map[string]uint64, len(rows))
	for _, row := range rows {
		byMinute[row.Date] = row.Count
	}

	loc := fromTime.Location()
	now := time.Now().In(loc)
	endMinute := now.Truncate(requestStatsBucketDuration)
	startMinute := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)

	data := make([]dto.RequestDailyItem, 0)
	for cursor := startMinute; !cursor.After(endMinute); cursor = cursor.Add(requestStatsBucketDuration) {
		key := cursor.Format("2006-01-02 15:04")
		data = append(data, dto.RequestDailyItem{
			Date:  key,
			Count: byMinute[key],
		})
	}
	return data
}

func buildUserInteractionDailySeries(fromTime time.Time, lastDates uint32, rows []dto.UserInteractionDailyItem) []dto.UserInteractionDailyItem {
	byDate := make(map[string]dto.UserInteractionDailyItem, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row
	}

	data := make([]dto.UserInteractionDailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		item := byDate[date]
		data = append(data, dto.UserInteractionDailyItem{
			Date:         date,
			Interactions: item.Interactions,
			ActiveUsers:  item.ActiveUsers,
		})
	}
	return data
}

func buildUserInteractionHourlySeries(fromTime time.Time, rows []dto.UserInteractionDailyItem) []dto.UserInteractionDailyItem {
	byHour := make(map[string]dto.UserInteractionDailyItem, len(rows))
	for _, row := range rows {
		byHour[row.Date] = row
	}

	loc := fromTime.Location()
	now := time.Now().In(loc)
	endHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	startHour := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)

	hours := int(endHour.Sub(startHour).Hours()) + 1
	if hours < 1 {
		hours = 1
	}

	data := make([]dto.UserInteractionDailyItem, 0, hours)
	for cursor := startHour; !cursor.After(endHour); cursor = cursor.Add(time.Hour) {
		key := cursor.Format("2006-01-02 15:00")
		item := byHour[key]
		data = append(data, dto.UserInteractionDailyItem{
			Date:         key,
			Interactions: item.Interactions,
			ActiveUsers:  item.ActiveUsers,
		})
	}
	return data
}

func buildUserInteractionMinuteSeries(fromTime time.Time, rows []dto.UserInteractionDailyItem) []dto.UserInteractionDailyItem {
	byMinute := make(map[string]dto.UserInteractionDailyItem, len(rows))
	for _, row := range rows {
		byMinute[row.Date] = row
	}

	loc := fromTime.Location()
	now := time.Now().In(loc)
	endMinute := now.Truncate(requestStatsBucketDuration)
	startMinute := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)

	data := make([]dto.UserInteractionDailyItem, 0)
	for cursor := startMinute; !cursor.After(endMinute); cursor = cursor.Add(requestStatsBucketDuration) {
		key := cursor.Format("2006-01-02 15:04")
		item := byMinute[key]
		data = append(data, dto.UserInteractionDailyItem{
			Date:         key,
			Interactions: item.Interactions,
			ActiveUsers:  item.ActiveUsers,
		})
	}
	return data
}
