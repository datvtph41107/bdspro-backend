package usecase

import (
	"context"
	"errors"
	"social/internal/domain"
	"social/internal/dto"
	iusecase "social/internal/interface"
	"social/internal/repo"
)

type ReportReasonUsecase struct {
	reportReasonRepo  repo.ReportReasonRepo
	permissionUsecase iusecase.PermissionUsecase
}

func NewReportReasonUsecase(
	reportReasonRepo repo.ReportReasonRepo,
	permissionUsecase iusecase.PermissionUsecase,
) *ReportReasonUsecase {
	return &ReportReasonUsecase{
		reportReasonRepo:  reportReasonRepo,
		permissionUsecase: permissionUsecase,
	}
}

// -- admin --
func (u *ReportReasonUsecase) CreateReportReason(ctx context.Context, reportReason *domain.ReportReason) (*domain.ReportReason, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, errors.New("permission denied")
	}
	_, err := u.reportReasonRepo.Create(ctx, reportReason)
	return reportReason, err
}

func (u *ReportReasonUsecase) UpdateReportReason(ctx context.Context, reportReason *domain.ReportReason) (*domain.ReportReason, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, errors.New("permission denied")
	}
	_, err := u.reportReasonRepo.Update(ctx, reportReason)
	return reportReason, err
}

func (u *ReportReasonUsecase) DeleteReportReason(ctx context.Context, id uint64) (*uint64, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, errors.New("permission denied")
	}
	err := u.reportReasonRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (u *ReportReasonUsecase) GetReportReason(ctx context.Context, dto *dto.ReportReasonRequest) ([]domain.ReportReason, int64, error) {
	return u.reportReasonRepo.GetList(ctx, dto)
}

// -- public --
func (u *ReportReasonUsecase) GetPublicReportReason(ctx context.Context) ([]domain.ReportReason, error) {
	return u.reportReasonRepo.GetPublic(ctx)
}
