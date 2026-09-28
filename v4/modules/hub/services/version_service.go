package services

import (
	"context"

	hubpb "pb/types/hub"

	"userv4/modules/hub/dto"
)

// VersionService owns app-bundle version flows for modules/hub.
type VersionService interface {
	CheckAppUpdate(ctx context.Context, appName string, req dto.CheckAppUpdateRequest) (*hubpb.CheckAppUpdateResponse, error)
	CreateVersionFromDev(ctx context.Context, req *hubpb.CreateVersionRequest) (*hubpb.CreateVersionResponse, error)
}
