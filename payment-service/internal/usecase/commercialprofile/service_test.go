package commercialprofile

import (
	"context"
	"errors"
	"testing"

	payment "payment/internal/domain/payment"
)

type repositoryStub struct {
	item payment.AdminOrder
	err  error
}

func (s repositoryStub) GetAdminOrder(context.Context, uint64) (payment.AdminOrder, error) {
	return s.item, s.err
}

func TestGetSeparatesFundsFromActivation(t *testing.T) {
	service := NewService(repositoryStub{item: payment.AdminOrder{
		Order:       payment.Order{ID: 9, Subject: payment.Subject{Kind: payment.SubjectProfile, ID: "42"}, Status: payment.OrderFundsConfirmed},
		Settlement:  &payment.Settlement{Status: payment.SettlementFundsConfirmed},
		Fulfillment: &payment.Fulfillment{Status: payment.FulfillmentRetry},
	}})
	projection, err := service.Get(context.Background(), Viewer{ProfileID: 42}, 9)
	if err != nil {
		t.Fatal(err)
	}
	if projection.ProgressState != ProgressActivating || projection.SettlementStatus != "funds_confirmed" || projection.FulfillmentStatus != "retry" {
		t.Fatalf("projection = %+v", projection)
	}
}

func TestGetRequiresOwningSubject(t *testing.T) {
	service := NewService(repositoryStub{item: payment.AdminOrder{
		Order: payment.Order{ID: 9, Subject: payment.Subject{Kind: payment.SubjectProfile, ID: "42"}},
	}})
	_, err := service.Get(context.Background(), Viewer{ProfileID: 7}, 9)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v, want ErrForbidden", err)
	}
}

func TestGetAllowsCurrentOrganizationSubject(t *testing.T) {
	organizationID := uint64(91)
	service := NewService(repositoryStub{item: payment.AdminOrder{
		Order: payment.Order{ID: 9, Subject: payment.Subject{Kind: payment.SubjectOrganization, ID: "91"}, Status: payment.OrderPendingFunds},
	}})
	projection, err := service.Get(context.Background(), Viewer{ProfileID: 42, OrganizationID: &organizationID}, 9)
	if err != nil || projection.ProgressState != ProgressAwaitingPayment {
		t.Fatalf("projection = %+v, error = %v", projection, err)
	}
}
