package operationbinding

import (
	_operationregistry "common/operationregistry"

	"errors"
	"fmt"
	"sort"
)

var (
	ErrBindingNotFound  = errors.New("operation binding not found")
	ErrResolverRequired = errors.New("domain resolver is required")
	ErrNotResolvable    = errors.New("transport is not resolvable")
)

/** Catalog is an explicit binding-registry input. */
type Catalog struct {
	Version  _operationregistry.Version
	Bindings []Binding
}

/** Registry stores immutable transport mapping decisions. */
type Registry struct {
	version     _operationregistry.Version
	ordered     []Binding
	byCode      map[_operationregistry.Code]Binding
	byTransport map[string]Binding
	byModule    map[_operationregistry.ModuleCode][]Binding
}

/** New validates and builds an immutable binding registry. */
func New(catalog Catalog, operations *_operationregistry.Registry) (*Registry, error) {
	if err := _operationregistry.ValidateVersion(catalog.Version); err != nil {
		return nil, fmt.Errorf("binding catalog version: %w", err)
	}

	byCode := make(map[_operationregistry.Code]Binding, len(catalog.Bindings))
	byTransport := map[string]Binding{}
	for _, candidate := range catalog.Bindings {
		binding := cloneBinding(candidate)
		if err := binding.Validate(operations); err != nil {
			return nil, fmt.Errorf("binding %q: %w", binding.Code, err)
		}
		if _, exists := byCode[binding.Code]; exists {
			return nil, fmt.Errorf("duplicate binding code %q", binding.Code)
		}
		for _, transport := range binding.Transports {
			key := transport.CanonicalKey()
			if existing, exists := byTransport[key]; exists {
				return nil, fmt.Errorf("transport %q is already owned by binding %q", key, existing.Code)
			}
			byTransport[key] = binding
		}
		byCode[binding.Code] = binding
	}

	ordered := make([]Binding, 0, len(byCode))
	byModule := map[_operationregistry.ModuleCode][]Binding{}
	for _, binding := range byCode {
		ordered = append(ordered, binding)
		byModule[binding.ModuleOwner] = append(byModule[binding.ModuleOwner], binding)
	}
	sortBindings(ordered)
	for module := range byModule {
		sortBindings(byModule[module])
	}

	return &Registry{
		version:     catalog.Version,
		ordered:     ordered,
		byCode:      byCode,
		byTransport: byTransport,
		byModule:    byModule,
	}, nil
}

/** Version returns the binding-catalog version. */
func (r *Registry) Version() _operationregistry.Version {
	if r == nil {
		return ""
	}
	return r.version
}

/** LookupCode returns one reviewed binding record. */
func (r *Registry) LookupCode(code _operationregistry.Code) (Binding, bool) {
	if r == nil {
		return Binding{}, false
	}
	binding, ok := r.byCode[code]
	return cloneBinding(binding), ok
}

/** LookupTransport returns the decision for one transport. */
func (r *Registry) LookupTransport(transport Transport) (Binding, bool) {
	if r == nil {
		return Binding{}, false
	}
	binding, ok := r.byTransport[transport.CanonicalKey()]
	return cloneBinding(binding), ok
}

/** ResolveDirect resolves a transport mapped to one operation. */
func (r *Registry) ResolveDirect(transport Transport) (_operationregistry.Code, error) {
	binding, ok := r.LookupTransport(transport)
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrBindingNotFound, transport.CanonicalKey())
	}
	switch binding.Disposition {
	case DispositionDirect:
		return binding.OperationCode, nil
	case DispositionResolver:
		return "", fmt.Errorf("%w: %s", ErrResolverRequired, binding.ResolverCode)
	default:
		return "", fmt.Errorf("%w: %s is %s", ErrNotResolvable, binding.Code, binding.Disposition)
	}
}

/** RequireReviewed verifies every transport has an explicit decision. */
func (r *Registry) RequireReviewed(transports []Transport) error {
	for _, transport := range transports {
		if err := transport.Validate(); err != nil {
			return err
		}
		if _, ok := r.LookupTransport(transport); !ok {
			return fmt.Errorf("%w: %s", ErrBindingNotFound, transport.CanonicalKey())
		}
	}
	return nil
}

/** RequireResolvable verifies every transport has an active target. */
func (r *Registry) RequireResolvable(transports []Transport) error {
	for _, transport := range transports {
		binding, ok := r.LookupTransport(transport)
		if !ok {
			return fmt.Errorf("%w: %s", ErrBindingNotFound, transport.CanonicalKey())
		}
		if binding.Disposition != DispositionDirect && binding.Disposition != DispositionResolver {
			return fmt.Errorf("%w: %s is %s", ErrNotResolvable, binding.Code, binding.Disposition)
		}
	}
	return nil
}

/** Bindings returns a stable immutable snapshot. */
func (r *Registry) Bindings() []Binding {
	if r == nil {
		return nil
	}
	result := make([]Binding, len(r.ordered))
	for index, binding := range r.ordered {
		result[index] = cloneBinding(binding)
	}
	return result
}

/** ByModule returns module-owned mapping decisions. */
func (r *Registry) ByModule(module _operationregistry.ModuleCode) []Binding {
	if r == nil {
		return nil
	}
	bindings := r.byModule[module]
	result := make([]Binding, len(bindings))
	for index, binding := range bindings {
		result[index] = cloneBinding(binding)
	}
	return result
}

func sortBindings(bindings []Binding) {
	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].Code < bindings[j].Code
	})
}
