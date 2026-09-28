package services

import (
	"context"

	hubpb "pb/types/hub"

	"userv4/modules/hub/dto"
)

// UserGuideService owns user-facing user-guide flows for modules/hub.
type UserGuideService interface {
	GetUserGuideSimple(ctx context.Context, req dto.GetUserGuideSimpleRequest) (*hubpb.GetUserGuideSimpleResponse, error)
	GetUserGuide(ctx context.Context, id uint64) (*hubpb.GetUserGuideResponse, error)
	GetUserGuideByKey(ctx context.Context, key string) (*hubpb.GetUserGuideByKeyResponse, error)
	GetAllUserGuidesByKey(ctx context.Context, req dto.GetAllUserGuidesByKeyRequest) (*hubpb.GetAllUserGuidesByKeyResponse, error)
	GetUserGuideGroups(ctx context.Context) (*hubpb.GetUserGuideGroupsResponse, error)
}
