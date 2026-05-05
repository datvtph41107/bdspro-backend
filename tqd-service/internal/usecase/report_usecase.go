package usecase

import (
	"context"
	"errors"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/datatypes"
)

type ReportUsecase interface {
	CreateAsync(ctx context.Context, userID uint64, reqDTO dto.ReportDTO) (uint64, error)
	GetStatus(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error)
	Download(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error)
	ListUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error)
	Delete(ctx context.Context, userID uint64, reportID uint64) error
	AdminList(ctx context.Context, userID *uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error)
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
	// if targetData != nil {
	// 	b, err := json.Marshal(targetData)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	targetSnapshot = datatypes.JSON(b)
	// }
	report := &domain.Report{
		UserID:         userID,
		ProblemReport:  reqDTO.ProblemReport,
		ReportType:     reqDTO.ReportType,
		Description:    reqDTO.Description,
		TargetID:       reqDTO.TargetId,
		TargetSnapshot: targetSnapshot,
		// Profile:        profile,
		// Format:         format,
		Status: 10,
	}
	if err := u.repo.Create(ctx, report); err != nil {
		return 0, err
	}
	// go u.processReport(report.ID, userID, reportType, targetID, targetData, profile, format, includeMapImage)
	return report.ID, nil
}

func (u *reportUsecase) processReport(reportID, userID uint64, reportType uint32, targetID *string, targetData interface{}, profile, format uint32, includeMapImage bool) {
	_ = u.repo.UpdateStatus(context.Background(), reportID, 20, nil)
	// TODO: generate actual report PDF
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

func (u *reportUsecase) AdminList(ctx context.Context, userID *uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error) {
	return u.repo.AdminList(ctx, userID, reportType, status, page, limit)
}

func (u *reportUsecase) AdminDelete(ctx context.Context, reportID uint64) error {
	return u.repo.Delete(ctx, reportID)
}
