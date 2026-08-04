package operationregistry

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	maxCodeLength        = 128
	maxModuleLength      = 64
	maxDescriptionLength = 500
)

/** Code is a stable business operation identifier. */
type Code string

/** ModuleCode identifies the module that owns an operation. */
type ModuleCode string

/** Version is a strict major.minor.patch semantic version. */
type Version string

/** Status defines the operation lifecycle. */
type Status string

const (
	StatusActive     Status = "active"
	StatusDeprecated Status = "deprecated"
	StatusRetired    Status = "retired"
)

/** ChargeMode defines when metering may be committed. */
type ChargeMode string

const (
	ChargeModeUnmetered         ChargeMode = "unmetered"
	ChargeModeCommitOnSuccess   ChargeMode = "commit_on_success"
	ChargeModeReserveThenCommit ChargeMode = "reserve_then_commit"
	ChargeModeMeasureAfter      ChargeMode = "measure_after_success"
)

/** SuccessCondition defines the business success checkpoint. */
type SuccessCondition string

const (
	SuccessHandlerSucceeded SuccessCondition = "handler_succeeded"
	SuccessBusinessOutcome  SuccessCondition = "business_outcome_confirmed"
	SuccessJobAccepted      SuccessCondition = "job_accepted"
	SuccessJobCompleted     SuccessCondition = "job_completed"
)

/** FailurePolicy defines the metering outcome after failure. */
type FailurePolicy string

const (
	FailureNoCharge          FailurePolicy = "no_charge"
	FailureCancelReservation FailurePolicy = "cancel_reservation"
	FailureRetainCharge      FailurePolicy = "retain_charge"
	FailureManualReconcile   FailurePolicy = "manual_reconcile"
)

/** Metering describes operation-level measurement semantics. */
type Metering struct {
	MeterCode        Code
	ChargeMode       ChargeMode
	SuccessCondition SuccessCondition
	FailurePolicy    FailurePolicy
}

/** Definition is the stable operation contract. */
type Definition struct {
	Code            Code
	ModuleOwner     ModuleCode
	FeatureCode     Code
	Metering        Metering
	Status          Status
	SemanticVersion Version
	Replacement     Code
	Description     string
}

/** IsMetered reports whether the operation has a commercial meter. */
func (d Definition) IsMetered() bool {
	return d.Metering.ChargeMode != ChargeModeUnmetered
}

/** Validate checks one operation definition. */
func (d Definition) Validate() error {
	if err := validateDottedCode(string(d.Code), "operation code"); err != nil {
		return err
	}
	if err := validateModuleCode(string(d.ModuleOwner)); err != nil {
		return err
	}
	if err := validateDottedCode(string(d.FeatureCode), "feature code"); err != nil {
		return err
	}
	if err := validateVersion(string(d.SemanticVersion)); err != nil {
		return err
	}
	if err := validateStatus(d.Status); err != nil {
		return err
	}
	if err := validateDescription(d.Description); err != nil {
		return err
	}
	if err := validateMetering(d.Metering); err != nil {
		return err
	}
	if d.Replacement != "" {
		if d.Status == StatusActive {
			return fmt.Errorf("active operation %q cannot declare replacement", d.Code)
		}
		if err := validateDottedCode(string(d.Replacement), "replacement code"); err != nil {
			return err
		}
		if d.Replacement == d.Code {
			return fmt.Errorf("operation %q cannot replace itself", d.Code)
		}
	}
	return nil
}

func validateMetering(m Metering) error {
	switch m.ChargeMode {
	case ChargeModeUnmetered:
		if m.MeterCode != "" {
			return fmt.Errorf("unmetered operation cannot declare meter code %q", m.MeterCode)
		}
	case ChargeModeCommitOnSuccess, ChargeModeReserveThenCommit, ChargeModeMeasureAfter:
		if err := validateDottedCode(string(m.MeterCode), "meter code"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported charge mode %q", m.ChargeMode)
	}

	switch m.SuccessCondition {
	case SuccessHandlerSucceeded, SuccessBusinessOutcome, SuccessJobAccepted, SuccessJobCompleted:
	default:
		return fmt.Errorf("unsupported success condition %q", m.SuccessCondition)
	}

	switch m.FailurePolicy {
	case FailureNoCharge, FailureCancelReservation, FailureRetainCharge, FailureManualReconcile:
	default:
		return fmt.Errorf("unsupported failure policy %q", m.FailurePolicy)
	}
	return nil
}

func validateStatus(status Status) error {
	switch status {
	case StatusActive, StatusDeprecated, StatusRetired:
		return nil
	default:
		return fmt.Errorf("unsupported operation status %q", status)
	}
}

func validateDottedCode(raw, field string) error {
	if raw == "" {
		return fmt.Errorf("%s is required", field)
	}
	if len(raw) > maxCodeLength {
		return fmt.Errorf("%s exceeds %d bytes", field, maxCodeLength)
	}
	if strings.TrimSpace(raw) != raw {
		return fmt.Errorf("%s cannot contain surrounding whitespace", field)
	}
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return fmt.Errorf("%s %q must contain at least two segments", field, raw)
	}
	for _, part := range parts {
		if err := validateToken(part); err != nil {
			return fmt.Errorf("invalid %s %q: %w", field, raw, err)
		}
	}
	return nil
}

func validateModuleCode(raw string) error {
	if raw == "" {
		return fmt.Errorf("module owner is required")
	}
	if len(raw) > maxModuleLength {
		return fmt.Errorf("module owner exceeds %d bytes", maxModuleLength)
	}
	if strings.Contains(raw, ".") {
		return fmt.Errorf("module owner %q must be one segment", raw)
	}
	if err := validateToken(raw); err != nil {
		return fmt.Errorf("invalid module owner %q: %w", raw, err)
	}
	return nil
}

func validateToken(raw string) error {
	if raw == "" {
		return fmt.Errorf("empty segment")
	}
	for i, r := range raw {
		if i == 0 {
			if r < 'a' || r > 'z' {
				return fmt.Errorf("segment must start with lowercase letter")
			}
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			continue
		}
		return fmt.Errorf("unsupported character %q", r)
	}
	return nil
}

func validateVersion(raw string) error {
	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		return fmt.Errorf("semantic version %q must use major.minor.patch", raw)
	}
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("semantic version %q has empty segment", raw)
		}
		if len(part) > 1 && part[0] == '0' {
			return fmt.Errorf("semantic version %q has leading zero", raw)
		}
		if _, err := strconv.ParseUint(part, 10, 32); err != nil {
			return fmt.Errorf("semantic version %q is invalid", raw)
		}
	}
	return nil
}

func validateDescription(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("description is required")
	}
	if len(raw) > maxDescriptionLength {
		return fmt.Errorf("description exceeds %d bytes", maxDescriptionLength)
	}
	for _, r := range raw {
		if unicode.IsControl(r) {
			return fmt.Errorf("description contains control character")
		}
	}
	return nil
}

/** ValidateCode validates a stable dotted code. */
func ValidateCode(code Code) error {
	return validateDottedCode(string(code), "code")
}

/** ValidateModuleCode validates one module owner code. */
func ValidateModuleCode(module ModuleCode) error {
	return validateModuleCode(string(module))
}

/** ValidateVersion validates one semantic version. */
func ValidateVersion(version Version) error {
	return validateVersion(string(version))
}
