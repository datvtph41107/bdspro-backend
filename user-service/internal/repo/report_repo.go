package repo

import (
	"context"
	"user/internal/dto"
	user_model "user/internal/models/user"
)

type ReportRepo interface {
	Create(ctx context.Context, report *user_model.Report) error
	GetByID(ctx context.Context, id uint64) (*user_model.Report, error)
	List(ctx context.Context, req dto.ReportListRequest) ([]user_model.Report, int64, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
	UpdateStatus(ctx context.Context, id uint64, status uint32, response string) error
	CreateProof(ctx context.Context, proof *user_model.ReportProof) error
	GetProofsByReportID(ctx context.Context, reportID uint64) ([]user_model.ReportProof, error)
	DeleteProofsByReportID(ctx context.Context, reportID uint64) error
	GetReasons(ctx context.Context) ([]user_model.ReportReason, error)
}
