package operationbinding

import (
	_operationregistry "common/operationregistry"

	"fmt"
	"strings"
	"unicode"
)

/** Disposition describes one reviewed mapping decision. */
type Disposition string

const (
	DispositionDirect            Disposition = "direct"
	DispositionResolver          Disposition = "resolver"
	DispositionInternalComponent Disposition = "internal_component"
	DispositionReviewRequired    Disposition = "review_required"
	DispositionSourceMissing     Disposition = "source_missing"
	DispositionRetiredTransport  Disposition = "retired_transport"
)

/** Binding is one explicit transport mapping decision. */
type Binding struct {
	Code          _operationregistry.Code
	ModuleOwner   _operationregistry.ModuleCode
	Transports    []Transport
	Disposition   Disposition
	OperationCode _operationregistry.Code
	ResolverCode  _operationregistry.Code
	Description   string
}

/** Validate checks one binding against the operation registry. */
func (b Binding) Validate(operations *_operationregistry.Registry) error {
	if err := _operationregistry.ValidateCode(b.Code); err != nil {
		return fmt.Errorf("binding code: %w", err)
	}
	if err := _operationregistry.ValidateModuleCode(b.ModuleOwner); err != nil {
		return err
	}
	if err := validateDescription(b.Description); err != nil {
		return err
	}

	if b.Disposition == DispositionSourceMissing {
		if len(b.Transports) != 0 {
			return fmt.Errorf("source-missing binding cannot declare transport")
		}
	} else if len(b.Transports) == 0 {
		return fmt.Errorf("binding %q requires at least one transport", b.Code)
	}

	seen := map[string]struct{}{}
	for _, transport := range b.Transports {
		if err := transport.Validate(); err != nil {
			return fmt.Errorf("binding %q: %w", b.Code, err)
		}
		key := transport.CanonicalKey()
		if _, exists := seen[key]; exists {
			return fmt.Errorf("binding %q repeats transport %q", b.Code, key)
		}
		seen[key] = struct{}{}
	}

	switch b.Disposition {
	case DispositionDirect:
		if b.ResolverCode != "" {
			return fmt.Errorf("direct binding cannot declare resolver")
		}
		if b.OperationCode == "" {
			return fmt.Errorf("direct binding requires operation code")
		}
		if operations == nil {
			return fmt.Errorf("direct binding requires operation registry")
		}
		operation, ok := operations.Lookup(b.OperationCode)
		if !ok {
			return fmt.Errorf("operation %q is not registered", b.OperationCode)
		}
		if operation.Status == _operationregistry.StatusRetired {
			return fmt.Errorf("operation %q is retired", b.OperationCode)
		}
		if operation.ModuleOwner != b.ModuleOwner {
			return fmt.Errorf("binding owner %q does not match operation owner %q", b.ModuleOwner, operation.ModuleOwner)
		}
	case DispositionResolver:
		if b.OperationCode != "" {
			return fmt.Errorf("resolver binding cannot declare direct operation")
		}
		if err := _operationregistry.ValidateCode(b.ResolverCode); err != nil {
			return fmt.Errorf("resolver code: %w", err)
		}
	case DispositionInternalComponent, DispositionReviewRequired, DispositionRetiredTransport:
		if b.OperationCode != "" || b.ResolverCode != "" {
			return fmt.Errorf("%s binding cannot declare target", b.Disposition)
		}
	case DispositionSourceMissing:
		if b.OperationCode != "" || b.ResolverCode != "" {
			return fmt.Errorf("source-missing binding cannot declare target")
		}
	default:
		return fmt.Errorf("unsupported binding disposition %q", b.Disposition)
	}
	return nil
}

func validateDescription(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("binding description is required")
	}
	if len(value) > 500 {
		return fmt.Errorf("binding description exceeds 500 bytes")
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return fmt.Errorf("binding description contains control character")
		}
	}
	return nil
}

func cloneBinding(binding Binding) Binding {
	binding.Transports = append([]Transport(nil), binding.Transports...)
	return binding
}
