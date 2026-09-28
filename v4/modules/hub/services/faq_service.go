package services

import (
	"context"

	hubpb "pb/types/hub"

	"userv4/modules/hub/dto"
)

// FAQService owns user-facing FAQ flows for modules/hub.
type FAQService interface {
	GetFAQSimple(ctx context.Context, req dto.GetFAQSimpleRequest) (*hubpb.GetFAQSimpleResponse, error)
	GetFAQSimpleDetail(ctx context.Context, id uint64) (*hubpb.FAQDetail, error)
}
