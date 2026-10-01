package repository

import (
	"context"
	"time"

	qh "userv4/modules/tqd/entity/qh"
	"userv4/modules/tqd/dto"
)

type IMapWorkspaceRepo interface {
	GetParcelWorkspacePreview(ctx context.Context, parcelID uint64) (*dto.ParcelWorkspacePreviewRow, error)
	GetParcelWorkspacePreviewsByIDs(ctx context.Context, parcelIDs []uint64) ([]dto.ParcelWorkspacePreviewRow, error)

	GetRegionWorkspacePreview(ctx context.Context, regionID uint64) (*dto.RegionWorkspacePreviewRow, error)
	GetRegionWorkspacePreviewsByIDs(ctx context.Context, regionIDs []uint64) ([]dto.RegionWorkspacePreviewRow, error)

	ListFollowedParcels(ctx context.Context, userID uint64, limit, offset int) ([]qh.QHUserFollowedParcel, int64, error)
	UpsertFollowedParcel(ctx context.Context, userID, parcelID uint64, note string) (*qh.QHUserFollowedParcel, error)
	RemoveFollowedParcel(ctx context.Context, userID, followID, parcelID uint64) error

	ListViewHistory(ctx context.Context, userID uint64, limit, offset int) ([]qh.QHUserViewHistory, int64, error)
	TrackViewHistory(ctx context.Context, event qh.QHUserViewEvent, dedupeWindow time.Duration) (*dto.TrackViewHistoryResultDTO, error)
	RemoveViewHistory(ctx context.Context, userID, historyID uint64) error
	ClearViewHistory(ctx context.Context, userID uint64) error

	ListGeneratedReports(ctx context.Context, userID uint64, limit, offset int) ([]qh.QHUserReported, int64, error)
	GetGeneratedReport(ctx context.Context, userID, reportID uint64) (*qh.QHUserReported, error)
	RemoveGeneratedReport(ctx context.Context, userID, reportID uint64) error
	RegenerateGeneratedReport(ctx context.Context, userID, reportID uint64) error
	UpdateGeneratedReportShareURL(ctx context.Context, userID, reportID uint64, shareURL string) error
}
