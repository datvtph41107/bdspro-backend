package admin

import (
	"context"
	"testing"

	payment "payment/internal/domain/payment"
)

type repositoryStub struct {
	orderQuery       OrderQuery
	fulfillmentQuery FulfillmentQuery
	redriveCommand   payment.RedriveCommand
	changed          bool
	replay           bool
	status           payment.FulfillmentStatus
	err              error
}

func (r *repositoryStub) ListAdminOrders(_ context.Context, query OrderQuery) (OrderPage, error) {
	r.orderQuery = query
	return OrderPage{Page: query.Page, PageSize: query.PageSize}, r.err
}
func (r *repositoryStub) GetAdminOrder(context.Context, uint64) (payment.AdminOrder, error) {
	return payment.AdminOrder{}, r.err
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
	service := NewService(repository)
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
	service := NewService(repository)
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
	service := NewService(&repositoryStub{})
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
