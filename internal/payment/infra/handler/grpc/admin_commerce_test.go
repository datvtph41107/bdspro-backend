package commercegrpc

import (
	"context"
	"errors"
	"testing"

	"common/identity"
	adminusecase "payment/internal/usecase/admin"
	paymentpb "pb/types/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type adminPermissionStub struct {
	requested []string
	err       error
}

func (s *adminPermissionStub) HasPermissions(_ context.Context, codes []string) error {
	s.requested = append([]string(nil), codes...)
	return s.err
}

func paymentAdminContext(t *testing.T) context.Context {
	t.Helper()
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "gateway-service"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindActor(ctx, identity.Actor{ProfileID: 42, TokenType: "ACCESS"})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func TestPaymentAdminRequiresTrustedActorAndPermission(t *testing.T) {
	authorizer := &adminPermissionStub{err: errors.New("denied")}
	handler := NewAdminCommerceHandler(nil, authorizer)
	_, err := handler.ListPaymentOrders(paymentAdminContext(t), &paymentpb.ListAdminPaymentOrdersRequest{})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("ListPaymentOrders() code = %v, want PermissionDenied", status.Code(err))
	}
	if len(authorizer.requested) != 1 || authorizer.requested[0] != adminusecase.PermissionOrderView {
		t.Fatalf("requested permissions = %v", authorizer.requested)
	}

	_, err = handler.ListPaymentOrders(context.Background(), &paymentpb.ListAdminPaymentOrdersRequest{})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("untrusted ListPaymentOrders() code = %v", status.Code(err))
	}
}

func TestPaymentRedriveRequiresDedicatedPermissionBeforeCommand(t *testing.T) {
	authorizer := &adminPermissionStub{err: errors.New("denied")}
	handler := NewAdminCommerceHandler(nil, authorizer)
	_, err := handler.RedriveFulfillment(paymentAdminContext(t), &paymentpb.RedriveAdminFulfillmentRequest{FulfillmentId: 3})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("RedriveFulfillment() code = %v", status.Code(err))
	}
	if len(authorizer.requested) != 1 || authorizer.requested[0] != adminusecase.PermissionFulfillmentRedrive {
		t.Fatalf("requested permissions = %v", authorizer.requested)
	}
}

func TestPaymentConfirmationRequiresDedicatedPermissionBeforeCommand(t *testing.T) {
	authorizer := &adminPermissionStub{err: errors.New("denied")}
	handler := NewAdminCommerceHandler(nil, authorizer)
	_, err := handler.ConfirmPaymentFunds(paymentAdminContext(t), &paymentpb.ConfirmAdminPaymentFundsRequest{OrderId: 3, Reason: "bank receipt checked"})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("ConfirmPaymentFunds() code = %v", status.Code(err))
	}
	if len(authorizer.requested) != 1 || authorizer.requested[0] != adminusecase.PermissionFundsConfirm {
		t.Fatalf("requested permissions = %v", authorizer.requested)
	}
}
