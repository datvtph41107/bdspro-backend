// Package operation defines the server-owned semantic action selected for a request.
package operation

import (
	"common/codeformat"
	"errors"
)

const MaxCodeLength = codeformat.MaxLength

var ErrInvalidCode = errors.New("operation code is invalid")

// Code is the stable semantic identifier of an action.
// It is independent of transport names, Operation-ID, and commercial policy.
type Code string

// Parse validates the serialized form of an operation code.
func Parse(raw string) (Code, error) {
	code := Code(raw)
	if !code.IsValid() {
		return "", ErrInvalidCode
	}

	return code, nil
}

// IsValid reports whether the code uses the stable dotted identifier syntax.
func (c Code) IsValid() bool {
	return codeformat.IsValid(string(c))
}
