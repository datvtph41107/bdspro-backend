package usecase

import (
	_errors "common/errors"
	"common/identity"
	"context"
	"math"
	"strings"

	"organization/env"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type OrganizationAuthUsecase interface {
	HasPermission(ctx context.Context, userId uint32, organizationId uint32, permission string) (bool, error)
	Authorize(ctx context.Context, userID uint64, organizationID uint32, permission string) (bool, error)
	IsMemberOrganization(ctx context.Context, userId uint32, organizationId uint32) (bool, error)
	IsSystemAdmin(ctx context.Context) error
}

type organizationAuthUsecase struct {
	organizationMemberRepository repository.OrganizationMemberRepository
	organizationRoleRepository   repository.OrganizationRoleRepository
	organizationRepository       repository.OrganizationRepository
}

func NewOrganizationAuthUsecase(
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationRepository repository.OrganizationRepository,
) OrganizationAuthUsecase {
	return &organizationAuthUsecase{
		organizationMemberRepository: organizationMemberRepository,
		organizationRoleRepository:   organizationRoleRepository,
		organizationRepository:       organizationRepository,
	}
}
func (u *organizationAuthUsecase) HasPermission(ctx context.Context, userId uint32, organizationId uint32, permission string) (bool, error) {
	if organizationId == 0 {
		// Global account role chỉ hợp lệ cho policy hệ thống; tuyệt đối không
		// dùng nó thay organization role khi organizationID khác 0.
		return strings.EqualFold(strings.TrimSpace(permission), env.ADMIN_ROLE_KEY) && isSystemAdminActor(ctx, uint64(userId)), nil
	}
	return u.Authorize(ctx, uint64(userId), organizationId, permission)
}

// Authorize là policy owner cho quyền trong organization. Mọi dependency lỗi,
// membership không active, role lệch organization hoặc permission thiếu đều
// fail closed; caller không được tự suy diễn quyền từ role global trong token.
func (u *organizationAuthUsecase) Authorize(ctx context.Context, userID uint64, organizationID uint32, permission string) (bool, error) {
	permission = strings.TrimSpace(permission)
	if userID == 0 || organizationID == 0 || permission == "" {
		return false, nil
	}

	member, err := u.organizationMemberRepository.FindByUserIDAndOrganizationID(ctx, userID, organizationID)
	if err != nil {
		return false, err
	}
	if member == nil || member.Status != entity.OrganizationMemberStatusActive || member.RemovedAt != nil {
		return false, nil
	}
	if member.RoleId == 0 {
		// Compatibility production: code cũ tạo admin role nhưng không gán
		// role_id vào owner member. Chỉ owner_id chính thức được phép đi qua;
		// active member thường có role_id=0 vẫn bị từ chối.
		if u.organizationRepository == nil {
			return false, nil
		}
		organization, err := u.organizationRepository.FindById(ctx, organizationID)
		if err != nil {
			return false, err
		}
		return organization != nil && organization.OwnerId == userID, nil
	}
	if member.RoleId > math.MaxUint32 {
		return false, nil
	}

	role, err := u.organizationRoleRepository.FindById(ctx, uint32(member.RoleId))
	if err != nil {
		return false, err
	}
	if role == nil || role.OrganizationId != organizationID {
		return false, nil
	}

	// Role admin do organization tạo ra là owner của organization và có toàn
	// quyền. Các role còn lại bắt buộc có permission explicit trong DB.
	if strings.EqualFold(strings.TrimSpace(role.Key), env.ADMIN_ROLE_KEY) {
		return true, nil
	}
	for _, candidate := range role.Permissions {
		if candidate != nil && strings.EqualFold(strings.TrimSpace(candidate.Key), permission) {
			return true, nil
		}
	}
	return false, nil
}

func (u *organizationAuthUsecase) IsMemberOrganization(ctx context.Context, userId uint32, organizationId uint32) (bool, error) {
	organizationMember, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, userId, organizationId)
	if err != nil {
		return false, err
	}
	return organizationMember != nil, nil
}

// IsSystemAdmin chỉ tin global role đã được transport xác thực và bind vào Actor.
func (u *organizationAuthUsecase) IsSystemAdmin(ctx context.Context) error {
	if isSystemAdminActor(ctx, 0) {
		return nil
	}
	return _errors.InvalidRequest(&errdetails.ErrorInfo{
		Reason: "Bạn không phải quản trị viên hệ thống",
	})
}

func isSystemAdminActor(ctx context.Context, expectedProfileID uint64) bool {
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(strings.TrimSpace(actor.Role), env.ADMIN_ROLE_KEY) {
		return false
	}
	return expectedProfileID == 0 || actor.ProfileID == expectedProfileID
}
