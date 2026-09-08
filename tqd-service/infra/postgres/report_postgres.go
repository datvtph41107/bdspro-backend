package postgres

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) repo.ReportRepository {
	return &reportRepo{db: db}
}

func (r *reportRepo) Create(ctx context.Context, report *domain.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepo) Update(ctx context.Context, report *domain.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepo) UpdateAssignee(ctx context.Context, id, assigneeID uint64, qaStatus uint32) error {
	return r.db.WithContext(ctx).Model(&domain.Report{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"assignee_id": assigneeID,
			"qa_status":   qaStatus,
		}).Error
}

func (r *reportRepo) UpdateStatus(ctx context.Context, id uint64, status uint32, errorMsg *string) error {
	updates := map[string]interface{}{"status": status}
	if errorMsg != nil {
		updates["error_message"] = *errorMsg
	}
	if status == 30 {
		updates["completed_at"] = gorm.Expr("NOW()")
	}
	return r.db.WithContext(ctx).Model(&domain.Report{}).Where("id = ?", id).Updates(updates).Error
}

func (r *reportRepo) UpdateFileInfo(ctx context.Context, id uint64, fileURL *string, fileSize *int64, fileHash *string, completedAt *string) error {
	updates := map[string]interface{}{
		"file_url":  fileURL,
		"file_size": fileSize,
		"file_hash": fileHash,
		"status":    30,
	}
	if completedAt != nil {
		updates["completed_at"] = *completedAt
	} else {
		updates["completed_at"] = gorm.Expr("NOW()")
	}
	return r.db.WithContext(ctx).Model(&domain.Report{}).Where("id = ?", id).Updates(updates).Error
}

