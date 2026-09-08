package recovery

import (
	"context"
	"strings"
	"time"

	payment "payment/internal/domain/payment"
)

type Clock func() time.Time

type Store interface {
	RedriveFulfillment(ctx context.Context, command payment.RedriveCommand, now time.Time) (changed bool, replay bool, err error)
}

type Service struct {
	store Store
	now   Clock
}

func NewService(store Store, now Clock) *Service { return &Service{store: store, now: now} }

func (s *Service) RedriveFulfillment(ctx context.Context, command payment.RedriveCommand) (changed bool, replay bool, err error) {
	if s == nil || s.store == nil || s.now == nil || command.FulfillmentID == 0 ||
		strings.TrimSpace(command.CommandKey) == "" || strings.TrimSpace(command.ActorID) == "" {
		return false, false, payment.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return false, false, err
	}
	return s.store.RedriveFulfillment(ctx, command, s.now().UTC())
}
