package handler

import (
	"common/identity"
	"context"
	"testing"

	"organization/env"
	organizationpb "pb/types/organization"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type organizationAuthUsecaseStub struct {
	allowed        bool
	err            error
	profileID      uint64
	organizationID uint32
	permission     string
}

func (s *organizationAuthUsecaseStub) HasPermission(context.Context, uint32, uint32, string) (bool, error) {
	return false, nil
}

func (s *organizationAuthUsecaseStub) Authorize(_ context.Context, profileID uint64, organizationID uint32, permission string) (bool, error) {
	s.profileID = profileID
	s.organizationID = organizationID
	s.permission = permission
	return s.allowed, s.err
}

func (s *organizationAuthUsecaseStub) IsMemberOrganization(context.Context, uint32, uint32) (bool, error) {
	return false, nil
}

func (s *organizationAuthUsecaseStub) IsSystemAdmin(context.Context) error { return nil }

func TestAuthorizeOrganizationSubscriptionCheckout(t *testing.T) {
	t.Parallel()
	authorizer := &organizationAuthUsecaseStub{allowed: true}
	handler := &InternalOrganizationHandler{organizationAuthUsecase: authorizer}

	ctx := authorizedOrganizationCallerContext(t, 42)
	response, err := handler.AuthorizeOrganizationAction(ctx, &organizationpb.AuthorizeOrganizationActionRequest{
		OrganizationId: 7,
		ActorProfileId: 42,
		Action:         organizationpb.OrganizationAction_ORGANIZATION_ACTION_SUBSCRIPTION_CHECKOUT,
	})
	require.NoError(t, err)
	assert.True(t, response.GetAllowed())
	assert.Equal(t, uint64(42), authorizer.profileID)
	assert.Equal(t, uint32(7), authorizer.organizationID)
	assert.Equal(t, env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION, authorizer.permission)
}

func TestAuthorizeOrganizationActionRejectsUnspecifiedAction(t *testing.T) {
	t.Parallel()
	handler := &InternalOrganizationHandler{organizationAuthUsecase: &organizationAuthUsecaseStub{}}

	_, err := handler.AuthorizeOrganizationAction(authorizedOrganizationCallerContext(t, 42), &organizationpb.AuthorizeOrganizationActionRequest{
		OrganizationId: 7,
		ActorProfileId: 42,
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestAuthorizeOrganizationActionRejectsActorMismatch(t *testing.T) {
	t.Parallel()
	handler := &InternalOrganizationHandler{organizationAuthUsecase: &organizationAuthUsecaseStub{}}

	_, err := handler.AuthorizeOrganizationAction(authorizedOrganizationCallerContext(t, 99), &organizationpb.AuthorizeOrganizationActionRequest{
		OrganizationId: 7,
		ActorProfileId: 42,
		Action:         organizationpb.OrganizationAction_ORGANIZATION_ACTION_SUBSCRIPTION_CHECKOUT,
	})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthorizeOrganizationActionRejectsUntrustedService(t *testing.T) {
	t.Parallel()
	handler := &InternalOrganizationHandler{organizationAuthUsecase: &organizationAuthUsecaseStub{}}
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "unknown-service"})
	require.NoError(t, err)
	ctx, err = identity.BindActor(ctx, identity.Actor{ProfileID: 42})
	require.NoError(t, err)

	_, err = handler.AuthorizeOrganizationAction(ctx, &organizationpb.AuthorizeOrganizationActionRequest{
		OrganizationId: 7,
		ActorProfileId: 42,
		Action:         organizationpb.OrganizationAction_ORGANIZATION_ACTION_SUBSCRIPTION_CHECKOUT,
	})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func authorizedOrganizationCallerContext(t *testing.T, profileID uint64) context.Context {
	t.Helper()
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: env.USER_SERVICE_ID})
	require.NoError(t, err)
	ctx, err = identity.BindActor(ctx, identity.Actor{ProfileID: profileID})
	require.NoError(t, err)
	return ctx
}
