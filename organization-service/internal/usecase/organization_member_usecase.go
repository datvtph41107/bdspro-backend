package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationMemberUsecase interface {
	ExistByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (bool, error)
}

type organizationMemberUsecase struct {
	organizationMemberRepository repository.OrganizationMemberRepository
}

func NewOrganizationMemberUsecase(organizationMemberRepository repository.OrganizationMemberRepository) OrganizationMemberUsecase {
	return &organizationMemberUsecase{organizationMemberRepository: organizationMemberRepository}
}

func (u *organizationMemberUsecase) ExistByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (bool, error) {
	return u.organizationMemberRepository.ExistByUserIdAndOrganizationId(ctx, userId, organizationId)
}
