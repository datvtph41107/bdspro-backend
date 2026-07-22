package apperr

import (
	stdErrors "errors"
	"fmt"
)

func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}

	var ae *AppError
	if stdErrors.As(err, &ae) {
		return ae.WithOp(op)
	}

	return Internal(op, err)
}

func Wrapf(op, format string, args ...any) error {
	return Wrap(op, fmt.Errorf(format, args...))
}
