package kernel

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicateComponent = errors.New("duplicate component")
	ErrMissingDependency  = errors.New("missing dependency")
	ErrCycleDependency    = errors.New("cycle dependency")
	ErrInvalidID          = errors.New("invalid component id")
)

type LifecycleError struct {
	Component ComponentID
	Op        string
	Err       error
}

func (e *LifecycleError) Error() string {
	return fmt.Sprintf("component=%s op=%s: %v", e.Component, e.Op, e.Err)
}

func (e *LifecycleError) Unwrap() error {
	return e.Err
}
