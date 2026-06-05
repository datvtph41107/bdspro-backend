package repo

import (
	"context"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
)

type IMapWorkspaceRepo interface {
	// Parcel spatial preview
	GetParcelWorkspacePreview(ctx context.Context, parcelID uint64) (*dto.ParcelWorkspacePreviewRow, error)
	GetParcelWorkspacePreviewsByIDs(ctx context.Context, parcelIDs []uint64) ([]dto.ParcelWorkspacePreviewRow, error)

	// Followed parcels
	ListFollowedParcels(ctx context.Context, userID uint64, limit, offset int) ([]qh_domain.QHUserFollowedParcel, int64, error)
	UpsertFollowedParcel(ctx context.Context, userID, parcelID uint64, note string) (*qh_domain.QHUserFollowedParcel, error)
	RemoveFollowedParcel(ctx context.Context, userID, followID, parcelID uint64) error

	// View history
	ListViewHistory(ctx context.Context, userID uint64, limit, offset int) ([]qh_domain.QHUserViewHistory, int64, error)
	AddViewHistory(ctx context.Context, item qh_domain.QHUserViewHistory) (*qh_domain.QHUserViewHistory, error)
	RemoveViewHistory(ctx context.Context, userID, historyID uint64) error
	ClearViewHistory(ctx context.Context, userID uint64) error

	// Reports
	ListGeneratedReports(ctx context.Context, userID uint64, limit, offset int) ([]qh_domain.QHUserReported, int64, error)
	GetGeneratedReport(ctx context.Context, userID, reportID uint64) (*qh_domain.QHUserReported, error)
	RemoveGeneratedReport(ctx context.Context, userID, reportID uint64) error
	UpdateGeneratedReportStatus(ctx context.Context, userID, reportID uint64, status uint32) error
	UpdateGeneratedReportShareURL(ctx context.Context, userID, reportID uint64, shareURL string) error
}
