package postgres

import (
	"time"

	domain "payment/internal/domain/payment"
)

func orderToModel(order domain.Order) orderModel {
	return orderModel{
		ID:          order.ID,
		SubjectKind: string(order.Subject.Kind), SubjectID: order.Subject.ID,
		ProductCode: order.Terms.ProductCode, PlanCode: order.Terms.PlanCode,
		PlanVersionID: order.Terms.PlanVersionID, PlanVersion: order.Terms.PlanVersion,
		TierRank: order.Terms.TierRank, SubscriptionTermDays: order.Terms.SubscriptionTermDays,
		TermsChecksum: order.Terms.TermsChecksum,
		Currency:      string(order.Terms.Price.Currency), AmountMinor: order.Terms.Price.AmountMinor,
		CommandKey: order.CommandKey, CommandFingerprint: order.CommandFingerprint,
		Reference: order.Reference, Status: string(order.Status), ExpiresAt: order.ExpiresAt,
		FundsConfirmedAt: order.FundsConfirmedAt, CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
	}
}

func orderFromModel(row orderModel) domain.Order {
	return domain.Order{
		ID:      row.ID,
		Subject: domain.Subject{Kind: domain.SubjectKind(row.SubjectKind), ID: row.SubjectID},
		Terms: domain.CommercialTerms{
			ProductCode: row.ProductCode, PlanCode: row.PlanCode,
			PlanVersionID: row.PlanVersionID, PlanVersion: row.PlanVersion,
			TierRank: row.TierRank, SubscriptionTermDays: row.SubscriptionTermDays,
			TermsChecksum: row.TermsChecksum,
			Price:         domain.Money{Currency: domain.CurrencyCode(row.Currency), AmountMinor: row.AmountMinor},
		},
		CommandKey: row.CommandKey, CommandFingerprint: row.CommandFingerprint,
		Reference: row.Reference, Status: domain.OrderStatus(row.Status), ExpiresAt: row.ExpiresAt,
		FundsConfirmedAt: row.FundsConfirmedAt, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}

func settlementToModel(settlement domain.Settlement) settlementModel {
	var orderID *uint64
	if settlement.OrderID != 0 {
		id := settlement.OrderID
		orderID = &id
	}
	return settlementModel{
		ID: settlement.ID, OrderID: orderID,
		Provider: settlement.Provider, ProviderTransactionID: settlement.ProviderTransactionID,
		Reference: settlement.Reference,
		Currency:  string(settlement.Amount.Currency), AmountMinor: settlement.Amount.AmountMinor,
		OccurredAt: settlement.OccurredAt, EvidenceHash: settlement.EvidenceHash,
		Status: string(settlement.Status), ReviewReason: settlement.ReviewReason,
		CreatedAt: settlement.CreatedAt, UpdatedAt: settlement.CreatedAt,
	}
}

func settlementFromModel(row settlementModel) domain.Settlement {
	var orderID uint64
	if row.OrderID != nil {
		orderID = *row.OrderID
	}
	return domain.Settlement{
		ID: row.ID, OrderID: orderID,
		Provider: row.Provider, ProviderTransactionID: row.ProviderTransactionID,
		Reference:  row.Reference,
		Amount:     domain.Money{Currency: domain.CurrencyCode(row.Currency), AmountMinor: row.AmountMinor},
		OccurredAt: row.OccurredAt, EvidenceHash: row.EvidenceHash,
		Status: domain.SettlementStatus(row.Status), ReviewReason: row.ReviewReason,
		CreatedAt: row.CreatedAt,
	}
}

func fulfillmentFromModel(row fulfillmentModel) domain.Fulfillment {
	lease := timeOrZero(row.LeaseUntil)
	return domain.Fulfillment{
		ID: row.ID, OrderID: row.OrderID, Status: domain.FulfillmentStatus(row.Status),
		AttemptCount: row.AttemptCount, AvailableAt: row.AvailableAt,
		LockedBy: row.LockedBy, ClaimVersion: row.ClaimVersion, LeaseUntil: lease,
		LastError: row.LastError, CompletedAt: row.CompletedAt,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}

func timeOrZero(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}

func paymentAttemptToModel(attempt domain.PaymentAttempt) paymentAttemptModel {
	return paymentAttemptModel{
		ID: attempt.ID, OrderID: attempt.OrderID,
		CommandKey: attempt.CommandKey, CommandFingerprint: attempt.CommandFingerprint,
		Method: string(attempt.Method), Provider: attempt.Provider, ProviderReference: attempt.ProviderReference,
		Status: string(attempt.Status), NextActionKind: string(attempt.NextAction.Kind),
		RedirectURL: attempt.NextAction.RedirectURL, QRPayload: attempt.NextAction.QRPayload,
		ActionExpiresAt: attempt.NextAction.ExpiresAt, ExpiresAt: attempt.ExpiresAt,
		FailureCode: attempt.FailureCode, FailureDetail: attempt.FailureDetail,
		CreatedAt: attempt.CreatedAt, UpdatedAt: attempt.UpdatedAt,
	}
}

func paymentAttemptFromModel(row paymentAttemptModel) domain.PaymentAttempt {
	return domain.PaymentAttempt{
		ID: row.ID, OrderID: row.OrderID,
		CommandKey: row.CommandKey, CommandFingerprint: row.CommandFingerprint,
		Method: domain.PaymentMethod(row.Method), Provider: row.Provider, ProviderReference: row.ProviderReference,
		Status:     domain.AttemptStatus(row.Status),
		NextAction: domain.NextAction{Kind: domain.NextActionKind(row.NextActionKind), RedirectURL: row.RedirectURL, QRPayload: row.QRPayload, ExpiresAt: row.ActionExpiresAt},
		ExpiresAt:  row.ExpiresAt, FailureCode: row.FailureCode, FailureDetail: row.FailureDetail,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}