func (r *reportRepo) GetByID(ctx context.Context, id uint64) (*domain.Report, error) {
	var report domain.Report
	err := r.db.WithContext(ctx).First(&report, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

func (r *reportRepo) ListByUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.Report{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	if reportType != nil {
		query = query.Where("report_type = ?", *reportType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reports).Error
	return reports, total, err
}

func (r *reportRepo) AdminList(ctx context.Context, filter repo.AdminReportFilter) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.Report{}).Where("deleted_at IS NULL")
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.ReportType != nil {
		query = query.Where("report_type = ?", *filter.ReportType)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ProblemReport != nil {
		query = query.Where("problem_report = ?", *filter.ProblemReport)
	}
	if len(filter.QaStatuses) > 0 {
		query = query.Where("qa_status IN ?", filter.QaStatuses)
	} else if filter.QaStatus != nil {
		query = query.Where("qa_status = ?", *filter.QaStatus)
	}
	if len(filter.ExcludeQaStatuses) > 0 {
		query = query.Where("qa_status NOT IN ?", filter.ExcludeQaStatuses)
	}
	if filter.Severity != nil {
		query = query.Where("severity = ?", *filter.Severity)
	}
	if filter.UnassignedOnly {
		query = query.Where("assignee_id IS NULL")
	} else if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filter.AssigneeID)
	}
	if filter.Q != "" {
		like := "%" + filter.Q + "%"
		query = query.Where("(title ILIKE ? OR left(description, 500) ILIKE ?)", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	err := query.
		Select("id, user_id, title, problem_report, report_type, status, qa_status, severity, assignee_id, support_ticket_id, target_id, linkage, created_at, closed_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&reports).Error
	return reports, total, err
}

func (r *reportRepo) AdminQueueSummary(ctx context.Context, assigneeID uint64) (repo.QueueSummary, error) {
	out := repo.QueueSummary{}
	base := r.db.WithContext(ctx).Model(&domain.Report{}).Where("deleted_at IS NULL")

	mineStatuses := []uint32{10, 20, 30, 40} // submitted, verifying, need_more_info, verified
	if err := base.Session(&gorm.Session{}).
		Where("assignee_id = ? AND qa_status IN ?", assigneeID, mineStatuses).
		Count(&out.Mine).Error; err != nil {
		return out, err
	}
	if err := base.Session(&gorm.Session{}).
		Where("assignee_id IS NULL AND qa_status = ?", 10).
		Count(&out.Unassigned).Error; err != nil {
		return out, err
	}
	if err := base.Session(&gorm.Session{}).
		Where("qa_status = ?", 50).
		Count(&out.DataFix).Error; err != nil {
		return out, err
	}
	return out, nil
}

func (r *reportRepo) AdminSummary(ctx context.Context) (map[uint32]int64, error) {
	type row struct {
		QaStatus uint32
		Count    int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&domain.Report{}).
		Select("qa_status as qa_status, count(*) as count").
		Where("deleted_at IS NULL").
		Group("qa_status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := map[uint32]int64{}
	for _, r0 := range rows {
		out[r0.QaStatus] = r0.Count
	}
	return out, nil
}

func (r *reportRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Report{}, id).Error
}

func (r *reportRepo) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at IS NOT NULL AND expires_at < NOW()").Delete(&domain.Report{}).Error
}

func (r *reportRepo) LookupParcelLabels(ctx context.Context, ids []uint64) (map[uint64]string, error) {
	out := map[uint64]string{}
	if len(ids) == 0 {
		return out, nil
	}
	type row struct {
		ID          uint64  `gorm:"column:id"`
		AddressText string  `gorm:"column:address_text"`
		MapNumber   *string `gorm:"column:map_number"`
		LandNumber  *string `gorm:"column:land_number"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT id,
			COALESCE(address_text, '') AS address_text,
			map_number,
			land_number
		FROM parcels
		WHERE deleted_at IS NULL AND id IN ?
	`, ids).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		addr := strings.TrimSpace(strings.TrimPrefix(row.AddressText, ","))
		sheet := ""
		if row.MapNumber != nil || row.LandNumber != nil {
			m, l := "?", "?"
			if row.MapNumber != nil && strings.TrimSpace(*row.MapNumber) != "" {
				m = strings.TrimSpace(*row.MapNumber)
			}
			if row.LandNumber != nil && strings.TrimSpace(*row.LandNumber) != "" {
				l = strings.TrimSpace(*row.LandNumber)
			}
			sheet = "Tờ " + m + " · Thửa " + l
		}
		label := addr
		if addr != "" && sheet != "" {
			label = addr + " (" + sheet + ")"
		} else if sheet != "" {
			label = sheet + " (#" + strconv.FormatUint(row.ID, 10) + ")"
		} else if label == "" {
			label = "Thửa #" + strconv.FormatUint(row.ID, 10)
		}
		out[row.ID] = label
	}
	return out, nil
}

func (r *reportRepo) LookupLayerLabels(ctx context.Context, ids []uint64) (map[uint64]string, error) {
	out := map[uint64]string{}
	if len(ids) == 0 {
		return out, nil
	}
	type row struct {
		ID          uint64 `gorm:"column:id"`
		DisplayName string `gorm:"column:display_name"`
		Name        string `gorm:"column:name"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT id, COALESCE(display_name, '') AS display_name, COALESCE(name, '') AS name
		FROM qh_layers
		WHERE id IN ?
	`, ids).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		label := strings.TrimSpace(row.DisplayName)
		if label == "" {
			label = strings.TrimSpace(row.Name)
		}
		if label == "" {
			label = "#" + strconv.FormatUint(row.ID, 10)
		}
		out[row.ID] = label
	}
	return out, nil
}

func (r *reportRepo) AddEvent(ctx context.Context, event *domain.ReportEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *reportRepo) ListEventsByReport(ctx context.Context, reportID uint64, limit int) ([]domain.ReportEvent, error) {
	var events []domain.ReportEvent
	q := r.db.WithContext(ctx).Where("report_id = ?", reportID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&events).Error
	return events, err
}

func (r *reportRepo) ListEvents(ctx context.Context, filter repo.AdminEventFilter) ([]domain.ReportEvent, int64, error) {
	var events []domain.ReportEvent
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.ReportEvent{})
	if filter.ReportID != nil {
		query = query.Where("report_id = ?", *filter.ReportID)
	}
	if filter.ActorID != nil {
		query = query.Where("actor_id = ?", *filter.ActorID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.From != nil && *filter.From != "" {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil && *filter.To != "" {
		query = query.Where("created_at <= ?", *filter.To)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&events).Error
	return events, total, err
}
