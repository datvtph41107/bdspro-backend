package mapper

import (
	_models "common/models"
	_utils "common/utils"
	"context"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

type VersionMapper struct{}

func NewVersionMapper() *VersionMapper {
	return &VersionMapper{}
}

func (m *VersionMapper) CreateRequestToEntity(req *hubpb.CreateVersionRequest) *domain.VersionEntity {
	return &domain.VersionEntity{
		AppName:      req.GetAppName(),
		Platform:     req.GetPlatform(),
		VersionName:  req.GetVersionName(),
		BuildNumber:  req.GetBuildNumber(),
		ForceUpdate:  req.GetForceUpdate(),
		Active:       req.GetActive(),
		DownloadURL:  req.GetDownloadUrl(),
		ReleaseNotes: req.GetReleaseNotes(),
		BundleName:   req.GetBundleName(),
		BundleSize:   req.GetBundleSize(),
	}
}

func (m *VersionMapper) UpdateRequestToEntity(ctx context.Context, req *hubpb.UpdateVersionRequest) *domain.VersionEntity {
	entity := &domain.VersionEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.GetId(),
		},
		ForceUpdate:  req.GetForceUpdate(),
		Active:       req.GetActive(),
		ReleaseNotes: req.GetReleaseNotes(),
	}
	return entity
}

func (m *VersionMapper) EntityToProto(entity *domain.VersionEntity) *hubpb.VersionItem {
	return &hubpb.VersionItem{
		Id:           entity.ID,
		AppName:      entity.AppName,
		Platform:     entity.Platform,
		VersionName:  entity.VersionName,
		BuildNumber:  entity.BuildNumber,
		ForceUpdate:  entity.ForceUpdate,
		Active:       entity.Active,
		DownloadUrl:  entity.DownloadURL,
		ReleaseNotes: entity.ReleaseNotes,
		CreatedAt:    _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(entity.UpdatedAt),
		BundleName:   entity.BundleName,
		BundleSize:   entity.BundleSize,
	}
}

func (m *VersionMapper) EntitiesToProto(entities []*domain.VersionEntity) []*hubpb.VersionItem {
	result := make([]*hubpb.VersionItem, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.EntityToProto(entity))
	}
	return result
}
