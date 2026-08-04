package operationregistry

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrOperationNotFound  = errors.New("operation not found")
	ErrOperationNotActive = errors.New("operation is not active")
)

/** Catalog is an explicit registry input. */
type Catalog struct {
	Version     Version
	Definitions []Definition
}

/** Registry is an immutable operation catalog. */
type Registry struct {
	version  Version
	byCode   map[Code]Definition
	ordered  []Definition
	byModule map[ModuleCode][]Definition
}

/** New validates and builds an immutable registry. */
func New(catalog Catalog) (*Registry, error) {
	if err := validateVersion(string(catalog.Version)); err != nil {
		return nil, fmt.Errorf("catalog version: %w", err)
	}

	byCode := make(map[Code]Definition, len(catalog.Definitions))
	for _, definition := range catalog.Definitions {
		if err := definition.Validate(); err != nil {
			return nil, fmt.Errorf("operation %q: %w", definition.Code, err)
		}
		if _, exists := byCode[definition.Code]; exists {
			return nil, fmt.Errorf("duplicate operation code %q", definition.Code)
		}
		byCode[definition.Code] = definition
	}

	if err := validateReplacements(byCode); err != nil {
		return nil, err
	}

	ordered := make([]Definition, 0, len(byCode))
	byModule := make(map[ModuleCode][]Definition)
	for _, definition := range byCode {
		ordered = append(ordered, definition)
		byModule[definition.ModuleOwner] = append(byModule[definition.ModuleOwner], definition)
	}
	sortDefinitions(ordered)
	for module := range byModule {
		sortDefinitions(byModule[module])
	}

	return &Registry{
		version:  catalog.Version,
		byCode:   byCode,
		ordered:  ordered,
		byModule: byModule,
	}, nil
}

/** Version returns the catalog semantic version. */
func (r *Registry) Version() Version {
	if r == nil {
		return ""
	}
	return r.version
}

/** Lookup returns one operation definition. */
func (r *Registry) Lookup(code Code) (Definition, bool) {
	if r == nil {
		return Definition{}, false
	}
	definition, ok := r.byCode[code]
	return definition, ok
}

/** RequireActive returns only an active operation. */
func (r *Registry) RequireActive(code Code) (Definition, error) {
	definition, ok := r.Lookup(code)
	if !ok {
		return Definition{}, fmt.Errorf("%w: %s", ErrOperationNotFound, code)
	}
	if definition.Status != StatusActive {
		return Definition{}, fmt.Errorf("%w: %s is %s", ErrOperationNotActive, code, definition.Status)
	}
	return definition, nil
}

/** Definitions returns a stable sorted snapshot. */
func (r *Registry) Definitions() []Definition {
	if r == nil {
		return nil
	}
	return append([]Definition(nil), r.ordered...)
}

/** ByModule returns a stable module-owned snapshot. */
func (r *Registry) ByModule(module ModuleCode) []Definition {
	if r == nil {
		return nil
	}
	return append([]Definition(nil), r.byModule[module]...)
}

func validateReplacements(byCode map[Code]Definition) error {
	for code, definition := range byCode {
		if definition.Replacement == "" {
			continue
		}
		target, ok := byCode[definition.Replacement]
		if !ok {
			return fmt.Errorf("operation %q replacement %q is not registered", code, definition.Replacement)
		}
		if target.Status == StatusRetired {
			return fmt.Errorf("operation %q replacement %q is retired", code, definition.Replacement)
		}
	}

	for code := range byCode {
		visited := map[Code]bool{}
		current := code
		for current != "" {
			if visited[current] {
				return fmt.Errorf("replacement cycle detected at %q", current)
			}
			visited[current] = true
			current = byCode[current].Replacement
		}
	}
	return nil
}

func sortDefinitions(definitions []Definition) {
	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].Code < definitions[j].Code
	})
}
