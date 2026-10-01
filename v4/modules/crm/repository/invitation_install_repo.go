package repository

import (
	"context"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/dto"
)

type InvitationInstallRepo interface {
	SendInvitation(ctx context.Context, req []*domain.AppInvitedEntity) ([]*domain.AppInvitedEntity, error)
	ListInvitation(ctx context.Context, profileId uint64, req *dto.ListInvitationRequest) ([]domain.AppInvitedEntity, int64, error)
}