package admin

import (
	"context"
	"strings"
	"testing"
	"time"

	payment "payment/internal/domain/payment"
	"payment/internal/usecase/settlement"
)

type repositoryStub struct {
	orderQuery       OrderQuery
	fulfillmentQuery FulfillmentQuery
	redriveCommand   payment.RedriveCommand
	changed          bool
	replay           bool
	status           payment.FulfillmentStatus
	err              error
	order            payment.AdminOrder
}

func (r *repositoryStub) ListAdminOrders(_ context.Context, query OrderQuery) (OrderPage, error) {
	r.orderQuery = query
	return OrderPage{Page: query.Page, PageSize: query.PageSize}, r.err
}
func (r *repositoryStub) GetAdminOrder(context.Context, uint64) (payment.AdminOrder, error) {
	return r.order, r.err
}
func (r *repositoryStub) ListAdminFulfillments(_ context.Context, query FulfillmentQuery) (FulfillmentPage, error) {
	r.fulfillmentQuery = query
	return FulfillmentPage{Page: query.Page, PageSize: query.PageSize}, r.err
}
func (r *repositoryStub) RedriveAdminFulfillment(_ context.Context, command payment.RedriveCommand) (bool, bool, payment.FulfillmentStatus, error) {
	r.redriveCommand = command
	return r.changed, r.replay, r.status, r.err
}

func TestListOrdersNormalizesAdminQuery(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository, nil, nil)
	page, err := service.ListOrders(context.Background(), OrderQuery{SubjectID: " 42 ", Reference: " PAY- "})
	if err != nil {
		t.Fatalf("ListOrders() error = %v", err)
	}
	if page.Page != 1 || page.PageSize != 20 || repository.orderQuery.SubjectID != "42" || repository.orderQuery.Reference != "PAY-" {
		t.Fatalf("normalized query = %+v, page = %+v", repository.orderQuery, page)
	}
}

func TestRedrivePreservesIdempotencyEvidence(t *testing.T) {
	repository := &repositoryStub{changed: true, status: payment.FulfillmentPending}
	service := NewService(repository, nil, nil)
	changed, replay, current, err := service.Redrive(context.Background(), RedriveCommand{
		FulfillmentID: 9, CommandKey: " command-1 ", ActorID: " 42 ", Reason: " dependency recovered ",
	})
	if err != nil {
		t.Fatalf("Redrive() error = %v", err)
	}
	if !changed || replay || current != payment.FulfillmentPending {
		t.Fatalf("Redrive() = %v, %v, %q", changed, replay, current)
	}
	if repository.redriveCommand.CommandKey != "command-1" || repository.redriveCommand.ActorID != "42" || repository.redriveCommand.Reason != "dependency recovered" {
		t.Fatalf("redrive command = %+v", repository.redriveCommand)
	}
}

func TestRedriveRejectsMissingActorOrCommandKey(t *testing.T) {
	service := NewService(&repositoryStub{}, nil, nil)
	for _, command := range []RedriveCommand{
		{FulfillmentID: 1, ActorID: "42", Reason: "reason"},
		{FulfillmentID: 1, CommandKey: "key", Reason: "reason"},
		{CommandKey: "key", ActorID: "42", Reason: "reason"},
		{FulfillmentID: 1, CommandKey: "key", ActorID: "42"},
	} {
		if _, _, _, err := service.Redrive(context.Background(), command); err != payment.ErrInvalidCommand {
			t.Fatalf("Redrive(%+v) error = %v", command, err)
		}
	}
}

type operatorSettlementStub struct {
	event  payment.ProviderEvent
	effect payment.CommandEffect
	result settlement.Result
	err    error
}

func (s *operatorSettlementStub) RecordOperatorFundsConfirmation(_ context.Context, event payment.ProviderEvent, effect payment.CommandEffect) (settlement.Result, error) {
	s.event, s.effect = event, effect
	return s.result, s.err
}

func TestConfirmFundsUsesFrozenOrderTermsAndOperatorEvidence(t *testing.T) {
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	order := payment.AdminOrder{
		Order:    payment.Order{ID: 7, Reference: "QHP-REF", Status: payment.OrderPendingFunds, ExpiresAt: now.Add(time.Hour), Terms: payment.CommercialTerms{Price: payment.Money{Currency: payment.CurrencyVND, AmountMinor: 99000}}},
		Attempts: []payment.PaymentAttempt{{ID: 1, OrderID: 7}},
	}
	repository := &repositoryStub{order: order}
	recorder := &operatorSettlementStub{}
	service := NewService(repository, recorder, func() time.Time { return now })
	result, err := service.ConfirmFunds(context.Background(), ConfirmFundsCommand{OrderID: 7, CommandKey: " command-1 ", ActorID: " 42 ", Reason: " bank receipt checked "})
	if err != nil {
		t.Fatalf("ConfirmFunds() error = %v", err)
	}
	if !result.Changed || result.Replay {
		t.Fatalf("ConfirmFunds() result = %+v", result)
	}
	if recorder.event.Reference != "QHP-REF" || recorder.event.Amount.AmountMinor != 99000 || recorder.event.Amount.Currency != payment.CurrencyVND {
		t.Fatalf("provider event did not use frozen order terms: %+v", recorder.event)
	}
	if recorder.effect.ActorID != "42" || recorder.effect.Reason != "bank receipt checked" || recorder.effect.CommandKey != "command-1" {
		t.Fatalf("operator effect = %+v", recorder.effect)
	}
	if !strings.Contains(recorder.event.TransactionID, "order.7.command-1") {
		t.Fatalf("transaction id = %q", recorder.event.TransactionID)
	}
}

func TestConfirmFundsRejectsOrderWithoutPaymentAttempt(t *testing.T) {
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	service := NewService(&repositoryStub{order: payment.AdminOrder{Order: payment.Order{ID: 7, Status: payment.OrderPendingFunds, ExpiresAt: now.Add(time.Hour)}}}, &operatorSettlementStub{}, func() time.Time { return now })
	_, err := service.ConfirmFunds(context.Background(), ConfirmFundsCommand{OrderID: 7, CommandKey: "key", ActorID: "42", Reason: "receipt"})
	if err != payment.ErrFundsConfirmationNotAllowed {
		t.Fatalf("ConfirmFunds() error = %v", err)
	}
}
