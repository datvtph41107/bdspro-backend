package validator

import (
	"regexp"
	"strings"
	"time"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

// DirectorySourceValidator validates DirectorySource requests
type DirectorySourceValidator struct{}

// NewDirectorySourceValidator creates a new DirectorySourceValidator
func NewDirectorySourceValidator() *DirectorySourceValidator {
	return &DirectorySourceValidator{}
}

// ValidateCreateRequest validates CreateDirectorySourceRequestDTO
func (v *DirectorySourceValidator) ValidateCreateRequest(req *dto.DirectorySourceDTO) error {
	if err := v.validateBasicFields(req.Name, req.Code); err != nil {
		return err
	}

	// if err := v.validateOptionalFields(req.Description, req.Category, req.Type, req.Icon, req.Color); err != nil {
	// 	return err
	// }

	if err := v.validateAmountFields(req.ExpectedAmount, req.ActualAmount); err != nil {
		return err
	}

	if err := v.validateDateFields(req.StartDate); err != nil {
		return err
	}

	if err := v.validatePaymentFields(req.PaymentMethod, req.Frequency); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateRequest validates UpdateDirectorySourceRequestDTO
func (v *DirectorySourceValidator) ValidateUpdateRequest(req *domain.DirectorySource) error {
	if err := v.validateBasicFields(req.Name, req.Code); err != nil {
		return err
	}

	// if err := v.validateOptionalFields(req.Description, req.Category, req.Type, req.Icon, req.Color); err != nil {
	// 	return err
	// }

	if err := v.validateAmountFields(req.ExpectedAmount, req.ActualAmount); err != nil {
		return err
	}

	if err := v.validateDateFields(req.StartDate); err != nil {
		return err
	}

	if err := v.validatePaymentFields(req.PaymentMethod, req.Frequency); err != nil {
		return err
	}

	return nil
}

// validateBasicFields validates basic required fields
func (v *DirectorySourceValidator) validateBasicFields(name, code string) error {
	// Validate name
	if strings.TrimSpace(name) == "" {
		return &ValidationError{
			Field:   "name",
			Message: "Name is required",
		}
	}

	if len(name) > 255 {
		return &ValidationError{
			Field:   "name",
			Message: "Name must not exceed 255 characters",
		}
	}

	// Validate code
	if strings.TrimSpace(code) == "" {
		return &ValidationError{
			Field:   "code",
			Message: "Code is required",
		}
	}

	if len(code) > 50 {
		return &ValidationError{
			Field:   "code",
			Message: "Code must not exceed 50 characters",
		}
	}

	// Validate code format (alphanumeric and dash only)
	codeRegex := regexp.MustCompile(`^[A-Z0-9-]+$`)
	if !codeRegex.MatchString(code) {
		return &ValidationError{
			Field:   "code",
			Message: "Code must contain only uppercase letters, numbers, and dashes",
		}
	}

	return nil
}

// validateOptionalFields validates optional fields
func (v *DirectorySourceValidator) validateOptionalFields(description, category, typeStr, icon, color string) error {
	// Validate description
	if len(description) > 1000 {
		return &ValidationError{
			Field:   "description",
			Message: "Description must not exceed 1000 characters",
		}
	}

	// Validate category
	if category != "" && len(category) > 100 {
		return &ValidationError{
			Field:   "category",
			Message: "Category must not exceed 100 characters",
		}
	}

	// Validate type
	if typeStr != "" && len(typeStr) > 50 {
		return &ValidationError{
			Field:   "type",
			Message: "Type must not exceed 50 characters",
		}
	}

	// Validate color (hex color code)
	if color != "" {
		colorRegex := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
		if !colorRegex.MatchString(color) {
			return &ValidationError{
				Field:   "color",
				Message: "Color must be a valid hex color code (e.g., #FF0000)",
			}
		}
	}

	// Validate icon
	if len(icon) > 100 {
		return &ValidationError{
			Field:   "icon",
			Message: "Icon must not exceed 100 characters",
		}
	}

	return nil
}

// validateAmountFields validates amount fields
func (v *DirectorySourceValidator) validateAmountFields(expectedAmount, actualAmount int64) error {
	// Validate expected amount
	if expectedAmount < 0 {
		return &ValidationError{
			Field:   "expectedAmount",
			Message: "Expected amount must not be negative",
		}
	}

	// Validate actual amount
	if actualAmount < 0 {
		return &ValidationError{
			Field:   "actualAmount",
			Message: "Actual amount must not be negative",
		}
	}

	return nil
}

// validateDateFields validates date fields
func (v *DirectorySourceValidator) validateDateFields(startDate *time.Time) error {
	// Start date is optional, no validation needed for nil
	return nil
}

// validatePaymentFields validates payment related fields
func (v *DirectorySourceValidator) validatePaymentFields(paymentMethod, frequency string) error {
	// Validate payment method
	if paymentMethod != "" && len(paymentMethod) > 50 {
		return &ValidationError{
			Field:   "paymentMethod",
			Message: "Payment method must not exceed 50 characters",
		}
	}

	// Validate frequency
	if frequency != "" && len(frequency) > 50 {
		return &ValidationError{
			Field:   "frequency",
			Message: "Frequency must not exceed 50 characters",
		}
	}

	// Validate frequency values
	validFrequencies := []string{"daily", "weekly", "monthly", "quarterly", "yearly"}
	if frequency != "" {
		valid := false
		for _, validFreq := range validFrequencies {
			if frequency == validFreq {
				valid = true
				break
			}
		}
		if !valid {
			return &ValidationError{
				Field:   "frequency",
				Message: "Frequency must be one of: daily, weekly, monthly, quarterly, yearly",
			}
		}
	}

	return nil
}
