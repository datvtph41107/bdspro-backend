package grpc

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"common/identity"
	"common/request"
	organizationpb "pb/types/organization"
	userpb "pb/types/user"
	"user/internal/usecase/subscription/checkout"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscriptionCheckoutHandler struct {
	userpb.UnimplementedCheckoutServiceServer
	service                *checkout.Service
	organizationAuthorizer organizationCheckoutAuthorizer
}

type organizationCheckoutAuthorizer interface {
	AuthorizeOrganizationAction(ctx context.Context, organizationID, actorProfileID uint64, action organizationpb.OrganizationAction) (bool, error)
}

func NewSubscriptionCheckoutHandler(service *checkout.Service, organizationAuthorizer organizationCheckoutAuthorizer) *SubscriptionCheckoutHandler {
	return &SubscriptionCheckoutHandler{service: service, organizationAuthorizer: organizationAuthorizer}
}

func checkoutSubjectFromContext(ctx context.Context, requested string, authorizer organizationCheckoutAuthorizer) (checkout.Subject, error) {
	actor, ok := identity.ActorFromContext(ctx)
	if !ok {
		return checkout.Subject{}, status.Error(codes.Unauthenticated, "authenticated actor required")
	}
	switch checkout.SubjectKind(strings.TrimSpace(requested)) {
	case checkout.SubjectProfile:
		if actor.ProfileID == 0 {
			return checkout.Subject{}, status.Error(codes.Unauthenticated, "profile identity required")
		}
		return checkout.Subject{Kind: checkout.SubjectProfile, ID: strconv.FormatUint(actor.ProfileID, 10)}, nil
	case checkout.SubjectOrganization:
		if actor.OrganizationID == nil || *actor.OrganizationID == 0 {
			return checkout.Subject{}, status.Error(codes.PermissionDenied, "organization subject unavailable")
		}
		if actor.ProfileID == 0 {
			return checkout.Subject{}, status.Error(codes.Unauthenticated, "profile identity required")
		}
		if authorizer == nil {
			return checkout.Subject{}, status.Error(codes.Unavailable, "organization authorization unavailable")
		}
		allowed, err := authorizer.AuthorizeOrganizationAction(ctx, *actor.OrganizationID, actor.ProfileID, organizationpb.OrganizationAction_ORGANIZATION_ACTION_SUBSCRIPTION_CHECKOUT)
		if err != nil {
			return checkout.Subject{}, status.Error(codes.Unavailable, "organization authorization unavailable")
		}
		if !allowed {
			return checkout.Subject{}, status.Error(codes.PermissionDenied, "organization subscription checkout permission required")
		}
		return checkout.Subject{Kind: checkout.SubjectOrganization, ID: strconv.FormatUint(*actor.OrganizationID, 10)}, nil
	default:
		return checkout.Subject{}, status.Error(codes.InvalidArgument, "invalid subject kind")
	}
}
func checkoutCommandKey(ctx context.Context) (string, error) {
	v, ok := request.IdempotencyKeyFromContext(ctx)
	if !ok {
		return "", status.Error(codes.InvalidArgument, "Idempotency-Key is required")
	}
	return v, nil
}

func (h *SubscriptionCheckoutHandler) UpdateCheckoutContact(ctx context.Context, req *userpb.UpdateCheckoutContactRequest) (*userpb.CheckoutContact, error) {
	if h == nil || h.service == nil || req == nil {
		return nil, status.Error(codes.FailedPrecondition, "checkout contact unavailable")
	}
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 {
		return nil, status.Error(codes.Unauthenticated, "authenticated profile required")
	}
	contact, err := h.service.UpdateContact(ctx, actor.ProfileID, checkout.Contact{FullName: req.GetFullName(), Email: req.GetEmail()})
	if err != nil {
		return nil, checkoutError(err)
	}
	return &userpb.CheckoutContact{FullName: contact.FullName, Email: contact.Email, Phone: contact.Phone}, nil
}

func (h *SubscriptionCheckoutHandler) Checkout(ctx context.Context, req *userpb.CheckoutRequest) (*userpb.CheckoutResponse, error) {
	if h == nil || h.service == nil || req == nil {
		return nil, status.Error(codes.FailedPrecondition, "checkout unavailable")
	}
	subject, err := checkoutSubjectFromContext(ctx, req.GetSubjectKind(), h.organizationAuthorizer)
	if err != nil {
		return nil, err
	}
	key, err := checkoutCommandKey(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.service.Checkout(ctx, subject, req.GetPlanCode(), key)
	if err != nil {
		return nil, checkoutError(err)
	}
	return &userpb.CheckoutResponse{OrderId: result.Order.ID, Reference: result.Order.Reference, Status: result.Order.Status, Created: result.Created}, nil
}
func (h *SubscriptionCheckoutHandler) CreatePaymentAttempt(ctx context.Context, req *userpb.CheckoutPaymentAttemptRequest) (*userpb.CheckoutPaymentAttemptResponse, error) {
	if h == nil || h.service == nil || req == nil {
		return nil, status.Error(codes.FailedPrecondition, "payment attempt unavailable")
	}
	subject, err := checkoutSubjectFromContext(ctx, req.GetSubjectKind(), h.organizationAuthorizer)
	if err != nil {
		return nil, err
	}
	key, err := checkoutCommandKey(ctx)
	if err != nil {
		return nil, err
	}
	a, err := h.service.CreatePaymentAttempt(ctx, checkout.CreatePaymentAttemptCommand{Subject: subject, OrderID: req.GetOrderId(), Method: req.GetMethod(), CommandKey: key})
	if err != nil {
		return nil, checkoutAttemptError(err)
	}
	next := &userpb.CheckoutNextAction{Kind: a.NextActionKind, RedirectUrl: a.RedirectURL, QrPayload: a.QRPayload}
	exp := ""
	if !a.ExpiresAt.IsZero() {
		exp = a.ExpiresAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")
	}
	next.ExpiresAt = exp
	return &userpb.CheckoutPaymentAttemptResponse{AttemptId: a.ID, OrderId: a.OrderID, Method: a.Method, Provider: a.Provider, ProviderReference: a.ProviderReference, Status: a.Status, NextAction: next, ExpiresAt: exp, Created: a.Created}, nil
}
func checkoutError(err error) error {
	switch {
	case errors.Is(err, checkout.ErrInvalidCommand), errors.Is(err, checkout.ErrSubjectScope):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, checkout.ErrCheckoutCommandConflict):
		return status.Error(codes.Aborted, err.Error())
	case errors.Is(err, checkout.ErrSamePlan), errors.Is(err, checkout.ErrPendingChange), errors.Is(err, checkout.ErrDowngradeRequiresSchedule):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, checkout.ErrPlanTermsUnavailable), errors.Is(err, checkout.ErrAmbiguousPlanVersion), errors.Is(err, checkout.ErrAmbiguousRecurringPrice):
		return status.Error(codes.Unavailable, "checkout plan terms unavailable")
	default:
		return status.Error(codes.Internal, "checkout failed")
	}
}

func checkoutAttemptError(err error) error {
	if errors.Is(err, checkout.ErrInvalidCommand) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, checkout.ErrCheckoutCommandConflict) {
		return status.Error(codes.Aborted, err.Error())
	}
	return status.Error(codes.FailedPrecondition, err.Error())
}
