package checkout

import (
	"context"
	"strings"
)

// Service điều phối checkout nhưng không tự quyết định số tiền hay quyền
// sử dụng. Giá và điều kiện gói lấy từ PlanVersion của User Service; Payment
// Service chỉ nhận order đã đóng băng thông tin này để tạo và đối soát thanh
// toán. Nhờ vậy retry không bị đổi giá khi catalog được publish phiên bản mới.
type Service struct {
	plans          PlanResolver
	subscriptions  SubscriptionReader
	commands       CommandStore
	payment        PaymentOrderPort
	attempts       PaymentAttemptPort
	paymentMethods PaymentMethodReader
	contacts       ContactStore
	now            Clock
}

func NewService(
	plans PlanResolver,
	subscriptions SubscriptionReader,
	commands CommandStore,
	payment PaymentOrderPort,
	attempts PaymentAttemptPort,
	paymentMethods PaymentMethodReader,
	contacts ContactStore,
	now Clock,
) *Service {
	return &Service{
		plans:          plans,
		subscriptions:  subscriptions,
		commands:       commands,
		payment:        payment,
		attempts:       attempts,
		paymentMethods: paymentMethods,
		contacts:       contacts,
		now:            now,
	}
}

func (s *Service) ListPaymentMethods(ctx context.Context) (PaymentMethods, error) {
	if s == nil || s.paymentMethods == nil {
		return PaymentMethods{}, ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return PaymentMethods{}, err
	}
	return s.paymentMethods.ListMethods(ctx)
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

	// Kiểm tra replay bền vững trước khi đọc trạng thái hiện tại. Một retry
	// mạng phải trả lại command cũ, không được tạo order thứ hai.
	// Terms đã đóng băng vẫn được giữ khi catalog publish version mới hoặc
	// thanh toán đầu tiên đã thay đổi subscription.
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

	current, found, err := s.subscriptions.CurrentForProduct(ctx, subject, terms.ProductID, at)
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

	// CreateCommand là điểm phân xử idempotency. Trong hai request đồng thời,
	// command thắng cùng bộ terms đã đóng băng là sự thật cho mọi caller có
	// cùng Subject và CommandKey.
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
