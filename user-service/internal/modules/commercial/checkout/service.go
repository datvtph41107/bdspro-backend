package checkout

import (
	"context"
	"strings"
)

type Service struct {
	catalog       CatalogResolver
	subscriptions SubscriptionReader
	commands      CommandStore
	payment       PaymentOrderPort
	now           Clock
}

func NewService(
	catalog CatalogResolver,
	subscriptions SubscriptionReader,
	commands CommandStore,
	payment PaymentOrderPort,
	now Clock,
) *Service {
	return &Service{
		catalog:       catalog,
		subscriptions: subscriptions,
		commands:      commands,
		payment:       payment,
		now:           now,
	}
}

func (s *Service) Checkout(ctx context.Context, subject Subject, planCode, commandKey string) (Result, error) {
	planCode = strings.TrimSpace(planCode)
	commandKey = strings.TrimSpace(commandKey)
	if s == nil || s.catalog == nil || s.subscriptions == nil || s.commands == nil ||
		s.payment == nil || s.now == nil || !subject.IsValid() || planCode == "" || commandKey == "" {
		return Result{}, ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	// Durable replay is checked before reading current Catalog/Subscription state.
	// A successful first Checkout freezes the server-owned commercial snapshot;
	// network retry must not become a different command because a successor plan
	// became effective or the first payment already changed the Subscription.
	existing, found, err := s.commands.FindByCommand(ctx, subject, commandKey)
	if err != nil {
		return Result{}, err
	}
	if found {
		return s.replay(ctx, subject, planCode, commandKey, existing)
	}

	at := s.now().UTC()
	terms, err := s.catalog.ResolveCheckoutTerms(ctx, planCode, at)
	if err != nil {
		return Result{}, err
	}
	if !terms.IsValid() {
		return Result{}, ErrCommercialTermsUnavailable
	}
	if !terms.Allows(subject) {
		return Result{}, ErrSubjectScope
	}

	current, found, err := s.subscriptions.CurrentForProduct(ctx, subject, terms.ProductID)
	if err != nil {
		return Result{}, err
	}
	if found {
		if current.ID == 0 || current.PlanVersionID == 0 || current.TierRank <= 0 {
			return Result{}, ErrCommercialTermsUnavailable
		}
		if current.Status == "pending" || current.PendingPlanVersionID != 0 {
			return Result{}, ErrPendingChange
		}
		change, err := ClassifyPlanChange(current.TierRank, terms.TierRank)
		if err != nil {
			return Result{}, err
		}
		switch change {
		case PlanChangeSame:
			return Result{}, ErrSamePlan
		case PlanChangeDowngrade:
			return Result{}, ErrDowngradeRequiresSchedule
		case PlanChangeUpgrade:
			// Continue to persist the frozen Checkout command.
		default:
			return Result{}, ErrCommercialTermsUnavailable
		}
	}

	command := CheckoutCommand{
		Subject:    subject,
		PlanCode:   planCode,
		Terms:      terms,
		CommandKey: commandKey,
		CreatedAt:  at,
	}
	if !command.IsValid() {
		return Result{}, ErrCommercialTermsUnavailable
	}

	// CreateCommand is an idempotent arbitration point. On a concurrent retry,
	// the winner's frozen terms become authoritative for every caller sharing
	// the same Subject + CommandKey.
	saved, _, err := s.commands.CreateCommand(ctx, command)
	if err != nil {
		return Result{}, err
	}
	return s.replay(ctx, subject, planCode, commandKey, saved)
}

func (s *Service) replay(
	ctx context.Context,
	subject Subject,
	planCode string,
	commandKey string,
	command CheckoutCommand,
) (Result, error) {
	if !command.IsValid() || command.Subject != subject ||
		strings.TrimSpace(command.PlanCode) != planCode ||
		strings.TrimSpace(command.CommandKey) != commandKey {
		return Result{}, ErrCheckoutCommandConflict
	}

	order, created, err := s.payment.CreateOrder(ctx, CreatePaymentOrderCommand{
		Subject:    subject,
		Terms:      command.Terms,
		CommandKey: commandKey,
	})
	if err != nil {
		return Result{}, err
	}
	return Result{Order: order, Created: created, Terms: command.Terms}, nil
}
