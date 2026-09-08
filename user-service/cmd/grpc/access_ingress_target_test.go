package main

import (
	"common/identity"
	"common/rpc"
	"context"
	authpb "pb/types/auth"
	userpb "pb/types/user"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const adminCommercialListMethod = "/userpb.AdminCommercialService/ListSubscriptions"

func adminCommercialAssertionConfig(now time.Time) rpc.ServiceAssertionConfig {
	return rpc.ServiceAssertionConfig{
		VerificationSecrets: map[string]string{
			"gateway-service": "test-secret",
		},
		VerificationKeysConfigured: true,
		Now: func() time.Time {
			return now
		},
	}
}

func signedAdminCommercialContext(
	t *testing.T,
	now time.Time,
) context.Context {
	t.Helper()

	md := metadata.Pairs(
		rpc.CallerKindMetadataKey,
		string(identity.CallerUser),
		rpc.ProfileIDMetadataKey,
		"42",
		rpc.TokenTypeMetadataKey,
		"ACCESS",
	)

	err := rpc.SignServiceAssertion(
		adminCommercialListMethod,
		md,
		rpc.ServiceAssertionConfig{
			ServiceID: "gateway-service",
			Secret:    "test-secret",
			Now: func() time.Time {
				return now
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return metadata.NewIncomingContext(
		context.Background(),
		md,
	)
}

func TestAccessIngressRejectsUnsignedAdminCommercialCall(
	t *testing.T,
) {
	now := time.Unix(1_700_000_000, 0)

	called := false
	_, err := accessIngress(
		rpc.TransportConfig{
			ServiceAssertion: adminCommercialAssertionConfig(now),
		},
	)(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: adminCommercialListMethod,
		},
		func(
			context.Context,
			interface{},
		) (interface{}, error) {
			called = true
			return nil, nil
		},
	)

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf(
			"unsigned call code = %v, want %v (err=%v)",
			status.Code(err),
			codes.Unauthenticated,
			err,
		)
	}
	if called {
		t.Fatal(
			"unsigned AdminCommercial request reached handler",
		)
	}
}

func TestAccessIngressRejectsTamperedAdminCommercialIdentity(
	t *testing.T,
) {
	now := time.Unix(1_700_000_000, 0)

	ctx := signedAdminCommercialContext(t, now)
	md, _ := metadata.FromIncomingContext(ctx)
	md = md.Copy()

	// Actor identity is inside the signed metadata contract.
	// Changing it after signing must invalidate the transport hop.
	md.Set(rpc.ProfileIDMetadataKey, "99")

	ctx = metadata.NewIncomingContext(
		context.Background(),
		md,
	)

	called := false
	_, err := accessIngress(
		rpc.TransportConfig{
			ServiceAssertion: adminCommercialAssertionConfig(now),
		},
	)(
		ctx,
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: adminCommercialListMethod,
		},
		func(
			context.Context,
			interface{},
		) (interface{}, error) {
			called = true
			return nil, nil
		},
	)

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf(
			"tampered call code = %v, want %v (err=%v)",
			status.Code(err),
			codes.Unauthenticated,
			err,
		)
	}
	if called {
		t.Fatal(
			"tampered AdminCommercial request reached handler",
		)
	}
}

func TestAccessIngressPromotesSignedAdminCommercialIdentity(
	t *testing.T,
) {
	now := time.Unix(1_700_000_000, 0)
	ctx := signedAdminCommercialContext(t, now)

	called := false
	_, err := accessIngress(
		rpc.TransportConfig{
			ServiceAssertion: adminCommercialAssertionConfig(now),
		},
	)(
		ctx,
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: adminCommercialListMethod,
		},
		func(
			ctx context.Context,
			_ interface{},
		) (interface{}, error) {
			called = true

			caller, ok :=
				identity.ServiceCallerFromContext(ctx)
			if !ok ||
				caller.ServiceID != "gateway-service" {
				t.Fatalf(
					"service caller = %+v, %v",
					caller,
					ok,
				)
			}

			actor, ok :=
				identity.ActorFromContext(ctx)
			if !ok ||
				actor.ProfileID != 42 ||
				actor.TokenType != "ACCESS" {
				t.Fatalf(
					"actor = %+v, %v",
					actor,
					ok,
				)
			}

			return nil, nil
		},
	)

	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal(
			"signed AdminCommercial request did not reach handler",
		)
	}
}

func TestTrustedIdentityMethodIncludesAllAdminCommercialRPCs(
	t *testing.T,
) {
	methods := []string{
		"/userpb.AdminCommercialService/ListSubscriptions",
		"/userpb.AdminCommercialService/GetUserProjection",
		"/userpb.AdminCommercialService/ListUserProjections",
	}

	for _, method := range methods {
		if !trustedIdentityMethod(method) {
			t.Fatalf(
				"AdminCommercial method is not strict identity routed: %s",
				method,
			)
		}
	}
}

func TestTrustedIdentityMethodIncludesAllUserAdminRPCs(t *testing.T) {
	methods := []string{
		userpb.AdminUserProfileService_ListAllUsers_FullMethodName,
		userpb.AdminUserProfileService_LockUser_FullMethodName,
		userpb.AdminUserProfileService_UnLockUser_FullMethodName,
		userpb.AdminUserProfileService_ApproveUser_FullMethodName,
		userpb.AdminUserProfileService_SendWarningToUser_FullMethodName,
		userpb.AdminUserProfileService_GetUserStats_FullMethodName,
		userpb.AdminUserProfileService_CreateUser_FullMethodName,
		userpb.AdminUserProfileService_UpdateUser_FullMethodName,
		userpb.AdminUserProfileService_DeleteUser_FullMethodName,
		userpb.AdminUserProfileService_GetUserDetail_FullMethodName,
		userpb.AdminProfileService_ListAdmins_FullMethodName,
		userpb.AdminProfileService_CreateAdmin_FullMethodName,
		userpb.AdminProfileService_UpdateAdmin_FullMethodName,
		userpb.AdminProfileService_DeleteAdmin_FullMethodName,
		userpb.AdminProfileService_GetDetail_FullMethodName,
		authpb.AuthService_ListAdmins_FullMethodName,
	}
	for _, method := range methods {
		if !trustedIdentityMethod(method) {
			t.Fatalf("User Admin method is not strict identity routed: %s", method)
		}
	}
}

func TestTrustedIdentityMethodIncludesAdminOrganizationRPCs(t *testing.T) {
	methods := []string{
		"/userpb.AdminOrganizationService/ListOrganizations",
		"/userpb.AdminOrganizationService/GetOrganization",
		"/userpb.AdminOrganizationService/CreateOrganization",
		"/userpb.AdminOrganizationService/UpdateOrganization",
		"/userpb.AdminOrganizationService/ArchiveOrganization",
	}
	for _, method := range methods {
		if !trustedIdentityMethod(method) {
			t.Fatalf("Admin Organization method is not strict identity routed: %s", method)
		}
	}
}
