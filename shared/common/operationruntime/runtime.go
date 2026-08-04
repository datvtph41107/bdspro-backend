package operationruntime

import (
	_operationbinding "common/operationbinding"
	_operationcontext "common/operationcontext"
	_operationregistry "common/operationregistry"

	"context"
	"errors"
	"fmt"
)

var (
	ErrResolverNotConfigured = errors.New("operation resolver is not configured")
	ErrOperationUnavailable  = errors.New("operation is unavailable")
)

/** Resolver selects an operation from domain request semantics. */
type Resolver func(context.Context, interface{}) (_operationregistry.Code, error)

/** Runtime resolves approved transports into operation context. */
type Runtime struct {
	operations *_operationregistry.Registry
	bindings   *_operationbinding.Registry
	resolvers  map[_operationregistry.Code]Resolver
}

/** New creates one explicit operation runtime. */
func New(
	operations *_operationregistry.Registry,
	bindings *_operationbinding.Registry,
	resolvers map[_operationregistry.Code]Resolver,
) (*Runtime, error) {
	if operations == nil {
		return nil, errors.New("operation registry is required")
	}
	if bindings == nil {
		return nil, errors.New("operation binding registry is required")
	}

	copied := make(map[_operationregistry.Code]Resolver, len(resolvers))
	for code, resolver := range resolvers {
		if err := _operationregistry.ValidateCode(code); err != nil {
			return nil, fmt.Errorf("resolver code: %w", err)
		}
		if resolver == nil {
			return nil, fmt.Errorf("resolver %q is nil", code)
		}
		copied[code] = resolver
	}

	return &Runtime{
		operations: operations,
		bindings:   bindings,
		resolvers:  copied,
	}, nil
}

/** RequireReady validates startup coverage and resolver wiring. */
func (r *Runtime) RequireReady(transports []_operationbinding.Transport) error {
	if r == nil {
		return errors.New("operation runtime is nil")
	}
	if err := r.bindings.RequireResolvable(transports); err != nil {
		return err
	}
	for _, transport := range transports {
		binding, _ := r.bindings.LookupTransport(transport)
		if binding.Disposition == _operationbinding.DispositionResolver {
			if _, ok := r.resolvers[binding.ResolverCode]; !ok {
				return fmt.Errorf("%w: %s", ErrResolverNotConfigured, binding.ResolverCode)
			}
		}
	}
	return nil
}

/** Resolve selects one operation for a transport. */
func (r *Runtime) Resolve(
	ctx context.Context,
	req interface{},
	transport _operationbinding.Transport,
) (_operationcontext.Selection, bool, error) {
	if r == nil {
		return _operationcontext.Selection{}, false, errors.New("operation runtime is nil")
	}

	if existing, ok := _operationcontext.FromContext(ctx); ok {
		if _, err := r.operations.RequireActive(existing.Code); err != nil {
			return _operationcontext.Selection{}, false, fmt.Errorf("%w: %v", ErrOperationUnavailable, err)
		}
		return existing, true, nil
	}

	binding, ok := r.bindings.LookupTransport(transport)
	if !ok {
		return _operationcontext.Selection{}, false, nil
	}

	var (
		operationCode _operationregistry.Code
		source        _operationcontext.Source
	)

	switch binding.Disposition {
	case _operationbinding.DispositionDirect:
		operationCode = binding.OperationCode
		source = _operationcontext.SourceDirectBinding
	case _operationbinding.DispositionResolver:
		resolver, exists := r.resolvers[binding.ResolverCode]
		if !exists {
			return _operationcontext.Selection{}, false, fmt.Errorf("%w: %s", ErrResolverNotConfigured, binding.ResolverCode)
		}
		resolved, err := resolver(ctx, req)
		if err != nil {
			return _operationcontext.Selection{}, false, err
		}
		operationCode = resolved
		source = _operationcontext.SourceDomainResolver
	default:
		return _operationcontext.Selection{}, false, nil
	}

	definition, err := r.operations.RequireActive(operationCode)
	if err != nil {
		return _operationcontext.Selection{}, false, fmt.Errorf("%w: %v", ErrOperationUnavailable, err)
	}
	if definition.ModuleOwner != binding.ModuleOwner {
		return _operationcontext.Selection{}, false, fmt.Errorf(
			"operation owner %q does not match binding owner %q",
			definition.ModuleOwner,
			binding.ModuleOwner,
		)
	}

	selection, err := _operationcontext.New(operationCode, binding.Code, source)
	if err != nil {
		return _operationcontext.Selection{}, false, err
	}
	return selection, true, nil
}

/** Bind resolves and stores one operation in context. */
func (r *Runtime) Bind(
	ctx context.Context,
	req interface{},
	transport _operationbinding.Transport,
) (context.Context, bool, error) {
	selection, ok, err := r.Resolve(ctx, req, transport)
	if err != nil || !ok {
		return ctx, ok, err
	}
	bound, err := _operationcontext.Bind(ctx, selection)
	if err != nil {
		return ctx, false, err
	}
	return bound, true, nil
}
