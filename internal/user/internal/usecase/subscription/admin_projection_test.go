package subscription

import (
	"context"
	"errors"
	"testing"

	commonidentity "common/identity"
	subscriptiondomain "user/internal/domain/subscription"
)

type adminRepositoryStub struct {
	rows []subscriptiondomain.AdminUserProjection
}

type permissionAuthorizerStub struct {
	allowed bool
	err     error
}

func (s permissionAuthorizerStub) HasPermission(context.Context, uint64, string) (bool, error) {
	return s.allowed, s.err
}

func (r adminRepositoryStub) List(context.Context, AdminQuery) (AdminPage, error) {
	return AdminPage{}, nil
}
func (r adminRepositoryStub) ListEffectiveByProfiles(context.Context, []uint64) ([]subscriptiondomain.AdminUserProjection, error) {
	return r.rows, nil
}

func TestListUsersPreservesRequestOrderAndExplicitMissingSubscription(t *testing.T) {
	current := subscriptiondomain.AdminProjection{SubscriptionID: 9, SubjectKind: "profile", SubjectID: "2", PlanCode: "PRO"}
	service := NewAdminService(adminRepositoryStub{rows: []subscriptiondomain.AdminUserProjection{{ProfileID: 2, EffectiveSubscription: &current}}})
	result, err := service.ListUsers(context.Background(), []uint64{3, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 || result[0].ProfileID != 3 || result[0].EffectiveSubscription != nil {
		t.Fatalf("missing profile projection was fabricated: %+v", result)
	}
	if result[1].ProfileID != 2 || result[1].EffectiveSubscription == nil || result[1].EffectiveSubscription.PlanCode != "PRO" {
		t.Fatalf("effective projection lost: %+v", result)
	}
}

func TestListUsersRejectsInvalidBatch(t *testing.T) {
	service := NewAdminService(adminRepositoryStub{})
	if _, err := service.ListUsers(context.Background(), nil); err == nil {
		t.Fatal("empty batch must fail")
	}
	if _, err := service.ListUsers(context.Background(), []uint64{0}); err == nil {
		t.Fatal("zero profile id must fail")
	}
}

func trustedAdminContext(t *testing.T) context.Context {
	t.Helper()
	ctx, err := commonidentity.BindServiceCaller(context.Background(), commonidentity.ServiceCaller{ServiceID: "gateway-service"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = commonidentity.BindActor(ctx, commonidentity.Actor{ProfileID: 42, TokenType: "ACCESS"})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func TestRequireViewPermissionSeparatesTrustActorAndIAM(t *testing.T) {
	if _, err := RequireViewPermission(context.Background(), permissionAuthorizerStub{allowed: true}); !errors.Is(err, ErrAdminActorRequired) {
		t.Fatalf("untrusted request error = %v", err)
	}
	ctx := trustedAdminContext(t)
	if _, err := RequireViewPermission(ctx, nil); !errors.Is(err, ErrPermissionAuthorityUnavailable) {
		t.Fatalf("missing authority error = %v", err)
	}
	if _, err := RequireViewPermission(ctx, permissionAuthorizerStub{}); !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("denied assignment error = %v", err)
	}
	actorID, err := RequireViewPermission(ctx, permissionAuthorizerStub{allowed: true})
	if err != nil || actorID != 42 {
		t.Fatalf("allowed actor = %d, %v", actorID, err)
	}
}

type permissionAuthorizerSpy struct {
	allowed        bool
	err            error
	calls          int
	lastActorID    uint64
	lastPermission string
}

func (s *permissionAuthorizerSpy) HasPermission(
	_ context.Context,
	actorID uint64,
	permissionCode string,
) (bool, error) {
	s.calls++
	s.lastActorID = actorID
	s.lastPermission = permissionCode
	return s.allowed, s.err
}

func serviceCallerContext(t *testing.T) context.Context {
	t.Helper()

	ctx, err := commonidentity.BindServiceCaller(
		context.Background(),
		commonidentity.ServiceCaller{
			ServiceID: "gateway-service",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return ctx
}

func TestRequireViewPermissionKeepsTransportActorAndIAMIndependent(
	t *testing.T,
) {
	t.Run(
		"actor without verified transport caller",
		func(t *testing.T) {
			ctx, err := commonidentity.BindActor(
				context.Background(),
				commonidentity.Actor{
					ProfileID: 42,
					TokenType: "ACCESS",
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			authorizer :=
				&permissionAuthorizerSpy{allowed: true}

			_, err = RequireViewPermission(
				ctx,
				authorizer,
			)
			if !errors.Is(
				err,
				ErrAdminActorRequired,
			) {
				t.Fatalf(
					"error = %v, want %v",
					err,
					ErrAdminActorRequired,
				)
			}
			if authorizer.calls != 0 {
				t.Fatalf(
					"IAM called before transport trust: %d",
					authorizer.calls,
				)
			}
		},
	)

	t.Run(
		"verified transport without actor",
		func(t *testing.T) {
			authorizer :=
				&permissionAuthorizerSpy{allowed: true}

			_, err := RequireViewPermission(
				serviceCallerContext(t),
				authorizer,
			)
			if !errors.Is(
				err,
				ErrAdminActorRequired,
			) {
				t.Fatalf(
					"error = %v, want %v",
					err,
					ErrAdminActorRequired,
				)
			}
			if authorizer.calls != 0 {
				t.Fatalf(
					"IAM called before actor identity: %d",
					authorizer.calls,
				)
			}
		},
	)

	t.Run(
		"non access actor",
		func(t *testing.T) {
			ctx := serviceCallerContext(t)

			ctx, err :=
				commonidentity.BindActor(
					ctx,
					commonidentity.Actor{
						ProfileID: 42,
						TokenType: "REFRESH",
					},
				)
			if err != nil {
				t.Fatal(err)
			}

			authorizer :=
				&permissionAuthorizerSpy{allowed: true}

			_, err =
				RequireViewPermission(
					ctx,
					authorizer,
				)
			if !errors.Is(
				err,
				ErrAdminActorRequired,
			) {
				t.Fatalf(
					"error = %v, want %v",
					err,
					ErrAdminActorRequired,
				)
			}
			if authorizer.calls != 0 {
				t.Fatalf(
					"IAM called for non-ACCESS actor: %d",
					authorizer.calls,
				)
			}
		},
	)

	t.Run(
		"valid actor reaches exact IAM assignment",
		func(t *testing.T) {
			ctx := trustedAdminContext(t)
			authorizer :=
				&permissionAuthorizerSpy{
					allowed: true,
				}

			actorID, err :=
				RequireViewPermission(
					ctx,
					authorizer,
				)
			if err != nil {
				t.Fatal(err)
			}

			if actorID != 42 ||
				authorizer.calls != 1 ||
				authorizer.lastActorID != 42 ||
				authorizer.lastPermission != PermissionView {
				t.Fatalf(
					"IAM call = actor:%d calls:%d permission:%q returned actor:%d",
					authorizer.lastActorID,
					authorizer.calls,
					authorizer.lastPermission,
					actorID,
				)
			}
		},
	)
}
