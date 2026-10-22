package repository

import (
	"context"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
)

type InvitationInstallRepo interface {
	SendInvitation(ctx context.Context, req []*domain.AppInvitedEntity) ([]*domain.AppInvitedEntity, error)
	ListInvitation(ctx context.Context, profileId uint64, req *dto.ListInvitationRequest) ([]domain.AppInvitedEntity, int64, error)
}
