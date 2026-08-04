package operationcontext

import (
	_operationregistry "common/operationregistry"

	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	CodeMetadataKey    = "x-qhpro-operation-code"
	BindingMetadataKey = "x-qhpro-operation-binding"
	SourceMetadataKey  = "x-qhpro-operation-source"
)

/** Source identifies how the operation was selected. */
type Source string

const (
	SourceDirectBinding  Source = "direct_binding"
	SourceDomainResolver Source = "domain_resolver"
)

var ErrOperationAlreadySelected = errors.New("operation is already selected")

/** Selection is one server-resolved business operation. */
type Selection struct {
	Code        _operationregistry.Code
	BindingCode _operationregistry.Code
	Source      Source
}

type contextKey struct{}

/** New validates one operation selection. */
func New(code, bindingCode _operationregistry.Code, source Source) (Selection, error) {
	if err := _operationregistry.ValidateCode(code); err != nil {
		return Selection{}, fmt.Errorf("operation code: %w", err)
	}
	if err := _operationregistry.ValidateCode(bindingCode); err != nil {
		return Selection{}, fmt.Errorf("binding code: %w", err)
	}
	if source != SourceDirectBinding && source != SourceDomainResolver {
		return Selection{}, fmt.Errorf("unsupported operation source %q", source)
	}
	return Selection{Code: code, BindingCode: bindingCode, Source: source}, nil
}

/** Bind stores one immutable operation selection. */
func Bind(ctx context.Context, selection Selection) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	validated, err := New(selection.Code, selection.BindingCode, selection.Source)
	if err != nil {
		return ctx, err
	}
	if existing, ok := FromContext(ctx); ok {
		if existing == validated {
			return ctx, nil
		}
		return ctx, ErrOperationAlreadySelected
	}
	return context.WithValue(ctx, contextKey{}, validated), nil
}

/** FromContext returns the selected operation. */
func FromContext(ctx context.Context) (Selection, bool) {
	if ctx == nil {
		return Selection{}, false
	}
	selection, ok := ctx.Value(contextKey{}).(Selection)
	if !ok {
		return Selection{}, false
	}
	if _, err := New(selection.Code, selection.BindingCode, selection.Source); err != nil {
		return Selection{}, false
	}
	return selection, true
}

/** AppendMetadata serializes one selected operation. */
func AppendMetadata(md metadata.MD, selection Selection) {
	if md == nil {
		return
	}
	validated, err := New(selection.Code, selection.BindingCode, selection.Source)
	if err != nil {
		return
	}
	md.Set(CodeMetadataKey, string(validated.Code))
	md.Set(BindingMetadataKey, string(validated.BindingCode))
	md.Set(SourceMetadataKey, string(validated.Source))
}

/** FromMetadata validates one transported operation selection. */
func FromMetadata(md metadata.MD) (Selection, bool, error) {
	code, hasCode, err := single(md, CodeMetadataKey)
	if err != nil {
		return Selection{}, false, err
	}
	binding, hasBinding, err := single(md, BindingMetadataKey)
	if err != nil {
		return Selection{}, false, err
	}
	source, hasSource, err := single(md, SourceMetadataKey)
	if err != nil {
		return Selection{}, false, err
	}
	if !hasCode && !hasBinding && !hasSource {
		return Selection{}, false, nil
	}
	if !hasCode || !hasBinding || !hasSource {
		return Selection{}, false, errors.New("incomplete operation metadata")
	}
	selection, err := New(
		_operationregistry.Code(code),
		_operationregistry.Code(binding),
		Source(source),
	)
	if err != nil {
		return Selection{}, false, err
	}
	return selection, true, nil
}

func single(md metadata.MD, key string) (string, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) != 1 || values[0] == "" || strings.TrimSpace(values[0]) != values[0] {
		return "", false, fmt.Errorf("invalid %s metadata", key)
	}
	return values[0], true, nil
}
