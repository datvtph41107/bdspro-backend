package seo_strategy

import "testing"

func TestValidateQueryEvidenceStatusTransition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		current QueryEvidenceStatus
		target  QueryEvidenceStatus
		allowed bool
	}{
		{"active_to_superseded", QueryEvidenceStatusActive, QueryEvidenceStatusSuperseded, true},
		{"active_to_rejected", QueryEvidenceStatusActive, QueryEvidenceStatusRejected, true},
		{"superseded_to_rejected", QueryEvidenceStatusSuperseded, QueryEvidenceStatusRejected, true},
		{"active_to_active", QueryEvidenceStatusActive, QueryEvidenceStatusActive, false},
		{"superseded_to_active", QueryEvidenceStatusSuperseded, QueryEvidenceStatusActive, false},
		{"superseded_to_superseded", QueryEvidenceStatusSuperseded, QueryEvidenceStatusSuperseded, false},
		{"rejected_to_active", QueryEvidenceStatusRejected, QueryEvidenceStatusActive, false},
		{"rejected_to_superseded", QueryEvidenceStatusRejected, QueryEvidenceStatusSuperseded, false},
		{"rejected_to_rejected", QueryEvidenceStatusRejected, QueryEvidenceStatusRejected, false},
		{"invalid_current", QueryEvidenceStatus("invalid"), QueryEvidenceStatusRejected, false},
		{"invalid_target", QueryEvidenceStatusActive, QueryEvidenceStatus("invalid"), false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateQueryEvidenceStatusTransition(tt.current, tt.target)
			if tt.allowed && err != nil {
				t.Fatalf("transition phải hợp lệ: %v", err)
			}
			if !tt.allowed && err == nil {
				t.Fatal("transition phải bị từ chối")
			}
		})
	}
}
