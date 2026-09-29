package commercialprofilegrpc

import (
	"context"
	"testing"
	"time"

	"common/identity"
	payment "payment/internal/domain/payment"
	commercialprofile "payment/internal/usecase/commercialprofile"
	paymentpb "pb/types/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repositoryStub struct{ item payment.AdminOrder }

func (s repositoryStub) GetAdminOrder(context.Context, uint64) (payment.AdminOrder, error) {
	return s.item, nil
}

func TestGetMyPaymentOrderRequiresAccessActor(t *testing.T) {
	handler := New(commercialprofile.NewService(repositoryStub{}))
	_, err := handler.GetMyPaymentOrder(context.Background(), &paymentpb.GetMyPaymentOrderRequest{OrderId: 9})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("error = %v, want Unauthenticated", err)
	}
}

func TestGetMyPaymentOrderReturnsOwnerProgress(t *testing.T) {
	now := time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)
	handler := New(commercialprofile.NewService(repositoryStub{item: payment.AdminOrder{
		Order: payment.Order{
			ID: 9, Reference: "QHPRO-9", Subject: payment.Subject{Kind: payment.SubjectProfile, ID: "42"},
			Terms:  payment.CommercialTerms{ProductCode: "QHPRO", PlanCode: "BASIC", Price: payment.Money{Currency: payment.CurrencyVND, AmountMinor: 99000}},
			Status: payment.OrderFundsConfirmed, FundsConfirmedAt: &now, UpdatedAt: now,
		},
		Fulfillment: &payment.Fulfillment{Status: payment.FulfillmentRetry},
	}}))
	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42, TokenType: "ACCESS"})
	if err != nil {
		t.Fatal(err)
	}
	response, err := handler.GetMyPaymentOrder(ctx, &paymentpb.GetMyPaymentOrderRequest{OrderId: 9})
	if err != nil {
		t.Fatal(err)
	}
	if response.GetProgressState() != commercialprofile.ProgressActivating ||
		response.GetOrderStatus() != string(payment.OrderFundsConfirmed) ||
		response.GetFulfillmentStatus() != string(payment.FulfillmentRetry) {
		t.Fatalf("response = %+v", response)
	}
}

func TestGetMyPaymentOrderHidesAnotherProfile(t *testing.T) {
	handler := New(commercialprofile.NewService(repositoryStub{item: payment.AdminOrder{
		Order: payment.Order{ID: 9, Subject: payment.Subject{Kind: payment.SubjectProfile, ID: "7"}},
	}}))
	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42, TokenType: "ACCESS"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = handler.GetMyPaymentOrder(ctx, &paymentpb.GetMyPaymentOrderRequest{OrderId: 9})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("error = %v, want PermissionDenied", err)
	}
}
