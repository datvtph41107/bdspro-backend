package grpc

import (
	"common/identity"
	"context"
	"errors"
	organizationpb "pb/types/organization"
	"testing"

	"user/internal/usecase/subscription/checkout"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type organizationCheckoutAuthorizerStub struct {
	allowed        bool
	err            error
	organizationID uint64
	profileID      uint64
	action         organizationpb.OrganizationAction
}

func (s *organizationCheckoutAuthorizerStub) AuthorizeOrganizationAction(_ context.Context, organizationID, actorProfileID uint64, action organizationpb.OrganizationAction) (bool, error) {
	s.organizationID = organizationID
	s.profileID = actorProfileID
	s.action = action
	return s.allowed, s.err
}

func TestCheckoutSubjectFromContextUsesOrganizationAuthority(t *testing.T) {
	organizationID := uint64(7)
	ctx, err := identity.BindActor(context.Background(), identity.Actor{
		ProfileID:      42,
		OrganizationID: &organizationID,
		// Global role cố ý để trống: role này không phải organization role.
	})
	if err != nil {
		t.Fatal(err)
	}
	authorizer := &organizationCheckoutAuthorizerStub{allowed: true}

	subject, err := checkoutSubjectFromContext(ctx, string(checkout.SubjectOrganization), authorizer)
	if err != nil {
		t.Fatalf("checkoutSubjectFromContext() error = %v", err)
	}
	if subject != (checkout.Subject{Kind: checkout.SubjectOrganization, ID: "7"}) {
		t.Fatalf("subject = %+v", subject)
	}
	if authorizer.organizationID != 7 || authorizer.profileID != 42 {
		t.Fatalf("authorization identity = organization:%d profile:%d", authorizer.organizationID, authorizer.profileID)
	}
	if authorizer.action != organizationpb.OrganizationAction_ORGANIZATION_ACTION_SUBSCRIPTION_CHECKOUT {
		t.Fatalf("authorization action = %v", authorizer.action)
	}
}

func TestCheckoutSubjectFromContextRejectsOrganizationWithoutPermission(t *testing.T) {
	organizationID := uint64(7)
	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42, OrganizationID: &organizationID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = checkoutSubjectFromContext(ctx, string(checkout.SubjectOrganization), &organizationCheckoutAuthorizerStub{})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("code = %v, want %v", status.Code(err), codes.PermissionDenied)
	}
}

func TestCheckoutSubjectFromContextFailsClosedWhenOrganizationUnavailable(t *testing.T) {
	organizationID := uint64(7)
	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42, OrganizationID: &organizationID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = checkoutSubjectFromContext(ctx, string(checkout.SubjectOrganization), &organizationCheckoutAuthorizerStub{err: errors.New("rpc unavailable")})
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("code = %v, want %v", status.Code(err), codes.Unavailable)
	}
}

func TestCheckoutSubjectFromContextProfileDoesNotCallOrganization(t *testing.T) {
	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42})
	if err != nil {
		t.Fatal(err)
	}
	subject, err := checkoutSubjectFromContext(ctx, string(checkout.SubjectProfile), nil)
	if err != nil {
		t.Fatalf("checkoutSubjectFromContext() error = %v", err)
	}
	if subject != (checkout.Subject{Kind: checkout.SubjectProfile, ID: "42"}) {
		t.Fatalf("subject = %+v", subject)
	}
}
