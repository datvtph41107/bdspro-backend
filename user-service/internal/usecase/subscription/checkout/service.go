package checkout

import (
	"context"
	"strings"
)

type Service struct {
	plans         PlanResolver
	subscriptions SubscriptionReader
	commands      CommandStore
	payment       PaymentOrderPort
	attempts      PaymentAttemptPort
	contacts      ContactStore
	now           Clock
}

func NewService(
	plans PlanResolver,
	subscriptions SubscriptionReader,
	commands CommandStore,
	payment PaymentOrderPort,
	attempts PaymentAttemptPort,
	contacts ContactStore,
	now Clock,
) *Service {
	return &Service{
		plans:         plans,
		subscriptions: subscriptions,
		commands:      commands,
		payment:       payment,
		attempts:      attempts,
		contacts:      contacts,
		now:           now,
	}
}

func (s *Service) UpdateContact(ctx context.Context, profileID uint64, contact Contact) (Contact, error) {
	contact.FullName = strings.TrimSpace(contact.FullName)
	contact.Email = strings.ToLower(strings.TrimSpace(contact.Email))
	if s == nil || s.contacts == nil || profileID == 0 || !contact.IsValidInput() {
		return Contact{}, ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return Contact{}, err
	}
	return s.contacts.UpdateCheckoutContact(ctx, profileID, contact)
}

func (s *Service) Checkout(ctx context.Context, subject Subject, planCode, commandKey string) (Result, error) {
	planCode = strings.TrimSpace(planCode)
	commandKey = strings.TrimSpace(commandKey)
	if s == nil || s.plans == nil || s.subscriptions == nil || s.commands == nil ||
		s.payment == nil || s.now == nil || !subject.IsValid() || planCode == "" || commandKey == "" {
		return Result{}, ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	requestHash := RequestHashFor(subject, planCode)

	// Durable replay is checked before reading current Plan/Subscription state.
	// A successful first Checkout freezes the server-owned plan snapshot;
	// network retry must not become a different command because a successor plan
	// became effective or the first payment already changed the Subscription.
	existing, found, err := s.commands.FindByCommand(ctx, subject, commandKey)
	if err != nil {
		return Result{}, err
	}
	if found {
		return s.replay(ctx, subject, planCode, requestHash, commandKey, existing)
	}

	at := s.now().UTC()
	terms, err := s.plans.ResolveCheckoutTerms(ctx, planCode, at)
	if err != nil {
		return Result{}, err
	}
	if !terms.IsValid() {
		return Result{}, ErrPlanTermsUnavailable
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
			return Result{}, ErrPlanTermsUnavailable
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
			return Result{}, ErrPlanTermsUnavailable
		}
	}

	command := CheckoutCommand{
		Subject:     subject,
		PlanCode:    planCode,
		RequestHash: requestHash,
		Terms:       terms,
		CommandKey:  commandKey,
		CreatedAt:   at,
	}
	if !command.IsValid() {
		return Result{}, ErrPlanTermsUnavailable
	}

	// CreateCommand is an idempotent arbitration point. On a concurrent retry,
	// the winner's frozen terms become authoritative for every caller sharing
	// the same Subject + CommandKey.
	saved, _, err := s.commands.CreateCommand(ctx, command)
	if err != nil {
		return Result{}, err
	}
	return s.replay(ctx, subject, planCode, requestHash, commandKey, saved)
}

func (s *Service) replay(
	ctx context.Context,
	subject Subject,
	planCode string,
	requestHash string,
	commandKey string,
	command CheckoutCommand,
) (Result, error) {
	if !command.IsValid() || command.Subject != subject ||
		strings.TrimSpace(command.PlanCode) != planCode ||
		strings.TrimSpace(command.RequestHash) != requestHash ||
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
