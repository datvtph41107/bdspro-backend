package usecase

import (
	"context"
	"errors"
	"testing"

	_errors "common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"user/internal"
	"user/internal/dto"
)

type organizationMembershipProviderStub struct {
	member *dto.OrganizationMember
	err    error
}

func (s organizationMembershipProviderStub) GetOrganizationMember(
	context.Context,
	uint64,
	uint64,
) (*dto.OrganizationMember, error) {
	return s.member, s.err
}

func TestValidateOrganizationMembership(t *testing.T) {
	providerFailure := errors.New("organization service unavailable")

	tests := []struct {
		name        string
		provider    organizationMembershipProvider
		profileID   uint64
		orgID       uint64
		wantSpec    _errors.Spec
		wantErrIs   error
		wantNoError bool
	}{
		{
			name: "active member is accepted",
			provider: organizationMembershipProviderStub{member: &dto.OrganizationMember{
				OrganizationId: 7,
				UserId:         42,
				Status:         organizationMemberStatusActive,
			}},
			profileID:   42,
			orgID:       7,
			wantNoError: true,
		},
		{
			name:      "missing membership is forbidden",
			provider:  organizationMembershipProviderStub{},
			profileID: 42,
			orgID:     7,
			wantSpec:  service.OrganizationMembershipInactive,
		},
		{
			name: "inactive membership is forbidden",
			provider: organizationMembershipProviderStub{member: &dto.OrganizationMember{
				OrganizationId: 7,
				UserId:         42,
				Status:         20,
			}},
			profileID: 42,
			orgID:     7,
			wantSpec:  service.OrganizationMembershipInactive,
		},
		{
			name: "mismatched member identity is forbidden",
			provider: organizationMembershipProviderStub{member: &dto.OrganizationMember{
				OrganizationId: 7,
				UserId:         99,
				Status:         organizationMemberStatusActive,
			}},
			profileID: 42,
			orgID:     7,
			wantSpec:  service.OrganizationMembershipInactive,
		},
		{
			name:      "provider failure is preserved",
			provider:  organizationMembershipProviderStub{err: providerFailure},
			profileID: 42,
			orgID:     7,
			wantErrIs: providerFailure,
		},
		{
			name:      "missing provider fails closed",
			profileID: 42,
			orgID:     7,
			wantSpec:  service.OrganizationMembershipAuthenticationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &AuthUsecase{OrganizationClient: tt.provider}
			err := uc.validateOrganizationMembership(context.Background(), tt.profileID, tt.orgID)
			if tt.wantNoError {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				return
			}
			application, ok := _errors.As(err)
			require.True(t, ok)
			assert.Equal(t, tt.wantSpec.Key(), application.Key())
			assert.Equal(t, tt.wantSpec.RPCCode(), application.RPCCode())
			assert.Equal(t, tt.wantSpec.LegacyCode(), application.Spec().LegacyCode())
		})
	}
}
