package checkout

import (
	"context"
	"time"
)

type Clock func() time.Time

type CatalogResolver interface {
	ResolveCheckoutTerms(ctx context.Context, planCode string, at time.Time) (PlanTerms, error)
}

type SubscriptionReader interface {
	CurrentForProduct(ctx context.Context, subject Subject, productID uint64) (CurrentSubscription, bool, error)
}

// CommandStore owns the durable User-side resolution of a Checkout command.
// Persisting the frozen server terms before calling Payment makes a network
// replay independent from later Catalog or Subscription state changes.
type CommandStore interface {
	FindByCommand(ctx context.Context, subject Subject, commandKey string) (CheckoutCommand, bool, error)
	CreateCommand(ctx context.Context, command CheckoutCommand) (CheckoutCommand, bool, error)
}

type CreatePaymentOrderCommand struct {
	Subject    Subject
	Terms      PlanTerms
	CommandKey string
}

// PaymentOrderPort is a business command boundary. Transport, Payment address,
// protobuf and timeout policy belong to a later adapter/runtime checkpoint.
type PaymentOrderPort interface {
	CreateOrder(ctx context.Context, command CreatePaymentOrderCommand) (PaymentOrder, bool, error)
}
