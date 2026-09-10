package handler

import (
	"context"
	"errors"
	"testing"

	"common/identity"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/internal/usecase/useradmin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userAdminAuthorizerStub struct {
	allowed bool
	err     error
}

func (s userAdminAuthorizerStub) HasPermission(context.Context, uint64, string) (bool, error) {
	return s.allowed, s.err
}

func TestUserAdminHandlersFailClosedBeforeBusinessLogic(t *testing.T) {
	ctx := trustedUserAdminContext(t)
	profile := &AdminUserProfileHandler{}
	admin := &AdminHandler{}
	auth := &AuthHandler{}

	calls := []struct {
		name string
		call func() error
	}{
		{"list users", func() error { _, err := profile.ListAllUsers(ctx, &userpb.ListAllUsersRequest{}); return err }},
		{"lock user", func() error { _, err := profile.LockUser(ctx, &userpb.LockUserRequest{}); return err }},
		{"unlock user", func() error { _, err := profile.UnLockUser(ctx, &userpb.UnLockUserRequest{}); return err }},
		{"approve user", func() error { _, err := profile.ApproveUser(ctx, &userpb.ApproveUserRequest{}); return err }},
		{"user stats", func() error { _, err := profile.GetUserStats(ctx, &sharepb.GetStatsRequest{}); return err }},
		{"create user", func() error { _, err := profile.CreateUser(ctx, &userpb.CreateUserRequest{}); return err }},
		{"update user", func() error { _, err := profile.UpdateUser(ctx, &userpb.UpdateUserRequest{}); return err }},
		{"delete user", func() error { _, err := profile.DeleteUser(ctx, &userpb.DeleteUserRequest{}); return err }},
		{"user detail", func() error { _, err := profile.GetUserDetail(ctx, &userpb.GetUserDetailRequest{}); return err }},
		{"list admins", func() error { _, err := admin.ListAdmins(ctx, &userpb.AdminListRequest{}); return err }},
		{"create admin", func() error { _, err := admin.CreateAdmin(ctx, &userpb.AdminProfile{}); return err }},
		{"update admin", func() error { _, err := admin.UpdateAdmin(ctx, &userpb.AdminProfile{}); return err }},
		{"delete admin", func() error { _, err := admin.DeleteAdmin(ctx, &sharepb.IdRequest{}); return err }},
		{"admin detail", func() error { _, err := admin.GetDetail(ctx, &sharepb.IdRequest{}); return err }},
		{"auth admin list", func() error { _, err := auth.ListAdmins(ctx, &authpb.AdminListRequest{}); return err }},
	}

	for _, test := range calls {
		t.Run(test.name, func(t *testing.T) {
			if code := status.Code(test.call()); code != codes.Unavailable {
				t.Fatalf("nil IAM authority code = %v, want %v", code, codes.Unavailable)
			}
		})
	}
}

func TestUserAdminAuthorizationMapsIdentityAssignmentAndAuthorityErrors(t *testing.T) {
	ctx := trustedUserAdminContext(t)
	tests := []struct {
		name       string
		ctx        context.Context
		authorizer useradmin.PermissionAuthorizer
		want       codes.Code
	}{
		{"unsigned", context.Background(), userAdminAuthorizerStub{allowed: true}, codes.Unauthenticated},
		{"denied", ctx, userAdminAuthorizerStub{}, codes.PermissionDenied},
		{"authority unavailable", ctx, userAdminAuthorizerStub{err: errors.New("database unavailable")}, codes.Unavailable},
		{"allowed", ctx, userAdminAuthorizerStub{allowed: true}, codes.OK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if code := status.Code(requireUserAdminPermission(test.ctx, test.authorizer, useradmin.PermissionView)); code != test.want {
				t.Fatalf("code = %v, want %v", code, test.want)
			}
		})
	}
}

func trustedUserAdminContext(t *testing.T) context.Context {
	t.Helper()
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "gateway-service"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindCaller(ctx, identity.Caller{Kind: identity.CallerUser})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindActor(ctx, identity.Actor{ProfileID: 42, AuthID: 7, TokenType: "ACCESS"})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}
