package seo_strategy

import "fmt"

// ValidateQueryEvidenceStatusTransition khóa lifecycle governance của evidence.
// active -> superseded|rejected
// superseded -> rejected
// rejected là terminal.
func ValidateQueryEvidenceStatusTransition(
	current QueryEvidenceStatus,
	target QueryEvidenceStatus,
) error {
	if !IsValidQueryEvidenceStatus(current) {
		return fmt.Errorf("query evidence current status không hợp lệ: %s", current)
	}
	if !IsValidQueryEvidenceStatus(target) {
		return fmt.Errorf("query evidence target status không hợp lệ: %s", target)
	}

	switch current {
	case QueryEvidenceStatusActive:
		switch target {
		case QueryEvidenceStatusSuperseded, QueryEvidenceStatusRejected:
			return nil
		}
	case QueryEvidenceStatusSuperseded:
		if target == QueryEvidenceStatusRejected {
			return nil
		}
	case QueryEvidenceStatusRejected:
		// Terminal state.
	}

	return fmt.Errorf(
		"không thể chuyển query evidence từ %s sang %s",
		current,
		target,
	)
}
