package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"

	"gorm.io/datatypes"
)

type ReportUsecase interface {
	CreateAsync(ctx context.Context, userID uint64, reqDTO dto.ReportDTO) (uint64, error)
	GetStatus(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error)
	Download(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error)
	ListUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error)
	Delete(ctx context.Context, userID uint64, reportID uint64) error
	AdminList(ctx context.Context, filter repo.AdminReportFilter) ([]domain.Report, int64, error)
	AdminSummary(ctx context.Context) (map[uint32]int64, error)
	AdminGet(ctx context.Context, reportID uint64) (*domain.Report, error)
	AdminCreateInternal(ctx context.Context, actorID uint64, req dto.AdminCreateReportRequest) (*domain.Report, error)
	AdminAssign(ctx context.Context, reportID, assigneeID uint64) error
	AdminUpdateSeverity(ctx context.Context, reportID uint64, severity uint32) error
	AdminUpdateQaStatus(ctx context.Context, reportID uint64, qaStatus uint32) error
	AdminUpdateLinkage(ctx context.Context, reportID uint64, linkage map[string]interface{}) error
	AdminUpdateImages(ctx context.Context, reportID uint64, images []string) error
	AdminClose(ctx context.Context, reportID uint64, reject bool, note string) error
	AdminDelete(ctx context.Context, reportID uint64) error
}

type reportUsecase struct {
	repo repo.ReportRepository
}

func NewReportUsecase(repo repo.ReportRepository) ReportUsecase {
	return &reportUsecase{repo: repo}
}

func (u *reportUsecase) CreateAsync(ctx context.Context, userID uint64, reqDTO dto.ReportDTO) (uint64, error) {
	var targetSnapshot datatypes.JSON
	report := &domain.Report{
		UserID:         userID,
		ProblemReport:  reqDTO.ProblemReport,
		ReportType:     reqDTO.ReportType,
		Description:    reqDTO.Description,
		TargetID:       reqDTO.TargetId,
		TargetSnapshot: targetSnapshot,
		Status:         10,
		QaStatus:       enums.ReportQaStatusSubmitted,
		Severity:       enums.ReportSeverityMedium,
	}
	if err := u.repo.Create(ctx, report); err != nil {
		return 0, err
	}
	return report.ID, nil
}

func (u *reportUsecase) processReport(reportID, userID uint64, reportType uint32, targetID *string, targetData interface{}, profile, format uint32, includeMapImage bool) {
	_ = u.repo.UpdateStatus(context.Background(), reportID, 20, nil)
	fileURL := "https://storage.example.com/reports/" + string(rune(reportID)) + ".pdf"
	fileSize := int64(1024)
	fileHash := "dummyhash"
	_ = u.repo.UpdateFileInfo(context.Background(), reportID, &fileURL, &fileSize, &fileHash, nil)
}

func (u *reportUsecase) GetStatus(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error) {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return nil, errors.New("report not found")
	}
	if report.UserID != userID {
		return nil, errors.New("permission denied")
	}
	return report, nil
}

func (u *reportUsecase) Download(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error) {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return nil, errors.New("report not found")
	}
	if report.UserID != userID {
		return nil, errors.New("permission denied")
	}
	if report.Status != 30 {
		return nil, errors.New("report not ready")
	}
	return report, nil
}

func (u *reportUsecase) ListUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error) {
	return u.repo.ListByUser(ctx, userID, reportType, status, page, limit)
}

func (u *reportUsecase) Delete(ctx context.Context, userID uint64, reportID uint64) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return errors.New("report not found")
	}
	if report.UserID != userID {
		return errors.New("permission denied")
	}
	return u.repo.Delete(ctx, reportID)
}

func (u *reportUsecase) AdminList(ctx context.Context, filter repo.AdminReportFilter) ([]domain.Report, int64, error) {
	return u.repo.AdminList(ctx, filter)
}

func (u *reportUsecase) AdminSummary(ctx context.Context) (map[uint32]int64, error) {
	return u.repo.AdminSummary(ctx)
}

func (u *reportUsecase) AdminGet(ctx context.Context, reportID uint64) (*domain.Report, error) {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return nil, errors.New("report not found")
	}
	return report, nil
}

