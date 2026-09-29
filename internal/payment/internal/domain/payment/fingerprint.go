package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func OrderFingerprint(subject Subject, terms CommercialTerms) string {
	payload := strings.Join([]string{
		string(subject.Kind), subject.ID, terms.ProductCode, terms.PlanCode,
		fmt.Sprintf("%d", terms.PlanVersionID), terms.PlanVersion,
		fmt.Sprintf("%d", terms.TierRank), fmt.Sprintf("%d", terms.SubscriptionTermDays),
		terms.TermsChecksum, string(terms.Price.Currency), fmt.Sprintf("%d", terms.Price.AmountMinor),
	}, "\x00")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func ProviderEvidenceHash(event ProviderEvent) string {
	if value := strings.TrimSpace(event.EvidenceHash); value != "" {
		if len(value) == sha256.Size*2 {
			if _, err := hex.DecodeString(value); err == nil {
				return strings.ToLower(value)
			}
		}
		sum := sha256.Sum256([]byte(value))
		return hex.EncodeToString(sum[:])
	}
	payload := strings.Join([]string{
		event.Provider, event.TransactionID, event.Reference, string(event.Amount.Currency),
		fmt.Sprintf("%d", event.Amount.AmountMinor), event.OccurredAt.UTC().Format(time.RFC3339Nano), string(event.Type),
	}, "\x00")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func PaymentAttemptFingerprint(orderID uint64, method PaymentMethod, provider string) string {
	payload := strings.Join([]string{fmt.Sprintf("%d", orderID), string(method), strings.TrimSpace(provider)}, "\x00")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
