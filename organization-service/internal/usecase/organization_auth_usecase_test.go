package usecase

import (
	"common/identity"
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"organization/env"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type organizationMemberRepositoryStub struct {
	repository.OrganizationMemberRepository
	member *entity.OrganizationMember
	err    error
}

func (s organizationMemberRepositoryStub) FindByUserIDAndOrganizationID(context.Context, uint64, uint32) (*entity.OrganizationMember, error) {
	return s.member, s.err
}

type organizationRoleRepositoryStub struct {
	repository.OrganizationRoleRepository
	role *entity.OrganizationRole
	err  error
}

type organizationRepositoryStub struct {
	repository.OrganizationRepository
	organization *entity.Organization
	err          error
}

func (s organizationRepositoryStub) FindById(context.Context, uint32) (*entity.Organization, error) {
	return s.organization, s.err
}

func (s organizationRoleRepositoryStub) FindById(context.Context, uint32) (*entity.OrganizationRole, error) {
	return s.role, s.err
}

func TestOrganizationAuthorize(t *testing.T) {
	t.Parallel()

	removedAt := time.Now()
	databaseErr := errors.New("database unavailable")
	tests := []struct {
		name            string
		member          *entity.OrganizationMember
		memberErr       error
		role            *entity.OrganizationRole
		roleErr         error
		organization    *entity.Organization
		organizationErr error
		permission      string
		want            bool
		wantErr         error
	}{
		{
			name:       "organization admin owns all actions",
			member:     activeOrganizationMember(3),
			role:       &entity.OrganizationRole{Id: 3, OrganizationId: 7, Key: env.ADMIN_ROLE_KEY},
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
			want:       true,
		},
		{
			name:   "custom role needs explicit permission",
			member: activeOrganizationMember(4),
			role: &entity.OrganizationRole{Id: 4, OrganizationId: 7, Key: "finance", Permissions: []*entity.OrganizationPermission{
				{Key: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION},
			}},
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
			want:       true,
		},
		{
			name:       "custom role without permission is denied",
			member:     activeOrganizationMember(4),
			role:       &entity.OrganizationRole{Id: 4, OrganizationId: 7, Key: "member"},
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
		},
		{
			name:         "legacy owner without assigned role remains authorized",
			member:       activeOrganizationMember(0),
			organization: &entity.Organization{OwnerId: 42},
			permission:   env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
			want:         true,
		},
		{
			name:         "legacy roleless member who is not owner is denied",
			member:       activeOrganizationMember(0),
			organization: &entity.Organization{OwnerId: 99},
			permission:   env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
		},
		{
			name:            "legacy owner lookup error fails closed",
			member:          activeOrganizationMember(0),
			organizationErr: databaseErr,
			permission:      env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
			wantErr:         databaseErr,
		},
		{
			name:       "inactive membership is denied",
			member:     &entity.OrganizationMember{UserID: 42, OrganizationID: 7, RoleId: 3, Status: entity.OrganizationMemberStatusInactive},
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
		},
		{
			name:       "removed membership is denied",
			member:     &entity.OrganizationMember{UserID: 42, OrganizationID: 7, RoleId: 3, Status: entity.OrganizationMemberStatusActive, RemovedAt: &removedAt},
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
		},
		{
			name:       "role from another organization is denied",
			member:     activeOrganizationMember(3),
			role:       &entity.OrganizationRole{Id: 3, OrganizationId: 8, Key: env.ADMIN_ROLE_KEY},
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
		},
		{
			name:       "role id overflow is denied before lookup",
			member:     activeOrganizationMember(uint64(math.MaxUint32) + 1),
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
		},
		{
			name:       "membership dependency error fails closed",
			memberErr:  databaseErr,
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
			wantErr:    databaseErr,
		},
		{
			name:       "role dependency error fails closed",
			member:     activeOrganizationMember(3),
			roleErr:    databaseErr,
			permission: env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION,
			wantErr:    databaseErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := NewOrganizationAuthUsecase(
				organizationMemberRepositoryStub{member: tt.member, err: tt.memberErr},
				organizationRoleRepositoryStub{role: tt.role, err: tt.roleErr},
				organizationRepositoryStub{organization: tt.organization, err: tt.organizationErr},
			)

			allowed, err := service.Authorize(context.Background(), 42, 7, tt.permission)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, allowed)
		})
	}
}

func activeOrganizationMember(roleID uint64) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		UserID:         42,
		OrganizationID: 7,
		RoleId:         roleID,
		Status:         entity.OrganizationMemberStatusActive,
	}
}

func TestOrganizationSystemAdminUsesOnlyVerifiedGlobalRole(t *testing.T) {
	t.Parallel()
	service := NewOrganizationAuthUsecase(organizationMemberRepositoryStub{}, organizationRoleRepositoryStub{}, organizationRepositoryStub{})

	adminContext, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42, Role: env.ADMIN_ROLE_KEY})
	require.NoError(t, err)
	allowed, err := service.HasPermission(adminContext, 42, 0, env.ADMIN_ROLE_KEY)
	require.NoError(t, err)
	assert.True(t, allowed)
	require.NoError(t, service.IsSystemAdmin(adminContext))

	allowed, err = service.HasPermission(adminContext, 99, 0, env.ADMIN_ROLE_KEY)
	require.NoError(t, err)
	assert.False(t, allowed, "global admin evidence must belong to the same profile")

	memberContext, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42, Role: "member"})
	require.NoError(t, err)
	allowed, err = service.HasPermission(memberContext, 42, 0, env.ADMIN_ROLE_KEY)
	require.NoError(t, err)
	assert.False(t, allowed)
	require.Error(t, service.IsSystemAdmin(memberContext))
}
