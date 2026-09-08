package checkout

import (
	"context"
	"strings"
)

// CreatePaymentAttempt delegates one Payment attempt command through the
// consumer-owned PaymentAttemptPort. Transport must not call Payment directly.
func (s *Service) CreatePaymentAttempt(
	ctx context.Context,
	command CreatePaymentAttemptCommand,
) (PaymentAttempt, error) {
	command.Method = strings.TrimSpace(command.Method)
	command.CommandKey = strings.TrimSpace(command.CommandKey)
	if s == nil || s.attempts == nil || !command.Subject.IsValid() ||
		command.OrderID == 0 || command.Method == "" || command.CommandKey == "" {
		return PaymentAttempt{}, ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return PaymentAttempt{}, err
	}
	return s.attempts.CreateAttempt(ctx, command)
}
