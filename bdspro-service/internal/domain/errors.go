package domain

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrDuplicate        = errors.New("duplicate")
	ErrConcurrentUpdate = errors.New("ErrConcurrentUpdate")
)
