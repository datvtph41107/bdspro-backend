package usecase

import (
	"context"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
	"social/internal/repo"

	_utils "common/utils"
)

type ReportUsecase struct {
	reportRepo       repo.ReportRepo
	reportReasonRepo repo.ReportReasonRepo
	newsFeedRepo     repo.NewsFeedRepo
}

func NewReportUsecase(reportRepo repo.ReportRepo,
	reportReasonRepo repo.ReportReasonRepo,
	newsFeedRepo repo.NewsFeedRepo) *ReportUsecase {
	return &ReportUsecase{
		reportRepo:       reportRepo,
		reportReasonRepo: reportReasonRepo,
		newsFeedRepo:     newsFeedRepo,
	}
}

// hiện tại cho phép báo cáo bài viết công khai
func (u *ReportUsecase) CreateReport(ctx context.Context, report *domain.Report) (*domain.Report, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	report.UserID = profileID

	exist, err := u.reportReasonRepo.ExistedByID(ctx, report.ReasonID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, domain.ErrReportReasonNotFound
	}

	exist, err = u.reportRepo.ExistByTargetIdAndUserIdAndTargetType(ctx, report.TargetID, report.UserID, report.TargetType)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, domain.ErrReportAlreadySubmitted
	}
	if report.TargetType == enums.TargetTypeNewsFeed {

		exist, err = u.newsFeedRepo.UserAvaiableReact(ctx, profileID, report.TargetID)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, domain.ErrReportTargetUnavailable
		}
	}

	if report.TargetType == enums.TargetTypeProfile {

	}

	return u.reportRepo.Create(ctx, report)
}

func (u *ReportUsecase) GetReportMe(ctx context.Context, req *dto.ReportReasonRequest) ([]domain.Report, int64, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	reports, total, err := u.reportRepo.GetByCreatedBy(ctx, profileID, &req.Pagable)
	if err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}

func (u *ReportUsecase) UpdateReportStatus(ctx context.Context, id uint64, status enums.ReportStatus) error {
	// todo: check admin ở đây

	return u.reportRepo.UpdateStatus(ctx, id, status)
}
