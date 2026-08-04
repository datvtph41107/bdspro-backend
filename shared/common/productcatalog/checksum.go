package productcatalog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

type commercialTerms struct {
	ProductCode  string
	PlanCode     string
	Version      string
	DisplayName  string
	SubjectScope SubjectScope
	Entitlements []Entitlement
	Prices       []PriceItem
}

/** TermsChecksum fingerprints commercial terms without lifecycle state. */
func TermsChecksum(plan PlanVersion) (string, error) {
	if err := validatePlan(plan); err != nil {
		return "", err
	}

	entitlements := append([]Entitlement(nil), plan.Entitlements...)
	prices := append([]PriceItem(nil), plan.Prices...)
	sort.Slice(entitlements, func(i, j int) bool {
		return entitlements[i].Code < entitlements[j].Code
	})
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Code < prices[j].Code
	})

	payload, err := json.Marshal(commercialTerms{
		ProductCode:  plan.ProductCode,
		PlanCode:     plan.PlanCode,
		Version:      plan.Version,
		DisplayName:  plan.DisplayName,
		SubjectScope: plan.SubjectScope,
		Entitlements: entitlements,
		Prices:       prices,
	})
	if err != nil {
		return "", fmt.Errorf("marshal commercial terms: %w", err)
	}

	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
