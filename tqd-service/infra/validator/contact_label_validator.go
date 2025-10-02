package validator

import (
	"regexp"
	"strings"
	"tqd/internal/dto"
)

// ContactLabelValidator validates ContactLabel requests
type ContactLabelValidator struct{}

// NewContactLabelValidator creates a new ContactLabelValidator
func NewContactLabelValidator() *ContactLabelValidator {
	return &ContactLabelValidator{}
}

// ValidateCreateRequest validates CreateContactLabelRequestDTO
func (v *ContactLabelValidator) ValidateCreateRequest(req *dto.CreateContactLabelRequestDTO) error {
	if err := v.validateBasicFields(req.Name, req.Code); err != nil {
		return err
	}

	if err := v.validateOptionalFields(req.Description, req.Color, req.Icon); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateRequest validates UpdateContactLabelRequestDTO
func (v *ContactLabelValidator) ValidateUpdateRequest(req *dto.UpdateContactLabelRequestDTO) error {
	if err := v.validateBasicFields(req.Name, req.Code); err != nil {
		return err
	}

	if err := v.validateOptionalFields(req.Description, req.Color, req.Icon); err != nil {
		return err
	}

	return nil
}

// validateBasicFields validates basic required fields
func (v *ContactLabelValidator) validateBasicFields(name, code string) error {
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

	// Validate code format (alphanumeric and underscore only)
	codeRegex := regexp.MustCompile(`^[A-Z0-9_]+$`)
	if !codeRegex.MatchString(code) {
		return &ValidationError{
			Field:   "code",
			Message: "Code must contain only uppercase letters, numbers, and underscores",
		}
	}

	return nil
}

// validateOptionalFields validates optional fields
func (v *ContactLabelValidator) validateOptionalFields(description, color, icon string) error {
	// Validate description
	if len(description) > 500 {
		return &ValidationError{
			Field:   "description",
			Message: "Description must not exceed 500 characters",
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

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
