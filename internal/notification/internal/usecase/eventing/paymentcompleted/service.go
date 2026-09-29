package paymentcompleted

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	sharedevent "common/events/paymentcompleted"
	"notification/internal/domain/eventing"
)

var ErrInvalidEvent = errors.New("invalid payment completed event")

type ReactionStore interface {
	AcceptOnce(ctx context.Context, fact eventing.PaymentCompleted) (replay bool, err error)
}

type Service struct{ store ReactionStore }

func NewService(store ReactionStore) *Service { return &Service{store: store} }

func (s *Service) Handle(ctx context.Context, raw sharedevent.V1) (bool, error) {
	if s == nil || s.store == nil || !raw.IsValid() {
		return false, ErrInvalidEvent
	}
	fact := eventing.PaymentCompleted{
		EventID: strings.TrimSpace(raw.EventID), OrderID: raw.OrderID,
		SubjectKind: strings.TrimSpace(raw.SubjectKind), SubjectID: strings.TrimSpace(raw.SubjectID),
		ProductCode: strings.TrimSpace(raw.ProductCode), PlanCode: strings.TrimSpace(raw.PlanCode),
		PlanVersionID: raw.PlanVersionID, PlanVersion: strings.TrimSpace(raw.PlanVersion),
		Currency: strings.ToUpper(strings.TrimSpace(raw.Currency)), AmountMinor: raw.AmountMinor,
		FundsConfirmedAt: raw.FundsConfirmedAt.UTC(), CompletedAt: raw.CompletedAt.UTC(),
	}
	return s.store.AcceptOnce(ctx, fact)
}

func NotificationPayload(f eventing.PaymentCompleted) string {
	return fmt.Sprintf("%s|%s|%d|%s", f.PlanCode, f.Currency, f.AmountMinor, strconv.FormatUint(f.OrderID, 10))
}

func ShouldCreatePushIntent(f eventing.PaymentCompleted) bool {
	return strings.EqualFold(strings.TrimSpace(f.SubjectKind), "profile")
}
