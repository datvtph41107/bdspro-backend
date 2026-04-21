package http_resp

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	apperr "bdspro/core/pkg/errors"
)

// DecodeJSON decodes JSON body with safety checks:
// - size limit (1MB)
// - syntax error
// - unknown fields
// - multiple JSON objects
func DecodeJSON(r *http.Request, v interface{}) error {
	body := io.LimitReader(r.Body, 1<<20)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return apperr.ErrInvalidJSON.WithField("_syntax", "json_syntax", syntaxErr.Error(), nil)
		}
		if strings.Contains(err.Error(), "unknown field") {
			return apperr.ErrInvalidArgument.WithField("_", "unknown_field", err.Error(), nil)
		}
		return apperr.ErrInvalidJSON.WithCause(err)
	}

	// Check for extra data after the first JSON object
	if dec.More() {
		return apperr.ErrInvalidJSON.WithMessage("multiple JSON objects in request body")
	}
	return nil
}
