package grpcoperation

import (
	"common/businessoperation"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidMethod    = errors.New("gRPC method is invalid")
	ErrDuplicateMethod  = errors.New("gRPC method already has an operation")
	ErrInvalidOperation = errors.New("business operation is invalid")
)

/**
 * Registry giữ mapping giữa gRPC method và business operation.
 *
 * Mapping nằm ở TQD vì TQD hiểu method nào đại diện cho business action nào.
 */
type Registry struct {
	byMethod map[string]businessoperation.Code
}

/**
 * NewRegistry tạo registry rỗng.
 */
func NewRegistry() *Registry {
	return &Registry{byMethod: make(map[string]businessoperation.Code)}
}

/**
 * Add thêm một mapping gRPC method -> business operation.
 */
func (r *Registry) Add(method string, operation businessoperation.Code) error {
	if r == nil {
		return errors.New("registry is nil")
	}
	if !validMethod(method) {
		return fmt.Errorf("%w: %q", ErrInvalidMethod, method)
	}
	if !operation.IsValid() {
		return fmt.Errorf("%w: %q", ErrInvalidOperation, operation)
	}
	if _, exists := r.byMethod[method]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateMethod, method)
	}

	r.byMethod[method] = operation
	return nil
}

/**
 * Find trả về operation đã được khai báo cho gRPC method.
 */
func (r *Registry) Find(method string) (businessoperation.Code, bool) {
	if r == nil {
		return "", false
	}
	operation, ok := r.byMethod[method]
	return operation, ok
}

func validMethod(method string) bool {
	if method == "" || strings.TrimSpace(method) != method || method[0] != '/' {
		return false
	}
	parts := strings.Split(strings.TrimPrefix(method, "/"), "/")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}
