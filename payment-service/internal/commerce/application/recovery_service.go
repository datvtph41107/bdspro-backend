package application

import (
	"context"
	"strings"

	"payment/internal/commerce/domain"
)

type RecoveryService struct {
	store RecoveryStore
	now   Clock
}

func NewRecoveryService(store RecoveryStore, now Clock) *RecoveryService {
	return &RecoveryService{store: store, now: now}
}

func (s *RecoveryService) RedriveFulfillment(ctx context.Context, command domain.RedriveCommand) (changed bool, replay bool, err error) {
	if s == nil || s.store == nil || s.now == nil || command.FulfillmentID == 0 ||
		strings.TrimSpace(command.CommandKey) == "" || strings.TrimSpace(command.ActorID) == "" {
		return false, false, domain.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return false, false, err
	}
	return s.store.RedriveFulfillment(ctx, command, s.now().UTC())
}