func (u *reportUsecase) AdminCreateInternal(ctx context.Context, actorID uint64, req dto.AdminCreateReportRequest) (*domain.Report, error) {
	imagesJSON := datatypes.JSON([]byte("[]"))
	if len(req.Images) > 0 {
		b, err := json.Marshal(req.Images)
		if err != nil {
			return nil, err
		}
		imagesJSON = datatypes.JSON(b)
	}
	report := &domain.Report{
		UserID:          actorID,
		Title:           req.Title,
		Description:     req.Description,
		ProblemReport:   enums.ProblemReport(req.ProblemReport),
		ReportType:      enums.ReportType(req.ReportType),
		Severity:        enums.ReportSeverity(req.Severity),
		QaStatus:        enums.ReportQaStatusSubmitted,
		Status:          10,
		Profile:         10,
		Format:          10,
		SupportTicketID: req.SupportTicketID,
		TargetID:        req.TargetID,
		Images:          imagesJSON,
	}
	if req.Severity == 0 {
		report.Severity = enums.ReportSeverityMedium
	}
	if req.ProblemReport == 0 {
		report.ProblemReport = enums.ProblemReport_Other
	}
	if req.ReportType == 0 {
		report.ReportType = enums.ReportTypeParcel
	}
	if err := u.repo.Create(ctx, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (u *reportUsecase) AdminAssign(ctx context.Context, reportID, assigneeID uint64) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return errors.New("report not found")
	}
	report.AssigneeID = &assigneeID
	if report.QaStatus == enums.ReportQaStatusSubmitted {
		report.QaStatus = enums.ReportQaStatusVerifying
	}
	return u.repo.Update(ctx, report)
}

func (u *reportUsecase) AdminUpdateSeverity(ctx context.Context, reportID uint64, severity uint32) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return errors.New("report not found")
	}
	report.Severity = enums.ReportSeverity(severity)
	return u.repo.Update(ctx, report)
}

func (u *reportUsecase) AdminUpdateQaStatus(ctx context.Context, reportID uint64, qaStatus uint32) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return errors.New("report not found")
	}
	to := enums.ReportQaStatus(qaStatus)
	if !enums.CanTransitionReportQa(report.QaStatus, to) {
		return errors.New("STATUS_TRANSITION_INVALID")
	}
	report.QaStatus = to
	if to == enums.ReportQaStatusClosed || to == enums.ReportQaStatusRejected {
		now := time.Now()
		report.ClosedAt = &now
	}
	return u.repo.Update(ctx, report)
}

func (u *reportUsecase) AdminUpdateLinkage(ctx context.Context, reportID uint64, linkage map[string]interface{}) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return errors.New("report not found")
	}
	b, err := json.Marshal(linkage)
	if err != nil {
		return err
	}
	report.Linkage = datatypes.JSON(b)
	return u.repo.Update(ctx, report)
}

func (u *reportUsecase) AdminUpdateImages(ctx context.Context, reportID uint64, images []string) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return errors.New("report not found")
	}
	// Client gửi full array (thêm/xoá local rồi lưu) — thay thế, không merge
	next := sanitizeImageList(images)
	b, err := json.Marshal(next)
	if err != nil {
		return err
	}
	report.Images = datatypes.JSON(b)
	return u.repo.Update(ctx, report)
}

func unmarshalReportImages(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []string{}
	}
	return out
}

func sanitizeImageList(images []string) []string {
	seen := make(map[string]struct{}, len(images))
	out := make([]string, 0, len(images))
	for _, s := range images {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func mergeUniqueStrings(existing, incoming []string) []string {
	return sanitizeImageList(append(append([]string{}, existing...), incoming...))
}

func (u *reportUsecase) AdminClose(ctx context.Context, reportID uint64, reject bool, note string) error {
	to := enums.ReportQaStatusClosed
	if reject {
		to = enums.ReportQaStatusRejected
	}
	_ = note
	return u.AdminUpdateQaStatus(ctx, reportID, uint32(to))
}

func (u *reportUsecase) AdminDelete(ctx context.Context, reportID uint64) error {
	return u.repo.Delete(ctx, reportID)
}
