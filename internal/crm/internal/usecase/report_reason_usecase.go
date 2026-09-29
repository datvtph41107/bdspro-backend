package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
)

type ReportReasonUsecase struct {
	reportReasonRepo repo.ReportReasonRepo
}

func NewReportReasonUsecase(reportReasonRepo repo.ReportReasonRepo) *ReportReasonUsecase {
	return &ReportReasonUsecase{
		reportReasonRepo: reportReasonRepo,
	}
}

// CreateReportReason creates a new report reason
func (u *ReportReasonUsecase) CreateReportReason(ctx context.Context, req *dto.ReportReasonCreateRequest) (*dto.ReportReasonResponse, error) {
	reason := &domain.ReportReason{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	reason, err := u.reportReasonRepo.Create(ctx, reason)
	if err != nil {
		return nil, err
	}

	return &dto.ReportReasonResponse{
		ID:          reason.ID,
		Name:        reason.Name,
		Description: reason.Description,
		IsActive:    reason.IsActive,
		CreatedAt:   reason.CreatedAt,
		UpdatedAt:   reason.UpdatedAt,
	}, nil
}

// GetReportReasonByID gets a report reason by ID
func (u *ReportReasonUsecase) GetReportReasonByID(ctx context.Context, id uint64) (*dto.ReportReasonResponse, error) {
	reason, err := u.reportReasonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.ReportReasonResponse{
		ID:          reason.ID,
		Name:        reason.Name,
		Description: reason.Description,
		IsActive:    reason.IsActive,
		CreatedAt:   reason.CreatedAt,
		UpdatedAt:   reason.UpdatedAt,
	}, nil
}

// GetReportReasonList gets a list of report reasons
func (u *ReportReasonUsecase) GetReportReasonList(ctx context.Context, req *dto.ReportReasonListRequest) (*dto.ReportReasonListResponse, error) {
	reasons, err := u.reportReasonRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReportReasonResponse, len(reasons))
	for i, reason := range reasons {
		responses[i] = dto.ReportReasonResponse{
			ID:          reason.ID,
			Name:        reason.Name,
			Description: reason.Description,
			IsActive:    reason.IsActive,
			CreatedAt:   reason.CreatedAt,
			UpdatedAt:   reason.UpdatedAt,
		}
	}

	return &dto.ReportReasonListResponse{
		Data:  responses,
		Total: int64(len(reasons)),
	}, nil
}

// UpdateReportReason updates a report reason
func (u *ReportReasonUsecase) UpdateReportReason(ctx context.Context, req *dto.ReportReasonUpdateRequest) (*dto.ReportReasonResponse, error) {
	reason := &domain.ReportReason{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	err := u.reportReasonRepo.Update(ctx, reason)
	if err != nil {
		return nil, err
	}

	return &dto.ReportReasonResponse{
		ID:          reason.ID,
		Name:        reason.Name,
		Description: reason.Description,
		IsActive:    reason.IsActive,
		CreatedAt:   reason.CreatedAt,
		UpdatedAt:   reason.UpdatedAt,
	}, nil
}

// DeleteReportReason deletes a report reason
func (u *ReportReasonUsecase) DeleteReportReason(ctx context.Context, id uint64) error {
	return u.reportReasonRepo.Delete(ctx, id)
}

// GetAllReportReasons gets all report reasons
func (u *ReportReasonUsecase) GetAllReportReasons(ctx context.Context) ([]*dto.ReportReasonResponse, error) {
	reasons, err := u.reportReasonRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ReportReasonResponse, len(reasons))
	for i, reason := range reasons {
		responses[i] = &dto.ReportReasonResponse{
			ID:          reason.ID,
			Name:        reason.Name,
			Description: reason.Description,
			IsActive:    reason.IsActive,
			CreatedAt:   reason.CreatedAt,
			UpdatedAt:   reason.UpdatedAt,
		}
	}

	return responses, nil
}