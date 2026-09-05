package plan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

type termsPayload struct {
	ProductCode  string
	PlanCode     string
	Version      string
	DisplayName  string
	SubjectScope SubjectScope
	// omitempty preserves the exact pre-term checksum payload for legacy rows.
	// New publication still requires a positive term through Snapshot validation.
	SubscriptionTermDays int32 `json:",omitempty"`
	Entitlements         []Entitlement
	Operations           []OperationBinding
	Prices               []PriceItem
}

/** TermsChecksum fingerprints commercial terms without lifecycle state. */
func TermsChecksum(plan PlanVersion) (string, error) {
	if err := validatePlan(plan); err != nil {
		return "", err
	}

	entitlements := append([]Entitlement(nil), plan.Entitlements...)
	operations := append([]OperationBinding(nil), plan.Operations...)
	prices := append([]PriceItem(nil), plan.Prices...)
	sort.Slice(entitlements, func(i, j int) bool {
		return entitlements[i].Code < entitlements[j].Code
	})
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Code < prices[j].Code
	})
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Code < operations[j].Code
	})

	payload, err := json.Marshal(termsPayload{
		ProductCode:          plan.ProductCode,
		PlanCode:             plan.PlanCode,
		Version:              plan.Version,
		DisplayName:          plan.DisplayName,
		SubjectScope:         plan.SubjectScope,
		SubscriptionTermDays: plan.SubscriptionTermDays,
		Entitlements:         entitlements,
		Operations:           operations,
		Prices:               prices,
	})
	if err != nil {
		return "", fmt.Errorf("marshal commercial terms: %w", err)
	}

	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
