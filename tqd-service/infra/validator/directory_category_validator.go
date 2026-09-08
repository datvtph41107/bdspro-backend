package validator

import (
	_errors "common/errors"
	"regexp"
	"strings"
	"tqd/internal/domain"
)

// DirectoryCategoryValidator validates DirectoryCategory requests
type DirectoryCategoryValidator struct{}

// NewDirectoryCategoryValidator creates a new DirectoryCategoryValidator
func NewDirectoryCategoryValidator() *DirectoryCategoryValidator {
	return &DirectoryCategoryValidator{}
}

// ValidateCreateRequest validates domain DirectoryCategory for create
func (v *DirectoryCategoryValidator) ValidateCreateRequest(req *domain.DirectoryCategory) error {
	if err := v.validateBasicFields(req.Name, req.Code); err != nil {
		return err
	}

	if err := v.validateOptionalFields(req.Description, req.Icon, req.Color); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateRequest validates domain DirectoryCategory for update
func (v *DirectoryCategoryValidator) ValidateUpdateRequest(req *domain.DirectoryCategory) error {
	if err := v.validateBasicFields(req.Name, req.Code); err != nil {
		return err
	}

	if err := v.validateOptionalFields(req.Description, req.Icon, req.Color); err != nil {
		return err
	}

	return nil
}

// validateBasicFields validates basic required fields
func (v *DirectoryCategoryValidator) validateBasicFields(name, code string) error {
	// Validate name
	if strings.TrimSpace(name) == "" {
		return _errors.ReturnError(400, "name is required")
	}

	if len(name) > 255 {
		return _errors.ReturnError(400, "Name must not exceed 255 characters")
	}

	// Validate code
	if strings.TrimSpace(code) == "" {
		return _errors.ReturnError(400, "code is required")
	}

	if len(code) > 50 {
		return _errors.ReturnError(400, "Code must not exceed 50 characters")
	}

	// Validate code format (alphanumeric and underscore only)
	codeRegex := regexp.MustCompile(`^[A-Z0-9_]+$`)
	if !codeRegex.MatchString(code) {
		return _errors.ReturnError(400, "Code must contain only uppercase letters, numbers, and underscores")
	}

	return nil
}

// validateOptionalFields validates optional fields
func (v *DirectoryCategoryValidator) validateOptionalFields(description, icon, color string) error {
	// Validate description
	if len(description) > 500 {
		return _errors.ReturnError(400, "Description must not exceed 500 characters")
	}

	// Validate color (hex color code)
	if color != "" {
		colorRegex := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
		if !colorRegex.MatchString(color) {
			return &DirectoryCategoryValidationError{
				Field:   "color",
				Message: "Color must be a valid hex color code (e.g., #FF0000)",
			}
		}
	}

	// Validate icon
	if len(icon) > 100 {
		return &DirectoryCategoryValidationError{
			Field:   "icon",
			Message: "Icon must not exceed 100 characters",
		}
	}

	return nil
}

// DirectoryCategoryValidationError represents a validation error for directory category
type DirectoryCategoryValidationError struct {
	Field   string
	Message string
}

func (e *DirectoryCategoryValidationError) Error() string {
	return e.Message
}
