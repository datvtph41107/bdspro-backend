package usecase

import (
	"context"
	"errors"
	"social/internal/domain"
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
	exist, err := u.reportRepo.ExistByNewsFeedIdAndUserId(ctx, report.NewsFeedID, report.UserID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("bạn đã báo cáo bài viết này")
	}
	exist, err = u.reportReasonRepo.ExistedByID(ctx, report.ReasonID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("lý do báo cáo không tồn tại")
	}
	exist, err = u.newsFeedRepo.UserAvaiableReact(ctx, profileID, report.NewsFeedID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("bài viết không tồn tại hoặc không thể báo cáo")
	}
	return u.reportRepo.Create(ctx, report)
}
